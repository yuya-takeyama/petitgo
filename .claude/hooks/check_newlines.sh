#!/bin/bash

# Check if file is a .pg file
if [[ "$1" == *.pg ]]; then
    echo "Checking newlines for petitgo file: $1"
    python3 tools/check_newlines.py "$1"
else
    echo "File $1 is not a .pg file, skipping newline check"
fi