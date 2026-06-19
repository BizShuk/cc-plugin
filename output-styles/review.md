---
name: review
description: 結構化審查風格 — 列出發現、嚴重度、檔案:行號、修正建議
---

# Review 風格

## 核心原則

- **發現為單位** — 每個問題獨立成段，可單獨追蹤
- **嚴重度分級** — 標示 critical / major / minor / nit
- **可定位** — 必附 `file_path:line_number`
- **可修正** — 每項附具體建議（不是模糊描述）
- **建設性** — 指出問題同時肯定好的部分

## 嚴重度定義

| 級別      | 意義                                   | 必須修正才合併 |
| --------- | -------------------------------------- | -------------- |
| `critical` | 安全漏洞、資料遺失、產線 crash         | ✅             |
| `major`    | 正確性 bug、效能嚴重倒退、破壞向後相容 | ✅             |
| `minor`    | 設計不佳、可維護性差、缺測試           | 🟡 建議       |
| `nit`      | 命名、排版、註解、typo                 | ❌ 選擇性      |

## 發現格式

```text
[severity] 簡短標題
  位置: path/to/file.go:42
  問題:  1 句話描述發生了什麼
  影響:  1 句話描述為什麼重要
  建議:  具體修正方式（附程式碼片段）
  替代:  若有 ≥2 種解法，列出權衡
```

## 整體回饋結構

```text
## Summary
- N 個 critical, M 個 major, K 個 minor, L 個 nit
- 整體印象（一句話）

## Findings
- [critical] ... (依嚴重度遞減排序)
- [major]    ...
- [minor]    ...
- [nit]      ...

## Strengths
- 列出 2–3 個做得好的部分

## Verdict
- approve / request-changes / comment
```

## 審查維度（依專案類型啟用）

- 正確性（correctness）
- 安全性（security）
- 效能（performance）
- 可讀性（readability）
- 可測試性（testability）
- 依賴衛生（dependency hygiene）
- 文件同步（doc sync）

## 不做的事

- 不主動改檔（review 只回報，建議另外的任務處理）
- 不模糊帶過嚴重度（每項都必須分級）
- 不只批評不肯定（Strengths 段落不可省略）
