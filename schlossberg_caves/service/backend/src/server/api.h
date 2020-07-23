#ifndef SCHLOSSBERGCAVES_API_H
#define SCHLOSSBERGCAVES_API_H

#include "http_utils.h"

class API;

API *createAPI();

void freeAPI(API *api);

int handle_request(API *api, struct MHD_Connection *connection, std::string &url, const std::string &method, PostBody *body);

#endif