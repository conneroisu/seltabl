#!/bin/bash
# file: makefile.prod.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.prod.sh
# title: Running Production scripts
# description: This script runs the production generating scripts for the application.
#
# usage: make prod

gum spin --spinner dot --title "Running Tailwind Build" --show-output -- \
	make tailwind-build
gum spin --spinner dot --title "Generating Templ Files" --show-output -- \
	make templ-generate
gum spin --spinner dot --title "Building Application" --show-output -- \
	go build -ldflags "-X main.Environment=production" \
		-o ./bin/$APP_NAME ./main.go
