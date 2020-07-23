#ifndef SCHLOSSBERGCAVES_AST_H
#define SCHLOSSBERGCAVES_AST_H

#include <utility>
#include <vector>
#include <string>
#include "../Lexer.h"
#include "../Scope.h"


class TypeNode;

class StatementNode;

class ExpressionNode;

class VarDefStmtNode;

class BuildContext;

namespace llvm {
	class Value;

	class Type;

	class Module;
}


/**
 * The whole sourcecode file is represented as trees of ASTNode instances.
 * "Import" and "Definition" nodes are the roots.
 *
 * Each language construct has its subclass of this type. It is the responsibility of this subclass to:
 * - check that the type rules of this node (and its children) are not violated
 * - generate corresponding LLVM code
 */
class ASTNode {
public:
	const Token token;
	std::vector<ASTNode *> children;

	explicit ASTNode(Token token) : token(std::move(token)) {}

	ASTNode(Token token, std::initializer_list<ASTNode *> children) : token(std::move(token)), children(children) {}

	virtual ~ASTNode() {
		for (auto &child: children) {
			delete child;
		}
	}
};

/**
 * Definitions are: variables, global variables, functions, function arguments
 */
class DefiningNode {
protected:
	TypeNode *type = nullptr;

public:
	const std::string symbol;

	explicit DefiningNode(std::string symbol) : symbol(std::move(symbol)) {}

	/**
	 * Each definition has a saarlang type
	 */
	TypeNode *getType();

	virtual bool isWriteable() {
		return false;
	}

	/**
	 * Read the defined symbol
	 */
	virtual llvm::Value *getValue(BuildContext &context) = 0;

	/**
	 * Get address to the defined symbol (so that we can write to)
	 */
	virtual llvm::Value *getAddress(BuildContext &context) {
		throw std::exception();
	};
};


class ImportNode : public ASTNode {
public:
	const std::string filename;

	ImportNode(const Token &token, std::string filename) : ASTNode(token), filename(std::move(filename)) {}

	void print(std::ostream &out);
};


class GlobalDefinition : public ASTNode, public DefiningNode {
public:
	GlobalDefinition(const Token &token, const std::string &symbol, const std::initializer_list<ASTNode *> &children) :
			ASTNode(token, children), DefiningNode(symbol) {}

	virtual void print(std::ostream &out) = 0;

	/**
	 * Check types of this node and all children, report errors to diag.
	 */
	virtual void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) = 0;

	/**
	 * Add to scope, but do not check its type (useful for imported functions)
	 */
	virtual void declare(Diagnostic &diag, Scope<DefiningNode> &scope) = 0;

	/**
	 * Add generated LLVM code to build context
	 */
	virtual void generateCode(BuildContext &context) = 0;
};


class ConstantDefNode : public GlobalDefinition {
public:
	ConstantDefNode(const Token &token, const std::string &symbol, const std::initializer_list<ASTNode *> &children)
			: GlobalDefinition(token, symbol, children) {}

	void print(std::ostream &out) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void declare(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;

	llvm::Value *getValue(BuildContext &context) override;
};


class FunctionArgumentNode : public ASTNode, public DefiningNode {
public:
	FunctionArgumentNode(const Token &token, const std::string &symbol, TypeNode *type) : ASTNode(token, {(ASTNode *) type}), DefiningNode(symbol) {
		this->type = type;
	};

	llvm::Value *getValue(BuildContext &context) override;

	llvm::Value *argReference = nullptr;
};


class FunctionDefNode : public GlobalDefinition {
public:
	std::vector<FunctionArgumentNode *> arguments;
	std::vector<VarDefStmtNode *> variables;

	// {returntype, body}
	FunctionDefNode(const Token &token, const std::string &symbol) : GlobalDefinition(token, symbol, {}) {}

	~FunctionDefNode() override;

	void print(std::ostream &out) override;

	void checkType(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void declare(Diagnostic &diag, Scope<DefiningNode> &scope) override;

	void generateCode(BuildContext &context) override;

	llvm::Value *getValue(BuildContext &context) override;
};


/**
 * Expression: All constructs that inherently return a value (constant, operator, function call, ...)
 */
class ExpressionNode : public ASTNode {
public:
	ExpressionNode(const Token &token, const std::initializer_list<ASTNode *> &children)
			: ASTNode(token, children) {}

	virtual void print(std::ostream &out) = 0;

	TypeNode *type = nullptr;

	virtual TypeNode *checkType(Diagnostic &diag, Scope<DefiningNode> &scope) = 0;

	virtual bool isWriteable() {
		return false;
	}

	virtual llvm::Value *generateCode(BuildContext &context) = 0;

	virtual llvm::Value *generateAssignCode(BuildContext &context) {
		throw std::exception();
	}

};


/**
 * Type: int, byte, array of int, array of byte, function between these types
 */
class TypeNode : public ASTNode {
public:
	TokenType basicType;
	bool isArray;
	bool isFunction = false;

	TypeNode(const Token &token) : ASTNode(token), basicType(token.type), isArray(false) {};

	TypeNode(const Token &token, TokenType arraytype) : ASTNode(token), basicType(arraytype), isArray(true) {};

	TypeNode(const Token &token, TokenType type, bool isArray) : ASTNode(token), basicType(type), isArray(isArray) {};

	void print(std::ostream &out) {
		if (isArray) out << "lischd ";
		if (basicType == TT_INT) out << "int";
		else out << "byte";
	}

	virtual TypeNode *clone() {
		return new TypeNode(token, basicType, isArray);
	}

	/**
	 * numeric type := int, byte
	 */
	bool isNumeric() {
		return !isArray && !isFunction && (basicType == TT_INT || basicType == TT_BYTE);
	}

	/**
	 * Check if one type can be converted to another type
	 */
	bool isCompatible(const TypeNode *t2) {
		if (isFunction || t2->isFunction) return false;
		if (isArray != t2->isArray) return false;
		if (basicType == TT_INT || basicType == TT_BYTE)
			return t2->basicType == TT_INT || t2->basicType == TT_BYTE;
		return false;
	}

	virtual llvm::Type *getLLVMType(BuildContext &context);
};


class FunctionTypeNode : public TypeNode {
public:
	TypeNode *returnType;

	FunctionTypeNode(const Token &token, TypeNode *rettype) : TypeNode(token), returnType(rettype) {
		isFunction = true;
	}

	~FunctionTypeNode() override {
		delete returnType;
	}

	TypeNode *clone() override {
		auto ft = new FunctionTypeNode(token, returnType->clone());
		for (auto arg: children)
			ft->children.push_back(dynamic_cast<TypeNode *>(arg)->clone());
		return ft;
	}

	llvm::Type *getLLVMType(BuildContext &context) override;

};




static void indentation(std::ostream &out, int indent) {
	for (int i = 0; i < indent; i++) {
		out << "\t";
	}
}

#endif //SCHLOSSBERGCAVES_AST_H
