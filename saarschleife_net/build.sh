#!/usr/bin/env sh

set -eux

# Build the service - the "service" directory will later be used to install.
# Can be empty if you build everything on the vulnbox. 
# You can remove files here that should never lie on the box.

cd service/backend
rm -rf .idea .gradle build gradle* out .gitignore *.iml src/test

npm install

cd ../frontend
rm -rf .gradle build web .gitignore *.iml .idea

