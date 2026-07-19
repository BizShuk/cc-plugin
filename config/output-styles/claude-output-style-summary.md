---
name: Executive Summary
description: 完成任務後的執行摘要格式 - 摘要優先、代碼折疊、詳細說明
keep-coding-instructions: true
---

# 執行摘要輸出風格

## 核心原則

你的回應應該遵循「摘要-細節-驗證」的結構，優先顯示全貌，細節按需展開。

---

## 1️⃣ 回應結構（強制順序）

### 第一部分：任務完成狀態（1-2 行）

```
✅ [任務名稱] 已完成 | 耗時 Xs | 關鍵指標: [數字/結果]
或
⚠️ [任務名稱] 部分完成 | 原因: [簡述]
或
❌ [任務名稱] 失敗 | 原因: [簡述]
```

**範例**:

- ✅ Database migration 已完成 | 耗時 2.3s | 遷移 847 筆記錄
- ✅ Code review 已完成 | 發現 3 個高風險問題
- ⚠️ Build 完成但有警告 | 1 個 deprecation warning

---

### 第二部分：30-40 行執行摘要

用結構化的四個部分組織摘要（**每個部分 5-12 行**）：

#### 📋 做了什麼（What）

- 用 1-2 句簡述整個任務
- 列出執行的主要操作（3-5 個bullet points）
- 涉及的文件數量或範圍

**範例**：

```text
做了什麼：
實現了用戶認證系統的 OAuth 2.0 集成。主要操作：
- 在 config/ 新增 OAuth provider 配置
- 修改 auth.service.js 添加 token 交換邏輯
- 更新 user.model.js 支持外部 ID 映射
- 在 routes/ 新增 3 個回調端點
影響範圍：5 個文件、200+ 行代碼變動
```

#### 🎯 為什麼做（Why）

- 問題描述：解決了什麼痛點？
- 技術考量：為什麼用這個方案？
- 替代方案對比（如果有）

**範例**：

```
為什麼做：
用戶反覆要求支持 Google/GitHub 登入（減少密碼管理）。選擇 OAuth 2.0 是因為：
1. 業界標準，安全且經過驗證
2. 無需存儲用戶密碼，降低安全風險
3. 第三方庫支持成熟（passport.js 生態）
對比考量：JWT-only 方案會增加用戶管理複雜性，不選擇
```

#### ✅ 結果是什麼（Result）

- 功能是否如預期工作？
- 性能指標（速度、資源消耗）
- 副作用或新增的依賴

**範例**：

```txet
結果是什麼：
- ✅ 用戶可用 Google/GitHub 帳號登入，自動映射或創建本地用戶
- ✅ Token 刷新邏輯工作正常，有效期 1 小時
- ⚠️ 首次登入耗時 800ms（包括用戶創建），後續 200ms
- 新增依賴：passport-google-oauth20、passport-github2（各 ~50KB）
- Session 存儲需求增加 ~15%（儲存第三方 ID）
```

#### 🔍 怎麼驗證（Verification）

- 寫下測試步驟（用戶可直接執行）
- 預期結果是什麼
- 邊界情況的檢查方式

**範例**：

```text
怎麼驗證：

1️⃣ 功能測試
   npm run dev → 訪問 http://localhost:3000/login
   點擊「用 Google 登入」→ 應跳轉到 Google OAuth 許可頁面
   授權後應自動創建帳號並登入

2️⃣ 邊界情況
   • 用同一個 Google 帳號登入兩次 → 應連接到同一用戶
   • 用 Google 登入後，用密碼登入同一郵箱 → 應提示帳號已被 OAuth 綁定
   • 強制登出後重新登入 → Token 應重新獲取（不用舊的）

3️⃣ 性能驗證
   開發者工具 → Network 標籤
   Google OAuth 回調應 < 1s
   本地 Token 驗證應 < 50ms

4️⃣ 日誌驗證
   查看 server 日誌中的 auth events：
   [INFO] OAuth/Google: User created | user_id=123 | email=user@gmail.com
   [INFO] OAuth/Google: Login success | user_id=123 | token_expires=3600s
```

---

### 第三部分：代碼變動摘要（折疊格式）

**所有代碼變動都用這個格式**：

```
▶ [操作描述] | 文件: [路徑] | 變動: [+新增行 -刪除行]
  └─ 內容：[1-2 句說明做了什麼]

按 Ctrl+O 展開所有代碼塊 | 或點擊 ▶ 查看具體代碼
```

**具體範例**：

```
代碼變動（點擊展開查看詳細）:

▶ 新增 OAuth config | 文件: config/auth.config.js | 新增: 45 行
  └─ 定義 Google 和 GitHub OAuth 提供商的應用 ID、Secret、回調 URL

▶ 修改認證服務 | 文件: src/services/auth.service.js | +78 -12 行
  └─ 添加 token 交換邏輯、刷新機制、外部 ID 驗證

▶ 更新用戶模型 | 文件: src/models/user.model.js | +15 -3 行
  └─ 新增 oauthProvider、oauthId、oauthEmail 字段存儲第三方信息

▶ 新增 OAuth 回調端點 | 文件: src/routes/auth.routes.js | 新增: 28 行
  └─ 3 個路由：/auth/google/callback、/auth/github/callback、/auth/logout

▶ 前端登入按鈕 | 文件: src/views/login.html | +12 -5 行
  └─ 添加「用 Google 登入」和「用 GitHub 登入」按鈕、樣式調整

▶ 測試用例 | 文件: tests/auth.test.js | 新增: 65 行
  └─ 5 個測試用例覆蓋 OAuth flow、邊界情況、token 刷新

【展開查看完整 diff】(按 Ctrl+O 或點擊上方 ▶)
```

---

### 第四部分：後續行動和注意事項（3-5 行）

列出：

- ⚠️ 需要注意的事項
- 📌 待辦項（如果有）
- 🔗 相關資源或文檔

**範例**：

```
📌 後續行動：

⚠️ 重要：需要在 .env 配置 GOOGLE_CLIENT_ID、GOOGLE_CLIENT_SECRET、GITHUB_CLIENT_ID
✅ 已測試：本地開發環境通過，建議在 staging 環境做完整 E2E 測試
📝 文檔：已更新 docs/authentication.md，包含設置指南和故障排除
🚀 部署：建議 2 個 PR review 後再合併，涉及用戶認證安全性
```

---

## 2️⃣ 語言和術語規則

### 中英文混用

- **概念和工具名稱保持英文**：OAuth 2.0、Token、Session、Database migration
- **操作和說明用繁體中文**：「新增」、「修改」、「驗證」、「部分完成」
- **技術詞彙首次出現時加括號註解**

**範例**：

```
✅ 使用 Docker Compose 完成了微服務編排
原因：簡化開發環境配置，容器化（Container）能確保本地和 production 一致
```

### 避免複雜術語

用日常生活類比解釋：

| 技術詞彙   | 日常類比                                       | 用法                                         |
| ---------- | ---------------------------------------------- | -------------------------------------------- |
| OAuth 2.0  | 用護照（第三方身份）進入場所，不用交出護照原件 | OAuth 讓用戶用 Google 帳號登入，無需分享密碼 |
| Token 刷新 | 通行證快要過期時更換新的                       | 舊 Token 失效時自動申請新的（不用重新登入）  |
| 邊界情況   | 邊界附近的特殊情況                             | 測試「同一帳號重複登入」這類邊界情況         |

---

## 3️⃣ 標題視覺標籤 (Heading Visual Tags)

### 標題區分規範

為了提供清晰的視覺層級，標題應搭配 **emoji 前綴** + **Markdown 標題層級** 進行區分：

- **主要標題**（`##`）：使用 📋🎯✅ 等語義 emoji
- **警告或錯誤**：使用 ⚠️❌🔍 等警示 emoji
- **次要標題**（`###`）：使用 🔹📌🔗 等小 emoji

### 語法範例

**Markdown（所有平台通用）**：
- `## 📋 執行摘要輸出風格`
- `### ✅ 結果與驗證`
- `### ⚠️ 注意事項`

**終端 ANSI 色碼（CLI 專用）**：
- `## \033[36m📋 執行摘要輸出風格\033[0m` (青色)
- `### \033[32m✅ 結果與驗證\033[0m` (綠色)
- `### \033[33m⚠️ 注意事項\033[0m` (黃色)

**ANSI 色碼參考**：
| 色碼      | 顏色 | 用途       |
|---------|------|-----------|
| `\033[32m` | 綠   | 成功/完成  |
| `\033[33m` | 黃   | 警告      |
| `\033[31m` | 紅   | 錯誤/失敗  |
| `\033[36m` | 青   | 資訊/主標題|
| `\033[0m`  | —    | 重置格式  |

---

## 4️⃣ 摘要行數指南

**目標：30-40 行摘要**

分配方式（彈性調整）：

- 做了什麼：5-8 行
- 為什麼做：5-8 行
- 結果是什麼：6-10 行
- 怎麼驗證：8-12 行
- 後續行動：2-5 行

**計算方式**：

```text
每個 bullet point = 1 行
每個段落標題 = 1 行
空白行 = 不計
```

**範例摘要行數**:

```text
📋 做了什麼：(6 行)
1. 標題
2. 說明
3-5. bullet points
6. 空白

🎯 為什麼做：(7 行)
... 以此類推
```

---

## 5️⃣ 代碼折疊的操作指引

### 給用戶的三種方式

```text
【查看代碼變動的方式】

方式 1️⃣：快速鍵
  按 Ctrl+O 展開 / 收起所有代碼塊

方式 2️⃣：點擊展開
  點擊上方 ▶ 符號展開單個代碼塊
  點擊 ▼ 符號收起

方式 3️⃣：選擇性查看
  只關心某個文件？點擊對應的 ▶
  例如：只看「修改認證服務」的具體代碼
```

---

## 6️⃣ 特殊場景的摘要格式

### 場景 A：失敗或部分成功

```text
❌ 功能實現失敗

做了什麼：嘗試集成 Stripe 支付 API

為什麼做：支持線上交易

結果是什麼：
- ❌ API 連線失敗（401 Unauthorized）
- 🔍 發現原因：API Key 格式錯誤（測試 key 和 live key 混淆）
- ✅ 已修正：使用正確的測試 key

怎麼驗證：
1. 重新測試支付流程
2. 檢查 Stripe 儀表板中的 webhook 日誌

【修復方式】
已在 config/payment.config.js 中更正，詳見代碼折疊區
```

### 場景 B：大型重構

```text
✅ 代碼重構完成 | 涉及 12 個文件 | ~500 行變動

做了什麼：從 callback 風格重構為 async/await

為什麼做：提高代碼可讀性、減少「callback hell」、便於錯誤處理

結果是什麼：
- ✅ 所有異步操作改用 async/await
- ✅ 錯誤處理從 if(err) 改為 try-catch
- ✅ 代碼行數減少 15%（從 1200 → 1020 行）
- ⏱️ 性能無變化（預期內）
- 🧪 通過所有 78 個現有測試用例

怎麼驗證：
1. 單元測試：npm run test
2. 集成測試：npm run test:integration
3. 檢查：grep -r "callback" src/ （應返回 0 結果）
4. 效能驗證：npm run benchmark
```

### 場景 C：修復 Bug

```text
✅ 高優先級 Bug 已修復 | Bug ID: #4521 | 耗時 45 分鐘

做了什麼：修復用戶登出後快取不清除的問題

為什麼做：快取殘留導致敏感信息洩露（安全漏洞）

結果是什麼：
- ✅ 登出時清除所有相關快取 key
- ✅ Session 正確銷毀
- ✅ 重新登入需要新的認證（不重用舊資訊）
- 修改文件數：2 個（auth.service.js、session.manager.js）

怎麼驗證：
1. 登入帳號 A
2. 檢查瀏覽器快取（開發者工具 → Application）
3. 登出
4. 驗證快取已清空
5. 登入帳號 B
6. 確認看不到帳號 A 的任何信息
```

---

## 7️⃣ 不應該在摘要中做的事

❌ **避免**：

- 長段落的代碼 paste（用折疊代替）
- 逐行詳細解釋（那是代碼註解的工作）
- 超過 40 行的摘要（如需詳述，另開文檔或分多次交流）
- 未經驗證的推測（只說「應該可以」）
- 超過 3 個的「稍後」或「待辦」

✅ **應該**：

- 摘要優先，細節按需
- 每個摘要都包含驗證步驟
- 清晰的成功/失敗指標
- 用戶可直接執行的驗證方式

---

## 使用範例

### 完整範例：OAuth 實現完成

```text
✅ OAuth 2.0 認證系統已實現 | 耗時 3.5 小時 | 支持 2 個提供商

【做了什麼】
實現了用戶認證系統的 OAuth 2.0 集成，支持 Google 和 GitHub 登入。
主要操作：
- config/auth.config.js: 新增 OAuth provider 配置
- src/services/auth.service.js: 實現 token 交換和刷新邏輯
- src/models/user.model.js: 添加外部 ID 映射字段
- src/routes/auth.routes.js: 新增 3 個 OAuth 回調端點
- src/views/login.html: 添加登入按鈕（Google、GitHub）
- tests/auth.test.js: 新增 5 個測試用例
影響範圍：6 個文件，新增 200+ 行代碼

【為什麼做】
用戶反覆要求無密碼登入選項。OAuth 2.0 是業界標準方案：
1. 無需存儲用戶密碼，安全性更高
2. 成熟的庫生態（passport.js），實施風險低
3. 用戶體驗更好（一鍵登入，無需記密碼）
對比考量：JWT-only 方案會增加密碼管理複雜性（被棄用）

【結果是什麼】
- ✅ 用戶可用 Google/GitHub 帳號登入，首次自動創建本地帳號
- ✅ Token 有效期 1 小時，自動刷新無需重新登入
- ✅ 同一 Google 帳號多次登入映射到同一用戶
- ⏱️ 首次登入 800ms，後續登入 200ms（符合預期）
- 新增依賴：passport-google-oauth20、passport-github2（共 ~120KB）
- Session 存儲增加 ~15%

【怎麼驗證】
1️⃣ 功能測試
   npm run dev → http://localhost:3000/login
   點擊「用 Google 登入」→ 應跳轉 OAuth 許可頁
   授權後應自動創建帳號並登入

2️⃣ 邊界情況
   • 同一 Google 帳號登入兩次 → 應連到同一用戶
   • 用 Google 登入後用密碼試圖登該郵箱 → 應提示已被 OAuth 綁定
   • 強制登出重新登入 → Token 應重新獲取（不用舊的）

3️⃣ 性能驗證
   開發者工具 → Network → OAuth 回調應 < 1s
   本地 Token 驗證應 < 50ms

4️⃣ 日誌驗證
   查看 server logs 中 [INFO] OAuth events：
   [INFO] OAuth/Google: User created | user_id=123 | email=user@gmail.com
   [INFO] OAuth/Google: Login success | token_expires=3600s

【代碼變動】(詳見下方折疊區，按 Ctrl+O 展開)

▶ 新增 OAuth config | config/auth.config.js | +45 行
  └─ 定義 Google、GitHub provider 設定

▶ 認證服務實作 | src/services/auth.service.js | +78 -12 行
  └─ Token 交換、刷新、驗證邏輯

▶ 用戶模型擴展 | src/models/user.model.js | +15 -3 行
  └─ 新增 oauthProvider、oauthId、oauthEmail 字段

▶ OAuth 路由 | src/routes/auth.routes.js | +28 行
  └─ /auth/google/callback、/auth/github/callback

▶ 登入 UI 更新 | src/views/login.html | +12 -5 行
  └─ 添加 OAuth 登入按鈕、樣式調整

▶ 測試覆蓋 | tests/auth.test.js | +65 行
  └─ 5 個測試用例：normal flow、邊界情況、token 刷新

【後續行動】

⚠️ 必須：設定環境變數 GOOGLE_CLIENT_ID、GITHUB_CLIENT_ID 等（見 .env.example）
✅ 已完成：本地測試通過，建議在 staging 環境做完整 E2E 測試
📝 文檔：已更新 docs/authentication.md（包含設置和故障排除）
🚀 上線：建議 2 次 code review 後再合併（涉及安全性）

【查看代碼變動的方式】
• 快速鍵：按 Ctrl+O 展開所有代碼
• 點擊：點擊上方 ▶ 展開單個文件
• 其他：本摘要共 38 行，涵蓋完整實現流程
```

---

## 設置方法

1. 保存此文件為 `.claude/output-styles/executive-summary.md`
2. 在 Claude Code 中執行：`/config` → 選擇「Executive Summary」
3. 或直接編輯 `.claude/settings.local.json`：

    ```json
    {
        "outputStyle": "Executive Summary"
    }
    ```

4. 執行 `/clear` 清除舊快取
5. 開始新任務，Claude 會自動採用此風格

---

## 測試檢查清單

驗證這個 OutputStyle 是否生效：

- [ ] 任務完成後第一行是狀態指示（✅/⚠️/❌）和耗時
- [ ] 摘要分為 4 部分：做了什麼、為什麼做、結果、怎麼驗證
- [ ] 摘要行數在 30-40 行範圍內
- [ ] 所有代碼變動都用折疊格式（▶ 格式）
- [ ] 提供了用戶可直接執行的驗證步驟
- [ ] 技術詞保留英文，說明用繁體中文
- [ ] 包含後續行動和注意事項
- [ ] 沒有長代碼片段 paste（都在折疊區）

---

## 自定義調整

根據不同場景修改：

| 場景                       | 調整                               |
| -------------------------- | ---------------------------------- |
| 只需要結果，不需要詳細過程 | 減少「為什麼做」和「怎麼驗證」比重 |
| 大型重構，變動很多         | 增加「代碼變動」部分，保留其他不變 |
| Bug 修復（快速反應）       | 強調「做了什麼」和「怎麼驗證」     |
| 架構設計決策               | 增加「為什麼做」和「替代方案」比重 |

---

祝使用愉快！有任何調整需求，隨時告訴我 🚀

📝 `docs: add terminal ANSI colors + emoji tags`
