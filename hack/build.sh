#!/bin/bash -x
dir=$(dirname $0)

functions="$(dirname "$0")/.build-env.sh"
if [ -f "$functions" ]; then
  source "$functions"
fi

#: set usevcs flag if .git folder is found in root
#: note, that it's a fast check, blank .git folder might
#: break compiling
USEVCS=true
if [ ! -d ${dir}/../.git ]; then
  USEVCS=false
fi

go build \
  -buildvcs=${USEVCS} \
  -trimpath \
  -ldflags "${BUILD_OPTIONS}" \
  -o bin/${APP_NAME} .
