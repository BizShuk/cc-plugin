# Understand Anything — AI 程式碼庫理解插件 (Codebase Knowledge Graph Plugin)

> 把任何程式碼庫轉成可探索的互動式知識圖譜 (knowledge graph),用多代理管線 (multi-agent pipeline) 結合 Tree-sitter 靜態分析 + LLM 語意理解,產出可在儀表板中瀏覽的 `.understand-anything/knowledge-graph.json`。

## 用途 (What It Does)

`Understand Anything` 是 [Egonex-AI/Understand-Anything](https://github.com/Egonex-AI/Understand-Anything) 的本地整合版本 (以 submodule 形式收錄),提供 9 個 skill 與 9 個 agent,覆蓋從初次分析、互動探索到 onboarding 導引的完整工作流。

```tree
進入新 codebase?
  ├─ /understand            完整 7 階段分析,產出圖譜
  ├─ /understand-dashboard  開啟互動式視覺化儀表板
  ├─ /understand-chat       對圖譜提問 (Q&A)
  ├─ /understand-diff       變更影響範圍分析
  ├─ /understand-explain    深入解釋單一檔案/函式
  ├─ /understand-onboard    產出新人 onboarding 指南
  ├─ /understand-domain     業務領域 + 流程萃取
  └─ /understand-knowledge  解析 Karpathy-pattern LLM wiki 知識庫
```

## 套件結構 (Package Layout)

```tree
plugins/understand-anything/
├── .claude-plugin/                        Claude marketplace manifest
├── understand-anything-plugin/            Claude Code plugin 本體
│   ├── .claude-plugin/plugin.json        真正的 plugin manifest (name: understand-anything)
│   ├── agents/                            9 個 agent 定義 (Markdown prompt)
│   │   ├── project-scanner.md             掃描 + 語言/框架偵測
│   │   ├── file-analyzer.md               Tree-sitter + LLM 解析
│   │   ├── assemble-reviewer.md           合併後修補 cross-batch edges
│   │   ├── architecture-analyzer.md       架構分層
│   │   ├── tour-builder.md                教學路線生成
│   │   ├── graph-reviewer.md              LLM 完整驗證
│   │   ├── domain-analyzer.md             業務領域/流程
│   │   ├── article-analyzer.md            Wiki 隱含關聯萃取
│   │   └── knowledge-graph-guide.md       引導查詢已產生的圖
│   ├── skills/                            9 個 skill 入口
│   │   ├── understand/SKILL.md            主分析管線 (7 phases)
│   │   ├── understand-chat/               對圖譜 Q&A
│   │   ├── understand-dashboard/          啟動 React 儀表板
│   │   ├── understand-diff/               diff 影響分析
│   │   ├── understand-explain/            深入解釋
│   │   ├── understand-onboard/            onboarding 指南
│   │   ├── understand-domain/             業務領域
│   │   ├── understand-knowledge/          知識庫圖譜化
│   │   └── understand/                    7 階段管線核心
│   ├── src/                               TypeScript skill 邏輯
│   │   ├── context-builder.ts
│   │   ├── diff-analyzer.ts
│   │   ├── explain-builder.ts
│   │   ├── onboard-builder.ts
│   │   └── understand-chat.ts
│   ├── packages/                          pnpm workspace
│   │   ├── core/                          共用分析引擎 (types, persistence, tree-sitter, search, schema, tours)
│   │   └── dashboard/                     React + React Flow + Zustand + TailwindCSS v4
│   └── hooks/                             auto-update post-commit hook
├── install.sh / install.ps1               跨 15+ AI CLI 平台一鍵安裝
├── package.json                           pnpm workspace 根 (private)
└── README.md                              上游原始說明
```

## 核心管線 (`/understand` 7 Phases)

`skills/understand/SKILL.md` 編排完整管線。每個階段的產物都寫到 `.understand-anything/intermediate/`,**不丟回 context**,避免大型 codebase 撐爆 context window。

```tree
Phase 0   Pre-flight        解析參數 / git commit / worktree 重定向 / 語系偵測
Phase 0.5 Ignore             產生/讀取 .understandignore
Phase 1   SCAN              project-scanner 產出 scan-result.json + importMap
Phase 1.5 BATCH             compute-batches.mjs 計算語意批次
Phase 2   ANALYZE           file-analyzer (5 並行, 20-30 檔/批) → batch-*.json
Phase 3   ASSEMBLE REVIEW   assemble-reviewer 修補合併後的圖
Phase 4   ARCHITECTURE      architecture-analyzer 分層 (API/Service/Data/UI/...)
Phase 5   TOUR              tour-builder 5-15 步教學路線
Phase 6   REVIEW            inline 確定性驗證 或 LLM graph-reviewer
Phase 7   SAVE              寫入 knowledge-graph.json + 指紋基準
```

## 雙軌混合分析 (Hybrid Analysis)

```tree
原始碼
   │
   ├─► Tree-sitter WASM (deterministic, 可重現)
   │     ├─ imports / exports map
   │     ├─ function / class 定義
   │     ├─ call sites, inheritance
   │     └─ 指紋 (fingerprint) → 增量更新
   │
   └─► LLM agent (semantic, 因模型而異)
         ├─ plain-English 摘要 / tags
         ├─ 架構分層
         ├─ 業務領域對應
         ├─ 教學路線
         └─ 語言概念標註 (generics, closure, decorator, ...)
```

`★ Insight ─────────────────────────────────────`
結構性 side **可重現** (同樣程式碼永遠產出同樣 edges) — 這是為什麼結構可以拿來做 fingerprint;語意 side 捕捉意圖 (檔案「為何」存在,而不只是 import 什麼) — 這部分才需要 LLM 介入。兩者**互不重疊**,所以可以平行驗證。
`─────────────────────────────────────────────────`

## 知識圖譜 Schema

### 13 種節點類型 (Node Types)

| Type       | 說明                                     | ID Convention                              |
| ---------- | ---------------------------------------- | ------------------------------------------ |
| `file`     | 原始碼檔案                               | `file:<relative-path>`                     |
| `function` | 函式/方法                                | `function:<relative-path>:<name>`          |
| `class`    | 類別/介面/型別                           | `class:<relative-path>:<name>`             |
| `module`   | 邏輯模組/套件                            | `module:<name>`                            |
| `concept`  | 抽象概念/模式                            | `concept:<name>`                           |
| `config`   | 設定檔 (YAML, JSON, TOML, env)           | `config:<relative-path>`                   |
| `document` | 文件檔 (Markdown, RST, TXT)              | `document:<relative-path>`                 |
| `service`  | 可部署服務定義 (Dockerfile, K8s)         | `service:<relative-path>`                  |
| `table`    | 資料表/遷移                              | `table:<relative-path>:<table-name>`       |
| `endpoint` | API 端點/路由                            | `endpoint:<relative-path>:<endpoint-name>` |
| `pipeline` | CI/CD 管線                               | `pipeline:<relative-path>`                 |
| `schema`   | Schema 定義 (GraphQL, Protobuf, Prisma)  | `schema:<relative-path>`                   |
| `resource` | 基礎設施資源 (Terraform, CloudFormation) | `resource:<relative-path>`                 |

### 26 種邊類型 (Edge Types, 7 類)

| 類別                      | 邊類型                                                     |
| ------------------------- | ---------------------------------------------------------- |
| 結構 (Structural)         | `imports`, `exports`, `contains`, `inherits`, `implements` |
| 行為 (Behavioral)         | `calls`, `subscribes`, `publishes`, `middleware`           |
| 資料流 (Data flow)        | `reads_from`, `writes_to`, `transforms`, `validates`       |
| 依賴 (Dependencies)       | `depends_on`, `tested_by`, `configures`                    |
| 語意 (Semantic)           | `related`, `similar_to`                                    |
| 基礎設施 (Infrastructure) | `deploys`, `serves`, `provisions`, `triggers`              |
| Schema/資料 (Schema/Data) | `migrates`, `documents`, `routes`, `defines_schema`        |

### 邊權重 (Edge Weight Conventions)

| Edge Type                                                  | Weight        |
| ---------------------------------------------------------- | ------------- |
| `contains`                                                 | 1.0           |
| `inherits`, `implements`                                   | 0.9           |
| `calls`, `exports`, `defines_schema`                       | 0.8           |
| `imports`, `deploys`, `migrates`                           | 0.7           |
| `depends_on`, `configures`, `triggers`                     | 0.6           |
| `tested_by`, `documents`, `provisions`, `serves`, `routes` | 0.5           |
| 其他                                                       | 0.5 (default) |

## 增量更新機制 (Incremental)

靠兩種指紋判斷是否需要重分析:

1. **Git commit hash** — 寫在 `.understand-anything/meta.json`,跨 commit 比對
2. **Tree-sitter 結構指紋** — `build-fingerprints.mjs` 產出 baseline,改空白行不算 STRUCTURAL 變更

`/understand --auto-update` 會註冊 post-commit hook,每次 commit 自動只重分析有結構變更的檔案,大型 codebase 增量時間可省 ~157k tokens / ~158s。

```bash
/understand --auto-update      # 啟用 auto-update
/understand --no-auto-update   # 停用
/understand --full             # 強制完整重建
/understand --review           # 走 LLM 完整驗證 (預設是 inline 確定性驗證)
/understand --language zh      # 產出繁中/簡中內容
/understand src/frontend       # 限定子目錄
```

`★ Insight ─────────────────────────────────────`
**Worktree 陷阱**:`Phase 0` 自動偵測是否在 git worktree 內,若是會把 `.understand-anything/` 重定向到**主 repo root**。worktree 是暫存的,直接寫進去會在 session 結束時整個圖譜跟著丟。
`─────────────────────────────────────────────────`

## 儀表板 (Dashboard)

`packages/dashboard/` 是 React + TypeScript 視覺化前端:

```tree
React 18 + TypeScript strict
├── React Flow          ← 圖譜節點/邊渲染
├── Zustand             ← 狀態管理
├── TailwindCSS v4      ← 樣式
├── prism-react-renderer ← 程式碼檢視器
└── Vite                ← 打包
```

**設計風格**:深黑奢華 (`#0a0a0a` + `#d4a574` 金色),DM Serif Display 字型,75% 圖譜 + 360px 右側欄。

**互動**:

- 節點點擊 → plain-English 解釋 + 關係 + 引導學習
- Persona-Adaptive UI (junior dev / PM / power user 看到不同細節層級)
- 程式碼檢視從底部滑入,可展開成全螢幕 modal
- 側欄兩個 tab:`Info` (專案總覽 → 選節點時 NodeInfo → Learn persona 顯示 LearnPanel) + `Files` (樹狀檔案瀏覽器)
- Source viewer 透過 dev server `/file-content.json` 端點 + access token + 圖譜推導的 path allowlist 守門

## 跨平台支援 (Multi-Platform)

`install.sh` 支援 15+ AI CLI / IDE 平台:

| 平台                              | 安裝方式                                 |
| --------------------------------- | ---------------------------------------- |
| Claude Code                       | 原生 plugin marketplace                  |
| Cursor                            | 自動發現 (`.cursor-plugin/plugin.json`)  |
| VS Code + Copilot                 | 自動發現 (`.copilot-plugin/plugin.json`) |
| Codex                             | `install.sh codex`                       |
| OpenCode / OpenClaw / Antigravity | `install.sh <platform>`                  |
| Gemini CLI / Pi Agent / Vibe CLI  | `install.sh <platform>`                  |
| Hermes / Cline / KIMI CLI / Trae  | `install.sh <platform>`                  |
| Nanobot / Kiro                    | `install.sh <platform>`                  |

## 使用方式 (Quick Start)

```bash
# 在 Claude Code 安裝
/plugin marketplace add Egonex-AI/Understand-Anything
/plugin install understand-anything

# 在當前專案跑分析
/understand

# 開啟互動式儀表板
/understand-dashboard

# 對程式碼庫提問
/understand-chat 認證流程怎麼走?

# 分析當前變更影響
/understand-diff

# 深入解釋單一檔案
/understand-explain src/auth/login.ts

# 產出 onboarding 指南
/understand-onboard

# 萃取業務領域知識
/understand-domain

# 解析 LLM wiki 知識庫
/understand-knowledge ~/path/to/wiki
```

## 開發注意事項 (Gotchas)

1. **Tree-sitter 用 WASM 版** (`web-tree-sitter`) 而非 native binding — native 版在 darwin/arm64 + Node 24 會失敗
2. **Dashboard import 限制**: 儀表板**只能**從 `core` 的 subpath export (`./search`, `./types`, `./schema`) 引入,主 entry 會拉 Node.js 模組炸瀏覽器
3. **版本同步**: 推送時 5 個檔案必須同步 bump — `understand-anything-plugin/package.json` + `understand-anything-plugin/.claude-plugin/plugin.json` + 專案根的 `.claude-plugin/plugin.json` / `.cursor-plugin/plugin.json` / `.copilot-plugin/plugin.json`
4. **Marketplace.json 不放 version** — `plugins[]` entry 只接受 `name` 和 `source`,加其他欄位會壞掉 schema 驗證
5. **Incremental fingerprint 必須先生效才能寫 meta.json** — 否則 auto-update 會把每個 commit 都誤判為 STRUCTURAL → 永遠 `FULL_UPDATE`
6. **Local plugin 測試**: Claude Code 把 plugin 緩存在 `~/.claude/plugins/cache/...`,symlink 不行 (Search/Glob 工具追不到)。改本地後需 `cp -R ./understand-anything-plugin ~/.claude/plugins/cache/.../`,再開新 session
7. **Agent frontmatter 省略 `model` 欄位** — 讓各平台 fallback 到預設。`inherit` 是 Claude Code 專屬關鍵字,其他平台會當字面 model id 處理然後 `ProviderModelNotFoundError`

## 前置需求 (Prerequisites)

- Node.js ≥ 22 (開發用 v24)
- pnpm ≥ 10 (透過 `packageManager` 欄位固定)
- 被分析專案需為可讀取的目錄

## 輸出目錄 (Output Layout)

跑完 `/understand` 後,被分析的專案根目錄會出現:

```tree
<project-root>/.understand-anything/
├── knowledge-graph.json     ← 最終圖譜 (可 commit 進 git 給團隊用)
├── meta.json                ← git commit hash + 分析時間
├── config.json              ← autoUpdate + outputLanguage
├── .understandignore        ← 排除規則 (gitignore 語法)
├── intermediate/            ← 管線中間產物 (可加入 .gitignore)
│   ├── scan-result.json     ← 保留供 incremental 重用
│   ├── batches.json
│   ├── batch-*.json
│   ├── assembled-graph.json
│   ├── layers.json
│   ├── tour.json
│   └── review.json
└── tmp/                     ← 暫存
```

**最佳實踐**: 把 `.understand-anything/intermediate/` 與 `diff-overlay.json` 加入 `.gitignore`,其餘 commit 進 git — 新人不必重跑管線直接享用 (範例: `GoogleCloudPlatform/microservices-demo`)。>10MB 的圖譜用 `git-lfs`。

## 核心設計原則 (Why It Works)

1. **責任分離** (separation of concerns):靜態分析做擅長的事 (結構),LLM 做擅長的事 (語意)。兩者不重疊,可獨立驗證。
2. **Context 工程** (context engineering):agent 結果寫到磁碟不丟回 context,100 萬行 monorepo 也不會撐爆 context window。
3. **漸進式增強** (progressive enhancement):incremental + fingerprint + worktree redirect + auto-update hook,持續追蹤 codebase 演進。
4. **可移植** (portable):同一個 skill manifest 在 15+ AI CLI 上跑,只依賴 `agentskills.io` frontmatter 規範。
5. **可分享** (shareable):圖譜就一個 JSON,commit 進 git,新人不必重跑管線。

---

`★ Insight ─────────────────────────────────────`
這套不是「AI 讀程式碼給你解釋」,而是「**多代理管線把程式碼變成可重複使用的結構化圖譜資料庫**,AI 與人類都能 query 同一份 ground truth」。
`─────────────────────────────────────────────────`
