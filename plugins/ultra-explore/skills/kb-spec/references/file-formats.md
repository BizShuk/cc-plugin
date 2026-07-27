# kb-spec — 檔案格式與邊規則 (File Formats & Edge Rules)

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
