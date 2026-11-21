# Claude-Mem

來源: [thedotmack/claude-mem](https://github.com/thedotmack/claude-mem) (Apache-2.0, ~118k stars, TypeScript, 最新 release `v16.0.0` 2026-11-18) · [docs.claude-mem.ai](https://docs.claude-mem.ai/)

Claude Code 生態中最普及的`自動 session 記憶`plugin. 在主 agent 旁邊跑一個 observer agent, 把每個 tool call 轉成簡短 observation, session 結束時寫摘要; 下次 session 開場注入一份精簡 index, agent 需要時才展開細節 (progressive disclosure). 同時支援 Codex, Cursor, OpenCode, Antigravity, OpenClaw, Pi 等.

這個 repo `是 Claude Code plugin`: 由 hooks + 本機 worker service + MCP search tools + skills 組成. npm 上的 `claude-mem` 套件只有 SDK, 不會註冊 hooks.

## 架構

| 元件        | 內容                                                                                       |
| ----------- | ------------------------------------------------------------------------------------------ |
| Hooks       | 每個 agent 事件跑 `claude-mem hook <agent> <event>`, 把 tool 活動交給 worker                |
| Worker      | Bun 管理的本機 HTTP service; 跑 observer, 提供 viewer (`/viewer`) 與 My Work (`/`)          |
| Observer    | 第二個 LLM agent (Claude, OpenRouter 等 provider), 產生 observation 與摘要                  |
| Storage     | SQLite `~/.claude-mem/claude-mem.db` (FTS5 關鍵字) + Chroma (semantic, 經 `uv` 安裝)        |
| 設定        | `~/.claude-mem/settings.json`; `CLAUDE_MEM_MODE` 控制模式與語言 (例 `code--zh`)            |

## 搜尋 (MCP)

三步驟節省 token, 由 `mem-search` skill 驅動:

1. `search` - 取 index (每筆 ~50-100 tokens)
2. `timeline` - 看某筆前後脈絡
3. `get_observations` - 只對挑中的 ID 取完整內容 (每筆 ~500-1,000 tokens)

另有 `smart_search` / `smart_outline` / `smart_unfold` (tree-sitter 程式碼導覽), `save_memory` (主動記筆記), `work_state_*` (待辦與工作狀態).

## 安裝

```bash
npx claude-mem install            # 偵測本機 agents 並逐一設定
npx claude-mem doctor             # 檢查 Bun, uv, hook launcher, worker
npx claude-mem telemetry disable  # 關閉預設開啟的匿名 analytics
```

或在 Claude Code 內 `/plugin marketplace add thedotmack/claude-mem` + `/plugin install claude-mem`. 需求: Node.js 20.19+ / 22.12+, Bun, uv (缺少會自動安裝).

本機現況: `已安裝`, plugin cache 版本 `16.0.0`; Bun 從 GitHub release 裝在 `~/.bun/bin` (brew 會因 prefix 過長改走 source build).

## 適用情境

- 單人主要用 Claude Code, 想要零操作的跨 session 延續: 第二個 session 起自動帶入前次脈絡.
- 回答 "上次怎麼解的" 這類問題: 以 observation ID 引用過去工作.
- `/learn-codebase` 一次預讀整個 repo (~5 分鐘) 作為冷啟動.

## 注意

- Observer 本身是 LLM 呼叫, 每個 session 都有額外 token 成本; `code--chill` 模式只記重要事項以降低量.
- 安裝流程預設引導登入 claude-mem 並開 30 天 trial 的 hosted observer; 要純本機就帶 `--provider` 或設 `CLAUDE_MEM_ONLINE_OPTIN=false`.
- 匿名 telemetry `預設開啟`; Cloud Sync 為 opt-in.
- `<private>` 標籤包住的內容不會被儲存; 其餘 tool 輸入輸出都會送到 observer model.
- 與 [agentmemory](agentmemory.md), [mempalace](mempalace.md) 的 auto-save hooks 重疊, 同一個 agent 只掛一套.
