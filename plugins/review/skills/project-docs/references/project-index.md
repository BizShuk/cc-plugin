# .project_index 規範

`.project_index/`（`projects.json` + `INDEX.md`）為全工作區的機器可讀註冊表，
掃描`兩層`（根目錄專案 + 分類目錄下的專案），依 README/CLAUDE.md 自動探索；
新增 repo 只要符合統一介面即自動被收錄。`categories` 與每個專案的 `tags`
為人工維護，`rebuild` 不會覆寫。
