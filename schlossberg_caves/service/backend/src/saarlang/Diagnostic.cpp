#include <iostream>
#include "Diagnostic.h"
#include "Lexer.h"
#include "ast/ast.h"


using namespace std;

void Diagnostic::showpos(int row, int col) {
	row -= 1;
	if (row >= lines.size())
		return;
	if (col >= lines[row].size())
		return;
	ssize_t start = col - 32;
	ssize_t length = 100;
	if (start < 0)
		start = 0;
	if (start + length > lines[row].size())
		length = lines[row].size() - start;
	cerr << "> ";
	for (auto c: lines[row].substr(start, length)) {
		if (c == '\t') {
			cerr << "    ";
			col += 3;
		} else {
			cerr << c;
		}
	}
	cerr << endl;

	cerr << "  ";
	for (auto i = start; i < col; i++) {
		cerr << " ";
	}
	cerr << "^" << endl;
}

void Diagnostic::file_error(const std::string filename) {
	cerr << "[ERROR] File not found / not readable: \"" << filename << "\"" << endl;
	throw exception();
}

void Diagnostic::lex_error(const char *message, int row, int col) {
	cerr << "[ERROR] lexer in " << filename << " at " << row << ":" << col << ": " << message << endl;
	showpos(row, col);
	throw exception();
}

void Diagnostic::parse_error(const Token &token, const string message) {
	cerr << "[ERROR] parser in " << filename << " at " << token.line << ":" << token.col << ": " << message << " (" << token.type << ")" << endl;
	showpos(token.line, token.col);
	throw exception();
}

void Diagnostic::type_error(const Token &token, const std::string message, TypeNode *type, TypeNode *type2) {
	cerr << "[ERROR] types in " << filename << " at " << token.line << ":" << token.col << ": " << message << " ";
	if (type)
		type->print(cerr);
	if (type2) {
		cerr << " , ";
		type2->print(cerr);
	}
	cerr << endl;
	showpos(token.line, token.col);
	throw exception();
}

void Diagnostic::scope_error(const Token &token, const std::string message, const string &name) {
	cerr << "[ERROR] symbols in " << filename << " at " << token.line << ":" << token.col << ": " << message;
	if (!name.empty())
		cerr << " \"" << name << "\"";
	cerr << endl;
	showpos(token.line, token.col);
	throw exception();
}
