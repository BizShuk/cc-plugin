# Refactor Agent vs. 多個專門 Agent 的架構抉擇

這是一個非常棒的架構設計問題,涉及 **Agent Design Pattern(代理設計模式)** 的核心權衡。讓我用實際情境幫你分析。

---

## 兩種方案的本質差異

### 方案 A:單一 Refactor Agent + 多個 Skills

```tree
Refactor Agent
  ├── skill: mvc-patterns
  ├── skill: golang-performance
  ├── skill: database-optimization
  └── skill: code-smell-detection
```

**生活比喻**:像一位**全能型家庭醫師(General Practitioner)**,看診時根據症狀翻閱不同的醫學手冊——感冒翻呼吸科手冊,胃痛翻消化科手冊。一個人,多本參考書。

### 方案 B:多個專門 Agent

```tree
MVC Agent ─────┐
                ├──→ Orchestrator(協調者)
Performance ───┘
   Agent
```

**生活比喻**:像**醫院的專科醫師團隊**——心臟科醫師、腸胃科醫師各自獨立看診,複雜病例由主治醫師(orchestrator)協調會診。

---

## 關鍵權衡比較

| 面向                         | 方案 A:單一 Agent + Skills             | 方案 B:多個專門 Agents                    |
| ---------------------------- | -------------------------------------- | ----------------------------------------- |
| **Context 共享(上下文)**     | ✅ 同一個 context,跨領域知識能融會貫通 | ❌ Agent 間需明確傳遞,易遺漏細節          |
| **Token 成本**               | ✅ 較低,只載入相關 skill               | ❌ 較高,每個 agent 有自己的 system prompt |
| **專精程度(Specialization)** | ⚠️ 中等,受限於主 agent 的整體判斷      | ✅ 高,每個 agent 可深度調校               |
| **維護性(Maintainability)**  | ✅ 集中管理,改一處生效                 | ⚠️ 多處維護,但邊界清晰                    |
| **並行處理(Parallelism)**    | ❌ 序列處理                            | ✅ 可同時跑多個 agent                     |
| **錯誤隔離**                 | ❌ 一個 skill 出錯影響整體             | ✅ 單一 agent 失敗不拖垮其他              |
| **跨領域推理**               | ✅ 強(看到全貌)                        | ⚠️ 弱(各自為政)                           |

---

## 你的場景:重構(Refactoring)的特殊性

**重構任務有一個關鍵特性**:**跨領域決策耦合度高(High Cross-domain Coupling)**。

舉例說明: 假設你要重構一個 Golang MVC 應用的「使用者查詢 API」:

```go
// 原始碼
func GetUser(c *gin.Context) {
    userID := c.Param("id")
    var user User
    db.Where("id = ?", userID).First(&user)
    db.Model(&user).Related(&user.Orders)  // N+1 query!
    c.JSON(200, user)
}
```

這段程式碼同時牽涉:

- **MVC 層次**:Controller 直接操作 DB,違反分層
- **Golang 效能**:可能有 goroutine 利用機會
- **DB 效能**:N+1 query 問題
- **錯誤處理**:沒有處理 error

### 方案 A 的處理方式(推薦給此情境)

單一 Refactor Agent 看到完整程式碼,**一次性決策**:

> 「我要把 DB 查詢移到 Repository 層(MVC),同時改用 Preload 解決 N+1(效能),並加上 context 控制超時(Golang 最佳實務)」

→ **一個連貫的重構方案**

### 方案 B 的處理方式

```
MVC Agent: 「請把 DB 邏輯抽到 Repository」
   ↓
Performance Agent: 「等等,Repository 的查詢還是 N+1」
   ↓
MVC Agent: 「那要不要改回來?」
   ↓
... 互相打架 🥊
```

→ 需要強力的 Orchestrator 來協調衝突,**架構複雜度暴增**

---

## 我的推薦:混合架構(Hybrid Approach)

對於你的場景,我建議**「方案 A 為主,方案 B 為輔」**的混合策略:

### 推薦架構

```
┌─────────────────────────────────────┐
│      Refactor Agent (主導者)         │
│                                      │
│   Skills(知識庫):                   │
│   ├── mvc-patterns.md               │
│   ├── golang-idioms.md              │
│   ├── golang-performance.md         │
│   ├── db-optimization.md            │
│   └── testing-strategies.md         │
│                                      │
│   可選召喚的 Sub-agents:             │
│   ├── 🔬 Benchmark Agent(跑壓測)    │
│   └── 🧪 Test Generation Agent      │
└─────────────────────────────────────┘
```

**設計理念**:

- **知識用 Skills**:MVC 規則、Golang 最佳實務這些是**靜態知識**,適合 Skills
- **行動用 Sub-agents**:需要實際**執行任務**(跑 benchmark、產生大量測試)才獨立成 agent

---

## 何時該選方案 B(多個專門 Agents)

如果符合以下**任一**情況,選方案 B:

1. **任務本身可並行**:例如同時對 100 個 microservices 做掃描
2. **領域邊界極清晰**:例如「前端 agent」vs「後端 agent」幾乎不重疊
3. **不同 agent 需要不同模型**:例如效能調校用 Opus,簡單格式化用 Haiku 省錢
4. **需要不同權限**:例如 DB agent 才有 production 存取權
5. **團隊分工開發**:不同小組各自維護自己的 agent

---

## 給你的實作計畫

### 階段一:從單一 Refactor Agent 開始(第 1-2 週)

建立基礎結構:

```
refactor-agent/
  ├── SKILL.md                    # 主指引
  ├── skills/
  │   ├── mvc-golang/
  │   │   └── SKILL.md
  │   ├── golang-performance/
  │   │   ├── SKILL.md
  │   │   ├── concurrency.md
  │   │   ├── memory-allocation.md
  │   │   └── profiling-guide.md
  │   └── db-optimization/
  │       └── SKILL.md
  └── examples/                    # 重構前後對照範例
```

**先讓單一 agent 跑得起來**,觀察它的表現。

### 階段二:識別瓶頸(第 3-4 週)

實際使用後記錄問題:

- 是否常常「忘記」某個 skill?→ description 寫得不夠精準
- 是否 context 爆掉?→ 考慮拆分
- 是否某類任務特別慢?→ 考慮獨立成 sub-agent

### 階段三:選擇性拆分(第 5 週後)

只有當**證據明確**時才拆 agent:

- ✅ 「benchmark 跑很久,拖累對話」→ 拆出 Benchmark Agent
- ✅ 「測試生成需要不同的 prompt 風格」→ 拆出 Test Agent
- ❌ 「MVC 和效能感覺很不同」→ 這只是直覺,先別拆

---

## 核心原則:YAGNI for Agents

軟體工程有個原則叫 **YAGNI(You Aren't Gonna Need It,你不會需要它)**。套用到 agent 架構:

> **「先用最簡單的單一 agent + skills 解決,直到證明它真的不夠用」**

過早拆分多個 agents 會帶來:

- 🐛 **協調 bug**:agent 間溝通失敗
- 💸 **token 浪費**:重複傳遞 context
- 🤯 **除錯困難**:出錯時不知道是哪個 agent 的鍋

---

## 總結回答

| 你的問題        | 我的建議                                                  |
| --------------- | --------------------------------------------------------- |
| 哪個比較好?     | **方案 A(單一 Refactor Agent + 多 Skills)更適合重構場景** |
| 為什麼?         | 重構需要跨領域決策,單一 context 才能做出連貫判斷          |
| 何時考慮方案 B? | 任務可並行、領域真正獨立、或需要不同模型/權限時           |
| 最佳實踐?       | **混合架構**:知識用 Skills,獨立執行任務才拆 Sub-agent     |

要不要我幫你設計一個具體的 `SKILL.md` 範本,展示如何把 MVC + Golang Performance 整合在同一個 Refactor Agent 裡?
