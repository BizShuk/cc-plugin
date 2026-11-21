# web-design-guidelines

| 欄位 | 值 |
| --- | --- |
| Source | <https://github.com/vercel-labs/agent-skills/tree/main/skills/web-design-guidelines> (規則上游 <https://github.com/vercel-labs/web-interface-guidelines>) |
| Type | `skill` |
| Author / License | Vercel (`metadata.author: vercel`) / skill repo `未知` (無 LICENSE 檔, GitHub license API 回 404); 規則上游 MIT |
| Stats | `vercel-labs/agent-skills` 34,621 stars, 3,011 forks, last push 2026-11-20, latest release `agent-skills-5a37a63` (2026-11-20); skill 版本 `1.1.0`, 最後變更 2026-10-29; 上游 `web-interface-guidelines` 1,229 stars, last push 2026-10-05 (2026-11-21 via GitHub API) |
| Noted | 2026-11-21 |
| Brief | 給 coding agent 的靜態 UI 審查清單: 讀前端原始碼, 依 17 類規則輸出 `file:line - issue` 格式的問題列表 |

## 摘要 (Summary)

- 它是什麼: 兩個 Markdown 檔 (`SKILL.md` 51 行 + `references/command.md` 191 行) 組成的 review skill, 沒有任何 script, 由模型逐條對照規則審查使用者指定的檔案.
- 最大優勢: 規則內容具體到屬性與 API 名稱 (`aria-live`, `prefers-reduced-motion`, `Intl.DateTimeFormat`, `overscroll-behavior`), 輸出格式固定為 VS Code 可點擊的 `file:line`; v1.1.0 起規則改為內含並釘住上游 commit 與 SHA-256, 不再每次從 `main` 抓遠端內容.
- 最大限制: 純靜態閱讀, 沒有 axe / Lighthouse 之類的實測工具, 也沒有色彩對比度 (WCAG contrast ratio) 規則; 規則用語偏 React / Next.js / Tailwind (`htmlFor`, `<Link>`, `focus-visible:ring-*`, `line-clamp-*`, `suppressHydrationWarning`), 其他框架要自行對映.
- 與同類筆記差異: [impeccable](ai-unslop/impeccable.md), [unslop](ai-unslop/unslop.md), [taste-skill](ai-unslop/taste-skill.md) 處理的是視覺美感與 AI slop 風格; 本 skill 不談美感, 只檢查 a11y, 表單, 效能, i18n, hydration 等工程正確性.
- 建議 `採用`: 對 React / Next.js / Tailwind 前端做 review 時成本極低 (兩個文字檔, 不執行程式), 可作為 merge 前的 checklist; 需要對比度或執行期驗證時另配實測工具.

## 它是什麼 (What is it)

解決的問題: agent 寫完 UI 後缺少一份固定, 可重複套用的 Web 介面品質清單. 本 skill 把 Vercel 的 Web Interface Guidelines 包成 agentskills.io 格式, 讓 agent 在使用者說 `review my UI` 時自動套用.

| 元件 | 角色 |
| --- | --- |
| `skills/web-design-guidelines/SKILL.md` | frontmatter (`name`, `description`, `metadata.version: "1.1.0"`, `argument-hint: <file-or-pattern>`) 與流程: 讀規則, 讀目標檔, 逐條檢查, 依格式輸出 |
| `skills/web-design-guidelines/references/command.md` | 規則本體: 17 個 `###` 類別, 104 條 bullet (其中 `Anti-patterns` 14 條是前文規則的重述), 以及 `Output Format` 範例 |
| `skills/web-design-guidelines.zip` | 同內容的打包檔, 供 claude.ai 上傳; 已比對與資料夾內兩檔逐字相同 |
| 上游 `vercel-labs/web-interface-guidelines` | 規則來源 repo, 另附 `install.sh` 把 `command.md` 裝成各工具的 slash command |

使用方式: 使用者給檔案或 glob -> agent 讀 `references/command.md` -> 讀目標檔 -> 依檔案分組輸出 `path:line - 問題`, 全數通過的檔案標 `✓ pass`; 沒給檔案時先詢問 (`SKILL.md:51`).

## 視角評估 (Perspectives)

| 視角 | 優勢 (Advantage) | 劣勢 (Disadvantage) | 結論 |
| --- | --- | --- | --- |
| 觸發準確度 | `description` 列出 5 個明確觸發片語 | `audit design` / `check my site against best practices` 範圍寬, 與 impeccable `audit` 重疊 | 同時裝 impeccable 時要手動指名 |
| 安全邊界 | 無 script, 無工具宣告; v1.1.0 改為內含規則並釘 commit + SHA-256 | SHA-256 只供維護者更新時核對, 執行期不驗證; 被審查的檔案內容仍是不可信輸入 | 風險低, 適合直接安裝 |
| 規則覆蓋與輸出品質 | 17 類, 規則具體到屬性名; 輸出格式固定, 可點擊 | 無對比度數值規則, 無實測; 規則有框架偏向; 少數規則彼此張力 | 作為 checklist 足夠, 不能取代 a11y 實測 |
| 跨 Agent 可攜性 | agentskills.io 格式 + `npx skills` + zip; 上游 `install.sh` 支援 6 種工具 | slash command 版用 `$ARGUMENTS`, skill 版要靠 `SKILL.md` 指示代換 | 可攜性高 |
| 維護活躍度 | 母 repo 活躍; 上游規則 2026-08, 2026-10 仍有新增; 已釘版本等於上游 `main` | 釘版本後更新需人工四步驟; skill repo 無 LICENSE | 目前同步, 長期需留意漂移 |

### 觸發準確度

- 證據: `skills/web-design-guidelines/SKILL.md:3` 觸發片語為 `review my UI`, `check accessibility`, `audit design`, `review UX`, `check my site against best practices`; `SKILL.md:7` `argument-hint: <file-or-pattern>`; `SKILL.md:51` 無檔案時要求詢問.
- 優勢: 片語都是 UI review 語境, 不會在一般 coding 時誤觸; 需要目標檔案, 不會自行掃整個 repo.
- 劣勢: `audit design` 與 [impeccable](ai-unslop/impeccable.md) 的 `audit` / `critique` 指令語意重疊; `check my site` 可能被理解成檢查線上網站, 但本 skill 只讀原始碼 (`SKILL.md:17`, `Read the specified files`). 兩者同裝時誰被選中屬 `推論`, 未實測.
- 結論: 單獨安裝時觸發清楚; 與其他 UI skill 並存時以 skill 名稱明確呼叫.

### 安全邊界

- 證據: 資料夾只有兩個 `.md` 檔 (`find skills/web-design-guidelines -type f`), 無 script, 無 `allowed-tools`. 舊版 `1.0.0` (commit `ba46938`, `SKILL.md:16-29`) 要求每次審查前以 WebFetch 抓 `raw.githubusercontent.com/vercel-labs/web-interface-guidelines/main/command.md`. commit `f766b64` (2026-10-29, `fix(web-design-guidelines): bundle pinned guidelines instead of fetching from main`) 改為內含 `references/command.md`, 並在 `SKILL.md:23-29` 寫明 `Do not fetch rules from the network`, 釘住上游 commit `434b7f91364665f2f733b310ec54809bf8f37937` 與 SHA-256 `d246b026...e0234`.
- 實測 (2026-11-21): `shasum -a 256 references/command.md` 與從上游該 commit 下載的檔案雜湊相同, `diff` 無差異.
- 優勢: 移除了「每次執行都讀遠端 `main` 分支內容當指令」的供應鏈 / prompt injection 面; 安裝後完全離線.
- 劣勢: SHA-256 只出現在 `SKILL.md:34-41` 的維護者更新流程, agent 執行時不核對檔案雜湊; 審查對象是使用者的原始碼, 其中的註解或字串仍可能夾帶指令, skill 沒有額外防護 (`推論`, 一般 review skill 共通風險).
- 結論: 安全面在 v1.1.0 已收斂到純本機文字讀取; 仍在用 `1.0.0` 或上游 `install.sh` 舊式安裝者應更新.

### 規則覆蓋與輸出品質

- 證據: `references/command.md:14-167` 共 17 類: Accessibility, Focus States, Forms, Animation, Typography, Content Handling, Images, Performance, Navigation & State, Touch & Interaction, Safe Areas & Layout, Dark Mode & Theming, Locale & i18n, Hydration Safety, Hover & Interactive States, Content & Copy, Anti-patterns. bullet 數以 `awk` 計得 104 條, 其中 `Anti-patterns` (`command.md:154-167`) 14 條重述前文. 輸出格式見 `command.md:169-191`.
- 優勢:
  - 規則可直接對應程式碼: 如 `command.md:55` `Never transition: all`, `command.md:80` `<img>` 要明確 `width` / `height`, `command.md:126-127` 日期與數字用 `Intl.*`, `command.md:107` 桌面頁面 rubber-band 用 `overscroll-behavior: none` 並保留 coarse pointer 的 pull-to-refresh.
  - 輸出為 `file:line - issue` 且 `Skip explanation unless fix non-obvious. No preamble` (`command.md:191`), 結果短, 易於逐條修.
- 劣勢:
  - 無對比度數值規則: 全文提到 contrast 只有 `command.md:140` (互動狀態要比靜止狀態更明顯), 沒有 WCAG 4.5:1 之類的門檻; 也沒有任何執行期量測步驟.
  - 框架偏向: `command.md:30` `focus-visible:ring-*`, `command.md:42-43` `htmlFor`, `spellCheck={false}`, `command.md:73-74` `truncate`, `line-clamp-*`, `min-w-0`, `command.md:82` `priority`, `command.md:99` `useState` / nuqs, `command.md:131-135` Hydration Safety 全節針對 SSR. Vue / Svelte / 純 HTML 專案需自行轉譯.
  - 規則張力: `command.md:39` 要求 inputs 有 `autocomplete`, `command.md:48` 又要求非認證欄位 `autocomplete="off"`; `command.md:145` 要求 Title Case, 與非英文介面無關. 模型如何取捨屬 `推論`.
- 結論: 涵蓋面廣且具體, 適合 code review 階段; a11y 合規認證仍需 axe / Lighthouse 等實測.

### 跨 Agent 可攜性

- 證據: repo `README.md:3-5` 宣告遵循 agentskills.io 格式並上架 skills.sh; `AGENTS.md:92-104` 安裝方式為 `npx skills add vercel-labs/agent-skills --skill <name>`, Claude Code 手動 `cp -r skills/<name> ~/.claude/skills/`, claude.ai 上傳 zip 或貼上 `SKILL.md`; `skills.sh.json` 將它列在 `Design` 群組. 上游 `install.sh:29-83` 會把 `command.md` 裝到 Amp, Claude Code (`~/.claude/commands/`), Cursor, OpenCode, Windsurf (append 到 `global_rules.md`), Antigravity.
- 優勢: 同一份規則有 skill 與 slash command 兩種形式, 主流 coding agent 都有路徑; zip 與資料夾內容一致 (已 `diff` 驗證).
- 劣勢: `command.md:8` 使用 `$ARGUMENTS` 佔位符, 這是 slash command 語法; skill 形式需靠 `SKILL.md:32` 指示模型自行代換. 上游 `install.sh` 以 `curl ... main/command.md` 安裝, 拿到的是未釘版的 `main` (`install.sh:30,38,46,54`).
- 結論: 可攜性高; 優先用 `npx skills` 裝 skill 版, 以取得釘版內容.

### 維護活躍度

- 證據: skill 歷史 (`gh api repos/vercel-labs/agent-skills/commits?path=skills/web-design-guidelines`): `9e2435e` 2026-01-14 由 `vercel-design-guidelines` 改名, `42ecbef` 2026-01-15 加版本, `ba46938` 2026-01-16 修 frontmatter, `f766b64` 2026-10-29 改內含釘版, `53b6e3d` 2026-10-29 補更新步驟文件. 上游規則最近 commit: 2026-08-18 加 a11y media 規則與修正 curly quote, 2026-10-05 加 desktop overscroll 規則; 上游 `main` HEAD 為 `434b7f9` (2026-11-21 via GitHub API), 與釘住的 SHA 相同.
- 優勢: 母 repo 34,621 stars, 2026-11-20 仍有 push, 由 Vercel 官方維護; 規則仍在演進且目前釘版即最新.
- 劣勢: 釘版後更新需人工跑 `SKILL.md:36-41` 四步驟 (查 SHA, 下載, 改 SHA-256, 重建 zip), 上游新規則不會自動進來; `vercel-labs/agent-skills` 無 LICENSE 檔, 依 GitHub 預設屬保留所有權利 (`推論`), 只有規則上游是 MIT.
- 結論: 目前維護良好; 重新安裝時比對 `SKILL.md:27` 的 SHA 與上游 `main` 即可判斷是否落後.

## 細節 (Details)

- 安裝:

  ```bash
  npx skills add vercel-labs/agent-skills --skill web-design-guidelines
  # 或手動
  cp -r skills/web-design-guidelines ~/.claude/skills/
  ```

- 呼叫: 說 `review my UI src/components/**/*.tsx`, 或在支援 slash command 的工具以 `/web-design-guidelines <file-or-pattern>` 指名.
- 規模: 2 檔, 共 242 行; `command.md` 8,055 bytes (zip 內記錄長度).
- 本機現況 (2026-11-21): `~/.claude`, `~/.agents`, `~/.codex`, `~/.gemini` 與 `~/projects` 皆查無 `web-design*` 檔案, 尚未安裝.
- 取證方式: `git clone --depth 1` 至 scratchpad 空目錄 (HEAD `5a37a63`, 2026-11-20), 只讀不執行.

## 相關參考 (References)

- [impeccable](ai-unslop/impeccable.md) — 前端設計詞彙 + 反 AI slop 偵測, 含 `audit` 指令; 偏設計方向與美感
- [unslop](ai-unslop/unslop.md) — 把 vibecoded UI 改成有個性的設計
- [taste-skill](ai-unslop/taste-skill.md) — landing page 規則 + headless Chrome 實際檢查的 finish gate (有執行期驗證, 本 skill 沒有)
- 規則上游: <https://github.com/vercel-labs/web-interface-guidelines>
- skill 目錄: <https://skills.sh/vercel-labs/agent-skills>
- 改為釘版的 commit: <https://github.com/vercel-labs/agent-skills/commit/f766b64>
