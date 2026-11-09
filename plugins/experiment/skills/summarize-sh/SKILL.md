---
name: summarize-sh
description: >
    Use when the user names the `summarize` CLI (summarize.sh) or needs what
    only that tool adds — streamed terminal summaries, pipe-friendly text
    extraction, or podcast / audio / YouTube summaries that require
    transcription. For an ordinary "summarize this link or file", use
    content-summarizer. Triggers on: "summarize.sh", "summarize CLI", "用
    summarize 指令", "summarize this podcast", "transcribe and summarize".
version: "1.0.0"
allowed-tools: Bash, Read
user-invocable: true
disable-model-invocation: false
effort: low
metadata:
    type: reference
    platforms: [macos, linux]
    homepage: https://summarize.sh
---

# Summarize

`Source: https://summarize.sh`

`summarize` is a fast CLI and Chrome Side Panel for clean extraction and sharp summaries of web pages, files, YouTube videos, and podcasts.

## Installation

Ensure the CLI is installed globally:

```bash
npm i -g @steipete/summarize
```

## Usage Examples

### Web page summary (streamed Markdown to terminal)

```bash
summarize "https://example.com/article"
```

### YouTube (captions first, yt-dlp + Whisper as fallback)

```bash
summarize "https://youtu.be/I845O57ZSy4"
```

### Local file (PDF, image, audio, video) using specific model

```bash
summarize ./report.pdf --model openai/gpt-5-mini
```

### Extract cleaned text only (skips LLM, perfect for pipes)

```bash
summarize "https://example.com" --extract --format md | wc -w
```

### Clipboard / stdin

```bash
pbpaste | summarize -
```

### JSON envelope for scripts (includes prompt + metrics)

```bash
summarize "https://example.com" --json --metrics detailed
```

## Configuration

Save defaults in `~/.summarize/config.json` and override per-invocation with flags.

```json
{
    "model": "auto",
    "output": {
        "length": "xl",
        "language": "auto"
    },
    "cache": { "enabled": true, "maxMb": 512 }
}
```
