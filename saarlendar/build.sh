#!/usr/bin/env bash

set -eux

# Build the service - the "service" directory will later be used to install.
# Can be empty if you build everything on the vulnbox. 
# You can remove files here that should never lie on the box.

git clone https://github.com/svaarala/duktape.git service/src/duktape
cd service/src/duktape && git checkout a1a7bfafed0e8d88a9204a221855a515094d183e && git apply ../duktape.patch
cd ../../../

cd service && ./configure && make -j4  # build C binary


rm -rf src/duktape src/nginx-1.17.4 src/sandbox.c src/nginx_saarjs.c src/config ./configure ./Makefile src/duktape.patch
# 
# cd service && npm install && npm run build  # use npm to build a frontend
# rm -rf service/.idea  # remove files that should not be on vulnbox
