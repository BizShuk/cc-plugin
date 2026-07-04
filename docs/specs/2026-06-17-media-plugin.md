# 建立多媒體插件實作計畫 (Create Media Plugin Implementation Plan)

本計畫旨在新增一個名稱為 `media` 的模組化插件 (Modular Plugin)，用於多媒體影片生成與劇本創作。此插件將包含三個核心技能，以形成完整的劇本創作與影片生成提示詞工作流。

## 使用者審查項目 (User Review Required)

> [!IMPORTANT]
> - 插件結構：在 `plugins/media` 下建立完整的目錄，包含 `README.md`、`.claude-plugin/plugin.json` 以及三個獨立的技能目錄。
> - 技能實作：實作 `prompt-to-story-script`、`scene-to-video-prompt` 與 `character-setting` 三個技能。
> - 全域註冊：更新 `.claude-plugin/marketplace.json` 與 `CLAUDE.md`，將 `media` 插件註冊至系統中。
> - 格式約束：所有產出的 Markdown 內容（包括本計畫與技能文件）皆不使用 `bold` 語法，改用 `backtick` 進行強調。

## 開放性問題 (Open Questions)

> [!NOTE]
> 目前無開放性問題。

---

## 預期變更內容 (Proposed Changes)

### 建立多媒體插件 (Create Media Plugin)

#### [NEW] [README.md](../../plugins/media/README.md)

- 建立 `plugins/media/README.md`，說明 `media` 插件的用途與三個技能的工作流。

#### [NEW] [plugin.json](../../plugins/media/.claude-plugin/plugin.json)

- 建立 `plugins/media/.claude-plugin/plugin.json`，定義插件資訊與技能清單。

### 實作多媒體技能 (Implement Media Skills)

#### [NEW] [SKILL.md](../../plugins/media/skills/prompt-to-story-script/SKILL.md)

- 實作劇本創作優化技能。
- 功能：將簡單的想法擴展成具備 Setup、Turn、Payoff 故事弧線的畫面化腳本。

#### [NEW] [SKILL.md](../../plugins/media/skills/scene-to-video-prompt/SKILL.md)

- 實作分鏡提示詞轉換技能。
- 功能：將腳本拆分成多鏡分鏡，並按六要素（主體、動作、場景、鏡頭、光線、風格）輸出符合 Kling 等多分鏡介面的英文提示詞，並鎖定角色與場景的一致性。

#### [NEW] [SKILL.md](../../plugins/media/skills/character-setting/SKILL.md)

- 實作角色一致性設定技能。
- 功能：設計一套角色設定範本，產生詳細外觀與特徵描述，維持跨鏡頭的角色一致性。

### 更新系統配置 (Update System Configurations)

#### [MODIFY] [marketplace.json](../../.claude-plugin/marketplace.json)

- 在 `plugins` 陣列中新增 `media` 插件的註冊，包含其技能路徑。

#### [MODIFY] [CLAUDE.md](../../CLAUDE.md)

- 在專案結構與模組對應中加入 `media` 插件。
- 在 AI 技能與代理生態中更新說明。

#### [MODIFY] [README.md](../../README.md)

- 在 `AI 技能與代理生態` 部份加入 `media` 插件。

---

## 驗證計畫 (Verification Plan)

### 自動化測試
- 執行 `go test ./...` 確保現有的 Go 測試無損壞。

### 手動驗證
- 驗證 `plugins/media/.claude-plugin/plugin.json` 的 JSON 語法。
- 驗證 `.claude-plugin/marketplace.json` 的 JSON 語法。
- 檢查所有建立的 Markdown 檔案，確保無使用 `bold` 語法，全以 `backtick` 代替。
