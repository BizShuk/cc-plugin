# taste-skill (`Leonxlnx/taste-skill`)

| 欄位 | 值 |
| --- | --- |
| Source | <https://github.com/Leonxlnx/taste-skill> · <https://tasteskill.dev> |
| Type | `skill` (skill 名稱 `design-taste-frontend`, 資料夾 `skills/taste-skill/`) |
| Author / License | Leonxlnx (277 / ~300 commits 為作者本人) / MIT |
| Stats | 119,376 stars, 7,983 forks, 14 open issues + PR, last push 2026-11-08, latest release `v2.5.1` (2026-11-07) (2026-11-21 via GitHub API) |
| Noted | 2026-11-21 |
| Brief | 讓 coding agent 產出「像真公司官網」而非 AI 模板感的 landing page: 一套規則 + 素材 kit + 在 headless Chrome 實際檢查並擋住收尾的 finish gate |

## 摘要 (Summary)

- `taste-skill` 是專做`公司官網 / landing page` 的重量級 frontend skill: 381 行 `SKILL.md`, 22 份 references, 12 支 Node script, 39 MB 素材 (照片, icon, shader effect, 版型).
- 最大優勢是`可驗證的品質門檻`: 規則 (一屏一 section, 左右邊界對齊, 無圖說, 禁紫色...) 幾乎都對應到 checker id, 由 Stop hook 在真瀏覽器跑過才放行; 這是 [impeccable](impeccable.md) 與 [unslop](unslop.md) 沒有的強制層.
- 最大限制是`成本與安全`: 首頁前規定必讀約 25K words (`推論` 約 33K tokens), 每輪檢查要起 headless Chrome; 已公開但未修的 issue 包含本機 static server 可經 symlink 讀到專案外檔案, CDP debugging port 未驗證 (#202, #207).
- 適用: 對外行銷頁, 店家 / SaaS / 硬體官網. 不適用: dashboard, 內部工具, 既有 design system 的產品 UI (規則是為官網敘事寫的).
- 建議 `觀望`: 品質機制值得參考, 但等 2.5.2 安全修正 (PR #199) 合併後, 再在 scratchpad 專案試用.

## 它是什麼 (What is it)

解決的問題: agent 產生的網頁容易同質化 (置中 hero, 紫藍漸層, 卡片牆, 空泛文案). `taste-skill` 把「品味」拆成可執行的規則, 素材與檢查.

| 層 | 內容 | 位置 |
| --- | --- | --- |
| 規則 | 9 條 House rules + 依業種 (`software` / `gaming` / `tech` / `product` / `service` / `place`) 決定可用素材的對照表 | `SKILL.md` §1-3 |
| 抽籤 (Deal) | hero 版型, 配色, 字體配對, 動態 moment 由 script 隨機抽並記 history, 避免每頁長一樣 | `assets/layout/deal.mjs`, `scripts/start.mjs`, `scripts/moments.mjs` |
| 素材 kit | `starter.html` + `system.css/js`, 20 個 `fx-*` 元件, 31 個 shader effect, 66 張 Commons 照片, 929 個 icon | `assets/` |
| 文案層 | voice card, register, slop 詞檢查 | `references/copy*.md`, `scripts/copy-lint.mjs` |
| 檢查 | `consistency.mjs` (11,485 行), `check.mjs`, `perf.mjs`, `shoot.mjs` 在 1440 / 1024 / 820 / 768 / 390 / 320 寬度跑 headless Chrome | `scripts/` |
| Finish gate | `gate.mjs` 以 skill frontmatter 的 `Stop` hook 掛載, FAIL 未清就回 `{"decision":"block"}` | `SKILL.md:14-19`, `scripts/gate.mjs` |

使用流程 (`SKILL.md` §4): 讀 brief 判定業種 -> `gate.mjs --start` 武裝 gate -> 先寫完全部文案 -> 以「屏」為單位規劃 -> 抽籤 -> 依 recipe 組頁 -> 跑檢查與截圖 -> gate 到 `GATE: PASS` 或 `GATE: OPEN` (最多 3 輪) -> 依固定格式回報.

## 視角評估 (Perspectives)

| 視角 | 優勢 (Advantage) | 劣勢 (Disadvantage) | 結論 |
| --- | --- | --- | --- |
| 品質強制機制 | 規則對應 checker id, 真瀏覽器多寬度檢查, Stop hook 擋收尾 | 每輪起 Chrome, gate 單次最長 9.5 分鐘; 判斷規則有固定審美偏好 | 本 skill 的核心價值, 值得借鏡 |
| 觸發範圍與適用情境 | description 列出多種業種與框架, 觸發明確 | 904 字元的 description 涵蓋「review / polish 任何 landing page」, 容易在非官網任務被載入 | 只適合官網類任務 |
| Context / Token 成本 | references 依素材分檔, 只讀用到的 | 首頁前規定讀完 5 份 reference + `SKILL.md` 約 25K words | 成本高, 不適合當預設 skill |
| 安全邊界 | 零 npm 依賴, 不需 API key; hook 有防無限阻擋的上限 | 本機 server symlink 越界, CDP port 未驗證, 寫檔跟隨 symlink, hook 期間的網路存取未寫進文件 | 目前 `不建議`在含敏感檔案的專案使用 |
| 相依, 可攜性與維護 | Node built-in only, Claude Code / Codex / Antigravity 三種 manifest, CI 驗 frontmatter | 幾乎單人維護; 安全修正 PR 寫明「Not for merging yet」 | 活躍但 bus factor 低 |

### 品質強制機制 (Gate)

- 證據: `SKILL.md:298-363` checklist 每行都標出 `consistency.mjs` / `check.mjs` / `perf.mjs` 的檢查 id (例 `[hero-screen, hero-fill, fold-action]`); `gate.mjs` 的 hook 模式回 `{"decision":"block","reason":...}` (claude, codex) 或 `continue` (antigravity) (`scripts/gate.mjs:71-72,111`); 誤判可在 `.taste-skill/waivers.json` 豁免, reason 須 12 字元以上 (`scripts/gate.mjs:569`).
- 優勢: 「看起來不像 AI 做的」被轉成可重跑的檢查, 例如 section 高度, 邊界對齊, tap target 44px, 對比度, slop 詞. agent 不能只宣稱完成. Hook 設計避免卡死: 同一 build 超過「已跑輪數 + 4」次 block, 或資料夾 6 小時內 30 次後一律放行, `--done` 後也放行 (`scripts/gate.mjs:73-79`).
- 劣勢: 每次檢查起一個 headless Chrome (WebGL 走 CPU), gate 單次上限 9.5 分鐘 (`scripts/gate.mjs:80`); 規則帶強烈的固定審美 (例如禁紫色, 一屏一 section, 店家頁禁 shader), 品牌需求不同時要逐條 waive.
- 結論: 最值得參考的部分是「規則 -> checker id -> hook 擋收尾」這條鏈, 可以拿來對照本 workspace 的 frontend 類 skill.

### 觸發範圍與適用情境

- 證據: description 904 字元 (上限 1024), 觸發句為 build / design / redesign / restyle / polish / review 任何 landing page, company, product, shop, SaaS, game site, 以及使用者說頁面「看起來很 AI」 (`SKILL.md:3-13`). 業種表 (`SKILL.md:66-73`) 只涵蓋對外官網.
- 優勢: 業種分流清楚, 同一套規則能套到咖啡店 (以照片為主, 禁 effect) 與 dev tool (live UI mock + 一個 shader).
- 劣勢: 對 dashboard, 後台, 互動工具沒有對應的頁型; 「review 任何頁面」的觸發詞讓它可能在只想小修 CSS 時被載入, 並帶入整套規則 (`推論`, 未實測觸發率).
- 結論: 限定在官網類任務使用; 若安裝, 考慮在 project 層而非全域安裝.

### Context / Token 成本

- 證據: `SKILL.md:156` 要求第一頁前完整讀完 `copy.md`, `copy-voice.md`, `system.md`, `layout-v4.md`, `sections.md`; 這 5 份加 `SKILL.md` 共 25,183 words, 全部 22 份 references 共 67,961 words (2026-11-21 以 `wc -w` 計, 未實測 token).
- 優勢: references 依素材家族分檔, `SKILL.md` §3 表格指明每種素材對應哪份文件, 其餘按需讀取.
- 劣勢: 必讀量 `推論` 約 33K tokens (以 1.3 tokens/word 估), 加上每輪檢查輸出, 一個頁面會吃掉大量 context.
- 結論: 適合「一次做一個重要頁面」, 不適合高頻小改.

### 安全邊界

- 證據:
    - Stop hook 從多個候選路徑 `exec node gate.mjs` (`SKILL.md:18`); 未武裝 build 時只對判定為 tasteskill 頁面的檔案 block 一次並提示 (`scripts/gate.mjs:588-594`).
    - 網路存取: Wikimedia Commons API 與 `upload.wikimedia.org` (`scripts/photos.mjs:667,787`), `cdn.jsdelivr.net` 的 pinned icon 套件 (`scripts/icons.mjs:65-72`), Google Fonts (`scripts/start.mjs`).
    - 開放中的安全 issue: #202 本機 static server 以字串比對路徑, `effects.mjs` / `moments.mjs` / `perf.mjs` 用 `startsWith(root)` 無分隔符, symlink 可讀專案外檔案; Chrome 以 `--remote-debugging-port=0` 驅動, 本機其他程序可連上 (`scripts/shoot.mjs:321`). #207 輸出檔以 `writeFileSync` / `mkdirSync({recursive})` 寫固定檔名, 會跟隨 repo 內的 symlink. #211 README 未說明 Stop hook 期間 Chrome 會載入頁面引用的第三方資源.
    - 修正狀態: PR #199 (2.5.2) 含 `photos.mjs` 安全修正, 描述寫明「Not for merging yet」(2026-11-21 via GitHub API).
- 優勢: 零 npm 依賴 (所有 import 皆 `node:*`), 不需 API key, README 明載「installs nothing」(`README.md:86,137`); 每個 check 有 watchdog, gate 被殺時連同 Chrome 一起收掉 (`scripts/gate.mjs:79-80`).
- 劣勢: 上列 issue 在 clone 來的不可信 repo 中執行時風險最高; hook 會在每次 Stop 執行 Node script, 屬於自動執行的程式碼.
- 結論: 安全修正合併前, 只在不含機密的獨立資料夾使用.

### 相依, 可攜性與維護

- 證據: 需求 Node 22+, Chrome / Edge / Chromium, ffmpeg 選用 (`README.md:86`); repo 同時提供 `.claude-plugin/marketplace.json`, `.codex-plugin/plugin.json`, `.agents/plugins/marketplace.json`; `install.sh` 以 symlink 裝到 `~/.claude/skills` 或 `~/.codex/skills` (`install.sh:38-50`); CI 驗 frontmatter, UTF-8 與 script 語法 (`.github/workflows/check-skills.yml`). 2.0.0 (2026-10-13) 到 2.5.1 (2026-11-07) 共 9 個 release; contributors 前五名為 Leonxlnx 277, cursoragent 9, HugoSpy 5, bougieL 3, flexi767 2.
- 優勢: 三種 agent 的安裝管道齊全, `git pull` 即更新; release 節奏約每週一版.
- 劣勢: 幾乎單人維護 (bus factor 低); #200 回報 2.1.1 起缺 `.claude-plugin/plugin.json`, claude.ai 無法同步; `SKILL.md` 的 `name` (`design-taste-frontend`) 與資料夾名 (`taste-skill`) 不一致, 不符本 workspace 的 [skill-frontmatter](../../../../.claude/skills/skill-frontmatter/SKILL.md) 規範 (install.sh 以 frontmatter name 建 link, 所以安裝後名稱正確).
- 結論: 維護活躍, 但重大變更與安全修正都卡在作者一人.

## 細節 (Details)

安裝 (三選一):

```bash
npx skills add https://github.com/Leonxlnx/taste-skill --skill "design-taste-frontend"
# Claude Code 內
/plugin marketplace add Leonxlnx/taste-skill
/plugin install taste-skill@taste-skill
# 或 clone 後
./install.sh            # ~/.claude/skills; --codex / --all-agents
```

規模 (2026-11-21 shallow clone, 未實測執行):

| 項目 | 數量 / 大小 |
| --- | --- |
| Repo 全部 | 1,663 files, 69 MB |
| `skills/taste-skill/` | 1,631 files, 39 MB |
| `assets/photos` | 66 張照片, 12 MB |
| `assets/icons` | 932 files (929 個 demo SVG), 5.8 MB |
| `assets/objects` / `scenes` / `moments` | 6.3 MB / 3.1 MB / 2.4 MB |
| `assets/effects` / `components` | 31 個 shader effect / 20 個 `fx-*` 元件 |
| Scripts | 12 支 `.mjs`, 共 17,319 行 (`consistency.mjs` 佔 11,485) |

執行期副產物: 專案內 `.taste-skill/` (history, `gate.json`, `gate-blocks.json`, `waivers.json`, Chrome profile), `shots/` 截圖, `names.json`.

本機現況: `未安裝` (`~/.claude/skills`, `~/.codex/skills`, plugin cache 皆無).

注意事項:

- 與其他 Stop hook (例如 Claude Code `/goal`) 並用時, 兩個 hook 都可能 block; gate 自身有上限, 但兩者疊加的行為未驗證 (`推論`).
- House rule 9 規定地址不寫門牌與郵遞區號, 聯絡方式用 `.example` 網域; 真實客戶頁面上線前要手動換回.

## 相關參考 (References)

- [impeccable](impeccable.md) - 同樣反 AI slop, 偏設計詞彙與反模式偵測引擎, 支援 agent 數量更多, 沒有擋收尾的 gate
- [unslop](unslop.md) - 改寫既有 vibecoded UI, 輕量, 涵蓋 dashboard 與 app
- [hallmark](hallmark.md) - 同樣做新頁面, 以 50 個主題抽牌與零依賴 linter 驗證, edit hook 只建議不阻擋
- [De-AI spec](../../../../docs/specs/2026-11-18-de-ai.md) - 本 workspace 整理的 de-AI 設計概念
- `taste-skill-classics` (`README.md:166`) - 舊版 v1 skill 的分支 repo
- Issues: [#202](https://github.com/Leonxlnx/taste-skill/issues/202) server / CDP, [#207](https://github.com/Leonxlnx/taste-skill/issues/207) symlink 寫檔, [#211](https://github.com/Leonxlnx/taste-skill/issues/211) hook 期間網路存取, [#199](https://github.com/Leonxlnx/taste-skill/pull/199) 2.5.2 安全修正
