---
name: markitdown
description: >
    Use when converting files or URLs to Markdown for LLM consumption. Handles
    PDF, DOCX, PPTX, XLSX, HTML, CSV, JSON, XML, images, audio, YouTube,
    EPUB, ZIP, and Outlook messages. Triggers on: "convert to markdown",
    "extract text from", "read this PDF", "parse this document", or any
    file-to-text conversion need.
---

# MarkItDown

Python CLI and library (by Microsoft) that converts files and URLs to clean
Markdown. Optimized for LLM pipelines — preserves headings, lists, tables,
and links while staying token-efficient.

Source: https://github.com/microsoft/markitdown

## What It Works For

| Category  | Formats                                                  |
| --------- | -------------------------------------------------------- |
| Web       | URLs (fetches HTML → Markdown)                           |
| Documents | PDF, DOCX, PPTX, XLSX, XLS, EPUB                        |
| Data      | CSV, JSON, XML                                           |
| Media     | Images (EXIF + OCR), Audio (EXIF + speech transcription) |
| Video     | YouTube URLs (transcript extraction)                     |
| Archives  | ZIP (iterates over contents)                             |
| Other     | Outlook `.msg`, plain HTML files                         |

Not all formats are enabled by default — some require optional extras
(see Installation).

## Installation

Requires `Python 3.10+`.

### Check if installed

```bash
markitdown --version
```

### Install (all formats)

```bash
pip install 'markitdown[all]'
```

### Install (selective extras)

```bash
pip install 'markitdown[pdf,docx,pptx,xlsx]'
```

Available extras:

| Extra                      | What it enables                       |
| -------------------------- | ------------------------------------- |
| `[all]`                    | Everything below                      |
| `[pdf]`                    | PDF files                             |
| `[docx]`                   | Word documents                        |
| `[pptx]`                   | PowerPoint presentations              |
| `[xlsx]`                   | Excel spreadsheets                    |
| `[xls]`                    | Older Excel files                     |
| `[outlook]`                | Outlook `.msg` messages               |
| `[audio-transcription]`    | WAV / MP3 speech transcription        |
| `[youtube-transcription]`  | YouTube video transcript extraction   |
| `[az-doc-intel]`           | Azure Document Intelligence           |
| `[az-content-understanding]` | Azure Content Understanding         |

### Plugins

3rd-party plugins extend format support. Disabled by default.

```bash
# List installed plugins
markitdown --list-plugins

# Enable plugins for a conversion
markitdown --use-plugins file.pdf
```

Notable plugin: `markitdown-ocr` — OCR on embedded images in PDF/DOCX/PPTX/XLSX
via LLM Vision. Install with `pip install markitdown-ocr`.

## Usage (3 ways)

### 1. CLI — direct file or URL

```bash
# File → stdout
markitdown report.pdf

# URL → stdout
markitdown https://example.com/article

# File → output file
markitdown report.pdf -o report.md
```

### 2. CLI — pipe / stdin

```bash
# Pipe from another command
cat document.docx | markitdown

# Redirect
markitdown < presentation.pptx

# With format hint (when stdin has no extension)
cat data | markitdown -x .json
```

### 3. Python API

```python
from markitdown import MarkItDown

md = MarkItDown()

# Local file
result = md.convert("report.pdf")
print(result.text_content)

# URL
result = md.convert_url("https://example.com")
print(result.text_content)

# Stream
with open("doc.docx", "rb") as f:
    result = md.convert_stream(f, file_extension=".docx")
    print(result.text_content)
```

## CLI Options

Run `markitdown --help` for full details. Key flags:

| Flag                  | Purpose                                        |
| --------------------- | ---------------------------------------------- |
| `-o FILE`             | Write output to file instead of stdout         |
| `-x .EXT`             | Hint file extension (useful with stdin)         |
| `-m MIME`             | Hint MIME type                                 |
| `-c CHARSET`          | Hint charset (e.g. UTF-8)                      |
| `-p, --use-plugins`   | Enable 3rd-party plugins                       |
| `--keep-data-uris`    | Keep base64-encoded images (truncated by default) |
| `-d, --use-docintel`  | Use Azure Document Intelligence                |

## Expected Output

Output is `clean Markdown` text. Structure is preserved:

- Headings → `#`, `##`, `###`
- Lists → `-` or `1.`
- Tables → Markdown pipe tables
- Links → `[text](url)`
- Images → `![alt](src)` (data URIs truncated by default)
- Code blocks → fenced with language hints when detectable

Example (from a PDF):

```md
# Annual Report 2024

## Executive Summary

Revenue grew 15% year-over-year...

| Quarter | Revenue | Growth |
| ------- | ------- | ------ |
| Q1      | $2.1M   | 12%    |
| Q2      | $2.4M   | 18%    |
```

Output is meant for LLM consumption, not high-fidelity human rendering.

## Common Mistakes

| Mistake                                        | Fix                                              |
| ---------------------------------------------- | ------------------------------------------------ |
| Expecting JS-rendered content from URL         | `markitdown` fetches static HTML only             |
| Missing output for a format                    | Check if the matching extra is installed           |
| Truncated images in output                     | Use `--keep-data-uris` if you need them           |
| No YouTube transcript                          | Install `markitdown[youtube-transcription]`       |
| Running on login-gated URL                     | Use browser automation instead, then pipe HTML    |

## When NOT to Use

- JavaScript-heavy SPAs → use Playwright or browser automation first
- Login-gated pages → authenticate with browser, then pipe HTML to markitdown
- High-fidelity document rendering for humans → use dedicated viewers
- Real-time web scraping at scale → use a dedicated scraper
