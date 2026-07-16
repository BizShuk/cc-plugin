# 插件目錄說明 (Plugins Directory Overview)

`plugins/` 只收納八個本地 plugin；外部 plugin 由根目錄 `.claude-plugin/marketplace.json` 直接引用 GitHub，不假設存在同名本地資料夾。

## 本地插件 (Local Plugins)

| Plugin | 職責 |
| :--- | :--- |
| `experiment` | 候選 skill 沙盒 |
| `explore` | 摘要、轉檔、專案探索與路由 |
| `general` | 通用 skill、feature agent、hooks 與 output styles |
| `god` | 系統法則、LLM mechanics 與通用 operators |
| `review` | 審查、規劃、自演化與 review coordinator |
| `team` | Agent team、角色與 orchestration 設計 |
| `tools` | Apple Calendar、Mail、Notes、Reminders CLI 整合 |
| `ultra-explore` | 可續跑、可驗證的知識庫建構管道 |

每個本地 plugin 都有 `.claude-plugin/plugin.json` 與 `README.md`。`skills`／`agents` manifest 欄位保持空陣列，由標準目錄自動探索。

## 外部 Marketplace Plugins

根 marketplace 目前引用：`pm2`、`awesome-claude-code-subagents`、`gosdk`、`inf`、`superpowers`、`ip-incubation`、`understand-anything`、`last30days`、`ui-ux-pro-max-skill`。來源與順序以 `.claude-plugin/marketplace.json` 為唯一真理來源。

## Git Submodules

`.gitmodules` 只保留實際 gitlink：

- `pkg/system-prompts/CL4R1T4S`
- `pkg/tools/career-ops`

初始化與更新：

```bash
git submodule update --init --recursive
```

`plugins/explore/skills/summarize.sh/summarize` 等 ignored local clone 不屬於 `.gitmodules`，不得記錄為正式 submodule。
