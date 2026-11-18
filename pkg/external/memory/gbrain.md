# GBrain

來源: [garrytan/gbrain](https://github.com/garrytan/gbrain) (MIT, ~33k stars, TypeScript / Bun, 最新 release `v0.60.739.0` 2026-11-18) · [gbrain.io](https://gbrain.io/) (選用的付費 hosting)

Y Combinator CEO Garry Tan 自用的 agent brain, 定位是`個人 + 公司知識庫`, 讓所有 agent 回答前先查它 (brain-first). 與 claude-mem / agentmemory 捕捉 session 不同, GBrain 以 Markdown 頁面為核心: 個人筆記 (可直接吃 Obsidian vault), 公司 `facts.md` / `decisions.md`, 會議逐字稿, Gmail / Calendar / Slack 同步, 再加上 typed knowledge graph 與帶引用的合成回答.

這個 repo `不是 Claude skill`: 是 CLI + MCP server + 背景 daemon (enrichment, dream cycle). 另有 Claude Code / Codex plugin 附帶 curated skills.

## 架構

| 元件              | 內容                                                                                       |
| ----------------- | ------------------------------------------------------------------------------------------ |
| Engine / Storage  | 預設 PGLite (嵌入式 Postgres, 免 server); 共享或大型 brain 改用 Postgres + pgvector         |
| Search            | `gbrain search`: vector + keyword + RRF + source-tier boost + reranker; 不設 key 時只有 keyword |
| Synthesis         | `gbrain think`: 檢索後合成帶引用的答案, 並明列 brain `不知道`的部分 (gap analysis)          |
| Knowledge graph   | 本機寫入時以 pattern 抽 `[[people/alice]]` 這類引用成 typed edge, 不呼叫 LLM               |
| Schema packs      | 頁面型別可換: 預設 `gbrain-base-v2` (15 型別: person, company, deal, concept ...), 可 `schema detect` 從既有資料夾推導 |
| 背景工作          | 開著的 `gbrain serve` 閒置時跑免費 upkeep; `gbrain maintain --schedule` 排程 dedup, 修引用, 找矛盾 |

## 安裝

```bash
bun install -g github:garrytan/gbrain#stable   # 需 Bun 1.4.4+; npm 上的 gbrain 是無關套件, 不要裝
gbrain setup claude-code                       # 寫 MCP, hooks 與 CLAUDE.md 的 brain-first 區塊; --dry-run 預覽, --remove 復原
gbrain sources add notes --path ~/notes --federated --force
gbrain sync --all
gbrain doctor
```

Claude Code plugin: `/plugin marketplace add garrytan/gbrain` + `/plugin install gbrain@gbrain`. 遠端存取: `gbrain mcp expose` 發佈到 Tailscale tailnet. 本機現況: `未安裝`.

## 適用情境

- 有大量個人筆記 / 公司文件, 想讓 agent 回答 "我們為什麼決定 X" 並附來源.
- 人脈與公司 CRM 式查詢 ("明天跟 Alice 開會前要知道什麼"), 搭配 Gmail / Calendar 的 open-loop 追蹤 (`gbrain waiting`).
- 團隊共用的 company brain: OAuth scoped, 依 source 與操作授權.

## 注意

- `gbrain setup` 會改寫 `~/.claude/CLAUDE.md` 加入 brain-first 指令; 本 workspace 的 home 設定由 cc-plugin 產生, 直接改會在下次 `./run.sh` 被覆蓋.
- 發版極頻繁 (一天多次), 一律用 `#stable` channel.
- 設定 embedding / reranking / synthesis provider 後, 頁面文字會送往該雲端服務; keyless 模式下 harness model 仍看得到 recall 內容.
- 不是 data warehouse: analytics, CRM, 金流資料留在原系統.
- Always-on 版 (24/7 ingest + dream cycle) 需 8GB+ RAM 的主機與自己的 API 費用.
