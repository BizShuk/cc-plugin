"""
# cmd: python -m samples.all
# title: Run All Samples
# description: 執行所有範例
"""

import click

from . import client, basic, messages, tools, thinking, streaming, metadata, response, models, error, full_request


@click.command()
# cmd: python -m samples.all
def all():
    """執行所有範例"""
    click.echo("=" * 60)
    click.echo("MiniMax Anthropic API - Python SDK Examples")
    click.echo("=" * 60)

    click.echo("\n[1] Client Initialization")
    client.get_client()
    click.echo("  OK")

    click.echo("\n[2] Basic Requests")
    basic.minimal()
    basic.system()
    basic.temperature()
    basic.top_p()

    click.echo("\n[3] Message Content Types")
    messages.text()
    messages.image_url()
    messages.image_base64()
    messages.video_url()
    messages.video_mmfile()

    click.echo("\n[4] Tool Calling")
    tools.basic()
    tools.auto()
    tools.any()
    tools.forced()
    tools.result()

    click.echo("\n[5] Thinking / Reasoning")
    thinking.adaptive()
    thinking.disabled()
    thinking.budget()

    click.echo("\n[6] Streaming")
    streaming.streaming()

    click.echo("\n[7] Metadata& Token Counting")
    metadata.with_metadata()
    metadata.count_tokens()

    click.echo("\n[8] Response Attributes")
    response.response()

    click.echo("\n[9] Supported Models")
    models.models()

    click.echo("\n[10] Error Handling")
    error.error()

    click.echo("\n[11] Full Request")
    full_request.full_request()

    click.echo("\n" + "=" * 60)
    click.echo("All examples completed!")
    click.echo("=" * 60)


if __name__ == "__main__":
    all()
