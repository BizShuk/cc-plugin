# Impeccable (`pbakaus/impeccable`)

- Source: <https://github.com/pbakaus/impeccable> · Docs: <https://impeccable.style/docs>
- Author: Paul Bakaus · License: Apache-2.0 (起始自 Anthropic 的 [frontend-design](https://github.com/anthropics/skills/tree/main/skills/frontend-design) skill)
- Stats (2026-11-21 via GitHub API): 114,094 stars, 6,782 forks, last push 2026-11-21, CLI v4.10.0 (`package.json`), Rust engine v0.1.27 (`Cargo.toml`)
- 支援工具: Claude Code, Cursor, Codex CLI, GitHub Copilot, Gemini CLI, Google Antigravity, Grok Build, OpenCode, Pi, Hermes Agent, Devin, Replit Agent, Kiro, Trae, Rovo Dev, Qoder, Mistral Vibe, DeepSeek Harness, Veto 等

## 定位與核心功能 (Brief)

`Impeccable` 是一套專為 AI coding agent 打造的「前端設計詞彙系統與反模式偵測引擎」, 用於解決 AI 生成介面常見的「AI slop」(例如全站 `Inter` 字體, 紫藍漸層, 卡片套卡片, 彩色底灰字, 每個標題上方放圓角方塊 icon).

它不只是單一 `SKILL.md` prompt, 而是一整套跨工具鏈的系統:

- **24 個設計指令** (含 `craft` 棄用別名, 23 個主力指令):
  - `Plan & Set Up`: `init` (`PRODUCT.md`), `document` (`DESIGN.md`), `shape`, `extract`
  - `Steer / Enhance / Refine`: `bolder`, `quieter`, `distill`, `colorize`, `typeset`, `layout`, `animate`, `delight`, `overdrive`, `polish`, `harden`, `onboard`
  - `Evaluate & Fix`: `critique` (heuristic UX 評分), `audit` (a11y / 效能 / 響應式 P0-P3 報告), `clarify`, `adapt`, `optimize`
  - `Iterate Live`: `live` (本機 dev server 掛載瀏覽器 overlay 挑元素切換變體並回寫原始碼), `generate` (免手動挑選直接生成元素變體)
- **原生 App 支援**: 在 `PRODUCT.md` 的 `## Platform` 設為 `ios`, `android` 或 `adaptive` 時, 設計與 `audit` / `adapt` 會切換為 Apple HIG / Material 3 與 SwiftUI / Compose / React Native / Flutter 檢查; detector, `live` mode 與 edit hook 則維持 Web 專用.
- **獨立維運指令**: `/impeccable pin <cmd>` (註冊獨立斜線捷徑), `/impeccable hooks` (開關自動檢查 hook), `/impeccable doctor` (修復 context 過期與環境飄移), `/impeccable foundations` (產品基礎設施清單檢查).

---

## 與一般做法的核心差異 (How It Differs From General Approaches)

相較於一般自然語言提示詞, 單檔 UI skill (如 Anthropic 原生約 71 行的 `frontend-design`), 甚至單純靠 LLM 逐層改版的去 AI 化 skill (如 [unslop.md](unslop.md)), `Impeccable` 在五個架構層面上做出本質區隔:

| 面向 | 一般做法 / 單檔 Skill | Impeccable |
| :--- | :--- | :--- |
| **指令結構** | 單一系統提示詞 + 模糊形容詞 ("make it pop") | 24 個專責動詞, 各自對應獨立 playbook (`skill/reference/*.md`) |
| **場景區分** | 套用單一預設審美 (常把後台做成 landing page) | 依表面 (surface) 強制分流為 4 種 Mode: `Persuade` / `Operate` / `Read` / `Experience` |
| **品質驗證** | 讓 LLM 自評文字或截圖 (模型受限自身偏見且結果飄移) | Rust + WASM **62 條確定性規則** (Zero-LLM), 串接 Edit Hook, CLI/CI 與瀏覽器擴充 |
| **設計稿還原** | AI 看圖一次寫出大檔, 材質扁平化且重寫耗盡 turn | 機械化 `comp-spec` / `comp-diff` / `font-match` + 磁碟狀態機閘門 (`build-phase`) |
| **工作迴圈** | 單次對話重新交代脈絡, 靠文字反覆猜位置微調 | `PRODUCT.md` + `DESIGN.md` 跨 session 記憶, 加上 `live` mode 瀏覽器內直接挑變體寫回原始碼 |

### 1. 從「單一審美提示」變成「設計詞彙表 + Surface Mode 分流」

- **Surface Mode 強約束**: 不以產品類型, 而以「該頁面的使用者成功定義」決定模式 (`SKILL.src.md`):
  - `Persuade` (Landing / 行銷 / 定價): 設計即產品, 強調張力與真實視覺素材
  - `Operate` (App UI / Dashboard / 設定 / 編輯器): 掃視效率, 一致性與原生預期高於張揚表現
  - `Read` (文件 / 文章 / Changelog): 閱讀結構與可讀性優先
  - `Experience` (作品集 / 展間): 作品本身在首屏領銜, 介面退居其次
- **明確的 Craft Floor 與 Refuse 清單** (`skill/reference/craft-floor.md`): 在任何 UI 改動前必讀. 除了禁止重複等大卡片網格, eyebrow kicker, `01/02/03` 裝飾編號, 漸層文字, 無脈絡硬陰影 (`4px 4px 0`) 之外, 特別要求對「模型最常漏掉的瀏覽器預設表面」(`selection`, `caret`, custom scrollbar, focus ring, `underline-offset`, `tabular-nums`) 上品牌色階.

### 2. Zero-LLM 的 Rust + WASM 確定性偵測器 (三端共用同一套核心)

- 一般 skill 最大的盲點是「叫生成 slop 的同一個 LLM 去抓 slop」. Impeccable 將規則引擎抽成純 Rust workspace (`crates/foundation`, `crates/core`, `crates/html`, `crates/browser`, `crates/detect`), 無 Node 執行期依賴, 並以 `wasm-pack` 編譯同一套邏輯到瀏覽器端 (`crates/wasm`, `browser-bundle/`).
- 內建 **62 條規則** (`crates/live/assets/antipatterns.json`: `slop` 31 條 + `quality` 31 條; `error` 46 條 + `advisory` 16 條), 透過三種入口落地:
  1. **Agent Edit Hook (`crates/hook`)**: 在 Claude Code, Cursor, Codex, Grok Build, Hermes Agent 中, 每次修改 UI 檔案後自動執行, 於 Agent 宣稱完成前直接把違規項目退回修正.
  2. **CLI & CI Gate (`crates/detect`)**: 可掃靜態原始碼 (JSX/TSX/Svelte/Vue/HTML), 本地背景圖像素對比 (偵測白字壓亮圖), 或透過 CDP 連真實 Chrome/Edge 掃渲染後 DOM 與排版 (`impeccable detect http://localhost:5173`), 支援 `--sarif` (SARIF 2.1.0), `--github` 與 `--report-unused-ignores` (抓過期 inline waiver).
  3. **Chrome Extension (`extension/`)**: 直接在任何線上網頁疊加偵測報告.

### 3. 機械化的設計稿還原管線 (`Comp Fidelity Pipeline`)

依據 `docs/COMP-FIDELITY.md`, 團隊針對模型「看得見圖片卻畫不準, 又無法客觀自評相似度」的結構缺陷, 把還原流程從純文字說教改成硬體級腳本閘門 (`crates/comp`, `crates/comp-verbs`):

1. **`comp-spec` 10x10 格線規格化**: 先在設計稿壓上 10x10 格線, 要求模型將區塊具名標記為 `plate` / `image` / `texture` / `text` / `control` / `chrome`, 產生 `.impeccable/build/spec.json`; 任何未列於 spec 的視覺物件不得出現, 並以 `organic-clip-path` 與 `buried-raster` 兩條規則封殺「用大量多邊形 `clip-path` 或 SVG 假裝插畫」與「把生成的貼圖埋在幾乎不透明漸層下」的偷懶行為.
2. **`build-phase` 磁碟狀態機**: 強制鎖定 `comps -> spec -> plates -> hero -> sections -> motion -> responsive -> review` 階段推進; `hero` 門檻未達 `comp-diff` 0.72 分 (且不得有 `missing` 區塊) 即擋下, 不允許跳關硬寫整頁.
3. **`comp-diff` 像素與結構評分**: 依區塊類型加權計算灰度 SSIM, 色彩直方圖 / Lab 色盤, 高頻細節能量與水平 band 對齊, 輸出 `match` / `drift` / `missing` / `contradicted` 判決與視覺熱圖.
4. **`font-match` 字體指紋索引**: 不再讓模型靠感覺猜字體名稱, 而是預先對 Google Fonts 約 `3,100` 個拉丁字體家族在 `48px` 與 `14px` cap height 建立約 700 KB 的形狀指紋索引 (`skill/scripts/data/font-index.json`), 直接測量設計稿裁切文字的字寬比, x-height, 筆畫粗細與襯線比來排序候選字體, 再用實際文案重新渲染對比.

### 4. `Live Mode` 與跨 Session 專案記憶體

- **持久化專案雙檔**: `/impeccable init` 產出 `PRODUCT.md` (受眾, 目的, 約束, 平台), `/impeccable document` 產出 `DESIGN.md` (token 與視覺系統), 搭配 `.impeccable/surfaces/*.md` 跨對話保存設計脈絡.
- **瀏覽器內雙向熱迭代 (`crates/live`)**: 支援 Vite, Next.js, SvelteKit, Astro, Nuxt 與純 HTML. 在頁面中選取任意 DOM 元素即可並列預覽多種設計方向變體, 確認後直接寫回對應元件原始碼 (支援手動改文案並透過 session/CLI 回寫), 避免在聊天視窗用自然語言來回猜測微調.

### 5. 模型與 Harness 意識的條件編譯 (`Multi-Harness + Model-Specific Guardrails`)

- 不同 AI 模型與 CLI 的能力邊界不同, `Impeccable` 的 `skill/SKILL.src.md` 與參考文件不是靜態一份共用, 而是由建置腳本 (`scripts/lib/transformers/providers.js`) 依目標環境編譯至 18+ 個資料夾 (`.claude/`, `.cursor/`, `.agents/`, `.agent/`, `.gemini/`, `.grok/`, `.hermes/` 等):
  - **Harness 適配**: 例如針對 Claude Code 的 skill 列表 1% context window 上限提供精簡版 `claude-code-description`, 並刻意省略 `allowed-tools` 以避開非互動模式的啟用封鎖 (`issue #736`); 針對各家 hook manifest 格式 (`.claude/settings.json`, `.cursor/hooks.json`, Hermes plugin) 個別生成設定.
  - **針對特定模型家族的壞習慣下藥**: 在 `craft-floor.md` 內以條件標籤 (如 `<codex>...</codex>`) 僅對 Codex 注入專屬禁令, 例如壓制其過度負字距 (`-0.04em`), 1px 邊框配寬軟陰影的 `ghost card`, 粗糙手繪 SVG (`feTurbulence` / `loose-sketch`), `repeating-linear-gradient` 條紋與兩軸網格背景, 以及先命名概念再自我反諷的文案習慣.

---

## 安裝與使用摘要 (Install & Usage)

```bash
# 於專案根目錄自動偵測已安裝的 AI 工具並設定 Skill 與 Design Hook
npx impeccable install

# 單獨跑確定性偵測 (無需 LLM / API key)
npx impeccable detect src/
npx impeccable detect http://localhost:5173
npx impeccable detect --sarif src/ > impeccable.sarif

# 在 Agent 內初始化並執行
/impeccable init
/impeccable critique the pricing cards
/impeccable polish settings
```

- **何時選 Impeccable vs. [unslop.md](unslop.md)**:
  - `unslop` 輕量 (純 skill + Playwright script), 強項在快速對現有網頁或 URL 截圖跑 9 層 redline 修正與產出展示用 story video.
  - `Impeccable` 屬於工程化產品級設計系統: 需要長期維護 `PRODUCT.md` / `DESIGN.md`, 想在每次編輯或 GitHub CI 擋下 AI slop, 需要在瀏覽器內做 `live` 變體實驗, 或要把視覺設計稿高保真轉譯成前端程式碼時, `Impeccable` 的覆蓋度與確定性遠高於一般 skill.
