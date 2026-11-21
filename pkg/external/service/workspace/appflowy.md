# appflowy

來源: [AppFlowy-IO/AppFlowy](https://github.com/AppFlowy-IO/AppFlowy) (AGPL-3.0, ~79k stars) · [website](https://www.appflowy.com) · [docs](https://docs.appflowy.io)

開源的 Notion 替代品, 定位是 "AI collaborative workspace": 文件, wiki, database (table / kanban / calendar), 團隊協作與 AI 助手放在同一個 workspace, 主打 `資料自己掌控` (local-first, 可 self-host). 桌面 (macOS / Windows / Linux) 與手機 (iOS / Android) 都有 client.

這個 repo `不是 Claude skill`: 沒有 `SKILL.md`, 是一個完整的桌面 / 行動 app. 列在這裡是作為外部參考, 用於筆記 / 知識庫工具選型, 或研究 local-first 協作架構.

## 架構

| 部分 | 用途 |
| ---- | ---- |
| `frontend/appflowy_flutter/` | Flutter UI, 跨桌面與手機的單一 codebase |
| `frontend/rust-lib/` | Rust 核心, 以 `dart-ffi` 接給 Flutter; `lib-dispatch` 是 UI 與核心之間的 event dispatcher |
| `flowy-document` / `flowy-database2` / `flowy-folder` | 文件, database, workspace 樹三大領域, 每個領域另有 `-pub` crate 放對外 interface |
| `collab-integrate` | 以 CRDT (多人同時編輯不衝突的資料結構) 做即時協作與同步 |
| `flowy-sqlite` / `flowy-sqlite-vec` | 本機儲存, 外加 vector index 供 AI 搜尋 |
| `flowy-ai` | AI chat / 摘要 / ask AI, 內含 `mcp/` manager (可接 MCP server) |
| `flowy-server` / `flowy-user` | 雲端同步 (官方 cloud 或 self-host 的 AppFlowy Cloud) 與帳號 |

核心設計:

- UI 薄, 邏輯在 Rust: Flutter 只送 event, 業務規則與資料都在 Rust 核心, 同一份核心服務所有平台.
- 每個領域拆成實作 crate 與 `-pub` interface crate, 避免領域之間直接耦合.
- local-first: 資料先寫本機 SQLite, 再經 CRDT 同步; 離線可用.

## 安裝

- 使用者: 從 [releases](https://github.com/AppFlowy-IO/AppFlowy/releases) 下載, 或 FlatHub / Snapcraft / App Store / Play Store.
- Self-host: 依 [Self-Hosting Guide](https://appflowy.com/docs/Step-by-step-Self-Hosting-Guide---From-Zero-to-Production) 部署 AppFlowy Cloud (後端在另一個 repo).
- 從原始碼: 需 Flutter + Rust toolchain, 步驟見 [from-source](https://docs.appflowy.io/docs/documentation/appflowy/from-source).

## 適用情境

- 想要 Notion 式 workspace 但資料要留在自己機器或自己的 server (例如部署到 liva).
- 參考 Flutter + Rust FFI + CRDT 的跨平台 local-first app 架構.
- 需要 AI 搜尋筆記: 內建 vector index 與 MCP 接點.

## 注意

- `AGPL-3.0`: 修改後以網路服務提供給他人, 必須公開修改後的原始碼.
- AI 功能預設走 AppFlowy 雲端模型; 隱私敏感內容改用 `Local AI` (設定頁 Settings -> AI, 接本機 Ollama) 讓 chat 推論留在本機, `flowy-ai/src/local_ai` 是對應實作.
