# gws (Google Workspace CLI)

來源: [googleworkspace/cli](https://github.com/googleworkspace/cli) (Apache-2.0, Rust) · [Skills Index](https://github.com/googleworkspace/cli/blob/main/docs/skills.md) · [Releases](https://github.com/googleworkspace/cli/releases)

一個涵蓋所有 Google Workspace API（Drive、Gmail、Calendar、Sheets、Docs、Chat、Tasks 等）的統一 CLI，兼顧人類操作與 AI Agent 呼叫。核心採用 Rust 實作、支援動態讀取 Google Discovery Service 建構指令樹（不寫死靜態指令表），並內建 95 個 Agent Skills（`SKILL.md`）。

> 注意：此專案明確聲明**並非 Google 官方支援產品**。

## 安裝

需先具備 Google Cloud Project（供 OAuth 使用）與 Google 帳號：

```sh
# npm（自動下載對應 OS / Arch 的預編譯 binary，需 Node.js 18+）
npm install -g @googleworkspace/cli

# Homebrew（macOS / Linux）
brew install googleworkspace-cli

# Cargo（從原始碼編譯）
cargo install --git https://github.com/googleworkspace/cli --locked

# Nix
nix run github:googleworkspace/cli
```

亦可直接至 [GitHub Releases](https://github.com/googleworkspace/cli/releases) 下載對應平台的預編譯 binary 放入 `$PATH`。

### 安裝 Agent Skills

```sh
# 一次安裝全部 skills
npx skills add https://github.com/googleworkspace/cli

# 只安裝特定 skill（目錄名為上游實際名稱）
npx skills add https://github.com/googleworkspace/cli/tree/main/skills/gws-gmail
npx skills add https://github.com/googleworkspace/cli/tree/main/skills/gws-drive
```

## 驗證與 Credential 機制

憑證預設以 AES-256-GCM 加密儲存於 OS keyring（設定 `GOOGLE_WORKSPACE_CLI_KEYRING_BACKEND=file` 時改存 `~/.config/gws/.encryption_key`）。

| 使用情境 | 設定方式 |
| :--- | :--- |
| 有 `gcloud` CLI（互動式桌面） | `gws auth setup`（建 Cloud project、開 API、登入）→ 後續 `gws auth login` |
| 無 `gcloud`（手動 OAuth） | 於 GCP Console 建立 **Desktop app** OAuth client，存為 `~/.config/gws/client_secret.json`，將自己加入 Test users 後執行 `gws auth login` |
| 測試模式 Scope 上限警告 | 未驗證 App 上限約 25 個 scopes，勿全選 `recommended`（85+ scopes），改指定服務：`gws auth login -s drive,gmail,sheets` |
| Headless / CI | 在有瀏覽器機器執行 `gws auth export --unmasked > credentials.json`，遠端設 `GOOGLE_WORKSPACE_CLI_CREDENTIALS_FILE` |
| Service Account | 設定 `export GOOGLE_WORKSPACE_CLI_CREDENTIALS_FILE=/path/to/service-account.json` |
| 既有的 Access Token | 設定 `export GOOGLE_WORKSPACE_CLI_TOKEN=$(gcloud auth print-access-token)` |

憑證讀取優先序：`GOOGLE_WORKSPACE_CLI_TOKEN` > `GOOGLE_WORKSPACE_CLI_CREDENTIALS_FILE` > `gws auth login` 加密憑證 > `~/.config/gws/credentials.json` 明文憑證。

## 核心架構與 CLI 用法

`gws` 採兩階段解析（Two-phase parsing）：讀取 `argv[1]` 識別服務 → 抓取並快取（24h）該服務的 Discovery Document → 動態建立指令樹 → 再解析剩餘參數並執行。所有成功輸出、錯誤與下載中繼資料皆為結構化 JSON。

```sh
# 動態探索任何方法的請求 / 回應 schema
gws schema drive.files.list
gws schema gmail.users.messages.list

# 呼叫資源方法（--params 放 query/path 參數，--json 放 request body）
gws drive files list --params '{"pageSize": 10}'
gws sheets spreadsheets create --json '{"properties": {"title": "Q1 Budget"}}'

# 預覽請求不實際執行（--dry-run）與自動分頁串流（--page-all）
gws chat spaces messages create --params '{"parent": "spaces/xyz"}' --json '{"text": "Deploy complete."}' --dry-run
gws drive files list --params '{"pageSize": 100}' --page-all | jq -r '.files[].name'
```

### Exit Codes

| Code | 意義 | 常見原因 |
| :--- | :--- | :--- |
| `0` | Success | 正常完成 |
| `1` | API error | Google API 回傳 4xx / 5xx（含 `accessNotConfigured` 未啟用 API） |
| `2` | Auth error | 憑證缺少、過期或無效 |
| `3` | Validation error | 參數錯誤、未知服務或無效 flag |
| `4` | Discovery error | 無法取得 API Discovery schema |
| `5` | Internal error | 非預期內部錯誤 |

## 內建 Skills 清單（統一以 `gws-` 前綴索引）

上游 `skills/` 目錄共包含 95 個 skills，分為 Services（19）、Helpers（25）、Personas（10）與 Recipes（41）；其中 Personas 與 Recipes 上游目錄名為 `persona-*` / `recipe-*`，下表統一加上 `gws-` 前綴標示並附上上游目錄名。

### 1. Services（19 個核心 API Skills）

| Skill (`gws-` 前綴) | 說明 |
| :--- | :--- |
| `gws-shared` | 驗證模式、全域 flags 與輸出格式共用規範（其餘 `gws-*` 前置必讀） |
| `gws-drive` | Google Drive：管理檔案、資料夾與共用雲端硬碟 |
| `gws-sheets` | Google Sheets：讀寫試算表 |
| `gws-gmail` | Gmail：寄信、讀信與信箱管理 |
| `gws-calendar` | Google Calendar：管理日曆與活動 |
| `gws-admin-reports` | Google Workspace Admin SDK：稽核紀錄與用量報表 |
| `gws-docs` | Google Docs：讀寫文件 |
| `gws-slides` | Google Slides：讀寫簡報 |
| `gws-tasks` | Google Tasks：管理任務清單與待辦事項 |
| `gws-people` | Google People：管理聯絡人與個人檔案 |
| `gws-chat` | Google Chat：管理聊天室與訊息 |
| `gws-classroom` | Google Classroom：管理課程、名冊與作業 |
| `gws-forms` | Google Forms：讀寫表單 |
| `gws-keep` | Google Keep：管理筆記 |
| `gws-meet` | Google Meet：管理會議 |
| `gws-events` | Google Workspace Events：訂閱 Workspace 事件 |
| `gws-modelarmor` | Google Model Armor：過濾使用者輸入與模型輸出的安全檢查 |
| `gws-workflow` | Google Workflow：跨服務生產力工作流 |
| `gws-script` | Google Apps Script：管理指令碼專案 |

### 2. Helpers（25 個快捷指令 Skills）

| Skill (`gws-` 前綴) | 說明 |
| :--- | :--- |
| `gws-drive-upload` | 上傳檔案至 Drive 並自動帶入 metadata |
| `gws-sheets-append` | 新增一列資料至 Google Sheets |
| `gws-sheets-read` | 從 Google Sheets 讀取儲存格數值 |
| `gws-gmail-send` | 寄送 Gmail 郵件 |
| `gws-gmail-triage` | 顯示未讀收件匣摘要（寄件者、主旨、日期） |
| `gws-gmail-reply` | 回覆郵件（自動處理 thread） |
| `gws-gmail-reply-all` | 回覆全部（自動處理 thread） |
| `gws-gmail-forward` | 轉寄郵件給新收件者 |
| `gws-gmail-read` | 讀取郵件並抽取純文字/HTML 內文或 headers |
| `gws-gmail-watch` | 監聽新郵件並以 NDJSON 串流輸出 |
| `gws-calendar-insert` | 建立新日曆活動 |
| `gws-calendar-agenda` | 顯示跨日曆的近期活動 |
| `gws-docs-write` | 附加文字至 Google Docs |
| `gws-chat-send` | 傳送訊息至 Google Chat 空間 |
| `gws-events-subscribe` | 訂閱 Workspace 事件並串流為 NDJSON |
| `gws-events-renew` | 續期或重新啟用 Workspace Events 訂閱 |
| `gws-modelarmor-sanitize-prompt` | 透過 Model Armor template 淨化使用者 prompt |
| `gws-modelarmor-sanitize-response` | 透過 Model Armor template 淨化模型回應 |
| `gws-modelarmor-create-template` | 建立新的 Model Armor template |
| `gws-workflow-standup-report` | 彙整今日會議與未完成任務為站會摘要 |
| `gws-workflow-meeting-prep` | 準備下一場會議：議程、參與者與關聯文件 |
| `gws-workflow-email-to-task` | 將 Gmail 郵件轉換為 Google Tasks 待辦 |
| `gws-workflow-weekly-digest` | 每週摘要：本週會議與未讀郵件統計 |
| `gws-workflow-file-announce` | 在 Chat 空間公告一個 Drive 檔案 |
| `gws-script-push` | 上傳本機檔案至 Apps Script 專案 |

### 3. Personas（10 個角色型工作包，上游目錄為 `skills/persona-*`）

| Skill (`gws-` 前綴) | 上游目錄名 | 說明 |
| :--- | :--- | :--- |
| `gws-persona-exec-assistant` | `persona-exec-assistant` | 管理主管行程、收件匣與通訊 |
| `gws-persona-project-manager` | `persona-project-manager` | 專案協調：追蹤任務、排會議與共享文件 |
| `gws-persona-hr-coordinator` | `persona-hr-coordinator` | 人資流程：新人報到、公告與員工通訊 |
| `gws-persona-sales-ops` | `persona-sales-ops` | 銷售營運：追蹤商機、排通話與客戶聯繫 |
| `gws-persona-it-admin` | `persona-it-admin` | IT 管理：監控安全性與設定 Workspace |
| `gws-persona-content-creator` | `persona-content-creator` | 跨 Workspace 建立、整理與分發內容 |
| `gws-persona-customer-support` | `persona-customer-support` | 客服支援：追蹤工單、回覆與升級處理 |
| `gws-persona-event-coordinator` | `persona-event-coordinator` | 活動規劃與執行：排程、邀請與後勤 |
| `gws-persona-team-lead` | `persona-team-lead` | 團隊帶領：主持 standup、協調任務與溝通 |
| `gws-persona-researcher` | `persona-researcher` | 研究整理：管理參考資料、筆記與協作 |

### 4. Recipes（41 個多步驟任務腳本，上游目錄為 `skills/recipe-*`）

| Skill (`gws-` 前綴) | 上游目錄名 | 說明 |
| :--- | :--- | :--- |
| `gws-recipe-label-and-archive-emails` | `recipe-label-and-archive-emails` | 為符合條件的 Gmail 加上標籤並封存 |
| `gws-recipe-draft-email-from-doc` | `recipe-draft-email-from-doc` | 讀取 Google Doc 內容作為 Gmail 郵件內文 |
| `gws-recipe-organize-drive-folder` | `recipe-organize-drive-folder` | 建立 Drive 資料夾結構並搬移檔案 |
| `gws-recipe-share-folder-with-team` | `recipe-share-folder-with-team` | 將 Drive 資料夾與內容分享給協作者清單 |
| `gws-recipe-email-drive-link` | `recipe-email-drive-link` | 分享 Drive 檔案並寄信通知收件人連結 |
| `gws-recipe-create-doc-from-template` | `recipe-create-doc-from-template` | 複製 Docs 範本、填入內容並分享 |
| `gws-recipe-create-expense-tracker` | `recipe-create-expense-tracker` | 建立含欄位標題與初始資料的費用追蹤試算表 |
| `gws-recipe-copy-sheet-for-new-month` | `recipe-copy-sheet-for-new-month` | 複製 Sheets 範本分頁供新月份追蹤使用 |
| `gws-recipe-block-focus-time` | `recipe-block-focus-time` | 在 Calendar 建立週期性專注時段 |
| `gws-recipe-reschedule-meeting` | `recipe-reschedule-meeting` | 調整 Calendar 會議時間並自動通知參與者 |
| `gws-recipe-create-gmail-filter` | `recipe-create-gmail-filter` | 建立 Gmail filter 自動加標籤、加星號或分類 |
| `gws-recipe-schedule-recurring-event` | `recipe-schedule-recurring-event` | 建立含參與者的週期性日曆活動 |
| `gws-recipe-find-free-time` | `recipe-find-free-time` | 查詢多位使用者的 free/busy 狀態以找會議空檔 |
| `gws-recipe-bulk-download-folder` | `recipe-bulk-download-folder` | 列出並批次下載 Drive 資料夾內所有檔案 |
| `gws-recipe-find-large-files` | `recipe-find-large-files` | 找出佔用配額的大型 Drive 檔案 |
| `gws-recipe-create-shared-drive` | `recipe-create-shared-drive` | 建立共用雲端硬碟並加入不同角色成員 |
| `gws-recipe-log-deal-update` | `recipe-log-deal-update` | 將銷售商機狀態更新附加至 Sheets 追蹤表 |
| `gws-recipe-collect-form-responses` | `recipe-collect-form-responses` | 取得並檢視 Google Form 回覆內容 |
| `gws-recipe-post-mortem-setup` | `recipe-post-mortem-setup` | 建立 Docs 復盤文件、排定 Calendar 檢討會並於 Chat 通知 |
| `gws-recipe-create-task-list` | `recipe-create-task-list` | 建立新的 Tasks 清單與初始待辦事項 |
| `gws-recipe-review-overdue-tasks` | `recipe-review-overdue-tasks` | 找出已逾期的 Google Tasks |
| `gws-recipe-watch-drive-changes` | `recipe-watch-drive-changes` | 訂閱 Drive 檔案或資料夾變更通知 |
| `gws-recipe-create-classroom-course` | `recipe-create-classroom-course` | 建立 Classroom 課程並邀請學生 |
| `gws-recipe-create-meet-space` | `recipe-create-meet-space` | 建立 Google Meet 會議空間並分享加入連結 |
| `gws-recipe-review-meet-participants` | `recipe-review-meet-participants` | 檢視 Google Meet 與會者名單與參與時長 |
| `gws-recipe-create-presentation` | `recipe-create-presentation` | 建立新 Slides 簡報並新增初始投影片 |
| `gws-recipe-save-email-attachments` | `recipe-save-email-attachments` | 搜尋含附件的 Gmail 並存至 Drive 資料夾 |
| `gws-recipe-send-team-announcement` | `recipe-send-team-announcement` | 同時透過 Gmail 與 Chat 空間發送團隊公告 |
| `gws-recipe-create-feedback-form` | `recipe-create-feedback-form` | 建立意見回饋表單並以 Gmail 分享 |
| `gws-recipe-sync-contacts-to-sheet` | `recipe-sync-contacts-to-sheet` | 將 Google Contacts 匯出至 Sheets |
| `gws-recipe-share-event-materials` | `recipe-share-event-materials` | 將 Drive 檔案分享給 Calendar 活動所有參與者 |
| `gws-recipe-create-vacation-responder` | `recipe-create-vacation-responder` | 啟用指定日期與訊息的 Gmail 休假自動回覆 |
| `gws-recipe-create-events-from-sheet` | `recipe-create-events-from-sheet` | 從 Sheets 讀取活動資料並逐列建立 Calendar 活動 |
| `gws-recipe-plan-weekly-schedule` | `recipe-plan-weekly-schedule` | 檢視本週 Calendar、找出空檔並補入活動 |
| `gws-recipe-share-doc-and-notify` | `recipe-share-doc-and-notify` | 分享 Docs 編輯權限並寄信通知協作者連結 |
| `gws-recipe-backup-sheet-as-csv` | `recipe-backup-sheet-as-csv` | 將 Sheets 匯出為本機 CSV 備份 |
| `gws-recipe-save-email-to-doc` | `recipe-save-email-to-doc` | 將 Gmail 郵件內文存成 Google Doc 歸檔 |
| `gws-recipe-compare-sheet-tabs` | `recipe-compare-sheet-tabs` | 比對同份試算表內兩個分頁的資料差異 |
| `gws-recipe-batch-invite-to-event` | `recipe-batch-invite-to-event` | 批次將與會者加入既有 Calendar 活動並寄出通知 |
| `gws-recipe-forward-labeled-emails` | `recipe-forward-labeled-emails` | 搜尋特定標籤的 Gmail 郵件並轉寄至目標信箱 |
| `gws-recipe-generate-report-from-sheet` | `recipe-generate-report-from-sheet` | 從 Sheets 讀取數據並產出格式化的 Docs 報告 |
