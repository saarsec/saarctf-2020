#include "SaarlangModule.h"
#include "Parser.h"
#include "BuildContext.h"
#include <fstream>

#include <llvm/Support/TargetRegistry.h>
#include <llvm/Support/TargetSelect.h>
#include <llvm/Target/TargetMachine.h>
#include <llvm/IR/LegacyPassManager.h>
#include <llvm/Support/FileSystem.h>
#include <llvm/IR/Verifier.h>
#include <llvm/CodeGen/DIE.h>
#include <llvm/CodeGen/UnreachableBlockElim.h>
#include <llvm/Transforms/Scalar.h>

using namespace llvm;

llvm::Module *SaarlangModule::generateCode() {
	BuildContext context(filename);
	for (auto &def: definitions) {
		def->generateCode(context);
	}
	context.finalizeModule();

	legacy::PassManager manager;
	manager.add(createVerifierPass(true));
	manager.add(createDeadInstEliminationPass());
	manager.add(createCFGSimplificationPass());
	manager.run(*context.module);
	return context.module;
}



void SaarlangModule::resolveImports(Diagnostic &diag, SimpleModuleLoader &loader) {
	for (auto &node: imports) {
		SaarlangModule *m = loader.load(node->filename);
		for (auto &def: m->definitions) {
			def->declare(diag, globalSymbols);
		}
	}
}


std::string &SaarlangModule::prepareJITObject() {
	if (!llvmModule)
		llvmModule = generateCode();

	if (!asmModule.empty())
		return asmModule;

	{
		InitializeNativeTarget();
		InitializeNativeTargetAsmParser();
		InitializeNativeTargetAsmPrinter();

		raw_string_ostream stream(asmModule);
		buffer_ostream stream2(stream);
		stream2.SetUnbuffered();

		// Setup LLVM backend
		std::string error;
		std::string targetTriple = sys::getDefaultTargetTriple();
		auto target = TargetRegistry::lookupTarget(targetTriple, error);
		TargetOptions opt;
		opt.UseInitArray = 1;
		auto targetMachine = target->createTargetMachine(targetTriple, "generic", "", opt, Reloc::PIC_);
		targetMachine->setOptLevel(CodeGenOpt::Aggressive);

		// Generate code
		legacy::PassManager pass;
		pass.add(createVerifierPass(true));
		targetMachine->addPassesToEmitFile(pass, stream2, nullptr, TargetMachine::CGFT_ObjectFile);
		pass.run(*llvmModule);
		stream.flush();
	}

	return asmModule;
}

SaarlangModule::~SaarlangModule() {
	for (auto &it: imports) delete it;
	for (auto &it: definitions) delete it;
	if (llvmModule) delete llvmModule;
}


SaarlangModule *SimpleModuleLoader::load(const std::string &filename) {
	auto it = modules.find(filename);
	if (it != modules.end())
		return &it->second;

	std::ifstream fileinput(basepath + filename, std::ios::in);
	if (!fileinput.is_open())
		diag.file_error(filename);
	return preload(filename, fileinput);
}

SaarlangModule *SimpleModuleLoader::preload(const std::string &filename, std::istream &input) {
	diag.setFilename(filename);
	Lexer lexer(input, diag);
	lexer.lex();
	modules.emplace(filename, std::move(Parser::parseFile(filename, lexer)));
	return &modules[filename];
}
