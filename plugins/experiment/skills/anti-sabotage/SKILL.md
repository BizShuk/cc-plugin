---
name: anti-sabotage
description: >
    Use when reviewing changes whose tests may conceal production-only behavior,
    runtime condition branches, or cross-system state risks. Triggers on:
    "environment parity", "dormant branch", "test rigging", "internal sabotage",
    "測試造假", "內部破壞".
version: "1.0.0"
allowed-tools: Read, Bash, Grep, Glob
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: review
    platforms: [macos, linux]
---

# Skill：偵測與預防「測試造假型內部破壞」

(Detecting & Preventing Test-Rigging Insider Sabotage)

## 目的

在 review / QA 階段，抓出「測試會過、但真實條件下會壞」這類被`刻意隱藏`的行為。核心信念：`「測試通過」只代表「在我們測過的條件下沒事」，不等於「行為正確」。`

## 何時啟用（Trigger）

改動只要符合任一項，就套用這份清單：

- 依`執行期環境變數`（region、tenant、時間、流量百分比…）來分流的邏輯
- 會`寫入或改變跨系統共享狀態`的程式
- 作者`同時主導了實作與測試設計`

---

## 五項檢查（對應五層手法）

`1. 環境一致性 (Environment Parity)`
問：「測試環境跟正式環境在 region routing 上一模一樣嗎？」

- 要求測試矩陣明確覆蓋`真實的 region 值`，而不只是 over-range / 邊界值。
- 對任何「只在某條件下分流」的程式，強制要有一個會走到`正式環境那條路徑`的測試。

`2. 測試與實作分離 (Separation of Duties)`
寫這段邏輯的人，不該是唯一決定怎麼測它的人。

- QA 或第二人要能`獨立設計`測試、並能質疑測試矩陣的缺口。
- 必問一句：「有哪些`真實輸入`，是這組測試保證不會碰到的？」

`3. 找休眠的條件分支 (Dormant Branch)`
Review 時專門找「依執行期條件才會走的分支」。

- 對每個分支問：「正式環境下這條件什麼時候會成立？那條路被測到了嗎？」
- `沒被測到的條件分支，預設視為「未驗證」`，不是「應該沒事」。

`4. 跨系統狀態閉環 (Cross-system State Closure)`
只要改動會碰到共享狀態，文件就必須回答：

- 「還有哪些系統會`讀或改`這個狀態？這次改動之後它們的行為是什麼？」
- 由`負責那些下游系統的人 sign-off`。沒寫下游影響的，當作「沒做完」。

`5. 滅知識單點 (Reduce Bus Factor)`

- 任何「只有作者懂」的行為，要寫成文件並由`第二人複述一遍`確認理解一致。
- 如果一段功能只有一個人能解釋它的真正行為，這本身就是「升級審查」的訊號。

---

## 流程面的護欄（Process Guardrails）

- `金絲雀／灰度發布 (Canary / Staged Rollout)` ＋ `依 region 維度的監控`：先放一小部分真實流量，盯著那個 region 的錯誤率，不要一次全量上線。
- `上線後真實環境冒煙測試 (Production smoke test)`：用真實 region 值跑一遍關鍵路徑。
- `變更紀錄強制「下游影響」欄位`：留白不給過。
- `文化上把「通過」與「正確」分開講`：通過 = 「在我們測過的條件下沒事」，僅此而已。

---

## 一句話總結

> 內部破壞的招數，幾乎都是利用「`測試條件 ≠ 真實條件`」與「`沒人負責系統交接點`」這兩個縫隙。把這兩個縫隙補起來，大部分這類手法就現形了。
