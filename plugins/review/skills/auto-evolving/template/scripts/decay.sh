#!/bin/bash
# Trigger confidence decay for registered knowledge entries
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
python3 "$SCRIPT_DIR/kb_manager.py" decay
