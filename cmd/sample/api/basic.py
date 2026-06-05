"""
# cmd: python -m samples.basic
# title: Basic Requests
# description: 最小請求、系統提示詞、溫度、Top-P 等基本範例
"""

import click
from .client import get_client


@click.group()
def basic():
    """Basic request examples"""
    pass


@basic.command()
# cmd: python -m samples.basic minimal
@click.command()
def minimal():
    """最小請求 —只需要 model, messages, max_tokens"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "Hello, how are you?"}],
    )
    click.echo(f"Response: {response.content[0].text}")


@basic.command()
# cmd: python -m samples.basic system
@click.command()
def system():
    """系統提示詞"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        system="You are a helpful assistant specialized in geography.",
        messages=[{"role": "user", "content": "What is the capital of France?"}],
    )
    click.echo(f"Response: {response.content[0].text}")


@basic.command()
# cmd: python -m samples.basic temperature
@click.command()
def temperature():
    """溫度參數 (0.0 - 2.0)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        temperature=0.7,
        messages=[{"role": "user", "content": "Write a short poem about the sea."}],
    )
    click.echo(f"Response: {response.content[0].text}")


@basic.command()
# cmd: python -m samples.basic top_p
@click.command()
def top_p():
    """Top-P 核抽樣 (0.0 - 1.0)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        top_p=0.95,
        messages=[{"role": "user", "content": "Explain quantum physics in simple terms."}],
    )
    click.echo(f"Response: {response.content[0].text}")


if __name__ == "__main__":
    basic()
