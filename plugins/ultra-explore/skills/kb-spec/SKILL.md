---
name: kb-spec
description: Reference spec read by other kb-* skills; not directly invocable.
version: "1.1.0"
user-invocable: false
disable-model-invocation: true
metadata:
    type: reference
---

# kb-spec — 知識庫規範 (Knowledge Base Spec)

所有 `kb-*` 技能共用的單一事實來源：儲存佈局、檔案格式、真實性分級、去重與邊規則。

| 要找什麼 | 讀哪裡 |
| --- | --- |
| capture / source / entity 檔的 frontmatter 與章節 | [references/file-formats.md](references/file-formats.md) |
| 邊規則、relation 動詞、`_index.md` 結構 | [references/file-formats.md](references/file-formats.md) |
| `manifest.json` / `progress.json` / `STATUS.md` 樣板、分批大小 | [references/state-tracking.md](references/state-tracking.md) |
| 佈局、命名、truth tier、指紋去重 | 本檔以下各節 |

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
        ├── _raw/              # 工具產物：kb_history.py commit 清單 (commits.jsonl)
        ├── _diffs/            # 工具產物：每 ISO 週過濾後 diff
        ├── stats.json         # 工具產物：週分桶與作者統計
        ├── CHANGELOG.md       # 工具產物：週敘事（kb-ingest-history 回填）
        └── payments/          # zone 一律用資料夾
            └── service-a.md   # curated entity
```

工具產物（`_raw/`、`_diffs/`、`stats.json`、`CHANGELOG.md`）由
`kb-ingest-history` 內建的 `kb_history.py` 管道產生、供其佐證與回填 —
它們不是 entity、不是 capture，entity 掃描與驗證一律排除。

### 命名規則 (Naming)

- entity 檔名 `<entity_name>.md`，kebab-case，同一 `<proj>` 內跨 zone 全域唯一；
  wikilink 只在 `<proj>` 內解析，不跨專案
- capture 檔名 `<yyyy-mm-dd>-<slug>.md`
- `run-id` 一律含專案名（`<date>-<skill縮寫>-<project>`），全域 `_state/`
  才分得開各專案的 run
- web/chat/schema 來源歸屬：使用者指定專案，或依 `zone-hint` 對應；
  無法歸屬時建立獨立專案資料夾
- `_inbox/`、`_sources/`、`_state/` 內的檔案不是 entity，不參與 wikilink 邊統計

## 狀態追蹤三鐵律 (State Tracking) — 所有 kb-* 技能必守

本知識庫的目標規模是 1000+ 檔案的 codebase 與 10000+ 份文件。
任何一步都可能中斷，所以每一步的結果必須先落盤、再繼續。

1. `先寫計畫`：開始處理前，先把完整 item 清單分批寫進
   `_state/runs/<run-id>/manifest.json`，之後不再修改 manifest
2. `每批落盤`：每處理完一批，立即更新 `progress.json` 與 `STATUS.md`，
   才能開始下一批。禁止「全部做完再一次寫入」
3. `續跑先讀狀態`：啟動任何 run 前，先檢查 `_state/runs/` 是否有同來源的
   `in-progress` run；有就從第一個未完成的 batch 繼續，不要重做

三個狀態檔的完整樣板與各來源的分批大小見
[references/state-tracking.md](references/state-tracking.md)。

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
