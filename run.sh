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
if [ ! -f "$HOME/.config/litellm/litellm_config.yaml" ]; then
    cp "$(pwd)/pkg/litellm/litellm_config.sample.yaml" "$HOME/.config/litellm/litellm_config.yaml"
fi
ln -sf "$HOME/.config/litellm/litellm_config.yaml" "$(pwd)/config/" 


# SKILL: summarize
mkdir -p ~/.summarize/
if [ ! -f ~/.summarize/config.json ]; then
    cp "$(pwd)/pkg/summarize.sh/config.sample.json" ~/.summarize/config.json
fi
ln -sf ~/.summarize ./config/


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
ln -sf "$(pwd)/pkg/tokscale/settings.json" "$HOME/.config/tokscale/settings.json"
ln -sf "$HOME/.config/tokscale" "$(pwd)/config/"
