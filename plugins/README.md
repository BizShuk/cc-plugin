# 插件目錄說明 (Plugins Directory Overview)

`plugins/` 只收納七個本地 plugin；外部 plugin 由根目錄
`.claude-plugin/marketplace.json` 直接引用 GitHub，不假設存在同名本地資料夾。

## 本地插件 (Local Plugins)

| Plugin | 職責 |
| :--- | :--- |
| `experiment` | 候選 skill 沙盒 |
| `explore` | 摘要、轉檔、專案探索與路由 |
| `general` | 通用 skill、feature agent、hooks 與 output styles |
| `god` | 系統法則、LLM mechanics 與通用 operators |
| `review` | 審查、規劃、自演化與 review coordinator |
| `team` | Agent team、角色與 orchestration 設計 |
| `ultra-explore` | 可續跑、可驗證的知識庫建構管道 |

每個本地 plugin 都有 `.claude-plugin/plugin.json` 與 `README.md`，七個都登錄於
`.claude-plugin/marketplace.json`。本地 skill／agent 由目錄自動探索；manifest 的
`skills`／`agents` 只宣告歸屬於該 plugin 的`外部`來源。詳見根目錄 `CLAUDE.md`
的「檔案職責」與「慣例」章節。

## 外部 Marketplace Plugins

根 marketplace 目前引用五個 GitHub 來源：`gosdk`、`inf`、`ip-incubation`、`pm2`、`tools`。
`tools` 提供 Apple Calendar、Mail、Notes、Reminders CLI 整合。
來源與順序以 `.claude-plugin/marketplace.json` 為唯一真理來源。

## Git Submodules

`plugins/` 底下`沒有` submodule。`superpowers` 改由 superset VSCode plugin 安裝，
不再 vendored 進本 repo。全 repo 的 gitlink 只剩兩個，皆在 `pkg/` 且記於 `.gitmodules`：

- `pkg/prompt/system-prompts/CL4R1T4S`
- `pkg/tools/career-ops`

初始化與更新：

```bash
git submodule update --init --recursive
```

`plugins/experiment/skills/summarize.sh/` 是 skill 依賴的本機 clone，已由
`.gitignore` 排除，不得記錄為 submodule。
