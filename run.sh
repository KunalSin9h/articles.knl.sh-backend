#!/bin/bash

Color_Off='\033[0m' # Text Reset
Bold_Green='\033[1;32m' # Bold Green


SERVE_PORT=$1

if [[ $SERVE_PORT == "" ]]; then
    echo -e "Port not specified, using ${Bold_Green}default port as 5000${Color_Off}"
    SERVE_PORT=5000
fi

PASSWORD="1" \
    DB=./db/dev.db \
    PORT=$SERVE_PORT go run ./src/main.go
