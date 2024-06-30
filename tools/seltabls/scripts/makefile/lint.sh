#!/bin/bash
# file: makefile/lint.sh
# url: https://github.com/conneroisu/seltab/tools/seltab-lsp/scripts/makefile/lint.sh
# title: Linting Script
# description: This script runs the linting for the project.
#
# Usage: make js

staticcheck ./...

golangci-lint run

go vet ./...

revive -config .revive.toml ./...
