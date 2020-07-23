#ifndef SCHLOSSBERGCAVES_HTTP_UTILS_H
#define SCHLOSSBERGCAVES_HTTP_UTILS_H

#include <cstring>
#include <stdexcept>
#include "../../libraries/json.hpp"

struct MHD_Connection;
struct MHD_Response;


#define MHD_assert(x) do{ if ((x) == MHD_NO) std::cerr << "Failed (" << __FILE__ << ":" << __LINE__ << "): " << #x << endl; }while(false)

#define CONTENT_JSON "application/json; charset=utf-8"
#define CONTENT_BINARY "application/octet-stream"


static inline bool startswith(const std::string &str, const std::string &substr) {
	return strncmp(str.c_str(), substr.c_str(), substr.size()) == 0;
}

static inline bool startswith(const std::string &str, const char *substr) {
	return strncmp(str.c_str(), substr, strlen(substr)) == 0;
}

// Read a file from disk
std::vector<char> &getFile(std::string filename);


class HTTPUtilities {
public:
	static int respondNotFound(MHD_Connection *connection);

	static int respondWithStringConstant(MHD_Connection *connection, const char *str, unsigned int status = 200);

	static int respondWithString(MHD_Connection *connection, const std::string &str, unsigned int status = 200);

	static int respondWithJson(MHD_Connection *connection, const nlohmann::json &data);

	static int respondWithFile(MHD_Connection *connection, std::string filename, const char *contentType);

	/**
	 * Forks a new process that should handle this connection.
	 * - The parent process (result != 0) does not do anything with this connection anymore, it should return MHD_NO
	 * - The child process (result == 0) should send the response body to stdout / stderr, and terminate afterwards.
	 */
	static pid_t forkConnectionHandler(MHD_Connection *connection, bool dup_stdout, bool dup_stderr);

};

class AuthorizationException : public std::runtime_error {
public:
	AuthorizationException() : runtime_error("Authorization required") {};
};

class PostBody {
	std::string buffer;
public:
	void append(const char *data, size_t length);

	const inline std::string &str() const {
		return buffer;
	}

	const inline nlohmann::json to_json() const {
		return nlohmann::json::parse(buffer);
	}
};

#endif
