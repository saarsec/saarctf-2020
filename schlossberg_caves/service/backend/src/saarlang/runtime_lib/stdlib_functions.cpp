#include "saarlang.h"
#include "../JIT.h"


sl_int print(sl_int x) {
	printf("%ld ", x);
	return 0;
}

sl_int println(sl_int x) {
	printf("%ld\n", x);
	return 0;
}

sl_int println_as_str(sl_array_byte *x) {
	puts((char *) x->data);
	return 0;
}

static sl_int saarlang_version = 1337;


void import_stdlib_functions(JitEngine &engine) {
	// Override with runtime library symbols
	// Use 'holmol "stdlib.sl";' to see these functions
	engine.addSymbol("saarlang_version", &saarlang_version);
	engine.addFunction("sahmol", (void *) &print);
	engine.addFunction("sahmol_ln", (void *) &println);
	engine.addFunction("sahmol_as_str", (void *) &println_as_str);
	srand(time(0));
	engine.addFunction("ebbes", (void *) &rand);
}
