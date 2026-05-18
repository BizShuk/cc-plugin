# Hermes Agent

## Installation

```bash
# Installation - git clone(including curl) will use .hermes/hermes-agent/ as root
curl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh | bash
pip install hermes-agent


# Verify the installation
hermes doctor

# Configure model
hermes model
```

##

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
git clone https://github.com/nesquena/hermes-webui
cd hermes-webui
cp .env.docker.example .env
# Edit .env if your host UID isn't 1000 (e.g. macOS where UIDs start at 501)
brew install docker-compose docker-buildx
# docker buildx build.
docker-compose up -d # ❌ buildkit issue
# Open http://localhost:8787
```
