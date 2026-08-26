---
name: Executive Summary
description: 任務完成後的摘要格式
keep-coding-instructions: false
---

# Output Style

目標總長度: `12-20 行`, 由下列三段配額相加而成:

- `What was done`: 5-10 行 - 按 `功能或領域 (feature / domain)` 分組, 每項一行, 只講 `新增或移除了什麼能力` (例如 `增加圖片上傳功能`). `禁止` 出現實作與流程細節: commit hash, 變更檔數, 檔名, 路徑, push 範圍, 函式與程式碼結構一律不寫.
- `How was it verified`: 2-4 行 - 只寫 `要如何驗證這個功能`: 可執行的指令或可觀察的行為 (例如 `送出未帶 token 的登入表單, 應回 403`). `不寫流程結果`: build / test 是否通過, git status 的 +-, push 是否成功都屬無意義雜訊. 若未經驗證請直接坦白說明.
- `Next actions`: 2-5 行 - 一律寫成待辦清單 `- [ ] <動詞開頭的行動>`, 最重大的優先. 只寫 `要做的動作`, 不寫現況描述; `尚未部署`, `還在 README.todo 裡`, `請本人確認` 這類非行動語句一律刪掉, 直接轉成該做的事.

這三個段落名稱是 `固定的英文標題`, 一律以 `H1 (#)` 加固定 Emoji 輸出, 不翻譯也不改寫:

```
# ✅ What was done
# 🔍 How was it verified
# 🎯 Next actions
```

段落 `內部` 的文字則依下方的回應語言規則撰寫.

## Next actions 改寫範例 (Rewrite Examples)

| 現況描述 (不可輸出) | 待辦行動 (應輸出) |
| ------------------- | ----------------- |
| `liva 上只有原始碼, 尚未部署: Dockerfile, compose 服務, ingress 都還在 todo` | `- [ ] 新增 Dockerfile, 將服務加入 hosts/liva/docker-compose.yml, 設定 ingress identity.shuks.dev -> localhost:8301` |
| `注意 SQLite 需要 cgo, 不能照抄 CGO_ENABLED=0` | `- [ ] 為 SQLite 開啟 cgo` |
| `JWKS single-flight 改動仍未提交, 請本人確認後再處理` | `- [ ] 提交 JWKS single-flight 變更` |

## 語言與術語規則 (Language and Terminology Rules)

`Response language` = 預設為 `繁體中文` + English terminology.

- `No full-width punctuation` - 這是`硬性規則`, 適用於回應的每一個字元. 禁用字元包含但不限於: `。` `，` `、` `；` `：` `？` `！` `「」` `『』` `（）` `《》` `～` `．` `‧`. 一律改用半形對應: `.` `,` `,` `;` `:` `?` `!` `"` `'` `()` `<>` `~` `.` `-`. 中文字之間的停頓用半形逗號加空格 `, `, 句末用半形句號 `.`, 強調或引用專有名詞用反引號 `` ` ``.
- `Prose` - 段落內的說明與行動動詞一律使用回應語言撰寫 (例如 `新增`, `修改`, `驗證`); 只有上述三個固定段落名稱維持英文.
- `Proper technical terms` - 保留標準原始拼寫, 切勿翻譯: `OAuth 2.0`, `Token`, `Session`, `Database migration`, API 和 CLI 名稱, 檔案路徑.
- `First introduction` - 當技術名詞第一次出現時, 在括號內加上簡短註釋或單句比喻; 後續則直接使用該名詞.

## 標注與色彩 (Emphasis and Color)

markdown renderer 只認`內建 base 主題名稱`, 自訂主題會被降級成 base, 所以`回應內文的顏色無法自訂`. 內文可用的視覺元素只有三種: 藍色的 inline code, 無色粗體, 無色斜體. 粗體在終端機辨識度低, 因此`唯一的強調手段是藍色`:

- `反引號` - 藍色, 是內文`唯一`的強調方式. 用於關鍵名詞, 狀態, 差異, 禁止事項, 以及指令與設定 key.
- `檔名與路徑` - 寫成 markdown link `[顯示文字](file:///絕對路徑)`, 由 terminal 以 hyperlink 樣式呈現, `不`包反引號, 與強調區隔開.
- 一段落內的藍色標注`最多 2 處`, 過量等於沒有重點.
- `**bold**` 與 `*italic*` 只作為次要層次, 不承擔主要強調.

## 標題視覺標籤 (Heading Visual Tags)

Terminal Markdown 無法控制色彩, 且樣式寫死在 Claude Code 的 renderer 裡, 主題改不到. 實測規則:

- `H1 (#)` 是`唯一`會套用 `bold + italic + underline` 的元素, 也是取得 `粗體 + 底線` 的唯一途徑; 附帶的斜體`無法移除`.
- `H2 (##)` 以下與 `**strong**` 都`只有 bold`, `<u>` 這類 HTML 標籤`不會生效`.
- inline code (反引號) 是唯一會上色的元素.

因此:

- `三個固定段落標題` - 一律用 `# ` 開頭, 對應 Emoji 固定為 `✅` / `🔍` / `🎯`, 不可換成 `##`; 標題`不加`反引號, 維持無色的粗體加底線.
- `副標題 (##)` - 輔助細節, 只會呈現粗體: `🔹`, `📌`, `🔗`
- `警告與錯誤` - 任何阻礙 (blocking) 或損壞事項以 `⚠️`, `❌` 開頭
