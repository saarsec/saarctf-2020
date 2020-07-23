#!/usr/bin/env bash

# Add default options to gcc and g++. In particular: disable PIE

set -eu

function write_compiler_wrapper {
	echo '#!/bin/bash' > $2
	echo 'if [[ "${@#-shared}" = "$@" ]]; then' >> $2
	echo '  exec -a "$0"' "$1" '"$@" -fdiagnostics-color -fno-pie -no-pie -fpic' >> $2
	echo 'else' >> $2
	echo '  exec -a "$0"' "$1" '"$@" -fdiagnostics-color -fno-pie -fpic' >> $2
	echo 'fi' >> $2
	chmod +x "$2"
}

# Hook using the "alternatives" system
write_compiler_wrapper "/usr/bin/gcc-8" /usr/bin/x86_64-pc-linux-gnu-gcc-8
write_compiler_wrapper "/usr/bin/g++-8" /usr/bin/x86_64-pc-linux-gnu-g++-8
update-alternatives --install /usr/bin/cc cc /usr/bin/x86_64-pc-linux-gnu-gcc-8 30
update-alternatives --install /usr/bin/c++ c++ /usr/bin/x86_64-pc-linux-gnu-g++-8 30

# Hook using a custom symlink
rm /usr/bin/gcc-8 /usr/bin/g++-8
write_compiler_wrapper "/usr/bin/x86_64-linux-gnu-gcc-8" /usr/bin/gcc-8
write_compiler_wrapper "/usr/bin/x86_64-linux-gnu-g++-8" /usr/bin/g++-8
