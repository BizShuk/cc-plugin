"""
# cmd: python -m samples.thinking
# title: Thinking / Reasoning
# description: 思考模式範例 (adaptive, disabled, budget)
"""

import click
from .client import get_client


@click.group()
def thinking():
    """Thinking/reasoning examples"""
    pass


@thinking.command()
# cmd: python -m samples.thinking adaptive
@click.command()
def adaptive():
    """思考模式 - adaptive (自動調整思維預算)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "Solve this riddle: I have cities but no houses."}],
        thinking={"type": "adaptive"},
    )
    for block in response.content:
        if hasattr(block, "thinking"):
            click.echo(f"Thinking: {block.thinking}")
        if hasattr(block, "text"):
            click.echo(f"Text: {block.text}")


@thinking.command()
# cmd: python -m samples.thinking disabled
@click.command()
def disabled():
    """思考模式 - disabled (關閉思考)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "What is 2+2?"}],
        thinking={"type": "disabled"},
    )
    click.echo(f"Response: {response.content[0].text}")


@thinking.command()
# cmd: python -m samples.thinking budget
@click.command()
def budget():
    """思考模式 - 自訂 token 預算"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "Explain the theory of relativity."}],
        thinking={"type": "adaptive", "budget_tokens": 8000},
    )
    for block in response.content:
        if hasattr(block, "thinking"):
            click.echo(f"Thinking: {block.thinking[:100]}...")
        if hasattr(block, "text"):
            click.echo(f"Text: {block.text}")


if __name__ == "__main__":
    thinking()
