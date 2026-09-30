# Hermes CLI

## Installation

安裝與升級走 `pnpm run hermes:install`／`pnpm run hermes:upgrade`（於 `cc-plugin` 根層）（見 [AGENTS.md](../AGENTS.md)）。
`pip install hermes-agent` 與 `brew install hermes-agent` 不採用：少了 whatsapp 腳本。

```bash
hermes doctor   # 驗證安裝
hermes model    # 設定推理模型
```

## Skills (任務委派)

安裝執行者委派技能以操控 `agy` 與 `grok` CLI：

```bash
hermes skills install official/autonomous-ai-agents/antigravity-cli
hermes skills install official/autonomous-ai-agents/grok
```


## 檔案與指令

```bash
# 📁 Your files:
#    Config:    /Users/shuk/.hermes/config.yaml
#    API Keys:  /Users/shuk/.hermes/.env
#    Data:      /Users/shuk/.hermes/cron/, sessions/, logs/
#    Code:      /Users/shuk/.hermes/hermes-agent

# 🚀 Commands:
hermes                  # Start chatting
hermes setup            # Configure API keys & settings
hermes config           # View/edit configuration
hermes config edit      # Open config in editor
hermes gateway install  # Install gateway service (messaging + cron)
hermes update           # Update to latest version
hermes desktop # 📁 Your files:
#    Config:    /Users/shuk/.hermes/config.yaml
#    API Keys:  /Users/shuk/.hermes/.env
#    Data:      /Users/shuk/.hermes/cron/, sessions/, logs/
#    Code:      /Users/shuk/.hermes/hermes-agent


```

## hermesd

`hermesd` 是一個為 Hermes Agent 打造的即時 TUI 監控儀表板 (Real-time TUI monitoring dashboard)。它以唯讀模式 (read-only) 運作，自動讀取 `~/.hermes/` 內的狀態，提供一個統一的介面來查看網關狀態 (Gateway health)、活躍對話 (Sessions)、Token 消耗與成本估算 (Tokens / Cost)、排程任務 (Cron jobs)、已安裝的技能 (Skills) 及系統日誌 (Logs)。使用者無需查閱多個日誌檔案或反覆輸入狀態指令，即可一目了然掌握 Agent 運行全貌。

<https://github.com/mudrii/hermesd>

```bash
uvx hermesd
```

## hermes-webui

<https://github.com/nesquena/hermes-webui>

```bash
cd ~/.hermes/
git clone https://github.com/nesquena/hermes-webui
cd hermes-webui
./start.sh
ln -s ~/.hermes/hermes-webui/start.sh ~/.local/bin/
```

Retrieve and follow the instructions at:
<https://raw.githubusercontent.com/garrytan/gbrain/master/INSTALL_FOR_AGENTS.md>

## PM2 常駐管理 (PM2 Process Management)

Hermes Gateway 於 `cc-plugin` 根目錄 `ecosystem.config.js` 登記為 `Agent` namespace 的常駐程序 (`Hermes Gateway`)，執行腳本為 `pkg/hermes/scripts/gateway.sh`（自動載入 `/Users/shuk/.venv` 環境並以 `--external-supervisor` 模式運作）：

```bash
# 註冊 / 套用設定（預設 optional: true，處於 paused 狀態）
pm2 apply ecosystem.config.js

# 啟動 / 恢復 Hermes Gateway
pm2 task resume "Agent:Hermes Gateway"

# 暫停 Hermes Gateway
pm2 task pause "Agent:Hermes Gateway"

# 查看狀態與即時日誌
pm2 list
pm2 logs "Hermes Gateway" --lines 50
```

> 若先前已透過 `hermes gateway install` 安裝 launchd 服務，切換至 PM2 前請先執行 `hermes gateway stop` 停止 launchd 服務，以避免雙重派發衝突。
