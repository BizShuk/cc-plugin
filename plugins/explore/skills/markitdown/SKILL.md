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

```python
from markitdown import MarkItDown

md = MarkItDown()
result = md.convert("report.pdf")                # local file
result = md.convert_url("https://example.com")   # URL
print(result.text_content)
```

Key flags (`markitdown --help` for all): `-o FILE` output to file, `-x .EXT`
extension hint, `-p / --use-plugins`, `--keep-data-uris` (keep base64 images,
truncated by default), `-d / --use-docintel` (Azure Document Intelligence).

## Output

Clean Markdown for LLM consumption (not high-fidelity human rendering):
headings → `#`, lists → `-` / `1.`, tables → pipe tables, links → `[text](url)`,
images → `![alt](src)` (data URIs truncated by default).

## Workflow: 轉換後存入 Apple Notes (Convert → Save to Apple Notes)

把任意來源（URL / PDF / DOCX…）轉成 Markdown 後存進 Apple Notes 長期保存，搭配
`apple-notes` skill 寫入。

內容規則：

- 只留主要內容：去掉導覽列、頁尾、廣告、側欄等雜訊，只保留正文
- 保留所有連結／參考：markitdown 產出的 `[text](url)` 與引用清單不要刪
- Source 連結要帶頁面標題：在筆記開頭補一行 `Source: [頁面標題](原始 URL)`

```bash
# 1) 轉成 markdown
markitdown https://example.com/article -o /tmp/page.md

# 2) 清掉雜訊、保留正文與連結，並在開頭補上 `Source: [頁面標題](URL)`

# 3) 用 apple-notes skill 存入（-m 吃 markdown；-f 指定資料夾）
notes add -F /tmp/page.md -m -f Report
```

要在 Notes 內取得「可點擊超連結 + 原生標題高亮」，改產生 HTML 並用 `notes add -h`
（`<h1>` 原生標題、`<a href>` 超連結；見 `apple-notes` skill 的 Rich Text
Formatting 說明）。

> 純 URL 文章其實 `notes add -u <URL>` 就會自動抓取並清理；markitdown 路線的價值
> 在於 PDF / DOCX / PPTX 等非 URL 來源，或需要精準控制 Markdown 結構時。

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
