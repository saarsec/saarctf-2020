#include "http_utils.h"
#include <vector>
#include <string>
#include <fstream>
#include <iostream>
#include <unordered_map>
#include <microhttpd.h>

using namespace std;

vector<char> readFile(string &filename) {
	ifstream fs(filename, std::ios::in | std::ios::binary);
	return vector<char>(std::istreambuf_iterator<char>(fs), std::istreambuf_iterator<char>());
}

static unordered_map<string, vector<char>> fileCache;

vector<char> &getFile(string filename) {
	auto it = fileCache.find(filename);
	if (it == fileCache.end()) {
		return fileCache[filename] = readFile(filename);
	} else {
		return it->second;
	}
}

int HTTPUtilities::respondNotFound(MHD_Connection *connection) {
	const char *page = "404 Not found";
	auto response = MHD_create_response_from_buffer(strlen(page), (void *) page, MHD_RESPMEM_PERSISTENT);
	int ret = MHD_queue_response(connection, MHD_HTTP_NOT_FOUND, response);
	MHD_destroy_response(response);
	return ret;
}

int HTTPUtilities::respondWithStringConstant(MHD_Connection *connection, const char *str, unsigned int status) {
	auto response = MHD_create_response_from_buffer(strlen(str), (void *) str, MHD_RESPMEM_PERSISTENT);
	int ret = MHD_queue_response(connection, status, response);
	MHD_destroy_response(response);
	return ret;
}

int HTTPUtilities::respondWithString(MHD_Connection *connection, const string &str, unsigned int status) {
	auto response = MHD_create_response_from_buffer(str.length(), (void *) str.c_str(), MHD_RESPMEM_MUST_COPY);
	int ret = MHD_queue_response(connection, status, response);
	MHD_destroy_response(response);
	return ret;
}

int HTTPUtilities::respondWithJson(MHD_Connection *connection, const nlohmann::json &data) {
	string str = data.dump();
	auto response = MHD_create_response_from_buffer(str.length(), (void *) str.c_str(), MHD_RESPMEM_MUST_COPY);
	MHD_assert(MHD_add_response_header(response, "Content-Type", CONTENT_JSON));
	int ret = MHD_queue_response(connection, MHD_HTTP_OK, response);
	MHD_destroy_response(response);
	return ret;
}

int HTTPUtilities::respondWithFile(MHD_Connection *connection, string filename, const char *contentType) {
	if (fileCache.find(filename) == fileCache.end()) {
		fileCache[filename] = readFile(filename);
	}
	auto &data = fileCache[filename];
	auto response = MHD_create_response_from_buffer(data.size(), data.data(), MHD_RESPMEM_PERSISTENT);
	MHD_assert(response);
	MHD_assert(MHD_add_response_header(response, "Content-Type", contentType));
	int ret = MHD_queue_response(connection, MHD_HTTP_OK, response);
	MHD_destroy_response(response);
	return ret;
}

pid_t HTTPUtilities::forkConnectionHandler(MHD_Connection *connection, bool dup_stdout, bool dup_stderr) {
	auto info = (union MHD_ConnectionInfo *) MHD_get_connection_info(
			connection, MHD_CONNECTION_INFO_CONNECTION_FD);
	int fd = info->connect_fd;
	pid_t pid = fork();
	if (pid) {
		// Parent - close socket manually (and remove it from libmicrohttpd)
		close(fd);
		info->connect_fd = MHD_INVALID_SOCKET;
		return pid;
	} else {
		// Child process
		// Bind stdout / stderr to the connection
		if (dup_stdout) dup2(fd, 1);
		if (dup_stderr) dup2(fd, 2);
		// send HTTP headers
		string headers = "HTTP/1.1 200 OK\r\n"
				"Content-Type: text/plain; charset=utf-8\r\n"
				"Connection: Closed\r\n"
				"\r\n";
		write(fd, headers.c_str(), headers.size());
		return pid;
	}
}

void PostBody::append(const char *data, size_t length) {
	buffer += std::string(data, length);
}
