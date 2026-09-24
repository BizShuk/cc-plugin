# hermes — 術語表 (Terminology)

本專案無程式碼；出處欄指向上游 Hermes Agent 文件或本專案文件章節。起始規模偏小，待補 (To be extended)。

## 代理執行環境 (Agent Runtime)

| 術語 (Term) | 英文 (English) | 定義 (Definition) | 出處 (Source) |
| ----------- | -------------- | ----------------- | ------------- |
| Hermes | Hermes Agent | Nous Research 開源的常駐型 AI agent（Python），本專案部署與設定的對象；不是本專案自行開發的程式 | <https://github.com/NousResearch/hermes-agent> |
| Hermes 家目錄 | Hermes Home | 上游固定的執行期根目錄 `~/.hermes`，內含 `config.yaml`、`.env`、session 資料庫 `state.db` | 上游文件 Quick Start |
| 工作階段 | Session | 一段與 Hermes 的連續對話；以 `source` 欄位標記來自哪個通道 | 上游 `state.db` `sessions` 表 |
| 推理模型 | Model Provider | Hermes 自身用來規劃與對話的 LLM；與執行者使用的模型無關 | `hermes model` |

## 任務委派 (Task Delegation)

| 術語 (Term) | 英文 (English) | 定義 (Definition) | 出處 (Source) |
| ----------- | -------------- | ----------------- | ------------- |
| 委派 | Delegation | Hermes 把一件工作交給外部編碼代理 CLI 以非互動 (print mode) 方式執行，再回收結果 | 上游 bundled skill `autonomous-ai-agents/claude-code` |
| 執行者 | Executor | 被委派的外部代理 CLI；本專案限定 `agy`、`claude`、`grok` 三者 | README.md 任務委派 |
| agy | Antigravity CLI | Google Antigravity 的命令列代理 | `agy --help` |
| claude | Claude Code | Anthropic 的命令列編碼代理 | `claude --version` |
| grok | Grok Build | xAI 的命令列編碼代理 | `grok --help` |
| autop | autop | 工作區既有的 CLI facade，以單一入口啟動 `agy`／`claude`／`grok` 等 client；是委派的候選統一入口 | 工作區 `tools/autop` 專案 |

## 訊息通道 (Message Channel)

| 術語 (Term) | 英文 (English) | 定義 (Definition) | 出處 (Source) |
| ----------- | -------------- | ----------------- | ------------- |
| 閘道 | Gateway | Hermes 的常駐背景程序，同時連接所有已設定的通道、管理 session 並執行 cron | 上游文件 Messaging Gateway |
| 通道 | Channel | 使用者與 Hermes 對話的訊息平台（如 Feishu、DingTalk、WeCom）；上游稱 platform adapter | 上游文件 Messaging Gateway |
| 主通道 | Home Channel | 接收 cron 結果與主動通知的預設聊天室 | 上游文件 Feishu / Lark |
| 中國可達 | China Reachability | 通道在中國大陸網路環境下不需翻牆即可使用；本專案選通道的硬性條件 | README.md 訊息通道 |

## 縮寫 (Abbreviations)

| 縮寫 | 全稱 | 說明 |
| ---- | ---- | ---- |
| GFW | Great Firewall | 中國大陸的網路封鎖；Telegram、Slack、Discord、WhatsApp 等通道因此不可用 |
