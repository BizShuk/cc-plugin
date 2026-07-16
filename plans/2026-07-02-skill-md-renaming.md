# Feature: 插件技能檔案命名規範 (Plugin Skill Naming Convention)

## 問題陳述 (Problem Statement)

[plugins/experiment/skills/anti-sabotage/anti-sabotage.md](../plugins/experiment/skills/anti-sabotage/anti-sabotage.md) 不符合 `agentskills.io` 規範的 `SKILL.md` 命名，應改為 `SKILL.md` 並補上標準 YAML frontmatter。

## 受影響檔案 (Affected Files)

- `plugins/experiment/skills/anti-sabotage/anti-sabotage.md` → `plugins/experiment/skills/anti-sabotage/SKILL.md`

## 實作步驟 (Implementation Steps)

1. 將 `plugins/experiment/skills/anti-sabotage/anti-sabotage.md` 重新命名為 `SKILL.md`
2. 在檔案頂部新增符合 `agentskills.io` 規範的 YAML frontmatter（`name`、`description` 必填，建議加上 `version` 與 `allowed-tools`）
3. 確認 `plugins/experiment/.claude-plugin/plugin.json` 的 `skills` 陣列仍正確指向此技能
4. 確認 `.claude-plugin/marketplace.json` 的 `keywords` 與描述同步更新

## 驗證方式 (Verification)

- 執行 `npx skills add .` 確認無異常
- 手動測試 `cc-plugin` 啟動後，`/anti-sabotage` 技能仍可正常觸發
- `markdownlint` 與 YAML 解析無錯誤
