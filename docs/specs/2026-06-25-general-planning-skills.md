# 建立 General 規劃技能實作計畫 (Create General Planning Skills Implementation Plan)

本計畫旨在 `plugins/general` 目錄下新增兩項 AI 技能 (skills)，用以指導 AI Agent 進行「商業價值擴充」與「系統架構簡化、組織化與可擴充性」的規劃，並同步更新對應的配置與說明文件。

## 使用者審查項目 (User Review Required)

> [!IMPORTANT]
>
> 1. 所有 Markdown 文件（包括 `SKILL.md`、`README.md`、`walkthrough.md` 等）皆禁止使用 `**` (粗體) 語法，一律改用 `backtick` 進行高亮與強調。
> 2. 技能的 Frontmatter `name` 必須與其所在的子目錄名稱完全一致，使用 `kebab-case`。
> 3. 建立技能後，必須同步更新 `plugins/general/.claude-plugin/plugin.json` 與 `.claude-plugin/marketplace.json` 中的 `skills` 列表，以及 `plugins/general/README.md`。
> 4. 本計畫回覆與所有新撰寫之說明皆使用繁體中文，且專有名詞均會附上英文與圓括號。

## 開放性問題 (Open Questions)

> [!NOTE]
> 目前無開放性問題。兩個技能已實作至 `v1.1.0`，本計畫為設計紀錄；
> 各 `SKILL.md` 為唯一真實來源 (Single Source of Truth)，本文僅記設計意圖與優化決策。

## 優化紀錄 (Optimization Log — v1.1.0)

相較 `v1.0.0` 初版，兩技能依 `business-extract` 的品質標竿強化如下：

- `effort` 由 `medium` 升為 `high`：兩者皆需全庫掃描 + 策略綜整。
- 新增 `Step 1 蒐集脈絡與界定範圍`：明訂 folder/repo/file/text 的取材方式與噪音排除。
- 導入決策框架：`business-planner` 以 `RICE` 排序機會並用 `quadrantChart` 視覺化；
  `system-planner` 以「高改動 × 高耦合」交集鎖定重構 ROI 熱點。
- 內嵌 Mermaid 範例：機會優先矩陣、分層依賴圖。
- `Common Mistakes` / `Failure Modes` 改為表格，對齊 `business-extract` 風格。
- 強化落地：`business-planner` 增 `成功指標 (North Star)`；
  `system-planner` 增 `特徵測試安全網 + 絞殺榕 (Strangler-Fig) 漸進遷移`。
- 明訂三技能鏈路：`business-extract` → `business-planner` → `system-planner`。
- `business-planner` 的 `allowed-tools` 增 `WebSearch` 以支援競品/市場佐證。
- `僅規劃 (Only Planning)`：兩技能只產出規劃報告，不實作/不改碼/不改設定。
- 輸出目錄改為 `./plans/`（對齊全域慣例），取代原 `docs/proposals/`。
- 修正初版錯字（全形句號、簡體「配置驅動」）。

---

## 預期變更內容 (Proposed Changes)

### 實作規劃技能 (Implement Planning Skills)

#### [DONE] [SKILL.md](../../plugins/general/skills/business-planner/SKILL.md)

`business-planner` 已實作 (`v1.1.0`)。設計重點（完整內容見檔案）：

- 七步流程：`脈絡蒐集 → 核心價值 → 水平/垂直機會 → 痛點與市場契合 → RICE 排序 → 商業模式/MVP → 撰寫報告`。
- 決策框架：`RICE = Reach × Impact × Confidence ÷ Effort`，輔以 `quadrantChart` 機會優先矩陣。
- 報告結構升級為八章，新增 `成功指標 (North Star)` 與 `風險與假設`。
- `allowed-tools` 含 `WebSearch`（競品/市場佐證，須標來源與日期）。

#### [DONE] [SKILL.md](../../plugins/general/skills/system-planner/SKILL.md)

`system-planner` 已實作 (`v1.1.0`)。設計重點（完整內容見檔案）：

- 六步流程：`複雜度量測 → 解耦分層 → 目錄重整 → 插件化判斷 → 漸進遷移 → 撰寫計畫`。
- 量測先行：以 `git log` 改動頻率 + 扇入/扇出 + 循環相依鎖定「高改動 × 高耦合」熱點。
- 安全落地：重構前補特徵測試 (Characterization Test)，採絞殺榕 (Strangler-Fig) 分步可回滾。
- 內嵌分層依賴 Mermaid 圖；報告結構升級為七章，新增 `複雜度量測` 與 `風險與回滾`。
- 防過度設計：插件機制需先論證必要性，簡化以「移除」優先於「新增抽象」。

### 更新插件與系統配置 (Update Plugin & System Configurations)

#### [MODIFY] [plugin.json](../../plugins/general/.claude-plugin/plugin.json)

在 `skills` 列表中加入兩個新技能的相對路徑：
- `./skills/business-planner`
- `./skills/system-planner`

#### [MODIFY] [marketplace.json](../../.claude-plugin/marketplace.json)

在 `general` 插件的 `skills` 列表中更新這兩個新技能。

#### [MODIFY] [README.md](../../plugins/general/README.md)

更新 General 插件的 `README.md`，在 `## 技能 (Skills)` 表格中加入以下內容：
- `business-planner` | 分析既有系統業務，規劃如何擴充商業價值與業務範疇
- `system-planner` | 診斷現有系統複雜度，規劃如何進行架構簡化、目錄組織與可擴充性設計

#### [MODIFY] [CLAUDE.md](../../CLAUDE.md)

更新 `CLAUDE.md` 中的專案結構樹中關於 `plugins/general/skills/` 的註解，將這兩個新技能加進去。

---

## 驗證計畫 (Verification Plan)

### 自動化測試與驗證

1. 檢查 YAML frontmatter 語法與 `name` 欄位是否與子目錄名稱一致。
2. 執行 `npx skills add .` (若可用) 驗證技能加載。

### 手動驗證

1. 檢查所有新建立或修改的 Markdown 檔案，確保沒有使用 `**` 的粗體語法，全面替換為 `backtick`。
2. 確保沒有鬆散的技能檔案，即 `SKILL.md` 必須放置於其專屬的子目錄下。
3. 驗證產出的專有名詞是否有附帶英文及圓括號。
