# experiment plugin

## 用途 (Purpose)

`experiment` 是技能 (skills) 的沙盒 (sandbox) — 在技能正式合併到所屬的正式插件
（`explore`、`general` 等）之前，先集中在這裡觀察、重構、驗證。

`experiment` 不是穩定的對外介面。預期行為：

- 暫存中的技能可能被快速調整、替換或移除。
- 當技能穩定後，會搬遷到正式插件並從 `experiment` 移除。
- 不建議在 production workflow 中依賴此處的技能介面。

## 內容 (Contents)

技能一律位於 `plugins/experiment/skills/<name>/`，由目錄自動探索，manifest 不列舉。

| 技能 (Skill)      | 預期歸屬 (Destination) | 說明 (Notes)                                                |
| ----------------- | ---------------------- | ----------------------------------------------------------- |
| `anti-sabotage`   | `experiment`（常駐）   | 防範 Agent 自我破壞作業流程的檢查清單                       |
| `firecrawl`       | `explore`              | Firecrawl 網頁爬取 SDK                                      |
| `mermaid`         | `general`              | Mermaid.js 圖表語法指南                                      |
| `model-evaluator` | `general`              | 評估 LLM 模型品質／表現                                      |
| `playwright-cli`  | `explore`              | Playwright 瀏覽器自動化 CLI                                  |
| `scrapling`       | `explore`              | Scrapling 反爬網頁抓取框架                                   |
| `summarize-sh`    | `experiment`（常駐）   | 透過 `summarize.sh` CLI 摘要網頁／檔案／YouTube／Podcast     |

`skills/summarize.sh/` 是 `summarize-sh` 依賴的本機 clone（未追蹤、已 gitignore），
不是技能目錄。

另有一個歸屬本插件的`外部`技能，由 manifest 的 `skills` 宣告、skills cmd 自動取得：
`guangyuspace/codex-gamestudio-skill`。

## 結構 (Structure)

```tree
plugins/experiment/
├── .claude-plugin/plugin.json   # manifest（skills／agents 只列外部來源）
├── README.md                    # 本文件
└── skills/<name>/SKILL.md       # 每個技能一個子目錄
```

## 遷出流程 (Graduation Flow)

1. 將 `plugins/experiment/skills/<name>/` 移至目標插件的 `skills/` 下。
2. 在本 README 的「內容」表格移除該列，並在目標插件 README 補上。
3. 若 `~/.claude/skills/` 仍有指向舊路徑的符號連結，同步更新。

本地技能由目錄自動探索，遷移不需要改動 `plugin.json` 或 `marketplace.json`；
manifest 的 `skills`／`agents` 只在增減`外部`來源時才需要編輯。
