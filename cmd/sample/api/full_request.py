"""
# cmd: python -m samples.full_request
# title: Full Request Example
# description: 完整請求範例 (包含所有欄位)
"""

import click
from .client import get_client


@click.command()
# cmd: python -m samples.full_request
def full_request():
    """完整請求範例 (以 dict 格式)"""
    client = get_client()

    request_dict = {
        "model": "MiniMax-M3",
        "messages": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "What is the capital of France?"},
                ],
            },
            {
                "role": "assistant",
                "content": [
                    {"type": "text", "text": "The capital of France is Paris."},
                ],
            },
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "Tell me more about it."},
                ],
            },
        ],
        "max_tokens": 1024,
        "stream": False,
        "system": "You are a helpful assistant.",
        "temperature": 0.7,
        "top_p": 0.95,
        "tools": [
            {
                "name": "get_weather",
                "description": "Get current weather for a location",
                "input_schema": {
                    "type": "object",
                    "properties": {
                        "location": {"type": "string", "description": "City name"},
                        "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]},
                    },
                    "required": ["location"],
                },
            }
        ],
        "tool_choice": {"type": "auto"},
        "metadata": {
            "user_id": "user_123",
            "session_id": "sess_abc456",
            "agent_id": "agent_001",
            "tags": ["qa", "geography"],
        },
        "thinking": {"type": "adaptive", "budget_tokens": 8000},
    }

    response = client.messages.create(**request_dict)
    click.echo(f"Response: {response.content[0].text}")


if __name__ == "__main__":
    full_request()
