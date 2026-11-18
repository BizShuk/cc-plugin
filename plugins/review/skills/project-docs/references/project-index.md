# 專案地圖 (Project Map)

`~/projects/` 第一層的簡介。本檔`只列第一層`; 分類底下有哪些專案不寫在這裡,
找專案時才讀當下的 `README.md`, 所以子專案增刪不必改本檔。

## 第一層 (Top Level)

### 分類 (Category)

容器目錄本身不是落點。

| 分類 | 收什麼 |
| ---- | ------ |
| `ai` | 會自己跑推理迴圈的 agent, 以及寫 agent 用的 SDK 與 Claude Code 設定 |
| `data` | 把外部資料取得, 正規化, 儲存與呈現的服務 (行情, 訊息, 影片, 活動, 時間序列) |
| `env_setup` | 機器初始化與 dotfiles; 底下的 `bin`, `pkg`, `scripts` 是它的子模組 |
| `game` | 遊戲原型, 引擎, 地圖編輯器與遊戲設計研究 |
| `iphone` | iOS app, 以及 iPhone 與 Mac 之間的同步工具 |
| `platform` | 跨服務共用的基礎: 身分, 金流, 監控與部署, Go SDK, webhook/CI, 編輯器 |
| `playground` | 產品點子與技術實驗 |
| `product` | 可上線的產品與產品研究 (營養, 穿搭, 衝浪, 影片剪輯, 作品集) |
| `research` | 只交付 Markdown 的知識工作: 教材, 證據型研究, 選型調查 |
| `tools` | 獨立的 CLI 與本機工具 (行程管理, LLM proxy, 備份, 翻譯, 磁碟, 影像) |

### 根目錄專案 (Root Project)

本身就是落點。

| 專案 | 用途 |
| ---- | ---- |
| `default` | 個人紀錄的純檔案庫 |
| `social` | Instagram 與 Threads 的單帳號社群工作台 (`Relay`) |
| `web` | 獨立的瀏覽器互動動畫習作 |

`tmp/` 是暫存區, 不是專案, 也不是落點。

## 找專案 (Finding a Project)

1. 用上表挑出`一個`分類; 內容就是根目錄專案的話, 直接到第 3 步。
2. 讀該分類下每個專案的業務定義, 挑出唯一符合的那一個:

    ```bash
    for f in ~/projects/<category>/*/README.md; do
      echo "== ${f%/README.md}"; grep -v '^\s*$' "$f" | grep -v '^#' | head -3
    done
    ```

3. 讀選中專案的 `README.md`, `AGENTS.md`, `docs/terminology.md` 之後再動作。

不需要索引檔或腳本: 分類清單以 `ls ~/projects/` 為準, 專案定義以各自的 `README.md` 為準。
第一層若新增或移除目錄, 只要更新本檔表格裡的一列。

## 簡介來源規則 (Purpose Line)

第 2 步只看 README 的第一段, 所以`第一段必須是一句完整的用途說明`。
以下開頭會讓專案辨識失真, `audit` 時列為漂移:

- 連結或來源標註 (`[Wiki](...)`, `來源: ...`)
- 清單項目或表格列 (`- pkg/ — ...`, `| # | ... |`)
- HTML 註解或 badge (`<!-- ... -->`)
