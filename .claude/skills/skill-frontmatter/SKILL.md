---
name: skill-frontmatter
description: >
    Use when authoring or reviewing a `SKILL.md` YAML frontmatter in this
    workspace — choosing between the minimal / standard / full tier, naming a
    skill, or writing its description. Triggers on: "skill frontmatter",
    "SKILL.md 規範", "new skill", "frontmatter tier", "allowed-tools",
    "skill 命名規範".
---

# SKILL Frontmatter 規範 (Frontmatter Spec)

YAML frontmatter 分三個 tier，由簡至詳擇一使用：

| Tier       | 必填欄位              | 選填欄位                                                                                                  | 適用情境                       |
| ---------- | --------------------- | --------------------------------------------------------------------------------------------------------- | ------------------------------ |
| `minimal`  | `name`, `description` | —                                                                                                         | 參考文件、靜態知識             |
| `standard` | `name`, `description` | `version`, `allowed-tools`                                                                                | 一般 CLI 工具技能              |
| `full`     | `name`, `description` | `version`, `allowed-tools`, `user-invocable`, `disable-model-invocation`, `effort`, `context`, `metadata` | 需要控制模型呼叫行為的進階技能 |

額外規範：

- `name` 必須與所在子目錄名稱一致，使用 `kebab-case`
- `description` 採用 `>` 折疊式（`>` 或 `|`），長度 ≤ 1024 字元，且`必須`包含觸發詞（"Use when...", "Triggers on..."）
- `versio` 之類的拼字錯誤禁止（CI 將以 `yaml.Unmarshal` 驗證）
- Rules-style frontmatter（`trigger: always_on` + `globs` + `scope`）僅 `consistency` 與 `go-convention` 兩個常駐技能使用，其餘不得混用
- 標準 frontmatter 範例（`full` tier）：

```yaml
---
name: my-skill
description: >
    Use when ... Triggers on: "foo", "bar".
version: "1.0.0"
allowed-tools: Read, Bash, Glob
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: reference
    platforms: [macos, linux]
---
```
