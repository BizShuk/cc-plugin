#!/bin/bash


ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.claude/CLAUDE.md" 


ln -sf "$HOME/.claude"        "./config/"
ln -sf "$HOME/.claude-mem"    "./config/"
ln -sf "$HOME/.claude.json"   "./config/"