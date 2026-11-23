# 插件目錄說明 (Plugins Directory Overview)

`plugins/` 收納五個本地 plugin，其中 `default` 同時相容 Claude Code、Antigravity (`agy`) 與 Grok Build；外部 plugin 若有引用，由根目錄
`.claude-plugin/marketplace.json` 以 GitHub `owner/repo_name` 宣告。即使以
submodule checkout 到本地，`source` 仍寫 repo，不用 `./plugins/<name>`。

## 本地插件 (Local Plugins)

| Plugin | 職責 |
| :--- | :--- |
| `default` | 雙相容預設插件：output style / rule (`executive_summary.md`)、Stop 完成通知 hooks、全域規則 (`AGENTS.md`)、CLI 設定與 statusline，並整合通用技能 (`content-summarizer`、`daily-summary`、`markdownlint`、`markitdown`、`mermaid`、`system-one-two-composer`) 與 `feature` subagent |
| `experiment` | 候選 skill 沙盒 (`anti-sabotage`、`firecrawl`、`model-evaluator`、`playwright-cli`、`scrapling`、`summarize-sh` 與外部 skill) |
| `god` | 系統法則 (`system-laws`)、LLM mechanics (`llm-mechanics`)、領域探索 (`domain-exploration`) 與五大通用算子 (`universal-*`) |
| `review` | 審查、規劃、正典文件治理 (`project-docs`)、成長目標審查 (`review-target`)、復盤 (`session-retro`)、教學 (`tutorial`)、自演化 (`auto-evolving`) 與 `review-coordinator` 代理 |
| `team` | Agent team 規劃 (`team-design`)、角色提示生成 (`role-generator`) 與 12 份跨職能範例角色 (`roles/`) |

五個本地 plugin 都有 `.claude-plugin/plugin.json` 與 `README.md`，並登錄於
`.claude-plugin/marketplace.json`；`default` 另有 Antigravity 原生根目錄 `plugin.json` 與 `hooks.json`，由 `scripts/run.sh` 軟連結至 `~/.gemini/config/plugins/default`，並為 Claude Code 註冊與安裝 `default@shuk-cc-plugin`。預設 `skills/`／`agents/` 與 submodule 由目錄
自動探索；manifest 的 `skills`／`agents` 只宣告非預設路徑或外部來源。詳見根目錄
`AGENTS.md` 的「檔案職責」與「慣例」章節。

## 外部 Marketplace Plugins

根 `.claude-plugin/marketplace.json` 目前`沒有`引用外部 GitHub plugin，僅發布 `./plugins/` 內的本地插件。
來源與順序以 `.claude-plugin/marketplace.json` 為唯一真理來源。

## Git Submodules

`plugins/` 底下`沒有` submodule。全 repo 的 gitlink 皆位於 `pkg/` 且記於 `.gitmodules`：

- `pkg/prompt/system-prompts/CL4R1T4S`
- `pkg/tools/career-ops`（目前於工作區標記為刪除）

初始化與更新：

```bash
git submodule update --init --recursive
```

`.gitignore` 排除 `plugins/experiment/skills/summarize.sh/` 作為 skill 依賴的本機 clone 路徑，不得記錄為 submodule。
