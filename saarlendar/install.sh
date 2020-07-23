#!/usr/bin/env bash

set -eux

mv service/* "$INSTALL_DIR/"
mkdir "$INSTALL_DIR/logs"
mkdir "$INSTALL_DIR/users"
mkdir "$INSTALL_DIR/events"
chown -R "$SERVICENAME:$SERVICENAME" "$INSTALL_DIR"

service-add-advanced "$INSTALL_DIR/run.sh" "$INSTALL_DIR/" "saarlendar" <<EOF
LimitNPROC=1024
MemoryAccounting=true
MemoryMax=1024M
EOF

