# Hallmark (`nutlope/hallmark`)

| 欄位 | 值 |
| --- | --- |
| Source | <https://github.com/nutlope/hallmark> · <https://www.usehallmark.com> |
| Type | `skill` (skill 名稱 `hallmark`, 資料夾 `skills/hallmark/`, 另附 Claude Code plugin manifest) |
| Author / License | Nutlope / Together AI 出品 (`README.md:7`); 主要 commit 者 Luffixos 312 次 / MIT; 照片庫各自保留 CC0 / PD / CC BY / CC BY-SA |
| Stats | 58,655 stars, 3,197 forks, 51 open issues + PR (其中 22 個 PR), last push 2026-11-11, 無 GitHub release, `version: 2.0.0` (`SKILL.md:5`, `package.json`) (2026-11-21 via GitHub API) |
| Noted | 2026-11-21 |
| Brief | 讓 coding agent 從零做出`不像 AI 生成`的網頁: 一組不可豁免的 Floor 規則, 50 個主題材質抽牌, 零依賴的 linter 在真瀏覽器 8 種寬度驗證, 再由 fresh-context subagent 做最後審查 |

## 摘要 (Summary)

- `Hallmark` 是專做`新頁面`的反 AI slop skill: 154 行 `SKILL.md` + 24 份 references + 50 個主題, 21 種頁型, 4 種 genre + 10 支零依賴 Node script, 安裝後約 7.1 MB.
- 最大優勢是`三層 gate + 機械驗證`: linter 的 105 個 gate id 分為 Floor (47, 不可豁免) / Reflex (57, 需具名理由才可 waive) / Finish (不阻擋), 由 4,426 行的 `sloplint.mjs` 靜態檢查, `--render` 在 320-1920 共 8 種寬度實測, 再交給不帶對話紀錄的 subagent 做 hero 與 finish 兩輪審查.
- 最大限制是`範圍與成本`: description 明訂`不對既有設計動手`, 除非使用者點名 `/hallmark redesign`; 一次頁面建置的必讀量上限約 22.9K words (`推論` 約 30K tokens); 第三方 ablation (#74, 舊版) 顯示未點名時自動觸發 0/12.
- 與同目錄筆記的差異: 比 [taste-skill](taste-skill.md) 少了 Stop hook 強制擋收尾 (Hallmark 的 edit hook 只建議不阻擋); 比 [impeccable](impeccable.md) 的 Rust 偵測器輕, 但多了主題抽牌與照片庫; 比 [unslop](unslop.md) 偏"新建"而非"改寫既有 UI".
- 建議 `觀望`: 品質鏈 (Floor -> linter -> fresh-context review) 設計扎實且無 API key 依賴, 但本 workspace 已有 frontend-design 與 ui-ux-pro-max 類 skill; 需要做對外 landing page 時, 再在 scratchpad 專案以 `/hallmark <prompt>` 點名試用.

## 它是什麼 (What is it)

解決的問題: 模型做網頁會落入固定反射 (每個標題上方的 eyebrow, 追蹤大寫 label, 標題中一個斜體字, 紫藍漸層, 卡片套卡片, 捏造的 logo 牆與數據, 每頁結尾 email 訂閱). Hallmark 把這些拆成可檢查的 gate, 並用"抽主題"讓連續兩次建置不會收斂成同一張臉.

| 步驟 | 內容 | 位置 |
| --- | --- | --- |
| 0. Pre-flight | 讀專案的 `design.md`, 字體, 色票, framework, 品牌規格; 結果快取在 `.hallmark/preflight.json` | `SKILL.md:54-58` |
| 1. Ask once | 一次問 4 件事: Material (截圖, 照片, 真實數字), Audience, Use, Tone; headless / batch 時不問, 直接推論並寫進 stamp | `SKILL.md:60-66` |
| 2. Shape | 依訪客的 job (訂閱 / 預約 / 購買...) 讀 `anatomy.md`, 由 `theme.mjs draw` 抽頁型, hero 姿態, nav, footer | `SKILL.md:72-84` |
| 3. Direction | 使用者點名主題 -> 用之; 說 `custom` -> 改讀 `direction.md` 自由設計; 其他 -> `theme.mjs deal` 從 50 個主題抽 3 個再挑 1 | `SKILL.md:86-100` |
| 4-5. Load + Picks | 只載入規定的 references; 輸出一次 Picks block 作為使用者改向的窗口 | `SKILL.md:102-127` |
| 6. Build | scaffold 產生 `tokens.css` / `base.css`; 先做 hero 並單獨 lint + review, 通過後才做其餘 section; CSS 第一行寫 stamp | `SKILL.md:129-140` |
| 7. Sweep | `sloplint --render` 修到 0 FAIL, 走 judged gate, finish review, 最後列出每個 slot 該放什麼真實素材 | `SKILL.md:142-150` |

另有兩個`需點名`的 verb: `/hallmark redesign <path>` (保留內容與 wiring, 以 `sloplint --before` / `--kept` 確認連結, 欄位, handler 沒掉) 與 `/hallmark audit <path>` (只報告) (`SKILL.md:50`, `README.md:66-78`).

## 視角評估 (Perspectives)

| 視角 | 優勢 (Advantage) | 劣勢 (Disadvantage) | 結論 |
| --- | --- | --- | --- |
| 品質強制機制 | 三層 gate 與 linter `TIERS` 由 CI 交叉驗證; 8 種寬度實測; fresh-context subagent 審查 | edit hook 只建議不阻擋, 且 plugin 安裝不會自動掛上; ablation 顯示部分規則曾適得其反 | 核心價值, 機械層比 taste-skill 弱, 比 unslop 強 |
| 觸發準確度與適用範圍 | description 252 字元, 明確排除既有設計, 不會在小改 CSS 時誤入 | 自動觸發率低 (#74: 0/12, 舊版 + Haiku); 以官網 / landing page 為主 | 當作"點名才用"的 skill |
| Context / Token 成本 | references 分檔, `gates.md` § 5-6 延後到 Step 7; review prompt 由 subagent 讀, 主線程不讀 | 必讀上限約 22.9K words + 3 張主題圖 + 多輪截圖 | 適合一次做一個重要頁面 |
| 安全邊界與外部呼叫 | 零 npm 依賴; 外部網路只有 opt-in 的 Together 生圖與限定 `upload.wikimedia.org` 的照片下載; hook 寫入用 atomic rename | `install-hook.mjs --global` 會改 `~/.claude/settings.json`; 生圖會把 prompt 送到 Together API | 風險低, 生圖與 global hook 需刻意開啟 |
| 維護活躍度與可攜性 | 30 天內 75 commits, CI 跑 refcheck / theme validate / tests; 支援 Claude Code, Codex, Cursor, OpenCode | 無 release tag; 主要由一人 (Luffixos) 提交; 文件內部矛盾 (#62) 與 install 後斷鏈 (#76) 曾被回報 | 活躍, 但版本釘選只能靠 commit SHA |

### 品質強制機制 (Gate)

- 證據:
    - Floor 16 條寫在 `SKILL.md:16-35`, 每條附 gate 編號, 例如 `transition: all` 禁用 (10a), eyebrow 禁用 (54, 67), hero 在 1280x800 不捲動可見 (44b), 不得捏造數據 / 推薦 / logo 牆 (46a).
    - `sloplint.mjs:801-821` 的 `TIERS` map: floor 47, reflex 57, ledger 1 (2026-11-21 以 grep 計數); `refcheck.mjs` 在 CI 比對 `gates.md` 與 `TIERS`, 不一致即失敗 (`AGENTS.md` Layout 節, `.github/workflows/verify.yml`).
    - Render tier 寬度 `[320, 375, 414, 768, 1024, 1280, 1440, 1920]` (`sloplint-render.mjs:6`), hero 階段只跑 375 與 1280 (`sloplint-render.mjs:170`).
    - Fresh-context review: "The thread that built the page is its weakest reviewer", hero review 與 finish review 由無對話紀錄的 subagent 執行, 主線程不讀 review prompt; 依第一個字 `recapture` / `rebuild` / `fix` / `ship` 行動 (`references/gates.md:237-243`).
    - Edit-time hook: `lint-hook.mjs` 只檢查 CSS 有 `/* Hallmark` stamp 的檔案或 `.hallmark/` 專案內的元件, 跑 `--tier floor`, 以 `additionalContext` 回傳最多 10 條 FAIL, 永不阻擋 (`lint-hook.mjs:24-46`).
- 優勢: "不像 AI 做的"被轉成可重跑的檢查, 且 gate 表與 linter 由 CI 互相鎖定, 不會文件說一套程式做一套. Self-review 問題用 fresh-context subagent 處理, 並要求把審查結果寫進 stamp (`SKILL.md:148` "A fabricated pass is itself slop").
- 劣勢: hook 只建議不阻擋, 收尾仍靠模型自律; plugin 安裝不會掛 hook, 需另跑 `install-hook.mjs` (#131, ROADMAP Next). 第三方 ablation (#74, 2026-09-10, 舊版, 100 次) 顯示: 純黑白色使用率 100% -> 50% (有效), 但 `transition: all` 35% -> 57% (反向, 因模型用 skill 自己的 duration token 寫出 `transition: all var(--dur-fast)`); 現行版已把 `transition: all` 移進 Floor 第 2 條 (`SKILL.md:21`), 新版效果`未知`.
- 結論: 機械驗證完整度介於 [unslop](unslop.md) 與 [taste-skill](taste-skill.md) 之間; 若要強制, 需自行掛 hook.

### 觸發準確度與適用範圍

- 證據:
    - description: "Use when the user asks to build or design something new, or names Hallmark. Never use it on an existing design unless the user names Hallmark" (`SKILL.md:3`, 252 字元).
    - #74 實測 (`claude-haiku-4-5`, 舊版 description): 只描述工作的 prompt 自動觸發 0/12, 點名後 40/40.
    - 範圍分流: 單一元件走 component scope, dashboard / admin / settings 走 app screen (`SKILL.md:70`); ROADMAP 寫明 app screen 的專屬檢查屬 Next (`ROADMAP.md`).
- 優勢: 明確排除"既有設計", 避免在修 bug 或小改 CSS 時整套規則被載入; 對比 taste-skill 904 字元且涵蓋 review / polish 的 description, 誤觸發面小.
- 劣勢: 自動觸發率低 (`推論`: 新版 description 更窄, 未點名時大概率不會啟動); 主題與頁型 (21 種 macrostructure) 都為對外頁面設計, dashboard 類只有基本分流.
- 結論: 以 `/hallmark <prompt>` 明確點名使用, 不期待自動觸發.

### Context / Token 成本

- 證據: Step 4 規定頁面建置要載入 `SKILL.md` + `layout.md` + `anatomy.md` + `imagery.md` + `craft.md` + 1 份 genre + 1 份主題 + `gates.md` § 1-4, Step 6 再讀 `patterns.md` (`SKILL.md:104,132`). 以 `wc -w` 計: 2,889 + 2,996 + 1,268 + 2,516 + 3,301 + 約 250 + 約 1,200 + 7,285 (整份 `gates.md`, 為上限) + 1,284 = 約 22.9K words; 全部 references 共 100,840 words (2026-11-21, 未實測 token).
- 優勢: 分檔按需載入, `SKILL.md:108` 明寫"loading more is the largest avoidable cost"; review prompt 交給 subagent 讀, 主線程不負擔; 有 `--fast` 模式給 prototype (`SKILL.md:154`).
- 劣勢: 必讀量 `推論` 約 30K tokens (以 1.3 tokens/word 估), 加上 deal 時看 3 張主題圖, render tier 截圖與兩輪 review; ROADMAP Now 自承"Measure a build"尚未完成, 官方沒有單頁成本數字.
- 結論: 與 taste-skill (必讀約 25K words) 同量級; 適合重要頁面, 不適合高頻小改.

### 安全邊界與外部呼叫

- 證據:
    - 依賴: 所有 script 只 import `node:*` 模組; puppeteer-core 為選用, 從 `HALLMARK_PUPPETEER`, cwd, `npm root -g` 尋找 (`capture.mjs:42-43`, `sloplint-render.mjs:113-118`); README 要求另裝 `npm i -g puppeteer-core` 與 Chrome 才有 render tier (`README.md:91`).
    - 網路: `imagegen.mjs:6,60-71` 以 `TOGETHER_API_KEY` 呼叫 `https://api.together.xyz/v1/images/generations` (預設 `FLUX.1.1-pro`), 未設 key 時退回 CSS / placeholder; `theme.mjs:413-415` 只允許下載 `https://upload.wikimedia.org/` 的照片; render tier 會在頁面內 fetch 頁面引用的 stylesheet (`sloplint-render.mjs:159`).
    - Hook 寫入: `install-hook.mjs` 預設寫專案 `.claude/settings.json`, 加 `--global` 才寫 `~/.claude/settings.json`; 解析失敗即拒寫, 以 tmp + rename 原子寫入, 對 symlink 的 settings 寫到目標檔 (`install-hook.mjs:14-16,23-29,56`).
    - `capture.mjs --cdp <port>` 會連上本機已開 DevTools port 的 Electron / WebView 應用擷取畫面 (`capture.mjs:7-10,60`).
    - Codex 的 `workspace-write` sandbox 預設斷網, 生圖需明確同意 (`references/flow.md:114`).
- 優勢: 零 npm 依賴, 無必需 API key, 外部網路皆 opt-in 或白名單; 照片只用 Wikimedia Commons 並在頁面標註授權, 避開禁止 AI 挑選的 stock site (`README.md:114`).
- 劣勢: "Powered by Together AI"寫在 `SKILL.md:14` 與 manifest, 生圖路徑綁 Together (`推論`: 屬推廣導向的預設); 依本 repo `AGENTS.md` 檔案職責表, `~/.claude/settings.json` 由 `scripts/run.sh` 軟連結到 `plugins/default/config/settings/claude.json`, 因此執行 `install-hook.mjs --global` 會直接改寫 repo 內受版控的設定檔. 查過的來源 (script grep `fetch` / `child_process`, open issues) 未見類似 taste-skill #202 的路徑越界 issue.
- 結論: 風險低; 只用 project 層 hook, 不要對本 workspace 跑 `--global`.

### 維護活躍度與可攜性

- 證據:
    - 2026-10-21 之後 75 commits, 最後 commit 2026-11-11 (`ci: run every test file in one step`); contributors: Luffixos 312, Nutlope 12, rudycelekli 10, 其餘 ≤5 (2026-11-21 via GitHub API).
    - 無 release 與 tag (`releases/latest` 404); 版本只寫在 `SKILL.md`, `package.json`, plugin manifest (皆 2.0.0).
    - CI 跑 `refcheck.mjs`, 主題 validate, `index.json` 是否最新, 全部 `*.test.mjs` (`.github/workflows/verify.yml`); browser test 在 CI 一律 skip (`AGENTS.md` Commands 節).
    - Harness 差異表: Claude Code 有 edit hook 與 subagent; Codex 無 hook 也無 subagent, review 改為 in-thread; OpenCode 有 `@general` subagent (`references/flow.md:103-115`).
    - 已回報問題: #76 install 後 `SKILL.md` 的 `../../` 連結失效 (2026-11-21 grep 現行 `SKILL.md` 與 references 已無 `../../`); #62 依 `SKILL.md` 逐字執行時 nav 規則有三個互相矛盾的來源; #131 plugin 不會掛 hook.
- 優勢: 活躍, 有 CI 保護 gate 與文件一致性; 一份 skill 支援 4 種 harness, 差異集中在一張表.
- 劣勢: 無 release 可釘版本, 主要由單一 contributor 提交 (bus factor 低); 在 Codex 上失去 hook 與 fresh-context review 兩個關鍵機制.
- 結論: 試用時以 commit SHA 記錄版本; 以 Claude Code 為主要 harness.

## 細節 (Details)

安裝 (擇一):

```bash
# 全域, 所有 harness
npx skills add nutlope/hallmark -g
# 只裝 Claude Code
npx skills add nutlope/hallmark -g --agent claude-code
# Claude Code plugin (以 /hallmark:hallmark 叫用)
/plugin marketplace add nutlope/hallmark
/plugin install hallmark@hallmark
# 從 clone 安裝並掛專案層 edit hook
node skills/hallmark/scripts/install.mjs
node skills/hallmark/scripts/install-hook.mjs          # 寫 ./.claude/settings.json
node skills/hallmark/scripts/install-hook.mjs --remove
```

使用:

```bash
/hallmark a booking page for a mobile bike repair service
/hallmark custom <prompt>          # 不抽主題, 從 direction.md 自由設計
/hallmark redesign ./src/app       # 保留內容與 wiring
/hallmark audit ./index.html       # 只報告
node <skill-dir>/scripts/sloplint.mjs <file-or-dir> --render   # 單獨 lint 任何頁面
```

規模 (2026-11-21, clone commit `04977e1`, `wc` / `du` 計):

- `skills/hallmark/` 7.1 MB (references 5.9 MB, 其中 50 個主題 `.md` + 50 張 `.jpg` 共 5.6 MB); 照片庫 `photos.json` 755 張, 只存 URL 不存檔.
- 21 種 macrostructure, 4 種 genre (`atmospheric` / `editorial` / `modern-minimal` / `playful`), 9 個 sloplint fixture, `site/made/` 21 個範例頁.
- `eval/` 為 repo 內部 benchmark (Tier B 純 API pack, Tier A 驅動 `claude -p "/hallmark <brief>"`), 不隨 skill 發佈; 結果目錄 `runs/` 未入版控, 公開數字`未知`.

本機現況: `~/.claude/skills/hallmark` 與 `~/.agents/skills/hallmark` 不存在 (2026-11-21 `ls`), `未安裝`.

注意事項:

- 依賴: Node 18+; render tier 需 `puppeteer-core` (全域) 與 Chrome / Chromium, 缺少時只跑靜態檢查 (`README.md:91`).
- 每頁 CSS 第一行是 Hallmark stamp, 是 edit hook 判斷"這是 Hallmark 頁面"的依據; 移除 stamp 後 hook 不再檢查.
- 建置會在專案根目錄寫 `.hallmark/` (`preflight.json`, `log.json`), 主題輪替靠 `log.json` 避開最近用過的組合.

## 相關參考 (References)

- [taste-skill](taste-skill.md) - 同為官網類反 AI slop skill, 有 Stop hook 強制擋收尾, 但有未修的本機 server 安全 issue
- [impeccable](impeccable.md) - Rust + WASM 確定性偵測器與 24 個設計 verb, 偏既有產品的長期設計系統; Hallmark 在 README `Inspired by` 列出它
- [unslop](unslop.md) - 輕量改寫既有 vibecoded UI, 涵蓋 dashboard 與 app
- [De-AI spec](../../../../docs/specs/2026-11-18-de-ai.md) - 本 workspace 整理的 de-AI 設計概念
- Issues: [#74](https://github.com/nutlope/hallmark/issues/74) 第三方 ablation, [#62](https://github.com/nutlope/hallmark/issues/62) 文件矛盾, [#76](https://github.com/nutlope/hallmark/issues/76) install 後斷鏈, [#131](https://github.com/nutlope/hallmark/issues/131) plugin 未掛 hook
- 主題預覽: <https://www.usehallmark.com>
