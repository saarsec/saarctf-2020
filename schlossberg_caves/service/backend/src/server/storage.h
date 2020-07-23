#ifndef SCHLOSSBERGCAVES_STORAGE_H
#define SCHLOSSBERGCAVES_STORAGE_H

#include <string>
#include <map>
#include "../../libraries/json.hpp"
#include <experimental/filesystem>
#include <fstream>
#include <iostream>

namespace filesystem = std::experimental::filesystem;

template<class T>
class Storage {
	std::map<std::string, T> storage;
	std::string basepath;

	void save(T &t) {
		std::ofstream fs(basepath + "/" + t.id);
		fs << nlohmann::json(t);
	}

	void loadall() {
		std::cout << "Loading from " << basepath << std::endl;
		std::experimental::filesystem::v1::path path = basepath;
		if (!std::experimental::filesystem::v1::exists(path))
			std::experimental::filesystem::v1::create_directories(path);
		for (const auto &p : std::experimental::filesystem::v1::directory_iterator(path)) {
			try {
				std::ifstream fs(p.path(), std::ios_base::in);
				nlohmann::json j;
				fs >> j;
				T t = j;
				storage[t.getKey()] = t;
			} catch (const std::exception &e) {
				std::cout << "Failed reading " << p << std::endl;
			}
		}
		std::cout << "Loading done." << std::endl;
	}

public:
	explicit Storage(std::string basepath) : basepath(std::move(basepath)) {
		loadall();
	}

	T &store(T t) {
		std::string id = std::to_string(time(nullptr)) + "_" + std::to_string(rand());
		t.id = id;
		t.created = time(nullptr);
		storage[t.getKey()] = t;
		save(t);
		return storage.at(t.getKey());
	}

	T &load(const std::string &key) {
		return storage.at(key);
	}

	void update(T &t) {
		storage.at(t.getKey()) = t;
		save(t);
	}

	bool exists(const std::string &key) {
		return storage.find(key) != storage.end();
	}

	typename std::map<std::string, T>::iterator begin() {
		return storage.begin();
	}

	typename std::map<std::string, T>::iterator end() {
		return storage.end();
	}
};

#endif