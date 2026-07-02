---
name: universal-aggregate
description: >
  Use when composing parts into a higher-order whole — building an index,
  summary, matrix, dashboard, ontology, roll-up report, or assembling
  components into a system view — or when checking completeness of a
  collection. Triggers on: "aggregate", "roll up", "index", "summary of
  all", "big picture", "matrix", "ontology", "completeness", "彙總",
  "總覽", "拼圖".
version: "2.0.0"
metadata:
  type: pattern
  tier: philosophy
  operator: aggregate
---

# universal-aggregate — 聚合算子

`簽名 (Signature)`：{異類部件 ×N} → 高階整體 (whole)，N → 1（升維）

把部件組合成高於部件的整體：索引、矩陣、總覽、本體論 (ontology)、系統視圖。與 `universal-consolidate` 的分界：aggregate 的輸入是`不同`部件，輸出是`新維度`的整體；consolidate 的輸入是`同類`冗餘，輸出仍是同維度的典範。

## 四階段程序

### 1. 定座標 (Define Axes)

先選正交座標軸，再放內容。軸必須互不干涉（改一軸不迫使改另一軸）：

- 範例：`env / domain / entity / attribute` 四座標定位任何資料元素
- 範例：`法則 × 角色` 兩軸張出 system-laws 透鏡表

### 2. 對映 (Map Parts to Cells)

把每個部件放進格子。一個部件恰好一格；放不進去表示座標軸選錯，回到階段 1。

### 3. 負空間 (Negative Space)

空格即資訊。逐一空格判定：(a) 結構上不可能 (b) 存在但漏收 (c) 真正缺口。這是 aggregate 独有的價值 — 部件各自看不見「全體缺什麼」。

### 4. 產出整體 (Emit the Whole)

整體必須比部件之和多出東西：

- `索引 (index)`：可導航性
- `矩陣 (matrix)`：交叉洞察
- `總結 (summary)`：壓縮後的可傳遞性
- `本體論 (ontology)`：可驗證的座標系統（詳見 `references/ontology-template.md`）

## 典範範例：四大終極融合體

以 `融合法 × 角色 × 法則` 三軸聚合出的系統架構矩陣：

| 融合體 | 融合工具 | 角色 | 法則 | 系統意義 |
| :--- | :--- | :--- | :--- | :--- |
| 時空運行矩陣 | 正交性 | 建築師 ✖ 指揮家 | 空間 ✖ 時間 | 資源與時序互不干擾卻完美支撐 |
| 能量守恆引擎 | 同構性 | 推進者+修剪者+傳令官+指揮家 | 烈焰+冰霜+雷霆+重力 | 燃燒、凍結、傳導、引力統一調度 |
| 零信任稜鏡 | 催化劑 | 啟明者 ✖ 守密者 | 光明 ✖ 黑暗 | 以 Auth 為介面讓隱匿與可觀測共存 |
| 無限演化之輪 | 辯證循環 | 造物主→探險家→實驗家 | (精神+生命)→混沌→(因果+破壞) | 創造、未知與毀滅化為自我進化迴圈 |

架構完整性檢查：審視系統時確認四大融合體各有對應（命名空間與排程正交？算力與快取統一調度？監控與安全經 Auth 協作？CI/CD 成循環？）。

## 反模式

| 反模式 | 問題 | 修正 |
| :--- | :--- | :--- |
| 先收內容後定軸 | 格子互相重疊、無法歸位 | 座標軸先行 |
| 軸不正交 | 一個部件塞進多格 | 用「改 A 是否迫使改 B」驗軸 |
| 忽略空格 | 缺口不可見 | 負空間逐格判定 |
| 整體只是部件清單 | 沒有升維、沒有新資訊 | 產出交叉洞察或可導航結構 |

## 算子組合

`aggregate` 消費 `consolidate` 產出的乾淨部件；其負空間發現的缺口回饋給 `universal-generate`（生成缺件）— 這條回路由 `universal-evolve` 驅動。
