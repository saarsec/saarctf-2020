#ifndef SCHLOSSBERGCAVES_PARSER_H
#define SCHLOSSBERGCAVES_PARSER_H


#include "Lexer.h"
#include "ast/ast.h"

class SaarlangModule;

class Diagnostic;


class Parser {
	SaarlangModule &module;
	Lexer &lexer;
	Diagnostic &diag;

	std::unordered_map<TokenType, int, std::hash<int>> operatorPrecedence;

	void initOperators() {
		// higher numbers bind stronger
		operatorPrecedence[TT_ASSIGN] = 1;
		operatorPrecedence[TT_OR] = 2;
		operatorPrecedence[TT_AND] = 3;
		operatorPrecedence[TT_XOR] = 4;
		operatorPrecedence[TT_EQUAL] = 5;
		operatorPrecedence[TT_NOTEQUAL] = 5;
		operatorPrecedence[TT_LESS] = 6;
		operatorPrecedence[TT_LESSEQUAL] = 6;
		operatorPrecedence[TT_GREATER] = 6;
		operatorPrecedence[TT_GREATEREQUAL] = 6;
		operatorPrecedence[TT_PLUS] = 7;
		operatorPrecedence[TT_MINUS] = 7;
		operatorPrecedence[TT_MUL] = 8;
		operatorPrecedence[TT_DIV] = 8;
		operatorPrecedence[TT_MOD] = 8;
		operatorPrecedence[TT_AT] = 9;
	}

	// Grammar: see Lexer.h
	void parse();

	void parseImport();

	void parseConst();

	void parseFunction();

	StatementNode *parseStatement();

	ExpressionNode *parseExpression();

	ExpressionNode *parseExpressionComponent();

	ExpressionNode *parsePrimitiveExpression();;

	TypeNode *parseType();

	[[noreturn]] void reportUnexpectedToken();


public:
	Parser(SaarlangModule &module, Lexer &lexer) : module(module), lexer(lexer), diag(lexer.getDiagnostic()) {
		initOperators();
	}

	static SaarlangModule parseFile(std::string name, Lexer &lexer);
};


#endif //SCHLOSSBERGCAVES_PARSER_H
