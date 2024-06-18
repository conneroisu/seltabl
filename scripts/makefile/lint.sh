#!/bin/bash
# file: makefile.js.sh
# url: https://github.com/conneroisu/seltab/tools/seltab-lsp/scripts/makefile.js.sh
# title: Running Webpack
# description: This script runs Webpack to build the JavaScript files.
#
# Usage: make js

staticcheck ./...

golangci-lint run

go vet ./...

revive -config .revive.toml ./...
