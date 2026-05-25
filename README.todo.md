# 工作區改進待辦清單 (Workspace Improvements Todo List)

此檔案盤點了 `cc-plugin` 專案中所有可以優化的設定、說明文件與路徑問題。

## 設定與配置 (Configuration & Settings)

- 主題 (Topic)：`.claude-plugin/plugin.json` 中的路徑設定不一致
    - 原因 (Why)：檔案中的 `commands`、`hooks`、`mcpServers`、`outputStyles`、`themes` 與 `monitors` 的路徑與專案實際結構不符（例如專案中無 `commands` 資料夾，而 `hooks.json` 實際上是在 `hooks/` 目錄下而非 `config/` 目錄下）。
    - 方法 (How)：修正 `plugin.json` 內的路徑配置，使之符合目前專案的實際結構。

- 主題 (Topic)：`hooks/hooks.json` 僅為預留位置 (Placeholder)
    - 原因 (Why)：`PostToolUse` 鉤子目前只執行 `echo hook`，並無實質功能，無法在檔案編輯後自動進行排版或程式碼檢查。
    - 方法 (How)：將 `PostToolUse` 修改為能執行 `go fmt` 或 `golangci-lint` 的實用指令。

- 主題 (Topic)：缺乏自動建立日誌目錄 `logs` 的機制
    - 原因 (Why)：`monitors/monitors.json` 中設定了對 `./logs/error.log` 與 `./logs/access.log` 的監控。然而專案根目錄下預設並沒有 `logs` 目錄，在某些平台上直接執行 `tail` 會因為路徑不存在而失敗。
    - 方法 (How)：在專案中新增空目錄 `logs` 並加上 `logs/.gitkeep` 檔案納入 git版本控制，或者在 `run.sh` 裡加上 `mkdir -p logs`。

## 說明文件與文字修正 (Documentation & Typos)

- 主題 (Topic)：`CLAUDE.md` 的專案結構描述與實際不符
    - 原因 (Why)：`CLAUDE.md` 中將專案結構寫為 `my-plugin/` 子目錄，但實際上所有設定都直接位於根目錄下。這會對開發者理解專案架構造成誤導。
    - 方法 (How)：重新更新 `CLAUDE.md` 中的結構描述，移除 `my-plugin/` 的層級，使其符合當前專案根目錄結構。

- 主題 (Topic)：`plans/agent/refactor-agent-vs-specialized-agents.md` 的重疊章節
    - 原因 (Why)：檔案末尾出現了兩次完全相同的 `總結回答` 章節，造成了文件冗餘。
    - 方法 (How)：移除該檔案中重複多餘的章節與表格。

- 主題 (Topic)：`README.md` 缺少整體架構說明
    - 原因 (Why)：目前的 `README.md` 只列出了 MCP 與套件的安裝指令，缺乏對專案功能定位（如 `run.sh` 的作用、`pkg/` 的配置）的說明。
    - 方法 (How)：增補 `README.md` 的內容，詳細介紹此專案作為一個「全域設定配置庫 (Configuration Repository)」的定位與初始化指南。

## 路徑與環境適應性 (Paths & Environment Compatibility)

- 主題 (Topic)：`run.sh` 軟連結 (symbolic link) 建立強固性不足
    - 原因 (Why)：指令稿中使用了 `ln -sf` 將設定檔案軟連結至 `$HOME/.claude/` 與 `$HOME/.hermes/` 等，但在全新環境中，若這些目標目錄本身不存在，軟連結的建立會出錯或失敗。
    - 方法 (How)：在執行 `ln -sf` 前，先使用 `mkdir -p` 指令確保 `$HOME/.claude` 與 `$HOME/.hermes` 等目標目錄已被建立。

- 主題 (Topic)：`run.sh` 缺乏跨平台支援 (Cross-platform Compatibility)
    - 原因 (Why)：此專案目前的初始化指令稿僅為 Shell Script 格式（`run.sh`），無法在 Windows 原生環境下執行，限制了跨平台的便利性。
    - 方法 (How)：新增 Windows 相容的 PowerShell 指令稿（例如 `run.ps1`），或以 Node.js 撰寫跨平台的初始化工具。

- 主題 (Topic)：`run.sh` 反向連結本地路徑的安全隱憂
    - 原因 (Why)：腳本中包含將本機家目錄下的 `.gemini` 等敏感設定軟連結回本機專案目錄 `config/` 的操作，若不慎將此類軟連結或本機敏感資訊提交至 git，將可能造成安全風險。
        - 方法 (How)：在 `.gitignore` 中明確排除 `config/` 下因連結產生的本機目錄與設定檔，並在 `run.sh` 中加強目錄是否存在的條件判斷。
