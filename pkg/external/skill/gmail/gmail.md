# gmail

來源：`~/.gemini/config/plugins/gmail/`（Plugin `gmail` v0.0.7，包含 `skills/gmail/SKILL.md`、`scripts/gmail`、`scripts/gmail.cmd`）

透過 `gmail` Go CLI（`gmail_go version cl/1017793813`）直接操作 Gmail REST API，支援搜尋、閱讀郵件與討論串、下載附件、草稿、標籤、Emoji 反應與度假自動回覆查詢。

## 如何下載與初始化 `gmail` CLI

### 1. 透過 Launcher 初始化（標準做法）

首次使用前執行一次 launcher 腳本，會自動下載對應作業系統與架構的 pinned binary，並註冊 `gmail` 到 agent 的 `PATH`：

```bash
# macOS / Linux
~/.gemini/config/plugins/gmail/skills/gmail/scripts/gmail --version
```

```powershell
# Windows (PowerShell)
& "$env:USERPROFILE\.gemini\config\plugins\gmail\skills\gmail\scripts\gmail.cmd" --version
```

完成註冊後，後續一律直接呼叫 `gmail`（不可改用腳本路徑呼叫，否則無法取得 Antigravity 的使用者登入狀態）。

### 2. 下載來源與本機檔案配置

`scripts/gmail` 與 `scripts/gmail.cmd` 鎖定版本 `CLI_VERSION=1017793813`，從 Google Cloud Storage 下載：

```text
https://storage.googleapis.com/antigravity-public/plugin/cli/1017793813/<os>_<arch>/gmail
```

支援平台：

| OS | Arch | 下載路徑片段 |
| :--- | :--- | :--- |
| macOS (`Darwin`) | `arm64` / `x86_64` | `darwin_arm64/gmail`、`darwin_x86_64/gmail` |
| Linux (`Linux`) | `arm64` / `x86_64` | `linux_arm64/gmail`、`linux_x86_64/gmail` |
| Windows | `arm64` / `x86_64` | `windows_arm64/gmail`、`windows_x86_64/gmail` |

手動下載 binary（例如 macOS Apple Silicon）：

```bash
curl -fsSL https://storage.googleapis.com/antigravity-public/plugin/cli/1017793813/darwin_arm64/gmail -o gmail
chmod +x gmail
```

執行期檔案落點（以 `ANTIGRAVITY_APP_DATA_DIR` 為基準，預設 `~/.gemini/antigravity`；Antigravity CLI 環境為 `~/.gemini/antigravity-cli`）：

- **PATH Shim**：`$ANTIGRAVITY_APP_DATA_DIR/bin/gmail`（Unix 為指向 `scripts/gmail` 的 symlink；Windows 為 `gmail.cmd` wrapper），此目錄不存放 binary 本體或憑證。
- **Go Binary 本體**：`$ANTIGRAVITY_APP_DATA_DIR/plugin_data/gmail/bin/gmail-1017793813`（Windows 為 `.exe`）。

## Credential 與驗證機制

- **Antigravity Managed Sign-in**：授權狀態由 Antigravity plugin 帳號系統代管；呼叫註冊後的 `gmail` 時自動帶入對應帳號權限。
- **停用獨立 CLI 憑證指令**：在 sandbox 環境下不可執行 `gmail auth login` 或 `gmail auth status`（後者僅檢查獨立 CLI 本機儲存，會失敗）；改以 `list_plugin_accounts` 查詢已連結帳號。
- **未登入處理**：回傳 `"needs authentication"` 時，顯示系統提供的 plugin card 請使用者點擊 **Sign in** 後再重試；若回傳 `no credentials found: run gmail auth login`，代表此帳號或版本尚不支援 Gmail 登入，須回報並停止。
- **權限不足（403）**：若 403 錯誤明確指出缺少權限或 scope，屬於授權範圍問題，不要求重新登入。
- **多帳號切換**：`AGY_ACCOUNT=me@example.com gmail readonly search "is:unread" --max 5`。
- **安全原則**：不搜尋檔案系統或環境變數中的 token／憑證，也不向使用者索取 token。

## 指令規範與安全邊界

| 規則項目 | 說明 |
| :--- | :--- |
| **Namespace 分離** | `readonly` 僅讀取（信件、討論串、草稿、標籤、歷史紀錄、設定）；`mutate` 會修改雲端狀態或寄信 |
| **Unattended 模式** | 設定 `AGY_UNATTENDED` 時，CLI 層會硬性阻擋所有外寄指令（`send`、`reply`、`react`、`forward`、`send-with-attachment`、`send-draft`），可改用 `mutate create-draft` 建立草稿供人工審核寄出 |
| **高風險操作確認** | 外寄信件必須先確認收件人、主旨與內文；`mutate trash` 僅在使用者明確要求時使用；`mutate delete-draft` 為永久刪除（繞過垃圾桶），刪除前須先複查 `readonly drafts` 與 `in:sent` |
| **禁用殼層過濾** | 禁止將 `gmail` 輸出管線交給 `grep`、`egrep`、`awk`，一律改用 Gmail search operators |
| **數量限制 Flag** | 僅支援 `--max N`（`search`/`threads`/`drafts` 預設 10、`read` 預設 5）；不可使用 `--max_results`、`--limit`、`-n`、`-l` |
| **負號查詢跳脫** | 查詢字串以 `-` 開頭時，前面須加 `--`，例如 `gmail readonly search --max 5 -- "-in:trash -in:spam"` |

## 常用指令與 Recipes

### 搜尋與閱讀郵件

```bash
gmail readonly search "has:attachment newer_than:35d" --max 25
gmail readonly read "Dentist" --max 5
gmail readonly read "is:starred" --page-token "TOKEN"
gmail readonly get MESSAGE_ID
gmail readonly get-thread THREAD_ID
```

> 搜尋策略：先以最穩定的核心關鍵字搜尋，不把未經確認的日期寫死在查詢中；`readonly read` 預設輸出乾淨純文字，一般閱讀避免加 `--json`；支援標準 16 碼 hex ID、server ID 與 Web UI Sapinto ID。

### 檢查與下載附件

```bash
# 1. 取得信件 JSON 結構，從 payload.parts 找出 filename 與 body.attachmentId
gmail readonly get MESSAGE_ID --json

# 2. 下載附件至指定路徑
gmail readonly download-attachment MESSAGE_ID ATTACHMENT_ID ./output.pdf
```

### 寄信、回覆、草稿、標籤與設定

```bash
# 寄信與回覆（Markdown 加 --md，表格用 --html + inline CSS）
gmail mutate send --to "user@example.com" --subject "Hello" --body "Hi!"
gmail mutate send --to "a@example.com" --cc "b@example.com" --subject "Update" --md --body-file ./body.md
gmail mutate send-with-attachment --to "user@example.com" --subject "Report" --body-file ./body.md --md --file ./report.pdf
gmail mutate reply --message MESSAGE_ID --body "Thanks!" --all --include-history
gmail mutate forward MESSAGE_ID "other@example.com"

# 草稿
gmail readonly drafts --json
gmail mutate create-draft --to "user@example.com" --subject "Draft" --md --body-file ./body.md
gmail mutate update-draft DRAFT_ID --md --body-file ./body.md
gmail mutate send-draft DRAFT_ID
gmail mutate delete-draft DRAFT_ID

# 標籤與整理
gmail readonly labels
gmail mutate create-label "Parent/Child" --bg-color red --text-color white
gmail mutate add-label MESSAGE_ID LABEL_ID
gmail mutate mark-read MESSAGE_ID
gmail mutate archive --thread THREAD_ID
gmail mutate batch-modify "MSG_ID1,MSG_ID2" --add-labels "LABEL1" --remove-labels "LABEL2"

# Emoji 反應與設定
gmail readonly list-reactions --message MESSAGE_ID --json
gmail mutate react --message MESSAGE_ID --emoji "🚧"
gmail mutate delete-reaction --message MESSAGE_ID --emoji "🚧"
gmail readonly get-vacation --timezone "Asia/Taipei"
gmail readonly forwarding-addresses
gmail readonly get-autoforward
gmail readonly sendas
gmail readonly history 12345
```
