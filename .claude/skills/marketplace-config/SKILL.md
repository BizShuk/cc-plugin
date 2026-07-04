---
name: marketplace-config
description: >
    Use when configuring Claude Code plugin manifests — `marketplace.json`
    (multi-plugin catalog) or `plugin.json` (single plugin). Triggers on:
    "set up a new plugin", "register a skill", "add a plugin entry",
    "fix plugin.json", "marketplace.json schema", "github source format",
    "git-subdir source", "hook not firing", "skill not picked up".
allowed-tools: Read, Bash, Glob
user-invocable: true
disable-model-invocation: false
effort: low
context: fork
metadata:
    type: reference
    platforms: [macos, linux]
---

# marketplace-config

如何撰寫 `marketplace.json`（多插件目錄）與 `plugin.json`（單插件 manifest），給 Claude Code 與 vercel-labs/skills 共用。

## 檔案位置 (File Layout)

```tree
project-root/
├── .claude-plugin/
│   ├── marketplace.json       # 多插件目錄（每專案一份）
│   └── plugin.json            # 單插件 manifest（單插件時用）
└── plugins/<plugin-name>/
    ├── .claude-plugin/
    │   └── plugin.json        # 該插件的 manifest
    ├── hooks/                 # Stop / StopFailure 等（自動發現）
    │   ├── hooks.json
    │   └── <handler>.sh
    ├── agents/                # *.md（自動發現）
    ├── output-styles/         # *.md（自動發現）
    ├── skills/<skill-name>/
    │   └── SKILL.md           # frontmatter + body
    └── README.md

# Project-wide skill（不屬於任何 plugin）
.claude/skills/<skill-name>/SKILL.md
```

`hooks/`、`agents/`、`output-styles/` 都是**自動發現** — 不需要在 `plugin.json` 內列名。

## marketplace.json

最外層：

```json
{
    "name": "my-marketplace",
    "owner": { "name": "...", "email": "..." },
    "plugins": [ ... ]
}
```

每個 plugin 條目：

```json
{
    "name": "unique-plugin-id",
    "source": "./plugins/<id>",
    "description": "...",
    "keywords": ["..."],
    "agents": ["./agents/<name>.md"],
    "skills": [
        "./skills/<skill-name>",
        "./skills/<other>/SKILL.md"
    ]
}
```

`skills[]` 接受兩種路徑：

- `./skills/<dir>` — 指向整個技能目錄（內含 `SKILL.md`）
- `./skills/<file>.md` — 指向獨立 skill 檔案

## source 欄位的三種形式

```jsonc
// 1. 本地相對路徑（最常見）
"source": "./plugins/my-plugin"

// 2. GitHub shorthand — owner/repo 或 github:owner/repo
"source": "owner/my-plugin"

// 3. Object 形式（明確指定來源類型）
"source": {
    "source": "github",
    "repo": "owner/repo",
    "ref": "v1.2.0",       // 選填：branch 或 tag
    "sha": "abc123..."     // 選填：鎖定 commit（比 ref 更強）
}
```

Monorepo 內的單插件（git-subdir）：

```jsonc
"source": {
    "source": "git-subdir",
    "url": "https://github.com/owner/monorepo",
    "path": "plugins/my-plugin",
    "ref": "main"
}
```

Self-hosted git（url 形式）：

```jsonc
"source": {
    "source": "url",
    "url": "https://gitlab.com/team/my-plugin",
    "ref": "main"
}
```

`skills add .` 只解析本地路徑；object 形式需要走 `skills add owner/repo` 或 Claude Code plugin 流程。

## plugin.json（單插件）

```json
{
    "name": "my-plugin",
    "version": "1.0.0",
    "description": "...",
    "author": {
        "name": "...",
        "email": "..."
    },
    "homepage": "...",
    "repository": "...",
    "license": "MIT",
    "keywords": ["..."],
    "skills": ["./skills/<name>"]
}
```

`plugin.json` **不列** `hooks`、`agents`、`output-styles` — 從子目錄自動發現。

## SKILL.md frontmatter（必要）

```yaml
---
name: my-skill              # kebab-case，必須 = 目錄名
description: >              # folded 樣式，≤ 1024 字元
    Use when ... Triggers on: "foo", "bar".
version: "1.0.0"            # 選填
allowed-tools: Read, Bash   # 選填
user-invocable: false       # 選填
disable-model-invocation: false  # 選填
effort: medium              # 選填：low / medium / high
context: fork               # 選填
metadata:
    type: reference         # 選填
    platforms: [macos, linux]
---

# My Skill
...body...
```

強制欄位只有 `name` + `description`。`description` 必須含觸發詞（`Use when...` 或 `Triggers on:...`）。

## 路徑慣例 (Path Conventions)

- 技能路徑必須以 `./` 開頭（`./skills/<name>`）
- 本地 `source` 必須以 `./` 開頭（`./plugins/<name>`）
- 技能目錄名必須 = SKILL.md frontmatter 的 `name:`

## 驗證 (Validation)

```bash
# 結構與型別
jq -e . .claude-plugin/marketplace.json
jq -e . plugins/<name>/.claude-plugin/plugin.json

# 列出 marketplace 所有 source 類型
jq -r '.plugins[] | "  " + .name + ": " + (
  if (.source | type) == "string" then .source
  else (.source.source // "?") + "/" + (.source.repo // "?")
  end
)' .claude-plugin/marketplace.json

# 檢查 ignore 規則
git check-ignore -v .claude/skills/<skill-name>
```

## 常見錯誤 (Common Mistakes)

| 症狀 | 原因 | 解法 |
| --- | --- | --- |
| 技能沒被 Claude Code 抓取 | `name:` ≠ 目錄名 | frontmatter `name` 必須等於父目錄 |
| `source` 被 parser 拒絕 | 物件形式缺少 `source` 區分鍵 | 必須有 `source: "github" \| "url" \| "git-subdir"` |
| `skills add .` 沒列出 github-source 插件 | 預期行為 | `add .` 只處理本地；github-source 走 `add owner/repo` |
| Hook 從未觸發 | `hooks.json` 事件名稱錯誤 | 必須用 `Stop`、`StopFailure`、`UserPromptSubmit` 等 |
| 技能被 gitignore 排除 | 路徑在 `.claude/skills/*` 規則下 | 加 negation 行 `!.claude/skills/<skill>/` |
| 重複追蹤既有技能 | `marketplace.json` 內已有同名技能 | 用 `git ls-files` 與插件 manifest 比對 |
| Frontmatter `description` 缺觸發詞 | 沒寫 `Use when` / `Triggers on` | 加上去 — 否則模型不會自動呼叫 |

## 範例 (Examples)

最小 `marketplace.json`：

```json
{
    "name": "shuk-cc-plugin",
    "plugins": [
        {
            "name": "hello",
            "source": "./plugins/hello",
            "skills": ["./skills/greet"]
        }
    ]
}
```

最小 `plugin.json`：

```json
{
    "name": "hello",
    "version": "1.0.0",
    "skills": ["./skills/greet"]
}
```

含 github-source 與 git-subdir 混合的目錄：

```json
{
    "name": "mixed",
    "plugins": [
        { "name": "local-tool", "source": "./plugins/local-tool" },
        {
            "name": "external",
            "source": { "source": "github", "repo": "owner/external", "sha": "<pinned-sha>" }
        },
        {
            "name": "monorepo-plugin",
            "source": {
                "source": "git-subdir",
                "url": "https://github.com/owner/monorepo",
                "path": "plugins/x"
            }
        }
    ]
}
```

## 與其他技能 (Related Skills)

- `plugin-development` — 開發新插件的全流程（若存在）
- `vercel-skills-cli` — 用 `skills add` / `skills find` 操作
- `claude-code-hooks` — Hook 撰寫細節

## 設計原則

- **最小 manifest**：只列需要的欄位；auto-discover 子目錄
- **明確 source**：本地用 `./plugins/...`，github 用物件形式並鎖 `sha`
- **description 含觸發詞**：模型靠這個判斷何時呼叫技能