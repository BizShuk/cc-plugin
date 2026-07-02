---
name: system-laws
description: >
  Use when needing a lens checklist to examine any system or artifact —
  designing boundaries, scheduling, load balancing, tracing, resilience,
  caching, compute, events, observability, or secrets. Also use when a
  universal operator (generate/review/consolidate/aggregate/evolve) needs
  its lens table. Triggers on: "system laws", "lens check", "which law
  applies", "13 laws", "role responsibility", "missing concern",
  "architecture checklist".
version: "2.0.0"
metadata:
  type: reference
  tier: philosophy
---

# system-laws — 13 法則 × 10 角色透鏡表

任何系統（或任何產物）都可以用 13 條法則逐一透視。每條法則配對一個角色：法則描述`力 (force)`，角色描述`責任 (duty)`。缺格即候選缺陷。

> 本表是五大通用算子的共用透鏡：`universal-review` 用它找缺陷、`universal-aggregate` 用它找缺格、`universal-generate` 用它設約束。

## 總覽表

| # | 法則 | 角色 | 域 | 系統對應 | 違反症狀 |
| :- | :--- | :--- | :--- | :--- | :--- |
| 1 | 空間 (Space) | 建築師 (Architect) | 宇宙基石 | namespace、服務邊界、目錄結構、資源歸屬 | 共用 table、命名衝突、上帝物件 |
| 2 | 時間 (Time) | 指揮家 (Orchestrator) | 宇宙基石 | 排程、事件順序、工作流編排、timeout | race condition、deadlock、間歇性 bug |
| 3 | 重力 (Gravity) | 指揮家 (Orchestrator) | 宇宙基石 | 負載平衡、流量整形、auto-scaling、熔斷 | 熱點崩潰、雪崩效應、容量失衡 |
| 4 | 因果 (Causality) | 實驗家 (Experimenter) | 宇宙基石 | 分散式追蹤、稽核日誌、RCA、可重現性 | 無法重現的 bug、只能猜的除錯 |
| 5 | 混沌 (Chaos) | 探險家 (Explorer) | 系統意志 | 混沌工程、PoC、邊界測試、技術雷達 | 生產環境才首遇故障、不知極限 |
| 6 | 精神 (Spirit) | 造物主 (Creator) | 系統意志 | AI/LLM 整合、決策引擎、推理管道 | 只能執行預設規則、無法處理模糊 |
| 7 | 生命 (Life) | 造物主 (Creator) | 系統意志 | health check、自動重啟、優雅關閉、韌性 | 脆性系統、崩潰需人工介入 |
| 8 | 破壞 (Destruction) | 實驗家 (Experimenter) | 系統意志 | 壓力測試、技術債清償、deprecation | 腐朽累積、廢棄 API 無法移除 |
| 9 | 冰霜 (Frost) | 修剪者 (Pruner) | 元素力量 | 快取、快照、凍結狀態、retention 修剪 | 快取失效災難、資料無限膨脹 |
| 10 | 烈焰 (Flame) | 推進者 (Propeller) | 元素力量 | 運算、熱路徑、業務邏輯轉化、效能燃燒 | 只存不算、hot path 無優化 |
| 11 | 雷霆 (Thunder) | 傳令官 (Herald) | 元素力量 | 事件、觸發器、告警、webhook、訊息傳導 | 只能輪詢、無法即時反應 |
| 12 | 光明 (Light) | 啟明者 (Illuminator) | 狀態邊界 | logging、metrics、tracing、可觀測性 | 黑盒系統、出事看不到 |
| 13 | 黑暗 (Dark) | 守密者 (Secretkeeper) | 狀態邊界 | secrets 管理、加密、隱私、隱藏狀態 | 機密裸奔、隱藏知識未被攤開 |

## 四大域 (Four Domains)

- `宇宙基石 (Cosmic Foundation)`：空間、時間、重力、因果 — 系統存在的絕對網格
- `系統意志 (System Will)`：混沌、精神、生命、破壞 — 系統的自主性與演化力
- `元素力量 (Elemental Forces)`：冰霜、烈焰、雷霆 — 能量的凍結、燃燒與傳導
- `狀態邊界 (State Boundaries)`：光明、黑暗 — 可見與隱匿的對偶

## 各法則設計原則

### 1. 空間 × 建築師

1. 為每個領域劃定清晰邊界；邊界內自治，邊界外走明確介面
2. 命名即文件：好的名稱讓邊界自我解釋

### 2. 時間 × 指揮家

1. 時序依賴顯式化，不依賴執行速度的假設
2. 所有異步操作定義完成條件與 timeout；用 idempotency 讓重放安全

### 3. 重力 × 指揮家

1. 任何單點都應可替換或擴展，不能成為必要的「重力中心」
2. 在系統入口做流量整形；過載時優雅降級而非崩潰

### 4. 因果 × 實驗家

1. 每個請求注入唯一 trace ID，貫穿所有服務
2. 重要狀態變更都要有原因紀錄；測試驗證因果而不只驗證結果

### 5. 混沌 × 探險家

1. 主動在非生產環境製造混沌，勝過被動在生產遭遇
2. 每次探索設定安全邊界，防止混沌蔓延

### 6. 精神 × 造物主

1. 推理邏輯與執行邏輯分離，使其可獨立升級
2. 推理結果必須可追溯；低信心時降級至人工審核

### 7. 生命 × 造物主

1. 預設系統會失敗，然後設計恢復機制
2. 區分 liveness 與 readiness；保活機制不能成為單點故障

### 8. 破壞 × 實驗家

1. 破壞是演化的必要條件，不要恐懼刪除
2. 廢棄要有明確時間表與遷移路徑，不能無限期保留

### 9. 冰霜 × 修剪者

1. 凍結的副本（cache / snapshot / backup）本質同構，統一失效策略
2. 一切會生長的資料都要有 retention 修剪規則

### 10. 烈焰 × 推進者

1. 找出熱路徑集中燃燒算力，冷路徑降級處理
2. 轉化必須守恆：輸入與輸出可對帳，運算不憑空產生或吞噬資料

### 11. 雷霆 × 傳令官

1. 狀態變化以事件傳導，不讓消費者輪詢
2. 告警要可行動 (actionable)：每個 alert 對應明確的處置動作

### 12. 光明 × 啟明者

1. 每個元件對外發光：log、metric、trace 三位一體
2. 可觀測性在設計期注入，不是出事後補裝

### 13. 黑暗 × 守密者

1. 機密只存在於 secrets manager，程式碼與日誌中零出現
2. 把隱藏知識（未文件化的不變式、口耳相傳的規則）主動攤開成文件

## 用法：透鏡檢查 (Lens Check)

對任一節點（service、module、文件、計畫）逐法則提問：「我的『X 法則』對應到哪個元件？」每個空格判定為：

- (a) 真的不需要 — 記錄理由
- (b) 存在但未被提及 — 補上對映
- (c) 真正缺口 — 列為候選缺陷，標信心等級
