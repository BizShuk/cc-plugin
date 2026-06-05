"""
# cmd: python -m samples.error
# title: Error Handling
# description: 錯誤處理範例
"""

import click
from .client import get_client


@click.command()
# cmd: python -m samples.error
def error():
    """錯誤處理範例"""
    client = get_client()

    try:
        # 故意省略 required 欄位
        response = client.messages.create(
            model="MiniMax-M3",
            # max_tokens=1024,  # 缺少 max_tokens
            messages=[{"role": "user", "content": "Hello!"}],
        )
    except Exception as e:
        click.echo(f"Error Type: {type(e).__name__}")
        click.echo(f"Error Message: {e}")


if __name__ == "__main__":
    error()
