#ifndef SCHLOSSBERGCAVES_CAVEMAP_H
#define SCHLOSSBERGCAVES_CAVEMAP_H

#include <vector>
#include <string>

#define CAVE_ROCK ((uint8_t) 0)
#define CAVE_FREE ((uint8_t) 1)


class Treasure {
public:
	std::string name{};
	int x;
	int y;

	Treasure(std::string name, int x, int y) : name(std::move(name)), x(x), y(y) {}

	Treasure() : x(0), y(0) {}
};



class Position {
public:
	unsigned int x;
	unsigned int y;
};



class CaveMap {
	std::vector<uint8_t> map;
	std::vector<Treasure> treasures;
public:
	unsigned int width;
	unsigned int height;
	unsigned int startX;
	unsigned int startY;
	unsigned int posX;
	unsigned int posY;
	unsigned int allowedSteps;
	bool finished = false;

	explicit CaveMap(std::vector<char> tmpl, unsigned int allowedSteps, std::vector<Treasure> treasures = {});

	uint8_t get(unsigned int x, unsigned int y);

	/**
	 * Get random, but free position (means: not in a wall)
	 * @return
	 */
	Position getRandomPosition();

	/**
	 *
	 * @param direction (1=top, 2=right, 3=bottom, 4=left)
	 * @return true if successful
	 */
	bool move(int direction);

	Position getCurrentField();

	int finish(std::vector<Treasure>& finalTreasures);
};

#endif
