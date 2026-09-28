# business 模式 — 純業務分析

輸入不是完整 workspace，或明確要求純業務分析時，
跳過 `README.md` / `AGENTS.md` / symlinks，僅產出 `README.business.md`。

1. `Scope 界定` — folder/repo 用 Glob 鎖定 entry points；單一檔案直接讀取；
   純文字直接分析。
2. `Core Business & Operations` — 找出系統的存在理由、列出常見業務操作
   （業務動詞，不是函數名清單）。
3. `其餘章節` — 依 [readme-business.md](readme-business.md) 八章節完成。
4. `Write Report` — folder/repo 同步寫入兩個位置：
   - `<target>/README.business.md`
   - `~/projects/product/projects/<name>/README.business.md`
   既有檔案先讀取後合併，保留仍正確的內容。
