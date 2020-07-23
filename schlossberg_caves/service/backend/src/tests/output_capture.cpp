#include "output_capture.h"
#include "../saarlang/runtime_lib/saarlang.h"

using namespace std;

static string output;

void clearOutput() {
	output = "";
}

const string &getOutput() {
	return output;
}



static sl_int print(sl_int x) {
	char buffer[256];
	sprintf(buffer, "%ld ", x);
	output = output + string(buffer);
	return 0;
}

static sl_int println(sl_int x) {
	char buffer[256];
	sprintf(buffer, "%ld\n", x);
	output = output + string(buffer);
	return 0;
}

static sl_int println_as_str(sl_array_byte *x) {
	output = output + string((char *) x->data);
	return 0;
}



void override_print_functions(JitEngine &engine) {
	engine.addFunction("sahmol", (void *) &print);
	engine.addFunction("sahmol_ln", (void *) &println);
	engine.addFunction("sahmol_as_str", (void *) &println_as_str);
}