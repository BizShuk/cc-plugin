# LiteLLM 本機代理 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 cc-plugin 專案內建立一個 docker-compose 化的 LiteLLM proxy，集中存放 minimax 真實金鑰，讓本機工具透過單一本機金鑰使用 minimax 模型。

**Architecture:** 單一 litellm 容器讀取 `config.yaml`，同時暴露 Anthropic (`/v1/messages`) 與 OpenAI (`/v1/chat/completions`) 格式端點。静態 master key 認證；Postgres/Redis 寫在 docker-compose 的 `db` profile 內、預設不啟動，預留動態虛擬金鑰升級。

**Tech Stack:** Docker, docker-compose, LiteLLM proxy (`ghcr.io/berriai/litellm:main-stable`), YAML。

設計依據：`docs/superpowers/specs/2026-05-20-litellm-proxy-design.md`

---

## File Structure

所有檔案位於 cc-plugin 專案的 `litellm-proxy/` 子目錄，外加修改專案根目錄的 `.gitignore`。

- `litellm-proxy/config.yaml` — LiteLLM 模型清單、router、認證設定。minimax 啟用，其他 provider 註解預留。
- `litellm-proxy/docker-compose.yml` — litellm service 常駐；db/cache service 掛 `profiles: ["db"]`。
- `litellm-proxy/.env.example` — 金鑰欄位範本，進版控。
- `litellm-proxy/.env` — 真實金鑰，不進版控（由 `.gitignore` 排除，Task 6 建立）。
- `litellm-proxy/logs/.gitkeep` — 保留 log 掛載目錄。
- `litellm-proxy/README.md` — 啟動、客戶端接法、DB 升級說明。
- `.gitignore` — 新增 `litellm-proxy/.env` 與 `litellm-proxy/logs/` 排除規則。

---

## Task 1: 建立目錄骨架與更新 .gitignore

**Files:**

- Create: `litellm-proxy/logs/.gitkeep`
- Modify: `.gitignore`

- [ ] **Step 1: 建立目錄與 logs 佔位檔**

Run:

```bash
mkdir -p litellm-proxy/logs && touch litellm-proxy/logs/.gitkeep
```

- [ ] **Step 2: 在 .gitignore 末尾新增排除規則**

讀取專案根目錄 `.gitignore`，在檔案末尾新增以下兩行：

```sh
litellm-proxy/.env
litellm-proxy/logs/
```

註：`logs/.gitkeep` 仍會被追蹤，因為下一步以明確路徑 `git add`；`logs/` 規則只排除之後產生的 log 檔。

- [ ] **Step 3: 確認 git 狀態**

Run: `git status --short`
Expected: 顯示 `.gitignore` 已修改、`litellm-proxy/logs/.gitkeep` 為未追蹤檔。

- [ ] **Step 4: Commit**

```bash
git add .gitignore litellm-proxy/logs/.gitkeep
git commit -m "chore: scaffold litellm-proxy directory and gitignore rules"
```

---

## Task 2: 核對 minimax endpoint 並寫 config.yaml

**Files:**

- Create: `litellm-proxy/config.yaml`

- [ ] **Step 1: 核對 minimax 的 API endpoint 與模型名稱**

用 WebFetch 取得 minimax 官方 API 文件，確認三件事：(1) OpenAI 相容的 `api_base`（國際站與中國站不同）、(2) 目標模型的正式名稱、(3) LiteLLM 建議的 provider 前綴。

Run: WebFetch `https://www.minimax.io/platform/document/platform%20introduction` 提問「OpenAI-compatible base URL and the exact chat model id」。
若該頁無法取得，改查 LiteLLM 對 minimax 的支援頁：WebFetch `https://docs.litellm.ai/docs/providers` 找 minimax 條目。

把核對到的 `api_base` 與模型名稱記下，用於 Step 2。若無法確認，採用下列預設值並在 README 標註「待線上驗證」：`api_base: https://api.minimax.io/v1`、模型 `MiniMax-M2`。

- [ ] **Step 2: 寫 config.yaml**

建立 `litellm-proxy/config.yaml`，內容如下（`api_base` 與 `model` 用 Step 1 核對的值取代）：

```yaml
model_list:
    # ---- minimax (啟用) ----
    - model_name: minimax-m2
      litellm_params:
          model: openai/MiniMax-M2
          api_base: https://api.minimax.io/v1
          api_key: os.environ/MINIMAX_API_KEY

    # ---- openai (預留模板，註解) ----
    # - model_name: gpt-4o
    #   litellm_params:
    #     model: openai/gpt-4o
    #     api_key: os.environ/OPENAI_API_KEY

    # ---- anthropic (預留模板，註解) ----
    # - model_name: claude-sonnet
    #   litellm_params:
    #     model: anthropic/claude-sonnet-4-6
    #     api_key: os.environ/ANTHROPIC_API_KEY

router_settings:
    num_retries: 2
    # fallbacks: []                       # 加第二家 provider 後填入

litellm_settings:
    json_logs: true
    # success_callback: []                # 預留外部 logging 整合

general_settings:
    master_key: os.environ/LITELLM_MASTER_KEY
    # database_url: os.environ/DATABASE_URL    # 升級 DB 時解開
    # store_model_in_db: true                  # 升級 DB 時解開
```

- [ ] **Step 3: 驗證 YAML 語法**

Run: `python3 -c "import yaml; yaml.safe_load(open('litellm-proxy/config.yaml')); print('config.yaml OK')"`
Expected: 印出 `config.yaml OK`，無例外。

- [ ] **Step 4: Commit**

```bash
git add litellm-proxy/config.yaml
git commit -m "feat: add litellm proxy config with minimax model"
```

---

## Task 3: 寫 .env.example

**Files:**

- Create: `litellm-proxy/.env.example`

- [ ] **Step 1: 建立 .env.example**

建立 `litellm-proxy/.env.example`，內容如下：

```
# litellm proxy 環境變數範本
# 複製為 .env 並填入真實值：cp .env.example .env

# minimax 真實 API 金鑰（必填）
MINIMAX_API_KEY=

# 本機 master key — 所有本機工具用這把連 proxy
# 產生方式：python3 -c "import secrets; print('sk-local-master-' + secrets.token_hex(16))"
LITELLM_MASTER_KEY=

# ---- 升級 DB 時解開 ----
# DATABASE_URL=postgresql://litellm:litellm@db:5432/litellm

# ---- 其他 provider 金鑰（啟用對應 config.yaml 區塊時填）----
# OPENAI_API_KEY=
# ANTHROPIC_API_KEY=
```

- [ ] **Step 2: Commit**

```bash
git add litellm-proxy/.env.example
git commit -m "feat: add litellm proxy env template"
```

---

## Task 4: 寫 docker-compose.yml

**Files:**

- Create: `litellm-proxy/docker-compose.yml`

- [ ] **Step 1: 建立 docker-compose.yml**

建立 `litellm-proxy/docker-compose.yml`，內容如下：

```yaml
services:
    litellm:
        image: ghcr.io/berriai/litellm:main-stable
        ports:
            - "4000:4000"
        volumes:
            - ./config.yaml:/app/config.yaml
            - ./logs:/app/logs
        env_file: .env
        command: ["--config", "/app/config.yaml", "--port", "4000"]
        healthcheck:
            test:
                ["CMD", "curl", "-f", "http://localhost:4000/health/liveliness"]
            interval: 30s
            timeout: 5s
            retries: 3
        restart: unless-stopped

    db:
        image: postgres:16
        profiles: ["db"]
        environment:
            POSTGRES_USER: litellm
            POSTGRES_PASSWORD: litellm
            POSTGRES_DB: litellm
        volumes:
            - litellm-db:/var/lib/postgresql/data

    cache:
        image: redis:7
        profiles: ["db"]

volumes:
    litellm-db:
```

- [ ] **Step 2: 驗證 YAML 語法**

Run: `python3 -c "import yaml; yaml.safe_load(open('litellm-proxy/docker-compose.yml')); print('compose OK')"`
Expected: 印出 `compose OK`，無例外。

註：完整的 `docker compose config` 驗證需要 `.env` 存在，留待 Task 6。

- [ ] **Step 3: Commit**

```bash
git add litellm-proxy/docker-compose.yml
git commit -m "feat: add litellm proxy docker-compose with db profile"
```

---

## Task 5: 寫 README.md

**Files:**

- Create: `litellm-proxy/README.md`

- [ ] **Step 1: 建立 README.md**

建立 `litellm-proxy/README.md`，內容如下：

````markdown
# litellm-proxy

本機 LiteLLM 代理，集中存放 minimax 金鑰，讓本機工具透過單一本機金鑰使用 minimax 模型。

設計文件：`../docs/superpowers/specs/2026-05-20-litellm-proxy-design.md`

## 啟動

```bash
cp .env.example .env
# 編輯 .env：填入 MINIMAX_API_KEY，並產生 LITELLM_MASTER_KEY：
#   python3 -c "import secrets; print('sk-local-master-' + secrets.token_hex(16))"

docker compose up -d
curl http://localhost:4000/health
```

## 客戶端接法

Claude Code CLI（Anthropic 格式）：

```bash
export ANTHROPIC_BASE_URL=http://localhost:4000
export ANTHROPIC_AUTH_TOKEN=<LITELLM_MASTER_KEY>
export ANTHROPIC_MODEL=minimax-m2
claude -p "hi"
```

OpenAI 格式工具：

```
base_url = http://localhost:4000/v1
api_key  = <LITELLM_MASTER_KEY>
model    = minimax-m2
```

## 觀測

- log：`docker compose logs -f litellm`（JSON 格式），或看掛載的 `./logs`
- Claude Code 自身花費：可獨立用 `ccusage`，與本 proxy 無關

## 升級為動態虛擬金鑰（DB profile）

1. `.env` 解開 `DATABASE_URL`
2. `config.yaml` 解開 `general_settings` 的 `database_url` 與 `store_model_in_db`
3. `docker compose --profile db up -d`
4. 開 `http://localhost:4000/ui` 用 Admin UI 管理金鑰與查花費
````

- [ ] **Step 2: Commit**

```bash
git add litellm-proxy/README.md
git commit -m "docs: add litellm proxy usage readme"
```

---

## Task 6: 啟動 proxy 並驗證端點

**Files:**

- Create: `litellm-proxy/.env`（不進版控）

註：本 Task 需要真實的 `MINIMAX_API_KEY`。若執行時無法取得，停在 Step 1 並請使用者提供。

- [ ] **Step 1: 建立 .env 並填入金鑰**

Run:

```bash
cd litellm-proxy && cp .env.example .env
```

編輯 `litellm-proxy/.env`：

- `MINIMAX_API_KEY` 填入真實 minimax 金鑰
- `LITELLM_MASTER_KEY` 填入下列指令產生的值：

Run: `python3 -c "import secrets; print('sk-local-master-' + secrets.token_hex(16))"`

- [ ] **Step 2: 驗證 compose 設定**

Run: `cd litellm-proxy && docker compose config -q && echo "compose config OK"`
Expected: 印出 `compose config OK`，無錯誤。

- [ ] **Step 3: 啟動 proxy**

Run: `cd litellm-proxy && docker compose up -d`
Expected: `litellm` 容器啟動；`db`、`cache` 不啟動（在 db profile 內）。

Run: `docker compose ps`
Expected: 只有 `litellm` 一個 service 在執行。

- [ ] **Step 4: 驗證健康檢查**

Run: `curl -s http://localhost:4000/health/liveliness`
Expected: 回傳存活回應（如 `"I'm alive!"`）。容器剛起需數秒，必要時重試。

- [ ] **Step 5: 驗證模型清單**

Run:

```bash
source litellm-proxy/.env
curl -s -H "Authorization: Bearer $LITELLM_MASTER_KEY" http://localhost:4000/v1/models
```

Expected: JSON 含 `minimax-m2`。

- [ ] **Step 6: 驗證 OpenAI 格式端點**

Run:

```bash
source litellm-proxy/.env
curl -s http://localhost:4000/v1/chat/completions \
  -H "Authorization: Bearer $LITELLM_MASTER_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"minimax-m2","messages":[{"role":"user","content":"ping"}],"max_tokens":32}'
```

Expected: JSON 含 minimax 回應的 `choices[0].message.content`。

- [ ] **Step 7: 驗證 Anthropic 格式端點**

Run:

```bash
source litellm-proxy/.env
curl -s http://localhost:4000/v1/messages \
  -H "Authorization: Bearer $LITELLM_MASTER_KEY" \
  -H "anthropic-version: 2023-06-01" \
  -H "Content-Type: application/json" \
  -d '{"model":"minimax-m2","max_tokens":32,"messages":[{"role":"user","content":"ping"}]}'
```

Expected: JSON 含 `content[0].text`。

註：本 Task 不 commit — 唯一新增的 `.env` 已被 `.gitignore` 排除。

---

## Task 7: 接 Claude Code 並做端到端驗證

**Files:** 無檔案變更（僅環境變數與驗證）。

- [ ] **Step 1: 設定 Claude Code 環境變數**

Run:

```bash
source litellm-proxy/.env
export ANTHROPIC_BASE_URL=http://localhost:4000
export ANTHROPIC_AUTH_TOKEN=$LITELLM_MASTER_KEY
export ANTHROPIC_MODEL=minimax-m2
```

- [ ] **Step 2: 端到端測試**

Run: `claude -p "reply with the single word: pong"`
Expected: 回應為 `pong`（或含 pong），代表 Claude Code → proxy → minimax 鏈路通。

- [ ] **Step 3: 驗證請求確實走 proxy**

Run: `docker compose -f litellm-proxy/docker-compose.yml logs --tail 20 litellm`
Expected: log 出現對應 Step 2 的 `/v1/messages` 請求紀錄。

- [ ] **Step 4: 驗證 DB profile 預留可用（不長期啟用）**

Run:

```bash
cd litellm-proxy && docker compose --profile db up -d && docker compose ps
```

Expected: `litellm`、`db`、`cache` 三個 service 都在執行。

Run: `cd litellm-proxy && docker compose --profile db down`
Expected: 收掉 db profile 的容器，確認預留機制可運作。

註：本 Task 不 commit。

---

## Self-Review

Spec coverage（對照 `2026-05-20-litellm-proxy-design.md`）：

- 概述/動機 → README (Task 5)
- config.yaml（minimax 啟用 + provider 模板）→ Task 2
- .env / .env.example → Task 3、Task 6
- docker-compose.yml（db profile 預留）→ Task 4
- 認證模型（單一 master key）→ Task 2 `general_settings`、Task 6 Step 1
- Claude Code 接法 → Task 5 README、Task 7
- OpenAI 客戶端接法 → Task 5 README、Task 6 Step 6
- 錯誤處理（num_retries / healthcheck / restart）→ Task 2、Task 4
- 可觀測性（json_logs / logs 掛載）→ Task 2、Task 4、Task 5
- 測試方式 → Task 6、Task 7
- 檔案結構 → File Structure 段、Task 1
- DB 升級路徑 → Task 4 profile、Task 5 README、Task 7 Step 4
- minimax endpoint 不確定性 → Task 2 Step 1 明確核對步驟

無 placeholder；型別/名稱一致（`minimax-m2`、`LITELLM_MASTER_KEY`、`MINIMAX_API_KEY` 全程一致）。
