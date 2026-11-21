# openproject

來源: [opf/openproject](https://github.com/opf/openproject) (GPL-3.0, ~17.5k stars, default branch `dev`) · [website](https://www.openproject.org) · [docs](https://www.openproject.org/docs/)

開源的專案管理平台, 定位是 Jira / MS Project / Asana 的 self-host 替代品: work package (OpenProject 對 task / bug / feature 的統稱), agile board, Gantt, roadmap, 工時與成本, wiki, 會議記錄放在同一個 web app, 主打 `資料主權` (可完全自架). 分 Community (免費), Enterprise (SSO / LDAP / SCIM 與支援) 與 BIM (營建業) 三個 edition.

這個 repo `不是 Claude skill`: 沒有 `SKILL.md`, 是一個完整的 Rails web app. 列在這裡是因為它 `內建 MCP server`, 自架後 agent 可以直接讀寫專案與 work package; 也可作為專案管理工具選型參考.

## 架構

| 部分                                       | 用途                                                                                                                                                                                                          |
| ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `app/` `lib/` `config/` `db/`              | Ruby on Rails 核心 (Ruby 4.0), PostgreSQL 為必要資料庫                                                                                                                                                        |
| `frontend/`                                | TypeScript; 新 UI 改用 server-rendered Hotwire (Turbo + Stimulus) + GitHub Primer (ViewComponent), 舊 Angular 元件逐步遷成 custom element                                                                     |
| `modules/`                                 | 以 Rails engine 形式拆出的功能模組: `backlogs` `boards` `gantt` `team_planner` `meeting` `costs` `budgets` `wikis` `storages` (Nextcloud) `github_integration` `gitlab_integration` `bim` `webhooks` `llm` 等 |
| `app/services/mcp_tools/` `mcp_resources/` | 內建 MCP server 的 tool 與 resource 實作                                                                                                                                                                      |
| `modules/llm/`                             | LLM connection 管理: API key, 預設 model, 功能綁定, semantic index (語意搜尋索引)                                                                                                                             |
| `docs/api/apiv3/`                          | REST API v3 (HAL+JSON) 與 OpenAPI spec                                                                                                                                                                        |

核心設計:

- 一切以 `work package` 為中心: task, bug, milestone, phase 都是同一個實體的不同 type, 以 relation (blocks / follows / parent) 串成 Gantt 與依賴圖.
- 功能以 module 外掛, 核心保持精簡; 各 edition 的差異靠 module 與 enterprise token 開關.
- repo 自帶 `AGENTS.md` / `CLAUDE.md`, 規範 AI coding agent 的開發流程 (預設零註解, lefthook 檢查, RSpec).

## MCP Server

Endpoint 是 instance 的 `/mcp` (例如 `https://op.example.com/mcp`), 在 `Administration -> Artificial Intelligence (AI) -> Model Context Protocol (MCP)` 開關與調整.

- 認證: 個人用走 `API token` (需開啟 "Enable API tokens", token 會過期); 多人共用 / web client 走 `OAuth`, 在 `Administration -> Authentication -> OAuth applications` 建立含 `mcp` scope 的 confidential application, 或用相容的 OpenID Connect provider 發 token.
- Tools: `search_work_packages` `create_work_package` `update_work_package` `get_work_package_form` `create_work_package_comment` `list_work_package_comments` `create/update/delete/list_work_package_relation` `search_projects` `search_portfolios` `search_programs` `search_users` `search_versions` `search_labels` `search_custom_fields` `list_statuses` `list_types` `current_user`.
- Resources: project, work package, user, status, type, version, label, custom field.
- 每個 tool / resource 可單獨停用或改名 (例如把 "work package" 改叫 "work item", 讓 LLM 認得團隊用語).
- Response format 三選一: `Full` (text + structured, 最相容但吃 token), `Structured content only`, `Content only`.

Claude Code 接法 (以 API token 為例):

```bash
claude mcp add --transport http openproject https://op.example.com/mcp \
  --header "Authorization: Bearer <API_TOKEN>"
```

## 安裝

- 試用: [14 天 cloud trial](https://start.openproject.com/).
- Self-host: 官方提供 Docker (all-in-one 或 compose) 與 Linux package, 見 [installation guides](https://www.openproject.org/download-and-installation/).
- 開發環境: 見 repo 內 `docs/development/development-environment/` (native 或 `docker/dev/`).

## 適用情境

- 想要 Jira 等級的專案管理但資料要留在自己的 server (例如部署到 liva).
- 讓 agent 透過 MCP 查詢 / 建立 / 更新 work package, 把 TODO 或 incident 直接開成 ticket.
- 參考大型 Rails monorepo 如何以 module 拆功能, 以及如何把既有 REST API 包成 MCP server (output filter 精簡回應, 可由 admin 調整 tool 描述).

## 注意

- `GPL-3.0`: 散佈修改版必須公開原始碼; 只在內部自架使用則無此義務.
- MCP 的 `Bearer` header 會把 API token 寫進 MCP client 設定, 不要提交到 repo; 共用情境改用 OAuth.
- AI 文件 (`docs/system-admin-guide/ai/`) 中 text assistant 與總覽頁仍標 "to be updated", LLM 相關功能仍在發展中.
- MCP 是否需 Enterprise token: admin controller 內未見 enterprise gating (推論, 未實際部署驗證).
- Claude Code 指令中的 `--header` 寫法依 MCP CLI 慣例撰寫 (推論, 未對實際 instance 驗證).
