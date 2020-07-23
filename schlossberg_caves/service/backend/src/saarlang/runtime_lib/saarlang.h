#ifndef SCHLOSSBERGCAVES_SAARLANG_H
#define SCHLOSSBERGCAVES_SAARLANG_H

/*
 * SAARLANG RUNTIME LIBRARY - C/C++ INTERFACE
 */

#include <cstdint>
#include <vector>

class JitEngine;

class CaveMap;

class Position;

class Treasure;


// Saarlang types

typedef int64_t sl_int;
typedef uint8_t sl_byte;


// Saarlang arrays

typedef struct {
	uint64_t length;
	sl_int data[];
} sl_array_int;

typedef struct {
	uint64_t length;
	sl_byte data[];
} sl_array_byte;



// Runtime library functions / symbols

void import_array_functions(JitEngine &engine);

void import_stdlib_functions(JitEngine &engine);

void import_cave_functions(JitEngine &engine);

std::vector<Position> &getVisitedPath();

std::vector<Treasure> &getFoundTreasures();

static inline void importSaarlangLibrary(JitEngine &engine) {
	import_array_functions(engine);
	import_stdlib_functions(engine);
	import_cave_functions(engine);
}

// Runtime library interface (C++ side)

void setCurrentMap(CaveMap *cavemap);

#endif
