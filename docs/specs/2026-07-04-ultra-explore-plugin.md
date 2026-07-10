# Ultra-Explore Plugin — 設計計畫 (Design Plan)

建立 `plugins/ultra-explore/` 插件：從多來源建構可驗證、可增量更新的知識庫 (Knowledge Base)。

## 既有技能審查 (Review of Existing Skills)

| 參考技能               | 可取之處 (Keep)                                                    | 缺口 (Gap for KB)                              |
| ---------------------- | ------------------------------------------------------------------ | ---------------------------------------------- |
| `topology-builder`     | entity=檔案、wikilink 邊、evidence-before-edges、Phase 化 pipeline | 只建拓撲，無 provenance 分級、無 raw 層        |
| `project-explore`      | 高訊號檔案掃描、業務領域分組                                       | 輸出綁定 repo 根目錄，非集中知識庫             |
| `business-extract`     | 業務約束/風險/上下游萃取框架                                       | 一次性報告，不累積                             |
| `markitdown`           | 任意檔案/URL → Markdown                                            | 無 provenance、無去噪後入庫流程                |
| `content-summarizer`   | fetch-first、來源可追溯、fact 與 inference 分離                    | 摘要即終點，不落入知識庫                       |
| `project-route`        | index + staleness auto-rebuild 概念                                | —（借用 staleness 概念）                       |
| `tutorial`             | 教學文件輸出規範                                                   | 無知識庫作為輸入源                             |
| `understand-anything`  | 多 agent 分工、graph-reviewer 驗收概念                             | Node.js 重依賴、輸出 JSON graph 非 Markdown    |
| `team`                 | orchestrator/pipeline 編排模式                                     | —（借用 coordinator agent 模式）               |
| cc-plugin distill 管道 | 指紋去重 (SHA-256)、真實性門檻 (truth qualification)               | —（核心規則直接移植）                          |

## 核心設計 (Core Design)

兩層儲存：`raw captures`（無損入庫）與 `curated entities`（蒸餾後圖譜），中間以
真實性分級 (truth tier) 與指紋去重把關。

```text
來源 (Sources)                 管道 (Pipeline)              儲存 (Storage)
git repo ── kb-ingest-repo ─┐
web link ── kb-ingest-web ──┤→ _inbox/ (raw) → kb-distill → <zone>/entity.md
chat     ── kb-ingest-chat ─┤                  kb-connect →  邊 + Backlinks + _index.md
schema   ── kb-ingest-schema┘                  kb-verify  →  完整性/佐證/矛盾檢查
                                               kb-query   →  帶引用回答
```

- 真實性分級：`confirmed`（人工確認）> `firsthand`（第一人稱事實）>
  `corroborated`（2+ 獨立來源）> `candidate`（單一來源，不得進 curated 區）
- 指紋：`sha256(正規化文本)`，入庫前查重
- 邊規則完全沿用 `topology-builder`（方向、relation 動詞、edge grounding）
- 格式相容 `topology-builder` entity 檔，KB 可與 topologies 互通

## 技能清單 (Skills)

所有子技能 `disable-model-invocation: true` — 僅由 `/ultra-explore` 入口或
使用者手動觸發。

| Skill               | 職責                                                       |
| ------------------- | ---------------------------------------------------------- |
| `ultra-explore`     | 唯一全管道手動入口：解析來源 → 五階段 → 報告               |
| `kb-spec`           | 儲存佈局與檔案格式規範（單一事實來源，其他技能引用）       |
| `kb-ingest-repo`    | git repo 現狀 → captures + entity 候選                     |
| `kb-ingest-history` | git log 快取 + `last-commit` 游標增量 → 開發史 captures    |
| `kb-ingest-web`    | URL → 清理後 Markdown capture                              |
| `kb-ingest-chat`   | 對話紀錄 → 候選事實 captures                               |
| `kb-ingest-schema` | KV/RDB/MQ/檔案 schema → datastore captures                 |
| `kb-distill`       | inbox → curated entities（身分、去重、真實性分級、zone）   |
| `kb-connect`       | 建邊、Backlinks、`_index.md`、Mermaid 總覽                 |
| `kb-verify`        | 斷鏈/孤兒/佐證抽查/矛盾/過期檢查                           |
| `kb-query`         | 帶引用回答，缺口記入 Frontier                              |

Agent：`kb-coordinator` — 編排 ingest → distill → connect → verify 全管道。

## 驗收 (Acceptance)

- `plugins/ultra-explore/` 含 11 skills + 1 agent + plugin.json + README.md
- 所有 SKILL.md frontmatter 通過 yaml 解析，name 與目錄一致
- CLAUDE.md 專案結構與模組對應同步更新
