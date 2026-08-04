---
name: project-task
description: >
    Use when a project needs a standardized task runner via package.json — creating
    or updating npm scripts for dev, test, build, deploy across any language or
    stack. Triggers on: "create package.json", "add npm scripts", "project tasks",
    "run scripts", "task runner", "建立專案任務", "npm run".
version: "1.0.0"
allowed-tools: Read, Write, Bash
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: technique
    platforms: [macos, linux]
---

# 專案任務 (Project Task)

## 概要 (Overview)

產生一份 `package.json` 作為任意專案的`通用任務跑器 (universal task runner)` — 無論 Go、Python、Rust、多服務 monorepo 或混合技術棧皆適用。檔案不包含任何 JS 依賴；其唯一用途是讓 VSCode 的 NPM Scripts 面板與 CI/CD 管線有一個統一、可探索的介面。

核心原則：`頂層腳本 = 管線階段 (pipeline stages)；第二層腳本 = 元件 (components)`。

## 使用時機 (When to Use)

- 新專案尚無 `package.json` 時進行初始化
- 為非 JS 專案加入統一的 `dev`/`test`/`build`/`deploy` 進入點
- 在多語言 monorepo 中統一任務執行方式
- 讓專案相容 VSCode 內建的 NPM Scripts 面板

不適用情境：
- 專案本身是純 JS/TS 應用，已有完整含依賴的 `package.json`
- 使用者明確要求 Makefile、Taskfile 或 Justfile

## 腳本層級 (Script Hierarchy)

### 第一層 — 管線階段 (Pipeline Stages)

頂層類別 (top entry category) 對應標準 CI/CD 管線，每個類別皆可帶第二層子任務 (`<stage>:<component>`)。每個專案預設包含以下項目；確實不適用者可省略：

| 腳本 (Script) | 用途 | 執行方式 |
| --- | --- | --- |
| `dev` | 啟動本機開發環境 | 平行執行所有 `dev:*` |
| `test` | 執行完整測試套件 | 平行執行所有 `test:*` |
| `build` | 安裝依賴並產生建構產物（含 `npm install` 等安裝步驟） | 平行執行所有 `build:*` |
| `deploy` | 部署至目標環境 | 平行執行所有 `deploy:*` |
| `lint` | 執行所有 linter 與格式化工具（不另設 `format`，格式化屬於 lint） | 平行執行所有 `lint:*` |
| `clean` | 移除產出物 (artifacts)：output、build files、`node_modules` | 平行執行所有 `clean:*` |
| `destroy` | 終止執行中的服務 (kill running services) | 平行執行所有 `destroy:*` |
| `run` | 執行專案腳本 (project script)，腳本通常位於 `scripts/` 目錄 | 不聚合 — 個別執行 `run:<script>` |

管線相關的額外腳本（僅在需要時加入）：

| 腳本 (Script) | 用途 |
| --- | --- |
| `ci` | 完整 CI 管線：`lint && test && build` |

### 第二層 — 元件 (Components)

第二層腳本針對特定元件，與程式語言無關：

```
<stage>:<component>
```

範例：

```json
{
  "dev:api": "cd api && go run ./cmd/server",
  "dev:web": "cd web && npm run dev",
  "dev:worker": "cd worker && python -m worker",

  "test:api": "cd api && go test ./...",
  "test:web": "cd web && npm test",
  "test:e2e": "cd e2e && npx playwright test",

  "build:api": "cd api && go build -o bin/server ./cmd/server",
  "build:web": "cd web && npm install && npm run build",

  "deploy:api": "cd api && ./scripts/deploy.sh",
  "deploy:web": "cd web && ./scripts/deploy.sh",

  "lint:api": "cd api && gofmt -l . && go vet ./...",
  "lint:web": "cd web && npx eslint . && npx prettier --check .",

  "clean:api": "rm -rf api/bin",
  "clean:web": "rm -rf web/dist web/node_modules",

  "destroy:api": "pkill -f 'bin/server' || true",
  "destroy:web": "pkill -f 'vite' || true",

  "run:seed": "./scripts/seed.sh",
  "run:migrate": "./scripts/migrate.sh",
  "run:backfill": "python scripts/backfill.py"
}
```

`run` 類別的第二層以`腳本名稱`而非元件命名（`run:<script>`），對應 `scripts/` 目錄下的檔案；各腳本獨立執行，頂層不設 `run` 聚合腳本。

### 第三層 — 變體 (Variants)（少用）

僅當元件有明確不同的執行模式時使用：

```
<stage>:<component>:<variant>
```

範例：`test:api:unit`、`test:api:integration`、`deploy:web:staging`、`deploy:web:prod`

## 聚合模式 (Aggregation Pattern)

頂層腳本透過 `npm-run-all` 聚合其子項，預設皆平行執行 (`--parallel`)：

```json
{
  "dev": "npx npm-run-all --parallel dev:*",
  "test": "npx npm-run-all --parallel test:*",
  "build": "npx npm-run-all --parallel build:*",
  "deploy": "npx npm-run-all --parallel deploy:*",
  "lint": "npx npm-run-all --parallel lint:*",
  "clean": "npx npm-run-all --parallel clean:*",
  "destroy": "npx npm-run-all --parallel destroy:*"
}
```

`平行 vs 循序規則`：
- `預設平行`：所有頂層類別 (`dev`, `test`, `build`, `deploy`, `lint`, `clean`, `destroy` 等) 預設皆使用平行執行 (`--parallel`)。
- `特例循序`：僅當子項元件有明確的建構順序依賴時，才調整為循序 (`--sequential`)。
- `run 不聚合`：`run:*` 各腳本互不相關，只個別執行（`npm run run:seed`），不建立頂層 `run` 聚合腳本。

`單一元件捷徑`：若某階段只有一個元件，頂層腳本可直接呼叫，無需 `npm-run-all`：

```json
{
  "build": "go build -o bin/server ./cmd/server"
}
```

## package.json 範本 (Template)

```json
{
  "name": "<project-name>",
  "version": "0.0.0",
  "private": true,
  "description": "<一行專案描述>",
  "scripts": {
    "dev": "npx npm-run-all --parallel dev:*",
    "dev:api": "cd api && go run ./cmd/server",
    "dev:web": "cd web && npm run dev",

    "test": "npx npm-run-all --parallel test:*",
    "test:api": "cd api && go test ./...",
    "test:web": "cd web && npm test",

    "build": "npx npm-run-all --parallel build:*",
    "build:api": "cd api && go build -o bin/server ./cmd/server",
    "build:web": "cd web && npm install && npm run build",

    "deploy": "npx npm-run-all --parallel deploy:*",
    "deploy:api": "cd api && ./scripts/deploy.sh",
    "deploy:web": "cd web && ./scripts/deploy.sh",

    "lint": "npx npm-run-all --parallel lint:*",
    "lint:api": "cd api && gofmt -l . && go vet ./...",
    "lint:web": "cd web && npx eslint .",

    "clean": "npx npm-run-all --parallel clean:*",
    "clean:api": "rm -rf api/bin",
    "clean:web": "rm -rf web/dist web/node_modules",

    "destroy": "npx npm-run-all --parallel destroy:*",
    "destroy:api": "pkill -f 'bin/server' || true",
    "destroy:web": "pkill -f 'vite' || true",

    "run:seed": "./scripts/seed.sh",
    "run:migrate": "./scripts/migrate.sh"
  }
}
```

`必要欄位`：
- `name` — 專案目錄名稱，kebab-case
- `version` — `"0.0.0"`（非發布套件）
- `private` — `true`（防止意外 `npm publish`）
- `scripts` — 任務定義

`省略欄位`：
- 不加 `main`、`module`、`type` — 這不是 JS 套件
- 不加 `dependencies` 或 `devDependencies` — 工具透過 `npx` 或語言原生指令呼叫
- 除非使用者要求，否則不加 `license`

## 工作流程 (Workflow)

1. `掃描專案` — 辨識元件（含有獨立建構系統的目錄：`go.mod`、`Cargo.toml`、`pyproject.toml`、巢狀 `package.json` 等）
2. `偵測既有指令` — 檢查 `Makefile`、`run.sh`、`scripts/`、既有 `package.json`
3. `對應元件至腳本` — 分配 `<stage>:<component>` 名稱；`scripts/` 目錄下不屬於管線階段的獨立腳本對應為 `run:<script>`
4. `選擇聚合方式` — 單一元件 = 直接呼叫；多元件 = `npm-run-all`
5. 產生 `package.json` — 套用範本，填入探索到的指令
6. `保留既有內容` — 若 `package.json` 已存在，僅合併 `scripts`；不得刪除既有欄位

## 語言指令參照 (Language Command Reference)

主要語言為 Golang、Node.js 相關 (tsc、vite 等) 與 Python：

| 語言 | Dev | Test | Build（含安裝） | Lint | Clean |
| --- | --- | --- | --- | --- | --- |
| Go | `go run ./cmd/<name>` | `go test ./...` | `go build -o bin/<name> ./cmd/<name>` | `gofmt -l . && go vet ./...` | `rm -rf bin` |
| Node.js (tsc) | `npx tsc --watch` | `npm test` | `npm install && npx tsc` | `npx eslint .` | `rm -rf dist node_modules` |
| Node.js (vite) | `npm run dev` | `npm test` | `npm install && npm run build` | `npx eslint .` | `rm -rf dist node_modules` |
| Python | `python -m <module>` | `pytest` | `pip install -e . && python -m build` | `ruff check .` | `rm -rf dist build *.egg-info` |

其他語言（Rust、Java 等）依相同模式類推，不在預設參照範圍。

## VSCode 整合 (VSCode Integration)

`package.json` 存在於專案根目錄後，VSCode 自動：
- 在 `NPM Scripts` 面板（Explorer 側邊欄）顯示所有腳本
- 提供點擊即執行功能
- 支援 `Ctrl/Cmd+Shift+P` → "Tasks: Run Task" → "npm"

無需額外設定。

## 常見錯誤 (Common Mistakes)

| 錯誤 | 修正 |
| --- | --- |
| 加入 `"main": "index.js"` | 省略 `main` — 這不是 JS 進入點 |
| 指令放錯層級 | 管線階段 (pipeline stages) 在頂層；元件 (components) 在第二層 |
| 頂層用 `&&` 串接元件 | 使用 `npm-run-all` 以取得正確的平行處理與錯誤處理 |
| 忘記加 `"private": true` | 務必設定以防止意外 `npm publish` |
| 腳本中使用絕對路徑 | 使用專案根目錄的相對路徑搭配 `cd` |
| 加入 `node_modules` 依賴 | 保持無依賴；透過 `npx` 使用 `npm-run-all` |
| 另設 `format` 頂層腳本 | 格式化屬於 lint — 併入 `lint:*`，不獨立成類別 |
| `build` 未含依賴安裝 | `build:*` 應包含安裝步驟（如 `npm install`），確保 clean checkout 可直接建構 |
| 用 `clean` 停服務或用 `destroy` 刪檔案 | `clean` 只移除 artifacts；`destroy` 只終止執行中的服務 |
| 把 `scripts/` 的獨立腳本塞進 `dev`/`build` | 不屬於管線階段的腳本歸入 `run:<script>` |
