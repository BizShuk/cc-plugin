---
name: ultra-explore
description: >
    Single manual entry point for the ultra-explore plugin — runs the full
    knowledge base pipeline (ingest repo/history/web/chat/schema → distill →
    connect → verify) over the given sources, resumable at every step.
    Manual invocation only: run /ultra-explore with a list of sources, or say
    "ultra explore", "全量建庫". The model must never trigger this on its own.
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit, Agent, AskUserQuestion
user-invocable: true
disable-model-invocation: true
effort: high
---

# ultra-explore — 全管道入口

唯一的全管道手動入口。使用者用 `/ultra-explore <sources>` 觸發後，
本技能負責：解析來源 → 依序跑五階段 → 回報。所有格式與狀態規則以
`kb-spec` 為準（`<plugin>/skills/kb-spec/SKILL.md`），開工先讀它。

子技能（kb-*）全部 `disable-model-invocation: true` — 它們不會被模型
自主觸發，只由本入口（或使用者手動）驅動。執行方式：Read 對應
SKILL.md，照其 Checklist 逐步執行。

## 輸出位置 (Output Routing)

依 `kb-spec` 兩層佈局，硬性規定：

- 全域狀態與跨專案總覽 → `~/projects/product/`（`_state/`、`_index.md`）
- 每個 repo/專案的知識 → `~/projects/product/projects/<project>/`
  （`<project>` = repo 資料夾名，kebab-case）

## 輸入 (Input)

`/ultra-explore` 後接來源清單，型態自動判別：

| 輸入樣式                          | 判定       | 交給                                  |
| --------------------------------- | ---------- | ------------------------------------- |
| 資料夾（含多個 git repo）         | monofolder | 展開為多個 repo（見下節）             |
| 本地路徑含 `.git` / git URL       | repo       | `kb-ingest-repo` + `kb-ingest-history`|
| `http(s)://` URL / URL 清單檔     | web        | `kb-ingest-web`                       |
| 對話匯出檔 / 貼上對話 / 本次對話  | chat       | `kb-ingest-chat`                      |
| DDL/migration/broker 設定/DSN     | schema     | `kb-ingest-schema`                    |

- repo 來源自動同時排入 history（讀現狀 + 讀演進），除非使用者排除
- 無參數時：先讀 `<kb>/_state/STATUS.md`，有 in-progress run → 詢問續跑；
  全部 done → 走增量模式（見下）；當前資料夾含多個 git repo →
  monofolder 模式；都不是 → 用 AskUserQuestion 要來源清單

## Monofolder 模式 (Monofolder Mode)

單一觸發掃遍資料夾下所有 git repo，一 repo 一專案：

1. 探索 repo（跳過被上層 repo `.gitignore` 的資料目錄與 submodule）：

    ```bash
    find <folder> -maxdepth 4 -type d -name .git 2>/dev/null \
      | sed 's|/\.git$||' | while read -r r; do
        p=$(git -C "$(dirname "$r")" rev-parse --show-toplevel 2>/dev/null)
        if [ -n "$p" ] && [ "$p" != "$(cd "$r" && pwd)" ]; then
            git -C "$p" check-ignore -q "$r" && continue     # 被忽略的資料/vendored
            [ -f "$p/.gitmodules" ] && grep -q "$(basename "$r")" "$p/.gitmodules" \
              && continue                                     # submodule 歸屬上層
        fi
        echo "$r"
    done
    ```

2. 把探索結果清單先落盤到 `<kb>/_state/runs/<date>-mono-<folder>/manifest.json`
   （一 repo 一 item），使用者確認清單後才開跑（repo 數可能很多）
3. 每個 repo = 一個 `<project>`：`kb-ingest-repo` + `kb-ingest-history`
   成對排入，輸出到 `<kb>/projects/<project>/`；repo 內部檔案掃描一律走
   `git ls-files`（天然遵循該 repo 自己的 `.gitignore`，資料檔不入庫）
4. mono run 的 progress 以 repo 為單位計數：每完成一個 repo 的兩個
   ingest run，`progress.json` 的 `items_done` +1
5. Phase 2~4（distill/connect/verify）逐專案執行；Phase 5 更新全域
   `<kb>/_index.md` 專案註冊表

## 執行 (Execution)

### 派工模式（環境支援 Agent tool 時優先）

spawn `kb-coordinator` agent，把來源清單與 kb 根目錄交給它；
一個來源一個子代理平行 ingest，coordinator 依 phase gate 推進。

### 單機模式（fallback）

依序自己執行，每階段 Read 對應 SKILL.md 照做：

```md
- [ ] Phase 1 Ingest: 逐來源跑對應 kb-ingest-* 技能（含 history）
- [ ] Phase gate: grep -L '"status": "done"' <kb>/_state/runs/*/progress.json → 空才前進
- [ ] Phase 2 Distill: kb-distill（_inbox → curated entities）
- [ ] Phase 3 Connect: kb-connect（建邊 + Backlinks + _index.md）
- [ ] Phase 4 Verify: kb-verify（健檢報告落盤 _state/verify/）
- [ ] Phase 5 Report: 總結表 + 待裁決清單
```

鐵律（詳見 kb-spec 狀態追蹤）：每批落盤、續跑先讀狀態、
phase gate 未過不前進。一次 session 跑不完是正常的 —
把 `STATUS.md` 更新到位，下次 `/ultra-explore` 會自動接續。

## 增量模式 (Incremental)

kb 已存在且無 in-progress run 時：

1. history 來源：游標增量（只讀新 commit）
2. `_sources/` 中 `last-seen` 過期（>90 天）的來源重新 ingest
3. distill 只處理 `status: raw` captures
4. connect 只對新/變更 entities 建邊；Backlinks 與 `_index.md` 全圖重算
5. verify 照常全跑

## 回報 (Report)

```markdown
# Ultra-Explore Report — <date>

| Phase   | Runs | Done | Failed | 產出                   |
| ------- | ---- | ---- | ------ | ---------------------- |
| Ingest  | 5    | 4    | 1      | 214 captures           |
| Distill | 1    | 1    | 0      | 38 entities (5 zones)  |
| Connect | 1    | 1    | 0      | 156 edges              |
| Verify  | 1    | 1    | 0      | 2 conflicts, 3 orphans |

## 待裁決 (Pending Decisions)

## 下一步 (Next)
```

## Common Mistakes

| 錯誤                               | 修正                                         |
| ---------------------------------- | -------------------------------------------- |
| 模型在一般任務中自行觸發本技能     | 僅限使用者手動 `/ultra-explore`              |
| 跳過 phase gate 直接 distill       | 所有 ingest run 皆 done 才前進               |
| 忽略既有 in-progress run 重頭開跑  | Step 0 讀 STATUS.md，預設續跑                |
| repo 來源漏掉 history              | repo 一律同時排入 kb-ingest-history          |
| 只在對話回報、狀態不落盤           | STATUS.md 與 progress.json 是唯一事實來源    |
| monofolder 把 vendored/資料 repo 入庫 | 上層 `git check-ignore` 命中即跳過        |
| monofolder 未經確認直接掃數十 repo | repo 清單先落盤並經使用者確認                |
| 知識輸出寫到全域根                 | 知識一律在 projects/<project>/；全域只放狀態 |
