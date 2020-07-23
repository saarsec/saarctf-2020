#!/usr/bin/env bash
set -euxo pipefail

SCRIPTPATH="$(cd "$(dirname "$BASH_SOURCE")" && pwd)"

export INSTALL_DIR="/home/$SERVICENAME"
export PATH=$PATH:$SCRIPTPATH/commands

useradd -m "$SERVICENAME"
chmod 0750 "$INSTALL_DIR"

# do the actual service installation
bash /opt/install-service.sh

# clean-up
rm -rf /opt/service /opt/install.sh /opt/install-service.sh
rm -rf /tmp/*
