# general — 通用功能插件 (General Plugin)

通用、跨專案的 Claude Code 技能集合，涵蓋專案探索、知識圖譜、文件品質與每日工作彙整。

## 技能 (Skills)

| Skill | 用途 |
| --- | --- |
| `business-planner` | 一次規劃一個 feature 的商業價值：挖掘未開發新價值或解鎖隱藏價值，輸出 `plans/business-<feature_name>.md` |
| `daily-summary` | 彙整過去 24h 跨來源工作，產生工作日報並寫入 Apple Notes |
| `sort-todo` | 排序並格式化待辦清單 |
| `system-planner` | 一次規劃一個 feature 的系統架構：位置、邊界、介面與資料流，保持整體清晰且可擴充，輸出 `plans/architecture-<feature_name>.md` |

## daily-summary

彙整五個本機來源（過去 24h），合成「工作摘要 + 明確 TODO + 潛在 TODO」，寫入 Apple Notes `iCloud › Daily`，標題 `Daily Summary - {YYYY-MM-DD}`（同日存在則 append）。

來源與路徑：

```tree
1 claude-mem       ~/.claude-mem/claude-mem.db              (主來源)
2 claude sessions  ~/.claude/projects/<encoded>/*.jsonl     (raw, 通常略過)
3 antigravity      ~/.gemini/antigravity/brain/*/task.md    (TODO [x]/[ ])
                   ~/.gemini/antigravity/conversations/*.db (binary, 活動訊號)
4 hermes           ~/.hermes/state.db  (sessions+messages, channel: slack/whatsapp/cron/cli...)
5 git log          ~/projects/**/.git  --since=24h
→ 輸出            notes CLI → iCloud / Daily
```

使用方式：

```bash
# 技能開頭先跑預檢掃描（read-only），只處理 [OK] 來源
bash skills/daily-summary/scan.sh 24   # 參數 = 視窗小時數
```

觸發詞："daily summary"、"工作日報"、"summarize my work today"、"今天做了什麼"。

注意：此技能讀取本機路徑，僅適用於本機執行（cloud routine 讀不到這些來源）。
