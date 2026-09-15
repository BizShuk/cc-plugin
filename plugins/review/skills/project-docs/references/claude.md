# CLAUDE.md 模板與規則

歸屬判準見 [content-ownership.md](content-ownership.md)：`CLAUDE.md` 回答
`邊界是什麼、誰擁有什麼`，是結構樹、模組對應與架構決策的單一 owner。

## 模板 (Template)

```markdown
# <Project Name> — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

<實際目錄樹，2-3 層深>

## 技術棧 (Tech Stack)

- Language: <detected>
- Framework: <detected>
- Build tool: <detected>
- Key dependencies: <top 5-8 deps>

## 關鍵決策 (Key Decisions)

- Decision 1：為何選擇此做法（從程式碼模式推斷）
- Decision 2：...

## 模組對應 (Module Mapping)

把每個業務領域（從 README）對應到技術實作：

| 業務領域 (Domain) | 套件/模組 (Package/Module) | 進入點 (Entry Point) |
| ----------------- | -------------------------- | -------------------- |
| <Domain 1>        | `pkg/xxx`, `handler/yyy`   | `HandleXxx()`        |
| <Domain 2>        | `pkg/aaa`, `handler/bbb`   | `HandleAaa()`        |

## 開發指南 (Development Guide)

- Build: `<精確的 build 指令>`
- Test: `<精確的 test 指令，或註明無測試>`
- Deploy: `<可偵測的部署方式，或「未偵測到部署設定 (Not detected)」>`

前置需求、安裝步驟、排程與設定同步細節見 [docs/development.md](docs/development.md)；
指令一覽見 [docs/cli.md](docs/cli.md)。

## 慣例 (Conventions)

- Naming: <detected patterns>
- Error handling: <detected patterns>
- Logging: <detected patterns>
- Testing: <detected patterns>
```

## 規則 (Rules)

- 章節標題用繁體中文加英文括號
- `模組對應`的業務領域必須與 `README.md` 的領域`逐一對齊`，領域名稱是兩表的對照鍵
- `開發指南`只留 build／test／deploy 各一行；展開細節下放 `docs/development.md`，
  本檔留指標（判定門檻：章節超過約 25 行）
- `關鍵決策`只寫`仍然成立的不變式`；「已移除」「原 X 是」這類歷史敘述屬
  `docs/CHANGELOG.md`，不留在本檔
- 易腐計數（行數、檔案數、module 數、容量上限）一律不寫 —— 每次 commit 都會讓它變假
- 寫得出 pass/fail 的斷言不放 Markdown，落成測試或 `scripts/`，本檔只留一句規則 + 測試名稱
- 絕對路徑改相對路徑或 `$(git rev-parse --show-toplevel)`
- 不使用粗體強調，改用 `backtick`
