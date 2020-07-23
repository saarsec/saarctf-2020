#!/usr/bin/env bash

set -eux

# Install the service on a fresh vulnbox. Target should be /home/<servicename>
# You get:
# - $SERVICENAME
# - $INSTALL_DIR
# - An user account with your name ($SERVICENAME)

# 1. Install dependencies


# 2. Copy/move files
mv service/* "$INSTALL_DIR/"
chown -R "root:root" $INSTALL_DIR/*
chown -R "$SERVICENAME:$SERVICENAME" "$INSTALL_DIR/storage.db"

# 4. Configure startup for your service
# Typically use systemd for that:
# Install backend as systemd service
# Hint: you can use "service-add-simple '<command>' '<working directory>' '<description>'"
service-add-simple "$INSTALL_DIR/run.sh" "$INSTALL_DIR/" "Let's Schwenkrypt"

# Example: Cronjob that removes stored files after a while
# cronjob-add '*/6 * * * * find $INSTALL_DIR/data -mmin +45 -type f -delete'

