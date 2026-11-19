# default — 雙相容預設插件 (Default Plugin)

同時相容 Claude Code 與 Antigravity (`agy`) 的預設插件。兩邊共用同一份內容，平台專屬入口一律用軟連結指回本體，免除重複維護。

## 資料夾相容性 (Folder Compatibility)

| 路徑 | Claude Code | Antigravity (`agy`) | 說明 |
| :--- | :--- | :--- | :--- |
| `.claude-plugin/plugin.json` | ✅ manifest | — | Claude plugin 定義 |
| `plugin.json` | — | ✅ manifest | `agy` plugin 定義（含 `displayName`） |
| `AGENTS.md` | ⚠️ 非 plugin 元件 | ⚠️ 非 plugin 元件 | 跨專案全域規則；`scripts/run.sh` 軟連結至偵測到的 Claude／`agy`／Codex／Pi／Hermes／Grok 全域規則路徑 |
| `output-styles/` | ✅ output style | — | `executive_summary.md` 本體；安裝後以 `default:executive_summary` 出現，`plugins/default/config/settings/claude.json` 的 `outputStyle` 指名它 |
| `rules/` | — | ✅ 啟用時載入的 rule | `output-style.md` 軟連結至 `output-styles/executive_summary.md` |
| `skills/` | ✅ | ✅ | `<name>/SKILL.md`，兩邊格式相同，唯一可直接共用的元件 |
| `agents/` | ✅ subagent | ❌ 未支援 | `<name>.md`，含 frontmatter |
| `commands/` | ✅ slash command | ❌ 未支援 | `<name>.md`，以 `/default:<name>` 叫用 |
| `hooks/` | ✅ `hooks/hooks.json` | ❌ 未支援 | 事件 hook（例 Stop, PreToolUse） |
| `config/settings/` | ⚠️ 非 plugin 元件 | ⚠️ 非 plugin 元件 | Claude（`claude*.json` 各 provider 完整設定）、`agy`、Codex、Grok、Tokscale 的 CLI 設定檔；由 `scripts/run.sh` 軟連結至各家目錄，`scripts/aliases.sh` 與 `~/bin/claude*` 以 `--settings` 指名 |
| `config/statusline/` | ⚠️ 非 plugin 元件 | ⚠️ 非 plugin 元件 | 兩邊 plugin 都無法宣告 statusline，由 `scripts/run.sh` 軟連結，見下表 |

`agy` 欄的 ❌ 表示本機 8 個 `agy` plugin 皆未出現該元件，未經官方文件證實。MCP server 屬檔案而非資料夾：Claude 讀 `.mcp.json`、`agy` 讀 `mcp_config.json`，格式不同需各放一份。

空資料夾以 `.gitkeep` 保留。`不要`在元件資料夾內放 `README.md`：`agents/`、`commands/`、`rules/`、`output-styles/` 會把任何 `*.md` 當成元件載入。

### Statusline

| 檔案 | 使用者 | `scripts/run.sh` 連結目標 |
| :--- | :--- | :--- |
| `claude-statusline.json` | Claude Code（`ccstatusline`） | `~/.config/ccstatusline/settings.json` |
| `claude-statusline.md` | 人類 | — （`ccstatusline` 說明） |
| `agy-statusline` | Antigravity | `~/.local/bin/agy-statusline` |
| `grok-statusline` | Grok Build | `~/.grok/statusline.sh` |

## 安裝 (Install)

[`scripts/run.sh`](../../scripts/run.sh) 會同時處理兩邊：

- Claude Code：登錄本 repo 為 marketplace `shuk-cc-plugin` 並安裝 `default`；手動安裝為 `/plugin install default@shuk-cc-plugin`，再於 `/config` 選 `default:executive_summary`。
- Antigravity：軟連結至 `~/.gemini/config/plugins/default`，以 `agy plugin enable default` / `agy plugin disable default` 開關。
