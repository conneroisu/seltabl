#!/bin/bash
# file: makefile.test.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.test.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

go test -race -timeout 30s ./...

go test -coverprofile=coverage.out ./...
