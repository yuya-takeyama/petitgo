#!/bin/bash

# Check if file is a .pg or .sh file
if [[ "$1" == *.pg ]] || [[ "$1" == *.sh ]]; then
    echo "Checking newlines for file: $1"
    python3 tools/check_newlines.py "$1"
else
    echo "File $1 is not a .pg or .sh file, skipping newline check"
fi
