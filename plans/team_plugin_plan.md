# 建立團隊插件實作計畫 (Create Team Plugin Implementation Plan)

本計畫旨在新增一個名稱為 `team` 的模組化插件 (Modular Plugin)，用於協助使用者規劃、設計與管理 AI `代理團隊 (Agent Team)` 的架構、角色提示與編排邏輯。

## 使用者審查項目 (User Review Required)

> [!IMPORTANT]
> - 插件結構：在 `plugins/team` 下建立完整目錄，包含 `README.md`、`.claude-plugin/plugin.json` 以及三個獨立的技能目錄。
> - 技能實作：
>   1. `team-design`：規劃跨職能 AI 代理團隊架構。
>   2. `role-generator`：依據五大原則（身分、職責、思考方式、格式、限制）與企業文化生成 `系統提示 (System Prompt)`。
>   3. `orchestration-config`：設計代理團隊的編排（Pipeline 或 Orchestrator）與專案層級共享設定。
> - 全域註冊：更新 `.claude-plugin/marketplace.json` 與 `CLAUDE.md`，將 `team` 插件註冊至系統中。
> - 格式約束：所有產出的 Markdown 內容（包括本計畫與技能文件）皆不使用 `bold` 語法，全以 `backtick` 代替。

## 開放性問題 (Open Questions)

> [!NOTE]
> 目前無開放性問題。

---

## 預期變更內容 (Proposed Changes)

### 建立團隊插件 (Create Team Plugin)

#### [NEW] [README.md](../plugins/team/README.md)

- 建立 `plugins/team/README.md`，說明 `team` 插件的用途與三個技能的工作流。

#### [NEW] [plugin.json](../plugins/team/.claude-plugin/plugin.json)

- 建立 `plugins/team/.claude-plugin/plugin.json`，定義插件資訊與技能清單。

### 實作團隊技能 (Implement Team Skills)

#### [NEW] [SKILL.md](../plugins/team/skills/team-design/SKILL.md)

- 實作團隊架構設計技能。
- 功能：根據專案需求與交付目標，規劃跨職能團隊的角色編制與主要職責。

#### [NEW] [SKILL.md](../plugins/team/skills/role-generator/SKILL.md)

- 實作角色提示生成技能。
- 功能：依據五大原則設計 `系統提示 (System Prompt)`，並可選融合 Meta/Google/Amazon/TikTok 等大廠的文化特質與技能需求。

#### [NEW] [SKILL.md](../plugins/team/skills/orchestration-config/SKILL.md)

- 實作團隊編排與共享配置技能。
- 功能：規劃 `協調者模式 (Orchestrator Pattern)` 或 `流水線模式 (Pipeline Pattern)` 的編排邏輯，並設定專案層級的共享規則。

### 更新系統配置 (Update System Configurations)

#### [MODIFY] [marketplace.json](../.claude-plugin/marketplace.json)

- 在 `plugins` 陣列中新增 `team` 插件的註冊，包含其技能路徑。

#### [MODIFY] [CLAUDE.md](../CLAUDE.md)

- 在專案結構與模組對應中加入 `team` 插件。
- 在 AI 技能與代理生態中更新說明。

#### [MODIFY] [README.md](../README.md)

- 在 `AI 技能與代理生態` 部分加入 `team` 插件。

---

## 驗證計畫 (Verification Plan)

### 自動化測試
- 執行 `go test ./...` 確保現有的 Go 測試無損壞。

### 手動驗證
- 驗證 `plugins/team/.claude-plugin/plugin.json` 的 JSON 語法。
- 驗證 `.claude-plugin/marketplace.json` 的 JSON 語法。
- 檢查所有建立的 Markdown 檔案，確保無使用 `bold` 語法，全以 `backtick` 代替。
