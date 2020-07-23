#!/usr/bin/env bash

set -eux

# set database credentials to user-authentication on socket
sed -i 's|^$dbConnectionConfig\s*=.*|$dbConnectionConfig = "pgsql:dbname=mensaar;user=mensaar";|' service/db-config.php
sed -i 's|^CONNECTION = .*|CONNECTION = "dbname=mensaar user=mensaar"|' service/initDB.py

