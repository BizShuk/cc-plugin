"""
# cmd: python -m samples.tools
# title: Tool Calling
# description: 工具定義與呼叫範例
"""

import click
from .client import get_client


TOOL_DEFINITION = [
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
]


@click.group()
def tools():
    """Tool calling examples"""
    pass


@tools.command()
# cmd: python -m samples.tools basic
@click.command()
def basic():
    """工具定義與呼叫"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "What is the weather in Tokyo?"}],
        tools=TOOL_DEFINITION,
    )
    click.echo(f"Response: {response.content}")


@tools.command()
# cmd: python -m samples.tools auto
@click.command()
def auto():
    """工具選擇策略 - auto"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "What is the weather in Tokyo?"}],
        tools=TOOL_DEFINITION,
        tool_choice={"type": "auto"},
    )
    click.echo(f"Response: {response.content}")


@tools.command()
# cmd: python -m samples.tools any
@click.command()
def any():
    """工具選擇策略 - any (強制使用工具)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "What is the weather in Tokyo?"}],
        tools=TOOL_DEFINITION,
        tool_choice={"type": "any"},
    )
    click.echo(f"Response: {response.content}")


@tools.command()
# cmd: python -m samples.tools forced
@click.command()
def forced():
    """工具選擇策略 - forced (指定工具)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "What is the weather in Tokyo?"}],
        tools=TOOL_DEFINITION,
        tool_choice={"type": "forced", "name": "get_weather"},
    )
    click.echo(f"Response: {response.content}")


@tools.command()
# cmd: python -m samples.tools result
@click.command()
def result():
    """工具結果回傳 (tool result)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[
            {"role": "user", "content": "What is the weather in Tokyo?"},
            {
                "role": "assistant",
                "content": [
                    {
                        "type": "tool_use",
                        "id": "toolu_01HY5V",
                        "name": "get_weather",
                        "input": {"location": "Tokyo", "unit": "celsius"},
                    }
                ],
            },
            {
                "role": "user",
                "content": [
                    {
                        "type": "tool_result",
                        "tool_use_id": "toolu_01HY5V",
                        "content": '{"temperature": 22, "condition": "sunny", "humidity": 65}',
                        "is_error": False,
                    }
                ],
            },
        ],
        tools=TOOL_DEFINITION,
    )
    click.echo(f"Response: {response.content[0].text}")


if __name__ == "__main__":
    tools()
