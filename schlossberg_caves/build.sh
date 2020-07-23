#!/usr/bin/env bash

set -eux

# Build the service - the "service" directory will later be used to install.
# Can be empty if you build everything on the vulnbox. 
# You can remove files here that should never lie on the box.

if ! python3 > /dev/null 2>&1 ; then
	apt-get update
	apt-get install -y --no-install-recommends python3
fi

cd service/frontend
npm install
cd ..

python3 "dist.py"
# sed -i 's|<base href="/"|<base href="/schlossberg/"|' dist/frontend/index.html
