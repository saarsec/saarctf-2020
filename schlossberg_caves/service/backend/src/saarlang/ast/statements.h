#ifndef SCHLOSSBERGCAVES_STATEMENTS_H
#define SCHLOSSBERGCAVES_STATEMENTS_H

#include "ast.h"

/*
 * Statements are things that can be executed, but return no value. Example: "if", "while", "return", and variable definitions.
 */


class StatementNode : public ASTNode {
public:
	StatementNode(const Token &token, const std::initializer_list<ASTNode *> &children) : ASTNode(token, children) {}

	virtual void print(std::ostream &out, int indent) = 0;

	virtual void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) = 0;

	virtual void generateCode(BuildContext &context) = 0;
};

class BlockStmtNode : public StatementNode {
public:
	explicit BlockStmtNode(const Token &token) : StatementNode(token, {}) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;
};

class ReturnStmtNode : public StatementNode {
public:
	explicit ReturnStmtNode(const Token &token) : StatementNode(token, {}) {}

	ReturnStmtNode(const Token &token, ASTNode *value) : StatementNode(token, {value}) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;
};

class IfStmtNode : public StatementNode {
public:
	// {condition, body, [else case]}
	IfStmtNode(const Token &token, ASTNode *cond, ASTNode *body) : StatementNode(token, {cond, body}) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;
};

class WhileStmtNode : public StatementNode {
public:
	WhileStmtNode(const Token &token, ASTNode *cond, ASTNode *body) : StatementNode(token, {cond, body}) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;
};

class VarDefStmtNode : public StatementNode, public DefiningNode {
protected:
	llvm::Value *ref = nullptr;
public:
	// {type} or {type, init-expr}
	VarDefStmtNode(const Token &token, std::string &symbol, const std::initializer_list<ASTNode *> &children)
			: StatementNode(token, children), DefiningNode(symbol) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;

	virtual void generateCodeDefinition(BuildContext &context);

	llvm::Value *getValue(BuildContext &context) override;

	bool isWriteable() override;

	llvm::Value *getAddress(BuildContext &context) override;
};

class VarArrayStmtNode : public VarDefStmtNode {
public:
	uint64_t size;

	// {type}
	VarArrayStmtNode(const Token &token, std::string &symbol, ASTNode *type, const Token &size)
			: VarDefStmtNode(token, symbol, {type}), size(std::stoull(size.text)) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;

	void generateCodeDefinition(BuildContext &context) override;
};

/**
 * Statement that evaluates an expression and throws the result away. Example. "a = 1;"
 */
class ExpressionStmtNode : public StatementNode {
public:
	explicit ExpressionStmtNode(ASTNode *node) : StatementNode(node->token, {node}) {}

	void print(std::ostream &out, int indent) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;
};

#endif
