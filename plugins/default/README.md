# default — 雙相容預設插件 (Default Plugin)

同時相容 Claude Code、Antigravity (`agy`) 與 Grok Build 的預設插件。三家共用同一份內容，平台專屬入口一律用軟連結指回本體，免除重複維護。

## 資料夾相容性 (Folder Compatibility)

| 路徑 | Claude Code | Antigravity (`agy`) | 說明 |
| :--- | :--- | :--- | :--- |
| `.claude-plugin/plugin.json` | ✅ manifest | — | Claude plugin 定義 |
| `plugin.json` | — | ✅ manifest | `agy` plugin 定義（含 `displayName`） |
| `AGENTS.md` | ⚠️ 非 plugin 元件 | ⚠️ 非 plugin 元件 | 跨專案全域規則；`scripts/run.sh` 軟連結至偵測到的 Claude／`agy`／Codex／Pi／Hermes／Grok 全域規則路徑 |
| `output-styles/` | ✅ output style | — | `executive_summary.md` 本體；安裝後以 `default:executive_summary` 出現，`plugins/default/config/settings/claude.json` 的 `outputStyle` 指名它 |
| `rules/` | — | ✅ 啟用時載入的 rule | `output-style.md` 軟連結至 `output-styles/executive_summary.md` |
| `skills/` | ✅ | ✅ | `<name>/SKILL.md`，兩邊格式相同，唯一可直接共用的元件 |
| `agents/` | ✅ subagent | ✅ subagent | `<name>.md`，含 frontmatter；Grok 也讀，以 `default:<name>` 出現 |
| `commands/` | ✅ slash command | ❌ 未支援 | `<name>.md`，以 `/default:<name>` 叫用 |
| `hooks/` | ✅ `hooks/hooks.json` | ✅ 腳本本體 | `hooks/hooks.json` 是 Claude 格式（Grok 也讀這份）；`notify-done.sh` 是三家共用的腳本，見下方 Hooks |
| `hooks.json` | — | ✅ plugin root | `agy` 只從 plugin root 讀 `hooks.json`，且格式不同（具名 hook，`Stop` 為扁平清單），呼叫同一支腳本 |
| `config/settings/` | ⚠️ 非 plugin 元件 | ⚠️ 非 plugin 元件 | Claude（`claude*.json` 各 provider 完整設定）、`agy`、Codex、Grok、Tokscale 的 CLI 設定檔；由 `scripts/run.sh` 軟連結至各家目錄，`scripts/aliases.sh` 與 `~/bin/claude*` 以 `--settings` 指名 |
| `config/statusline/` | ⚠️ 非 plugin 元件 | ⚠️ 非 plugin 元件 | 兩邊 plugin 都無法宣告 statusline，由 `scripts/run.sh` 軟連結，見下表 |

`agy` 欄的 ❌ 表示本機 8 個 `agy` plugin 皆未出現該元件，未經官方文件證實。MCP server 屬檔案而非資料夾：Claude 讀 `.mcp.json`、`agy` 讀 `mcp_config.json`，格式不同需各放一份。

空資料夾以 `.gitkeep` 保留。`不要`在元件資料夾內放 `README.md`：`agents/`、`commands/`、`rules/`、`output-styles/` 會把任何 `*.md` 當成元件載入。

### Hooks

| Hook | 事件 | 行為 |
| :--- | :--- | :--- |
| `markdownlint.sh` | `PostToolUse`（寫入或編輯 `.md` 後） | 檢查全域是否已安裝 `markdownlint`（未安裝則執行 `npm install -g markdownlint-cli`），先以 `--fix` 自動修復（優先讀取專案本身的 `.markdownlint*`，否則套用 `hooks/markdownlint.json` 關閉 `MD013`），Claude 與 Grok 若仍有未修掉的錯誤會透過 `additionalContext` 回饋給模型；`agy` 僅靜默修復並回傳 `{}` |
| `notify-done.sh` | `Stop`（每回合最終回應完成） | 發 macOS 通知：標題 `<CLI> <專案資料夾名>`、內文為最終回應的第 1 行（`executive_summary` 規定的單行 recap；其他回應取第一個非標題文字行），上限 120 字；同一 CLI 同一目錄的新通知會取代舊的 |
| `focus-window.sh` | 點擊通知 | 把執行該 CLI 的視窗切到前景：宿主是 Antigravity IDE 時，從 IDE 的 `storage.json` 找出包含 CLI 工作目錄的最深層已開啟資料夾並聚焦該視窗；其他宿主 (以 `__CFBundleIdentifier` 辨識) 只把該 app 切到前景 |

`notify-done.sh` 在三家交付最終回應的方式不同，由腳本分流：Claude 讀 stdin `last_assistant_message`；Grok（以 `$GROK_HOOK_EVENT` 辨識）讀 `lastAssistantMessage`，略過 session 結束時的 `Stop`；`agy`（參數 `agy`）讀 `transcriptPath` 最後一筆有文字的 `PLANNER_RESPONSE`，略過 `USER_CANCELED`。依賴 macOS 內建的 `jq`、`open` 與 `terminal-notifier`（`brew install terminal-notifier`，`osascript` 的通知無法設定點擊動作）；非 macOS 直接略過。`markdownlint.sh` 同樣以三家分流讀取 `PostToolUse` 檔案路徑（Claude 讀 `tool_input.file_path`、Grok 讀 `tool_input`／`toolInput`、Antigravity 讀 `toolCall.args.Path`／`TargetFile`）。

### Statusline

| 檔案 | 使用者 | `scripts/run.sh` 連結目標 |
| :--- | :--- | :--- |
| `claude-statusline.json` | Claude Code（`ccstatusline`） | `~/.config/ccstatusline/settings.json` |
| `claude-statusline.md` | 人類 | — （`ccstatusline` 說明） |
| `agy-statusline` | Antigravity | `~/.local/bin/agy-statusline` |
| `grok-statusline` | Grok Build | `~/.grok/statusline.sh` |

## 安裝 (Install)

[`scripts/run.sh`](../../scripts/run.sh) 會處理 Claude 與 `agy`；Grok 由 marketplace 安裝：

- Claude Code：登錄本 repo 為 marketplace `shuk-cc-plugin` 並安裝 `default`；手動安裝為 `/plugin install default@shuk-cc-plugin`，再於 `/config` 選 `default:executive_summary`。
- Antigravity：軟連結至 `~/.gemini/config/plugins/default`，以 `agy plugin enable default` / `agy plugin disable default` 開關。
- Grok：`~/.grok/config.toml` 已登錄 marketplace `cc-plugin` 並啟用 `default`；Grok 安裝的是`複本`，改動後執行 `grok plugin update default`。

## 參考資源 (References)

- [VoltAgent/awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents)：社群整理的 Claude Code subagent 範本集，新增 `agents/` 元件前可先參考現成定義
