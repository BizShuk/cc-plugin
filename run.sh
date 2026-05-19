#!/bin/bash

# Claude
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.claude/CLAUDE.md" 

# Claude plugin
ln -sf "$HOME/.claude"        "./config/"
ln -sf "$HOME/.claude-mem"    "./config/"
ln -sf "$HOME/.claude.json"   "./config/"

# Hermes
ln -sf "$HOME/.hermes"   "./config/"
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.hermes/AGENTS.md" 
ln -sf "$(pwd)/pkg/hermes/MEMORY.md" $HOME/.hermes/
ln -sf "$(pwd)/pkg/hermes/USER.md" $HOME/.hermes/


