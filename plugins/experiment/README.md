# experiment plugin

## 用途 (Purpose)

`experiment` 是技能 (skills) 的沙盒 (sandbox) — 在技能正式合併到所屬的正式插件 (`explore`、`general`、`media` 等) 之前，先集中在這裡觀察、重構、驗證。

`experiment` 不是穩定的對外介面。預期行為：

- 暫存中的技能可能被快速調整、替換或移除。
- 當技能穩定後，會從 `experiment` 搬遷到正式插件 (例如 `explore` 或 `general`)，並從 `experiment` 移除。
- 不建議在 production workflow 中依賴此處的技能介面。

## 內容 (Contents)

`experiment` 目前收容以下十個技能：

| 技能 (Skill)        | 原位置 (Original Location)             | 預期正式歸屬 (Expected Destination) | 說明 (Notes)                                                |
| ------------------- | -------------------------------------- | ------------------------------------ | ----------------------------------------------------------- |
| `anti-sabotage`     | `plugins/general/skills/anti-sabotage` | `general`                            | 防範 Agent 自我破壞作業流程的檢查清單                       |
| `business-planner`  | `plugins/general/skills/business-planner` | `general`                          | 規劃單一功能的商業價值                                      |
| `firecrawl`         | `plugins/explore/skills/firecrawl`     | `explore`                            | Firecrawl 網頁爬取 SDK                                      |
| `markdownlint`      | `~/.claude/skills/markdownlint`        | `general`                            | Markdown 格式檢查規則子集 + 自訂限制                        |
| `mermaid`           | `plugins/general/skills/mermaid`      | `general`                            | Mermaid.js 圖表語法指南                                     |
| `model-evaluator`   | `plugins/general/skills/model-evaluator` | `general`                          | 評估 LLM 模型品質 / 表現                                  |
| `playwright-cli`    | `~/.claude/skills/playwright-cli`      | `explore`                            | Playwright 瀏覽器自動化 CLI                                 |
| `scrapling`         | `plugins/explore/skills/scrapling`     | `explore`                            | Scrapling 反爬網頁抓取框架                                  |
| `summarize.sh`      | `plugins/explore/skills/summarize.sh`  | `explore`                            | 透過 `summarize` CLI 摘要網頁 / 檔案 / YouTube / Podcast    |
| `system-planner`    | `plugins/general/skills/system-planner` | `general`                           | 規劃單一功能的系統架構                                      |

## 結構 (Structure)

```tree
plugins/experiment/
├── .claude-plugin/
│   └── plugin.json       # 插件 manifest — 列出所有收容技能
├── README.md             # 本文件
└── skills/               # 技能目錄 (每個技能一個子目錄)
    ├── business-planner/
    ├── firecrawl/
    ├── markdownlint/
    ├── mermaid/
    ├── playwright-cli/
    ├── scrapling/
    ├── summarize.sh/
    └── system-planner/
```

## 已知殘留副本 (Known Stale Copies)

由於 `playwright-cli` 與 `markdownlint` 是從 `~/.claude/skills/` (此處為 `~/.agents/skills/` 的符號連結) 移動過來的，它們在 `plugins/explore/skills/playwright-cli/` 與 `plugins/general/skills/markdownlint/` 仍有對應的舊複本。這些舊複本已從 `explore` 與 `general` 的 `plugin.json` 與 `marketplace.json` 移除登錄，但實體檔案仍存在於 `plugins/explore/skills/` 與 `plugins/general/skills/`，待手動確認後清理。

## 遷出流程 (Graduation Flow)

當某個技能準備升級為正式插件的一員時：

1. 將技能目錄從 `plugins/experiment/skills/<name>/` 移至目標插件 (例如 `plugins/explore/skills/<name>/`)。
2. 從 `plugins/experiment/.claude-plugin/plugin.json` 的 `skills` 陣列中移除該技能路徑。
3. 從 `.claude-plugin/marketplace.json` 的 `experiment` 條目 `skills` 陣列中移除，並加入目標插件的 `skills` 陣列。
4. 在本 README 的「內容」表格移除該列。
5. 若技能在 `~/.claude/skills/` 仍有符號連結，同步還原。