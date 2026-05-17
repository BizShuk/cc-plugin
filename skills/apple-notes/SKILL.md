---
name: apple-notes
description: "基於終端的筆記管理工具 (A command-line notes tool)，支援標籤 (tags)、附件 (attachments)，並透過 SQLite 或 Markdown 檔案儲存。"
homepage: https://github.com/harperreed/memo
metadata:
    {
        "openclaw":
            {
                "emoji": "📝",
                "os": ["darwin", "linux"],
                "requires": { "bins": ["memo"] },
                "install":
                    [
                        {
                            "id": "brew",
                            "kind": "brew",
                            "formula": "harperreed/tap/memo",
                            "bins": ["memo"],
                            "label": "透過 Homebrew 安裝 memo (Install memo via Homebrew)"
                        }
                    ]
            }
    }
---

# 蘋果筆記 (apple-notes / memo)

使用 `memo` 可以在終端機 (terminal) 直接管理筆記 (notes)。此工具以 Markdown 為主，支援建立 (create)、檢視 (view)、編輯 (edit)、刪除 (delete)、搜尋 (search)，並可加入標籤 (tags) 與附件 (attachments)。

## 初始設定 (Setup)

- 透過 Homebrew 安裝 (Install via Homebrew): `brew install harperreed/tap/memo`
- 透過原始碼安裝 (Install from source): `go install github.com/harperreed/memo/cmd/memo@latest`
- 透過 Github 下載二進位檔 (Download binary): 參考 [GitHub Releases](https://github.com/harperreed/memo/releases)

## 新增筆記 (Add Notes)

- 使用預設編輯器撰寫: `memo add "筆記標題 (Note Title)"`
- 直接提供內容: `memo add "Quick thought" --content "Remember to call mom"`
- 從檔案匯入: `memo add "Article Draft" --file draft.md`
- 加入標籤: `memo add "Project Ideas" --content "..." --tags "work,brainstorm"`
- 使用目前目錄標籤: `memo add "Project TODO" --content "..." --here`

## 檢視與搜尋筆記 (List and View Notes)

- 列出最近筆記: `memo list`
- 依標籤過濾: `memo list --tag work`
- 搜尋筆記內容 (支援 FTS5 模糊搜尋): `memo list --search "meeting"`
- 限制顯示數量: `memo list --limit 5`
- 列出包含目前目錄標籤的筆記: `memo list --here`
- 檢視特定筆記 (使用 ID 前綴，至少 6 個字元): `memo show abc123`

## 編輯與刪除筆記 (Edit and Delete Notes)

- 編輯筆記: `memo edit abc123`
- 刪除筆記: `memo rm abc123`
- 略過確認刪除: `memo rm abc123 --force`

## 標籤管理 (Manage Tags)

- 新增標籤到特定筆記: `memo tag add abc123 important`
- 移除標籤: `memo tag rm abc123 important`
- 列出所有標籤: `memo tag list`

## 附件管理 (Attachments)

- 附加檔案: `memo attach abc123 document.pdf`
- 提取附件: `memo attach get def456 --output ./downloads/`

## 匯出與匯入 (Export and Import)

- 匯出所有筆記為 JSON 格式: `memo export --format json --output backup.json`
- 匯出至 Markdown 資料夾: `memo export --format md --output ./notes/`
- 從 JSON 匯入: `memo import backup.json`
- 從 Markdown 資料夾匯入: `memo import ./notes/`

## 儲存後端 (Storage)

支援兩種儲存後端 (Storage backends)：

- **Markdown** (新使用者預設): 筆記以 markdown 檔案儲存於 `~/.local/share/memo/`
- **SQLite**: 筆記儲存於 `~/.local/share/memo/memo.db`

可於 `~/.config/memo/config.json` 進行設定：

```json
{
    "backend": "markdown",
    "data_dir": "~/.local/share/memo"
}
```

可透過 `memo migrate --to markdown` 或 `memo migrate --to sqlite` 來切換儲存模式。

## MCP 伺服器 (MCP Server)

內建 MCP 伺服器，可供 AI 助理 (AI assistant) 整合使用：
`memo mcp`

支援的工具 (Tools) 包含：`add_note`、`list_notes`、`get_note`、`update_note`、`delete_note`、`search_notes`、`add_tag`、`remove_tag`、`add_attachment`、`list_attachments`、`get_attachment`、`export_note` 等。
