#!/usr/bin/env bash

set -eux

# Build the service - the "service" directory will later be used to install.
# Can be empty if you build everything on the vulnbox. 
# You can remove files here that should never lie on the box.

# Examples:
# cd service && make -j4  # build C binary
# cd service && npm install && npm run build  # use npm to build a frontend
# rm -rf service/.idea  # remove files that should not be on vulnbox

# create build structure
pushd .
mkdir go
mkdir go/src
mkdir go/src/stahl4
cp -r service go/src/stahl4/

# set $GOPATH to repo go folder
export GOPATH=$PWD/go

cd go/src/stahl4

# start building
go build -o stahl4.bin ./service

popd
mv go/src/stahl4/stahl4.bin service/stahl4.bin
