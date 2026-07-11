# Claude 插件元資料與 PM2 註冊實作計畫

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 將 `creating-plugin-metadata` 技能重新命名為 `claude-plugin-metadata`，新增 `marketplace.json` 說明與實體範本，將其註冊至通用插件，並在工作區註冊表註冊 `pm2` 插件。

**Architecture:** 透過重命名技能目錄、更新 `SKILL.md` 內容、新增 `marketplace.json` 範例檔案，並更新 `plugins/general/` 的元資料與 README，最後在工作區的 `marketplace.json` 中註冊 `pm2` 插件以完成對接。

**Tech Stack:** `JSON`, `Markdown`

## Global Constraints
- 技能採用 `agentskills.io` 規範，YAML frontmatter 必須包含 `name` 與 `description`。
- Mermaid 邊線文字必須雙引號包覆。
- 不使用粗體，一律以 `backtick` 強調。
- 繁體中文為主，術語以 local language 搭配英文圓括號。

---

### Task 1: 技能目錄重命名與新增 marketplace.json 範本

**Files:**
- Create: `plugins/general/skills/claude-plugin-metadata/marketplace.json`
- Modify: `plugins/general/skills/creating-plugin-metadata` -> `plugins/general/skills/claude-plugin-metadata`

**Interfaces:**
- Consumes: None
- Produces: `plugins/general/skills/claude-plugin-metadata` 目錄與其中的 `marketplace.json`

- [ ] **Step 1: 重新命名技能目錄**
  執行命令將目錄更名：
  `mv plugins/general/skills/creating-plugin-metadata plugins/general/skills/claude-plugin-metadata`

- [ ] **Step 2: 建立 marketplace.json 範本**
  在 `plugins/general/skills/claude-plugin-metadata/marketplace.json` 寫入以下內容：
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

- [ ] **Step 3: 驗證 JSON 語法**
  執行：`jq . plugins/general/skills/claude-plugin-metadata/marketplace.json`
  預期：輸出格式化後的 JSON，且無錯誤。

- [ ] **Step 4: Commit 變更**
  ```bash
  git add plugins/general/skills/claude-plugin-metadata/marketplace.json
  git commit -m "feat: rename skill directory and add marketplace.json template"
  ```

---

### Task 2: 擴充並更新 SKILL.md

**Files:**
- Modify: `plugins/general/skills/claude-plugin-metadata/SKILL.md`

- [ ] **Step 1: 修改 SKILL.md 內容**
  修改 `plugins/general/skills/claude-plugin-metadata/SKILL.md` 的內容，更新 frontmatter、標題、與新增 `marketplace.json` 的說明：
  ```markdown
  ---
  name: claude-plugin-metadata
  description: >
      Use when authoring, initializing, or updating plugin.json and marketplace.json manifests for a Claude Code workspace plugin or individual skills. Triggers on: "create plugin", "update plugin", "plugin.json", "marketplace.json", "register skill".
  version: "1.0.0"
  allowed-tools: Read, Write, Bash
  user-invocable: true
  disable-model-invocation: false
  effort: medium
  context: fork
  metadata:
      type: reference
      platforms: [macos, linux]
  ---

  # Claude Plugin Metadata

  ## Overview
  A guide for authoring the `.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json` manifest files to define and register custom plugins, skills, and agents.

  ## When to Use
  - Initializing a new workspace plugin configuration.
  - Defining plugin name, description, author, repository, and keywords.
  - Registering or updating plugins in the workspace-wide marketplace.json.
  - Structuring sub-components (skills, agents) for discovery.

  When NOT to use:
  - Generating general Go project configuration or system settings.

  ## Manifest Schemas

  ### 1. Plugin Metadata Manifest (plugin.json)
  The plugin metadata manifest must be placed at `.claude-plugin/plugin.json` relative to the plugin root directory.

  #### Structure
  ```json
  {
    "name": "plugin-name",
    "version": "1.0.0",
    "description": "Short description of what the plugin does",
    "author": {
      "name": "Author Name",
      "email": "email@example.com",
      "url": "https://github.com/username"
    },
    "homepage": "https://github.com/username/repository",
    "repository": "https://github.com/username/repository",
    "license": "MIT",
    "keywords": [
      "keyword1",
      "keyword2"
    ],
    "agents": [],
    "skills": []
  }
  ```

  #### Key Fields
  - `name`: Unique name identifier for the plugin (lowercase, hyphens/numbers allowed).
  - `version`: SemVer formatted version string (e.g., `1.0.0`).
  - `description`: Plain text description summarizing the plugin's capabilities.
  - `author`: Nested object containing author contact information.
  - `homepage` / `repository`: URL links to project pages or git repositories.
  - `keywords`: Array of tags for discoverability.
  - `skills` / `agents`: Array of paths to specific skill/agent configurations if explicit loading is required; leave empty for auto-discovery by folder convention (i.e. `skills/` and `agents/` directories).

  ### 2. Workspace Marketplace Registry (marketplace.json)
  The workspace plugin marketplace manifest must be placed at `.claude-plugin/marketplace.json` relative to the workspace root directory.

  #### Structure
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

  #### Key Fields
  - `name`: Name identifier for the marketplace workspace.
  - `owner`: Nested object containing owner name and email.
  - `plugins`: Array of plugin registrations.
    - `name`: The name of the registered plugin.
    - `source`: The source location of the plugin. Can be a relative path or a github source object.

  ## Common Mistakes
  - Putting the manifest in the root directory directly (e.g. `plugin.json` instead of `.claude-plugin/plugin.json`).
  - Using uppercase letters or special characters in the `name` field.
  - Forgetting to include contact info under the `author` field.
  - Forgetting to update the `marketplace.json` when adding a new local plugin folder.
  ```

- [ ] **Step 2: Commit 變更**
  ```bash
  git add plugins/general/skills/claude-plugin-metadata/SKILL.md
  git commit -m "docs: update SKILL.md for claude-plugin-metadata"
  ```

---

### Task 3: 註冊技能與更新插件資訊

**Files:**
- Modify: `plugins/general/.claude-plugin/plugin.json`
- Modify: `plugins/general/README.md`

- [ ] **Step 1: 註冊技能至 general plugin**
  編輯 `plugins/general/.claude-plugin/plugin.json`，將 `claude-plugin-metadata` 加入 `skills` 陣列（按字母排序）：
  ```json
    "skills": [
      "./skills/claude-plugin-metadata",
      "./skills/daily-summary",
      "./skills/markdownlint",
      "./skills/sort-todo"
    ]
  ```

- [ ] **Step 2: 驗證 general plugin.json 語法**
  執行：`jq . plugins/general/.claude-plugin/plugin.json`
  預期：輸出格式化後的 JSON，且無錯誤。

- [ ] **Step 3: 更新 plugins/general/README.md**
  修改 `plugins/general/README.md` 的技能表格：
  將技能清單中的說明補上，並將檔案結構樹的說明修正為 `4 個技能目錄`：
  ```markdown
  ## 技能 (Skills)

  | Skill | 用途 |
  | --- | --- |
  | `claude-plugin-metadata` | 建立與更新 Claude 插件與技能的元資料（含 plugin.json 與 marketplace.json） |
  | `daily-summary` | 彙整過去 24h 跨來源工作，產生工作日報並寫入 Apple Notes |
  | `markdownlint` | Markdown 格式檢查（精選 rule + CUSTOM-01 no-bold），所有插件的 `.md` 檔通用 |
  | `sort-todo` | 排序並格式化待辦清單 |
  ```

  同時修改結構樹說明：
  ```markdown
  ├── skills/              # 4 個技能目錄
  ```

- [ ] **Step 4: Commit 變更**
  ```bash
  git add plugins/general/.claude-plugin/plugin.json plugins/general/README.md
  git commit -m "feat: register claude-plugin-metadata skill and update general plugin README"
  ```

---

### Task 4: 在工作區 marketplace.json 註冊 pm2 插件

**Files:**
- Modify: `.claude-plugin/marketplace.json`

- [ ] **Step 1: 註冊 pm2**
  編輯 `.claude-plugin/marketplace.json`，在 `plugins` 陣列中新增 `pm2` 的註冊資訊：
  ```json
          {
              "name": "pm2",
              "source": "../tmp/pm2"
          },
  ```

- [ ] **Step 2: 驗證 marketplace.json 語法**
  執行：`jq . .claude-plugin/marketplace.json`
  預期：輸出格式化後的 JSON，且無錯誤。

- [ ] **Step 3: Commit 變更**
  ```bash
  git add .claude-plugin/marketplace.json
  git commit -m "feat: register pm2 plugin in workspace marketplace.json"
  ```
