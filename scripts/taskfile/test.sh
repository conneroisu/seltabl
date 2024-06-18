#!/bin/bash
# file: taskfile.test.sh
# url: https://github.com/conneroisu/blob/main/seltab/tools/seltab-lsp/scripts/taskfile.test.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

gum spin --spinner dot --title "Running Tests" --show-output -- \
    go test -race -timeout 30s ./...

gum spin --spinner dot --title "Generating Coverage" --show-output -- \
    go test -coverprofile=coverage.out ./...

gocovsh coverage.out
