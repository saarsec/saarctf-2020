#ifndef SCHLOSSBERGCAVES_JIT_H
#define SCHLOSSBERGCAVES_JIT_H


#include <llvm/IR/Module.h>
#include <llvm/ExecutionEngine/ExecutionEngine.h>
#include "SaarlangModule.h"
#include "runtime_lib/saarlang.h"


void printLLVMIR(llvm::Module *m);


/**
 * Take LLVM and x86_64 files, put them together, execute them.
 * Call addModule() once for every source file.
 * Then use addFunction() / addSymbol() to register the runtime environment.
 */
class JitEngine {
private:
	llvm::ExecutionEngine *engine;
	std::vector<std::unique_ptr<llvm::MemoryBuffer>> buffers;
	std::vector<std::string> initFunctions;

public:
	JitEngine();

	void addModule(const std::string &object);

	void addModule(SaarlangModule *module) {
		addModule(module->prepareJITObject());
	}

	void addFunction(const std::string &name, void *functionAddr);

	void addSymbol(const std::string &name, void *addr);

	void init(bool stdlib = true);

	sl_int execute();

	sl_int executeIsolated();
};


#endif //SCHLOSSBERGCAVES_JIT_H
