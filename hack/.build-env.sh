#!/bin/bash -x
APP_NAME=gena
BUILD_OPTIONS="-w -s"

if [[ ! -d bin ]]; then
    mkdir bin
fi
