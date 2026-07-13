#!/bin/bash
# Trigger promotion of a knowledge entry
# Usage: ./promote.sh <item_id> [target_layer]
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -z "$1" ]; then
    echo "Usage: $0 <item_id> [target_layer]"
    echo "Example: $0 HY-001 principle"
    exit 1
fi

python3 "$SCRIPT_DIR/kb_manager.py" promote "$1" "$2"
