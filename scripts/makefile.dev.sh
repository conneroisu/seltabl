#!/bin/bash
# file: makefile.dev.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.dev.sh
# title: Running Development Scripts

shopt -s globstar
templ generate --watch --proxy="http://localhost:8080" --cmd="doppler run -- air"
# doppler run -- air
