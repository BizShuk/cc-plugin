# MemPalace

來源: [MemPalace/mempalace](https://github.com/MemPalace/mempalace) (MIT, ~61k stars, Python, 最新 release `v3.12.0` 2026-11-12, 預設 branch `develop`) · [PyPI](https://pypi.org/project/mempalace/) · [mempalaceofficial.com](https://mempalaceofficial.com/)

Local-first 的 AI 對話記憶, 主張`逐字儲存`: 不摘要, 不抽取, 不改寫, 只用 semantic search 取回原文. 以記憶宮殿 (method of loci) 結構化索引: 人與專案是 *wing*, 主題是 *room*, 原文放在 *drawer*, 所以搜尋可以限定範圍而非掃平面語料. LongMemEval R@5 raw `96.6%`, 全程零 API 呼叫.

這個 repo `不是 Claude skill`: 是 Python CLI + MCP server. 另提供 3 個 skills (`npx skills add MemPalace/mempalace`): `mempalace` (引導安裝與操作), `mempalace-recall` (回答前先搜尋), `mempalace-task` (logstream 委派); 裝 skill 不會順便裝 CLI.

## 架構

| 元件            | 內容                                                                                          |
| --------------- | --------------------------------------------------------------------------------------------- |
| Storage backend | 可插拔, 預設 ChromaDB; 另有 `sqlite_exact`, `rust_exact` (原生 Rust 掃同一個 SQLite 檔), `milvus`, `qdrant`, `pgvector` |
| Embedding       | 本機 ONNX: `embeddinggemma` (300m, 多語, 推薦) 或 `minilm` (英文, ~30MB); 也可 `openai-compat` 指向遠端 |
| Knowledge graph | 本機 SQLite 的 temporal entity-relationship graph, 每條關係帶有效期間 (add, query, invalidate, timeline) |
| MCP             | 47 個 tools: palace 讀寫, graph 操作, 跨 wing 導覽, agent diary, logstream 協調              |
| Agents          | 每個 specialist agent 有自己的 wing 與 diary, 以 `mempalace_list_agents` 動態發現               |

## 安裝與使用

```bash
uv tool install mempalace                              # 或 pipx install mempalace
mempalace init ~/projects/myapp --embedding-model embeddinggemma
mempalace mine ~/.claude/projects/ --mode convos       # 回填 Claude Code transcript
mempalace search "why did we switch to GraphQL"
mempalace wake-up                                      # 新 session 載入脈絡
mempalace sweep <transcript-dir>                       # 逐則訊息存成 drawer, 可重跑
```

也有 multi-arch Docker image `ghcr.io/mempalace/mempalace` (資料在 `/data`). 需求 Python 3.10+, ~300MB 磁碟給 embedding model. 本機現況: `未安裝`.

## 適用情境

- 要求`原文可追溯`, 不信任 LLM 摘要會失真的記憶場景.
- 完全離線或不想付 observer LLM 費用: 核心路徑不需 API key.
- 回填已存在的 Claude Code / Codex JSONL transcript, 建立歷史可搜尋庫.

## 注意

- 上游警告 Claude Code transcript `30 天後過期`; 沒接 auto-save hooks 就要先備份 `~/.claude/projects/*.jsonl` 再 `mine`.
- 官方來源只有 GitHub repo, PyPI 與 mempalaceofficial.com, 上游明確警告有冒名網站.
- Apple Silicon 上 embedding 預設跑 CPU (CoreML 會產生無效向量); 換 embedding model 需全部 re-embed.
- 用 `uv tool` / `pipx` 隔離安裝, 避免 `chromadb`, `grpcio` 污染全域 site-packages.
- 與 [claude-mem](claude-mem.md) 對比: claude-mem 是 LLM 壓縮的 observation, MemPalace 是逐字 drawer; 兩者 hooks 不要同時掛.
