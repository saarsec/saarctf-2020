cmake_minimum_required(VERSION 3.5)
project(SchlossbergCaves)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g")


# import LLVM
find_package(LLVM 7 REQUIRED CONFIG)
message(STATUS "Found LLVM ${LLVM_PACKAGE_VERSION}")
message(STATUS "Using LLVMConfig.cmake in: ${LLVM_DIR}")
add_definitions(${LLVM_DEFINITIONS})
include_directories(${LLVM_INCLUDE_DIRS})
link_directories(${LLVM_LIBRARY_DIRS})

if (${LLVM_STATIC_LINK})
	# static linking
	llvm_map_components_to_libnames(LLVM_LIBS core irreader X86Desc X86Info X86CodeGen X86AsmParser ExecutionEngine McJit)
	set(CMAKE_SHARED_LINKER_FLAGS "-Wl,--as-needed ${CMAKE_SHARED_LINKER_FLAGS}")
	MESSAGE(STATUS "Linking LLVM in a static way")
else (${LLVM_STATIC_LINK})
	# shared linking
	set(LLVM_LIBS "LLVM")
endif (${LLVM_STATIC_LINK})


set(SOURCES_COMPILER
		src/saarlang/Lexer.h
		src/saarlang/Diagnostic.cpp
		src/saarlang/Diagnostic.h
		src/saarlang/Parser.cpp
		src/saarlang/Parser.h
		src/saarlang/SaarlangModule.cpp
		src/saarlang/SaarlangModule.h
		src/saarlang/ast/ast.h
		src/saarlang/ast/ast.cpp
		src/saarlang/ast/statements.cpp
		src/saarlang/ast/statements.h
		src/saarlang/ast/expressions.cpp
		src/saarlang/ast/expressions.h
		src/saarlang/Scope.h
		src/saarlang/BuildContext.h
		src/saarlang/JIT.cpp
		src/saarlang/JIT.h
		src/saarlang/BuildContext.cpp
		src/saarlang/runtime_lib/saarlang.h
		src/saarlang/runtime_lib/array_functions.cpp
		src/saarlang/runtime_lib/stdlib_functions.cpp
		src/saarlang/runtime_lib/cave_functions.cpp
		src/caves/CaveMap.cpp
		src/caves/CaveMap.h
		src/saarlang/ast/ast_print.cpp)


set(SOURCES_SERVER
		src/server/httpserver.cpp
		src/server/models.cpp
		src/server/models.h
		src/server/storage.h
		src/server/http_utils.cpp
		src/server/http_utils.h
		src/server/api.cpp
		src/server/api.h
		)

add_executable(SchlossbergCaveServer ${SOURCES_SERVER} ${SOURCES_COMPILER})
target_link_libraries(SchlossbergCaveServer microhttpd stdc++fs ${LLVM_LIBS})


add_executable(SaarlangCompiler src/saarlang/compiler_main.cpp ${SOURCES_COMPILER})
target_link_libraries(SaarlangCompiler ${LLVM_LIBS})
