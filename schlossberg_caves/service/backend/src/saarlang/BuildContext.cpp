#include "BuildContext.h"

LLVMContext BuildContext::ctx;

Value *convert(BuildContext &context, Value *value, Type *type) {
	if (type->isIntegerTy()) {
		return context.builder.CreateIntCast(value, type, false);
	} else {
		return context.builder.CreateBitCast(value, type);
	}
}

BuildContext::BuildContext(const std::string &name) : builder(ctx) {
	module = new Module(StringRef(name), ctx);

	// Prepare init function (runs before execution)
	initFunction = cast<Function>(module->getOrInsertFunction("__saarlang_init_" + name, FunctionType::get(Type::getVoidTy(ctx), {}, false)));
	initFunction->setLinkage(Function::InternalLinkage);
	llvm::appendToGlobalCtors(*module, initFunction, 0xffff);
	initFunctionPos = BasicBlock::Create(ctx, "__saarlang_init_entry", initFunction);

	// Helper functions, defined by runtime library
	arrayBoundsViolated = cast<Function>(module->getOrInsertFunction("sl_array_bound_violation", FunctionType::get(Type::getVoidTy(ctx), false)));
}

void BuildContext::finalizeModule() {
	// Finalize init function
	builder.SetInsertPoint(initFunctionPos);
	builder.CreateRetVoid();
}
