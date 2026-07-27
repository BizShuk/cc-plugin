# README.md 模板與規則

## 模板 (Template)

```markdown
# <Project Name>

<1-2 句 elevator pitch：解決什麼業務問題>

## 業務領域 (Business Domains)

### <Domain 1 Name>

<2-3 句：此領域做什麼、為何存在、何時觸發>

`領域流程 (Domain Flow):`

1. <Step 1: entry point / trigger>
2. <Step 2: core processing>
3. <Step 3: outcome / side effects>

`核心實體 (Key Entities):` <Entity A>, <Entity B>, <Entity C>

`相關處理器 (Related Handlers):` <HandlerX>, <HandlerY>

---

### <Domain 2 Name>

<同上結構>

---

## 領域關聯 (Domain Relationships)

<描述領域之間如何互動。哪個領域的輸出是另一個領域的輸入？有沒有共用實體？>

## 使用方式 (Usage)

<主要 CLI commands、API endpoints、UI flows — 按領域分組>

## 改善建議 (Improvement Suggestions)

根據 codebase 分析：

- [ ] 建議 1：理由
- [ ] 建議 2：理由
- [ ] 建議 3：理由
```

## 規則 (Rules)

- 章節標題用繁體中文加英文括號
- 以 `業務領域 (Business Domain)` 為單位組織，不是以檔案或 handler 為單位
- 每個領域必須有：描述、流程、核心實體、相關處理器
- 領域流程要追溯真實程式碼路徑，不要抽象描述
- 使用專案中實際找到的 function/handler 名稱
- 改善建議必須具體可執行，根據真實發現；最少 3 個、最多 7 個
- 建議應涵蓋：領域邊界、缺漏的使用情境、資料流缺口
