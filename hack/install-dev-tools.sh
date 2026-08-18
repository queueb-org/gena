#!/usr/bin/env bash

set -euo pipefail

if (( $# != 0 )); then
  echo "Usage: $0" >&2
  exit 2
fi

project_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
tools_dir="$project_dir/bin/tools"
mkdir -p -- "$tools_dir"

# Versions match the tools validated locally for the initial release.
# Explicit versions keep tool dependencies separate from the application module.
dev_tools=(
  github.com/fzipp/gocyclo/cmd/gocyclo@latest
  github.com/gordonklaus/ineffassign@latest
  honnef.co/go/tools/cmd/staticcheck@latest
  golang.org/x/vuln/cmd/govulncheck@latest
)

for tool in "${dev_tools[@]}"; do
  printf 'Installing %s\n' "$tool"
  GOBIN="$tools_dir" go install "$tool"
done

printf 'Development tools installed in %s\n' "$tools_dir"
