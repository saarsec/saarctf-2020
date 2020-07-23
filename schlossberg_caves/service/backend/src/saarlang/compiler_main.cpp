

#include <fstream>
#include <iostream>
#include "Diagnostic.h"
#include "Lexer.h"
#include "Parser.h"
#include "SaarlangModule.h"
#include "JIT.h"
#include "../caves/CaveMap.h"

using namespace std;


static vector<char> readFile(string filename) {
	ifstream fs(filename, std::ios::in | std::ios::binary);
	return vector<char>(std::istreambuf_iterator<char>(fs), std::istreambuf_iterator<char>());
}


static string getFilename(string path) {
	std::size_t found = path.find_last_of("/\\");
	if (found != string::npos)
		return path.substr(found + 1);
	else
		return path;
}


/**
 * Throw a list of files into the Saarlang Compiler.
 * Not used by the service - just for your convenience.
 */
int main(int argc, char *argv[]) {
	if (argc < 2) {
		cerr << "USAGE: " << argv[0] << " <filename>..." << endl;
		return 1;
	}

	bool saveObjects = false;

	Diagnostic diag;
	SimpleModuleLoader loader(diag, "../include/");
	JitEngine je;

	for (int i = 1; i < argc; i++) {
		ifstream f(argv[i]);
		if (f.good()) {
			loader.preload(getFilename(argv[i]), f);
		}
	}

	for (int i = 1; i < argc; i++) {
		if (argv[i] == string("--save")) {
			saveObjects = true;
			continue;
		}

		// STEP 1: Lexer: sourcecode => Tokens
		// STEP 2: Parser: Tokens => AST (abstract syntax tree)
		auto module = loader.load(getFilename(argv[i]));
		module->resolveImports(diag, loader);
		// STEP 3: Typechecker: Check that the AST is valid (types match, variable names are defined, ...)
		module->checkTypes(diag);
		// STEP 4: Translate: AST => LLVM assembly
		je.addModule(module);
		if (saveObjects) {
			ofstream f(string(argv[i]) + ".o");
			f << module->prepareJITObject();
		}
	}

	// Load cave and prepare execution environment
	setCurrentMap(new CaveMap(readFile("../../data/cave-templates/schlossberg_1.cave"), 4000));

	// STEP 5: Execute: LLVM assembly => x86_64 bytecode
	je.init();
	sl_int result = je.executeIsolated();
	cout << "Result: " << result << endl;
}

