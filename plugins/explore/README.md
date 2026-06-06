# 探索插件 (Explore Plugin)

`explore` 插件整合 6 個 `抓取`、`爬蟲` 與 `內容摘要` 技能，用於從網頁、檔案、影片等來源擷取並整理內容。

## 技能清單 (Skills)

| 技能                 | 工具                          | 路徑                          |
| -------------------- | ----------------------------- | ----------------------------- |
| `content-summarizer` | workflow（底層 `markitdown`） | `./skills/content-summarizer` |
| `firecrawl`          | `firecrawl` CLI（雲端）       | `./skills/firecrawl`          |
| `markitdown`         | `markitdown` Python CLI       | `./skills/markitdown`         |
| `playwright-cli`     | `playwright-cli` 瀏覽器自動化 | `./skills/playwright-cli`     |
| `scrapling`          | `scrapling` Python 框架       | `./skills/scrapling`          |
| `summarize.sh`       | `summarize` CLI               | `./skills/summarize.sh`       |

## 技能比較 (Skill Comparison)

> 測試標的 (URL)：`https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices`
> 測試日期：2026-06-07

### 結論 (TL;DR)

| 面向             | 推薦                                                                 | 原因                             |
| ---------------- | -------------------------------------------------------------------- | -------------------------------- |
| 速度             | `markitdown`                                                         | `1.0s` 純 HTTP                   |
| 內容乾淨度       | `markitdown` / `content-summarizer`                                  | 已自動裁切 nav/footer            |
| 主文完整度       | `markitdown`                                                         | 結尾 mid-code-block 表示完整下載 |
| `JS` 渲染 / 反爬 | `playwright-cli`                                                     | 唯一真瀏覽器                     |
| 附加價值         | `content-summarizer`                                                 | `TL;DR` + 6 要點 + 3 商業提案    |
| 無需 API key     | `markitdown` / `scrapling` / `playwright-cli` / `content-summarizer` | 完全本地執行                     |

### 完整比較表 (Detailed Comparison)

| 技能                 | 工具/版本                | 耗時     | 字元數   | 詞數    | 主文含側欄 | 設定需求                                        |
| -------------------- | ------------------------ | -------- | -------- | ------- | ---------- | ----------------------------------------------- |
| `markitdown`         | `markitdown 0.1.6`       | `1.0s`   | `62,064` | `8,739` | ✗ 乾淨     | `pip install markitdown[all]`                   |
| `scrapling`          | `scrapling 0.4.8` (HTTP) | `1.97s`  | `64,239` | `8,510` | ⚠ 部分     | `pip install scrapling browserforge playwright` |
| `firecrawl`          | `firecrawl v1.18.0`      | `1.475s` | `70,488` | `8,552` | ⚠ 大量     | `FIRECRAWL_API_KEY` 或 stored credentials       |
| `playwright-cli`     | `playwright-cli 0.1.13`  | `54s`    | `56,132` | `7,595` | ✓ 完整     | 瀏覽器已安裝                                    |
| `content-summarizer` | workflow                 | `143s`   | `67,566` | `9,517` | ✗ 乾淨     | LLM provider（用於摘要）                        |
| `summarize.sh`       | `summarize` CLI          | —        | —        | —       | —          | ❌ LLM provider + prompt config                 |

### 速度階梯 (Speed Ladder)

```mermaid
graph LR
  A["markitdown<br/>1.0s"] --- B["firecrawl<br/>1.475s"]
  B --- C["scrapling<br/>1.97s"]
  C --- D["playwright-cli<br/>54s"]
  D --- E["content-summarizer<br/>143s"]
```

### 內容品質觀察 (Content Quality)

- `markitdown`：H1 起頭、code block 收尾，nav 殘留為零
- `content-summarizer`：同 `markitdown` 內容 + 結構化摘要與商業提案
- `scrapling`：首段含 `Loading...` 殘字、nav 為行內列表
- `firecrawl`：首段含整個 `Claude API Docs Home` nav、結尾含 `### Terms and policies` 頁尾
- `playwright-cli`：內容正確但 `eval document.body.innerText` 回傳的 JSON 字串內 `\n` 是字面跳脫

## 選用指引 (When to Use)

| 情境                              | 推薦                               |
| --------------------------------- | ---------------------------------- |
| 純粹抓文檔頁 markdown             | `markitdown`                       |
| 需要 `JS` 渲染或繞過反爬          | `playwright-cli`                   |
| 需要乾淨 HTTP 抓取 + 自訂轉檔     | `scrapling`                        |
| 已訂閱 `firecrawl`、要雲端託管    | `firecrawl`                        |
| 要「重點 + 行動提案」一站完成     | `content-summarizer`               |
| 想用 LLM 重寫為摘要、不需本地 LLM | `summarize.sh`（需先設定 API key） |

## 跳過的技能 (Skipped Skills)

- `summarize.sh`：設定檔 `~/.summarize/config.json` 缺 `prompt` 且所有 `apiKeys` 為空，本地 Ollama 也未運行 → `Invalid config file ... "prompt" must not be empty.` 停止
- `firecrawl`：原本擔心缺 `FIRECRAWL_API_KEY`，但 `firecrawl --status` 確認已認證（`998/1000` credits）→ 改為派代理成功

## 測試輸出 (Test Artifacts)

```text
/tmp/skill-comparison/
├── markitdown.md           62,064 B   8,739 w   904 L
├── scrapling.md            64,239 B   8,510 w
├── firecrawl.md            70,488 B   8,552 w   997 L
├── playwright-cli.md       56,132 B   7,595 w
├── content-summarizer.md   67,566 B   9,517 w
└── fetch_scrapling.py
```

## 相關連結 (See Also)

- 插件清單：`plugins/explore/.claude-plugin/plugin.json`
- 技能註冊表：`skills.json`
