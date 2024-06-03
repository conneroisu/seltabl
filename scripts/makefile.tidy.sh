#!/bin/bash
gum spin --spinner dot --title "Running Go Mod Tidy" --show-output -- \
	go mod tidy
