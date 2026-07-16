# Skill Metadata 標準化規格 (Skill Metadata Standardization)

Status: `completed`

## 目標 (Goal)

所有本地插件 skill 均以 `plugins/<plugin>/skills/<kebab-name>/SKILL.md` 儲存，frontmatter 可被機械驗證，插件元件由標準目錄自動探索。

## 目錄與命名契約 (Directory and Naming Contract)

- skill 入口固定為獨立子目錄下的 `SKILL.md`。
- frontmatter `name` 必須與目錄一致並使用 `kebab-case`。
- `description` 使用 YAML folded/literal style，長度不超過 1024 字元，且包含 `Use when` 或 `Triggers on`。
- `plugin.json` 的 `skills` 與 `agents` 維持空陣列；Claude Code 由 `skills/`、`agents/` 自動探索，不重複維護檔案清單。
- 每個本地 plugin 都必須有 `README.md`，並提及實際存在的 skill 與 agent。

## 實作結果 (Implementation Result)

- [x] `anti-sabotage/anti-sabotage.md` 改為 `anti-sabotage/SKILL.md` 並補齊 full-tier metadata。
- [x] `session_retro` 改為 `session-retro`。
- [x] experiment 的 `summarize.sh` skill 改為 `summarize-sh`；品牌文字仍保留 `summarize.sh`。
- [x] 修正缺少觸發詞或 folded description 的 frontmatter。
- [x] 新增 `plugins/tools/README.md` 並同步各 plugin README。
- [x] marketplace 補登本地 `god` plugin，並驗證本地與 GitHub source schema。
- [x] 新增 `scripts/pluginmeta`，支援 read-only check 與 `--write` 清理冗餘 manifest 陣列。

## 驗證 (Verification)

```bash
go run ./scripts/pluginmeta
go test ./scripts/pluginmeta -count=1
jq -e . .claude-plugin/marketplace.json plugins/*/.claude-plugin/plugin.json
```

同步工具的責任是掃描、驗證與維持 auto-discovery；不把實際 skill/agent 清單寫入 manifest。
