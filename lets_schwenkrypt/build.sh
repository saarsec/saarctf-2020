#!/usr/bin/env bash

set -eux

# Build the service - the "service" directory will later be used to install.
# Can be empty if you build everything on the vulnbox. 
# You can remove files here that should never lie on the box.

find -iname '*.pyc' -delete
find -iname '__pycache__' -delete
