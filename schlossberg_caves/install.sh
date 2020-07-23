#!/usr/bin/env bash

set -eux

# Install the service on a fresh vulnbox. Target should be /home/<servicename>
# You get:
# - $SERVICENAME
# - $INSTALL_DIR
# - An user account with your name


# 1. Install dependencies
apt-get update
apt-get install -y build-essential gcc g++ cmake libmicrohttpd-dev nginx
./gcc_configure.sh  # compiler patch to disable PIE
apt-get install -y llvm llvm-dev # will be llvm-7 on buster


# 2. Copy/move files
cp -r service/dist/backend service/dist/frontend service/dist/data "$INSTALL_DIR/"
cp $INSTALL_DIR/backend/samples/*.sl $INSTALL_DIR/
chown root:root "$INSTALL_DIR/backend" "$INSTALL_DIR/frontend" $INSTALL_DIR/*.sl
chown -R "$SERVICENAME:$SERVICENAME" "$INSTALL_DIR/data"
# Give nginx access rights on the homedir
usermod -aG $SERVICENAME www-data


# 3. Compile
pushd .
cd $INSTALL_DIR/backend
mkdir build
cd build
cmake ..
make -j 4
find -iname '*.o' -delete
popd


# 4. Configure the server
#nginx-location-add <<EOF
#	# Schlossberg Caves service
#	location /schlossberg/ {
#		charset utf-8;
#		alias $INSTALL_DIR/frontend/;
#	}
#	location /schlossberg/api {
#		rewrite /schlossberg/(.*) /\$1  break;
#		proxy_pass http://127.0.0.1:9081;
#	}
#EOF
cat - > /etc/nginx/sites-available/schlossberg <<EOF
server {
  listen 9080;

	# Schlossberg Caves service
	location / {
		charset utf-8;
		alias $INSTALL_DIR/frontend/;
	}
	location /api {
		proxy_pass http://127.0.0.1:9081;
	}
}
EOF
ln -s /etc/nginx/sites-available/schlossberg /etc/nginx/sites-enabled/schlossberg


# 5. Configure startup for your service
# Typically use systemd for that:
# Install backend as systemd service
service-add-advanced "$INSTALL_DIR/backend/build/SchlossbergCaveServer" "$INSTALL_DIR/backend/build/" "Schlossberg Cave Server" <<EOF
SystemCallFilter=~@debug kill
EOF

# Example: Cronjob that removes stored files after a while
cronjob-add '*/6 * * * * find $INSTALL_DIR/data/caves $INSTALL_DIR/data/users -mmin +45 -type f -delete'




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
