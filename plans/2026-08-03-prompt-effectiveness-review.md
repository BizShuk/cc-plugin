# Prompt Effectiveness Review — 全 repo prompt 審查

日期: 2026-08-03
範圍: 51 個 SKILL.md + 3 個 agent prompts + config/CLAUDE.global.md + output style + .claude/skills
方法: 4 個並行 read-only review agents, 逐檔閱讀並以 ls/find/grep 驗證引用路徑與事實
狀態: 唯讀審查, 尚未修改任何檔案

## 總量 (Token Cost Baseline)

- 全部 prompt 內容約 6,748 行; frontmatter 合計約 30 KB (~7-8K tokens 常駐每個 session)
- 常駐成本最高群: god 9 份 description (~850 tokens) 觸發詞過廣; ultra-explore 11 份
  `disable-model-invocation: true` 的 kb-* description 帶雙語 trigger 清單 (~1.5K tokens 完全無觸發作用)
- On-trigger 最重: `scrapling` (~6K tokens), `playwright-cli` (~3K), `topology-builder` (285 行)

## P0 — 錯誤事實與死引用 (會直接誤導模型)

1. `plugins/explore/skills/project-route/SKILL.md:36` — 分類只列 4 個, 實際 13 個;
   9 個分類對 routing 隱形; L120 又給 `collections` 規則自相矛盾; L64-71 宣稱的
   3 天 auto-rebuild 未生效 (index 已 12 天未更新)。`projects.json.categories` 同樣過期。
2. `plugins/god/skills/llm-mechanics/SKILL.md:29-37` — transformer 機制描述錯誤
   (否定詞「向量減法」、同音字靠「位置編碼 + POS tagging」), 教模型錯誤的自我模型。
   只有 L39-49 Thinking Mode 表有價值。
3. Dead references 全清單:
   - `universal-aggregate:47` → `references/ontology-template.md` 不在該 skill 目錄 (實在 `plugins/god/references/`)
   - `auto-evolving:151-159` → `engine/system_prompt.md` 等 3 檔不存在
   - `feature.md:153,155,156` → `@golang-refactor` agent 不存在 (死路)
   - `naming-convention:61-62` → `[[consistency]]`/`[[folder-structure]]` 無此 skill
   - `changelog:30` → `pip install -e .claude/skills/changelog/scripts` 路徑不存在
   - `topology-builder/scripts/workflow.js:15` → 寫死 `plugins/general/...` 舊路徑
   - `kb-verify:36,42` → 硬編 `~/projects/product/projects/` 絕對路徑
   - `skill-frontmatter:25` → 宣稱的 CI yaml.Unmarshal 驗證不存在; L26 引用兩個不存在 skill
   - `marketplace-config:253-257` → Related 3 個 skill 全部不存在
   - `plugins/tools/README.md:9-10` → `iphone-deploy`/`iphone-mac-permission` 不存在; 漏列外部 `bizshuk/autop`
   - `plugins/ultra-explore/README.md:122` → changelog 歸屬敘述過期
4. `sort-todo` 綁死別的專案: L93 引用 `src/todo/todoStore.ts`, L66-75 是某 VSCode
   extension 的 domain taxonomy — 在任何其他專案觸發都會套錯 section 名。
5. `review-coordinator.md` Part 5 範本用的 dimension 標籤 (`consistency` 等)
   對不上 Part 3 路由表的真實 skill 名。
6. `happy` skill: `$ARGUMENT[3]` 非法語法 + 非規範 frontmatter 欄位 + PostToolUse
   hook 在每次 Read 都 echo — 示範了壞範例且常駐 392 chars description。

## P1 — 大型重複 (單一事實多份 owner, 必然分岔)

1. `統一介面表三份 owner`: `CLAUDE.global.md` 原版 + `project-docs:39-52` 複本 +
   `docs-consolidation:112-120` 複本。修法: 只留 global 一份, 其餘改 pointer。
2. `topology-builder ⟷ kb-spec` 約 150 行逐字重複 (Entity/Edge format, verification
   shell loop 連 kb-verify:50-56 也是同一份)。修法: topology-builder 砍 ~130 行,
   格式規則一行指向 `kb-spec/references/file-formats.md`。
3. `kb-coordinator.md ⟷ ultra-explore/SKILL.md` 約 70 行逐字重複 (含相同假數字
   `214 captures`)。修法: coordinator 瘦身到 30 行內, 只留角色定義 + 失敗處理。
4. `changelog skill (891 LOC scripts)` 被 `kb_history.py` (376 行單檔) 完整取代,
   且缺 `disable-model-invocation: true` 會被意外自動觸發 + 安裝路徑是死的。
   修法: 整個目錄刪除。
5. `anti-sabotage.md` 是 SKILL.md:19-82 的逐字複本 — 刪。
6. `markitdown:21-29` 格式表逐字複製到 `content-summarizer:28-38` — markitdown
   擁有, summarizer 改連結。
7. `claude-plugin-metadata` 與 `.claude/skills/marketplace-config` 覆蓋同兩份
   manifest 且較淺; `marketplace-config:143-165` 又重複 `skill-frontmatter` — 收斂為單一 owner。

## P2 — 常駐 context 瘦身 (每 session 都省)

1. `config/CLAUDE.global.md` (96 行常駐): L62-83 內容歸屬表 + L85-96 檔案命名
   共 34 行只在寫文件時有用 → 移入 docs-consolidation skill, 留一行 pointer。
   L27-28 硬編 13 分類名 → 指向 projects.json。L13-14 CLAUDE.md 已由 harness
   自動注入, 只留 README.md 半句。L44 中英互相矛盾需重寫。
2. `Output style` 內部矛盾: 配額 15+8+2=25 ~ 20+12+5=37 對不上宣告的 30-40;
   `punchtuation` typo; no full-width punctuation 規則與檔案自身全形標點矛盾;
   L19 標題語言規則與 L11-13 英文section名矛盾; emoji 規則與 subagent harness 衝突。
3. `docs-consolidation` description 1627 chars (全 repo 最長) → 砍到 ~300。
4. 11 個 `disable-model-invocation: true` 的 kb-* skill 刪雙語 trigger 清單
   (~1.5K tokens/session, 零觸發損失); `kb-spec` description 改一行
   "Reference spec read by other kb-* skills; not directly invocable."
5. god 觸發詞過廣: `universal-generate` 的 "create/build/generate/draft" 與
   `universal-review` 的 "review/audit" 跟具體 skills 及內建 /review 搶觸發 →
   改為 scope-limiting 描述 (target-agnostic、無既有範例可循等)。
   `grand-unified-theory` 純索引零可執行步驟 → 降級為 plugins/god/README.md 章節。

## P3 — 結構重整 (progressive disclosure)

1. `scrapling` 23KB → ~5KB: 選項表移 `references/cli-options.md`, Python 區塊移
   `references/python-api.md`; 刪除損壞的 vendored `Scrapling/` clone (.git 損毀,
   git 看不見, 子目錄全空)。
2. `playwright-cli:30-181` ~150 行 CLI 目錄 → `references/commands.md`。
3. `apple-reminders` 是四兄弟中唯一無 references/ 者: L52-96 CLI dump 外移,
   補 Rules/Related, 對齊 apple-calendar 範本 (該檔為 exemplar)。
   注意: 不要抽共用 boilerplate — SKILL.md 無 include 機制, 對齊即可。
4. `system-laws:48-113` 26 行泛用建議移 references/, 留 13x10 法則表。
5. `auto-evolving` 的 Legacy Migration + 8 維評分表移 references/
   (且 L36 「模型評分不是驗證」與評分制自相矛盾, 需擇一)。
6. `markitdown:83-109` Apple Notes 教學 (27 行跨域污染) → apple-notes references。
7. `team-design` + `orchestration-config` 合併; `role-generator` 指向
   `plugins/team/roles/` 12 份現成 exemplar (目前 skill 看不見最好的 few-shot 素材);
   L17 大廠文化融合承諾未兌現 — 實作或刪除。
8. `topology-builder` 改 `disable-model-invocation: true` 並統一 storage root
   (目前 `~/projects/product/topologies/` 與 kb 的 `projects/` 兩座互不連通)。

## P4 — 慣例一致性

- `plans/` 檔名: `system-planner:83` 與 `business-planner:63` 違反
  `YYYY-MM-DD-<topic>.md`; business-planner L63/L69 雙路徑矛盾。
- `markdownlint` 範圍矛盾 (skill 限 experiment/, README 稱全插件通用);
  MD043/MD044 無 config 不可能 on; bold 違規: sort-todo 21 處、project-task 15 處、feature.md 6 處。
- `12條宇宙法則.md` (New Age blog 轉貼, 9.5KB) 無 skill 引用 → 刪。
- `model-evaluator` probes 是被記憶的知名 benchmark 題 (bat & ball, strawberry) — 量不到東西。
- frontmatter tier 缺漏: changelog/summarize-sh/sort-todo 缺 version 或 allowed-tools;
  team 三 skill `allowed-tools: []` 空陣列語意存疑。

## 品質標竿 (可作為改寫範本)

- `mermaid` — progressive disclosure 教科書 (選擇表 + 27 references)
- `apple-calendar` — 精準觸發 + 具體 Rules + references 分離
- `kb-ingest-history` — 每行有用, `${CLAUDE_PLUGIN_ROOT}` 正確用法
- `naming-convention`、`daily-summary`、`content-summarizer` — 具體步驟 + 反捏造設計
