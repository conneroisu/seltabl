#!/bin/bash
# file: makefile.coverage.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.coverage.sh
# title: Coverage Script
# description: This script runs the coverage testing for the project.

go test -coverprofile=coverage.out ./...

gocovsh
