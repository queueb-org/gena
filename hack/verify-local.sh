#!/usr/bin/env bash

# Project-wide verification. test.sh owns test execution and exclusions.
set -euo pipefail

if (( $# != 0 )); then
  echo "Usage: $0 (verifies the whole project)" >&2
  exit 2
fi

project_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
cd -- "$project_dir"

./hack/test.sh

awk '
  function fail(message) {
    print "coverage: " message > "/dev/stderr"
    failed = 1
  }
  NR == 1 {
    if ($0 !~ /^mode: (set|count|atomic)$/) fail("invalid profile header")
    next
  }
  {
    if (NF != 3 || $1 !~ /:[0-9]+\.[0-9]+,[0-9]+\.[0-9]+$/ ||
        $2 !~ /^[0-9]+$/ || $3 !~ /^[0-9]+$/) {
      fail("invalid profile record at line " NR)
      next
    }
    statements += $2
    if ($2 > 0 && $3 == 0) fail("uncovered: " $1 " (" $2 " statements)")
    if ($1 ~ /^queueb.org\/gena\/main.go:/) root_statements += $2
  }
  END {
    if (NR == 0 || statements == 0) fail("empty profile or no statements")
    if (root_statements == 0) fail("root main.go is missing from coverage")
    if (failed) exit 1
    printf "coverage gate: 100%% (%d statements)\n", statements
  }
' .coverage

./hack/lint.sh
