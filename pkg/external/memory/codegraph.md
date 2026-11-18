# CodeGraph

來源: [colbymchenry/codegraph](https://github.com/colbymchenry/codegraph) (MIT, ~78k stars, TypeScript + Rust kernel, 最新 release `v1.6.3` 2026-11-17) · npm `@colbymchenry/codegraph` · [文件站](https://colbymchenry.github.io/codegraph/)

預先建好的`程式碼 knowledge graph`, 透過 MCP 交給 coding agent. 它記的不是對話或決策, 而是 codebase 結構: 每個 symbol, call edge, import, 繼承關係. Agent 問 "X 怎麼運作" 時一次 `codegraph_explore` 就拿到相關 symbol 原文 + 呼叫路徑, 不必 grep / Read 一個個檔案爬.

這個 repo `不是 Claude skill`: 是自帶 runtime 的 CLI + MCP server (`codegraph serve --mcp`), `codegraph install` 會寫入各 agent 的 MCP 設定與使用指引. 放在 memory 分類是因為它是 agent 的`結構記憶`, 與其他四個的`經驗記憶`互補.

## 運作方式

1. Extraction - 內建 Rust kernel 以 tree-sitter 解析 20+ 語言, 抽出 node (function, class, method) 與 edge (calls, imports, extends, implements).
2. Storage - 寫入專案內 `.codegraph/codegraph.db` (SQLite + FTS5).
3. Resolution - 解析 call -> 定義, import -> 檔案, 以及 framework 特有連結 (Django / FastAPI / Express / NestJS / Rails 等 20+ 種 route -> handler; Swift <-> ObjC, React Native bridge, Tauri `invoke`).
4. Auto-sync - MCP server 以 FSEvents / inotify 監看檔案, 2 秒 debounce 後增量更新; 未同步的檔案會在回應前加 `⚠️` banner 提示 agent 直接 Read.

## Agent 介面

- MCP 預設只暴露`一個` tool `codegraph_explore` (實測單一強 tool 比工具選單更能引導 agent); `codegraph_node`, `codegraph_search`, `codegraph_callers`, `codegraph_callees`, `codegraph_impact` 等可用 `CODEGRAPH_MCP_TOOLS` 打開.
- CLI 同名對應: `codegraph explore <query>`, `callers`, `callees`, `impact <symbol>`, `affected [files...]` (找受影響的測試), `cycles` (循環 import).

## 安裝

```bash
curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh   # 或 npm i -g @colbymchenry/codegraph
codegraph install          # 設定 Claude Code, Codex, Cursor 等
cd your-project && codegraph init   # 建 .codegraph/ 並完成首次索引; --local-only 不進 git
codegraph status           # 檢查索引與 pending sync
```

本機現況: `未安裝`.

## 適用情境

- 大型或陌生 codebase 的架構問答; 上游 benchmark (Claude Code headless, 7 個 repo) 中位數 tool call 減 88%, token 減 62%, file read 歸零.
- 改動前做 impact analysis: 查某 symbol 的 callers 與影響半徑.
- CI 只跑受變更影響的測試 (`codegraph affected`).

## 注意

- 匿名 telemetry 預設開啟 (宣稱不含程式碼, 路徑, symbol 名稱); `codegraph telemetry off` 或 `CODEGRAPH_TELEMETRY=0` / `DO_NOT_TRACK=1` 關閉.
- `.codegraph/` 在專案根目錄, 記得確認 `.gitignore` 或用 `--local-only`.
- 與 claude-mem 的 `smart_search` / `smart_outline` 有部分重疊 (同為 tree-sitter), 但 claude-mem 是即時解析單檔, CodeGraph 是跨檔解析好的整張圖.
- 上游正在推 hosted 平台 (getcodegraph.com, PR 影響分析), 本機 CLI 仍完整可用且 100% 本地.
