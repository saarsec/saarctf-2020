#!/usr/bin/env bash

set -eux

# Install the service on a fresh vulnbox. Target should be /home/<servicename>
# You get:
# - $SERVICENAME
# - $INSTALL_DIR
# - An user account with your name ($SERVICENAME)

# 1. Copy/move files
mv service/stahl4.bin "$INSTALL_DIR"/stahl4
# leak handshake file to make exploit3 easier
mkdir "$INSTALL_DIR"/backup
mv service/internal/network/handshake.go "$INSTALL_DIR"/backup/handshake.go
# strip comments from handshake file
sed -i '/^\s*\/\//d' "$INSTALL_DIR"/backup/handshake.go
# set permissions
chown -R "$SERVICENAME:$SERVICENAME" "$INSTALL_DIR"

# 2. make service binary executable
chmod 755 "$INSTALL_DIR"/stahl4

# 3. Configure the server
# ...
# For example: 
# - adjust configs of related services (nginx/databases/...)
# - Build your service if there's source code on the box
# - ...
# 
# Useful commands:
# - nginx-location-add <<EOF
#   location {} # something you want to add to nginx default config (port 80)
#   EOF


# 4. Configure startup for your service
# Typically use systemd for that:
# Install backend as systemd service
# Hint: you can use "service-add-simple '<command>' '<working directory>' '<description>'"
service-add-advanced "$INSTALL_DIR"/stahl4 "$INSTALL_DIR" "Steel Industry 4.0" <<EOF
Restart=always
RestartSec=10
EOF

# Example: Cronjob that removes stored files after a while
#cronjob-add '*/6 * * * * find $INSTALL_DIR/data -mmin +45 -type f -delete'


# Configure ntp
echo '# NTP time server within VPN' >> /etc/ntp.conf
echo 'server  10.32.250.1' >> /etc/ntp.conf



# Example: Initialize Databases (PostgreSQL example)

# Example: 5. Startup database (CI DOCKER ONLY, not on vulnbox)
# if [ -f /.dockerenv ]; then
#     EXAMPLE for PostgreSQL: pg_ctlcluster 11 main start
# fi

# Example: 6. Configure PostgreSQL
# cp $INSTALL_DIR/*.sql /tmp/
# sudo -u postgres psql -v ON_ERROR_STOP=1 -f "/tmp/init.sql"

# Example: 7 Stop services (CI DOCKER ONLY)
# if [ -f /.dockerenv ]; then
#     pg_ctlcluster 11 main stop
# fi
