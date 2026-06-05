"""
# cmd: python -m samples.streaming
# title: Streaming Response
# description: 串流回應範例
"""

import click
from .client import get_client


@click.command()
# cmd: python -m samples.streaming
def streaming():
    """串流回應"""
    client = get_client()

    with client.messages.stream(
        model="MiniMax-M3",
        max_tokens=512,
        messages=[{"role": "user", "content": "Write a haiku about the ocean."}],
    ) as stream:
        for event in stream:
            if event.type == "content_block_delta":
                click.echo(event.delta.text, nl=False)
    click.echo()


if __name__ == "__main__":
    streaming()
