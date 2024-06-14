#!/bin/bash
# file: taskfile/fmt.sh
# url: https://github.com/conneroisu/seltabl/scripts/taskfile.fmt/sh
# title: Formatting Go Files
# description: This script formats the Go files using gofmt and golines.
#
# Usage: make fmt

gum spin --spinner dot --title "Formatting Go Files" --show-output -- \
    gofmt -w . && golines -w --max-len=79 .
