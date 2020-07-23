#!/bin/sh

cd /home/saarlendar/

exec ./nginx -p /home/saarlendar/ -c /home/saarlendar/config/nginx.conf -g 'daemon off;'
