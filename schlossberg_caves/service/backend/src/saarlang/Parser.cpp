#include <stack>
#include <utility>
#include "Parser.h"
#include "SaarlangModule.h"
#include "ast/statements.h"
#include "ast/expressions.h"

void Parser::parse() {
	// module := import* (functiondef|constdef)*
	while (lexer.peek().type == TT_IMPORT) {
		parseImport();
	}
	while (lexer.peek().type != TT_END) {
		switch (lexer.peek().type) {
			case TT_CONST:
				parseConst();
				break;
			case TT_FUNCTION:
				parseFunction();
				break;
			default:
				reportUnexpectedToken();
		}
	}
}

void Parser::parseImport() {
	// import "name";
	auto &token = lexer.expect(TT_IMPORT);
	auto filename = lexer.expect(TT_STRING).text;
	lexer.expect(TT_SEMICOLON);
	module.imports.push_back(new ImportNode(token, filename));
}

void Parser::parseConst() {
	// const <name>: <type> = <primitive expression>;
	auto &token = lexer.expect(TT_CONST);
	auto name = lexer.expect(TT_NAME).text;
	lexer.expect(TT_COLON);
	auto type = parseType();
	lexer.expect(TT_ASSIGN);
	auto expr = parseExpression();
	lexer.expect(TT_SEMICOLON);
	module.definitions.push_back(new ConstantDefNode(token, name, {type, expr}));
}

void Parser::parseFunction() {
	// function <name> (<name>: <type>, <name>: <type>...) "returning" <rtype>: <statement>
	auto &token = lexer.expect(TT_FUNCTION);
	auto name = lexer.expect(TT_NAME).text;
	auto node = new FunctionDefNode(token, name);
	lexer.expect(TT_PAREN_L);
	while (lexer.peek().type != TT_PAREN_R) {
		auto &token2 = lexer.expect(TT_NAME);
		auto arg = token2.text;
		lexer.expect(TT_COLON);
		auto type = parseType();
		node->arguments.push_back(new FunctionArgumentNode(token2, arg, type));
		if (lexer.peek().type != TT_PAREN_R)
			lexer.expect(TT_COMMA);
	}
	lexer.expect(TT_PAREN_R);
	lexer.expect(TT_RETURNING);
	auto rettype = parseType();
	lexer.expect(TT_COLON);
	auto body = parseStatement();
	node->children.push_back(rettype);
	node->children.push_back(body);
	module.definitions.push_back(node);
}

StatementNode *Parser::parseStatement() {
	if (lexer.peek().type == TT_BLOCKOPEN) {
		// {...}
		auto block = new BlockStmtNode(lexer.expect(TT_BLOCKOPEN));
		while (lexer.peek().type != TT_BLOCKCLOSE) {
			block->children.push_back(parseStatement());
		}
		lexer.expect(TT_BLOCKCLOSE);
		return block;
	}

	if (lexer.peek().type == TT_RETURN) {
		auto token = lexer.expect(TT_RETURN);
		auto value = parseExpression();
		lexer.expect(TT_SEMICOLON);
		return new ReturnStmtNode(token, value);
	}

	if (lexer.peek().type == TT_IF) {
		// IF <expression>: <body>
		// [ ELSE: <body2> ]
		auto &token = lexer.expect(TT_IF);
		auto expr = parseExpression();
		lexer.expect(TT_COLON);
		auto body = parseStatement();
		auto stmt = new IfStmtNode(token, expr, body);
		if (lexer.consumeIf(TT_ELSE)) {
			lexer.expect(TT_COLON);
			stmt->children.push_back(parseStatement());
		}
		return stmt;
	}

	if (lexer.peek().type == TT_WHILE) {
		auto &token = lexer.expect(TT_WHILE);
		auto expr = parseExpression();
		lexer.expect(TT_COLON);
		auto body = parseStatement();
		return new WhileStmtNode(token, expr, body);
	}

	if (lexer.peek().type == TT_VAR) {
		// VAR <name>: <type> = <expression>;
		// VAR <name>: <array-type> (<const>);
		// VAR <name>: <type>;
		auto &token = lexer.expect(TT_VAR);
		auto name = lexer.expect(TT_NAME).text;
		lexer.expect(TT_COLON);
		auto type = parseType();
		if (lexer.consumeIf(TT_ASSIGN)) {
			auto expr = parseExpression();
			lexer.expect(TT_SEMICOLON);
			return new VarDefStmtNode(token, name, {type, expr});
		} else if (lexer.consumeIf(TT_PAREN_L)) {
			auto size = lexer.expect(TT_NUMBER);
			lexer.expect(TT_PAREN_R);
			lexer.expect(TT_SEMICOLON);
			return new VarArrayStmtNode(token, name, type, size);
		} else {
			lexer.expect(TT_SEMICOLON);
			return new VarDefStmtNode(token, name, {type});
		}
	}

	auto expr = parseExpression();
	lexer.expect(TT_SEMICOLON);
	return new ExpressionStmtNode(expr);
}

ExpressionNode *Parser::parseExpression() {
	// handle binary operators (Shunting-yard algorithm)
	std::stack<ExpressionNode *> operands;
	std::stack<Token> operators;
	operands.push(parseExpressionComponent());
	while (operatorPrecedence.count(lexer.peek().type) > 0) {
		while (!operators.empty() &&
			   operatorPrecedence[operators.top().type] >= operatorPrecedence[lexer.peek().type]) {
			auto b = operands.top();
			operands.pop();
			auto a = operands.top();
			operands.pop();
			operands.push(new BinaryOperatorExprNode(operators.top(), a, b));
			operators.pop();
		}
		operators.push(lexer.consume());
		operands.push(parseExpressionComponent());
	}
	while (!operators.empty()) {
		auto b = operands.top();
		operands.pop();
		auto a = operands.top();
		operands.pop();
		operands.push(new BinaryOperatorExprNode(operators.top(), a, b));
		operators.pop();
	}
	return operands.top();
}

// Everything in an expression that doesn't contain binary operators
ExpressionNode *Parser::parseExpressionComponent() {
	// Intrinsics: runtime interaction symbols (like "new")
	if (lexer.peek().type > TT_INTRINSICS_START && lexer.peek().type < TT_INTRINSICS_END) {
		return new IntrinsicExprNode(lexer.consume());
	}

	switch (lexer.peek().type) {
		case TT_CALL: {
			auto &token = lexer.expect(TT_CALL);
			auto func = new CallExprNode(token, {parsePrimitiveExpression()});
			lexer.expect(TT_PAREN_L);
			while (lexer.peek().type != TT_PAREN_R) {
				func->children.push_back(parseExpression());
				if (lexer.peek().type != TT_PAREN_R)
					lexer.expect(TT_COMMA);
			}
			lexer.expect(TT_PAREN_R);
			return func;
		}
		case TT_NEW: {
			auto &token = lexer.expect(TT_NEW);
			auto type = parseType();
			lexer.expect(TT_PAREN_L);
			auto size = parseExpression();
			lexer.expect(TT_PAREN_R);
			return new NewArrayExprNode(token, {type, size});
		}
		case TT_LENGTH: {
			auto &token = lexer.expect(TT_LENGTH);
			auto arr = parsePrimitiveExpression();
			return new ArrayLengthExprNode(token, arr);
		}
		case TT_PAREN_L:
		case TT_NAME:
		case TT_NUMBER:
			return parsePrimitiveExpression();
		default:
			reportUnexpectedToken();
	}
}

// Symbol, constant or expression in parentheses
ExpressionNode *Parser::parsePrimitiveExpression() {
	switch (lexer.peek().type) {
		case TT_NAME:
			return new SymbolExprNode(lexer.consume());
		case TT_NUMBER:
			return new ConstantExprNode(lexer.consume());
		case TT_PAREN_L: {
			lexer.expect(TT_PAREN_L);
			auto result = parseExpression();
			lexer.expect(TT_PAREN_R);
			return result;
		}
		default:
			reportUnexpectedToken();
	}
}

TypeNode *Parser::parseType() {
	if (lexer.peek().type == TT_ARRAY) {
		auto &token = lexer.consume();
		if (lexer.peek().type == TT_INT || lexer.peek().type == TT_BYTE)
			return new TypeNode(token, lexer.consume().type);
		else
			reportUnexpectedToken();
	}
	if (lexer.peek().type == TT_INT || lexer.peek().type == TT_BYTE)
		return new TypeNode(lexer.consume());
	else
		reportUnexpectedToken();
}

void Parser::reportUnexpectedToken() {
	if (lexer.peek().type == TT_END)
		diag.parse_error(lexer.peek(), "Unexpected end of input");
	else
		diag.parse_error(lexer.peek(), "Unexpected token");
}

SaarlangModule Parser::parseFile(std::string name, Lexer &lexer) {
	SaarlangModule m(std::move(name));
	Parser parser(m, lexer);
	parser.parse();
	return std::move(m);
}
