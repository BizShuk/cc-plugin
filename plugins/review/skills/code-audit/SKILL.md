---
name: code-audit
description: >
    Use when auditing existing code or a diff for defects that are visible
    without running it. Covers cross-file consistency (contradictory business
    rules, drifting domain models, producer/consumer contract mismatches),
    security (secrets, input validation, authn/authz, injection, unsafe
    dependencies, sensitive data in logs), and identifier naming (casing drift,
    synonym sprawl, vague or misleading names).
    Triggers on: "review this code", "audit this diff", "check consistency",
    "is this secure", "security review", "any vulnerabilities", "review naming",
    "naming convention", "rename suggestions", "程式碼審查", "一致性檢查",
    "安全稽核", "命名規範", "code-audit".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: high
context: fork
metadata:
    type: review
    platforms: [macos, linux]
---

# code-audit

對`既有程式碼或一份 diff`做靜態審查，回報可被證據支持的缺陷。
三個`軸線 (axis)`各自獨立，可單獨執行也可一併執行：

| 軸線 (Axis)        | 問的問題                             | 參考文件                                             |
| ------------------ | ------------------------------------ | ---------------------------------------------------- |
| `一致性 (Consistency)` | 這份修改跟工作區其他地方矛盾嗎?       | [references/consistency.md](references/consistency.md) |
| `安全 (Security)`      | 這段程式碼會洩漏或被濫用嗎?           | [references/security.md](references/security.md)       |
| `命名 (Naming)`        | 讀名字能預測行為嗎? 同一概念同一個詞嗎? | [references/naming.md](references/naming.md)           |

## 核心契約 (Core Contract)

- `唯讀預設`：只提出建議，不改寫程式碼。使用者明確說 `fix` 才動手。
- `證據優先`：每筆發現必須引用 `file:line` 並指出它違反的規則，無法定位的疑慮不列入。
- `乾淨也是結果`：某個軸線查無問題就明確寫 `ok`，不得為了湊數而捏造發現。
- `慣例勝過外部指南`：先從 codebase 推斷既有慣例，再以該慣例為基準判斷，
  不套用與專案無關的風格書。
- `嚴重性看忽略成本`，不看發現的難度。

## 流程 (Procedure)

1. `界定目標`：未指定時依序解析 — 未提交變更 (`git diff`) → 分支與 base 的差異 →
   詢問使用者。確認後用一句話覆述審查範圍。
2. `選軸線`：依目標內容決定要跑哪些軸線，並記錄跳過的原因。
   - 任何程式碼變更 → 一致性 + 命名
   - 涉及輸入處理、認證授權、外部請求、機密、依賴清單 → 加上安全
   - 純文件變更 → 三軸皆跳過，改用 `[[project-docs]]`
3. `執行`：載入對應 reference，依其檢查鏡頭逐項比對。
4. `彙整`：跨軸線去重，同一行被多軸線命中時只留一筆並註記其他鏡頭。
5. `輸出`：依嚴重性排序回報，不寫檔。

## 嚴重性 (Severity)

| 等級      | 意義                                       |
| --------- | ------------------------------------------ |
| `blocker` | 破壞行為、資金、資料或造成外洩的矛盾與缺口 |
| `major`   | 實質缺陷，應在當前 PR 修復                 |
| `minor`   | 偏離慣例或衛生問題，可批次清理             |
| `nit`     | 化妝性問題，講一次即可                     |
| `ok`      | 該軸線已審查且無發現                       |

## 輸出格式 (Output Format)

```text
Code audit — <一句話範圍>
Ran: consistency, security, naming

[blocker] security     svc/auth.go:42 — bearer token 寫進 slog.Info
[major]   consistency  a.go:42 <-> b.go:88 — 同一規則兩處互斥
          -> also naming: 同概念在兩處叫 "tenant" 與 "account"
[minor]   naming       cmd/Read_Logic.go -> read_logic.go (檔名 casing 漂移)
[ok]      consistency  producer/consumer 合約一致

Top fix (value/effort): 移除 auth.go:42 的 token 日誌, 再統一 a.go/b.go 的規則.
```

## Common Mistakes

- `憑品味重構`：提出不減少歧義也不回復既有一致性的改名或改寫。
- `把架構問題當一致性問題`：目錄佈局、模組邊界、依賴衛生屬 `[[planner]]` 系統軸。
- `安全軸線漫無邊際`：列出與本次目標無關的通用威脅清單，而非指向具體 `file:line`。
- `跨軸線重複計數`：同一行在三個軸線各報一次，讓報告膨脹而失去排序意義。
- `擅自修復`：未經要求就改動程式碼。

## 相關技能 (Related Skills)

- `[[planner]]` 系統軸：目錄佈局、模組邊界、依賴衛生與架構規劃。
- `[[project-docs]]` 當落差在文件與程式碼之間，而非程式碼本身。
- `[[anti-sabotage]]` 當懷疑測試掩蓋了只在正式環境才成立的行為。
