#!/bin/bash

# Claude
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.claude/CLAUDE.md" 
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.gemini/GEMINI.md" 
ln -sf "$(pwd)/config/settings.json"    "$HOME/.claude/settings.json" 

# Global plugin
ln -sf "$HOME/.gemini"        "./config/"
ln -sf "$HOME/.claude"        "./config/"
ln -sf "$HOME/.claude-mem"    "./config/"
ln -sf "$HOME/.claude.json"   "./config/"

# Hermes
ln -sf "$HOME/.hermes"   "./config/"
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.hermes/AGENTS.md" 
ln -sf "$(pwd)/pkg/hermes/MEMORY.md" $HOME/.hermes/
ln -sf "$(pwd)/pkg/hermes/USER.md" $HOME/.hermes/


# LiteLLM
mkdir -p "$HOME/.config/litellm"
ln -sf "$HOME/.config/litellm/litellm_config.yaml" "$(pwd)/config/" 


# SKILL: summarize
mkdir -p ~/.summarize/
cp "$(pwd)/pkg/summarize.sh/config.sample.json" ~/.summarize/config.json
ln -s ~/.summarize ./config/
