"""
# cmd: python -m samples.response
# title: Response Attributes
# description: 完整回應屬性範例
"""

import click
from .client import get_client


@click.command()
# cmd: python -m samples.response
def response():
    """完整回應屬性"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "What is the capital of Japan?"}],
    )

    # 基本欄位
    click.echo(f"ID: {response.id}")
    click.echo(f"Type: {response.type}")
    click.echo(f"Role: {response.role}")
    click.echo(f"Model: {response.model}")
    click.echo(f"Stop Reason: {response.stop_reason}")

    # Content blocks
    for block in response.content:
        click.echo(f"Block Type: {block.type}")
        if hasattr(block, "text"):
            click.echo(f"Text: {block.text}")
        if hasattr(block, "thinking"):
            click.echo(f"Thinking: {block.thinking}")
        if hasattr(block, "name"):
            click.echo(f"Tool Name: {block.name}")
            click.echo(f"Tool Input: {block.input}")

    # Usage
    click.echo(f"Input Tokens: {response.usage.input_tokens}")
    click.echo(f"Output Tokens: {response.usage.output_tokens}")
    click.echo(f"Total Tokens: {response.usage.total_tokens}")

    # Metrics (if available)
    if hasattr(response, "metrics") and response.metrics:
        click.echo(f"Latency (ms): {response.metrics.latency_ms}")
        click.echo(f"Tokens/Second: {response.metrics.tokens_per_second}")


if __name__ == "__main__":
    response()
