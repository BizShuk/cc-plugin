# AgentMemory

來源: [rohitg00/agentmemory](https://github.com/rohitg00/agentmemory) (Apache-2.0, ~30k stars, TypeScript, 最新 release `v0.9.33` 2026-11-09) · npm `@agentmemory/agentmemory`

給 coding agent 用的`跨 agent 共用`長期記憶 server. 以 hooks 自動捕捉 session 事件 (prompt, tool call, 結果), 去除 secret 後存成 observation, 下次 session 透過 hybrid search (BM25 + vector + knowledge graph) 在 token budget 內注入回去. 自述是 Karpathy "LLM Wiki" 模式的實作, 額外加上 confidence scoring, 記憶 lifecycle 與 knowledge graph.

這個 repo `不是 Claude skill`: 是一個常駐的本機 server (建在 `iii-engine` runtime 之上), 再以各 agent 的 plugin / hooks / MCP / REST 接入. 附帶 17 個 skills (`npx skills add rohitg00/agentmemory`) 只是教 agent 何時去查記憶.

## 架構

| 元件          | 內容                                                                                          |
| ------------- | --------------------------------------------------------------------------------------------- |
| Server        | REST API `http://localhost:3111`, 預設只綁 loopback, protected 端點需 `~/.agentmemory/secret` bearer |
| Viewer        | `http://localhost:3113`, 瀏覽 / 編輯 / 匯出 / forget 記憶, 並可 replay 整段 session            |
| MCP           | 54 個 tools; `AGENTMEMORY_TOOLS=core` 縮成 8 個 (`memory_save`, `memory_recall`, `memory_smart_search` ...) |
| Storage       | macOS `~/Library/Application Support/agentmemory`, 可用 `--data-dir` 改; 支援 Redis backend    |
| 設定          | `~/.agentmemory/.env`                                                                         |

記憶分四層 consolidation: Working (原始 observation) -> Episodic (session 摘要) -> Semantic (事實與 pattern) -> Procedural (流程與決策). 依 Ebbinghaus 曲線衰減, 常被讀取的強化, 過期的自動淘汰, 並偵測矛盾.

## 運作流程

- `PostToolUse`: 穩定 event ID 寫入 capture inbox (server 不在時進 offline spool) -> privacy filter 去 secret -> 存原文 -> 壓縮 -> 有 embedding provider 才做 vector -> BM25 索引.
- `Stop` / `SessionEnd`: 增量摘要; `GRAPH_EXTRACTION_ENABLED=true` 才抽 knowledge graph.
- `SessionStart`: 載入 project profile + hybrid search, 預設注入上限 2000 tokens.

## 搜尋模式

- Keyless (預設): 不設 API key 就只有 BM25; LongMemEval-S 上 recall-all@5 `84.7%`, any-hit@5 `96.2%`.
- 本機 embedding: 內建 `@huggingface/transformers`, 首次下載 `Xenova/all-MiniLM-L6-v2`.
- 遠端 embedding: Gemini `gemini-embedding-001`, OpenAI `text-embedding-3-small` 等.

## 安裝

```bash
npx -y @agentmemory/agentmemory@latest        # 互動式 setup, 寫設定並啟動 server
npx -y @agentmemory/agentmemory@latest demo   # 塞 3 個範例 session 驗證 recall
agentmemory connect <agent>                   # 接其他 agent
agentmemory doctor                            # 診斷
```

Claude Code 接法: server 另開 terminal 跑, 再 `/plugin marketplace add rohitg00/agentmemory` + `/plugin install agentmemory`, 提供 12 個 hooks + MCP. 需求 Node.js 20+. 本機現況: `未安裝`.

## 適用情境

- 同時用多個 agent (Claude Code, Codex, Cursor, Hermes, OpenClaw ...) 想共用`同一份`專案記憶.
- 需要多 agent 協調原語: `memory_lease` (排他租約), `memory_signal_send` (agent 間訊息), `memory_action_*` (帶相依的工作項).
- 想要可視化檢查 agent 到底記了什麼, 並能逐筆刪除.

## 注意

- 常駐 server + 兩個 port (`3111`, `3113`); `--instance 1` 會改用 `3211` / `3213`, 多開時要避開.
- Context 自動注入是 opt-in, 未開啟時只捕捉不注入.
- 與 [claude-mem](claude-mem.md) 功能高度重疊 (hooks 捕捉 + 壓縮 + 注入), 兩者`不要`同時掛在 Claude Code, 否則每個 tool call 被捕捉兩次.
- `@agentmemory/mcp` 只是 shim, 連不到 server 時拿不到完整 54 個 tools.
