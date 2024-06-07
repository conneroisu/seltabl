#!/bin/bash
# file: makefile.fmt.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.fmt.sh
# title: Formatting Go Files
# description: This script formats the Go files using gofmt and golines.
#
# Usage: make fmt

gofmt -w .

golines -w --max-len=79 .
