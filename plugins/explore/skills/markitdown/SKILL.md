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

Python CLI and library (by Microsoft) that converts files and URLs to clean,
token-efficient Markdown for LLM pipelines — preserves headings, lists, tables,
and links.

Source: [GitHub - microsoft/markitdown: Python tool for converting files and office documents to Markdown](https://github.com/microsoft/markitdown)

## What It Works For

| Category  | Formats                                                  |
| --------- | -------------------------------------------------------- |
| Web       | URLs (fetches static HTML → Markdown)                    |
| Documents | PDF, DOCX, PPTX, XLSX, XLS, EPUB                          |
| Data      | CSV, JSON, XML                                           |
| Media     | Images (EXIF + OCR), Audio (EXIF + transcription)        |
| Video     | YouTube URLs (transcript extraction)                     |
| Archives  | ZIP (iterates contents)                                  |
| Other     | Outlook `.msg`, plain HTML                               |

Some formats need optional extras (see Installation).

## Installation

Requires `Python 3.10+`.

```bash
markitdown --version                            # check if installed
pip install 'markitdown[all]'                   # all formats
pip install 'markitdown[pdf,docx,pptx,xlsx]'    # selective extras
```

Extras: `[all]`, `[pdf]`, `[docx]`, `[pptx]`, `[xlsx]`, `[xls]`, `[outlook]`,
`[audio-transcription]`, `[youtube-transcription]`, `[az-doc-intel]`. A missing
format usually means its extra isn't installed.

Plugins (disabled by default) extend support — list with `markitdown
--list-plugins`, enable per run with `--use-plugins`. Notable:
`markitdown-ocr` (LLM-Vision OCR on embedded images; `pip install markitdown-ocr`).

## Usage

```bash
# CLI: file / URL → stdout or file
markitdown report.pdf
markitdown https://example.com/article
markitdown report.pdf -o report.md

# CLI: stdin (add -x .EXT when input has no extension)
cat document.docx | markitdown
cat data | markitdown -x .json
```

Python 函式庫用法（`MarkItDown().convert(...)`）見
[references/python-api.md](references/python-api.md)。

Key flags (`markitdown --help` for all): `-o FILE` output to file, `-x .EXT`
extension hint, `-p / --use-plugins`, `--keep-data-uris` (keep base64 images,
truncated by default), `-d / --use-docintel` (Azure Document Intelligence).

## Output

Clean Markdown for LLM consumption (not high-fidelity human rendering):
headings → `#`, lists → `-` / `1.`, tables → pipe tables, links → `[text](url)`,
images → `![alt](src)` (data URIs truncated by default).

## Common Mistakes

| Mistake                                | Fix                                          |
| -------------------------------------- | -------------------------------------------- |
| Expecting JS-rendered content from URL | `markitdown` fetches static HTML only        |
| Missing output for a format            | Check the matching extra is installed        |
| Truncated images in output             | Use `--keep-data-uris`                       |
| No YouTube transcript                  | Install `markitdown[youtube-transcription]`  |
| Login-gated URL                        | Authenticate in browser, then pipe HTML in   |

## When NOT to Use

- JavaScript-heavy SPAs → use Playwright / browser automation first
- Login-gated pages → authenticate with browser, then pipe HTML to markitdown
- High-fidelity rendering for humans → use dedicated viewers
- Real-time web scraping at scale → use a dedicated scraper

## Related

- `[[apple-notes]]` 轉檔後存入 Apple Notes 的流程歸它擁有，見該技能的
  `references/workflows.md`
