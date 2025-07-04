#!/bin/bash

# Check if npx and prettier are available
if ! command -v npx &> /dev/null; then
    echo "npx not found, skipping prettier formatting"
    exit 0
fi

# Check if file is a markdown file
if [[ "$1" == *.md ]]; then
    echo "Running prettier on $1"
    npx prettier --write "$1"
else
    echo "File $1 is not a markdown file, skipping prettier"
fi