#ifndef SCHLOSSBERGCAVES_EXPRESSIONS_H
#define SCHLOSSBERGCAVES_EXPRESSIONS_H

#include "ast.h"
#include "../Scope.h"

/*
 * An expression is something that can be evaluated to a value. Constants, variable usages, binary operators and more.
 */


class SymbolExprNode : public ExpressionNode {
	DefiningNode *usedSymbol = nullptr;
public:
	SymbolExprNode(const Token &token) : ExpressionNode(token, {}) {}

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;

	/**
	 * Writeable expressions: variable usage and array access
	 */
	bool isWriteable() override;

	/**
	 * Generate code that returns a pointer to this expression's result, so that it can be written.
	 */
	llvm::Value *generateAssignCode(BuildContext &context) override;
};

class ConstantExprNode : public ExpressionNode {
public:
	explicit ConstantExprNode(const Token &token) : ExpressionNode(token, {}) {}

	~ConstantExprNode() override;

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;
};

class CallExprNode : public ExpressionNode {
public:
	// {function, arg1, arg2, ...}
	CallExprNode(const Token &token, ASTNode *function) : ExpressionNode(token, {function}) {}

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;
};

class NewArrayExprNode : public ExpressionNode {
public:
	NewArrayExprNode(const Token &token, const std::initializer_list<ASTNode *> &children) : ExpressionNode(token, children) {};

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;
};

class ArrayLengthExprNode : public ExpressionNode {
public:
	ArrayLengthExprNode(const Token &token, ASTNode *arr) : ExpressionNode(token, {arr}) {};

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;
};

class BinaryOperatorExprNode : public ExpressionNode {
public:
	BinaryOperatorExprNode(const Token &token, ASTNode *a, ASTNode *b) : ExpressionNode(token, {a, b}) {}

	~BinaryOperatorExprNode() override;

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;

	bool isWriteable() override;

	llvm::Value *generateAssignCode(BuildContext &context) override;
};

class IntrinsicExprNode : public ExpressionNode {
public:
	explicit IntrinsicExprNode(const Token &token) : ExpressionNode(token, {}) {};

	void print(std::ostream &out) override;

	TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	llvm::Value *generateCode(BuildContext &context) override;
};


#endif
