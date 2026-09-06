#!/bin/bash
target=./...
if [ ! -z $1 ]; then
target=$1
fi

go test ${target} -tags debug -coverprofile=.coverage-all -vet=off \
  && cat .coverage-all | egrep -v "zz_|generated.pb.go|cmd/|pkg/generated" > .coverage \
  && go tool cover -func=.coverage
