#include <cstring>
#include <microhttpd.h>
#include <iostream>
#include <csignal>
#include <netinet/in.h>
#include "http_utils.h"
#include "api.h"

using namespace std;


int answer_to_connection(API *api, struct MHD_Connection *connection,
						 const char *c_url, const char *c_method,
						 const char *version, const char *upload_data,
						 size_t *upload_data_size, PostBody **con_cls) {
	// POST body handling
	// Multiple calls of this function collect body data, the last call (with size = 0) starts request processing.
	if (strcmp(c_method, "POST") == 0) {
		if (*con_cls == nullptr) {
			// First call - create body instance
			*con_cls = new PostBody();
			return MHD_YES;
		} else if (*upload_data_size) {
			(*con_cls)->append(upload_data, *upload_data_size);
			*upload_data_size = 0;
			return MHD_YES;
		}
		// else: Last call - data has been collected. Proceed with request handling.
	}

	// Delegate request to API and catch common errors
	string url(c_url);
	string method(c_method);
	try {
		return handle_request(api, connection, url, method, *con_cls);
	} catch (const AuthorizationException &e) {
		cerr << "Auth error: " << e.what() << endl;
		return HTTPUtilities::respondWithStringConstant(connection, e.what(), MHD_HTTP_FORBIDDEN);
	} catch (const out_of_range &e) {
		cerr << "Out of range: " << e.what() << endl;
		return HTTPUtilities::respondWithStringConstant(connection, e.what(), MHD_HTTP_NOT_FOUND);
	}
}



static void terminate_request(void *cls, struct MHD_Connection *connection, PostBody **con_cls,
							  enum MHD_RequestTerminationCode toe) {
	// Free post body (if one exists)
	if (*con_cls) {
		delete *con_cls;
		*con_cls = nullptr;
	}
}


static struct MHD_Daemon *server_daemon = nullptr;
static API *api;


static void terminateServer(int signum) {
	cout << "Caught signal " << signum << ", terminating." << endl;
	MHD_stop_daemon(server_daemon);
	freeAPI(api);
	exit(0);
}


/**
 * Spin up an HTTP server that answers REST commands, and handles requests in "api.cpp" (API::handle_request).
 */
int main() {
	// Handle waiting child processes (prevents "Zombies")
	struct sigaction sigchld_action{};
	sigchld_action.sa_handler = SIG_DFL;
	sigchld_action.sa_flags = SA_NOCLDWAIT;
	sigaction(SIGCHLD, &sigchld_action, nullptr);

	// Init random generator
	srand(time(0));

	// Spin up microhttpd
	uint16_t port = 9081;
	struct sockaddr_in loopback_addr{};
	memset(&loopback_addr, 0, sizeof(loopback_addr));
	loopback_addr.sin_family = AF_INET;
	loopback_addr.sin_port = htons(port);
	loopback_addr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
	api = createAPI();

	server_daemon = MHD_start_daemon(MHD_USE_SELECT_INTERNALLY, port, nullptr, nullptr,
									 (MHD_AccessHandlerCallback) answer_to_connection, api,
									 MHD_OPTION_SOCK_ADDR, (struct sockaddr *) (&loopback_addr),
									 MHD_OPTION_NOTIFY_COMPLETED, terminate_request, nullptr,
									 MHD_OPTION_END);
	if (nullptr == server_daemon) {
		cerr << "Couldn't create server. Port " << port << " already taken?" << endl;
		return 1;
	}

	signal(SIGINT, terminateServer);
	signal(SIGTERM, terminateServer);

	cout << "Listing on port " << port << " ..." << endl;
	// (void) getchar();
	while (true) sleep(1000);

	MHD_stop_daemon(server_daemon);
	freeAPI(api);
	return 0;
}