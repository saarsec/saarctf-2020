#include "../../libraries/json.hpp"
#include "models.h"

/*
 * Conversion between User/Cave and their json representation
 */

nlohmann::json Cave::to_public_json() const {
	return nlohmann::json{{"id",             id},
						  {"name",           name},
						  {"template_id",    template_id},
						  {"owner",          owner},
						  {"created",        created},
						  {"treasure_count", treasures.size()}};
}

void to_json(nlohmann::json &j, const Treasure &treasure) {
	j = nlohmann::json{{"name", treasure.name},
					   {"x",    treasure.x},
					   {"y",    treasure.y}};
}

void from_json(const nlohmann::json &j, Treasure &treasure) {
	treasure.name = j.at("name").get<std::string>();
	treasure.x = j.at("x").get<int>();
	treasure.y = j.at("y").get<int>();
}

void to_json(nlohmann::json &j, const Cave &cave) {
	j = nlohmann::json{{"id",             cave.id},
					   {"name",           cave.name},
					   {"template_id",    cave.template_id},
					   {"owner",          cave.owner},
					   {"created",        cave.created},
					   {"treasure_count", cave.treasures.size()},
					   {"treasures",      cave.treasures}};
}

void from_json(const nlohmann::json &j, Cave &cave) {
	cave.id = j.at("id").get<std::string>();
	cave.name = j.at("name").get<std::string>();
	cave.template_id = j.at("template_id").get<int>();
	cave.owner = j.at("owner").get<std::string>();
	cave.created = j.at("created").get<int>();
	cave.treasures = j.at("treasures").get<std::vector<Treasure>>();
}

void to_json(nlohmann::json &j, const User &user) {
	j = nlohmann::json{{"id",       user.id},
					   {"username", user.username},
					   {"password", user.password},
					   {"created",  user.created}};
}

void from_json(const nlohmann::json &j, User &user) {
	user.id = j.at("id").get<std::string>();
	user.username = j.at("username").get<std::string>();
	user.password = j.at("password").get<std::string>();
	user.created = j.at("created").get<int>();
}

void to_json(nlohmann::json &j, const Position &field) {
	j = nlohmann::json({{"x", field.x},
						{"y", field.y}});
}
