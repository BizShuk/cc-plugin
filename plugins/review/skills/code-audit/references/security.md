# 安全軸 (Security Axis)

安全軸問的是：`這段程式碼會洩漏資料, 還是會被濫用?`
本軸是`防禦性審查`：找出既有程式碼的弱點並提出修補方向，不產出攻擊工具或利用步驟。

每筆發現都必須指向具體 `file:line`；通用威脅清單不算發現。

## 檢查鏡頭 (Lenses)

### 機密與憑證 (Secrets and Credentials)

- 原始碼、測試資料、設定檔或範例檔內出現真實 API key、token、密碼、私鑰。
- 機密寫進版本控制的 `.env`、fixture 或 commit message。
- 機密以明文出現在日誌、錯誤訊息、URL query string 或 panic stack。
- 憑證缺乏輪替路徑：硬編死值且無環境變數或 secret store 來源。

### 輸入驗證 (Input Validation)

- 外部輸入（HTTP body/query、CLI 參數、檔案、佇列訊息）未驗證型別、長度或範圍即使用。
- 使用者控制字串進入 SQL、shell、模板、正規表示式或序列化器而未參數化或跳脫。
- 路徑由輸入拼接卻未正規化，導致可跳出預期目錄 (path traversal)。
- 未設上限的解析：無限制的 body size、解壓縮、遞迴深度或分頁筆數。

### 認證與授權 (Authn and Authz)

- 端點或指令缺少身分驗證，或驗證只在前端/中介層做而非資源層。
- 有身分卻無`資源歸屬`檢查：能讀寫別人的物件 (IDOR)。
- 角色、方案 (plan)、租戶邊界的判斷散落多處且不一致（與一致性軸交叉）。
- Token 驗證省略簽章、期限、audience 或 issuer 其中之一。

### 敏感資料處理 (Sensitive Data Handling)

- 個資、金流、健康資料未遮罩即寫入日誌、metrics label 或 error report。
- 回應體洩漏內部細節：stack trace、SQL 語句、內部主機名或檔案路徑。
- 快取或暫存檔存放敏感資料且未設權限或清理。

### 傳輸與加密 (Transport and Crypto)

- 明文 HTTP、關閉憑證驗證 (`InsecureSkipVerify`)、自訂 TLS 而未設定最低版本。
- 自行實作加密或雜湊；密碼使用快速雜湊而非 bcrypt/argon2 這類慢雜湊。
- 使用非密碼學安全的亂數產生 token、session id 或重設碼。

### 依賴與供應鏈 (Dependencies)

- 依賴版本未釘定，或使用已知有漏洞的版本。
- 引入來源不明、長期未維護的套件以取得一個小功能。
- 安裝腳本或 build 階段從網路抓取未驗證校驗碼的產物。

### 執行環境 (Runtime Surface)

- 除錯端點 (pprof, debug handler) 或管理介面在正式環境仍開放。
- 檔案、socket 或暫存目錄權限過寬。
- 錯誤處理吞掉例外後`繼續執行`，讓失敗的檢查等同通過（fail-open）。

## 執行步驟 (Procedure)

1. 先定位`信任邊界`：哪裡接收外部輸入、哪裡驗證身分、哪裡對外發出請求。
2. 沿邊界追資料流：輸入 -> 驗證 -> 使用 -> 輸出/日誌，找出跳過驗證的分支。
3. 針對機密與依賴做全域掃描（`grep` 常見 key 樣式、檢視依賴清單）。
4. 每筆發現寫成 `風險 -> 觸發條件 -> 修補方向`，附 `file:line` 與嚴重性。

## 邊界 (Boundary)

- 本軸`不`產出可直接使用的 exploit、攻擊腳本或繞過偵測的方法。
- 需要實際執行才能判定的漏洞（競態、記憶體破壞）標為`待驗證`並說明如何驗證。
- 測試掩蓋正式環境行為的疑慮交給 `[[anti-sabotage]]`。
