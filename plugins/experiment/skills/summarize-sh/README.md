# summarize.sh

`Source: https://summarize.sh`

`summarize` is a fast CLI and Chrome Side Panel for clean extraction and sharp summaries of web pages, files, YouTube videos, and podcasts.

## Installation

```bash
npm i -g @steipete/summarize
```

## Features

- `Real extraction`: Readability for articles, `markitdown` for files, Firecrawl as a fallback when sites fight back.
- `Media-aware`: YouTube and podcast pages prefer published transcripts, then `yt-dlp` + Whisper, then optional ONNX models (Parakeet/Canary).
- `Provider-agnostic models`: xAI, OpenAI, Google, Anthropic, NVIDIA, Z.AI, OpenRouter, GitHub Copilot, Ollama (local) — plus local CLI providers.
- `Shaped output`: Streamed ANSI Markdown for terminals, plain text for pipes, JSON envelope for scripts, ANSI-stripped under `--no-color`.
- `Slides for video`: `--slides` extracts scene-change keyframes and renders them inline or saves them to disk.
- `Stays local`: Optional daemon + Chrome Side Panel pair the CLI with the active tab.

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
