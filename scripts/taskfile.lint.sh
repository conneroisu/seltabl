#!/bin/bash
# file: taskfile.test.sh
# url: https://github.com/conneroisu/seltabl/scripts/taskfile.test.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

staticcheck ./...

golangci-lint run

go vet ./...

revive -config .revive.toml ./...
