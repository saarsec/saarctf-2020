#ifndef SCHLOSSBERGCAVES_MODULE_H
#define SCHLOSSBERGCAVES_MODULE_H


#include <sstream>
#include <utility>
#include "ast/ast.h"

class SimpleModuleLoader;


/**
 * Represents: Single sourcecode file, topmost AST node
 */
class SaarlangModule {
protected:
	std::vector<ImportNode *> imports;
	std::vector<GlobalDefinition *> definitions;
	Scope<DefiningNode> globalSymbols;
	std::string filename;

	llvm::Module *llvmModule = nullptr;
	std::string asmModule;

public:
	SaarlangModule() = default;

	explicit SaarlangModule(std::string name) : filename(std::move(name)) {}

	~SaarlangModule();

	SaarlangModule(const SaarlangModule &) = delete;

	SaarlangModule(SaarlangModule &&m) noexcept : imports(std::move(m.imports)), definitions(std::move(m.definitions)),
												  globalSymbols(std::move(m.globalSymbols)),
												  filename(std::move(m.filename)) {
		m.imports.clear();
		m.definitions.clear();
	}

	friend class Parser;

	void print(std::ostream &out) {
		for (const auto &node: imports) {
			node->print(out);
		}
		for (const auto &node: definitions) {
			out << "\n";
			node->print(out);
		}
	}

	/**
	 * Load imports, parse them and declare their definitions.
	 */
	void resolveImports(Diagnostic &diag, SimpleModuleLoader &loader);

	/**
	 * Recursive typechecking process over the AST
	 */
	void checkTypes(Diagnostic &diag) {
		diag.setFilename(filename);
		for (const auto &node: definitions) {
			node->checkType(diag, globalSymbols);
		}
	}

	/**
	 * AST => LLVM
	 */
	llvm::Module *generateCode();

	/**
	 * LLVM => x86_64 bytecode, as ELF object file
	 */
	std::string &prepareJITObject();

};


/**
 * Load sourcecode files and store references to the generated SaarlangModule instances
 */
class SimpleModuleLoader {
	std::unordered_map<std::string, SaarlangModule> modules;
	std::string basepath;
	Diagnostic &diag;
public:
	SimpleModuleLoader(Diagnostic &diag, std::string basepath) : basepath(std::move(basepath)), diag(diag) {}

	/**
	 * Load a module from a given input, and remember it as "filename"
	 */
	SaarlangModule *preload(const std::string &filename, std::istream &input);

	/**
	 * Load a module from memory, later ::load() calls will not search on disk
	 */
	SaarlangModule *preload(const std::string &filename, const std::string &input) {
		std::stringstream ss(input);
		return preload(filename, ss);
	}

	/**
	 * Load an actual file
	 */
	SaarlangModule *load(const std::string &filename);
};


#endif
