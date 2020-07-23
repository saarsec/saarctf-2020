#ifndef SCHLOSSBERGCAVES_DIAGNOSTIC_H
#define SCHLOSSBERGCAVES_DIAGNOSTIC_H

#include <string>
#include <vector>

class TypeNode;

class Token;

/**
 * Prints error messages. And that's it.
 */
class Diagnostic {
	std::vector<std::string> lines;
	std::string filename;

	void showpos(int row, int col);

public:
	int addLine(std::string &line) {
		lines.push_back(line);
		return lines.size();
	}

	void setFilename(const std::string &filename) {
		this->filename = filename;
	}

	[[noreturn]] void file_error(std::string filename);

	[[noreturn]] void lex_error(const char *message, int row, int col);

	[[noreturn]] void parse_error(const Token &token, std::string message);

	[[noreturn]] void type_error(const Token &token, std::string message, TypeNode *type = nullptr, TypeNode *type2 = nullptr);

	[[noreturn]] void scope_error(const Token &token, std::string message, const std::string &name = "");

};


#endif
