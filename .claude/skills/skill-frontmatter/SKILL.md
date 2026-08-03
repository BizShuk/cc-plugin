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
- 欄位名稱拼字必須精確（`versio`、`allowed_tools` 之類的錯字會讓該欄位被靜默忽略，
  不會報錯）；改完自行以 `python3 -c "import yaml,sys;yaml.safe_load(...)"` 之類的方式驗一次
- Rules-style frontmatter（`trigger: always_on` + `globs` + `scope`）不是本 workspace
  的慣例，不得與上表三個 tier 混用
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
