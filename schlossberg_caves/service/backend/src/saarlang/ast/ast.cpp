#include <llvm/Support/raw_ostream.h>
#include "ast.h"
#include "statements.h"
#include "../BuildContext.h"


using namespace std;

/*
 * CONTAINS:
 * Typechecking and codegen for definitions (var, function, ...)
 * Type handling (TypeNode)
 */


void ConstantDefNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	declare(diag, scope);
	auto exprtype = dynamic_cast<ExpressionNode *>(children[1])->checkType(diag, scope);
	if (!type->isCompatible(exprtype)) {
		diag.type_error(token, "Incompatible types", type, exprtype);
	}
}

void ConstantDefNode::declare(Diagnostic &diag, Scope<DefiningNode> &scope) {
	type = dynamic_cast<TypeNode *>(children[0]);
	if (scope.getFromTopScope(symbol))
		diag.scope_error(token, "Already defined", symbol);
	scope.set(symbol, this);
}

void ConstantDefNode::generateCode(BuildContext &context) {
	// define variable
	auto vartype = type->getLLVMType(context);
	auto variable = context.module->getOrInsertGlobal(symbol, vartype);
	cast<GlobalVariable>(variable)->setInitializer(ConstantAggregateZero::get(vartype));

	// initial value
	context.builder.SetInsertPoint(&context.initFunction->getEntryBlock());
	Value *value = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);
	context.builder.CreateStore(value, variable);
}

Value *ConstantDefNode::getValue(BuildContext &context) {
	auto var = context.module->getOrInsertGlobal(symbol, type->getLLVMType(context));
	return context.builder.CreateLoad(var);
}



void FunctionDefNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	declare(diag, scope);
	scope.push();
	scope.set("__current_function", this);
	for (auto &arg: arguments) {
		if (scope.getFromTopScope(arg->symbol)) {
			diag.scope_error(arg->token, "Symbol already defined", symbol);
		}
		scope.set(arg->symbol, arg);
	}
	dynamic_cast<StatementNode *>(children.back())->checkType(diag, scope);
	scope.pop();
}

void FunctionDefNode::declare(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto functionType = type = new FunctionTypeNode(token, dynamic_cast<TypeNode *>(children[0])->clone());
	for (auto &arg: arguments) {
		functionType->children.push_back(arg->getType()->clone());
	}
	if (scope.getFromTopScope(symbol))
		diag.scope_error(token, "Symbol already defined", symbol);
	scope.set(symbol, this);
}

FunctionDefNode::~FunctionDefNode() {
	for (auto ptr: arguments)
		delete ptr;
	delete type;
}

void FunctionDefNode::generateCode(BuildContext &context) {
	auto function = cast<Function>(context.module->getOrInsertFunction(symbol, cast<FunctionType>(type->getLLVMType(context))));
	auto entry = BasicBlock::Create(context.ctx, symbol + "_entry", function);
	context.builder.SetInsertPoint(entry);

	int i = 0;
	for (auto &arg : function->args()) {
		auto argnode = dynamic_cast<FunctionArgumentNode *>(arguments[i++]);
		arg.setName(argnode->symbol);
		argnode->argReference = &arg;
	}

	for (auto &var: variables) {
		var->generateCodeDefinition(context);
	}

	dynamic_cast<StatementNode *>(children[1])->generateCode(context);
	// Create default return
	Type *returnType = dynamic_cast<TypeNode *>(children[0])->getLLVMType(context);
	Value *defaultReturn = returnType->isPointerTy() ? ConstantPointerNull::get(cast<PointerType>(returnType)) : ConstantInt::get(returnType, 0);
	context.builder.CreateRet(defaultReturn);
}

Value *FunctionDefNode::getValue(BuildContext &context) {
	return context.module->getOrInsertFunction(symbol, cast<FunctionType>(type->getLLVMType(context)));
}




TypeNode *DefiningNode::getType() {
	return type;
}

llvm::Value *FunctionArgumentNode::getValue(BuildContext &context) {
	return argReference;
}

llvm::Type *TypeNode::getLLVMType(BuildContext &context) {
	if (basicType == TT_INT)
		return isArray ? context.tIntArray : (Type *) context.tInt;
	else if (basicType == TT_BYTE)
		return isArray ? context.tByteArray : (Type *) context.tByte;
	else
		throw exception();
}

llvm::Type *FunctionTypeNode::getLLVMType(BuildContext &context) {
	vector<Type *> params;
	for (auto &p: children) {
		params.push_back(dynamic_cast<TypeNode *>(p)->getLLVMType(context));
	}
	return FunctionType::get(returnType->getLLVMType(context), params, false);
}
