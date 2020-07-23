#!/usr/bin/env bash

exec socat -s -T20 TCP-LISTEN:4711,fork,reuseaddr EXEC:"python3 -u server.py",setsid
