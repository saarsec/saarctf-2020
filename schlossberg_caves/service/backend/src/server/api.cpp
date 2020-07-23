#include "api.h"
#include <microhttpd.h>
#include <iostream>
#include <fstream>
#include "../../libraries/json.hpp"
#include "../../libraries/picosha2.h"
#include <unordered_map>
#include "models.h"
#include "storage.h"
#include "../saarlang/Diagnostic.h"
#include "../saarlang/SaarlangModule.h"
#include "../saarlang/JIT.h"


using namespace std;
using namespace nlohmann;

#define DATA_DIR string("../../data/")
#define SESSION_COOKIE_NAME "schlossberg_session"


class SessionHandler {
private:
	// session-id => username
	unordered_map<string, string> sessions;
public:
	/**
	 * @param username
	 * @return Sessionid of the new session
	 */
	string login(string &username) {
		string sessionid(16, 0);
		std::generate_n(sessionid.begin(), 16, []() { return 'a' + rand() % 26; });
		sessions[sessionid] = username;
		return sessionid;
	}

	/**
	 * @param sessionid
	 * @return username or "" if the session is invalid
	 */
	string get(const char *sessionid) {
		if (!sessionid)
			return "";
		auto it = sessions.find(sessionid);
		if (it == sessions.end())
			return "";
		return it->second;
	}

	/**
	 * @param connection
	 * @return username or "" if the session is invalid
	 */
	string get(MHD_Connection *connection) {
		auto cookie = MHD_lookup_connection_value(connection, MHD_COOKIE_KIND, SESSION_COOKIE_NAME);
		return get(cookie);
	}

	/**
	 *
	 * @param sessionid
	 * @return username of the logged in user
	 * @throws AuthorizationException if no session / session is invalid
	 */
	string require(const char *sessionid) {
		if (!sessionid)
			throw AuthorizationException();
		auto it = sessions.find(sessionid);
		if (it == sessions.end())
			throw AuthorizationException();
		return it->second;
	}

	/**
	 * @param connection
	 * @return username of the logged in user
	 * @throws AuthorizationException if no session / session is invalid
	 */
	string require(MHD_Connection *connection) {
		auto cookie = MHD_lookup_connection_value(connection, MHD_COOKIE_KIND, SESSION_COOKIE_NAME);
		return require(cookie);
	}
};


/**
 * Hand out cave templates
 */
class CaveTemplateApi {
public:
	int list(MHD_Connection *connection, string &url) {
		return HTTPUtilities::respondWithFile(connection, DATA_DIR + "cave-templates/caves.json", CONTENT_JSON);
	}

	int get(MHD_Connection *connection, string &url, const string &params) {
		auto id = stoi(params);
		if (id >= 1 && id <= 50) {
			return HTTPUtilities::respondWithFile(connection,
												  DATA_DIR + "cave-templates/schlossberg_" + to_string(id) + ".cave",
												  CONTENT_BINARY);
		} else {
			return HTTPUtilities::respondNotFound(connection);
		}
	}
};


class CaveApi {
private:
	SessionHandler &sessionHandler;
	Storage<Cave> caves;

public:
	explicit CaveApi(SessionHandler &sessionHandler) : sessionHandler(sessionHandler), caves(DATA_DIR + "caves") {}

	int list(MHD_Connection *connection, string &url) {
		json response = json::array({});
		auto deadline = time(nullptr) - 20 * 60;
		for (const auto &it : caves) {
			response += it.second.to_public_json();
		}
		return HTTPUtilities::respondWithJson(connection, response);
	}

	int get(MHD_Connection *connection, string &url, const string &params) {
		auto cave = caves.load(params);
		if (sessionHandler.get(connection) == cave.owner) {
			return HTTPUtilities::respondWithJson(connection, cave);
		} else {
			return HTTPUtilities::respondWithJson(connection, cave.to_public_json());
		}
	}

	int rent(MHD_Connection *connection, string &url, PostBody *params) {
		json body = params->to_json();
		auto username = sessionHandler.require(connection);
		auto cave = caves.store(Cave(body["name"], body["template"], username));
		return HTTPUtilities::respondWithJson(connection, cave);
	}

	int hideTreasure(MHD_Connection *connection, string &url, PostBody *params) {
		json body = params->to_json();
		auto username = sessionHandler.require(connection);
		auto cave = caves.load(body["cave_id"]);
		if (cave.owner != username)
			throw AuthorizationException();

		auto tmpl = getFile(DATA_DIR + "cave-templates/schlossberg_" + to_string(cave.template_id) + ".cave");
		CaveMap caveMap(tmpl, 0, cave.treasures);
		for (auto &name: body["names"]) {
			Position pos = caveMap.getRandomPosition();
			cave.treasures.emplace_back(name, pos.x, pos.y);
		}
		caves.update(cave);
		return HTTPUtilities::respondWithJson(connection, cave);
	}

	int visit(MHD_Connection *connection, string &url, PostBody *params) {
		json body = params->to_json();
		auto cave = caves.load(body["cave_id"]);

		if (HTTPUtilities::forkConnectionHandler(connection, true, true))
			return MHD_NO;

		try {
			// Put code files in the module loader
			Diagnostic diag;
			SimpleModuleLoader loader(diag, "../include/");
			for (auto it = body["files"].begin(); it != body["files"].end(); ++it) {
				auto ss = new stringstream(it.value().get<string>());
				loader.preload(it.key(), *ss);
			}
			// Compile and load code into JIT engine
			JitEngine jitEngine;
			json signatures = "{}"_json;
			for (auto it = body["files"].begin(); it != body["files"].end(); ++it) {
				SaarlangModule *module = loader.load(it.key());
				module->resolveImports(diag, loader);
				module->checkTypes(diag);
				auto &code = module->prepareJITObject();
				signatures[it.key()] = picosha2::hash256_hex_string(code);
				jitEngine.addModule(code);
			}
			jitEngine.init();
			cout << "CODE SIGNATURES: " << signatures << endl; // Trusted computing

			// Prepare execution environment
			auto tmpl = getFile(DATA_DIR + "cave-templates/schlossberg_" + to_string(cave.template_id) + ".cave");
			setCurrentMap(new CaveMap(tmpl, 12345, cave.treasures));

			// Run code
			auto result = jitEngine.executeIsolated();
			cout << "Result: " << result << endl;

			// Send result
			cout << "VISITED PATH: ";
			cout << json({{"path",      getVisitedPath()},
						  {"treasures", getFoundTreasures()}});
		} catch (const std::exception &e) {
			cout << "Abort.";
		}

		exit(0);
	}
};


class UserApi {
private:
	SessionHandler &sessionHandler;
	Storage<User> users;

	int respondWithLogin(MHD_Connection *connection, string username, const string &sessionid) {
		string str = (json{{"username", username}}).dump();
		auto response = MHD_create_response_from_buffer(str.length(), (void *) str.c_str(), MHD_RESPMEM_MUST_COPY);
		MHD_assert(MHD_add_response_header(response, "Content-Type", CONTENT_JSON));
		auto cookie = SESSION_COOKIE_NAME "=" + sessionid + "; HttpOnly; Path=/";
		MHD_assert(MHD_add_response_header(response, "Set-Cookie", cookie.c_str()));
		int ret = MHD_queue_response(connection, MHD_HTTP_OK, response);
		MHD_destroy_response(response);
		return ret;
	}

public:
	explicit UserApi(SessionHandler &sessionHandler) : sessionHandler(sessionHandler), users(DATA_DIR + "users") {}

	int login(MHD_Connection *connection, string &url, PostBody *params) {
		json body = params->to_json();

		auto user = users.load(body["username"]);
		string hash = picosha2::hash256_hex_string(body.at("password").get<string>());
		if (user.password != hash) {
			throw AuthorizationException();
		}

		auto sessionId = sessionHandler.login(user.username);
		return respondWithLogin(connection, user.username, sessionId);
	}

	int registerUser(MHD_Connection *connection, string &url, PostBody *params) {
		json body = params->to_json();
		string username = body["username"];
		if (users.exists(username)) {
			return HTTPUtilities::respondWithStringConstant(connection, "Username taken", MHD_HTTP_INTERNAL_SERVER_ERROR);
		}
		string hash = picosha2::hash256_hex_string(body.at("password").get<string>());
		users.store(User(username, hash));
		auto sessionid = sessionHandler.login(username);
		return respondWithLogin(connection, username, sessionid);
	}

	int current(MHD_Connection *connection, string &url) {
		auto username = sessionHandler.get(connection);
		if (username.empty()) {
			return HTTPUtilities::respondWithJson(connection, {{"username", nullptr}});
		} else {
			return HTTPUtilities::respondWithJson(connection, {{"username", username}});
		}
	}

	int logout(MHD_Connection *connection, string &url, PostBody *params) {
		auto response = MHD_create_response_from_buffer(0, (void *) "", MHD_RESPMEM_PERSISTENT);
		MHD_assert(MHD_add_response_header(response, "Set-Cookie", SESSION_COOKIE_NAME
						   "=; HttpOnly; Path=/"));
		int ret = MHD_queue_response(connection, MHD_HTTP_OK, response);
		MHD_destroy_response(response);
		return ret;
	}

};


class API {
private:
	SessionHandler sessions;
	CaveTemplateApi caveTemplateApi;
	CaveApi caveApi;
	UserApi userApi;

public:
	API() : caveApi(sessions), userApi(sessions) {}

	int handle_request(struct MHD_Connection *connection, string &url, const string &method, PostBody *body) {
		// All requests arrive here. Routing happens here.
		cout << method << " " << url << endl;

		try {

			if (method == "GET") {
				if (url == "/api/templates/list")
					return caveTemplateApi.list(connection, url);

				if (startswith(url, "/api/templates/"))
					return caveTemplateApi.get(connection, url, url.substr(strlen("/api/templates/")));

				if (url == "/api/caves/list")
					return caveApi.list(connection, url);

				if (startswith(url, "/api/caves/"))
					return caveApi.get(connection, url, url.substr(strlen("/api/caves/")));

				if (url == "/api/users/current")
					return userApi.current(connection, url);
			}

			if (method == "POST") {
				if (url == "/api/caves/rent")
					return caveApi.rent(connection, url, body);

				if (url == "/api/caves/hide-treasures")
					return caveApi.hideTreasure(connection, url, body);

				if (url == "/api/visit")
					return caveApi.visit(connection, url, body);

				if (url == "/api/users/login")
					return userApi.login(connection, url, body);

				if (url == "/api/users/register")
					return userApi.registerUser(connection, url, body);

				if (url == "/api/users/logout")
					return userApi.logout(connection, url, body);
			}

			return HTTPUtilities::respondNotFound(connection);

		} catch (const std::exception &e) {
			cerr << "Exception: " << e.what() << endl;
			return HTTPUtilities::respondWithString(connection, e.what(), 500);
		}
	}
};


API *createAPI() {
	return new API();
}

void freeAPI(API *api) {
	delete api;
}

int handle_request(API *api, struct MHD_Connection *connection, string &url, const string &method, PostBody *body) {
	return api->handle_request(connection, url, method, body);
}
