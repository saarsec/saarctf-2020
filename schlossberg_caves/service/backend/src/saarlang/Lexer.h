#ifndef SCHLOSSBERGCAVES_LEXER_H
#define SCHLOSSBERGCAVES_LEXER_H


#include <string>
#include <istream>
#include <utility>
#include <vector>
#include <unordered_map>
#include "Diagnostic.h"

/*
module := import* (functiondef|constdef)*
import := "holmol" "<filename>";
constdef := "const" <name>: <type> = <expression>;
functiondef := "eija" <name> (<name>: <type>, <name>: <type>...) "gebbtserick" <rtype>: <statement>
statement  :=   { <statement>* }
              | "falls" <expression>: <statement> ["sonschd": <statement>]    // if
              | "solang" <expression>: <statement>    // while
              | <expression>;
              | "var" <name>: <type> [= <expression>];
              | "var" <name>: <array-type> (<constant>);
              | "serick" [<expression>];    // return
expression :=   <name>
              | <constant>
              | <expression> <binop> <expression>
              | "mach" <function> (<expression>*)    // function call
              | "neie" "lischd" <type> (<expression>)    // new array
<binop> := + | - | * | / | % | ^ | = | < | > | <= | >= | == | != | @ | unn | odder
<type>  := int | byte | "lischd" <type>
 */


// ORDER IS IMPORTANT
enum TokenType {
	TT_NONE,
	TT_END,

	TT_NAME,
	TT_NUMBER,
	TT_STRING,

	// Operators
	TT_PLUS,
	TT_MINUS,
	TT_MUL,
	TT_DIV,
	TT_MOD,
	TT_XOR,
	TT_AND,
	TT_OR,
	TT_ASSIGN,
	TT_EQUAL,
	TT_LESS,
	TT_LESSEQUAL,
	TT_GREATER,
	TT_GREATEREQUAL,
	TT_NOT,
	TT_NOTEQUAL,
	TT_AT,

	// Syntax
	TT_BLOCKOPEN,
	TT_BLOCKCLOSE,
	TT_COLON,
	TT_SEMICOLON,
	TT_COMMA,
	TT_PAREN_L,
	TT_PAREN_R,

	// Keywords
	TT_IMPORT,
	TT_CONST,
	TT_FUNCTION,
	TT_RETURNING,
	TT_RETURN,
	TT_IF,
	TT_ELSE,
	TT_WHILE,
	TT_VAR,
	TT_CALL,
	TT_NEW,
	TT_LENGTH,
	TT_ARRAY,
	TT_INT,
	TT_BYTE,

	// Intrinsics / Language Features (defined in cave_functions)
	TT_INTRINSICS_START,
	TT_RUFF, // "up"
	TT_RUNNER, // "down"
	TT_RIWWER, // "left"
	TT_DONIWWER, // "right"
	TT_FERDISCH, // "finish"
	TT_WO_X, // "where x"
	TT_WO_Y, // "where y"
	TT_INTRINSICS_END
};


class Token {
public:
	TokenType type;
	std::string text;
	int line;
	int col;

	Token(TokenType type, std::string text, int line, int col) : type(type), text(std::move(text)), line(line), col(col) {}
};


class Lexer {
	std::istream &input;
	Diagnostic &diag;

	std::vector<Token> tokens;
	int nextToken = 0;

	std::string line;
	int line_number = 0;

	std::unordered_map<std::string, TokenType> keywords;
public:
	Lexer(std::istream &input, Diagnostic &diag) : input(input), diag(diag) {
		initKeywords();
	}

	void initKeywords() {
		keywords["holmol"] = TT_IMPORT;
		keywords["const"] = TT_CONST;
		keywords["eijo"] = TT_FUNCTION;
		keywords["eija"] = TT_FUNCTION;
		keywords["gebbtserick"] = TT_RETURNING;
		keywords["serick"] = TT_RETURN;
		keywords["falls"] = TT_IF;
		keywords["sonschd"] = TT_ELSE;
		keywords["solang"] = TT_WHILE;
		keywords["var"] = TT_VAR;
		keywords["mach"] = TT_CALL;
		keywords["neie"] = TT_NEW;
		keywords["grees"] = TT_LENGTH;
		keywords["lischd"] = TT_ARRAY;
		keywords["int"] = TT_INT;
		keywords["byte"] = TT_BYTE;
		keywords["unn"] = TT_AND;
		keywords["odder"] = TT_OR;

		keywords["ruff"] = TT_RUFF;
		keywords["runner"] = TT_RUNNER;
		keywords["riwwer"] = TT_RIWWER;
		keywords["doniwwer"] = TT_DONIWWER;
		keywords["ferdisch"] = TT_FERDISCH;
		keywords["wo_x"] = TT_WO_X;
		keywords["wo_y"] = TT_WO_Y;
	}

	Diagnostic &getDiagnostic() {
		return diag;
	}

	// Interface to retrieve tokens

	const std::vector<Token> getTokens() {
		return tokens;
	}

	/**
	 * @return next token
	 */
	const Token &peek() {
		return tokens[nextToken];
	}

	/**
	 * @return next token, and move ahead
	 */
	const Token &consume() {
		return tokens[nextToken++];
	}

	/**
	 * Read the next token (and move ahead) if it has a given type, otherwise do nothing
	 */
	bool consumeIf(TokenType tt) {
		if (tokens[nextToken].type == tt) {
			nextToken++;
			return true;
		} else {
			return false;
		}
	}

	/**
	 * Read the next token (and move ahead) if it has a given type, otherwise raise an error
	 */
	const Token &expect(TokenType tt) {
		if (tokens[nextToken].type != tt) {
			diag.parse_error(tokens[nextToken], "Unexpected token, expected " + std::to_string(tt));
		}
		return tokens[nextToken++];
	}


	// Code to generate tokens from source code

	void lex() {
		// Fill "tokens" vector
		int col = 0;
		bool isInComment = false;
		while (std::getline(input, line)) {
			line_number = diag.addLine(line);
			col = 0;
			while (col < line.size()) {
				// next 2 characters
				char c1 = line[col];
				char c2 = (char) (col + 1 < line.size() ? line[col + 1] : '\x00');

				if (isInComment) {
					if (c1 == '*' && c2 == '/') {
						isInComment = false;
						col += 2;
					} else {
						col++;
					}
				} else if (isalpha(c1) || c1 == '_') {
					col = lexWord(col);

				} else if (isdigit(c1)) {
					col = lexConstant(col);

				} else if (c1 == '"') {
					col = lexString(col);

				} else if (c1 == '/' && c2 == '/') {
					break; // ignore rest of line

				} else if (c1 == '/' && c2 == '*') {
					col += 2;
					isInComment = true;

				} else if (isspace(c1)) {
					col++;

				} else {
					// + | - | * | / | % | ^ | = | < | > | == | !=   and {}:;
					TokenType tt = TT_NONE;
					int opsize = 1;
					switch (c1) {
						case '+':
							tt = TT_PLUS;
							break;
						case '-':
							tt = TT_MINUS;
							break;
						case '*':
							tt = TT_MUL;
							break;
						case '/':
							tt = TT_DIV;
							break;
						case '%':
							tt = TT_MOD;
							break;
						case '^':
							tt = TT_XOR;
							break;
						case '<':
						case '>':
						case '=':
						case '!':
							tt = c1 == '<' ? TT_LESS : c1 == '>' ? TT_GREATER : c1 == '=' ? TT_ASSIGN : TT_NOT;
							if (c2 == '=') {
								tt = (TokenType) (((int) tt) + 1);
								opsize++;
							}
							break;
						case '@':
							tt = TT_AT;
							break;
						case '{':
							tt = TT_BLOCKOPEN;
							break;
						case '}':
							tt = TT_BLOCKCLOSE;
							break;
						case '(':
							tt = TT_PAREN_L;
							break;
						case ')':
							tt = TT_PAREN_R;
							break;
						case ':':
							tt = TT_COLON;
							break;
						case ';':
							tt = TT_SEMICOLON;
							break;
						case ',':
							tt = TT_COMMA;
							break;
						default:
							diag.lex_error("Invalid character", line_number, col);
					}
					tokens.emplace_back(tt, line.substr(col, opsize), line_number, col);
					col += opsize;
				}
			}
		}
		tokens.emplace_back(TT_END, "<end>", line_number, col);
	}

	int lexWord(int col) {
		// read a word, which is either a symbol name ("var123") or a keyword ("if", ...)
		int initial = col;
		while (col < line.size() && (isalnum(line[col]) || line[col] == '_')) {
			col++;
		}
		auto str = line.substr(initial, col - initial);
		auto tt = keywords.find(str);
		if (tt == keywords.end()) {
			tokens.emplace_back(TT_NAME, str, line_number, initial);
		} else {
			tokens.emplace_back(tt->second, str, line_number, initial);
		}
		return col;
	}

	int lexConstant(int col) {
		// lex 1, 123, ...
		int initial = col;
		while (col < line.size() && isdigit(line[col])) {
			col++;
		}
		tokens.emplace_back(TT_NUMBER, line.substr(initial, col - initial), line_number, initial);
		return col;
	}

	int lexString(int col) {
		// lex "abc.def"
		col++;
		int initial = col;
		while (col < line.size() && line[col] != '"') {
			col++;
		}
		if (col >= line.size())
			diag.lex_error("Unterminated string literal", line_number, initial);
		tokens.emplace_back(TT_STRING, line.substr(initial, col - initial), line_number, initial);
		return col + 1;
	}

};


#endif
