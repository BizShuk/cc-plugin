---
name: topology-builder
description: >
    Use when asked to build a connection topology or knowledge graph across
    one or more sources — repo, folder, GitHub, database, articles, messages —
    producing one markdown file per entity (sections as dimensions, wikilink
    edges) under ~/projects/product/topologies. Triggers on: "build knowledge
    graph", "connection topology", "entity map", "service topology",
    "建立知識圖譜", "連結拓撲", "服務關聯圖".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit, Agent, Workflow, WebFetch
user-invocable: true
disable-model-invocation: false
effort: high
---

# topology-builder

以 agent team 從多種來源建立連結拓撲 (connection topology) 知識圖譜：
實體 (entity) = 檔案、維度 (dimension) = 章節、邊 (edge) = Obsidian wikilink。

## Overview

核心原則一：身分先於連結 (identity before edges)。先確立每個實體的唯一身分與
歸屬 zone，再萃取內部維度，最後才建邊；順序顛倒會產生重複實體與斷鏈。

核心原則二：證據先於邊 (evidence before edges)、信號先於噪音 (signal over noise)。
每條邊必須有發起方自身來源的直接佐證，禁止以命名相似、分層慣例或「看起來
合理」推斷；無佐證寧可不建。實體、維度與邊只收對「系統如何運作」有解釋力的
單元，基礎設施雜訊一律在 Discover 階段排除。

## When to Use

- 跨多來源（repo、folder、GitHub、資料庫、文章、訊息）建立知識圖譜或服務關聯圖
- 把散落的系統知識整理成可瀏覽、可回溯、可增量更新的拓撲文件庫
- 不適用：單一 repo 技術文件與業務價值分析（皆用 `project-docs`）、
  程式碼呼叫圖即時查詢（用 codegraph）

## Storage Layout

預設根目錄 `~/projects/product/topologies/`，使用者指定時以指定為準：

```text
topologies/
├── _index.md            # 全圖總覽：註冊表 + Mermaid + Frontier + Unlinked（檔尾）
├── payments/            # zone/wing/department 一律用資料夾
│   ├── service-a.md
│   └── billing-db.md
└── core/
    └── service-b.md
```

- 檔名 `<entity_name>.md`，kebab-case，與 frontmatter `name` 一致
- 檔名全域唯一（跨 zone 不可重複）— wikilink 以檔名解析，不含路徑
- 完整可驗證的小型範例見本技能目錄 `references/`
  （3 entities、2 zones，含 `_index.md`）

## Entity File Format

```markdown
---
name: service-a
type: service
zone: payments
tags: [billing, async]
aliases: [svc-a, billing-service]
sources:
    - type: repo
      ref: github.com/org/service-a
    - type: messages
      ref: "slack #payments 2026-05"
---

# Service A

一句話定位：替誰提供什麼能力。

## Method 1

kind: method

此維度做什麼（1~3 句）。

References:

- calls [[service-b#Method 2]] — 下游驗證
- writes-to [[billing-db#Invoices]]

## Billing Cycle

kind: concept

References:

- depends-on [[service-b#Method 3]]
- uses [[#Method 1]]

## External Sources

- [API 文件](https://example.com/api)

## Backlinks

<!-- auto-generated: do not hand-edit -->

- calls ← [[service-b#Method 1]]
```

- `type` 取值：`service` | `module` | `datastore` | `external-api` |
  `article` | `channel` | `team`
- 每個維度章節標題下第一行為 kind 標註，取值四種：
  `kind: concept`（領域概念/規則）、`kind: method`（操作/行為）、
  `kind: state`（狀態/生命週期）、`kind: interface`（對外介面）。
  不敷使用時選最接近者，於內文補充說明
- kind 以章節內容主體判定，不以實作形式判定：描述規則/政策者是
  `concept`，即使以函數實作；執行動作者是 `method`，即使動作對象是
  狀態資料（讀寫狀態的操作是 `method`，`state` 保留給狀態本身與生命週期）
- `## External Sources` 與 `## Backlinks` 是固定章節、不是維度：
  不加 kind、不計入 2~12 維度上限、不放 `References:` 邊
- `## Backlinks` 由 Phase 5 重建，標記註解必須保留

## Edge Format

- 邊一律寫在維度章節的 `References:` 清單：
  `- <relation> [[entity#Section]] — 補充說明（可省略）`
- 邊的方向 = 行為發起者 → 接受者。「被呼叫/被依賴」是 Backlinks 的職責，
  禁止寫成正向邊（`calls [[x]] — 被 x 呼叫` 即方向顛倒，屬錯誤）
- relation 動詞固定 kebab-case：`calls`, `uses`, `reads-from`, `writes-to`,
  `publishes-to`, `subscribes-to`, `depends-on`, `mentions`, `owned-by`；
  清單不敷使用時退用 `mentions`，並在補充說明寫明實際關係
- `Section` 必須與目標檔案的 `##` 標題逐字一致（含大小寫與空白）
- 同檔引用寫 `[[#Section]]`
- `References:` 只放 entity 維度之間的邊；引用目標不是 entity
  （設定鍵、環境變數、未生效的旗標等）時寫在維度內文描述，不建邊
- 外部 URL 不用 wikilink，放 `## External Sources` 或維度行內標準連結
- 基礎設施依賴過濾 (Dependency Filtering)：禁止將 logger、config、error
  wrappers、utils 等通用基礎設施元件建立正向邊。此類調用直接在維度的內文以
  文字描述，不放進 `References:`，避免拓撲圖充斥無業務/架構意義的雜訊邊

### 邊的真實性 (Edge Grounding)

- 每條邊必須有發起方「自身來源」的直接佐證：程式呼叫、import、SQL、設定鍵
  讀取、訊息內文。佐證寫進補充說明（`檔案:行`、API、表名）；指不出佐證即不建邊
- 禁止推斷式建邊：命名相似（`*Service` 互連）、分層慣例（同層必相連）、
  「應該會呼叫」皆非佐證。Phase 4 讀 sibling 檔只為抄標題，
  sibling 檔存在本身不是邊存在的依據
- 有弱訊號但無法確認的關係：不寫成正向邊，於維度內文以一句話描述，
  或列入 `_index.md` 的 `Frontier` 待下一輪確認
- `mentions` 是最後手段：每個維度 `mentions` 邊以 1~2 條為限；超出代表
  該關係應收斂為更精確的動詞，或根本不該建

## Identity Rules

- entity 是可獨立部署、維運或閱讀的單位：服務、系統、資料庫、文集、頻道、
  團隊。單一程式檔或函數永遠不是 entity，只能歸併為某 entity 的維度
- 薄轉接層（type alias、re-export、shim 檔）併入其包裝的實作 entity，
  原名進 `aliases`
- 每個 entity 維度章節 2~12 個；超過先做概念歸併，
  只剩 1 個時重新檢視是否應併入其他 entity
- canonical name 確立後同步寫入檔名與 frontmatter `name`；其餘稱呼進 `aliases`
- 去重針對「同一實體的多次發現」，不是「相似程式邏輯」：行為不同的平行實作
  各自成為所屬 entity 的維度，互以 `mentions` 交叉引用即可
- 噪音排除規則 (Noise Exclusion Rules)：
  - 排除任何測試檔案（例如 `*_test.go`, `*.spec.ts`）、模擬測試物件（mocks/fakes/stubs）與測試用 fixtures。
  - 排除單純的基礎設施與工具庫（如 `logger`, `json-parser`, `config-loader`, `db-driver` 等），它們不應被識別為 entity。
  - 排除標準庫與非專案特有的第三方依賴包（如 `cobra`, `lodash`）。
  - 對於不同 zone 的相似命名實體（例如 `payments/handler` 與 `auth/handler`），絕對不可因為名稱相似而合併為同一個 entity，必須以 context 區分並於檔名中加上 zone 區隔（例如 `payments-handler` 與 `auth-handler`）。

## Pipeline

| Phase      | 執行者                  | 產出                                                       |
| :--------- | :---------------------- | :--------------------------------------------------------- |
| 1 Discover | 每種來源一個 agent 平行 | entity 候選：name/type/zone 提示/aliases/證據              |
| 2 Identify | 主迴圈（不派 agent）    | 去重合併、定 canonical name 與 zone、建立 frontmatter 骨架 |
| 3 Extract  | 每個 entity 一個 agent  | 內部 tags/concepts/methods → 維度章節                      |
| 4 Connect  | 每個 entity 一個 agent  | 維度間有向邊、External Sources、標記 frontier              |
| 5 Verify   | 主迴圈 + 一個複核 agent | 斷鏈/孤兒/重名檢查、Backlinks 重建、`_index.md`            |

- Phase 3 與 Phase 4 之間有同步點：所有維度章節落檔後才開始建邊。
  Phase 4 引用標題時以磁碟上的目標檔案為唯一事實來源（Read 後逐字複製），
  禁止憑記憶拼寫——這就是同步點存在的原因
- 探索深度：以種子 entity 起算最多 2 hop；之外的新實體不建檔，
  列入 `_index.md` 的 `Frontier` 清單供下一輪輸入
- 無法派工的環境（單一執行者）依相同順序逐 phase 完成，規則不變；
  可在 Phase 4 順手寫 Backlinks，但 Phase 5 仍須由全圖正向邊重算覆寫

## Source Playbook

| 來源     | 切入方式                               | entity 候選                 |
| :------- | :------------------------------------- | :-------------------------- |
| repo     | entry points、部署單元、codegraph      | 服務、模組、資料庫          |
| folder   | 文件與設定掃描                         | 系統、文集                  |
| github   | `gh api`：README、issues、dependents   | 服務、external-api          |
| database | schema 與 table 清單                   | datastore（table 群為維度） |
| articles | markitdown / content-summarizer 轉文字 | article（論點為維度）       |
| messages | Slack / mail 匯出                      | channel、team（主題為維度） |

## Verification

Phase 5 必跑，不得以「圖很小」省略：

```bash
root=~/projects/product/topologies
# 1. 檔名全域唯一
find "$root" -name '*.md' ! -name '_index.md' -exec basename {} \; | sort | uniq -d
# 2. wikilink 斷鏈：目標檔或標題不存在
grep -rhoE '\[\[[^]]+\]\]' "$root" --include='*.md' | sort -u |
    sed 's/^\[\[//;s/\]\]$//' | while IFS='#' read -r name sec; do
    [ -z "$name" ] && continue
    f=$(find "$root" -name "$name.md" | head -1)
    [ -z "$f" ] && { echo "missing entity: $name"; continue; }
    [ -n "$sec" ] && ! grep -qxF "## $sec" "$f" && echo "missing heading: $name#$sec"
done
# 3. 維度章節標題下第一個非空行必須是合法 kind 標註
find "$root" -name '*.md' ! -name '_index.md' | while read -r f; do
    awk -v F="$f" '
        /^## /{h=substr($0,4); want=(h!="External Sources" && h!="Backlinks"); next}
        want && NF{ if ($0 !~ /^kind: (concept|method|state|interface)$/) print F": "h; want=0 }
    ' "$f"
done
```

另需完成：

- Backlinks 一致性：`## Backlinks` 只能由全圖正向邊重算推導，禁止手寫或
  憑記憶補。重算後若出現「無對應正向邊」的既有條目，即為錯誤——
  刪除它，或確認方向顛倒後改寫為正向邊
- 更新 `_index.md` 的註冊表與 Mermaid 總覽（zone 為 subgraph，
  邊聚合到 entity 層級）
- `_index.md` 必須以 `## Unlinked` 章節收尾，依方向分兩列，條目用
  `[[entity]]` wikilink。只計各維度 `References:` 清單中的跨實體正向邊；
  `Backlinks` 區段與同檔 `[[#Section]]` 一律不計
- `無入邊 (no inbound)` 列：沒有任何實體連向它；
  `無出邊 (no outbound)` 列：它未連向任何實體
- Unlinked 兩列同時出現者加註 `(orphan)`；任一列為空時明寫「無 (None)」
- 邊佐證抽查 (edge grounding audit)：複核 agent 隨機抽 N 條正向邊，要求指出
  發起方來源的佐證位置（`檔案:行`、表名、API）；指不出佐證者刪除該邊，
  或降級為維度內文描述。報告列出被刪/降級的邊
- 樞紐噪音偵測 (hub-noise detection)：入邊+出邊合計異常高、且關係多為
  `mentions` 的 entity 列為噪音候選；確認是基礎設施工具（logger、metrics）
  者拆除其邊或併為實作細節，報告供人工裁決

## Workflow Script

支援 Workflow tool 的環境可直接執行本目錄的 `workflow.js`：

```text
Workflow({
  scriptPath: "<skill_dir>/workflow.js",
  args: {
    root: "~/projects/product/topologies",
    sources: [{ "type": "repo", "target": "/path/to/repo" }]
  }
})
```

不支援時依 Pipeline 章節以 Agent tool 派工。

## Common Mistakes

| 錯誤                                           | 修正                                           |
| ---------------------------------------------- | ---------------------------------------------- |
| 一個程式檔或函數當成 entity（粒度爆炸）        | entity 是服務/系統層級；函數歸併入維度         |
| 用標準連結 `[B2](service-b.md#method2)` 表示邊 | 一律 `[[service-b#Method 2]]`，逐字對應標題    |
| 用檔名前綴（`cmd-export-*`）模擬分組           | zone 一律用實體資料夾                          |
| 身分資訊寫成自由文字 metadata                  | 一律 YAML frontmatter                          |
| 邊只有連結、沒有關係動詞                       | `- <relation> [[target#Section]] — 原因`       |
| 維度章節未標 kind 或自創取值                   | 標題下第一行 `kind:`，取值僅四種               |
| 跳過驗證與 Backlinks                           | Phase 5 必跑；Backlinks 是自動生成區段         |
| 把「被 X 呼叫」寫成正向邊 `calls [[x]]`        | 方向 = 發起者 → 接受者；反向關係交給 Backlinks |
| Backlinks 手寫、與正向邊不一致                 | 一律由全圖正向邊重算推導                       |
| 無限制向外爬連結                               | 2 hop 上限，之外進 Frontier                    |
| 誤將 common helper 或 logger 視為實體並建邊     | 排除 infrastructure noise，僅為實質業務/架構依賴建邊 |
| 不同 zone 的相似名稱 (例如 handler) 誤判合併   | 依據 business context 保持獨立，檔名與名稱加上 zone 做區隔 |
| 憑命名相似/分層慣例推斷邊（無來源佐證）         | 邊須有發起方來源佐證；指不出佐證即不建或降級為內文     |
| `mentions` 邊浮濫成為樞紐噪音                    | 每維度 `mentions` ≤ 2 條，否則收斂為精確動詞或刪除      |

## Failure Modes

| 情境                       | 動作                                                 |
| -------------------------- | ---------------------------------------------------- |
| 某來源不可達               | 該 source 標註 `status: unreachable`，不阻斷其他來源 |
| 跨 zone 同名但確為不同實體 | 名稱加 zone 後綴（`service-a-payments`）並記錄 alias |
| 目標目錄已有同名 entity 檔 | 讀取後增量合併，保留仍正確的維度與邊                 |
| 圖過大無法一輪完成         | 以 zone 分批；`_index.md` 記錄未涵蓋範圍與 Frontier  |
