#!/bin/bash
# file: makefile.test.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.test.sh
# title: Test Script
# description: This script runs the test for the project.
# 
# usage: make test

gum spin --spinner dot --title "Running Go Test With Race" --show-output -- \
    go test -race -v -timeout 30s ./...
# gum spin --spinner dot --title "Running Go Test With Coverage" --show-output -- \
    doppler run -- go test -coverprofile=coverage.out ./... 
# gum spin --spinner dot --title "Running Make Lint" --show-output -- \
    # make lint
gocovsh
