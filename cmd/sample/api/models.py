"""
# cmd: python -m samples.models
# title: Supported Models
# description: 支援的模型列表
"""

import click


@click.command()
# cmd: python -m samples.models
def models():
    """支援的模型列表"""
    model_list = [
        "MiniMax-M3",
        "MiniMax-M2.7",
        "MiniMax-M2.7-highspeed",
        "MiniMax-M2.5",
        "MiniMax-M2.5-highspeed",
        "MiniMax-M2.1",
        "MiniMax-M2.1-highspeed",
        "MiniMax-M2",
    ]
    for model in model_list:
        click.echo(f"  - {model}")


if __name__ == "__main__":
    models()
