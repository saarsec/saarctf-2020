
#include "duktape/prep/nondebug/duktape.h"

#define SAARJS_RESULT_SIZE 1024*1024

uint saarjs_buf_offset;

char saarjs_buf[SAARJS_RESULT_SIZE];


#define SAARJS_FETCH_BUF 1024*1024

uint saarjs_fetchbuf_offset;

char saarjs_fetchbuf[SAARJS_FETCH_BUF];


static void check_input(duk_context *ctx, const char *input,  uint64_t *size){
	uint8_t is_valid = 1;
	uint64_t i = 0;
	for (; input[i] != 0; ++i){
		switch (input[i]){
			case '\n':
			case '\t':
			case '\r':
			case '\'':
				is_valid = 0;
			default:
				continue;
		}
	}
	size[0] = i;
	if (!is_valid){
		size[0] = 0;	// set size to 0 to indicate an invalid input
		duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, input); // https://duktape.org/api.html#duk_push_error_object
	}
}

static duk_ret_t native_print(duk_context *ctx) {
	duk_push_string(ctx, " ");
	duk_insert(ctx, 0);
	duk_join(ctx, duk_get_top(ctx) - 1);
	saarjs_buf_offset += snprintf(saarjs_buf + saarjs_buf_offset, SAARJS_RESULT_SIZE - saarjs_buf_offset, "%s\n", duk_safe_to_string(ctx, -1));
	duk_dup_top(ctx);
	return 1;
}

static duk_ret_t native_system(duk_context *ctx) {
	duk_push_string(ctx, " ");
	duk_insert(ctx, 0);
	duk_join(ctx, duk_get_top(ctx) - 1);
	const char *cmd = duk_safe_to_string(ctx, -1);
	FILE *fp = popen(cmd, "r");
	if (fp == NULL) {
		duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, "Failed to execute command: %s", cmd);
		return 1;
	}
	if (SAARJS_RESULT_SIZE-1-saarjs_buf_offset > 0)
		saarjs_buf_offset += fread(saarjs_buf + saarjs_buf_offset, 1, SAARJS_RESULT_SIZE-1-saarjs_buf_offset, fp);
	pclose(fp);

	return 0;
}

static duk_ret_t native_include(duk_context *ctx) {
	duk_push_string(ctx, " ");
	duk_insert(ctx, 0);
	duk_join(ctx, duk_get_top(ctx) - 1);
	const char *filename = duk_safe_to_string(ctx, -1);
	uint64_t len = 0;
	check_input(ctx, filename, &len);
	if (!len)
		return 1;
	FILE *f = fopen(filename, "r");
	if (!f){
		duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, "include(): Failed to open: %s", filename);
		return 1;
	}
	fseek(f, 0, SEEK_END);
	long fsize = ftell(f);
	fseek(f, 0, SEEK_SET);

	char *string = malloc(fsize + 1);
	int i;
	if ((i = fread(string, 1, fsize, f)) <= 0)
		return 0;
	string[i] = 0;
	fclose(f);

	duk_eval_string_noresult(ctx, string);
	free(string);
	return 1;
}


static duk_ret_t native_readfile(duk_context *ctx) {
	duk_push_string(ctx, " ");
	duk_insert(ctx, 0);
	duk_join(ctx, duk_get_top(ctx) - 1);
	const char *filename = duk_safe_to_string(ctx, -1);
	uint64_t len = 0;
	check_input(ctx, filename, &len);
	if (!len)
		return 1;
	FILE *f = fopen(filename, "r");
	if (!f) {
		duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, "File not found: %s", filename);
		return 1;
	}
	fseek(f, 0, SEEK_END);
	long fsize = ftell(f);
	fseek(f, 0, SEEK_SET);

	char *string = malloc(fsize + 1);
	int i;
	if ((i = fread(string, 1, fsize, f)) <= 0)
		return 0;
	string[i] = 0;
	fclose(f);

	duk_push_string(ctx, string);
	free(string);
	return 1;
}

const char *SESSION_TEMPLATE = "/tmp/saarlendar_session_XXXXXX";

char *session = NULL;

const char * FETCH_FMT = "timeout 2 wget -O- -q --load-cookies=%s --save-cookies=%s --keep-session-cookies --no-http-keep-alive --no-check-certificate -U 'saarlendar-bot' -- '%s'";

static duk_ret_t native_fetch(duk_context *ctx) {
	if (session == NULL) {
		session = malloc(strlen(SESSION_TEMPLATE) + 1);
		if (session == 0) {
			duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, "Failed to create session file: %s", session);
			return 1;
		}
		memcpy(session, SESSION_TEMPLATE, strlen(SESSION_TEMPLATE) + 1);
		if (mkstemp(session) == -1){
			duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, "Failed to create session file: %s", session);
			free(session);
			session = NULL;
			return 1;
		}
		perror(session);
	}
	duk_push_string(ctx, " ");
	duk_insert(ctx, 0);
	duk_join(ctx, duk_get_top(ctx) - 1);
	const char *url = duk_safe_to_string(ctx, -1);
	uint64_t len = 0;
	check_input(ctx, url, &len);
	if (len) {
		char *cmd = malloc(len+strlen(FETCH_FMT)+2*strlen(SESSION_TEMPLATE)+1);
		snprintf(cmd, len+strlen(FETCH_FMT)+2*strlen(SESSION_TEMPLATE), FETCH_FMT, session, session, url);
		perror(cmd);
		FILE *fp = popen(cmd, "r");
		free(cmd);
		if (fp == NULL) {
			duk_push_error_object(ctx, DUK_ERR_TYPE_ERROR, "Failed to execute command: %s", cmd);
			return 1;
		}
		char * result = "";
		if (SAARJS_FETCH_BUF-1-saarjs_fetchbuf_offset > 0){
			result = saarjs_fetchbuf + saarjs_fetchbuf_offset;
			saarjs_fetchbuf_offset += fread(saarjs_fetchbuf + saarjs_fetchbuf_offset, 1, SAARJS_FETCH_BUF-1-saarjs_fetchbuf_offset, fp);
			saarjs_fetchbuf[saarjs_fetchbuf_offset++] = 0;
		}
		duk_push_string(ctx, result);
		pclose(fp);
	}
	return 1;
}

void saarjs_exec(char *code) {
	duk_context *ctx = duk_create_heap_default();

    saarjs_buf_offset = 0;

	strcpy(saarjs_buf, "\n");

	duk_push_c_function(ctx, native_print, DUK_VARARGS);
	duk_put_global_string(ctx, "print");

	duk_push_c_function(ctx, native_include, DUK_VARARGS);
	duk_put_global_string(ctx, "include");

	duk_push_c_function(ctx, native_readfile, DUK_VARARGS);
	duk_put_global_string(ctx, "readfile");

	duk_push_c_function(ctx, native_system, DUK_VARARGS);
	duk_put_global_string(ctx, "system");

	duk_push_c_function(ctx, native_fetch, DUK_VARARGS);
	duk_put_global_string(ctx, "fetch");
    
	duk_eval_string_noresult(ctx, "include('/home/saarlendar/saarjs-runtime/core.js')");

	duk_eval_string_noresult(ctx, code);

	duk_destroy_heap(ctx);

	if (session != NULL){
		remove(session);
		free(session);
		session = NULL;
	}
}