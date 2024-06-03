#!/bin/bash
# file: makefile.fmt.sh
# url: https://github.com/conneroisu/seltabl/scripts/makefile.fmt.sh
# title: Formatting Go Files
# description: This script formats the Go files using gofmt and golines.
#
# Usage: make fmt


targets=(
	"*.go"
	"**/*.go"
	"**/**/*.go"
	"**/**/**/*.go"
	"**/**/**/**/*.go"
	"**/**/**/**/**/*.go"
)

# For each of the targets, run gofmt and goline.
for target in "${targets[@]}"; do
	if ls $target &>/dev/null; then
		if ! command -v gum &>/dev/null; then
			echo "formatting $target with gofmt"
			gofmt -w $target
			echo "formatting $target with golines"
			goline -w --max-len=79 $target
		else
			gum spin --spinner dot --title "Formatting Go Files with 'go fmt' in $target" --show-output -- \
				go fmt $target
			gum spin --spinner dot --title "Formatting Go Files with 'golines' in $target" --show-output -- \
				golines -w --max-len=79 $target
		fi
	else
		continue
	fi
done
# gum spin --spinner dot --title "Formatting Go Files" --show-output -- \
#         go fmt *.go
# # gum spin --spinner dot --title "Formatting Go Files" --show-output -- \
# # golines --max-len=80 *.go **/*.go **/**/*.go --dry-run >
# OUTPUT=$(golines --max-len=80 *.go **/*.go **/**/*.go --dry-run)
# gum pager < <(echo "$OUTPUT") && gum confirm "Commit changes?" && gum spin --spinner dot --title "Formatting Go Files" --show-output -- \
#         golines -w --max-len=80 *.go **/*.go **/**/*.go
#
# OUTPUT=$(caddy fmt --diff ./deploy/**/**/Caddyfile)
gum pager < <(echo "$OUTPUT") && gum confirm "Commit changes?" && gum spin --spinner dot --title "Formatting Caddyfiles Files" --show-output -- \
	caddy fmt ./deploy/**/**/Caddyfile --overwrite
