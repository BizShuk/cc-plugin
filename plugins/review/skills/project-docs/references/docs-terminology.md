# docs/terminology.md 模板與規則

專案的`術語單一定義來源 (single source of truth for terms)`。README、CLAUDE.md、
程式碼、commit message 一律引用此處用詞；同一概念不得有第二種說法。

## 模板 (Template)

```markdown
# <Project Name> — 術語表 (Terminology)

## <領域 Domain 1>

| 術語 (Term) | 英文 (English) | 定義 (Definition) | 出處 (Source) |
| ----------- | -------------- | ----------------- | ------------- |
| 記憶蒸餾    | Distillation   | 從多來源抽取候選記憶並去重後寫入儲存的流程 | `cmd/memory/distill.go` |

## 狀態值 (Status Values)

| 狀態 | 字面值 (Literal) | 語意 | 出處 |
| ---- | ---------------- | ---- | ---- |

## 縮寫 (Abbreviations)

| 縮寫 | 全稱 | 說明 |
| ---- | ---- | ---- |
```

## 規則 (Rules)

- 每筆術語必須有`出處 (Source)`：檔案路徑、識別符、或文件章節；查無出處者不得列入
- 狀態值以程式碼中的`字面值`為準（如 `"pending"`），不得寫美化過的說法
- 只收`領域名詞`與專案自訂縮寫；通用技術詞（HTTP、JSON、CLI）不收
- 一個概念一行；發現同義詞時挑一個為正名，其餘列入定義欄註明「舊稱」
- 依領域分節，與 `README.md` 的業務領域分節對齊
- 起始規模 10-30 筆；找不到足夠術語時寫最小表格並註明「待補 (To be extended)」
