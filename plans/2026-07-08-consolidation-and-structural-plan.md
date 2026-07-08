# 2026-07-08 一致性清理與結構性整合計畫 (Consolidation & Structural Plan)

> 本計畫由 `consistency` 審查流程產出，整合 3 條 explore 線索的發現：
> Go 核心代碼、插件生態、計畫/文件/部署配置。目標：移除冗餘、修正結構、提升可擴展性。

## 結論 (Conclusion)

- **本專案有 4 條互相糾纏的不一致**：Go 核心代碼 (cmd/model)、插件目錄 (plugins)、計畫文件 (plans/docs)、部署配置 (ecosystem/run.sh/.gitmodules/.vscode)。
- **本計畫定義 4 個 Phase + 1 個 Add-on**，循序處理：先修核心 → 再清插件 → 再收文件 → 最後修部署。
- **不重做 10 個 `architecture-*` / 4 個 Phase 計畫的內容**，而是把它們尚未落實的項目收斂成一份 actionable checklist。
- 進度估算從當前 0% 推進到 80% (核心代碼部分完成、插件手冊自動產生、文件單一來源、部署可移植)。

---

## 範圍 (Scope)

**包含 (In-Scope)**：

- Go 核心 (`cmd/`, `model/`, `config/`)：資源釋放、讀寫去重、`topology` 移出、單一 `StateStore`、設定外化、結構化日誌、依賴清理
- 插件 (`plugins/`)：manifest 自動產生、孤立 skill 清理、命名一致
- 文件 (`plans/`, `docs/`, `README.md`, `README.todo`)：去重、單一來源、phase 進度真實化
- 部署 (`.gitmodules`, `ecosystem.config.js`, `run.sh`, `.claude/settings.local.json`, `.vscode/`)：路徑可移植、submodule 對齊

**排除 (Out-of-Scope)**：

- `plugins/understand-anything/` (uncommitted, 自帶 CLAUDE.md, 屬於獨立子專案)
- `pkg/bytedance/`, `pkg/hermes/`, `pkg/litellm/`, `pkg/usage/`, `pkg/prompt/`, `pkg/ccstatusline/` (資源型設定，未在審查範圍)
- 業務領域變更 (記憶管道語意、distill 流程)

---

## 現況盤點 (Current State Inventory)

### Go 核心 (`cmd/`, `model/`, `config/`) — P0/P1

| 編號 | 檔案:行 | 問題 | 嚴重度 |
| ---- | ------- | ---- | ------ |
| C1 | `cmd/write_agentmemory.go:30-41` | `defer resp.Body.Close()` 在迴圈內堆疊，連線/FD 洩漏 | **P0** |
| C2 | `cmd/distill.go:76-83` | `corroboration` 在迴圈內被覆寫，僅留最後一次計數 | **P0** |
| C3 | `cmd/export/mempalace.go:170-179` | 同一 room 多個 drawer 互相覆寫，僅最後一個留存 | **P0** |
| C4 | `cmd/read_logic.go:73`, `cmd/export/claudemem.go:26` | `gorm.Open(sqlite)` 未呼叫 `sqlDB.Close()` | **P0** |
| C5 | `cmd/export/mempalace.go:41-44` | `defer sqlDB.Close()` 條件式掛在 `err == nil` 上 | **P0** |
| C6 | `cmd/read_logic.go:60-103` vs `cmd/export/claudemem.go:15-56` | 兩個幾乎相同的 claude-mem 讀取函數 | **P1** |
| C7 | `cmd/read_logic.go:16-58` vs `cmd/export/gbrain.go:15-61` | 兩個幾乎相同的 gbrain 讀取函數 | **P1** |
| C8 | `cmd/root.go:13-19` vs `model/store.go:198-204` | 兩份完全相同的 `expandPath` 實作 | **P1** |
| C9 | `model/topology.go`, `model/topology_ops.go` | 知識圖譜解析器放在 `model/`，但 `cmd/` 沒人引用，違反「model = 資料」邊界 | **P1** |
| C10 | `config/default_settings.json` 為 `{}` | 與 `config.go:18-27` 的 10 個 `viper.SetDefault` 雙源真相 | **P1** |
| C11 | `main.go:9` + `cmd/root.go:24` | `config.Init()` 雙重呼叫 | **P1** |
| C12 | `model/store.go:92, 98, 128` | `AlreadyDistilled`/`MarkDistilled`/`DropDistilled` 僅測試用，正式代碼呼叫 `MarkDistilledBatch` | **P1** |
| C13 | `model/store.go:177-183` | `Close()` 回傳值被所有 caller 忽略，連線池未調校 | **P1** |
| C14 | `cmd/*.go` 多處 | 17 處 `viper.Get*` 直接呼叫，無型別化設定物件 | **P1** |
| C15 | `cmd/distill.go:34-155` | 120 行 `RunE` 混合業務邏輯與 cobra 樣板 | **P1** |
| C16 | `cmd/ollama.go:58-59` | `OllamaService.Timeout` 欄位從未讀取 | **P2** |
| C17 | `cmd/ollama.go:55` | 區域變數 `model` 遮蔽套件名 | **P3** |
| C18 | `cmd/distill.go:12-26` | `QualifiesForTruth` 是業務邏輯，卻放在 `cmd/` | **P2** |
| C19 | 6+ 處 | `"gbrain-working"`, `"claude-mem"` 字串常數散落 | **P2** |
| C20 | `model/store.go:50-52` + `162-175` | `AutoMigrate` 與 `Reset` 兩處各別列舉 model，新增第四個 table 需改兩處 | **P2** |
| C21 | `model/cursor.go:22` vs `model/store.go:110` | `Distilled` 與 `DistilledItem` 兩型別描述同概念 | **P2** |
| C22 | `model/topology.go:58` | `Findings` 為 `[]string` — 無型別、無嚴重度、無位置 | **P2** |
| C23 | `cmd/state.go` | 純 re-export 包裝，與 `cmd/export/*` 直接呼叫 `model.NewStateStore()` 風格不一致 | **P2** |
| C24 | `cmd/state_test.go` | 測試 `model.StateStore` 卻放在 `cmd/`；缺 `model/store_test.go` | **P2** |
| C25 | `cmd/retain.go:25-26`, `cmd/write_agentmemory.go:82` | `--max-age-days`, `--prune-gbrain-dir`, `--url` 旗標定義卻從未讀取 | **P2** |
| C26 | `cmd/ollama.go:63-70` | 1200-token prompt 內嵌 Go 二進位，無版本、無外部覆寫 | **P2** |
| C27 | `cmd/export/mempalace.go:303-319` | `sanitizePathComponent` / `quoteContent` 通用工具散落在命令檔 | **P3** |
| C28 | `go.mod:3` | `go 1.26.3` — Go 1.26 尚未發行 | **P1** |
| C29 | `go.mod` | `gorm.io/driver/mysql`, `anthropic-sdk-go`, `invopop/jsonschema` 直接依賴卻無引用 | **P1** |
| C30 | `go.mod` | `gosdk` 為單一呼叫 (`config.Default`)，拖入 500MB 傳遞依賴 | **P1** |

### 插件生態 (`plugins/`) — 全部 10 個本地插件

| 編號 | 插件 | 問題 | 嚴重度 |
| ---- | ---- | ---- | ------ |
| P1 | 全部 10 個 | `plugin.json` 的 `skills`/`agents` 全部為 `[]` 或缺少，與實際內容完全脫節 | **WARNING** |
| P2 | `experiment` | `skills/anti-sabotage/anti-sabotage.md` 為鬆散 `.md`，缺 `SKILL.md` 與 frontmatter | **BLOCKER** |
| P3 | `general` | `skills/changelog/` 是空資料夾，僅剩 `scripts/__pycache__` 與 `.egg-info` 建置殘留 | **BLOCKER** |
| P4 | `general` | `agents/feature.md` 引用 4 個不存在的 skill (`golang-*`) | **WARNING** |
| P5 | `general` | 與 `ultra-explore` 都有 `changelog` skill (重複) | **WARNING** |
| P6 | `tools` | 無 `README.md` (唯一例外) | **BLOCKER** |
| P7 | `tools` | 資料夾名為 `tools` 但 4 個 skill 仍維持 `apple-*` 前綴 | **WARNING** |
| P8 | `review` | `agents` 陣列用 `./agents/review-coordinator.md` 相對路徑，與其他插件風格不一致 | **NIT** |
| P9 | `ultra-explore` | `kb-coordinator` agent 的 `skills` 列表缺 `topology-builder` 與 `ultra-explore` 入口 | **WARNING** |
| P10 | `ultra-explore` | `plugin.json` 內 `agents: []`，但實際有 `kb-coordinator` | **WARNING** |
| P11 | `team` | 隨附 `roles/` 目錄 14 個檔案，manifest 與 README 未提及 | **WARNING** |
| P12 | `media` | 根目錄 `voice/`, `seedance-2.0/`, `VocCPM_setup.md` 為鬆散整合 | **WARNING** |
| P13 | `experiment` | README 結構表與實際不符 (列 8 個 skill 卻只有 7 個；標示的 3 個其實在 `general/`) | **WARNING** |
| P14 | `understand-anything` | 雙 `plugin.json` (外層與內 `understand-anything-plugin/`) 兩份相同名稱版本，loader 行為不明 | **WARNING** |
| P15 | `understand-anything` | 內層 `understand/frameworks/` 與 `languages/` 為鬆散 `.md` | **NIT** |
| P16 | `plugins/understand-anything.md` | 16KB 插件專屬文檔卻放在 `plugins/` 根目錄，應在插件目錄內 | **WARNING** |
| P17 | `plugins/README.md` | 「External Plugins」表列 7 個外部插件，6 個本地目錄不存在；3 處真理來源 (README/marketplace.json/plugins/) 不一致 | **WARNING** |
| P18 | `.claude-plugin/marketplace.json` | 列 15 個插件，5 個目錄不存在 (含 `gosdk`, `inf`, `superpowers` 等) | **WARNING** |

### 文件 (`plans/`, `docs/`, `README.md`, `README.todo`)

| 編號 | 位置 | 問題 | 嚴重度 |
| ---- | ---- | ---- | ------ |
| D1 | `plans/` 根 | 7 個主題各被寫兩次 (`topic.md` 與 `architecture-topic.md`)，日期相同，後者 2× 大小且含失效連結 | **WARNING** |
| D2 | `plans/2026-07-02-*.md` | 全部 20 個根 plans 同日 (`2026-07-02`)，破壞日期前綴的時序意義 | **NIT** |
| D3 | 7 個 feature plans | 連結到已刪除的 `architecture-cc-plugin.md` 等 6 個檔 | **WARNING** |
| D4 | `README.todo:35` | 引用 `plans/prompt-todo.md`，但檔案不存在 | **WARNING** |
| D5 | `README.todo` | Phase 4 標 100% 但兩個 item (`skill-md-renaming`, `topology-extraction`) 實際未完成；所有 checkbox 為 `[ ]` | **WARNING** |
| D6 | `README.md` 末段 | 「改善建議」區塊為 5 個開放 issue，3 個已被 plan 涵蓋，雙重出處 | **WARNING** |
| D7 | `README.md` 60-69 | 插件清單 (9 個) 與 `plugins/` 實際 (10 個) 不符；`base/apple/tmp` 不存在，`tools/ultra-explore/understand-anything` 未列 | **WARNING** |
| D8 | `docs/specs/` vs `plans/` | 兩者皆含 dated 設計文件，無共同命名/分類 | **NIT** |
| D9 | `plans/memory/2026-05-31-agent-memory-system.md` (63KB) vs `docs/superpowers/specs/2026-05-31-agent-memory-system-design.md` | 同日同主題兩份 spec | **WARNING** |
| D10 | `docs/principles/prompt/` | 雙語 prompt 工程教材，與本專案業務無明顯連結 | **NIT** |

### 部署 (`.gitmodules`, `ecosystem.config.js`, `run.sh`, `.claude/`, `.vscode/`)

| 編號 | 檔案 | 問題 | 嚴重度 |
| ---- | ---- | ---- | ------ |
| X1 | `.gitmodules` | 6 個 submodule，3 個未初始化 (`pkg/system-prompts/CL4R1T4S`, `pkg/tools/career-ops`, `plugins/media/seedance-2.0`) | **WARNING** |
| X2 | `.gitmodules` | `pkg/system_prompts/CL4R1T4S` (底線) 為 typo，實際路徑為 `pkg/system-prompts/CL4R1T4S` (連字號) | **WARNING** |
| X3 | `.gitmodules` | 同 `CL4R1T4S` 兩條目 (`system_prompts` 與 `system-prompts`)，其中一條為死條目 | **WARNING** |
| X4 | `ecosystem.config.js` | 硬編 `/Users/shuk/...` 路徑，與本機 `/Users/bytedance/` 不符 | **BLOCKER** |
| X5 | `ecosystem.config.js` | 3 個 app 缺 `cwd` 欄位 | **WARNING** |
| X6 | `ecosystem.config.js` | `cron: "10 0-9 * *"` (10 次/夜) 觸發機制無註解 | **NIT** |
| X7 | `run.sh` | `ln -sf` 重複區塊 (`~/.codex` 出現兩次) | **WARNING** |
| X8 | `run.sh` | 無 `set -e`，部分失敗時留下半配置狀態 | **WARNING** |
| X9 | `run.sh` | Claude/Gemini/Hermes 連結無 `.bak` 守衛 (CCStatusline/Tokscale 有)，不一致 | **WARNING** |
| X10 | `run.sh:39` | `~/.claude.json` 目標路徑非標準 (應為 `~/.claude/settings.json`) | **WARNING** |
| X11 | `.claude/settings.local.json` | 硬編 `/Users/shuk/...` 權限路徑 | **BLOCKER** |
| X12 | `.claude/settings.local.json` | 引用不存在的 `.mcp.json` | **WARNING** |
| X13 | `.claude/settings.local.json` | `outputStyle: "Test style"` 接線不明 | **WARNING** |
| X14 | `.vscode/launch.json` | 引用 `read-claudemem` Go 配置，但 `cmd/main.go` 沒有該 entry point | **WARNING** |
| X15 | `.vscode/settings.json` | `"-1": true` 為無效鍵；`files.exclude` 用 `false` 與意圖相反 | **WARNING** |
| X16 | `.gitignore` | 死規則：`plans/implementation_plan.md`, `plans/walkthrough.md`, `pkg/model/continue/Continue.*.vsix` | **NIT** |
| X17 | `.gitignore` | `plugins/superpower/**` (單數) 與實際 `plugins/superpowers/` (複數) 不符 | **WARNING** |
| X18 | `.gitignore` | `plugins/understand-anything/**` 排除未追蹤的目錄，造成無法加入 submodule | **WARNING** |
| X19 | 根目錄 `skills/`, `commands/`, `prompt/gpt-post/` | 為空或未追蹤的遺留 scaffolding | **NIT** |
| X20 | `.lsp.json`, `.mcp.json` | 引用於 CLAUDE.md 但檔案不存在 | **NIT** |

---

## Phase 1 — Go 核心資源安全 (P0 修復)

> 目標：消除資源洩漏與明顯正確性 bug。預估工作量：1 人天。

- [ ] **C1** 修正 `cmd/write_agentmemory.go` 迴圈內 `defer`：將 `defer` 移出迴圈 (函數結尾) 或於每次迭代結束顯式 `resp.Body.Close()`
- [ ] **C2** 修正 `cmd/distill.go:76-83` `corroboration` 覆寫：先收集所有 ref 的 count，再以 `max` 或 `len(unique sources)` 計算
- [ ] **C3** 修正 `cmd/export/mempalace.go:170-179` 同 room 多 drawer 覆寫：改為 append/merge 模式
- [ ] **C4/C5** 統一 `gorm.Open(sqlite)` 連線釋放：所有 reader 函數結尾呼叫 `sqlDB, _ := db.DB(); defer sqlDB.Close()`，移除 C5 條件式 defer
- [ ] **C13** 為 `StateStore.Close()` 加連線池調校 (`SetMaxOpenConns(1)` for SQLite, `SetMaxIdleConns(1)`, `SetConnMaxLifetime`)
- [ ] **C28** 修正 `go.mod:3` 至 `go 1.25.0` (或目前發行版本)
- [ ] 跑 `go test ./... -count=1` 確認全綠

## Phase 2 — Go 核心結構重構 (消除冗餘與層次不清)

> 目標：拆分 `cmd/` 業務邏輯、定義介面、收斂設定存取。預估工作量：3-4 人天。

- [ ] **C6/C7** 合併 gbrain/claude-mem 重複讀取函數：以 `cmd/export/*` 的參數化版本 (`gbrainRead(s, dir, fromCursor)`) 為主，`cmd/read_logic.go` 改為 thin wrapper
- [ ] **C8** 移除 `cmd/root.go:13-19` 的 `expandPath`，統一呼叫 `model.ExpandPath` (或反之，擇一並刪另一份)
- [ ] **C9** 將 `model/topology.go`, `model/topology_ops.go` 移至 `pkg/topology/` (新套件)，加 `cmd/topology` 進入點
- [ ] **C10** 將 `config.go:18-27` 10 個 `viper.SetDefault` 全部遷移至 `config/default_settings.json`，移除 `WithDefaultValue` 後的死路徑
- [ ] **C11** 移除 `main.go:9` 重複呼叫 (`cmd/root.go.init()` 已執行)
- [ ] **C12** 移除 `AlreadyDistilled` / `MarkDistilled` / `DropDistilled` (僅測試用) 或將 `MarkDistilledBatch` 改為內部呼叫
- [ ] **C14** 引入 `internal/config` 套件：型別化設定物件 (`type Config struct { State StateConfig; LLM LLMConfig; Sources SourcesConfig; Stores StoresConfig }`)，所有 `viper.Get*` 改為 `cfg.Field`
- [ ] **C15** 將 `cmd/distill.go` 120 行 `RunE` 拆為 `pipeline.Run(ctx, cfg) error` (新套件 `internal/distill/`)
- [ ] **C18** 將 `QualifiesForTruth` 從 `cmd/distill.go` 移至 `internal/distill/`
- [ ] **C19** 引入 `internal/source` 常數：`const SourceGBrain = "gbrain-working"; const SourceClaudeMem = "claude-mem"`
- [ ] **C20** 集中 `schema.Models()` 函式：所有 GORM model 註冊於一處
- [ ] **C21** 合併 `Distilled` 與 `DistilledItem` (讓 `DistilledAt` 為 `*int64` 或共用型別)
- [ ] **C22** 改 `Topology.Findings` 為 `[]Finding{ Message, Severity, Location }` 結構
- [ ] **C23** 決定 `cmd/state.go` 去留：刪除或補上方法；統一 `cmd/export/*` 改用 `cmd.state.NewStateStore()`
- [ ] **C24** 將 `cmd/state_test.go` 移至 `model/store_test.go`
- [ ] **C25** 移除 `cmd/retain.go:25-26`, `cmd/write_agentmemory.go:82` 死旗標
- [ ] **C26** 將 1200-token prompt 抽至 `pkg/prompt/extract.go` (或 `~/.config/cc-plugin/prompts/extract.txt`)
- [ ] **C16/C17** 移除 `OllamaService.Timeout` 未使用欄位；將區域變數 `model` 改名 `modelName`
- [ ] **C27** 將 `sanitizePathComponent` / `quoteContent` 移至 `pkg/text/`
- [ ] **C29** 從 `go.mod` 移除 `gorm.io/driver/mysql`, `anthropic-sdk-go`, `invopop/jsonschema` (未引用)；跑 `go mod tidy`
- [ ] **C30** 評估 `gosdk` 替代：若 `config.Default` 可由 5 行 viper 程式碼取代，移除 `gosdk` 直接依賴

## Phase 3 — 插件生態清理 (P 編號)

> 目標：所有 plugin manifest 與實際內容一致；孤立 skill 處理；命名收斂。預估工作量：2-3 人天。

- [ ] **P1** 寫 `scripts/sync_plugin_manifest.go` 或 `scripts/sync_plugin_manifest.sh`：掃描 `plugins/*/skills/` 與 `plugins/*/agents/`，自動更新每個 `plugin.json` 的 `skills` 與 `agents` 陣列；CI 驗證
- [ ] **P2** `experiment/skills/anti-sabotage/anti-sabotage.md` → 重新命名為 `SKILL.md` 並補 YAML frontmatter
- [ ] **P3** 決定 `general/skills/changelog/`：刪除整個目錄 (因 `ultra-explore/skills/changelog/` 為主要副本)
- [ ] **P4** 修正 `general/agents/feature.md` 的 dangling skill 引用
- [ ] **P5** 確認 P3 處理後 `changelog` 僅存在於 `ultra-explore`
- [ ] **P6** 為 `plugins/tools/` 補上 `README.md` (說明 Apple 整合策略)
- [ ] **P7** 評估 `tools/` 資料夾更名：`tools/apple/` (子目錄) 或保留 `tools/` 但加 `keywords: ["apple", "macos"]` 標示當前範疇
- [ ] **P8** 統一 `review/.claude-plugin/plugin.json` 的 `agents` 路徑：改為 `./agents/review-coordinator` (無 `.md`) 或全部插件一致
- [ ] **P9** 補 `kb-coordinator` agent 的 `skills` 列表：加入 `topology-builder`, `ultra-explore`
- [ ] **P10** 修正 `ultra-explore/.claude-plugin/plugin.json` 的 `agents` 陣列
- [ ] **P11** `team/README.md` 補上 `roles/` 目錄說明
- [ ] **P12** `media/` 根目錄的鬆散檔：移入 `media/skills/voice/` 或建立 `media/integrations/voice/` 與 `media/integrations/seedance/`
- [ ] **P13** 修正 `experiment/README.md` 結構表使其與實際一致
- [ ] **P14** 釐清 `understand-anything` 雙 `plugin.json` 行為：保留內層 `understand-anything-plugin/.claude-plugin/plugin.json` (loader 從 skills 自動掃描)，刪除外層空 manifest
- [ ] **P16** 將 `plugins/understand-anything.md` 移入 `plugins/understand-anything/README.zh-TW.md`
- [ ] **P17** 重整 `plugins/README.md`：分離「本地插件」與「外部插件/Submodule」兩表；本地插件加入 `tools, ultra-explore, understand-anything, experiment`；外部 6 個非本地目錄加註 N/A
- [ ] **P18** 修正 `.claude-plugin/marketplace.json`：標示 `awesome-claude-code-subagents`, `gosdk`, `inf`, `superpowers`, `last30days`, `ui-ux-pro-max-skill` 為外部來源 (已正確) 並驗證 5 個本地子模組 (`tools, explore, general, experiment, review, media, team`) 路徑存在

## Phase 4 — 文件單一來源 (D 編號)

> 目標：消除 plans/docs/README 三邊重複。預估工作量：1-2 人天。

- [ ] **D1** 合併 7 對 `topic.md` 與 `architecture-topic.md`：
    - 保留較長的 `architecture-*` 版本 (5.5KB+)
    - 將較短 `topic.md` 重新命名為 `YYYY-MM-DD-topic.md` 並補進度狀態 (D5)
    - 移動至 `docs/specs/` 統一存放已批准設計；`plans/` 只保留當前進行中項目
- [ ] **D2** 為同日多檔案加副時間戳 (例如 `2026-07-02T14-arch-X.md`)，或合併為單檔
- [ ] **D3** 全面掃描 plans 中的失效 `[link](architecture-*.md)`：移除或替換為 `archive://` 形式
- [ ] **D4** 補上 `plans/prompt-todo.md` (從 README.todo 描述重建) 或從 README.todo 移除該引用
- [ ] **D5** 重整 `README.todo`：
    - 為每個項目加 `Status: [ ] not-started / [~] in-progress / [x] done` frontmatter
    - Phase 4 標 0% (而非 100%) 直到 skill-md-renaming 與 topology-extraction 真正完成
    - 移除「Archive」段的 6 個已刪檔名 (或改為「已歸檔，連結至 git log」)
- [ ] **D6** 將 `README.md` 末段「改善建議」整段移除 (內容已收錄於 plans/ 與 README.todo)
- [ ] **D7** 重寫 `README.md` 60-69 插件清單使其與 `plugins/` 實際 10 個目錄一致
- [ ] **D8** 制定 `docs/specs/` 與 `plans/` 分類規則：`docs/specs/` 為 approved design (有 status: accepted)；`plans/` 為 WIP/未來項目
- [ ] **D9** 合併 `plans/memory/2026-05-31-agent-memory-system.md` 與 `docs/superpowers/specs/2026-05-31-agent-memory-system-design.md`：保留 `plans/` 版本 (63KB 較完整) 並將 docs 指向 plans
- [ ] **D10** 評估 `docs/principles/prompt/`：若為通用教育資源且不影響本專案，標 `.gitignore` 或移至外部 repo

## Phase 5 — 部署可移植性 (X 編號)

> 目標：任何人 clone 後能 `go build && ./run.sh` 直接運行。預估工作量：1 人天。

- [ ] **X1/X2/X3** 清理 `.gitmodules`：
    - 移除 typo 條目 `pkg/system_prompts/CL4R1T4S` (底線)
    - 對未初始化的 3 條 (`CL4R1T4S` 連字號, `career-ops`, `seedance-2.0`) 執行 `git submodule update --init` 或從 `.gitmodules` 移除
- [ ] **X4** 將 `ecosystem.config.js` 內 `/Users/shuk/...` 改為 `${HOME}` 變數或 `path.resolve(__dirname, ...)`
- [ ] **X5** 為 3 個 pm2 app 補 `cwd: process.cwd()`
- [ ] **X6** 為 `cron: "10 0-9 * *"` 加註解說明用途
- [ ] **X7** 移除 `run.sh` 重複 `~/.codex` 區塊
- [ ] **X8** 為 `run.sh` 開頭加 `set -euo pipefail`
- [ ] **X9** 為 Claude/Gemini/Hermes 連結加 `.bak` 守衛 (與 CCStatusline/Tokscale 一致)
- [ ] **X10** 修正 `run.sh:39` 將 `~/.claude.json` 改為 `~/.claude/settings.json` (或移除此行)
- [ ] **X11** 將 `.claude/settings.local.json` 內 `/Users/shuk/...` 改為 `${workspaceFolder}` 或本機實際路徑
- [ ] **X12/X20** 建立 `.mcp.json` 範本 (或從 `.claude/settings.local.json` 移除 `enableAllProjectMcpServers` 與 `enabledMcpjsonServers`)
- [ ] **X13** 為 `.claude/output-styles/test-style.md` 與 `outputStyle` 設定補接線文件
- [ ] **X14** 修正 `.vscode/launch.json`：移除 `read-claudemem` 配置或建立對應 `cmd/read_claudemem.go` 進入點
- [ ] **X15** 修正 `.vscode/settings.json`：移除 `"-1": true` 鍵；`files.exclude` 改為 `true` (意圖為隱藏)
- [ ] **X16** 從 `.gitignore` 移除死規則
- [ ] **X17** `.gitignore` 將 `plugins/superpower/**` 改為 `plugins/superpowers/**`
- [ ] **X18** 從 `.gitignore` 暫時移除 `plugins/understand-anything/**` 直到該 plugin 決定 submodule 化或留在主 repo
- [ ] **X19** 移除根目錄空 `skills/`, `commands/`, `prompt/` 殘留

---

## Add-on — Plugin Manifest 自動同步腳本

> 目標：P1 永久方案。預估工作量：0.5 人天。

- [ ] 寫 `scripts/sync-plugin-manifest.sh` (或 `cmd/sync-plugin-manifest.go`)：
    - 掃描 `plugins/*/skills/*/SKILL.md` 與 `plugins/*/agents/*.md`
    - 自動產生 `plugins/*/.claude-plugin/plugin.json` 的 `skills` 與 `agents` 陣列
    - 驗證 frontmatter 完整性 (name 與目錄名一致、description 含觸發詞)
- [ ] 加 `make sync-plugin-manifest` 目標至 `Makefile` (或 `run.sh`)
- [ ] 整合至 `pre-commit` hook (`.git/hooks/pre-commit`)

---

## 執行順序 (Execution Order)

```tree
Phase 1 (P0 安全) ── 1 day ──► Phase 2 (結構重構) ── 3-4 days ──► Phase 3 (插件) ── 2-3 days
                                                                              │
                                              ┌───────────────────────────────┘
                                              ▼
                                          Phase 4 (文件) ── 1-2 days ──► Phase 5 (部署) ── 1 day
                                                                                          │
                                                                                          ▼
                                                                              Add-on (manifest sync) ── 0.5 day
```

**總工作量估算**：8.5–11.5 人天 (單人) 或 4-5 個 sprint (3 人小組)。

## 風險 (Risks)

| 風險 | 影響 | 緩解 |
| ---- | ---- | ---- |
| Phase 2 重構破壞既有測試 | 中 | 每次重構跑 `go test ./...`；先擴展測試再改實作 |
| Phase 3 變更 plugin 名稱影響已安裝用戶 | 低 | 透過 `keywords` 而非 skill 名稱做向下相容；`apple-*` 保留別名 |
| Phase 4 合併 plans 失去時序紀錄 | 中 | 保留 git history 與 `archive/` 子目錄 |
| Phase 5 改 `run.sh` 行為破壞用戶現有 symlink | 中 | 加 `--dry-run` 模式 + `.bak` 守衛 |

## 不在本計畫範圍 (Non-Goals)

- 新功能 (例如新增第二個 LLM provider、新增第二個來源)
- 插件業務邏輯變更 (例如修改 `ultra-explore` 流程)
- `understand-anything` 內部結構 (其 CLAUDE.md 已自管)
- `pkg/bytedance/`, `pkg/hermes/`, `pkg/litellm/`, `pkg/usage/`, `pkg/prompt/`, `pkg/ccstatusline/` (外部工具設定，未在審查範圍)
- `data/`, `localtest/`, `tmp/` 內容調整 (instance-specific 資料)

## 驗收 (Acceptance Criteria)

- [ ] `go test ./... -count=1` 全綠
- [ ] `go vet ./...` 無警告
- [ ] `golangci-lint run` (若已配置) 通過
- [ ] 每個 `plugins/*/skills/` 子目錄都有 `SKILL.md` 含 frontmatter
- [ ] 每個 `plugins/*/.claude-plugin/plugin.json` 的 `skills` 與 `agents` 與實際檔案清單一致
- [ ] `git status` 乾淨 (除 `plugins/understand-anything/` 與 `config/settings.json` 等個人設定)
- [ ] 新機器 clone repo 後 `go build && ./run.sh && cc-plugin distill` 一次成功
- [ ] `README.todo` phase 進度反映實際完成度
- [ ] `README.md` 插件清單與 `plugins/` 與 `marketplace.json` 三邊一致
