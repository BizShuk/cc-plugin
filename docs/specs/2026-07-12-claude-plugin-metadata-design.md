# 設計文件：Claude 插件元資料 (Claude Plugin Metadata) 與 PM2 註冊

本文件規劃將既有的 `creating-plugin-metadata` 技能重新命名為 `claude-plugin-metadata`，同時擴充其內容以支援建立與更新 Claude 插件與技能 (包含 `plugin.json` 與 `marketplace.json`)。此外，將在工作區級註冊表 `.claude-plugin/marketplace.json` 中註冊 `pm2` 本地插件。

## 變更項目

### 1. 技能目錄重命名
- 原目錄：`plugins/general/skills/creating-plugin-metadata`
- 新目錄：`plugins/general/skills/claude-plugin-metadata`

### 2. 擴充技能說明 (SKILL.md)
- 修改 frontmatter 的 `name` 為 `claude-plugin-metadata`。
- 修改 `description` 為：`Use when authoring, initializing, or updating plugin.json and marketplace.json manifests for a Claude Code workspace plugin or individual skills.`。
- 擴充內文，新增 `marketplace.json` 的結構、欄位說明與 schema 指引。
- 說明如何在建立與更新插件時，同時維護這兩個元資料檔案。

### 3. 新增實體範本 (marketplace.json)
- 於 `plugins/general/skills/claude-plugin-metadata/` 底下新增 `marketplace.json` 檔案。
- 內容結構如下：
```json
{
  "name": "workspace-marketplace",
  "owner": {
    "name": "Developer Name",
    "email": "developer@example.com"
  },
  "plugins": [
    {
      "name": "plugin-name",
      "source": "./plugins/plugin-name"
    },
    {
      "name": "external-plugin",
      "source": {
        "source": "github",
        "repo": "username/repository"
      }
    }
  ]
}
```

### 4. 註冊該技能至通用插件
- 於 `plugins/general/.claude-plugin/plugin.json` 的 `skills` 陣列中新增 `"./skills/claude-plugin-metadata"`。

### 5. 註冊 PM2 插件至工作區
- 於 `.claude-plugin/marketplace.json` 的 `plugins` 陣列中新增 `pm2`：
```json
        {
            "name": "pm2",
            "source": "../tmp/pm2"
        }
```

### 6. 更新說明文件
- 更新 `plugins/general/README.md` 的 `skills` 表格，將 `creating-plugin-metadata` 替換為 `claude-plugin-metadata` 並修改用途說明，同時更新檔案結構樹。

## 驗證計畫
1. 檢查 `plugins/general/.claude-plugin/plugin.json` 與 `.claude-plugin/marketplace.json` 的 JSON 格式無誤。
2. 確認重新命名後，Marksman LSP 診斷無錯誤。
3. 執行本專案測試，確認一切功能正常。
