#include "CaveMap.h"


struct __attribute__ ((packed)) CaveTemplate {
	unsigned int width;
	unsigned int height;
	unsigned int startRow;
	unsigned int startCol;
	uint8_t data[];
};

CaveMap::CaveMap(const std::vector<char> tmpl, unsigned int allowedSteps, const std::vector<Treasure> treasures) :
		allowedSteps(allowedSteps), treasures(treasures) {
	// Parse cave template
	auto caveTemplate = (const CaveTemplate *) tmpl.data();
	width = caveTemplate->width;
	height = caveTemplate->height;
	startX = posX = caveTemplate->startCol;
	startY = posY = caveTemplate->startRow;
	map.reserve(width * height);
	int i = 0;
	int j = 0;
	while (i < width * height) {
		uint8_t c = caveTemplate->data[j++];
		for (unsigned k = 0; k < 8; k++) {
			map[i++] = ((c >> k) & 1) ? CAVE_FREE : CAVE_ROCK;
		}
	}
}

uint8_t CaveMap::get(const unsigned int x, const unsigned int y) {
	return map[x + y * width];
}

bool CaveMap::move(int direction) {
	if (direction == 1 && posY == 0) return false;
	if (direction == 2 && posX == width - 1) return false;
	if (direction == 3 && posY == height - 1) return false;
	if (direction == 4 && posX == 0) return false;
	unsigned int newX = posX;
	unsigned int newY = posY;
	if (direction == 1) newY--;
	else if (direction == 2) newX++;
	else if (direction == 3) newY++;
	else if (direction == 4) newX--;
	if (get(newX, newY) != CAVE_FREE) return false;
	posX = newX;
	posY = newY;
	return true;
}

Position CaveMap::getCurrentField() {
	return {posX, posY};
}

int CaveMap::finish(std::vector<Treasure> &finalTreasures) {
	finished = true;
	for (auto treasure: treasures) {
		if (treasure.x == posX && treasure.y == posY)
			finalTreasures.push_back(treasure);
	}
	return finalTreasures.size();
}

Position CaveMap::getRandomPosition() {
	while (true) {
		int i = rand() % (width * height);
		if (map[i] == CAVE_FREE) {
			return {x: i % width, y: i / width};
		}
	}
}
