# gmail-workflows

來源：`.agents/skills/gmail-workflows/SKILL.md`（`claude-office-skills`，MIT，version `1.0.0`）

基於 n8n workflow templates 與 Google Apps Script 的 Gmail 自動化工作流設計指南，涵蓋附件管理、郵件自動分類、封存與 Google Drive 整合。

## 宣告介面與依賴

- **MCP Server**：`google-workspace-mcp`
- **MCP Tools**：`gmail_search`、`gmail_get_attachments`、`gmail_apply_label`、`gmail_archive`、`gdrive_upload`、`gdrive_create_folder`
- **支援語言**：`en`、`zh`

## 四大核心工作流範本

| 範本 | 流程摘要 | 關鍵設計 |
| :--- | :--- | :--- |
| **Gmail Attachment Manager** | Trigger → Filter → Extract Attachments → Upload to Drive → Apply Label & Archive → Notify | 限制寄件網域與檔案大小（如 25 MB），依 `{year}/{month}` 建資料夾，啟用 hash／檔名去重 |
| **Invoice Auto-Archiver** | 偵測帳單/發票 → 抽取 PDF → OCR 解析欄位 → 存至 Drive → 更新試算表 → 封存加星號 | 支援多語關鍵字（`invoice`、`bill`、`statement`、`付款`、`发票`），寫入 Invoice Tracker |
| **Client Communication Organizer** | 依網域、主旨、未讀狀態分流 → 套用標籤、轉寄、存附件、建任務或發緊急通知 | 支援條件式路由（Conditional Routing）與 Slack / SMS 通知 |
| **Email Analytics Dashboard** | 定期統計收發件數、平均回覆時間、未讀與附件數、熱門寄件者與忙碌時段 | 輸出至 Google Sheets 並發送 Slack 週報 |

## 實作路徑與模式

- **n8n**：組合 `n8n-nodes-base.gmailTrigger` → `if`（`hasAttachment`）→ `gmail`（`getAttachments`）→ `googleDrive`（`upload`）。
- **Google Apps Script**：使用 `GmailApp.search('has:attachment is:unread')` 抽取附件並透過 `DriveApp` 存檔重命名（`yyyy-MM-dd_from_filename`），搭配 `ScriptApp.newTrigger().timeBased().everyMinutes(5)` 定時執行。
- **常見編排模式**：`Filter → Process → Organize → Notify`、每日批次處理（Batch Processing）、依檔案類型／寄件網域／檔案大小的條件路由。

## 安全與除錯要點

- **最小權限 OAuth Scopes**：僅讀取使用 `gmail.readonly`、標籤／封存使用 `gmail.modify`、寫入 Drive 使用 `drive.file`。
- **資料隱私**：不將郵件內文寫入 log，憑證採安全儲存，並設定資料保留政策。
- **常見問題對策**：偵測不到附件時檢查 MIME 過濾與輪詢頻率；重複檔案以 hash 或檔名去重；遇 API rate limit 時降低觸發頻率並改採批次處理。
