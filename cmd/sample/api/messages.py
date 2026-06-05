"""
# cmd: python -m samples.messages
# title: Message Content Types
# description: 文字、圖片、影片內容區塊範例
"""

import click
from .client import get_client


@click.group()
def messages():
    """Message content type examples"""
    pass


@messages.command()
# cmd: python -m samples.messages text
@click.command()
def text():
    """文字內容區塊 (text content block)"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "What is the capital of France?"},
                    {"type": "text", "text": "What about Germany?"},
                ],
            }
        ],
    )
    click.echo(f"Response: {response.content[0].text}")


@messages.command()
# cmd: python -m samples.messages image_url
@click.command()
def image_url():
    """
    圖片內容區塊 (M3 only)
    - 支援格式: JPEG/PNG/GIF/WEBP
    - 最大10 MB (512 MB via Files API)
    """
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "What is in this image?"},
                    {
                        "type": "image",
                        "source": {
                            "type": "url",
                            "media_type": "image/jpeg",
                            "url": "https://example.com/photo.jpg",
                        },
                    },
                ],
            }
        ],
    )
    click.echo(f"Response: {response.content[0].text}")


@messages.command()
# cmd: python -m samples.messages image_base64
@click.command()
def image_base64():
    """圖片 base64 編碼"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "Describe this image."},
                    {
                        "type": "image",
                        "source": {
                            "type": "base64",
                            "media_type": "image/png",
                            "data": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
                        },
                    },
                ],
            }
        ],
    )
    click.echo(f"Response: {response.content[0].text}")


@messages.command()
# cmd: python -m samples.messages video_url
@click.command()
def video_url():
    """
    影片內容區塊 (M3 only)
    - 最大50 MB (512 MB via Files API)
    """
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "What happens in this video?"},
                    {
                        "type": "video",
                        "source": {
                            "type": "url",
                            "media_type": "video/mp4",
                            "url": "https://example.com/video.mp4",
                        },
                    },
                ],
            }
        ],
    )
    click.echo(f"Response: {response.content[0].text}")


@messages.command()
# cmd: python -m samples.messages video_mmfile
@click.command()
def video_mmfile():
    """使用 mm_file://引用已上傳的檔案"""
    client = get_client()
    response = client.messages.create(
        model="MiniMax-M3",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "What happens in this video?"},
                    {
                        "type": "video",
                        "source": {
                            "type": "mm_file",
                            "media_type": "video/mp4",
                            "data": "file_id_12345",
                        },
                    },
                ],
            }
        ],
    )
    click.echo(f"Response: {response.content[0].text}")


if __name__ == "__main__":
    messages()
