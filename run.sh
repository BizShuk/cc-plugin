#!/bin/bash

# Ensure home directories exist
mkdir -p "$HOME/.claude"
mkdir -p "$HOME/.codex"
mkdir -p "$HOME/.gemini"
mkdir -p "$HOME/.claude-mem"
mkdir -p "$HOME/.hermes"

# Create local directories
mkdir -p config
mkdir -p logs

# Claude
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.claude/CLAUDE.md" 
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.gemini/GEMINI.md" 
ln -sf "$(pwd)/config/CLAUDE.global.md" "$HOME/.codex/AGENTS.md" 
ln -sf "$(pwd)/config/settings.json"    "$HOME/.claude/settings.json" 
ln -sf "$(pwd)/config/config.toml"    "$HOME/.codex/config.toml" 

# Global plugin configuration links back to local project (only link if the target actually exists)
if [ -d "$HOME/.gemini" ]; then
    ln -sf "$HOME/.gemini"        "./tmp/"
fi
if [ -d "$HOME/.claude" ]; then
    ln -sf "$HOME/.claude"        "./tmp/"
fi
if [ -d "$HOME/.codex" ]; then
    ln -sf "$HOME/.codex"        "./tmp/"
fi
if [ -d "$HOME/.claude-mem" ]; then
    ln -sf "$HOME/.claude-mem"    "./tmp/"
fi
if [ -f "$HOME/.claude.json" ]; then
    ln -sf "$HOME/.claude.json"   "./tmp/"
fi
if [ -d "$HOME/.cli-proxy-api" ]; then
    ln -sf "$HOME/.cli-proxy-api"   "./tmp/"
fi
if [ -d "$HOME/.auth2api" ]; then
    ln -sf "$HOME/.auth2api"   "./tmp/"
fi
if [ -d "$HOME/.hermes" ]; then
    ln -sf "$HOME/.hermes"        "./tmp/"
fi
if [ -d "$HOME/.gbrain" ]; then
    ln -sf "$HOME/.gbrain"        "./tmp/"
fi

if [ -d "$HOME/.mempalace" ]; then
    ln -sf "$HOME/.mempalace"        "./tmp/"
fi

if [ -d "$HOME/.config/cc-plugin" ]; then
    ln -sf "$HOME/.config/cc-plugin"        "./tmp/"
fi

if [ -d "$HOME/.agentmemory" ]; then
    ln -sf "$HOME/.agentmemory"        "./tmp/"
fi

if [ -d "$HOME/.paperclip" ]; then
    ln -sf "$HOME/.paperclip"        "./tmp/"
fi

if [ -d "$HOME/.codex" ]; then
    ln -sf "$HOME/.codex"        "./tmp/"
fi

if [ -d "$HOME/.config/opencode" ]; then
    ln -sf "$HOME/.config/opencode"        "./tmp/"
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

ln -sf "$HOME/.config/litellm/litellm_config.yaml" "./tmp/" 

# SKILL: summarize
mkdir -p "$HOME/.summarize"
if [ ! -f "$HOME/.summarize/config.json" ]; then
    cp "$(pwd)/pkg/summarize.sh/config.sample.json" "$HOME/.summarize/config.json"
fi
ln -sf "$HOME/.summarize" "./tmp/"

# CCStatusline
mkdir -p "$HOME/.config/ccstatusline"
if [ -f "$HOME/.config/ccstatusline/settings.json" ] && [ ! -L "$HOME/.config/ccstatusline/settings.json" ]; then
    mv "$HOME/.config/ccstatusline/settings.json" "$HOME/.config/ccstatusline/settings.json.bak"
fi
ln -sf "$(pwd)/pkg/ccstatusline/settings.json" "$HOME/.config/ccstatusline/settings.json"

ln -sf "$HOME/.config/ccstatusline" "./tmp/"

# Tokscale
mkdir -p "$HOME/.config/tokscale"
if [ -f "$HOME/.config/tokscale/settings.json" ] && [ ! -L "$HOME/.config/tokscale/settings.json" ]; then
    mv "$HOME/.config/tokscale/settings.json" "$HOME/.config/tokscale/settings.json.bak"
fi
ln -sf "$(pwd)/pkg/usage/tokscale/settings.json" "$HOME/.config/tokscale/settings.json"
ln -sf "$HOME/.config/tokscale" "./tmp/"


