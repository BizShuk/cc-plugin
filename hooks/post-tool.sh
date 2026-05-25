#!/bin/bash
# Read the JSON from stdin
PAYLOAD=$(cat)

# Extract tool name and target file path
TOOL_NAME=$(echo "$PAYLOAD" | jq -r '.tool_name // empty')
TARGET_FILE=$(echo "$PAYLOAD" | jq -r '.tool_input.TargetFile // .tool_input.AbsolutePath // .tool_input.path // empty')

# If target file is empty, try to search for any value ending in .go in the tool_input
if [ -z "$TARGET_FILE" ]; then
  TARGET_FILE=$(echo "$PAYLOAD" | jq -r '.tool_input | to_entries[] | select(.value | type == "string") | .value' | grep -E '\.go$' | head -n 1)
fi

# If we have a target file and it ends with .go and exists
if [ -n "$TARGET_FILE" ] && [[ "$TARGET_FILE" == *.go ]] && [ -f "$TARGET_FILE" ]; then
  # Run go fmt
  go fmt "$TARGET_FILE" >/dev/null 2>&1
  
  # If golangci-lint is installed, run it
  if command -v golangci-lint &> /dev/null; then
    golangci-lint run --fix "$TARGET_FILE" >/dev/null 2>&1
  fi
fi

exit 0
