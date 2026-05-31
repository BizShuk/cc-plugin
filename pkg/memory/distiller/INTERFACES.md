# 外部工具介面設計文件 (External Tool Interfaces)

本文件記錄了各個外部記憶工具（`claude-mem`、`agentmemory`、`mempalace`、`gbrain/working` 與 `hermes`）的確切整合介面、路徑與寫入方式。

---

## 1. claude-mem

### 實體路徑 (Physical Path)
* `SQLite` 資料庫路徑：`/Users/shuk/.claude-mem/claude-mem.db`

### 資料表結構 (Table Schema - `observations`)
我們主要讀取 `observations` 資料表：

| 欄位名稱 (Column Name) | 型態 (Type) | 說明 (Description) |
| --- | --- | --- |
| `id` | `INTEGER PRIMARY KEY` | 自動遞增 ID |
| `memory_session_id` | `TEXT` | 對話階段 ID (Session ID) |
| `project` | `TEXT` | 專案名稱 |
| `text` | `TEXT` | 原始文字內容 |
| `type` | `TEXT` | 觀察類型 (例如：`discovery`) |
| `title` | `TEXT` | 記憶標題 |
| `subtitle` | `TEXT` | 記憶副標題 |
| `facts` | `TEXT` | 事實列表 (JSON string, e.g. `["fact 1", "fact 2"]`) |
| `narrative` | `TEXT` | 詳細敘述文字 |
| `concepts` | `TEXT` | 關聯概念 (JSON string) |
| `created_at` | `TEXT` | ISO-8601 時間字串 (e.g. `2026-05-16T14:28:22.924Z`) |
| `created_at_epoch` | `INTEGER` | Epoch 毫秒時間戳記 (Epoch Milliseconds) |

### 查詢語法 (Query Example)
```sql
SELECT id, memory_session_id, project, title, subtitle, facts, narrative, created_at_epoch 
FROM observations 
WHERE created_at_epoch > ? 
ORDER BY created_at_epoch ASC;
```

---

## 2. agentmemory (long-term 查詢層)

### 整合介面 (Integration Interface)
* 套件類型：npm 套件 (`@agentmemory/agentmemory`)，本地啟動於 `http://localhost:3111`。
* 寫入端點 (Write Endpoint)：`POST http://localhost:3111/agentmemory/remember`

### 寫入資料載荷 (Write Payload)
```json
{
  "content": "記憶內容字串 (Memory content string)",
  "concepts": ["主題/標籤1", "主題/標籤2"],
  "files": ["相關檔案路徑1", "相關檔案路徑2"]
}
```

### 驗證與狀態查詢 (Health & Status Check)
* 伺服器健康度 (Health)：`GET http://localhost:3111/agentmemory/health`
* CLI 狀態指令 (Status Command)：`npx @agentmemory/agentmemory status`

---

## 3. mempalace (真實持久層)

### 整合介面 (Integration Interface)
* 套件類型：`uv` tool 安裝之 Python CLI (`/Users/shuk/.local/bin/mempalace`)。
* 儲存目錄 (Storage Directory)：`~/.mempalace/palace/`
* 寫入機制 (Ingestion Mechanism)：檔案掃描與挖掘 (Directory mining)。

### 寫入流程 (Write Flow)
1. 在暫存目錄下依 Room 建立 Markdown 檔案 (例如 `temp_dir/<room_name>/<fact_id>.md`)。
2. 執行 CLI 指令將檔案挖掘並寫入資料庫：
   ```bash
   mempalace mine <temp_dir> --wing <wing_name>
   ```

### 查詢指令 (Search Command)
```bash
mempalace search "關鍵字"
```

---

## 4. gbrain/working (對話關聯層)

### 實體路徑 (Physical Path)
* `gbrain` 知識庫目錄：`/Users/shuk/brain/`
* `working` 對話暫存目錄：`/Users/shuk/brain/working/` (儲存每位聯絡人或主題的 `*.md` running notes)

---

## 5. hermes 接線點 (hermes Integration Points)

我們可以使用以下三種方式之一將 `gbrain/working` 的對話關聯與 `agentmemory` 長期記憶注入 `hermes` 的對話脈絡：

### 擴充點 1：全域提示詞 (AGENTS.md)
* 設定檔路徑：`~/.hermes/profiles/devops/AGENTS.md` (軟連結自 `config/CLAUDE.global.md`)。
* 方式：直接在提示詞中指示 `hermes` 主動讀取與更新 `/Users/shuk/brain/working/<topic>.md` 的內容。

### 擴充點 2：事件鉤子 (Hooks)
* 設定檔路徑：`~/.hermes/profiles/devops/hooks/`
* 運作方式：建立 `HOOK.yaml` 與 `handler.py`，訂閱 `agent:start` 與 `agent:end` 事件。
  * `agent:start`：在對話開始前載入 `working note` 及 `agentmemory` 並寫入對話快照。
  * `agent:end`：在對話結束後將新的互動摘要追加回 `working note`。

### 擴充點 3：指令封裝 (CLI Wrapper)
* 方式：建立一個封裝腳本 (例如 `hermes-wrap.sh`) 包裹原生的 `hermes` 執行指令。在啟動前與結束後執行 python 腳本讀寫 `/Users/shuk/brain/working/`。
