# 建立 Markdown 格式檢查技能實作計畫 (Create Markdownlint Skill Implementation Plan)

本計畫旨在新增一個名稱為 `markdownlint` 的技能，該技能基於 `markdownlint` (v0.40.0) 規則。此技能專為模型調用設計，無須使用者手動調用，並在生成或修改 `plugins/general/skills` 目錄下的 Markdown 檔案時自動觸發，以確保所有 Markdown 技能文件皆符合標準格式規範。

## 使用者審查項目 (User Review Required)

> [!IMPORTANT]
> - 技能建立：建立新技能檔案 `plugins/general/skills/markdownlint/SKILL.md`，內容整合 `markdownlint` (v0.40.0) 的主要規則。
> - 僅模型調用：將 `user-invocable` 設定為 `false`，並在描述中指定觸發時機。
> - 插件註冊：更新 `plugins/general/.claude-plugin/plugin.json` 中的 `skills` 列表，加入新技能的註冊。
> - 格式規範：依據專案規範，所有產出的 Markdown 內容（包括計畫與技能本體）皆不使用 `bold` 語法，改用 `backtick` 進行高亮。

## 開放性問題 (Open Questions)

> [!NOTE]
> 目前無開放性問題。

---

## 預期變更內容 (Proposed Changes)

### 格式檢查技能與註冊 (Linting Skill & Registration)

#### [NEW] [SKILL.md](../../plugins/general/skills/markdownlint/SKILL.md)

- 建立 `plugins/general/skills/markdownlint/SKILL.md`。
- 設定 `YAML frontmatter`：
  - `name`: `markdownlint`
  - `description`: 描述其僅在生成或修改 `plugins/general/skills` 目錄下的 Markdown 檔案時觸發。
  - `user-invocable`: `false`
  - `disable-model-invocation`: `false`
- 內容：整理並呈現 `markdownlint` (v0.40.0) 中核心且常見的規格（如 `MD001`、`MD009`、`MD012`、`MD022`、`MD025`、`MD031`、`MD032`、`MD040` 、 `MD041` 、 `MD047` 等），並針對本專案特有的 `Don't use bold, use backtick` 規則進行特別約束。

#### [MODIFY] [plugin.json](../../plugins/general/.claude-plugin/plugin.json)

- 在 `skills` 陣列中新增 `"./skills/markdownlint"`。

#### [MODIFY] [marketplace.json](../../.claude-plugin/marketplace.json)

- 同步 `general` 插件的 `skills` 陣列，將 `"./skills/anti-sabotage-skill.md"` 更正為目錄型 `"./skills/anti-sabotage"`，移除已不存在的 `"./skills/superpower"`，並新增 `"./skills/markdownlint"`，同時依字母排序排列，使其與 `general` 插件的設定檔完全一致。

---

## 驗證計畫 (Verification Plan)

### 自動化測試
- 使用 `go test ./...` 執行既有的 Go 測試，確認專案功能未受損壞。

### 手動驗證
- 驗證 `plugin.json` 的語法正確性。
- 檢查 `SKILL.md` 中的 `YAML frontmatter` 結構。
- 檢查所有產生的 Markdown 檔案，確認沒有使用 `bold` 語法，且高亮皆為 `backtick`。
