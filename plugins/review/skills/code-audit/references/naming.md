# 命名軸 (Naming Axis)

命名軸問的是：`讀名字能預測行為嗎? 同一個概念是否永遠用同一個詞?`

好名字的兩個條件：讀者能從名字推測行為，且同一概念在任何地方都穿同一件衣服。
判斷前先推斷 codebase 既有慣例 —— `既有慣例勝過任何外部風格書`。

## 執行步驟 (Procedure)

1. 依類別（檔案、型別、函式、變數、常數、設定鍵、API 端點）抽樣既有名稱，
   推斷實際慣例：casing、詞序、前綴。
2. `grep` 本次變更引入的新名稱，以及它觸及的每個概念的同義詞。
3. 以下表鏡頭逐一評分，基準是步驟 1 推斷出的慣例，不是通用指南。
4. 輸出 `current -> suggested` 並附一行理由，只提出`淨收益`的改名。

## 檢查鏡頭 (Dimensions)

| 鏡頭 (Dimension)       | 失敗樣態 (Failure)                                        |
| ---------------------- | --------------------------------------------------------- |
| 大小寫 (Casing)        | `userId` 與 `user_id` 並存; 慣例套用不均                  |
| 同義詞漂移 (Synonym)   | 同一操作混用 `fetch` / `get` / `load`                     |
| 語義空洞 (Vagueness)   | `data`, `handle`, `manager`, `util` 之類說不出任何事的名稱 |
| 誤導 (Misleading)      | 名字暗示的行為程式碼並未執行                              |
| 縮寫 (Abbreviation)    | 用了 codebase 其他地方一律拼全的縮寫                      |
| 對稱 (Symmetry)        | 有 `open` 沒有 `close`; 有 `start` 沒有 `stop`            |
| 尺度相稱 (Scope fit)   | 極小作用域用長名, 或廣域 API 用極簡名                     |

## 輸出範例 (Output)

```text
Naming review (推斷慣例: Go exported = PascalCase, 檔名 = snake_case)
- cmd/Read_Logic.go      -> read_logic.go (檔名 casing 漂移)
- func GetUser/FetchUser -> 同一操作只留一個動詞 (同義詞漂移)
- var d                  -> decoded (40 行作用域內語義空洞)
```

`不要為品味而改名`。每一項建議都必須減少歧義, 或回復 codebase 自身的一致性。

## 邊界 (Boundary)

識別字以外的術語漂移（文件與程式碼用詞不一）屬 [consistency.md](consistency.md)；
目錄與資料夾命名屬 `[[planner]]` 的系統軸。
