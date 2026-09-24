# hermes — 技術脈絡 (Technical Context)

本專案不開發 agent，而是`部署與設定`上游 Hermes Agent；它是 `cc-plugin` 的子專案 (`pkg/hermes/`)。業務定義見 [README.md](README.md)，
用詞以 [docs/terminology.md](docs/terminology.md) 為準。

## 專案結構 (Project Structure)

```tree
pkg/hermes/
├── README.md            # 業務定義
├── CLAUDE.md            # 技術脈絡（本檔）
├── AGENTS.md -> CLAUDE.md
├── README.business.md   # 業務價值分析
├── README.todo          # 待辦（含未定的設計決策）
├── scripts/setup.sh     # 冪等安裝／升級上游 Hermes
├── scripts/setup-webui.sh # 冪等安裝／升級上游 Hermes WebUI（clone 至 ~/.hermes/hermes-webui）
├── scripts/run.sh       # 冪等設定同步（樣板、.env key、tmp/config），不安裝也不啟動服務
├── config/              # 與 ~/.hermes 同形；config.yaml 版控並 symlink，.env.example 只列 key
│   ├── workspace/       # ~/workspace 指向此處
│   └── memory/          # Hermes 記憶（→ ~/.hermes/memories）；USER.md 是指向 config/CLAUDE.global.md 的 symlink
├── persona/<name>/      # 人設 (persona)：IDENTITY.md、SOUL.md，另可帶 skills/
└── docs/
    ├── terminology.md
    ├── cli.md           # 常用指令與周邊工具（hermesd、hermes-webui）
    ├── slack.md         # Slack 通道設定步驟
    └── memory/
```

## 技術棧 (Tech Stack)

- Agent: Hermes Agent（Python，上游 `NousResearch/hermes-agent`）
- Executors: `agy`（Antigravity CLI）、`claude`（Claude Code）、`grok`（Grok Build）
- Delegation facade 候選: `autop`（工作區既有 Go CLI，已支援上述三個 client）
- Process manager: 未決定（pm2 或上游 `hermes gateway` 自帶的服務安裝）
- Channel: 未決定（見 README.md 訊息通道）

## 關鍵決策 (Key Decisions)

- `只設定不 fork`：Hermes 以上游安裝為準，本專案只放設定樣板、自訂 skill 與安裝腳本；
  不修改上游原始碼，升級走上游安裝程序。
- `執行期根目錄由上游固定`：Hermes 家目錄是 `~/.hermes`，不遵循工作區 `~/.config/<app>/` 慣例
  （該慣例來自 gosdk，Hermes 不是 gosdk 應用）。本專案的樣板與記憶檔由 `run:setup` 同步過去。
- `通道必須中國可達`：候選限於上游已有 adapter 的中國大陸原生平台；
  優先考慮 WebSocket 長連線型（Feishu、DingTalk），本機不需對外開放 port。
- `委派一律非互動`：執行者以 print／headless 模式呼叫，不依賴 PTY。

## 模組對應 (Module Mapping)

| 業務領域 (Domain) | 套件/模組 (Package/Module) | 進入點 (Entry Point) |
| ----------------- | -------------------------- | -------------------- |
| 代理執行環境 | 上游 Hermes（`~/.hermes`） | `hermes`、`hermes setup` |
| 任務委派 | 上游 bundled / official skills (`autonomous-ai-agents/*`) | `hermes skills install`、`claude -p`、`agy`、`grok`（或 `autop -c <client>`） |
| 訊息通道 | 上游閘道 | `hermes gateway setup`、`hermes gateway start` |

## 開發指南 (Development Guide)

- 任務一律在 `cc-plugin` 根層以 `hermes:*` 前綴執行（本子專案不自帶 `package.json`）
- Install: `pnpm run hermes:install` —— 缺 `hermes` 才以上游安裝器安裝（`--skip-setup`）
- Install WebUI: `pnpm run hermes:install:webui` —— 缺 WebUI 才 clone `nesquena/hermes-webui`；`pnpm run hermes:webui start|status|stop` 管理 daemon（:8787）
- Upgrade: `pnpm run hermes:upgrade`（WebUI：`hermes:upgrade:webui`） —— 對既有安裝跑 `hermes update`
- Setup: `pnpm run hermes:setup` —— `config/config.yaml`、`config/memory/`（→ `~/.hermes/memories`）與人設 `SOUL.md` 連結進 `~/.hermes`，
  `~/workspace` 連結到 `config/workspace/`，
  `config/` 樣板缺檔才複製，`.env` 只補缺少的 key，並連結 `tmp/config -> ~/.hermes`；
  `cc-plugin` 的 `run:setup` 會一併叫用
- Persona: `pnpm run hermes:persona <name>` —— 把 `~/.hermes/SOUL.md` 指向 `persona/<name>/SOUL.md`；
  預設 `yuna-v2`，可用 `HERMES_PERSONA` 覆寫，但下次 `run:setup` 會切回預設
- Lint: 根層 `pnpm run lint` 含 shellcheck
- Build: 不適用（無程式碼）
- Test: 未偵測到 (Not detected)
- Deploy: 未偵測到部署設定 (Not detected)

## 消費端契約 (Consumer Contracts)

- `daily-summary` skill 讀取 `~/.hermes/state.db` 的 `sessions`／`messages` 表，
  以 `sessions.source` 區分通道。更換 Hermes 家目錄或關閉 session 紀錄會讓日報少掉 Hermes 來源。

## 慣例 (Conventions)

- 機密（bot token、API key）只放 `~/.hermes/.env`，不進版控；專案內只放 `.env.example`。
- 版控與提交在 `cc-plugin` repo；`tmp/` 由 `cc-plugin` 根層 `.gitignore` 排除。
