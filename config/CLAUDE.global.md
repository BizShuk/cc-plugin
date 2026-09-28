# Global Rule

- 用繁體中文 + English Terminology 輸出
- 輸出分層: 變更摘要保持高階抽象 (high-level), 只說功能與差異 (例如 `增加圖片上傳功能`), 不講實作細節 (例如 `xxx.js 改了什麼`); 驗證步驟與後續行動必須具體, 寫出對象, 位置 (檔案, host, endpoint) 與目標值
- don't use double negative phrasing

## Principles

- Do not preserve backward compatibility. Remove obsolete paths instead of adding compatibility layers, fallbacks, or migrations.
- Choose the simplest implementation that fully meets the current requirements. Avoid speculative abstractions, configuration, and indirection.
- Grow the system in layers. Start from the smallest version that works end to end, and add each new capability on top of a product that already works. Never trade a working product for unfinished components.
- Keep components modular and concerns clearly separated.
- Prefer established, well-maintained libraries when they reduce overall complexity or improve reliability. Do not reimplement common functionality without a clear reason.
- Lean on the dependencies already in the project before writing your own implementation or adding packages. Do not assume a library lacks a capability without checking its documentation and types.
- Make architectural decisions for the long term. Do not accept a stopgap that only works for now and is meant to be replaced later.
- Opportunistic cleanup: when a minor change can align code with an existing pattern or make the overall structure cleaner and clearer, make the change, and mark it as an extra item in the change summary with a one-line reason (e.g., service logic split between the repo root and `svc/` — consolidate it into `svc/`).
- one file one responsibility. one package/folder one domain
- 遇到執行錯誤時，先嘗試修復，最多重試 5 次；若仍無法解決則明確報錯並停止, 多次遭遇相同錯誤/問題時，將解法記錄至 Memory
- Don't use git worktree or branch, just work on master branch, unless explictly mentioned.

### 上下文 (Context)

- 載入 `@./CLAUDE.md` 作為專案結構
- 載入 `@./README.md` 作為業務核心

## 專案結構 (Project Structure)

`~/projects/` 目錄佈局、每個 repo 必備的`統一介面 (Unified Interface)`、內容歸屬判準與
`plans/`／`docs/specs/` 的 `YYYY-MM-DD-<topic>.md` 檔名規範，全部由 `[[project-docs]]`
技能單一擁有（`references/unified-interface.md`、`references/content-ownership.md`）——
建立、撰寫或稽核專案文件前先讀它。
