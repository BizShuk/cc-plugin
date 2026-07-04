---
name: kb-spec
description: >
    Reference spec for the knowledge base storage layout, capture/entity file
    formats, truth tiers, fingerprint dedup, and edge rules. Use when reading or
    writing any file under the knowledge base root, or when another kb-* skill
    needs the canonical format. Triggers on: "kb spec", "knowledge base format",
    "知識庫格式", "知識庫規範".
version: "1.0.0"
user-invocable: false
disable-model-invocation: true
metadata:
    type: reference
---

# kb-spec — 知識庫規範 (Knowledge Base Spec)

所有 `kb-*` 技能共用的單一事實來源：儲存佈局、檔案格式、真實性分級、去重與邊規則。

## 儲存佈局 (Storage Layout)

兩層結構：狀態集中在全域根、知識依專案分庫。使用者指定路徑時以指定為準。

路徑佔位符（全部 kb-* 技能通用）：

- `<kb>` = 全域根 `~/projects/product/`（狀態與跨專案總覽）
- `<proj>` = 專案知識庫 `~/projects/product/projects/<project>/`
  （`<project>` = repo 資料夾名或使用者命名，kebab-case）

```text
~/projects/product/            # <kb> 全域根
├── _index.md                  # 跨專案總覽：專案註冊表 + 各庫健康度
├── _state/                    # 全部 run 的狀態（集中一處，不分散到各專案）
│   ├── STATUS.md              # 儀表板：所有 run 的階段與進度
│   ├── cache/                 # 來源快取（如 git log 全文，只追加）
│   │   └── history-service-a.log
│   ├── verify/                # 健檢報告（<date>-<project>.md）
│   └── runs/<run-id>/
│       ├── manifest.json      # 不可變工作計畫：item 清單 + 分批
│       ├── progress.json      # 可變進度：每批狀態與計數
│       └── log.md             # 人可讀事件記錄（錯誤、跳過原因）
└── projects/
    └── service-a/             # <proj> 專案知識庫（每個 repo 一個）
        ├── _index.md          # 專案內總覽：註冊表 + Mermaid + Frontier + Unlinked
        ├── _inbox/            # raw captures，待蒸餾
        │   └── 2026-07-04-payment-flow.md
        ├── _sources/          # 來源登記，供 corroboration 計數
        │   └── repo-service-a.md
        └── payments/          # zone 一律用資料夾
            └── service-a.md   # curated entity
```

- entity 檔名 `<entity_name>.md`，kebab-case，同一 `<proj>` 內跨 zone 全域唯一；
  wikilink 只在 `<proj>` 內解析，不跨專案
- capture 檔名 `<yyyy-mm-dd>-<slug>.md`
- `run-id` 一律含專案名（`<date>-<skill縮寫>-<project>`），全域 `_state/`
  才分得開各專案的 run
- web/chat/schema 來源歸屬：使用者指定專案，或依 `zone-hint` 對應；
  無法歸屬時建立獨立專案資料夾
- `_inbox/`、`_sources/`、`_state/` 內的檔案不是 entity，不參與 wikilink 邊統計

## 狀態追蹤 (State Tracking) — 所有 kb-* 技能必守

本知識庫的目標規模是 1000+ 檔案的 codebase 與 10000+ 份文件。
任何一步都可能中斷，所以每一步的結果必須先落盤、再繼續。三條鐵律：

1. `先寫計畫`：開始處理前，先把完整 item 清單分批寫進
   `_state/runs/<run-id>/manifest.json`，之後不再修改 manifest
2. `每批落盤`：每處理完一批，立即更新 `progress.json` 與 `STATUS.md`，
   才能開始下一批。禁止「全部做完再一次寫入」
3. `續跑先讀狀態`：啟動任何 run 前，先檢查 `_state/runs/` 是否有同來源的
   `in-progress` run；有就從第一個未完成的 batch 繼續，不要重做

`run-id` 格式：`<yyyy-mm-dd>-<skill縮寫>-<slug>`，如 `2026-07-04-repo-service-a`。

### manifest.json（開跑時寫一次，之後唯讀）

```json
{
  "run": "2026-07-04-repo-service-a",
  "skill": "kb-ingest-repo",
  "source": "github.com/org/service-a",
  "created": "2026-07-04",
  "total_items": 1240,
  "batch_size": 50,
  "batches": [
    { "id": 1, "items": ["cmd/main.go", "model/order.go"] }
  ]
}
```

### progress.json（每批更新一次）

```json
{
  "run": "2026-07-04-repo-service-a",
  "status": "in-progress",
  "batches_done": 3,
  "batches_total": 25,
  "items_done": 150,
  "items_skipped": 4,
  "outputs": ["_inbox/2026-07-04-order-lifecycle.md"],
  "last_batch_finished": "2026-07-04"
}
```

- `status` 取值：`in-progress` | `done` | `failed` | `paused`
- 跳過的 item 必須在 `log.md` 寫一行原因（如 `SKIP vendor/lib.go — 排除清單`）

### STATUS.md（儀表板，每批同步更新）

```markdown
# KB Pipeline Status

| Run                        | Skill          | Progress | Status      | Updated    |
| -------------------------- | -------------- | -------- | ----------- | ---------- |
| 2026-07-04-repo-service-a  | kb-ingest-repo | 3/25     | in-progress | 2026-07-04 |
```

### 分批建議 (Batch Sizing)

| 來源             | 每批大小       |
| ---------------- | -------------- |
| repo 原始碼      | 50 檔          |
| git log 歷史     | 300 commits    |
| 文件 (web/檔案)  | 20 份          |
| 對話紀錄         | 200 則訊息     |
| schema 物件      | 30 表/topic    |
| distill captures | 20 個 captures |

## 真實性分級 (Truth Tiers)

| Tier           | 定義                               | 可進 curated 區 |
| -------------- | ---------------------------------- | --------------- |
| `confirmed`    | 人工確認過的事實                   | 是              |
| `firsthand`    | 第一人稱事實/經驗（使用者親述）    | 是              |
| `corroborated` | 2+ 獨立來源佐證（查 `_sources/`）  | 是              |
| `candidate`    | 單一來源、未確認                   | 否 — 留在 inbox |

`candidate` 升級路徑：新 capture 帶來第二個獨立來源 → `corroborated`；
使用者確認 → `confirmed`。降級：來源被推翻時整段移回 `_inbox/` 或標記 rejected。

## 指紋去重 (Fingerprint Dedup)

- `fingerprint = sha256(正規化文本)`：小寫化、壓縮空白、去標點後計算
- 入庫前先查重：`grep -r "sha256:<hash>" <proj>/_inbox <proj>/_sources`
- 命中且來源相同 → 跳過；命中但來源不同 → 不新建 capture，改在既有 capture 的
  `sources:` 追加來源（這是 corroboration 訊號）

## Capture 檔格式 (`_inbox/`)

```markdown
---
name: 2026-07-04-payment-flow
sources:
    - type: chat            # repo | history | web | chat | schema | file
      ref: "slack #payments 2026-07-03"
fingerprint: "sha256:ab12..."
captured: 2026-07-04
status: raw                 # raw | distilled | rejected
truth: candidate
zone-hint: payments
---

# 內容 (Content)

清理後全文或高訊號摘錄，保留原始語言。

# 候選事實 (Candidate Facts)

- 每條一句、可獨立驗證的事實陳述
- 附行內佐證位置（檔案:行、訊息時間戳、URL 段落）
```

## Source 檔格式 (`_sources/`)

```markdown
---
name: repo-service-a
type: repo
ref: github.com/org/service-a
reliability: high           # high | medium | low
last-seen: 2026-07-04
last-commit: "e3fb8ca"      # 僅 history 來源使用 — 增量游標
---

# repo-service-a

一句話描述此來源。

## Captures

- [[2026-07-04-payment-flow]]
```

## Entity 檔格式（curated 區）

與 `topology-builder` 完全相容，外加 truth 標註：

```markdown
---
name: service-a
type: service
zone: payments
tags: [billing]
aliases: [svc-a]
sources:
    - type: repo
      ref: github.com/org/service-a
---

# Service A

一句話定位。

## Billing Cycle

kind: concept
truth: corroborated

此維度做什麼（1~3 句）。

References:

- calls [[service-b#Method 2]] — 佐證：cmd/main.go:42

Sources: [[repo-service-a]], [[2026-07-04-payment-flow]]

## External Sources

- [API 文件](https://example.com/api)

## Backlinks

<!-- auto-generated: do not hand-edit -->
```

- `type` 取值：`service` | `module` | `datastore` | `external-api` | `article` |
  `channel` | `team` | `concept` | `decision` | `person`
- 維度標題下第一行 `kind:`（`concept` | `method` | `state` | `interface`），
  第二行 `truth:`（四種 tier 之一，`candidate` 禁止出現在此區）
- 每個維度結尾必有 `Sources:` 一行，wikilink 指向 `_sources/` 或 `_inbox/` 檔，
  或行內 URL；無來源的維度視為未佐證，`kb-verify` 會標記
- 維度數 2~12；`## External Sources` 與 `## Backlinks` 為固定章節，不計入

## 邊規則 (Edge Rules)

沿用 `topology-builder`，重點不變：

- 邊寫在維度 `References:` 清單：`- <relation> [[entity#Section]] — 佐證`
- 方向 = 發起者 → 接受者；反向關係由 Backlinks 重算，禁止手寫
- relation 動詞：`calls`, `uses`, `reads-from`, `writes-to`, `publishes-to`,
  `subscribes-to`, `depends-on`, `mentions`, `owned-by`, `supersedes`,
  `contradicts`；`mentions` 每維度 ≤ 2 條
- `supersedes` / `contradicts` 為 KB 新增：知識演進與矛盾標記，
  `contradicts` 邊存在時 `kb-verify` 必列入報告要求裁決
- 每條邊必須有發起方自身來源的直接佐證；指不出佐證即不建
- 基礎設施雜訊（logger、config、utils）不建邊

## `_index.md` 結構

專案層 `<proj>/_index.md`：

1. 註冊表：entity 清單（zone、type、維度數、最舊 truth tier）
2. Mermaid 總覽：zone 為 subgraph，邊聚合到 entity 層級
3. `## Frontier`：2-hop 外新實體、candidate 事實、待確認關係
4. `## Unlinked`：無入邊/無出邊清單（沿用 topology-builder 規則）

全域層 `<kb>/_index.md`（跨專案總覽，由入口/coordinator 在收尾更新）：

1. 專案註冊表：`| project | entities | edges | 最近 verify | 待裁決數 |`
2. `## Frontier (Global)`：跨專案缺口與未歸屬來源
