#include <cstdlib>
#include <iostream>
#include "saarlang.h"
#include "../JIT.h"
#include "../../caves/CaveMap.h"


using namespace std;

static CaveMap *cave = nullptr;

static uint64_t steps = 0;

static vector<Position> path;

static vector<Treasure> treasures;


void setCurrentMap(CaveMap *cavemap) {
	cave = cavemap;
	steps = 0;
	path.clear();
	path.push_back(cave->getCurrentField());
}

vector<Position> &getVisitedPath() {
	return path;
}

vector<Treasure> &getFoundTreasures() {
	return treasures;
}


static void sl_assert(int condition, const char *msg) {
	if (!condition) {
		std::cerr << "ERROR: " << msg << std::endl;
		throw std::exception();
	}
}

sl_int move(int direction) {
	sl_assert(!cave->finished, "Your tour is already finished");
	sl_assert(steps < cave->allowedSteps, "Number of steps exceeded");
	sl_assert(cave->move(direction), "You ran into a wall");
	steps++;
	path.push_back(cave->getCurrentField());
	return 1;
}

sl_int move_up() { return move(1); }

sl_int move_right() { return move(2); }

sl_int move_down() { return move(3); }

sl_int move_left() { return move(4); }

sl_int finish() {
	sl_assert(!cave->finished, "Your tour is already finished");
	return cave->finish(treasures);
}

sl_int currentX() {
	return cave->posX;
}

sl_int currentY() {
	return cave->posY;
}

// stdlib
sl_int wallAt(sl_int x, sl_int y) {
	sl_assert(0 <= x && x < cave->width, "width out of range");
	sl_assert(0 <= y && y < cave->height, "height out of range");
	return cave->get(x, y) != CAVE_FREE;
}


void import_cave_functions(JitEngine &engine) {
	// Intrinsics
	engine.addFunction("sl_ruff", (void *) &move_up);
	engine.addFunction("sl_doniwwer", (void *) &move_right);
	engine.addFunction("sl_runner", (void *) &move_down);
	engine.addFunction("sl_riwwer", (void *) &move_left);
	engine.addFunction("sl_ferdisch", (void *) &finish);
	engine.addFunction("sl_wo_x", (void *) &currentX);
	engine.addFunction("sl_wo_y", (void *) &currentY);
	// Libary function (use 'holmol "cavelib.sl";' to see it)
	engine.addFunction("issdowas", (void *) &wallAt);
}
