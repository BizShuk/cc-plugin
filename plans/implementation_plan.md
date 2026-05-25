# cc-plugin 工作區改進實作計畫 (cc-plugin Workspace Improvements Implementation Plan)

此計畫旨在解決 `README.todo.md` 中列出的各項設定、說明文件、強固性與跨平台支援問題。

## 使用者審查需求 (User Review Required)

> [!IMPORTANT]
> - 我們將新增 `run.ps1` PowerShell 指令稿以支援 Windows 系統。
> - 我們將在 `hooks/hooks.json` 中配置 Go 檔案自動格式化鉤子，在編輯或新增 Go 檔案時自動執行 `go fmt`。

## 開放性問題 (Open Questions)

無。

## 預期變更 (Proposed Changes)

---

### 設定與配置 (Configuration & Settings)

#### [MODIFY] [plugin.json](file:///Users/shuk/projects/cc-plugin/.claude-plugin/plugin.json)
- 移除不存在的 `commands` 與 `outputStyles` 設定。
- 將 `hooks` 路徑從 `./config/hooks.json` 修正為 `./hooks/hooks.json`。
- 將 `mcpServers` 路徑從 `.mcp-config.json` 修正為 `.mcp.json`。
- 修正 `experimental` 中的 `monitors` 路徑為 `./monitors/monitors.json`，並移除不存在的 `themes` 設定。

#### [MODIFY] [hooks.json](file:///Users/shuk/projects/cc-plugin/hooks/hooks.json)
- 更新 `PostToolUse` 鉤子，在 `Write` 與 `Edit` 工具完成後，執行格式化腳本 `hooks/post-tool.sh`。

#### [NEW] [post-tool.sh](file:///Users/shuk/projects/cc-plugin/hooks/post-tool.sh)
- 新增 `PostToolUse` 的處理腳本，偵測被寫入或編輯的檔案是否為 `.go` 檔案，若是則執行 `go fmt`。

---

### 說明文件與文字修正 (Documentation & Typos)

#### [MODIFY] [CLAUDE.md](file:///Users/shuk/projects/cc-plugin/CLAUDE.md)
- 修正專案結構描述，移除 `my-plugin/` 的層級，改以根目錄為基準描述。

#### [MODIFY] [refactor-agent-vs-specialized-agents.md](file:///Users/shuk/projects/cc-plugin/plans/agent/refactor-agent-vs-specialized-agents.md)
- 移除檔案末尾重複多餘的 `總結回答` 章節。

#### [MODIFY] [README.md](file:///Users/shuk/projects/cc-plugin/README.md)
- 增補整體架構說明，詳細介紹本專案作為「全域設定配置庫 (Configuration Repository)」的定位與初始化指南。

---

### 路徑與環境適應性 (Paths & Environment Compatibility)

#### [MODIFY] [run.sh](file:///Users/shuk/projects/cc-plugin/run.sh)
- 在執行 `ln -sf` 前，先使用 `mkdir -p` 指令確保 `$HOME/.claude`、`$HOME/.gemini` 與 `$HOME/.hermes` 等目標目錄已被建立。
- 加上 `mkdir -p logs` 以自動建立日誌目錄。
- 加強對家目錄下 `.gemini`、`.claude`、`.hermes` 等目錄是否存在之條件判斷，避免建立無效的軟連結。

#### [NEW] [run.ps1](file:///Users/shuk/projects/cc-plugin/run.ps1)
- 新增與 `run.sh` 功能對等且具備目錄防禦性建立、路徑檢查的 PowerShell 初始化指令稿，以支援 Windows 原生環境。

#### [MODIFY] [.gitignore](file:///Users/shuk/projects/cc-plugin/.gitignore)
- 調整 `logs/` 排除規則，並確保 `logs/.gitkeep` 可以被提交。

#### [NEW] [.gitkeep](file:///Users/shuk/projects/cc-plugin/logs/.gitkeep)
- 在 `logs/` 目錄下建立 `.gitkeep` 檔案，確保空目錄可納入版本控制。

---

## 驗證計畫 (Verification Plan)

### 自動化測試 (Automated Tests)
- 執行 `bash run.sh` 與 `pwsh -File run.ps1`（如環境支援）確保連結建立成功且無報錯。
- 透過模擬 `PostToolUse` 的 payload 來測試 `hooks/post-tool.sh`，確保當傳入非 Go 檔案時不處理，傳入 Go 檔案時執行 `go fmt`。

### 手動驗證 (Manual Verification)
- 檢查家目錄下的軟連結是否正確指向專案目錄。
- 檢查 `.gitignore` 狀態。
