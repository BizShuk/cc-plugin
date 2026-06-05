"""
# cmd: python -m samples.metadata
# title: Metadata& Token Counting
# description: 中繼資料與 Token 計數範例
"""

import click
from .client import get_client


@click.group()
def metadata():
    """Metadata and token counting examples"""
    pass


@metadata.command()
# cmd: python -m samples.metadata with_metadata
@click.command()
def with_metadata():
    """中繼資料"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[{"role": "user", "content": "Hello!"}],
        metadata={
            "user_id": "user_123",
            "session_id": "sess_abc456",
            "agent_id": "agent_001",
            "tags": ["qa", "test"],
        },
    )
    click.echo(f"Response: {response.content[0].text}")


@metadata.command()
# cmd: python -m samples.metadata count_tokens
@click.command()
def count_tokens():
    """估算輸入 token 數量"""
    client = get_client()
    result = client.messages.count_tokens(
        model="MiniMax-M3",
        messages=[{"role": "user", "content": "Hello, how are you today?"}],
    )
    click.echo(f"Input tokens: {result.input_tokens}")


if __name__ == "__main__":
    metadata()
