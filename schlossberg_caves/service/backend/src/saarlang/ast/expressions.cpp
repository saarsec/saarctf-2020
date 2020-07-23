#include <llvm/Support/raw_ostream.h>
#include "expressions.h"
#include "../BuildContext.h"

#define assert_numeric(t) {TypeNode* x = t; if (!x->isNumeric()) diag.type_error(token, "Expected numeric type:", x);}

/*
 * Contains: typechecking and codegen for expressions
 */


TypeNode *SymbolExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	usedSymbol = scope.get(token.text);
	if (!usedSymbol)
		diag.scope_error(token, "Undefined reference", token.text);
	return type = usedSymbol->getType();
}

Value *SymbolExprNode::generateCode(BuildContext &context) {
	return usedSymbol->getValue(context);
}

bool SymbolExprNode::isWriteable() {
	return usedSymbol->isWriteable();
}

Value *SymbolExprNode::generateAssignCode(BuildContext &context) {
	return usedSymbol->getAddress(context);
}



TypeNode *ConstantExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	return type = new TypeNode(token, TT_INT, false);
}

ConstantExprNode::~ConstantExprNode() {
	delete type;
}

llvm::Value *ConstantExprNode::generateCode(BuildContext &context) {
	auto val = std::stoll(token.text);
	return ConstantInt::get(context.tInt, val, true);
}



TypeNode *CallExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto calleeType = dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
	if (!calleeType->isFunction)
		diag.type_error(token, "Function type expected", calleeType);
	auto functionType = dynamic_cast<FunctionTypeNode *>(calleeType);
	if (functionType->children.size() != children.size() - 1)
		diag.type_error(token, "Argument number doesn't match", functionType);
	for (int i = 0; i < functionType->children.size(); i++) {
		auto paramtype = dynamic_cast<TypeNode *>(functionType->children[i]);
		auto argtype = dynamic_cast<ExpressionNode *>(children[i + 1])->checkType(diag, scope);
		if (!paramtype->isCompatible(argtype))
			diag.type_error(token, "Incompatible types", paramtype, argtype);
	}
	return type = functionType->returnType;
}

Value *CallExprNode::generateCode(BuildContext &context) {
	auto func = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
	FunctionType *functype = cast<FunctionType>(dynamic_cast<ExpressionNode *>(children[0])->type->getLLVMType(context));
	std::vector<Value *> params;
	for (auto i = 1; i < children.size(); i++) {
		auto param = dynamic_cast<ExpressionNode *>(children[i])->generateCode(context);
		params.push_back(convert(context, param, functype->params()[i - 1]));
	}
	return context.builder.CreateCall(func, params);
}




TypeNode *NewArrayExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	type = dynamic_cast<TypeNode *>(children[0]);
	if (!type->isArray)
		diag.type_error(token, "Expected lischd type, but was ", type);
	auto sizetype = dynamic_cast<ExpressionNode *>(children[1])->checkType(diag, scope);
	if (!sizetype->isNumeric())
		diag.type_error(token, "Expected numeric type", type);
	return type;
}

Value *NewArrayExprNode::generateCode(BuildContext &context) {
	auto size = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);
	auto arrayType = type->getLLVMType(context);
	auto funcname = "sl_new_array_" + std::string(type->basicType == TT_INT ? "int" : "byte");
	auto function = context.module->getOrInsertFunction(funcname, FunctionType::get(arrayType, {context.tInt}, false));
	return context.builder.CreateCall(function, {size});
}



TypeNode *ArrayLengthExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto subtype = dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
	if (!subtype->isArray)
		diag.type_error(token, "must be lischd", subtype);
	return type = new TypeNode(token, TT_INT, false);
}

llvm::Value *ArrayLengthExprNode::generateCode(BuildContext &context) {
	auto arr = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
	auto zero = ConstantInt::get(Type::getInt32Ty(context.ctx), 0);
	auto sizeAddr = context.builder.CreateGEP(arr, {zero, zero});
	return context.builder.CreateLoad(sizeAddr);
}




TypeNode *BinaryOperatorExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	auto t1 = dynamic_cast<ExpressionNode *>(children[0])->checkType(diag, scope);
	auto t2 = dynamic_cast<ExpressionNode *>(children[1])->checkType(diag, scope);
	switch (token.type) {
		case TT_PLUS:
		case TT_MINUS:
		case TT_MUL:
		case TT_DIV:
		case TT_MOD:
		case TT_XOR:
		case TT_AND:
		case TT_OR:
			assert_numeric(t1);
			assert_numeric(t2);
			return type = new TypeNode(token, TT_INT, false);
		case TT_NOTEQUAL:
		case TT_GREATEREQUAL:
		case TT_GREATER:
		case TT_LESSEQUAL:
		case TT_LESS:
		case TT_EQUAL:
			assert_numeric(t1);
			assert_numeric(t2);
			return type = new TypeNode(token, TT_BYTE, false);
		case TT_ASSIGN:
			if (!t1->isCompatible(t2))
				diag.type_error(token, "Incompatible types:", t1, t2);
			if (!dynamic_cast<ExpressionNode *>(children[0])->isWriteable())
				diag.type_error(token, "Left side is not writeable!");
			return type = t2->clone();
		case TT_AT:
			if (!t1->isArray)
				diag.type_error(token, "Lischd type expected", t1);
			assert_numeric(t2);
			return type = new TypeNode(token, t1->basicType, false);
		default:
			diag.type_error(token, "Unknown operator");
	}
}

BinaryOperatorExprNode::~BinaryOperatorExprNode() {
	delete type;
}

Value *BinaryOperatorExprNode::generateCode(BuildContext &context) {
	if (token.type == TT_ASSIGN) {
		Value *a = dynamic_cast<ExpressionNode *>(children[0])->generateAssignCode(context);
		Value *b = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);
		b = convert(context, b, a->getType()->getPointerElementType());
		context.builder.CreateStore(b, a);
		return b;
	}

	if (token.type == TT_AT) {
		Value *a = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
		Value *b = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);

		auto zero = ConstantInt::get(Type::getInt32Ty(context.ctx), 0);
		auto one = ConstantInt::get(Type::getInt32Ty(context.ctx), 1);

		// check array bounds
		auto sizeField = context.builder.CreateGEP(a, {zero, zero});
		auto size = context.builder.CreateLoad(sizeField);
		auto isInBounds = context.builder.CreateICmpULT(convert(context, b, context.tInt), size);
		auto blockInBounds = BasicBlock::Create(context.ctx, "in_bounds", context.builder.GetInsertBlock()->getParent());
		auto blockOutOfBounds = BasicBlock::Create(context.ctx, "out_of_bounds", context.builder.GetInsertBlock()->getParent());
		context.builder.CreateCondBr(isInBounds, blockInBounds, blockOutOfBounds);
		context.builder.SetInsertPoint(blockOutOfBounds);
		context.builder.CreateCall(context.arrayBoundsViolated);
		context.builder.CreateUnreachable();

		context.builder.SetInsertPoint(blockInBounds);
		if (!b->getType()->isIntegerTy(32))
			b = convert(context, b, Type::getInt32Ty(context.ctx));
		auto ptr = context.builder.CreateGEP(a, {zero, one, b});
		return context.builder.CreateLoad(ptr);
	}

	Value *a = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
	Value *b = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);
	if (!a->getType()->isIntegerTy(64))
		a = context.builder.CreateZExt(a, context.tInt);
	if (!b->getType()->isIntegerTy(64))
		b = context.builder.CreateZExt(b, context.tInt);

	switch (token.type) {
		case TT_PLUS:
			return context.builder.CreateAdd(a, b);
		case TT_MINUS:
			return context.builder.CreateSub(a, b);
		case TT_MUL:
			return context.builder.CreateMul(a, b);
		case TT_DIV:
			return context.builder.CreateSDiv(a, b);
		case TT_MOD:
			return context.builder.CreateSRem(a, b);
		case TT_XOR:
			return context.builder.CreateXor(a, b);
		case TT_NOTEQUAL:
			return context.builder.CreateICmpNE(a, b);
		case TT_GREATEREQUAL:
			return context.builder.CreateICmpSGE(a, b);
		case TT_GREATER:
			return context.builder.CreateICmpSGT(a, b);
		case TT_LESSEQUAL:
			return context.builder.CreateICmpSLE(a, b);
		case TT_LESS:
			return context.builder.CreateICmpSLT(a, b);
		case TT_EQUAL:
			return context.builder.CreateICmpEQ(a, b);
		case TT_AND:
			return context.builder.CreateAnd(a, b);
		case TT_OR:
			return context.builder.CreateOr(a, b);
		default:
			throw std::exception();
	}
}

bool BinaryOperatorExprNode::isWriteable() {
	return token.type == TT_AT;
}

Value *BinaryOperatorExprNode::generateAssignCode(BuildContext &context) {
	if (token.type != TT_AT)
		throw std::exception();
	Value *a = dynamic_cast<ExpressionNode *>(children[0])->generateCode(context);
	Value *b = dynamic_cast<ExpressionNode *>(children[1])->generateCode(context);

	auto zero = ConstantInt::get(Type::getInt32Ty(context.ctx), 0);
	auto one = ConstantInt::get(Type::getInt32Ty(context.ctx), 1);

	// check array bounds
	auto sizeField = context.builder.CreateGEP(a, {zero, zero});
	auto size = context.builder.CreateLoad(sizeField);
	auto isInBounds = context.builder.CreateICmpULT(convert(context, b, context.tInt), size);
	auto blockInBounds = BasicBlock::Create(context.ctx, "in_bounds", context.builder.GetInsertBlock()->getParent());
	auto blockOutOfBounds = BasicBlock::Create(context.ctx, "out_of_bounds", context.builder.GetInsertBlock()->getParent());
	context.builder.CreateCondBr(isInBounds, blockInBounds, blockOutOfBounds);
	context.builder.SetInsertPoint(blockOutOfBounds);
	context.builder.CreateCall(context.arrayBoundsViolated);
	context.builder.CreateUnreachable();

	context.builder.SetInsertPoint(blockInBounds);
	if (!b->getType()->isIntegerTy(32))
		b = convert(context, b, Type::getInt32Ty(context.ctx));
	return context.builder.CreateGEP(a, {zero, one, b});
}




TypeNode *IntrinsicExprNode::checkType(Diagnostic &diag, Scope<DefiningNode> &scope) {
	return type = new TypeNode(token, TT_INT, false);
}

llvm::Value *IntrinsicExprNode::generateCode(BuildContext &context) {
	auto funcname = "sl_" + token.text;
	lowercase(funcname);
	auto intrinsic = context.module->getOrInsertFunction(funcname, FunctionType::get(context.tInt, false));
	return context.builder.CreateCall(intrinsic);
}

