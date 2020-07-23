#define CATCH_CONFIG_MAIN

#include "catch1.hpp"
#include "../saarlang/Diagnostic.h"
#include "../saarlang/SaarlangModule.h"
#include "../saarlang/JIT.h"
#include "../caves/CaveMap.h"
#include "output_capture.h"

using namespace std;

static vector<char> readFile(string filename) {
	ifstream fs(filename, std::ios::in | std::ios::binary);
	return vector<char>(std::istreambuf_iterator<char>(fs), std::istreambuf_iterator<char>());
}


void parseFile(const string &fname) {
	Diagnostic diag;
	SimpleModuleLoader loader(diag, "../include/");
	auto module = loader.load("../" + fname);
	REQUIRE(module);
}

void parsePrintParseFile(const string &fname) {
	Diagnostic diag;
	SimpleModuleLoader loader(diag, "../include/");
	auto module = loader.load("../" + fname);
	REQUIRE(module);
	ostringstream ss;
	module->print(ss);
	auto module2 = loader.preload(fname + ".2", ss.str());
	REQUIRE(module2);
}

void typecheckFile(const string &fname) {
	Diagnostic diag;
	SimpleModuleLoader loader(diag, "../include/");
	auto module = loader.load("../" + fname);
	REQUIRE(module);
	module->resolveImports(diag, loader);
	module->checkTypes(diag);
}

void generateCodeFile(const string &fname) {
	Diagnostic diag;
	SimpleModuleLoader loader(diag, "../include/");
	auto module = loader.load("../" + fname);
	REQUIRE(module);
	module->resolveImports(diag, loader);
	module->checkTypes(diag);
	auto s = module->prepareJITObject();
	REQUIRE(s.length() > 0);
	REQUIRE(s.substr(0, 4) == "\177ELF");
}

const string &executeCodeFile(const string &fname, int returncode) {
	Diagnostic diag;
	SimpleModuleLoader loader(diag, "../include/");
	auto module = loader.load("../" + fname);
	REQUIRE(module);
	module->resolveImports(diag, loader);
	module->checkTypes(diag);
	auto object = module->prepareJITObject();
	REQUIRE(object.length() > 0);
	REQUIRE(object.substr(0, 4) == "\177ELF");

	JitEngine jit;
	jit.addModule(object);

	setCurrentMap(new CaveMap(readFile("../../data/cave-templates/schlossberg_1.cave"), 4000));
	importSaarlangLibrary(jit);
	override_print_functions(jit);
	jit.init(false);
	clearOutput();

	auto code = jit.execute();
	REQUIRE(code == returncode);
	return getOutput();
}



TEST_CASE("Parse explore_cave_random.sl") {
	parseFile("samples/explore_cave_random.sl");
}

TEST_CASE("ParsePrintParse explore_cave_random.sl") {
	parsePrintParseFile("samples/explore_cave_random.sl");
}

TEST_CASE("Typecheck explore_cave_random.sl") {
	typecheckFile("samples/explore_cave_random.sl");
}

TEST_CASE("Codegen explore_cave_random.sl") {
	generateCodeFile("samples/explore_cave_random.sl");
}

TEST_CASE("Execute explore_cave_random.sl") {
	auto result = executeCodeFile("samples/explore_cave_random.sl", 0);
	REQUIRE(result.empty());
}


TEST_CASE("Parse language_features.sl") {
	parseFile("samples/language_features.sl");
}

TEST_CASE("ParsePrintParse language_features.sl") {
	parsePrintParseFile("samples/language_features.sl");
}

TEST_CASE("Typecheck language_features.sl") {
	typecheckFile("samples/language_features.sl");
}

TEST_CASE("Codegen language_features.sl") {
	generateCodeFile("samples/language_features.sl");
}

TEST_CASE("Execute language_features.sl") {
	auto result = executeCodeFile("samples/language_features.sl", 1337);
	REQUIRE(result == "3\n246\n103\n720 450\n1\n");
}
