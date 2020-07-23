#include "JIT.h"
#include "BuildContext.h"
#include <llvm/IRReader/IRReader.h>
#include <llvm/Support/TargetRegistry.h>
#include <llvm/Support/TargetSelect.h>
#include <llvm/IR/LegacyPassManager.h>
#include <llvm/Object/ObjectFile.h>
#include <iostream>
#include <csignal>
#include <unistd.h>
#include <sys/time.h>

using namespace llvm;


static inline bool startswith(const std::string &str, const char *substr) {
	auto len = strlen(substr);
	return str.size() >= len && strncmp(str.c_str(), substr, len) == 0;
}


void printLLVMIR(Module *m) {
	outs() << "Module size = " << m->size() << "\n";
	outs() << "\n-------------------------------------------------------------\n";
	outs() << *m << "\n";
	outs() << "-------------------------------------------------------------\n\n";
}

JitEngine::JitEngine() {
	InitializeNativeTarget();
	InitializeNativeTargetAsmParser();
	InitializeNativeTargetAsmPrinter();

	std::string error;
	std::string targetTriple = sys::getDefaultTargetTriple();
	auto target = TargetRegistry::lookupTarget(targetTriple, error);
	TargetOptions opt;
	opt.UseInitArray = 1;
	auto targetMachine = target->createTargetMachine(targetTriple, "generic", "", opt, Reloc::PIC_);
	targetMachine->setOptLevel(CodeGenOpt::Aggressive);

	// Engine must be initialized with some module - so we give it an empty one
	auto emptyModule = llvm::make_unique<Module>("__empty", BuildContext::ctx);
	emptyModule->setTargetTriple(targetTriple);
	EngineBuilder builder(std::move(emptyModule));

	// outs() << "Target: " << targetMachine->getTargetTriple().str() << "\n";
	engine = builder.create(targetMachine);
	engine->DisableLazyCompilation(true);
	engine->DisableSymbolSearching(false);
}

void JitEngine::addModule(const std::string &object) {
	buffers.push_back(MemoryBuffer::getMemBuffer(object, "", false));
	auto objectFile = object::ObjectFile::createObjectFile(buffers.back()->getMemBufferRef());
	if (!objectFile) {
		llvm::logAllUnhandledErrors(objectFile.takeError(), errs(), "Jit loading failed:\n");
	}
	// Find initializer functions in this file
	for (auto &s: objectFile.get()->symbols()) {
		auto name = s.getName().get().str();
		if (startswith(name, "__saarlang_init_")) {
			initFunctions.push_back(name);
		}
	}
	// Add object file and all of its symbols into the execution engine
	engine->addObjectFile(std::move(objectFile.get()));
}

void JitEngine::init(bool stdlib) {
	if (stdlib)
		importSaarlangLibrary(*this);
	engine->finalizeObject();

	// call __saarlang_init_*()
	engine->runStaticConstructorsDestructors(false);
	for (const auto &initname: initFunctions) {
		auto init = (void (*)()) engine->getFunctionAddress(initname);
		init();
	}
}


sl_int JitEngine::execute() {
	// call main()
	auto main = (int64_t (*)()) engine->getFunctionAddress("main");
	if (!main) {
		std::cerr << "No main function defined!" << std::endl;
		throw std::exception();
	}
	std::cout << "--- Saarlang execution starts ---\n";
	auto result = main();
	return result;
}

sl_int JitEngine::executeIsolated() {
	// Timeout - 1,2 seconds
	itimerval interval = {.it_interval = {0, 0}, .it_value = {1, 200000}};
	setitimer(ITIMER_REAL, &interval, nullptr);
	// Prevent DoS (limit memory usage, cpu time, ...)
	auto cmd = "prlimit --cpu=10 --data=48000000:48000000 --stack=10000000:10000000 --nproc=1024:1024 --pid " + std::to_string(getpid());
	system(cmd.c_str());
	return execute();
}

void JitEngine::addFunction(const std::string &name, void *functionAddr) {
	// New symbols overwrite old symbols - call this after all objects have been loaded
	engine->addGlobalMapping(name, (uint64_t) functionAddr);
}

void JitEngine::addSymbol(const std::string &name, void *addr) {
	// New symbols overwrite old symbols - call this after all objects have been loaded
	engine->addGlobalMapping(name, (uint64_t) addr);
}



static void llvm_error_handler(void *user_data, const std::string &reason, bool gen_crash_diag) {
	std::cerr << "LLVM ERROR: \"" << reason << "\"" << std::endl;
	std::cerr << "gen_crash_diag = " << gen_crash_diag << std::endl;
	throw std::runtime_error(reason);
}

__attribute__((constructor))
static void init_llvm_error_handler() {
	llvm::install_fatal_error_handler(llvm_error_handler, nullptr);
}
