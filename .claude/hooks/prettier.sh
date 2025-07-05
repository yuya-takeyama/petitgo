#!/bin/bash

# Check if npx and prettier are available
if ! command -v npx &> /dev/null; then
    echo "npx not found, skipping prettier formatting"
    exit 0
fi

# Check if file is a markdown, JSON, YAML, or YML file
if [[ "$1" == *.md ]] || [[ "$1" == *.json ]] || [[ "$1" == *.yaml ]] || [[ "$1" == *.yml ]]; then
    echo "Running prettier on $1"
    npx prettier --write "$1"
else
    echo "File $1 is not a supported file type, skipping prettier"
fi
