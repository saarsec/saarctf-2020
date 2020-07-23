#ifndef SCHLOSSBERGCAVES_SCOPE_H
#define SCHLOSSBERGCAVES_SCOPE_H

#include <string>
#include <algorithm>
#include "ast/ast.h"


// language is case insensitive
static inline void lowercase(std::string &str) {
	std::transform(str.begin(), str.end(), str.begin(), ::tolower);
}


/**
 * Hierarchical mapping of symbol name => definition
 */
template<class T>
class Scope {
	std::vector<std::unordered_map<std::string, T *>> scopes;

public:
	Scope() {
		scopes.emplace_back();
	}

	void push() {
		scopes.emplace_back();
	}

	void pop() {
		scopes.pop_back();
	}

	T *get(std::string key) {
		lowercase(key);
		for (ssize_t i = scopes.size() - 1; i >= 0; i--) {
			const auto &it = scopes[i].find(key);
			if (it != scopes[i].end())
				return it->second;
		}
		return nullptr;
	}

	T *getFromTopScope(std::string key) {
		lowercase(key);
		const auto &it = scopes.back().find(key);
		if (it != scopes.back().end())
			return it->second;
		return nullptr;
	}

	void set(std::string key, T *value) {
		lowercase(key);
		scopes.back()[key] = value;
	}
};

#endif //SCHLOSSBERGCAVES_SCOPE_H
