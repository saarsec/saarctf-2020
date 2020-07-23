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
chown -R "root:$SERVICENAME" "$INSTALL_DIR"
# remove all that executable bits (not needed)
find /home/mensaar/ -type f -exec chmod -x {} \;

# 3. Configure the server
# configure php-fpm pool
sed -i 's|^user = .*|user = mensaar|' /etc/php/*/fpm/pool.d/www.conf
sed -i 's|^group = .*|group = mensaar|' /etc/php/*/fpm/pool.d/www.conf
echo 'cgi.fix_pathinfo=0' >> /etc/php/*/fpm/php.ini

# configure the nginx
nginx-location-add <<'EOF'
	# MenSaar service
	location /mensaar/ {
		charset utf-8;
		alias /home/mensaar/;
		index index.php;

		if ($uri ~ "\.php$") {
			fastcgi_pass unix:/run/php/php7.3-fpm.sock;
		}
		include fastcgi_params;
		fastcgi_param  SCRIPT_FILENAME  /home/$fastcgi_script_name;
		fastcgi_param  SCRIPT_NAME  $fastcgi_script_name;
	}
EOF

# Give nginx access rights on the homedir
usermod -aG mensaar www-data

# Configure php-fpm
mkdir -p /etc/systemd/system/php7.3-fpm.service.d
cat > /etc/systemd/system/php7.3-fpm.service.d/override.conf <<'EOF'
[Service]
MemoryAccounting=true
MemoryMax=1024M
LimitNPROC=1024
EOF



# Initialize Databases (PostgreSQL example)

# 5. Startup database (CI DOCKER ONLY, not on vulnbox)
if [ -f /.dockerenv ]; then
	pg_ctlcluster 11 main start
fi

# 6. Configure PostgreSQL
sudo -u postgres psql -v ON_ERROR_STOP=1 <<EOF
DROP DATABASE IF EXISTS mensaar;
DROP ROLE IF EXISTS mensaar;
CREATE ROLE mensaar LOGIN VALID UNTIL 'infinity';
CREATE DATABASE mensaar WITH ENCODING='UTF8' CONNECTION LIMIT=-1;
EOF
sudo -u $SERVICENAME python3 -u "$INSTALL_DIR/initDB.py"

# 7 Stop services (CI DOCKER ONLY)
if [ -f /.dockerenv ]; then
	pg_ctlcluster 11 main stop
fi
