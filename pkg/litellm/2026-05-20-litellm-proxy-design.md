# LiteLLM 本機代理架構設計 (LiteLLM Local Proxy)

日期：2026-05-20
狀態：設計核准，待寫實作計畫

## 概述

以 `LiteLLM proxy` 作為本機統一 LLM 閘道。真實的 `minimax` API 金鑰只集中存放在 proxy 端（透過環境變數注入），本機所有工具（Claude Code CLI、IDE plugin、腳本）一律連到 `localhost:4000`，用一組本機金鑰認證。真實的供應商金鑰不再散落於各個工具的設定檔。

核心動機：把 `minimax` 金鑰收進 proxy 後，任何本機情境都能透過本機金鑰使用 `minimax` 模型，金鑰輪替或更換供應商只需改 proxy 一處。

## 需求決策

| 項目 | 決策 |
|------|------|
| 範圍 | `minimax` 為主，config 結構預留多 provider 擴充 |
| 客戶端格式 | 同時暴露 Anthropic 格式 (`/v1/messages`) 與 OpenAI 格式 (`/v1/chat/completions`) |
| 金鑰管理 | 静態金鑰起步（單一 master key），預留 DB 升級至動態虛擬金鑰 |
| 部署 | Docker / docker-compose |
| 預留機制 | `方法 C` — compose profiles，DB 與 cache service 寫在檔案內、預設不啟動 |
| 存放位置 | cc-plugin 專案內的 `litellm-proxy/` 子目錄 |
| 可觀測性 | 兩階段，對齊 DB profile；不自寫工具 |

## 架構

```
            本機工具                          litellm-proxy 容器
  ┌──────────────────────────┐        ┌──────────────────────────────┐
  │ Claude Code CLI          │        │  LiteLLM Proxy  :4000        │
  │  ANTHROPIC_BASE_URL ─────┼──┐     │  ┌────────────────────────┐  │
  │  /v1/messages            │  │     │  │ 認證 (Bearer=master)   │  │
  ├──────────────────────────┤  ├────▶│  │ 限流 (記憶體模式)      │  │      ┌──────────┐
  │ IDE plugin / 腳本 / SDK  │  │     │  │ Router (retry/fallback)│  ├─────▶│ minimax  │
  │  OpenAI base_url ────────┼──┘     │  │ litellm SDK 格式轉換   │  │      │   API    │
  │  /v1/chat/completions    │        │  │ 非同步 log → stdout    │  │      └──────────┘
  └──────────────────────────┘        │  └────────────────────────┘  │
                                      └──────────────────────────────┘
                                       預留 (profiles: ["db"]，預設不起)
                                       ┌────────────┐  ┌────────────┐
                                       │ Postgres   │  │ Redis      │
                                       └────────────┘  └────────────┘
```

七層請求流程（對應 LiteLLM 官方架構）：

1. 入口：本機工具送請求到 `localhost:4000`
2. 認證：檢查 `Bearer` token 是否等於本機 master key
3. 限流：記憶體模式下做 rpm/tpm 限制（無 DB 階段為單機限流）
4. 路由：依模型名稱導向對應的 model group
5. 負載與容錯：Router 處理 retry，`fallbacks` 結構就緒（單一 provider 階段為空）
6. API 轉換：litellm SDK 把請求轉成 `minimax` 的 API 格式並執行呼叫
7. 非同步後處理：log 以 JSON 寫到 stdout，由 Docker 收集

## 元件

### litellm proxy 容器

- 映像：`ghcr.io/berriai/litellm:main-stable`
- 對外 port：`4000`
- 讀取掛載進來的 `config.yaml`
- 暴露端點：`/v1/chat/completions`、`/v1/messages`、`/v1/models`、`/health`
- compose 設 healthcheck 打 `/health/liveliness`

### config.yaml

`model_list` 依 provider 分組，`minimax` 區塊啟用，其他 provider 以註解模板預留。`general_settings` 內的 `database_url` 與 `store_model_in_db` 以註解形式預留，升級時解開即可。

```yaml
model_list:
  # ---- minimax (啟用) ----
  - model_name: minimax-m2
    litellm_params:
      model: openai/MiniMax-M2          # 以 OpenAI 相容 provider 接入
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

實作時以 `minimax` 官方文件核對 `api_base` 與模型名稱（國際站與中國站 endpoint 不同），並確認 LiteLLM 對該模型的最佳 provider 前綴。

### .env 與 .env.example

`.env` 由 `.gitignore` 排除，存放真實金鑰：

```
MINIMAX_API_KEY=<真實 minimax 金鑰>
LITELLM_MASTER_KEY=sk-local-master-<隨機字串>
# DATABASE_URL=postgresql://litellm:litellm@db:5432/litellm   # 升級 DB 時解開
```

`.env.example` 進版控，所有 provider 金鑰欄位以註解列出供參考。

### docker-compose.yml

`litellm` service 常駐；`db`(Postgres) 與 `cache`(Redis) 掛在 `profiles: ["db"]` 之下，預設 `docker compose up` 不啟動它們。

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
      test: ["CMD", "curl", "-f", "http://localhost:4000/health/liveliness"]
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

預設啟動：`docker compose up -d`（只有 litellm）。
升級啟動：`docker compose --profile db up -d`（litellm + Postgres + Redis）。

## 認證模型

静態階段只有一把金鑰 `LITELLM_MASTER_KEY`（例如 `sk-local-master-xxxx`），本機所有工具共用。LiteLLM 沒有「多把静態虛擬金鑰」的概念 — 多把可獨立管理、各有預算與限流的虛擬金鑰必須有 DB 支援。這正是預留 `--profile db` 的理由：升級後即可動態建立 / 撤銷 per-工具的虛擬金鑰。

## 客戶端接法

### Claude Code CLI（Anthropic 格式）

```
ANTHROPIC_BASE_URL=http://localhost:4000
ANTHROPIC_AUTH_TOKEN=<LITELLM_MASTER_KEY>
ANTHROPIC_MODEL=minimax-m2
```

Claude Code 送 Anthropic 格式請求到 `/v1/messages`，LiteLLM 內部轉成 `minimax` 的 OpenAI 相容格式呼叫。

### OpenAI 格式工具

```
base_url=http://localhost:4000/v1
api_key=<LITELLM_MASTER_KEY>
model=minimax-m2
```

## 錯誤處理

- Router 設 `num_retries: 2`，暫時性錯誤自動重試
- `fallbacks` 清單現在留空，結構就緒；加入第二家 provider 後填入即生效
- compose 對 litellm 設 healthcheck 打 `/health/liveliness`，失敗時容器標記為 unhealthy
- `restart: unless-stopped` 確保容器異常退出後自動重啟
- `minimax` 端故障時，LiteLLM 將上游錯誤透傳給客戶端；多 provider 階段則由 fallback 接手

## 可觀測性（兩階段，對齊 DB profile）

階段一（静態金鑰、無 DB）：

- `litellm_settings.json_logs: true`，log 以 JSON 寫 stdout，由 `docker compose logs` 查看；`./logs` 目錄掛載供保存檔案
- Claude Code 自身花費可獨立用 `ccusage` 查看 — 它讀的是 Claude Code 本機 JSONL session log，與 proxy 無關，照常運作
- 此階段沒有跨工具的 proxy 層儀表板，屬正常

階段二（`docker compose --profile db`）：

- 啟用 DB 後，LiteLLM Admin UI `/ui` 自動可用，即內建的 ccusage 等價物
- 多了 `/spend/logs`、per-key 花費、用量報表

不自寫 usage / dashboard 工具。Prometheus `/metrics` + Grafana 路徑保留但不納入本設計（個人用過重）。

## 測試方式

1. `curl http://localhost:4000/health` — 健康檢查
2. `curl -H "Authorization: Bearer <master key>" http://localhost:4000/v1/models` — 列出模型
3. OpenAI 格式：`curl /v1/chat/completions` 帶 `model: minimax-m2`，確認回應
4. Anthropic 格式：`curl /v1/messages` 帶 `model: minimax-m2`，確認回應
5. 端到端：設好 Claude Code 環境變數後執行 `claude -p "hi"`，確認走 proxy 並由 `minimax` 回應
6. 升級驗證：`docker compose --profile db up -d` 後確認 `/ui` 可開啟

## 檔案結構

```
cc-plugin/
└── litellm-proxy/
    ├── docker-compose.yml
    ├── config.yaml
    ├── .env.example          # 進版控
    ├── .env                  # gitignore 排除
    ├── logs/                 # gitignore 排除
    └── README.md             # 啟動、接客戶端、升級 DB 的說明
```

`.gitignore` 需新增 `litellm-proxy/.env` 與 `litellm-proxy/logs/`。

## 升級路徑（静態金鑰 → 動態虛擬金鑰）

1. `.env` 解開 `DATABASE_URL`
2. `config.yaml` 解開 `general_settings.database_url` 與 `store_model_in_db`
3. `docker compose --profile db up -d`
4. 透過 `/ui` 或 `/key/generate` API 為各工具建立獨立虛擬金鑰，取代共用 master key

無需重建容器或改寫架構，僅解註解與切換 profile。

## 非目標

- 不做網路 / 多裝置共用（僅本機 localhost）
- 不在本設計納入動態虛擬金鑰的實際啟用（僅預留）
- 不自寫 usage / 花費追蹤工具
- 不納入 Prometheus / Grafana 監控堆疊
- 不處理 minimax 以外 provider 的實際接入（僅 config 模板預留）
