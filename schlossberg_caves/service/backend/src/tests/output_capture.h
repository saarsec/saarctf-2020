#ifndef SCHLOSSBERGCAVES_OUTPUT_CAPTURE_H
#define SCHLOSSBERGCAVES_OUTPUT_CAPTURE_H

#include "../saarlang/JIT.h"

void clearOutput();

const std::string &getOutput();

void override_print_functions(JitEngine &engine);

#endif //SCHLOSSBERGCAVES_OUTPUT_CAPTURE_H
