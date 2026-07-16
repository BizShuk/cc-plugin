---
name: session-retro
description: >
    Use when the user asks for a retrospective, post-mortem, or 復盤 of the current AI
    session — analyzing skills used, token cost per module, error rate and root causes,
    trust/review-boundary gaps, then turning findings into SMART goals (optionally executed).
    Triggers on "retro", "retrospective", "post-mortem", "復盤", "session 檢討",
    "分析這個 session", "what can be done better".
version: "1.0.0"
allowed-tools: Read, Write, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: high
metadata:
    type: reference
    platforms: [macos, linux]
---

# Session 復盤 (Session Retro)

## 概述 (Overview)

本技能指導 AI 以「效能審計員 (performance auditor)」角色，對`當前 session 的實際日誌`做深度復盤：
量化 skill/token/錯誤，歸因失效模式，稽核委託邊界，最後轉成可執行的 `SMART 目標`。
核心立場：`對數據誠實`——量不到就標為推估，不編造精確數字；不把 AI 的錯誤擬人化，
一律視為`防護欄缺失 (missing guardrails)` 或`提示詞工程瑕疵`。

## 使用時機 (When to Use)

- 使用者要求 `retro / 復盤 / post-mortem / 分析這個 session`。
- 一段密集協作（多功能、多輪迭代）後，要沉澱可重用資產與待改善項。
- 使用者問「哪裡可以做得更好」「token 花在哪」「錯誤率多少」。

## 核心規範 (Core Specification)

### 1. 資料蒐集 (Gather from Session Logs)

`必須從本 session 的實際過程取數`，不要泛談。可用來源：

- 對話中出現的 `Skill 呼叫`、`ToolSearch`、工具失敗訊息
- `git log` / `git diff --stat`（本 session 產出的檔案與行數）
- 截圖/驗證輪數（視覺任務的 token 大戶）
- 錯誤/退回/重試的具體對話段落

```bash
git -C <repo> log --oneline -20
git -C <repo> diff --stat HEAD~N   # 若已 commit
```

### 2. 四項統計指標 (The Four Metrics)

每次復盤`至少`回答這四題，能量化就量化：

| 指標                 | 內容                                    | 誠實要求                          |
| -------------------- | --------------------------------------- | --------------------------------- |
| Skills 使用          | 哪些 skill 被觸發、命中率、是否退回備援 | 標出空轉/失敗的 skill             |
| 可重用資產           | 可抽成 skill/module/memory 的產物       | 指向實際檔案路徑                  |
| Token 用量（模組級） | 各模組相對成本排序                      | `harness 不暴露精確值 → 用 proxy` |
| 錯誤率與清單         | 分母、錯誤數、逐條歸因                  | 區分「真錯誤」與「預期迭代」      |

`Token proxy 公式`：以 `截圖張數 × 對話往返 × 手寫 LOC` 推估相對成本，`明講是推估`。
常見洞察：大型資料檔用腳本`直接寫磁碟`可繞開輸出 token；截圖讀取常是真正大戶。

### 3. 失效模式歸因 (Failure Mode Attribution)

錯誤`先分兩層`，再歸因，`勿把預期的迭代（如視覺調校）算成錯誤`：

- `A 類 · 外部/環境`：資料源失敗、API 拒絕、extension 未連線、port 衝突
- `B 類 · 自引缺陷`：程式 bug、邏輯漏洞、API 誤用

歸因類別（技術導向，不擬人化）：
`上下文/提示詞不足` · `似真幻覺 (fail-plausible hallucination)` ·
`系統環境差異` · `人類過度信任而缺乏審查`。

每條錯誤記錄：`症狀 → 歸因 → 如何抓到`。特別標註`有無漏到「宣稱完成」之後`
（session 內閉環 = 健康；漏到事後 = 防護欄缺口）。

### 4. 委託審查 (Delegation Audit)

- 交給 AI 的任務是否`超出能力邊界`（如無真值可驗的判斷）？
- 人類是否對`高風險任務`（金流/刪除/對外發佈/權限）放棄實質把關？
- 哪個維度的`審查時間趨近於零`（例：視覺回歸只靠事後人眼，headless 下無法察覺）？

### 5. SMART 目標並執行 (SMART Goals, then Run)

每個建議轉成 `Specific / Measurable / Achievable / Relevant / Time-bound`。
`若使用者要求執行`，就實際做出來（建 skill、寫 harness、改文件），
`並回頭驗證`（跑一次、比對輸出），而非只列清單。落地產物優先於論述。

## 輸出格式 (Output Format)

先答使用者明確提問（結論優先），再依下列格式；`不使用粗體，一律 backtick 強調`：

```markdown
## 一、統計指標

（Skills 使用 / 可重用資產 / Token proxy / 錯誤率清單，用表格）

🤖 二、AI 協作健康度診斷
（效能提升場景、最大系統性摩擦、Token 成本浪費、品質風險）

⚠️ 三、信任與風險邊界評估
（審查機制漏洞、審查時間趨近零的環節、效能/覆蓋漂移、委託邊界是否合理）

## 四、SMART 目標

（每項含五要素；若執行，附產物路徑與驗證結果）
```

復盤全文另存 `docs/memory/YYYY-MM-DD-session-retro.md`，重點併入該 repo 的 `CLAUDE.md`。

## 常見錯誤 (Common Mistakes)

- 泛談「AI 很有幫助」而`未從本 session 實際取數`。
- 編造精確 token 數字，而非誠實標為 `proxy 推估`。
- 把`預期的迭代`（視覺調校、參數試錯）誤算成錯誤，灌水錯誤率。
- 把錯誤`擬人化`（「AI 粗心」），而非歸因到防護欄缺失/提示詞瑕疵。
- SMART 目標流於口號，未`轉成可執行動作`或未在使用者要求時實作與驗證。
- 忽略 `token 成本結構`（截圖迴圈、腳本落磁碟 vs 貼進 token）。
