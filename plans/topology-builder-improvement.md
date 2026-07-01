# [Goal Description]

對 `topology-builder` 進行 review 並改進其拓撲準確性 (topology accuracy) 與噪音排除 (noise exclusion) 機制。主要是為了避免在拓撲圖中引入不必要的基礎設施噪音（例如 logging、config、helpers 等瑣碎檔案或依賴）、防止跨 zone 的同名實體誤合併，以及確保連線邊僅反映實質的業務/架構呼叫，而非低階程式碼引用。

## User Review Required

> [!IMPORTANT]
> 此項改進會修改 `topology-builder` 的全域規則與 workflow prompts。若現有拓撲圖中已包含基礎設施節點（如 common utilities），在重新執行 topology-builder 時，它們會被判定為 noise 而遭到排除或不再建邊。

## Proposed Changes

### topology-builder 技能 (topology-builder skill)

---

#### [MODIFY] [SKILL.md](../plugins/general/skills/topology-builder/SKILL.md)

1. 在 `Identity Rules` 區段新增 `Noise Exclusion Rules` (噪音排除規則) 子區段，明確界定哪些不是實體：
   - 排除測試檔案與 mocks。
   - 排除 trivial helper packages / files (如 logger, json utils, config parsers)。
   - 排除外部第三方框架 (如 cobra, go-homedir, lodash)。
2. 在 `Edge Format` 區段新增 `Dependency Filtering Rules` (依賴過濾規則) 子區段：
   - 僅能對實質的業務或架構依賴建邊 (如 RPC 呼叫、資料庫寫入、訊息佇列發布)。
   - 禁止對基礎設施調用建邊 (如 `calls [[logger#Log]]` 或 `reads-from [[config#Load]]`)。
3. 更新 `Common Mistakes` 表格，新增以下幾點：
   - 誤將 common helper 或 logger 視為實體並建邊。
   - 將不同 zone 的相似名稱 (例如不同服務的 `handler`) 誤判為相同實體。

#### [MODIFY] [workflow.js](../plugins/general/skills/topology-builder/workflow.js)

1. 在 `RULES` 常數中同步新增 `noise exclusion`、`topology accuracy` 與 `alias merging` 的指導規則。
2. 調整 `Discover` phase 的 prompt，明確指示 Agent 排除測試、mocks 及基礎設施 helper 套件。
3. 調整 `Identify` phase 的 prompt，加強對於 context / zone 區分的要求，避免假陽性 (false-positive) 的實體合併。
4. 調整 `Extract` phase 的 prompt，指示 Agent 排除瑣碎的內部函式 (getters/setters/boilerplate)，只專注於關鍵業務維度。
5. 調整 `Connect` phase 的 prompt，指示 Agent 只針對實質業務與系統流轉依賴建邊，過濾無謂 of utility 調用邊。

---

## Verification Plan

### Automated Tests
由於這是一個以 LLM 驅動的 workflow 腳本，我們將針對其 `Verify` 階段腳本與手動測試進行驗證：
- 執行 `go test ./...` 確保沒有破壞 Go 專案原本的測試。
- 透過 mock 資料或在 references/ 測試環境下，驗證 Verification 腳本對 noise 檢測的相容性。

### Manual Verification
- 檢視修改後的 `SKILL.md` 與 `workflow.js` 代碼，確保其符合 `agentskills.io` 規範（如 frontmatter 各必填欄位、格式等）。
- 執行 Verification 腳本（於 `SKILL.md` 底部的 3 個 bash 檢查指令），確認驗證規則本身依然能正常執行。
