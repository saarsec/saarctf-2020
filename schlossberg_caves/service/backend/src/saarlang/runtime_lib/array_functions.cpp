#include <cstdlib>
#include <iostream>
#include "saarlang.h"
#include "../JIT.h"

const uint64_t max_array_size = 1024 * 1024;
const uint64_t max_memory = 1024 * 1024 * 7;

uint64_t used_memory = 0;


static void sl_assert(int condition, const char *msg) {
	if (!condition) {
		std::cerr << "ERROR: " << msg << std::endl;
		throw std::exception();
	}
}

void sl_array_bound_violation() {
	std::cerr << "Array access out of bounds" << std::endl;
	throw std::exception();
}

sl_array_byte *sl_new_array_byte(sl_int size) {
	sl_assert(size >= 0, "size >= 0");
	uint64_t memsize = sizeof(uint64_t) + size * sizeof(sl_byte);
	sl_assert(memsize <= max_array_size, "Trying to reserve too much memory");
	used_memory += memsize;
	sl_assert(used_memory <= max_memory, "Reserving too much memory");

	auto array = (sl_array_byte *) malloc(memsize);
	sl_assert(array != nullptr, "malloc() failed");
	array->length = size;
	return array;
}

sl_array_int *sl_new_array_int(sl_int size) {
	sl_assert(size >= 0, "size >= 0");
	uint64_t memsize = sizeof(uint64_t) + size * sizeof(sl_int);
	sl_assert(memsize <= max_array_size, "Trying to reserve too much memory");
	used_memory += memsize;
	sl_assert(used_memory <= max_memory, "Reserving too much memory");

	auto array = (sl_array_int *) malloc(memsize);
	sl_assert(array != nullptr, "malloc() failed");
	array->length = size;
	return array;
}




void import_array_functions(JitEngine &engine) {
	engine.addFunction("sl_array_bound_violation", (void *) &sl_array_bound_violation);
	engine.addFunction("sl_new_array_byte", (void *) &sl_new_array_byte);
	engine.addFunction("sl_new_array_int", (void *) &sl_new_array_int);
}
