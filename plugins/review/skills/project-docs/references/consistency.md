# 偵測一致性 (Consistency Detection)

文件在程式碼改動的當下就開始漂移。本技能檢測兩個維度的一致性：

- `縱向:` 每份文件 vs 程式碼（文件宣稱的事實是否仍為真）
- `橫向:` 文件 vs 文件（同一概念在不同文件中是否一致）

## 檢測流程

1. 從所有文件抽出可驗證的宣稱 — 路徑、指令、模組對應、術語、狀態值、領域名稱。
2. `縱向驗證:` 逐條對 repo 驗證，以 `doc says X → actually Y` 回報。
3. `橫向驗證:` 交叉比對文件間的引用，以 `A says X / B says Y` 回報。

## 縱向：文件 vs 程式碼

| 文件 | 檢查項 |
| ---- | ------ |
| `README.md` | 業務領域仍對應到真實 handler/module；CLI/API 指令仍可執行 |
| `AGENTS.md` | 目錄樹 diff 真實目錄；模組對應 entry point 仍存在；build/deploy 指令有效；設定路徑一致 |
| `docs/terminology.md` | 術語出處路徑仍存在；狀態字面值與程式碼 enum/const 一致 |
| `README.business.md` | 狀態名稱能在程式中找到；上游服務在程式碼中有對應呼叫 |
| 統一介面 | 必備檔案存在：`AGENTS.md`、`README.todo`、`docs/memory/`、`docs/terminology.md`；`CLAUDE.md` 是指向 `AGENTS.md` 的 symlink |

## 橫向：文件 vs 文件

| 比對對 | 檢查項 |
| ------ | ------ |
| `README.md` ↔ `AGENTS.md` | README 的業務領域必須全數出現在 CLAUDE 的模組對應表；反之模組對應不得列出 README 未定義的領域 |
| `README.md` ↔ `docs/terminology.md` | README 中出現的領域名詞必須在術語表有定義；術語表的領域分節必須與 README 業務領域對齊 |
| `AGENTS.md` ↔ `docs/terminology.md` | CLAUDE 使用的技術術語（若為領域詞）必須與術語表一致；不得出現同義詞漂移 |
| `README.business.md` ↔ `README.md` | business 的業務目的/常見操作必須與 README 業務領域對應；不得出現 README 未提及的領域 |
| `README.business.md` ↔ `docs/terminology.md` | 狀態值名稱必須與術語表的狀態值章節一致 |
| 所有文件 ↔ `docs/terminology.md` | 同一概念在所有文件中只能使用術語表定義的正名，不得有第二種說法 |

## audit 輸出格式

```text
Doc consistency review — <scope>

[縱向] 文件 vs 程式碼:
- AGENTS.md tree: omits plugins/god/ and plugins/team/ (both exist)
- README.md: "uses MySQL" → code uses SQLite (model/store.go)
- README.business.md: state "archived" not found in code

[橫向] 文件 vs 文件:
- README ↔ CLAUDE: README 列「資料匯出」領域，CLAUDE 模組對應表無此領域
- README ↔ terminology: README 用「蒸餾管道」，CLAUDE 用「distiller pipeline」— 術語表未定義
- business ↔ terminology: business 狀態 "pending" vs 術語表 "waiting" — 應統一

[統一介面] 缺件:
- docs/terminology.md: missing (必備)
```
