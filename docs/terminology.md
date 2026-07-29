# 術語表 (Terminology)

本檔是 cc-plugin 領域名詞的`單一定義來源`。README.md 描述流程、CLAUDE.md 描述邊界，
兩者都只引用此處的名詞，不重複定義。

## 插件生態 (Plugin Ecosystem)

| 名詞 | 定義 |
| :--- | :--- |
| `Plugin` | `plugins/<name>/` 一個資料夾，含 `.claude-plugin/plugin.json` 與 `README.md`；是 skill／agent／hook／MCP／LSP 的分發單位 |
| `Skill` | `plugins/<p>/skills/<name>/SKILL.md`，符合 agentskills.io 規範的能力定義；由目錄自動探索，manifest 不列舉 |
| `Agent` | `plugins/<p>/agents/<name>.md`，具獨立 system prompt 與工具集的子代理定義 |
| `Hook` | `hooks/hooks.json` 註冊的事件腳本（如 `Stop` / `StopFailure`），由 harness 執行而非模型 |
| `Marketplace` | `.claude-plugin/marketplace.json`，plugin 來源目錄；本地用相對路徑、外部用 GitHub source |
| `Manifest` | `plugin.json`；`skills`／`agents` 只列歸屬本 plugin 的`外部`來源（本地由目錄自動探索），hooks／MCP／LSP 等 metadata 亦在此宣告 |
| `Submodule` | `.gitmodules` 記錄的 gitlink（vendored 外部 repo），與 skill 依賴的本機 clone 不同 |

## 設定層 (Configuration Layers)

| 名詞 | 定義 |
| :--- | :--- |
| `Active settings` | `config/settings.json`，經 `run.sh` 軟連結至 `~/.claude/settings.json` 的啟用設定 |
| `Provider settings` | `config/<provider>.json`（`openai`／`xai`／`proxy`／`minimax`／`llmbox`），各自完整自足；Claude Code 一次只讀一個檔，故共用區塊刻意重複 |
| `Global rule` | `config/CLAUDE.global.md`，同時軟連結為 Claude／Gemini／Codex／Hermes 的全域指令 |
| `Viper default` | `config/config.go` 的 `viper.SetDefault`，cc-plugin CLI 執行期設定的唯一預設來源 |

## 記憶蒸餾 (Memory Distillation)

| 名詞 | 定義 |
| :--- | :--- |
| `Observation` | 來源儲存（`gbrain`、`claude-mem`）的一筆原始紀錄，蒸餾管道的輸入 |
| `Candidate` | LLM 從 observation 提取出、尚未分類寫入的記憶候選 |
| `Memory` | 通過提取的一般記憶，寫入 `agentmemory` API |
| `Fact` | 通過`真實性門檻`的第一人稱事實，額外寫入 `mempalace`；`agentmemory ⊋ mempalace` |
| `真實性門檻 (Truth Qualification)` | 人類確認、第一人稱事實／經驗、或 2+ 來源佐證，三者滿其一才升級為 `Fact` |
| `指紋 (Fingerprint)` | 正規化文本 + 排序實體的 SHA-256 雜湊，用於跨來源去重 |
| `Cursor` | `StateStore` 中每個來源的讀取位置；`claude-mem` 匯出使用獨立的 `claude-mem-export` 遊標與 `observations.id` 順序 |
| `StateStore` | `model/store.go` 的 GORM + SQLite 狀態儲存，記錄 `Cursor` / `Seen` / `Distilled` |

## 知識圖譜 (Topology)

| 名詞 | 定義 |
| :--- | :--- |
| `Entity` | topology 中的一個節點，對應一份 Markdown 檔案 |
| `Edge` | entity 之間的 wikilink 關係；logger／config／helper 等噪音邊被排除 |
| `Backlink` | 反向邊，由 `cc-plugin topology rewrite` 重算並寫回 `_index.md` |
| `raw / curated` | `ultra-explore` 知識庫的兩層儲存：未經處理的擷取層與經蒸餾分級的正典層 |
