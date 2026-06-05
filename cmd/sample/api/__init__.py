"""
MiniMax Anthropic API - Python SDK Samples
============================================
samples package - individual example modules

Usage:
    # Set environment
    export ANTHROPIC_BASE_URL=https://api.minimax.io/anthropic
    export ANTHROPIC_API_KEY=${YOUR_API_KEY}

    # Run all samples
    python -m samples.all

    # Run specific sample
    python -m samples.basic.minimal
"""

__all__ = [
    "client",
    "basic",
    "messages",
    "tools",
    "thinking",
    "streaming",
    "metadata",
    "response",
    "models",
    "error",
    "full_request",
]
