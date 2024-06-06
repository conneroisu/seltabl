#!/bin/bash
# file: makefile.js.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.js.sh
# title: Running Webpack
# description: This script runs Webpack to build the JavaScript files.
#
# Usage: make js

gum spin --spinner dot --title "Running Static Check" --show-output -- \
	staticcheck ./...

gum spin --spinner dot --title "Running GolangCI Lint" --show-output -- \
	golangci-lint run

gum spin --spinner dot --title "Running GoVet" --show-output -- \
	go vet ./...

gum spin --spinner dot --title "Running Revive" --show-output -- \
	revive -config revive.toml ./...
