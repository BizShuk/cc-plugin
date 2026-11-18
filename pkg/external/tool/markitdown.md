# MarkItDown

來源: [microsoft/markitdown](https://github.com/microsoft/markitdown) (MIT, ~196k stars, Python, 最新 release `v0.1.9` 2026-10-21) · [PyPI](https://pypi.org/project/markitdown/) · [markitdown-mcp](https://github.com/microsoft/markitdown/tree/main/packages/markitdown-mcp)

Microsoft (AutoGen 團隊) 維護的輕量 Python CLI + library, 把各種檔案與 URL 轉成 Markdown, 供 LLM 與文字分析 pipeline 使用. 定位接近 [textract](https://github.com/deanmalmgren/textract), 差別在於保留文件結構 (heading, list, table, link). 輸出是給機器讀的 token-efficient 文字, `不是` 給人看的高保真轉檔.

這個 repo `不是 Claude skill`: 是一個 Python monorepo, `packages/` 下有 `markitdown` (核心), `markitdown-mcp` (MCP server), `markitdown-ocr` (OCR plugin), `markitdown-sample-plugin` (plugin 範本). 本 repo 內的使用指引在 [plugins/explore/skills/markitdown](../../../plugins/explore/skills/markitdown/SKILL.md) skill; 這份文件只記上游是什麼.

## 支援格式

| 類別      | 格式                                                    | 需要的 extra                                   |
| --------- | ------------------------------------------------------- | ---------------------------------------------- |
| Office    | Word, PowerPoint, Excel (`.xlsx` / 舊版 `.xls`)         | `[docx]` `[pptx]` `[xlsx]` `[xls]`             |
| 文件      | PDF, EPUB, HTML                                         | `[pdf]`                                        |
| 文字資料  | CSV, JSON, XML                                          | 內建                                           |
| 郵件      | Outlook `.msg`                                          | `[outlook]`                                    |
| 圖片      | EXIF metadata + LLM caption                             | 需傳入 `llm_client`                            |
| 音訊      | EXIF metadata + speech transcription (wav / mp3)        | `[audio-transcription]`                        |
| 影片      | YouTube URL 字幕                                        | `[youtube-transcription]`                      |
| 壓縮檔    | ZIP (逐一轉換內容)                                      | 內建                                           |

`[all]` 一次裝齊. 圖片 caption 與語音轉錄依賴外部雲端服務.

## 進階轉換後端

| 後端                         | 啟用方式                                        | 特色                                                                          | 成本            |
| ---------------------------- | ----------------------------------------------- | ----------------------------------------------------------------------------- | --------------- |
| Built-in converters          | 預設                                            | 離線, 依格式抽取                                                              | 本機運算        |
| `markitdown-ocr` plugin      | `pip install markitdown-ocr` + `--use-plugins`  | 用 LLM Vision 對 PDF / DOCX / PPTX / XLSX 內嵌圖片做 OCR; 沒給 `llm_client` 時`靜默跳過` | LLM API 費用    |
| Azure Document Intelligence  | `-d -e <endpoint>` 或 `MARKITDOWN_DOCINTEL_ENDPOINT` | 雲端 layout / OCR, 可抽 query fields 成 YAML front matter                    | 按次計費        |
| Azure Content Understanding  | `--use-cu` 或 `MARKITDOWN_CU_ENDPOINT`          | 唯一支援影片; 自訂 analyzer 抽結構化欄位 (發票金額, 合約條款)                 | 按次計費        |

## Agent 介面

- CLI: `markitdown file.pdf -o out.md`, 也吃 stdin (`cat file.pdf | markitdown`); `--list-plugins` / `--use-plugins` 控制第三方 plugin (預設關閉).
- Python: `MarkItDown(enable_plugins=False).convert(path).markdown`; 安全起見改用最窄的 `convert_local()` / `convert_stream()` / `convert_response()`.
- MCP: `markitdown-mcp` 提供單一 tool `convert_to_markdown(uri)`, 接受 `http:` / `https:` / `file:` / `data:` URI; 預設 STDIO, `--http --host 127.0.0.1 --port 3001` 改走 Streamable HTTP + SSE.

Claude Code 接法 (`.mcp.json`, 推論, 未實測):

```json
{
  "mcpServers": {
    "markitdown": { "command": "markitdown-mcp" }
  }
}
```

## 安裝

- 需求: Python 3.10 - 3.14, 建議放 venv.
- `pip install 'markitdown[all]'`, 或只裝需要的 extra: `pip install 'markitdown[pdf,docx,pptx]'`.
- MCP: `pip install markitdown-mcp`.
- Docker: `docker build -t markitdown:latest .` 後 `docker run --rm -i markitdown:latest < file.pdf > out.md`.
- 本機現況: `~/.venv/bin/markitdown` 為 `0.1.6`, 落後上游 `v0.1.9`.

## 適用情境

- Read tool 打不開的格式 (DOCX, PPTX, XLSX, EPUB, `.msg`) 先轉 Markdown 再交給 agent 讀.
- RAG / 摘要 pipeline 的前處理: 批次把文件夾轉成 `.md`.
- 掃描版 PDF 或圖片很多的簡報, 搭 `markitdown-ocr` 或 Azure Document Intelligence 補文字.

## 注意

- MarkItDown 以`當前 process 的權限`做 I/O, 等同 `open()` / `requests.get()`; `convert()` 會接受本機路徑, 遠端 URI 與 stream. 處理不可信輸入時要自行限制路徑, URI scheme 與網段 (擋 loopback, link-local, metadata service), 否則有 SSRF 與任意讀檔風險.
- `markitdown-mcp` 只設計給本機可信 agent 使用; HTTP 模式預設綁 `localhost`, 文件未提供驗證機制, `不要` 綁到其他介面.
- 想要完全離線, 只裝需要的 extra, 並避開 LLM caption, 語音轉錄與 Azure 後端.
- 上游明確不收 web server, REST API, GUI 這類應用; 新格式建議做成第三方 plugin (GitHub tag `#markitdown-plugin`).
