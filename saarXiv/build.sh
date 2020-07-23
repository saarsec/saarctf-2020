#!/usr/bin/env bash

set -eux

# Build the service - the "service" directory will later be used to install.
# Can be empty if you build everything on the vulnbox. 
# You can remove files here that should never lie on the box.

# Examples:
# cd service && make -j4  # build C binary
# cd service && npm install && npm run build  # use npm to build a frontend
# rm -rf service/.idea  # remove files that should not be on vulnbox

cd service
dotnet publish saarXiv/saarXiv.csproj -c Release --self-contained false
# find result in service/saarXiv/bin/Release/netcoreapp3.1/publish/
