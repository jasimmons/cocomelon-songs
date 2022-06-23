#!/bin/bash

export TYPE="mysql"
export HOST="${DB_HOSTNAME:?missing}"
export PORT="${DB_PORT:?missing}"
export DATABASE="cocomelon"
export LOGIN="${DB_USERNAME:?missing}"
export PASSWORD="${DB_PASSWORD:?missing}"

exec ./shmig -m migrations/ "$@"
