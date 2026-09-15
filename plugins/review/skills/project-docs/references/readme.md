# README.md 模板與規則

歸屬判準見 [content-ownership.md](content-ownership.md)：`README.md` 回答
`為什麼用它、怎麼開始`，不擁有 entry point、結構樹與待辦。

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

---

### <Domain 2 Name>

<同上結構>

---

## 領域關聯 (Domain Relationships)

<描述領域之間如何互動。哪個領域的輸出是另一個領域的輸入？有沒有共用實體？>

## 使用方式 (Usage)

<每個領域 1-2 個 quick start 範例：最常用的一行指令或呼叫>

完整參考見 [docs/cli.md](docs/cli.md)。
```

## 規則 (Rules)

- 章節標題用繁體中文加英文括號
- 以 `業務領域 (Business Domain)` 為單位組織，不是以檔案或 handler 為單位
- 每個領域必須有：描述、流程、核心實體
- 領域流程要追溯真實程式碼路徑，不要抽象描述
- `不寫 entry point 與型別名`：handler／function 名稱由 `CLAUDE.md` 模組對應表單一擁有，
  領域名稱即兩表的對照鍵。需要先讀原始碼才懂的名詞出現在這裡就是越界
- `使用方式`只留 quick start；章節超過約 25 行即下放 `docs/cli.md`，本檔留指標
- `不寫待辦`：分析出的改善建議寫進 `README.todo`，不進 `README.md`
- 結構樹、架構決策一律不複製，用一行指向 `CLAUDE.md`
