#!/bin/bash -e
ALLOWED_CYCLE_LEVEL=20

functions="$(dirname "$0")/.build-env.sh"
if [ -f "$functions" ]; then
  # shellcheck disable=SC1090
  source "$functions"
fi

# Format locally; CI must separately reject unformatted sources.
go fmt ./... || exit "$?"

# Report dependency changes, but leave go.mod and go.sum untouched.
go mod tidy -diff || exit "$?"

#: run ./hack/install-dev-tools.sh once to get these apps.
gocyclo -over ${ALLOWED_CYCLE_LEVEL} -ignore "zz_.*.go|generated.pb.go|pkg/generated.*.go" .
go vet ./...
ineffassign ./...
staticcheck ./...
govulncheck ./...

if [ -f ./bin/${APP_NAME} ]; then
  echo -e "\nrunning govulncheck over the binary: ./bin/${APP_NAME}"
  govulncheck -mode=binary ./bin/${APP_NAME}
fi
