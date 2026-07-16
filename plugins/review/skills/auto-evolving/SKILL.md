---
name: auto-evolving
description: >
    Use when a workspace should autonomously discover one high-leverage improvement, design it, apply a safe coherent update, verify the result, and consolidate durable knowledge into the project's canonical files. Scan broadly across user, business, domain, data, architecture, code, testing, reliability, operations, performance, cost, security, privacy, documentation, developer experience, and future adaptability, but converge on one main flow. Triggers on: "auto-evolving", "evolve this workspace", "improve this project end to end", "think design update", "自主改善", "自我演化", "廣域改善".
version: "2.0.0"
allowed-tools: Read, Write, Edit, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: high
context: workspace
metadata:
    type: methodology
    platforms: [macos, linux]
---

# auto-evolving

廣域觀察工作空間，收斂成一個主提案，完成 `THINK → DESIGN → UPDATE → VERIFY → CONSOLIDATE` 閉環。

## 核心契約 (Core Contract)

每個週期必須符合：

```text
1 cycle = 1 evidence-backed gap
        -> 1 selected proposal
        -> 1 coherent workspace update
        -> 1 validation result
        -> 1 canonical knowledge update
```

- 將多種觀點當成暫時的分析鏡頭 (lens)，不得建立永久候選分支。
- 每次只推進一個最高槓桿問題；新發現留給下一週期，不分叉目前週期。
- 以 workspace 的原始碼、測試、`README.md`、`CLAUDE.md`、`docs/` 與 `README.todo` 為主流程及真實來源 (source of truth)。
- 不建立 `branches/`、`foundation/`、`cycles/` 或平行知識樹，也不以模板覆蓋 workspace 根文件。
- 模型評分不是驗證。只有可重現的測試、觀察、使用結果或明確人工確認，才能升格為穩定知識。

## 執行模式 (Execution Mode)

- `full-cycle`：使用者呼叫本技能、要求 evolve/run/apply/update，且未限制為唯讀時，預設執行完整閉環並修改授權範圍內的 workspace 檔案。
- `plan-only`：使用者明確要求 review、analysis、ideas 或 plan only 時，只完成 `THINK → DESIGN`，不寫檔。
- `resume`：若已有相同主題的進行中計畫，先驗證其仍符合現況，再沿該計畫繼續；不得另開競爭分支。

外部訊息、部署、付費操作、權限變更、破壞性或不可逆操作，以及超出使用者範圍的重大架構變更，必須停在 `DESIGN` 並取得批准。

## 1. THINK — 廣域觀察與單點收斂

### 錨定 workspace

1. 讀取 workspace instructions、`README.md`、`CLAUDE.md`、`README.todo` 與相關 `docs/`。
2. 檢視相關程式碼、設定、測試、執行結果、git diff 與近期歷史；保留使用者既有未提交變更。
3. 用一句話分別寫出目前的產品承諾、技術主流程、使用者當前目標與已知限制。
4. 先搜尋既有實作、計畫與待辦，避免重新命名或重複建立同一概念。

### 廣域鏡頭

掃描所有鏡頭，但每個鏡頭最多保留一項有證據的發現；沒有證據就標記 `no finding`，不得湊數。

| 鏡頭 | 核心問題 |
| --- | --- |
| 使用者與工作流 | 哪個摩擦、缺口或人工步驟最妨礙使用者取得結果？ |
| 業務與產品 | 哪項改動最能提升價值、採用、留存、收益或降本？ |
| 領域、知識與資料 | 哪個概念、規則、資料契約或假設不完整或互相矛盾？ |
| 系統與依賴 | 哪個邊界、資料流、耦合或擴充點阻礙演進？ |
| 程式品質與測試 | 哪個複雜度、重複、缺測或錯誤處理缺口最危險？ |
| 可靠性、效能與運維 | 哪個故障模式、觀測缺口、成本或瓶頸最值得先解？ |
| 安全、隱私與合規 | 哪個信任邊界、權限、敏感資料或規範風險未被處理？ |
| 文件與開發者體驗 | 哪個脈絡、介面、命名或工具摩擦使維護容易出錯？ |
| 生態與未來適應性 | 哪個既有資產可產生新能力，或哪項鎖定會限制下一步？ |

### 候選與決策

1. 將重複發現合併為最多三個候選；每個候選必須附具體檔案、測試、日誌、行為或明確待驗證假設。
2. 先套用硬門檻。候選必須：
   - 直接符合 workspace 的業務與技術範圍。
   - 有至少兩項 workspace 證據，或一個可重現的測試／日誌證據。
   - 有可驗收結果、驗證方式與最小可逆落地路徑。
   - 不依賴未獲授權的外部或受保護操作。
3. 對通過門檻的候選以 0–5 評分，再換算為 100 分：

| 維度 | 權重 |
| --- | ---: |
| 使用者／業務價值 | 20 |
| 系統槓桿 | 15 |
| Workspace 契合度 | 15 |
| 證據與信心 | 15 |
| 急迫性／風險降低 | 10 |
| 可行性／可逆性 | 10 |
| 可驗證性／學習價值 | 10 |
| 跨面向正向效益 | 5 |

評分後執行以下收斂規則：

1. 只選最高分主提案。低於 70 分時，將它縮小成可驗證實驗；仍無法通過門檻則以 `no-change` 結束，不強迫修改。
2. 分數差小於 5 時，依序選擇較小可逆變更、較強證據、較高契合度、較低成本者。
3. 落選候選只在當次決策表記一行淘汰理由，不建立檔案或獨立待辦。確有獨立價值且已有證據者，最多新增一項下一週期 seed。

## 2. DESIGN — 統一主提案

為選定提案建立一份整合設計，不為各鏡頭建立不同方案。

設計必須包含：

1. `Outcome`：受益者、可觀察結果與成功門檻。
2. `Scope`：in scope、out of scope 與不變量。
3. `Evidence`：現況證據、根因與被否決的假設。
4. `Placement`：檔案位置、模組邊界、依賴方向、資料或控制流；介面不超過五個。
5. `Impact`：只列與主提案相關的廣域鏡頭影響及權衡。
6. `Landing`：三至七個可獨立驗證、可回滾的步驟。
7. `Acceptance`：驗收準則、驗證命令、失敗處理與回滾方法。
8. `Consolidation`：完成後每項知識的唯一歸屬檔案。

多檔案、跨模組或高風險變更先寫入 `plans/YYYY-MM-DD-evolve-<topic>.md`。小型且清楚的修改可在當次回覆列出內部設計後直接實作。`plan-only` 模式只在回覆中交付設計，除非使用者明確要求建立計畫檔。

## 3. UPDATE — 實作主提案

1. 只實作設計中的最小完整變更 (smallest coherent change)。
2. 遵循既有分層、命名、介面與設定慣例；不得為方便另創平行入口或自訂設定。
3. 同步新增或調整足以證明行為的測試、fixture、觀測或文件。
4. 若實作揭露另一個問題，記為下一週期 seed；除非它阻擋驗收，否則不擴大本週期。
5. 不覆蓋或回復使用者既有變更。若重疊無法安全合併，停止並說明衝突。

## 4. VERIFY — 驗證結果

1. 依序執行最小針對性檢查、相關套件測試，再執行成本合理的廣泛驗證。
2. 對照每項 acceptance criterion，記錄實際命令、結果與證據。
3. 檢查最終 diff 是否只包含主提案，並確認沒有新增重複入口、孤立文件或反向依賴。
4. 發生錯誤時先修復，最多重試五次。仍失敗則標記 `blocked`，不得把假設寫成已驗證知識。

## 5. CONSOLIDATE — 回到 workspace 主流程

驗證完成後，將結果增量更新到既有 canonical owner；同一事實只保留一份，其他位置以連結引用。

| 知識類型 | Canonical owner |
| --- | --- |
| 實際行為與合約 | 原始碼、設定、schema 與測試 |
| 業務定義、使用方式、domain flow | `README.md` |
| 技術結構、介面、依賴、關鍵決策與慣例 | `CLAUDE.md` |
| 已完成且仍有效的設計 | `docs/specs/YYYY-MM-DD-<topic>.md` |
| 決策緣由、驗證證據、反例與 retrospective | `docs/memory/YYYY-MM-DD-<topic>.md` |
| 尚未完成且可執行的工作 | `README.todo` |

- 完成的計畫依 workspace 慣例移入或濃縮至 `docs/specs/`，移除已完成待辦與矛盾的舊說明。
- 未驗證概念只留在進行中計畫或一項具驗證條件的 TODO，不得進入穩定文件。
- 驗證成功的概念直接融入 canonical owner，不另存 `principle` 或 `axiom` 複本。
- 被證偽或取代的概念從 canonical owner 移除；只有能防止重犯時才在 `docs/memory/` 保留簡短原因。
- 若變更未影響某 canonical file，記錄 `checked — no change`，不要為了顯示同步而製造內容。

## 舊版分支遷移 (Legacy Migration)

若目標 workspace 已存在舊版 `engine/`、`branches/`、`foundation/` 或 `cycles/`：

1. 將其視為唯讀歷史輸入，不再新增分支內容。
2. 去重現有候選，依本流程只選一個仍有效主提案。
3. 將已驗證知識增量合併至 canonical owner；未驗證項目回到單一計畫或 TODO。
4. 將必要決策歷史濃縮成一份 `docs/memory/` 記錄。
5. 只有使用者明確要求遷移或清理時，才封存或刪除舊目錄。

## 輸出契約 (Output Contract)

每次回覆只呈現一個主提案：

```markdown
# Evolution Cycle — <topic>

Status: plan-only | applied | no-change | blocked

## Think

- Selected gap and evidence
- Relevant lens findings
- Candidate score and decision reason

## Design

- Outcome, scope, placement, landing steps, acceptance and rollback

## Update

- Changed files and behavior, or proposed changes in plan-only mode

## Verify

- Command, result and evidence

## Consolidate

- Canonical files updated or checked
- One knowledge result: hypothesis | verified | invalidated
```

不得輸出四套平行方案、永久候選 branches、未落地的空泛概念，或把模型評分宣稱為實證。
