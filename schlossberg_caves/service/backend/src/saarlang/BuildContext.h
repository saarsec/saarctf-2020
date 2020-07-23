#ifndef SCHLOSSBERGCAVES_BUILDCONTEXT_H
#define SCHLOSSBERGCAVES_BUILDCONTEXT_H

#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/Module.h>
#include <llvm/Transforms/Utils/ModuleUtils.h>
#include <llvm/Support/raw_ostream.h>

using namespace llvm;


/**
 * Container for all compilation-relevant things (LLVM module and builder classes, references to types and functions, ...)
 */
class BuildContext {
public:
	Module *module;

	static LLVMContext ctx;
	IRBuilder<> builder;
	IntegerType *tInt = IntegerType::getInt64Ty(ctx);
	IntegerType *tByte = IntegerType::getInt8Ty(ctx);
	IntegerType *tBool = IntegerType::getInt1Ty(ctx);
	// array = ptr to {size, [0], [1], ...}
	Type *tIntArray = StructType::create({tInt, ArrayType::get(tInt, 0)})->getPointerTo();
	Type *tByteArray = StructType::create({tInt, ArrayType::get(tByte, 0)})->getPointerTo();
	Function* arrayBoundsViolated;

	// This code runs before the execution starts. Append to the initFunctionPos block.
	Function *initFunction;
	BasicBlock *initFunctionPos;

	BuildContext(const std::string &name);

	void finalizeModule();
};

/**
 * Perform automatic type conversion on LLVM level (byte->int and more). Type checking happens on AST level.
 */
Value *convert(BuildContext &context, Value *value, Type *type);

#endif
