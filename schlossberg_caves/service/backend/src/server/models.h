#ifndef SCHLOSSBERGCAVES_MODELS_H
#define SCHLOSSBERGCAVES_MODELS_H

#include <string>
#include "../caves/CaveMap.h"

class Cave {
public:
	std::string id;
	time_t created;
	std::string name;
	int template_id;
	std::string owner;

	std::vector<Treasure> treasures;

	Cave() : template_id(0), created(0) {};

	Cave(std::string name, int template_id, std::string owner)
			: name(std::move(name)), created(time(nullptr)), template_id(template_id), owner(std::move(owner)) {}

	const inline std::string &getKey() const {
		return id;
	}

	nlohmann::json to_public_json() const;
};

class User {
public:
	std::string id;
	std::string username;
	std::string password;
	int created;

	User(std::string username, std::string password)
			: username(std::move(username)), password(std::move(password)), created(0) {}

	User() : created(0) {}

	const inline std::string &getKey() const {
		return username;
	}
};


// Methods to serialize / deserialize these classes

void to_json(nlohmann::json &j, const Treasure &treasure);

void from_json(const nlohmann::json &j, Treasure &treasure);

void to_json(nlohmann::json &j, const Cave &cave);

void from_json(const nlohmann::json &j, Cave &cave);

void to_json(nlohmann::json &j, const User &user);

void from_json(const nlohmann::json &j, User &user);

void to_json(nlohmann::json &j, const Position &field);

#endif
