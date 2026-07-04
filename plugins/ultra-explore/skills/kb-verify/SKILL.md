---
name: kb-verify
description: >
    Use when verifying knowledge base integrity — broken wikilinks, orphans,
    missing kind/truth labels, candidate facts leaked into the curated zone,
    unsourced dimensions, edge grounding audits, contradiction and staleness
    reports. Persists a dated verification report for tracking. Triggers on:
    "verify the knowledge base", "kb health check", "驗證知識庫",
    "知識庫健檢", "kb verify".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: high
context: fork
---

# kb-verify — 知識庫健檢

對 curated 區與狀態層做完整性驗證，產出落盤的健檢報告。
規模再小也必跑全部檢查；規模大時腳本照樣一次跑完（grep 級成本）。
格式規則依 `kb-spec`。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 1: 跑結構檢查腳本（A1~A4）
- [ ] Step 2: 跑品質檢查（B1~B4，抽查類）
- [ ] Step 3: 跑狀態層檢查（C1~C3）
- [ ] Step 4: 寫報告到 <kb>/_state/verify/<date>-<project>.md + 更新 STATUS.md
- [ ] Step 5: 回報結論與待裁決清單
```

驗證範圍以「專案」為單位：使用者指定 `<project>` 就只驗該庫；未指定則
`ls ~/projects/product/projects/` 逐專案跑一輪（wikilink 不跨專案，
逐專案驗才不會誤報跨庫斷鏈）。以下 `$root` 為單一專案根 `<proj>`；
curated 檔案 = 排除 `_inbox/`、`_sources/`、`_index.md` 後的 `*.md`。

## Step 1 — 結構檢查 (Structural, 全量腳本)

```bash
root=~/projects/product/projects/<project>
curated() { find "$root" -name '*.md' ! -path '*/_inbox/*' \
  ! -path '*/_sources/*' ! -name '_index.md' ! -name 'CHANGELOG.md'; }

# A1 檔名全域唯一
curated | xargs -n1 basename | sort | uniq -d

# A2 wikilink 斷鏈：目標檔或標題不存在（含 _inbox/_sources 目標）
grep -rhoE '\[\[[^]]+\]\]' "$root" --include='*.md' | sort -u |
    sed 's/^\[\[//;s/\]\]$//' | while IFS='#' read -r name sec; do
    [ -z "$name" ] && continue
    f=$(find "$root" -name "$name.md" | head -1)
    [ -z "$f" ] && { echo "missing target: $name"; continue; }
    [ -n "$sec" ] && ! grep -qxF "## $sec" "$f" && echo "missing heading: $name#$sec"
done

# A3 維度標註：標題下第一行 kind、第二行 truth，取值合法
curated | while read -r f; do
    awk -v F="$f" '
        /^## /{h=substr($0,4); want=(h!="External Sources" && h!="Backlinks"); n=0; next}
        want && NF { n++
            if (n==1 && $0 !~ /^kind: (concept|method|state|interface)$/) print F": "h" (bad kind)"
            if (n==2 && $0 !~ /^truth: (confirmed|firsthand|corroborated)$/) print F": "h" (bad truth)"
            if (n>=2) want=0 }
    ' "$f"
done
# 注意：curated 區 truth 出現 candidate 即違規（A3 會列為 bad truth）

# A4 維度必有 Sources: 行
curated | while read -r f; do
    awk -v F="$f" '
        /^## /{ if (want && !src) print F": "h" (no Sources)"
            h=substr($0,4); want=(h!="External Sources" && h!="Backlinks"); src=0; next }
        want && /^Sources:/ {src=1}
        END { if (want && !src) print F": "h" (no Sources)" }
    ' "$f"
done
```

## Step 2 — 品質檢查 (Quality, 抽查)

- `B1 邊佐證抽查`：隨機抽 `max(10, 邊總數的 5%)` 條正向邊，逐條確認補充說明
  指得出佐證位置（`檔案:行`、表名、時間戳、引句）。指不出 → 列入報告：
  刪除或降級為維度內文
- `B2 Backlinks 一致性`：抽 10 個 entity，其 Backlinks 每條都要有對應正向邊；
  無 → 列入報告（kb-connect 重跑 Step 3）
- `B3 矛盾清單`：`grep -rn 'contradicts \[\[' $root` 全列出 + `_index.md`
  Frontier 的 CONFLICT 條目 — 這些必須進報告的待裁決區，禁止靜默留存
- `B4 樞紐噪音`：入邊+出邊合計異常高且多為 `mentions` 的 entity
  → 噪音候選，報告供人工裁決

## Step 3 — 狀態層檢查 (State Layer)

- `C1 殭屍 run`：`_state/runs/` 內 `status: in-progress` 且
  `last_batch_finished` 超過 7 天 → 列出，建議續跑或標 `paused`
- `C2 capture 一致性`：`status: distilled` 的 capture 應被至少一個 entity 的
  `Sources:` 引用；`status: raw` 積壓超過 100 個 → 建議跑 `kb-distill`
- `C3 來源過期`：`_sources/` 中 `last-seen` 超過 90 天的 repo/schema 來源
  → 列入 staleness 清單，建議重新 ingest 校正

## Step 4 — 落盤報告 (Persist Report)

寫入 `_state/verify/<yyyy-mm-dd>-<project>.md`：

```markdown
# KB 健檢報告 (Verification Report) — <date>

## 摘要 (Summary)

| 檢查        | 結果        | 數量 |
| ----------- | ----------- | ---- |
| A1 重名     | pass / fail | 0    |
| A2 斷鏈     | ...         |      |

## 待裁決 (Pending Decisions)

- CONFLICT: ...（來源、兩方主張、建議）

## 待修復 (Fix List)

- [ ] ...（可直接交給對應 kb-* 技能的具體修復項）
```

同步更新 `STATUS.md` 加一列 verify run。與上一份報告比較：
新增/解決的問題數寫進摘要（趨勢追蹤）。

## Step 5 — 回報 (Report)

結論先行：pass/fail 總覽、待裁決清單、修復建議（各交給哪個 kb-* 技能）。
發現可自動修復項（斷鏈因標題改名、Backlinks 不一致）時列出但不擅自修改 —
本技能唯讀 curated 區，只寫報告。

## Common Mistakes

| 錯誤                         | 修正                                     |
| ---------------------------- | ---------------------------------------- |
| 「圖很小」跳過部分檢查       | A1~C3 全跑，成本是 grep 級               |
| 邊佐證抽查只抽 2~3 條        | max(10, 5%) 為下限                       |
| 發現問題直接動手改 entity    | 唯讀；報告列修復項交對應技能             |
| 矛盾條目靜默留存             | B3 全列入待裁決區                        |
| 報告只印在對話不落盤         | 必寫 <kb>/_state/verify/<date>-<project>.md 供趨勢追蹤  |
