# general — 通用插件 (General Plugin)

通用、跨專案的 Claude Code 插件。整合：

- 5 個本地跨場景技能（插件 metadata、每日彙整、文件檢查、待辦排序、專案任務）
- 1 個外部技能（`bizshuk/mdserver`）
- 1 個 `feature` 功能實作代理
- 1 個 always-on hook（agent loop 結束時的終端機鈴聲）
- 1 份 output-style 樣板（`brief`）

源自兩個先前的插件合併：

- 原 `general`（規劃、彙整類技能）
- 原 `base`（hooks、output-styles、`markdownlint` 技能）

## 技能 (Skills)

| Skill | 用途 |
| --- | --- |
| `claude-plugin-metadata` | 建立與更新 Claude 插件與技能的元資料（含 plugin.json 與 marketplace.json） |
| `daily-summary` | 彙整過去 24h 跨來源工作，產生工作日報並寫入 Apple Notes |
| `markdownlint` | Markdown 格式檢查（精選 rule + CUSTOM-01 no-bold）；`範圍限 plugins/experiment/skills/` 下的 `.md` 檔 |
| `mdserver` | 從外部來源 `bizshuk/mdserver` 安裝，用於 preview、serve 與本機瀏覽 Markdown 目錄 |
| `sort-todo` | 排序並格式化 `.todo` 待辦清單（P0/P1/P2 三鏡頭評分 + 依目標 repo 的 feature 目錄推導章節；`專案無關 (project-agnostic)`） |

## 代理 (Agents)

| Agent | 用途 |
| :--- | :--- |
| `feature` (`./agents/feature.md`) | 規劃、實作並驗證端到端功能 |

## Hooks

| Event          | 行為                  | 目的                                       |
| -------------- | --------------------- | ------------------------------------------ |
| `Stop`         | `printf '\a'` → stderr | 正常完成時通知                              |
| `StopFailure`  | `printf '\a'` → stderr | 失敗/中斷時通知（語意與 `Stop` 相同）        |

兩個事件共用同一支 `stop-bell.sh`，行為無差異。

## Output Styles

本資料夾只保留`一份樣板 (single template)`，作為撰寫自訂 `outputStyle` 的起點：

| 檔案       | 風格 | 長度    | 適合情境                     |
| ---------- | ---- | ------- | ---------------------------- |
| `brief.md` | 精簡 | 1–10 行 | 日常查詢、簡單任務、結論先行 |

實際啟用的風格由 `config/settings.json` 的 `outputStyle` 指定，風格檔本體放在
`config/output-styles/`。寫法與放置位置詳見 `output-styles/README.md`。

## 檔案結構

```text
plugins/general/
├── .claude-plugin/
│   └── plugin.json      # manifest
├── .lsp.json            # marksman LSP 設定（Markdown）
├── agents/
│   └── feature.md       # feature agent
├── hooks/
│   ├── hooks.json       # Stop / StopFailure 註冊
│   └── stop-bell.sh     # 終端機 bell 實作
├── output-styles/
│   ├── brief.md         # 唯一樣板
│   └── README.md
├── skills/              # 5 個技能目錄
└── README.md
```

## daily-summary（細節）

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
bash skills/daily-summary/scripts/scan.sh 24   # 參數 = 視窗小時數
```

觸發詞："daily summary"、"工作日報"、"summarize my work today"、"今天做了什麼"。

注意：此技能讀取本機路徑，僅適用於本機執行（cloud routine 讀不到這些來源）。

## 疑難排解

| 現象                  | 原因                              | 解法                                     |
| --------------------- | --------------------------------- | ---------------------------------------- |
| 沒聽到鈴聲            | 終端機關閉了 bell (`visible-bell`) | 終端機偏好設定開啟 audible bell           |
| 終端機顯示亂碼        | 不支援 BEL                        | 確認 `$TERM` 有效；腳本會 fallback 到 `tput bel` |
| Hook 沒執行           | Plugin 未啟用                      | 檢查 `enabledPlugins` 設定               |

## 設計原則

- **零依賴**：bell 僅使用 `bash` 內建 + 可選 `tput`
- **冪等**：每次觸發只輸出一個 BEL
- **靜默失敗**：若 `tput` 不可用，fallback 到 `printf '\a'`，不輸出錯誤
- **無 stdout 污染**：bell 寫到 stderr，避免干擾 hook 輸出 JSON 解析
