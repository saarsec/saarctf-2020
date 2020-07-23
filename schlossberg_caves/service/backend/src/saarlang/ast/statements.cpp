#include <llvm/Support/raw_ostream.h>
#include "statements.h"
#include "../BuildContext.h"

#define assert_numeric(t) {TypeNode* x = t; if (!x->isNumeric()) diag.type_error(token, "Expected numeric type:", x);}


void ExpressionStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
}

void ExpressionStmtNode::generateCode(BuildContext &context) {
	dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
}



void IfStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto condtype = dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
	assert_numeric(condtype);
	dynamic_cast<StatementNode *>(children[1])->checkType(diag, scope);
	if (children.size() > 2)
		dynamic_cast<StatementNode *>(children[2])->checkType(diag, scope);
}

void IfStmtNode::generateCode(BuildContext &context) {
	Value *cond = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
	cond = context.builder.CreateICmpNE(cond, ConstantInt::get(cond->getType(), 0, false));
	BasicBlock *caseThen = BasicBlock::Create(context.ctx, "if_then", context.builder.GetInsertBlock()->getParent());
	BasicBlock *caseFinal = BasicBlock::Create(context.ctx, "if_final", context.builder.GetInsertBlock()->getParent());
	BasicBlock *caseElse;
	if (children.size() > 2) {
		caseElse = BasicBlock::Create(context.ctx, "if_else", context.builder.GetInsertBlock()->getParent());
	} else {
		caseElse = caseFinal;
	}
	context.builder.CreateCondBr(cond, caseThen, caseElse);

	// then
	context.builder.SetInsertPoint(caseThen);
	dynamic_cast<StatementNode *>(children[1])->generateCode(context);
	context.builder.CreateBr(caseFinal);

	// else
	if (children.size() > 2) {
		context.builder.SetInsertPoint(caseElse);
		dynamic_cast<StatementNode *>(children[2])->generateCode(context);
		context.builder.CreateBr(caseFinal);
	}

	context.builder.SetInsertPoint(caseFinal);
}



void WhileStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto condtype = dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
	assert_numeric(condtype);
	dynamic_cast<StatementNode *>(children[1])->checkType(diag, scope);
}

void WhileStmtNode::generateCode(BuildContext &context) {
	BasicBlock *blockCondition = BasicBlock::Create(context.ctx, "while_cond",
													context.builder.GetInsertBlock()->getParent());
	BasicBlock *blockBody = BasicBlock::Create(context.ctx, "while_body",
											   context.builder.GetInsertBlock()->getParent());
	BasicBlock *blockFinal = BasicBlock::Create(context.ctx, "while_final",
												context.builder.GetInsertBlock()->getParent());
	context.builder.CreateBr(blockCondition);
	context.builder.SetInsertPoint(blockCondition);
	Value *cond = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
	cond = context.builder.CreateICmpNE(cond, ConstantInt::get(cond->getType(), 0, false));
	context.builder.CreateCondBr(cond, blockBody, blockFinal);

	context.builder.SetInsertPoint(blockBody);
	dynamic_cast<StatementNode *>(children[1])->generateCode(context);
	context.builder.CreateBr(blockCondition);

	context.builder.SetInsertPoint(blockFinal);
}



void VarDefStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	type = dynamic_cast<TypeNode *>(children[0]);
	if (children.size() > 1) {
		auto inittype = dynamic_cast<ExpressionNode *>(children[1])->checkType(diag, scope);
		if (!type->isCompatible(inittype))
			diag.type_error(token, "Incompatible types", type, inittype);
	}

	if (scope.getFromTopScope(symbol))
		diag.scope_error(token, "Variable name already taken", symbol);
	scope.set(symbol, this);

	dynamic_cast<FunctionDefNode *>(scope.get("__current_function"))->variables.push_back(this);
}

void VarDefStmtNode::generateCodeDefinition(BuildContext &context) {
	ref = context.builder.CreateAlloca(type->getLLVMType(context), nullptr, "var_" + symbol);
}

void VarDefStmtNode::generateCode(BuildContext &context) {
	if (children.size() > 1) {
		Value *initialValue = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);
		context.builder.CreateStore(convert(context, initialValue, type->getLLVMType(context)), ref);
	}
}

Value *VarDefStmtNode::getValue(BuildContext &context) {
	return context.builder.CreateLoad(ref);
}

bool VarDefStmtNode::isWriteable() {
	return true;
}

Value *VarDefStmtNode::getAddress(BuildContext &context) {
	return ref;
}



void VarArrayStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	type = dynamic_cast<TypeNode *>(children[0]);
	if (!type->isArray)
		diag.type_error(token, "Must be a lischd type", type);

	if (scope.getFromTopScope(symbol))
		diag.scope_error(token, "Variable name already taken", symbol);
	scope.set(symbol, this);

	dynamic_cast<FunctionDefNode *>(scope.get("__current_function"))->variables.push_back(this);
}

void VarArrayStmtNode::generateCode(BuildContext &context) {
}

void VarArrayStmtNode::generateCodeDefinition(BuildContext &context) {
	// reserve array ptr
	VarDefStmtNode::generateCodeDefinition(context);

	// reserve array
	auto arrayStructType = ref->getType()->getPointerElementType()->getPointerElementType();
	auto basicType = arrayStructType->getStructElementType(1)->getArrayElementType();
	auto sizedArrayType = StructType::create({context.tInt, ArrayType::get(basicType, size)});
	Value *array = context.builder.CreateAlloca(sizedArrayType);
	array = context.builder.CreateBitCast(array, ref->getType()->getPointerElementType());
	context.builder.CreateStore(array, ref);

	// fill array with data
	auto zero = ConstantInt::get(Type::getInt32Ty(context.ctx), 0);
	auto one = ConstantInt::get(Type::getInt32Ty(context.ctx), 1);
	context.builder.CreateStore(ConstantInt::get(context.tInt, size), context.builder.CreateGEP(array, {zero, zero}));
}



void ReturnStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto functionType = dynamic_cast<FunctionTypeNode *>(scope.get("__current_function")->getType());
	if (!children.empty()) {
		auto returnValueType = dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
		if (!functionType->returnType->isCompatible(returnValueType))
			diag.type_error(token, "Return type doesn't match", functionType->returnType, returnValueType);
	}
}

void ReturnStmtNode::generateCode(BuildContext &context) {
	if (!children.empty()) {
		Value *val = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
		Type *rettype = context.builder.getCurrentFunctionReturnType();
		context.builder.CreateRet(convert(context, val, rettype));
	} else {
		context.builder.CreateRetVoid();
	}
	auto currentFunction = context.builder.GetInsertBlock()->getParent();
	context.builder.SetInsertPoint(BasicBlock::Create(context.ctx, "deadcode", currentFunction));
}



void BlockStmtNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	scope.push();
	for (auto &stmt: children) {
		dynamic_cast<StatementNode *>(stmt)->checkType(diag, scope);
	}
	scope.pop();
}

void BlockStmtNode::generateCode(BuildContext &context) {
	for (auto &stmt: children) {
		dynamic_cast<StatementNode *>(stmt)->generateCode(context);
	}
}
