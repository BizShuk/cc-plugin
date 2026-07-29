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

產生一份 `package.json` 作為任意專案的**通用任務跑器 (universal task runner)** — 無論 Go、Python、Rust、多服務 monorepo 或混合技術棧皆適用。檔案不包含任何 JS 依賴；其唯一用途是讓 VSCode 的 NPM Scripts 面板與 CI/CD 管線有一個統一、可探索的介面。

核心原則：**頂層腳本 = 管線階段 (pipeline stages)；第二層腳本 = 元件 (components)。**

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

對應標準 CI/CD 管線。每個專案預設包含以下四項；確實不適用者可省略：

| 腳本 (Script) | 用途 | 執行方式 |
| --- | --- | --- |
| `dev` | 啟動本機開發環境 | 平行執行所有 `dev:*` |
| `test` | 執行完整測試套件 | 平行執行所有 `test:*` |
| `build` | 產生建構產物 | 平行執行所有 `build:*` |
| `deploy` | 部署至目標環境 | 平行執行所有 `deploy:*` |

管線相關的額外腳本（僅在需要時加入）：

| 腳本 (Script) | 用途 |
| --- | --- |
| `lint` | 執行所有 linter (`lint:*`) |
| `format` | 執行所有格式化工具 (`format:*`) |
| `clean` | 清除產生的產物 |
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
  "build:web": "cd web && npm run build",

  "deploy:api": "cd api && ./scripts/deploy.sh",
  "deploy:web": "cd web && ./scripts/deploy.sh"
}
```

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
  "deploy": "npx npm-run-all --parallel deploy:*"
}
```

**平行 vs 循序規則：**
- **預設平行**：所有頂層階段 (`dev`, `test`, `build`, `deploy` 等) 預設皆使用平行執行 (`--parallel`)。
- **特例循序**：僅當子項元件有明確的建構順序依賴時，才調整為循序 (`--sequential`)。

**單一元件捷徑：** 若某階段只有一個元件，頂層腳本可直接呼叫，無需 `npm-run-all`：

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
    "build:web": "cd web && npm run build",

    "deploy": "npx npm-run-all --parallel deploy:*",
    "deploy:api": "cd api && ./scripts/deploy.sh",
    "deploy:web": "cd web && ./scripts/deploy.sh"
  }
}
```

**必要欄位：**
- `name` — 專案目錄名稱，kebab-case
- `version` — `"0.0.0"`（非發布套件）
- `private` — `true`（防止意外 `npm publish`）
- `scripts` — 任務定義

**省略欄位：**
- 不加 `main`、`module`、`type` — 這不是 JS 套件
- 不加 `dependencies` 或 `devDependencies` — 工具透過 `npx` 或語言原生指令呼叫
- 除非使用者要求，否則不加 `license`

## 工作流程 (Workflow)

1. **掃描專案** — 辨識元件（含有獨立建構系統的目錄：`go.mod`、`Cargo.toml`、`pyproject.toml`、巢狀 `package.json` 等）
2. **偵測既有指令** — 檢查 `Makefile`、`run.sh`、`scripts/`、既有 `package.json`
3. **對應元件至腳本** — 分配 `<stage>:<component>` 名稱
4. **選擇聚合方式** — 單一元件 = 直接呼叫；多元件 = `npm-run-all`
5. **產生 `package.json`** — 套用範本，填入探索到的指令
6. **保留既有內容** — 若 `package.json` 已存在，僅合併 `scripts`；不得刪除既有欄位

## 語言指令參照 (Language Command Reference)

| 語言 | Dev | Test | Build |
| --- | --- | --- | --- |
| Go | `go run ./cmd/<name>` | `go test ./...` | `go build -o bin/<name> ./cmd/<name>` |
| Python | `python -m <module>` | `pytest` | `python -m build` |
| Rust | `cargo run` | `cargo test` | `cargo build --release` |
| Java/Kotlin | `./gradlew run` | `./gradlew test` | `./gradlew build` |
| C/C++ | `make run` | `make test` | `make build` |
| JS/TS (巢狀) | `npm run dev` | `npm test` | `npm run build` |

## VSCode 整合 (VSCode Integration)

`package.json` 存在於專案根目錄後，VSCode 自動：
- 在 **NPM Scripts** 面板（Explorer 側邊欄）顯示所有腳本
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
