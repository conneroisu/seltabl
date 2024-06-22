#!/bin/bash
# Name: makefile.docs.sh
# Url: 
# https://github.com/conneroisu/seltabl/main/scripts/makefile.docs.sh
# 
# Description: A script to generate the go docs for the project.
# 
# Usage: make docs

mkdir docs
golds -s -gen -wdpkgs-listing=promoted -dir=./docs -footer=verbose+qrcode
xdg-open ./docs/index.html
