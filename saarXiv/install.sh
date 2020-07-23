#!/usr/bin/env bash

set -eux

# Install the service on a fresh vulnbox. Target should be /home/<servicename>
# You get:
# - $SERVICENAME
# - $INSTALL_DIR
# - An user account with your name ($SERVICENAME)

# 1. Install dependencies - dotnet core aspnet
# moved to Dockerfile


# 2. Copy/move files
# /bin/ is our binary file directory (owned by root)
# /app.db is the database (owned by service user)
# some symlinks so that the program finds its files
mkdir "$INSTALL_DIR/bin"
mv service/saarXiv/bin/Release/netcoreapp3.1/publish/* "$INSTALL_DIR/bin/"
chown -R root "$INSTALL_DIR/bin"
ln -s $INSTALL_DIR/bin/appsettings.json $INSTALL_DIR/
ln -s $INSTALL_DIR/bin/wwwroot $INSTALL_DIR/

# enable sudo for nobody:$SERVICENAME without password
cat << EOF > /etc/sudoers.d/sudo_"$SERVICENAME"
$SERVICENAME ALL=(nobody:$SERVICENAME) NOPASSWD: $(which pdflatex)
EOF

# create pdflatex wrapper script
cat << EOF > "$INSTALL_DIR/pdflatex_wrapper"
#!/bin/bash
# drop privileges to nobody user and $SERVICENAME group
# and also add a timeout
chmod -R g+w ./data
timeout 60s sudo -u nobody -g $SERVICENAME $(which pdflatex) \$@ >/dev/null 2>&1
EOF
chmod a+x "$INSTALL_DIR/pdflatex_wrapper"

#chown -R "$SERVICENAME:$SERVICENAME" "$INSTALL_DIR"

# 3. TODO Configure the server
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


# 4. TODO Configure startup for your service
# Typically use systemd for that:
# Install backend as systemd service
# Hint: you can use "service-add-simple '<command>' '<working directory>' '<description>'"
service-add-advanced "dotnet $INSTALL_DIR/bin/saarXiv.dll --urls=http://0.0.0.0:5000" "$INSTALL_DIR/" "saarXiv webservice" <<EOF
LimitNPROC=1024
MemoryAccounting=true
MemoryMax=1024M
EOF

# Example: Cronjob that removes stored files after a while
# cronjob-add '*/6 * * * * find $INSTALL_DIR/data -mmin +45 -type f -delete'



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
