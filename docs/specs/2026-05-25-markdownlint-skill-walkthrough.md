# Markdown 格式檢查技能建立驗證報告 (Markdown Linting Skill Creation Walkthrough)

本報告記錄了基於 `markdownlint` (v0.40.0) 規則為 `general` 插件新增 `markdownlint` 技能的實作與驗證結果，並於後續使用 `markitdown` 提取官方完整規則更新技能，且使用專案本地的 `npx markdownlint` 命令行工具。

## 已完成的變更 (Changes Made)

### 格式檢查技能與規則更新

- 建立了 [SKILL.md](../../plugins/general/skills/markdownlint/SKILL.md) 技能檔案。
  - 技能 frontmatter 配置 `user-invocable: false`，表示僅能由模型調用。
  - 描述中明訂在生成或修改 `plugins/general/skills` 目錄下的 Markdown 檔案時觸發。
  - 使用 `markitdown` 工具從官方下載並提取完整的 `markdownlint` (v0.40.0) 規則，將原本的規則表擴充為 54 條規則（包含 53 條標準規則與 1 條專案自訂規則 `CUSTOM-01`）。
  - 所有規則的描述均符合專案格式規範（不使用 `bold` 語法，一律改用 `backtick` 進行標記）。
  - 於 `SKILL.md` 中新增了 `Usage` 區段，指引如何使用本地 `npx markdownlint` 工具執行檢查，避免全域依賴。
  - 於 `Quick Reference Table` 中增加了 `Status` 欄位，標明各規則啟閉狀態。除了 `MD013` 被設為 `off` 之外，其餘規則（包含 `CUSTOM-01`）皆為 `on`，該狀態與專案的 [.markdownlint.json](../../.markdownlint.json) 設定完全一致。

### 本地命令安裝與驗證

- 移除了系統中的全域 `markdownlint-cli`。
- 在專案本地安裝了 `markdownlint-cli` 以提供執行檔。
- 成功驗證以 `npx markdownlint` 的安裝與運作，版本為 `0.48.0`。
- 於專案根目錄建立了 [.markdownlint.json](../../.markdownlint.json) 設定檔，停用了 `MD013` (line-length) 規則以適應長表格和說明文字。

### 插件與市場設定 (Plugin & Marketplace Configurations)

- 修正了 [plugin.json](../../plugins/general/.claude-plugin/plugin.json) 註冊資訊：
  - 在 `skills` 陣列中按字母順序添加了 `"./skills/markdownlint"`。
  - 在 `keywords` 中加入了 `"markdownlint"` 關鍵字以利探索。
- 修正了 [marketplace.json](../../.claude-plugin/marketplace.json) 市場註冊資訊：
  - 在 `general` 插件的 `skills` 列表中將 `"./skills/anti-sabotage-skill.md"` 更正為目錄型 `"./skills/anti-sabotage"`。
  - 移除了不存在的 `"./skills/superpower"`。
  - 新增了 `"./skills/markdownlint"`，並依字母順序重新排序，使其與插件的主設定檔完全一致。

---

## 驗證結果 (Validation Results)

### 自動化測試

- 執行 `go test ./...` 順利通過，所有既有測試未受影響。
```text
ok  	github.com/bizshuk/cc-plugin/cmd	(cached)
ok  	github.com/bizshuk/cc-plugin/cmd/export	(cached)
```

### 命令與規則校驗驗證

- 執行 `npx markdownlint --version` 回傳：
```text
0.48.0
```
- 執行 `npx markdownlint plugins/general/skills/markdownlint/SKILL.md` 通過，無任何語法錯誤。

### 格式規範手動檢驗

- 已檢查新增的 [SKILL.md](../../plugins/general/skills/markdownlint/SKILL.md) 與 `plans/` 底下的所有計畫文件，確認全數使用 `backtick` 作為高亮標記，完全無使用雙星號 `**` 粗體格式。
