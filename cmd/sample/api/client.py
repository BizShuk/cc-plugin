"""
# cmd: python -m samples.client
# title: Client Initialization
# description: 建立 MiniMax API 用戶端
"""

import os
from anthropic import Anthropic

def get_client() -> Anthropic:
    """建立 MiniMax API 用戶端"""
    return Anthropic(
        base_url="https://api.minimax.io/anthropic",
        api_key=os.environ.get("ANTHROPIC_API_KEY", ""),
    )

if __name__ == "__main__":
    client = get_client()
    print(f"Client initialized: {client}")
    print(f"Base URL: {client.base_url}")
