#!/bin/bash

# Read JSON from stdin
json=$(cat)

# Extract file_path using jq
file_path=$(echo "$json" | jq -r '.tool_input.file_path // empty')

# If file_path ends with .go, format it
if [[ -n "$file_path" && "$file_path" == *.go ]]; then
    gofmt -w "$file_path"
fi