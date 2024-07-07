#!/bin/bash
# file: taskfile.clean.sh
# url: https://github.com/conneroisu/seltab/tools/seltab-lsp/scripts/taskfile.clean.sh
# title: Cleaning Script
# description: This script cleans the project

task install

# if there is a tmp folder, delete it
if [ -d "tmp" ]; then
    rm -rf tmp
fi

# if there is a bin folder, delete it
if [ -d "bin" ]; then
    rm -rf bin
fi

# if there is a node_modules folder, delete it
if [ -d "node_modules" ]; then
    rm -rf node_modules
fi

# if there is a node_modules in a subfolder, delete it
if [ -d "data/javascript/node_modules" ]; then
    rm -rf data/javascript/node_modules
fi

# if there is a coverage.out file, delete it
if [ -f "coverage.out" ]; then
    rm -rf coverage.out
fi

# ask to remove the ~/.config/seltabls/state.log file
if [ -f "$HOME/.config/seltabls/state.log" ]; then
    echo "Do you want to remove the ~/.config/seltabls/state.log file? (y/n)"
    read -r answer
    if [ "$answer" = "y" ]; then
        rm -rf "$HOME/.config/seltabls/state.log"
    fi
fi

# ask to remove the ~/.config/seltabls/uri.sqlite file
if [ -f "$HOME/.config/seltabls/uri.sqlite" ]; then
    echo "Do you want to remove the ~/.config/seltabls/uri.sqlite file? (y/n)"
    read -r answer
    if [ "$answer" = "y" ]; then
        rm -rf "$HOME/.config/seltabls/uri.sqlite"
    fi
fi
