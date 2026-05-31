#!/bin/bash

# Ensure home directories exist
mkdir -p "$HOME/.claude"
mkdir -p "$HOME/.gemini"
mkdir -p "$HOME/.claude-mem"
mkdir -p "$HOME/.hermes"

# Create local directories
mkdir -p config
mkdir -p logs

# Claude
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.claude/CLAUDE.md" 
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.gemini/GEMINI.md" 
ln -sf "$(pwd)/config/settings.json"    "$HOME/.claude/settings.json" 

# Global plugin configuration links back to local project (only link if the target actually exists)
if [ -d "$HOME/.gemini" ]; then
    ln -sf "$HOME/.gemini"        "./config/"
fi
if [ -d "$HOME/.claude" ]; then
    ln -sf "$HOME/.claude"        "./config/"
fi
if [ -d "$HOME/.claude-mem" ]; then
    ln -sf "$HOME/.claude-mem"    "./config/"
fi
if [ -f "$HOME/.claude.json" ]; then
    ln -sf "$HOME/.claude.json"   "./config/"
fi
if [ -d "$HOME/.hermes" ]; then
    ln -sf "$HOME/.hermes"        "./config/"
fi
if [ -d "$HOME/.gbrain" ]; then
    ln -sf "$HOME/.gbrain"        "./config/"
fi

if [ -d "$HOME/.mempalace" ]; then
    ln -sf "$HOME/.mempalace"        "./config/"
fi




# Hermes
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.hermes/AGENTS.md" 
ln -sf "$(pwd)/pkg/hermes/MEMORY.md" "$HOME/.hermes/MEMORY.md"
ln -sf "$(pwd)/pkg/hermes/USER.md" "$HOME/.hermes/USER.md"

# LiteLLM
mkdir -p "$HOME/.config/litellm"
if [ ! -f "$HOME/.config/litellm/litellm_config.yaml" ]; then
    cp "$(pwd)/pkg/litellm/litellm_config.sample.yaml" "$HOME/.config/litellm/litellm_config.yaml"
fi
ln -sf "$HOME/.config/litellm/litellm_config.yaml" "$(pwd)/config/" 

# SKILL: summarize
mkdir -p "$HOME/.summarize"
if [ ! -f "$HOME/.summarize/config.json" ]; then
    cp "$(pwd)/pkg/summarize.sh/config.sample.json" "$HOME/.summarize/config.json"
fi
ln -sf "$HOME/.summarize" "./config/"

# CCStatusline
mkdir -p "$HOME/.config/ccstatusline"
if [ -f "$HOME/.config/ccstatusline/settings.json" ] && [ ! -L "$HOME/.config/ccstatusline/settings.json" ]; then
    mv "$HOME/.config/ccstatusline/settings.json" "$HOME/.config/ccstatusline/settings.json.bak"
fi
ln -sf "$(pwd)/pkg/ccstatusline/settings.json" "$HOME/.config/ccstatusline/settings.json"
ln -sf "$HOME/.config/ccstatusline" "$(pwd)/config/"

# Tokscale
mkdir -p "$HOME/.config/tokscale"
if [ -f "$HOME/.config/tokscale/settings.json" ] && [ ! -L "$HOME/.config/tokscale/settings.json" ]; then
    mv "$HOME/.config/tokscale/settings.json" "$HOME/.config/tokscale/settings.json.bak"
fi
ln -sf "$(pwd)/pkg/usage/tokscale/settings.json" "$HOME/.config/tokscale/settings.json"
ln -sf "$HOME/.config/tokscale" "$(pwd)/config/"


