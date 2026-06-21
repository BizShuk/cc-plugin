# Output Styles

Claude Code 的 `outputStyle` 自訂風格檔案，提供三種切換選擇。

## 本資料夾內容

| 檔案          | 風格 | 長度       | 適合情境                     |
| ------------- | ---- | ---------- | ---------------------------- |
| `brief.md`    | 精簡 | 1–10 行    | 日常查詢、簡單任務、結論先行 |
| `detailed.md` | 詳盡 | 5 段固定   | 學習、設計、文件、複雜 bug   |
| `review.md`   | 審查 | 結構化清單 | Code Review、重構前評估      |

## 放置位置 — User vs Project

Claude Code 在兩個位置讀取 `output-styles/`，優先序：User → Project。

### 1. User-level（全域，所有專案生效）

```text
~/.claude/output-styles/
```

- 影響所有專案
- 適合：跨專案常用的通用風格（如 `brief`）
- 不會被 git 追蹤

### 2. Project-level（單一專案）

```text
<project>/.claude/output-styles/
```

- 只影響當前專案
- 適合：團隊共用的審查/設計風格
- 隨 git 追蹤、可分享

### 優先序範例

| User 有 `brief.md` | Project 有 `review.md` | 結果                     |
| ------------------ | ---------------------- | ------------------------ |
| ✅                 | ❌                     | User 的 `brief` 生效     |
| ❌                 | ✅                     | Project 的 `review` 生效 |
| ✅                 | ✅                     | Project 覆蓋 User        |

## 切換方式

1. **對話中切換**：`/config` → Output style → 選風格
2. **手動設值** — 編輯 `~/.claude/settings.json`：

    ```json
    {
        "outputStyle": "brief"
    }
    ```

3. **Project 預設** — 編輯 `.claude/settings.json`（提交到 git）：

    ```json
    {
        "outputStyle": "review"
    }
    ```

## 風格檔格式

每個 `*.md` 必須含 YAML frontmatter：

```yaml
---
name: my-style
description: 風格說明
---
# 風格本體
規則、偏好、長度限制等
```

`name` 會出現在 `/config` 選單中；`description` 為輔助說明。

## 新增風格

1. 決定放置位置（User / Project）
2. 建立 `<dir>/<name>.md`
3. 重新執行 `/config` 即可看到新選項

## 注意事項

- 同名時 Project 覆蓋 User
- 修改後需重新選擇風格才會生效
- 此資料夾目前位於 `output-styles/` 而非標準的 `.claude/output-styles/`，若 Claude Code 版本未支援根目錄讀取，請改放到 `.claude/output-styles/` 或建立 symlink：

    ```bash
    ln -s ../output-styles .claude/output-styles
    ```
