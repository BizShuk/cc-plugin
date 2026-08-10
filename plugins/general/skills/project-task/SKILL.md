---
name: project-task
description: >
    Use when a project needs a standardized task runner via package.json — creating
    or updating npm scripts for dev, test, build, deploy, release across any
    language or stack. Triggers on: "create package.json", "add npm scripts",
    "project tasks", "run scripts", "task runner", "建立專案任務", "npm run".
version: "1.2.0"
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
- 為非 JS 專案加入統一的 `dev`/`test`/`build`/`deploy`/`release` 進入點
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
| `dev` | 重新建置並在本機拋棄式環境跑起來 | 平行執行所有 `dev:*` |
| `test` | 執行完整測試套件 | 平行執行所有 `test:*` |
| `build` | 安裝依賴並產生建構產物（含 `npm install` 等安裝步驟） | 平行執行所有 `build:*` |
| `deploy` | 重新建置並裝進自己實際在用的環境 | 平行執行所有 `deploy:*` |
| `release` | 送上對外通路 (public channel) | 平行執行所有 `release:*` |
| `lint` | 執行所有 linter 與格式化工具（不另設 `format`，格式化屬於 lint） | 平行執行所有 `lint:*` |
| `clean` | 移除產出物 (artifacts)：output、build files、`node_modules` | 平行執行所有 `clean:*` |
| `destroy` | 終止執行中的服務 (kill running services) | 平行執行所有 `destroy:*` |
| `run` | 執行專案腳本 (project script)，腳本通常位於 `scripts/` 目錄 | 不聚合 — 個別執行 `run:<script>` |

管線相關的額外腳本（僅在需要時加入）：

| 腳本 (Script) | 用途 |
| --- | --- |
| `ci` | 完整 CI 管線：`lint && test && build` |

### 執行目標階梯 (Run Target Ladder)

`dev`、`deploy`、`release` 是同一份產物的三個`落點 (destination)`，不是三個抽象層級。
判準只有一個：`誰會看到這份產物`。

| 腳本 | 落點 | 誰看得到 | 典型 |
| --- | --- | --- | --- |
| `dev` | 本機拋棄式環境 | 只有這次終端機 | iOS Simulator、本機跑起來的 desktop app、dev server |
| `deploy` | 自己實際在用的環境 | 自己，且會留著繼續用 | 已配對的 iPhone、`/Applications`、自架 server |
| `release` | 對外通路 | 別人 | App Store Connect / Mac App Store、npm publish、正式環境 |

`dev 與 deploy 一律重新建置`。這兩條路徑存在的理由就是驗證`剛剛改的東西`；
沿用上一次的產物會讓改動看起來沒生效，而且不會有任何錯誤訊息——這是最難察覺的
假陰性。要「不重建只啟動」是明確的例外，用旗標表達（`--no-build`），不是預設。

`release 一律 Release configuration + archive`。Debug 產物帶著 debug entitlement
與未最佳化的 binary，通路端會直接退件。

`不要把對外上架掛在 deploy`。`deploy` 是`我把它裝到我的裝置上`，`release` 是
`我把它交出去`；兩者的失敗代價差一個數量級，混在同一個名字下遲早會誤觸。
對應的腳本檔名也照這個分野命名：`scripts/release.sh`，不是 `scripts/deploy.sh`。

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

  "release:api": "cd api && ./scripts/release.sh",
  "release:web": "cd web && ./scripts/release.sh",

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

範例：`test:api:unit`、`test:api:integration`、`release:web:staging`、`release:web:prod`

## 聚合模式 (Aggregation Pattern)

頂層腳本透過 `npm-run-all` 聚合其子項，預設皆平行執行 (`--parallel`)：

```json
{
  "dev": "npx npm-run-all --parallel dev:*",
  "test": "npx npm-run-all --parallel test:*",
  "build": "npx npm-run-all --parallel build:*",
  "deploy": "npx npm-run-all --parallel deploy:*",
  "release": "npx npm-run-all --parallel release:*",
  "lint": "npx npm-run-all --parallel lint:*",
  "clean": "npx npm-run-all --parallel clean:*",
  "destroy": "npx npm-run-all --parallel destroy:*"
}
```

`平行 vs 循序規則`：
- `預設平行`：所有頂層類別 (`dev`, `test`, `build`, `deploy`, `release`, `lint`, `clean`, `destroy` 等) 預設皆使用平行執行 (`--parallel`)。
- `特例循序`：僅當子項元件有明確的建構順序依賴時，才調整為循序 (`--sequential`)。
- `run 不聚合`：`run:*` 各腳本互不相關，只個別執行（`npm run run:seed`），不建立頂層 `run` 聚合腳本。

`單一元件規則`：若某階段只有一個元件，頂層腳本`直接呼叫`，`不要`建立第二層名稱：

```json
{
  "build": "go build -o bin/server ./cmd/server"
}
```

第二層後綴的用途是`區分`，不是`標示`。只有一個元件時，`build:mac` 相對於 `build`
不帶任何新資訊，只是同一件事的別名，讀的人還得多跳一層才知道實際指令。
單一平台專案維持 `build`；等真的出現第二個元件再拆。

判準：`拆出 <stage>:<component> 的理由必須是「有另一個 <stage>:<other> 存在」`。

### 呼叫方向 (Call Direction)

`聚合腳本只能被人叫，不能被腳本叫`。`scripts/` 下的 shell 腳本、第二層任務、
CI 設定，一律引用`具體元件` (`build:web`)，`絕不`引用頂層聚合名稱 (`build`)。

理由是聚合腳本的語意會隨元件增減而改變。今天 `build` 等於 web build，
明天多一個 `build:ios`，所有寫著 `npm run build` 的呼叫點就在無預警下
變成跑完整個矩陣——原本幾秒的 web build 變成連 Xcode 一起跑。
若該腳本本身又被 `build:ios` 呼叫，就直接構成`無窮遞迴 (infinite recursion)`。

因此`把既有頂層腳本拆成 fan-out 時，必須同時回頭修正所有呼叫點`：

```bash
# 拆分前：build 就是 web build，腳本這樣寫是對的
grep -rn "npm run build" scripts/

# 拆分後：每一處都要改指向具體元件
npm run build:web
```

漏掉這一步的症狀是「重複建置」或「指令跑不完」，而不是明顯的錯誤訊息，
所以要在拆分的同一次改動裡處理完，不能留到之後。

### 子項自給自足 (Self-Sufficient Components)

當元件 B 的產物`包含`元件 A 的產物（例：iOS bundle 內含 web build），
讓 `build:B` 自帶 `build:A`，使它在 clean checkout 單獨執行也能成功；
頂層則以 `--sequential` 依序跑 `build:A build:B`。

代價是聚合執行時 A 會被建置兩次。這是刻意的取捨：
`子項可獨立執行`比`省掉一次重複建置`重要，因為單獨跑 `build:B` 是日常操作，
跑完整聚合則相對少見。若 A 的建置成本高到無法接受，才改為讓 `build:B`
依賴前一階段的產物，並在文件中明確寫出這個順序前提。

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

    "release": "npx npm-run-all --parallel release:*",
    "release:api": "cd api && ./scripts/release.sh",
    "release:web": "cd web && ./scripts/release.sh",

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
   - 依`執行目標階梯`分辨 `dev` / `deploy` / `release`：模擬器與本機跑 → `dev`；
     實機、`/Applications`、自架環境 → `deploy`；App Store 與其他對外通路 → `release`
4. `選擇聚合方式` — 單一元件 = 直接呼叫且不加後綴；多元件 = `npm-run-all`
5. 產生 `package.json` — 套用範本，填入探索到的指令
6. `保留既有內容` — 若 `package.json` 已存在，僅合併 `scripts`；不得刪除既有欄位
7. `修正呼叫點` — 若本次把既有頂層腳本拆成 fan-out，`grep` 整個 repo（`scripts/`、
   CI 設定、README、進行中的 `plans/`）找出引用舊名稱之處，逐一改指向具體元件

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
| 只有一個元件卻加 `<stage>:<component>` 後綴 | 單一元件直接用頂層名稱；後綴是為了區分，不是為了標示平台 |
| `scripts/` 或子任務內呼叫頂層聚合腳本 | 一律引用具體元件（`build:web`）；聚合腳本的語意會隨元件增減而改變 |
| 拆成 fan-out 後沒回頭改舊呼叫點 | 同一次改動內 `grep` 全 repo 修正，否則靜默重複建置或無窮遞迴 |
| 把 App Store／npm publish 掛在 `deploy` | 對外通路一律 `release`；`deploy` 只到自己在用的裝置或環境 |
| `dev` / `deploy` 沿用既有產物 | 兩者都必須重新建置，否則改動不會反映且沒有錯誤訊息 |
| 實機安裝取名 `dev:ios`、`dev:device` | 裝到會留著用的裝置是 `deploy`，不是 `dev` 的一個變體 |
| 上架腳本取名 `scripts/deploy.sh` | 檔名跟著階段走：對外上架是 `scripts/release.sh` |
