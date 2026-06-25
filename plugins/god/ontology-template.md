# 系統本體論與資料定義範本 (System Ontology & Data Definition Template)

> 一份可直接放進 repo 的活文件 (living document)。
> 目的：把「結構掃得出來、語意挖不出來」的問題，變成一套可填、可驗證、可被 CI 守護的清單。
> 核心心法：**結構讓機器掃，語意靠 LLM 提問 + 人類回答，事實靠樣本資料驗證。**

---

## 0. 座標系統 (Coordinate System) — 先讀這段

任何一個資料元素的「全名」需要四個正交 (orthogonal) 座標：

```
env / domain / entity / attribute
例： prod / banking / Account / balance
```

| 座標 | 中文 | 它框住的是 | 一句話判別 |
|------|------|-----------|-----------|
| `env` | 環境 | **同義不同份** 的拷貝 | 互相複製資料，語意還成立 → env 差異 |
| `domain` | 領域 / 命名空間 | **同名不同義** 的字典 | 互相複製資料，語意壞掉、需翻譯 → domain 差異 |
| `entity` | 實體 | 一個有身分的東西 (table / struct) | — |
| `attribute` | 屬性 | 實體身上的一個欄位 | — |

> **env 與 domain 不是大環套小環，是兩條獨立的軸。**
> `prod.user_id` vs `dev.user_id` = 同一概念兩份拷貝（env）。
> `auth.account_id` vs `banking.account_id` = 同名兩種概念（domain）。

### 建議的 repo 檔案佈局 (suggested layout)

```
/ontology
  ├── README.md                  ← 本檔（說明 + 範本）
  ├── containers.yaml            ← §4a 容器層級
  ├── elements.yaml              ← §4b 外在元素類型 + §4c 元素透鏡
  ├── env-shape.yaml             ← §5 環境標準形狀 + diff
  ├── lens-city.yaml             ← §6 城市透鏡
  └── /domains
        ├── auth.glossary.yaml   ← §2 各領域一份詞彙表
        ├── banking.glossary.yaml
        ├── _identifiers.yaml    ← §3 跨領域識別碼基數表
        └── _context-map.yaml    ← 領域間映射
  └── /catalog
        └── banking.Account.yaml ← §1 每個實體一份資料定義
```

---

## 1. 資料定義模板 (Data Element Definition)

> 每個資料元素（一個 storage 欄位、一個 config key、一個 method 的輸入/輸出）的「身分證」。
> 最值錢的兩欄是 `semantic`（語意）和 `confidence.unknowns`（已知的未知）。

### 1a. Schema（空白範本，複製用）

```yaml
# catalog/<domain>.<Entity>.yaml
entity: ""                    # 實體名 entity name
domain: ""                    # 屬於哪個領域 domain / namespace
location: ""                  # 程式位置：service / table / struct
description: ""               # 實體層級的一句話描述
attributes:
  - name: ""                  # 1. 欄位名 name
    semantic: ""              # 2. 語意 semantic：代表現實中的什麼（最容易缺）
    type: ""                  # 3. 型別 type
    value_domain: ""          # 3. 值域 value domain：0/1 欄位務必寫「0=? 1=?」
    category: ""              # 4. 類別：flag | accumulator | immutable | metadata | reference | mutable
    lifecycle:                # 5. 生命週期 lifecycle
      created_by: ""          #    誰建立
      updated_by: ""          #    誰更新
      update_when: ""         #    何時更新
      deletable: false        #    能否刪除
    mutability: ""            # 6. immutable | append-only | mutable
    invariants: []            # 7. 不變式：永遠要成立的規則
    source: []                # 8. 寫入者 service/method（資料血緣上游）
    sink: []                  # 8. 讀取者 service/method（資料血緣下游）
    env_diff: ""              # 9. dev 與 prod 語意/值是否不同
    confidence:               # 10. 信心與未知 ← 隱藏知識先停放在這
      level: ""               #     high | medium | low | guessed-by-llm
      unknowns: []            #     明確列出「我不確定的部分」
```

### 1b. 填好的範例（你的 banking 情境）

```yaml
# catalog/banking.Account.yaml
entity: Account
domain: banking
location: services/banking/internal/store/account.go (struct Account)
description: 銀行帳戶；一個客戶可持有多個帳戶 (1 user : N accounts)
attributes:
  - name: balance
    semantic: 帳戶當前餘額，單位為「分」(cents)，避免浮點誤差
    type: int64
    value_domain: ">= 0；透支帳戶可為負，見 invariants"
    category: accumulator
    lifecycle:
      created_by: AccountService.Open()
      updated_by: LedgerService.Post()   # 只有分錄服務能改
      update_when: 每筆交易入帳時
      deletable: false
    mutability: mutable
    invariants:
      - "balance 的變動必須等於對應 ledger entries 的總和（雙式簿記）"
      - "非透支帳戶 balance 永不為負（DB constraint + service guard）"
    source: [LedgerService.Post]
    sink: [StatementService, FraudService, MobileAPI]
    env_diff: "dev 使用測試貨幣與假交易；prod 為真實金流"
    confidence:
      level: medium
      unknowns:
        - "透支帳戶的下限是否有規則？目前程式沒看到 floor 檢查"
        - "balance 是否曾被批次校正腳本直接 UPDATE 過？需查稽核紀錄"

  - name: is_frozen
    semantic: 帳戶是否被凍結（不可交易）
    type: bool / tinyint(1)
    value_domain: "0 = 正常可交易；1 = 已凍結禁止出金"   # 0/1 一定要寫清楚
    category: flag
    lifecycle:
      created_by: "預設 0"
      updated_by: [ComplianceService, FraudService]
      update_when: 風控或法遵觸發
      deletable: false
    mutability: mutable
    invariants:
      - "is_frozen=1 時，任何出金 method 必須拒絕"
    source: [ComplianceService, FraudService]
    sink: [LedgerService, MobileAPI]
    env_diff: "語意相同"
    confidence:
      level: low
      unknowns:
        - "解凍是誰有權限？流程沒寫下來（隱藏知識）"
        - "是否有第三態（如 partial freeze）被偷偷用其他欄位表達？"
```

---

## 2. 分領域詞彙表 (Per-Domain Glossary)

> 每個**限界上下文 (Bounded Context)** 一份。同一個詞在不同領域分開定義。
> 專門用來抓兩個語意陷阱：**同形異義 (homonym，同名不同義)**、**異形同義 (synonym，異名同義)**。

### 2a. Schema

```yaml
# domains/<domain>.glossary.yaml
domain: ""                    # 領域名
purpose: ""                   # 這個領域負責什麼
terms:
  - term: ""                  # 詞
    means_here: ""            # 在「本領域」的意思
    not_to_confuse_with: []   # 同形異義警告：別跟哪個領域的同名詞搞混
    synonyms: []              # 異形同義：本領域內/外的別名
```

### 2b. 填好的範例（auth vs banking 的 `account`）

```yaml
# domains/auth.glossary.yaml
domain: auth
purpose: 處理登入身分與授權
terms:
  - term: account
    means_here: "一組登入憑證（帳號 + 密碼），等同一個可登入的身分"
    not_to_confuse_with:
      - "banking.account = 金融帳戶（與登入身分非 1:1）"
    synonyms: [login, credential]
  - term: user_id
    means_here: "登入身分的唯一識別碼"
    synonyms: []
```

```yaml
# domains/banking.glossary.yaml
domain: banking
purpose: 處理資金、帳戶與分錄
terms:
  - term: account
    means_here: "一個金融帳戶（支票/儲蓄/信用卡），持有餘額"
    not_to_confuse_with:
      - "auth.account = 登入身分（一個人可有多個金融帳戶）"
    synonyms: []
  - term: customer
    means_here: "持有一個或多個帳戶的客戶（對映到 auth.user）"
    synonyms: [account_holder]
```

> **`bank` 法則**：同一個拼字，是領域在決定它的意思（河岸 vs 銀行）。詞彙表就是「指定查哪本字典」。

---

## 3. 跨領域識別碼基數表 (Cross-Domain Identifier Cardinality)

> 回答你的核心問題：`user_id == account_id` 嗎？
> **這不是普世事實，是領域屬性。** 對每個跨領域共用 ID，逼問：相等是「不變式」還是「巧合」？

### 3a. Schema

```yaml
# domains/_identifiers.yaml
identifiers:
  - concept_a: ""             # 識別碼 A（含領域）
    concept_b: ""             # 識別碼 B（含領域）
    cardinality: ""           # 1:1 | 1:N | N:M
    holds_in: []              # 在哪些領域/情境成立
    breaks_in: []             # 在哪些領域/情境不成立
    relationship_kind: ""     # invariant（機制保證） | coincidence（剛好成立）
    linking_table: ""         # N:M 時，靠哪張關係表/欄位連起來
    confidence: ""            # high | medium | low
    notes: ""
```

### 3b. 填好的範例

```yaml
# domains/_identifiers.yaml
identifiers:
  - concept_a: auth.user_id
    concept_b: banking.account_id
    cardinality: "1:N"
    holds_in:
      - "簡單 app：一人一帳號，1:1 巧合成立"
    breaks_in:
      - "banking：一個 user 有支票/儲蓄/信用卡多個 account"
      - "joint account：N users : M accounts（聯名戶）"
    relationship_kind: coincidence   # 小領域裡剛好成立，長大就破
    linking_table: "banking.account_holders (user_id, account_id, role)"
    confidence: high
    notes: >
      千萬別在程式裡寫死 user_id == account_id。
      Costco 家庭卡比喻：一張卡(account)掛多人(users)，
      一個人(user)也可有多家銀行卡(accounts)。

  - concept_a: auth.user_id
    concept_b: banking.customer.id
    cardinality: "1:1"
    holds_in: ["全系統"]
    breaks_in: []
    relationship_kind: invariant     # 由註冊流程機制保證
    linking_table: "banking.customer.auth_user_id (FK + unique)"
    confidence: medium
    notes: "需驗證：是否真有 unique constraint，還是只靠應用層約定？"
```

---

## 4. 頂級元素推論 (Top-Level Element Inference)

> 在「某個節點」上，用幾大頂級元素去**象徵這裡應該會有的東西**。
> 缺了的格子 = 隱藏或缺失元素的候選 (candidate)。分三層：容器、外在元素類型、元素透鏡。

### 4a. 容器層級 (Containment Levels)

> 像俄羅斯娃娃。**環境是服務的容器，不是兄弟。** 服務不會自己知道環境，是啟動時被注入 (injected) 的。

```yaml
# containers.yaml
containers:
  - level: platform           # 平台 / 宇宙
    contains: [environment]
  - level: environment        # 環境（dev/prod）— 由 ENV 注入，是「框」
    injected_via: "env var ENV / config injection / build flag"
    contains: [service]
    instances: [dev, staging, prod]
  - level: service            # 服務
    contains: [method, internal_state]
    consumes: [storage, external_service, middleware, config, sdk]  # 外圍資源
  - level: method             # 方法
  - level: internal_state     # 內部狀態
```

> **原則**：環境是「同義不同份」的容器；每個實體身上的「環境戳記」是容器蓋上去的，不是自己長出來的。

### 4b. 外在元素類型 (Element Types at the Service Layer)

```yaml
# elements.yaml — part 1
element_types:
  - type: method              # 運算邏輯
  - type: internal_state      # 內部狀態
  - type: contract            # 對外契約（interface），Go 的天然關聯邊界
  - type: middleware          # 連線管道
  - type: storage             # 資源依賴（視為 external service）
  - type: external_service    # 外部服務
  - type: config              # 設定
  - type: metadata            # 元資料
```

### 4c. 元素透鏡 (Elemental Lens) ★ 你補充的玩法

> 用「空間/時間/光/暗/五行/波」當**最原始的角色檢查表**。
> 在任一節點問：「我的『X 元素』在哪？」缺格即候選隱藏元素。比城市透鏡更抽象，能抓到城市透鏡漏掉的東西。

```yaml
# elements.yaml — part 2
elemental_lens:
  - element: 空間 Space
    symbolizes: 容器 / 位置 / 容量
    maps_to: "storage location, memory, capacity, namespace"
    missing_means: "資料無處可放 / 容量上限未定義"
  - element: 時間 Time
    symbolizes: 流轉 / 順序
    maps_to: "timestamps, TTL, versioning, ordering, scheduling, lifecycle"
    missing_means: "無法追溯何時、無法保證順序"
  - element: 光 Light
    symbolizes: 可見 / 外顯
    maps_to: "logging, metrics, tracing, public API, observability"
    missing_means: "黑盒，出事看不到（缺可觀測性）"
  - element: 暗 Dark
    symbolizes: 隱藏 / 未知
    maps_to: "secrets, private state, hidden invariants, unknown unknowns"
    missing_means: "機密沒管理 / 隱藏知識沒被攤開"
  - element: 火 Fire
    symbolizes: 運算 / 轉化
    maps_to: "compute, business logic, transformation methods"
    missing_means: "只存不算，沒有處理能力"
  - element: 水 Water
    symbolizes: 流動 / 傳遞
    maps_to: "data flow, streams, pipelines, messaging, ETL"
    missing_means: "資料卡住不流動"
  - element: 土 Earth
    symbolizes: 根基 / 持久
    maps_to: "persistent storage, stable schema, immutable base"
    missing_means: "沒有真相來源 (source of truth)"
  - element: 木 Wood
    symbolizes: 生長 / 累積
    maps_to: "accumulators, append-only logs, growth, scaling"
    missing_means: "無法累積歷史"
  - element: 雷 Thunder
    symbolizes: 觸發 / 突變
    maps_to: "events, triggers, interrupts, alerts, webhooks"
    missing_means: "只能輪詢 (polling)，無法即時反應"
  - element: 金 Metal
    symbolizes: 契約 / 邊界
    maps_to: "interfaces, schemas, type system, validation, auth boundary"
    missing_means: "邊界鬆散，契約不明"
  - element: 波 Wave
    symbolizes: 傳播 / 非同步
    maps_to: "async, pub-sub, eventual consistency, propagation delay"
    missing_means: "誤把分散式當同步，忽略延遲與失序"

# 進階：用五行相生相剋看「元素類別之間的張力」（可選）
wuxing_relationships:
  generating(相生):     # 健康的資料流循環
    - "木生火：累積資料餵給運算"
    - "火生土：運算產出落地為持久狀態"
    - "土生金：穩定 schema 定義出契約"
    - "金生水：契約讓資料得以流動"
    - "水生木：流動的資料回頭餵養累積"
  overcoming(相剋):     # 系統性張力 / 需設限之處
    - "金剋木：僵硬契約限制無界成長 → schema 版本化"
    - "木剋土：失控累積侵蝕儲存 → 保留策略 retention"
    - "土剋水：過度持久阻塞流動 → 讀寫分離"
    - "水剋火：流量淹沒運算 → 背壓 backpressure"
    - "火剋金：重運算熔毀契約 → breaking change 控管"
```

> 用法提示語（接你之前的 LLM 發掘流程）：
> 「把這個 service 當作一個節點，用元素透鏡逐一檢查：空間/時間/光/暗/火水土木雷金波 各對應到它的哪個元件？哪些元素是空的？每個空格判斷它是 (a) 真的不需要、(b) 存在但我沒提到、(c) 真正缺口，各標信心。」

---

## 5. 環境標準形狀 + Diff 檢查表 (Env Standard Shape + Diff)

> **平行宇宙原則**：dev 與 prod 應「物理法則相同、內容物不同」，即結構同構 (structurally isomorphic)。
> 每個結構差異只有兩種命運：**故意的（必須記錄）** 或 **沒人記得的（風險旗標）**。

### 5a. 標準形狀 Schema

```yaml
# env-shape.yaml — part 1
standard_shape:                # 一個「合格環境」該有的元素類型清單
  required_elements:
    - storage
    - cache
    - message_queue
    - secrets_manager
    - logging
    - external_service_clients
    - config_source
  note: "每個 environment 都應具備上述每一類至少一個實例"
```

### 5b. Diff 檢查表 Schema

```yaml
# env-shape.yaml — part 2
env_diff_checklist:
  - element: ""               # 哪個元素類型
    in_dev: ""                # dev 的實例
    in_prod: ""               # prod 的實例
    structurally_same: true   # 結構是否同構
    difference: ""            # 差在哪
    verdict: ""               # intentional（已記錄原因） | unknown（待查，風險）
    reason_or_risk: ""
```

### 5c. 填好的範例

```yaml
env_diff_checklist:
  - element: cache
    in_dev: "無（直接打 DB）"
    in_prod: "Redis 叢集"
    structurally_same: false
    difference: "prod 多一層快取，dev 沒有"
    verdict: unknown          # ← 經典「dev 正常 prod 爆」來源
    reason_or_risk: "快取失效/穿透問題只會在 prod 出現，dev 測不到"
  - element: secrets_manager
    in_dev: ".env 檔"
    in_prod: "Vault"
    structurally_same: false
    difference: "取密機制不同"
    verdict: intentional
    reason_or_risk: "dev 求方便、prod 求安全；已記錄"
  - element: storage
    in_dev: "Postgres (docker)"
    in_prod: "Postgres (managed)"
    structurally_same: true
    difference: "僅資料內容不同（env 差異，非 domain 差異）"
    verdict: intentional
    reason_or_risk: "符合平行宇宙原則"
```

---

## 6. 城市透鏡對映表 (City Lens Mapping)

> 用一個**完整、運作正常**的領域當比喻，把元件一一對映。比喻裡有、你系統對不上的格子 = 候選漏網元素。
> 城市透鏡最直觀，建議第一個用。

```yaml
# lens-city.yaml
lens: city
roles:
  - role: 居民 residents
    maps_to: services
    present: true
    note: ""
  - role: 道路 roads
    maps_to: "network / middleware / message bus"
    present: true
  - role: 水電 utilities
    maps_to: "config / secrets"
    present: true
  - role: 倉庫 warehouses
    maps_to: storage
    present: true
  - role: 海關邊境 customs/border
    maps_to: "對外 API 邊界 / gateway"
    present: unknown          # ← 對映後發現空格就標 unknown
  - role: 警察 police
    maps_to: "authn / authz 驗證授權"
    present: unknown          # 常見漏網：有倉庫卻沒警察
    candidate_hidden_element: true
  - role: 市政檔案 city records
    maps_to: metadata
    present: true
  - role: 電網 power grid
    maps_to: "基礎設施 / 排程 / 健康檢查"
    present: unknown

# 其他可換用的透鏡（每個都是一份角色檢查表）
alternative_lenses:
  organism: "器官=service, 血液=資料流, 神經=事件佇列, 免疫=安全, DNA=schema/config, 記憶=storage, 代謝=背景排程"
  theatre: "演員=service, 劇本=config, 舞台=環境, 道具=資源, 導演=orchestrator, 觀眾=clients"
```

---

## 透鏡如何協作 (How the Lenses Work Together)

| 透鏡 | 抓什麼 | 強項 |
|------|--------|------|
| 元素透鏡 (§4c) | 最原始的功能缺口（連時間/可觀測性都沒有） | 最抽象，抓深層缺漏 |
| 城市透鏡 (§6) | 角色缺口（有倉庫沒警察） | 最直觀，易上手 |
| Env diff (§5) | 環境間的結構不對稱 | 抓「dev 正常 prod 爆」 |
| 詞彙表 + 基數表 (§2,§3) | 語意衝突（同名不同義、假相等） | 抓跨領域整合地雷 |

> **共同鐵則**：以上所有由 LLM 產出的對映、猜測、空格，全部是**假設 (hypotheses)**，不是事實。
> 最後一定要拿去對程式碼、對樣本資料、對 env diff、對團隊驗證。LLM 是「會問好問題的新顧問」，不是下結論的人。

---

## 落地計畫 (Execution Plan)

1. **選試點**：挑一個最關鍵或最混亂的 service，別一次吃整個系統。
2. **建座標**：先填 `containers.yaml`，確認環境注入點（ENV 從哪讀）。
3. **機械盤點**：用 Go AST / grep 掃 struct tag、interface、config key，灌進 §1 模板的可自動填欄位。
4. **切領域**：建各領域 `*.glossary.yaml`（§2），標出同形異義衝突。
5. **盤點 ID**：填 `_identifiers.yaml`（§3），每個共用 ID 標「不變式 vs 巧合」。
6. **跑透鏡**：依序套元素透鏡（§4c）→ 城市透鏡（§6），把空格交給 LLM 做負空間推理產生問題清單。
7. **Env diff**：填 `env-shape.yaml`（§5），每個差異標 intentional / unknown。
8. **人工補語意**：把 LLM 的問題清單丟給懂的人，回填 `semantic` / `invariants` / `unknowns`。
9. **驗證**：跑 query 對樣本資料反驗（0/1 分布、immutable 沒被改、accumulator 單調遞增）。
10. **沉澱 + CI**：把成果存進 repo，加 CI 檢查防止 schema 改了文件沒跟（避免定義漂移 definition drift）。
11. **換下一個 service，重複。**
