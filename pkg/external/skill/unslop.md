# unslop

來源: [heydamir/unslop](https://github.com/heydamir/unslop) (CC BY 4.0, 作者 Damir Sokolovsky, 引用或改編需標註作者並連回原 repo) · [Unslopping series](https://x.com/heydamir/status/2115063463063998736)

一個 Claude Code / Codex / Google Antigravity 通用的 skill, 把 `vibecoded UI` (AI 隨手生成, 一看就是 AI 做的介面) 改成有自己個性的設計. 它先找出 slop (紫藍漸層, glow, emoji label, 喊口號的文案, 不一致的間距), 再從產品本身推導 direction, 像設計師做 redline review 一樣逐層修正. 適用 Next.js, React, Vue, Svelte, Tailwind 或純 HTML / CSS, 涵蓋 dashboard, app 與 landing page; 只有網址或截圖時也能做.

設計概念 (direction 五欄, 9 個 pass, slop 特徵清單, template check) 已整理在 [De-AI spec](../../../docs/specs/2026-11-18-de-ai.md), 本文只講怎麼裝, 怎麼用, 以及 skill 的運作邊界.

## 內容

| 部分 | 用途 |
| ---- | ---- |
| `skills/unslop/SKILL.md` | skill 本體: 流程, direction, 9 個 pass 的規則, story video, 禁止事項與 checklist |
| `references/signs.md` | audit 用的 slop 特徵清單 |
| `references/tokens.md` | 每個 pass 的確切數值: spacing scale, neutral ramp (graphite / warm / ink), 字體家族, radius, easing 與 duration |
| `references/visuals.md` | 哪些視覺風格有效, 如何挑主題, image prompt 配方與擺放方式 |
| `references/refs.md` + `refs/*.webp` | 完成範例 (dashboard / hero after, 雕版風, 3D 場景), 只借品質不借題材 |
| `references/effects.md` | code 做的視覺處理, 例如用字元畫出圖像 |
| `references/rtl.md` | 阿拉伯文, 波斯文, 烏爾都文, 希伯來文的右到左規則, 與 pass 衝突時以它為準 |
| `scripts/snapshot.mjs` | 把線上網站存成可編輯的 HTML 加截圖 |
| `scripts/story.mjs` | 把 story file 錄成 unslop 過程影片 |
| `examples/` | vibecoded 練習畫面 (`dashboard.html`, `hero.html`, `store.html`) |

## 安裝

Claude Code 以 plugin 安裝 (全域):

```sh
claude plugin install unslop --marketplace heydamir/unslop
# 舊版沒有 --marketplace 時:
claude plugin marketplace add heydamir/unslop
claude plugin install unslop@heydamir
```

- 只裝在單一專案: 在專案根目錄執行並加 `--scope local`; 團隊共用則加 `--scope project` 並 commit `.claude/settings.json`.
- 裝完開新 session 或執行 `/reload-plugins`; 更新用 `claude plugin update unslop@heydamir`.
- 若 `~/.claude/skills/` 或 `.claude/skills/` 曾手動複製過 unslop, 要刪掉, 否則 `/unslop` 會跑到舊副本.

其他 agent:

```sh
# Codex
$skill-installer install https://github.com/heydamir/unslop/tree/main/skills/unslop
# 或手動: cp -r unslop/skills/unslop ~/.agents/skills/
# Google Antigravity
cp -r unslop/skills/unslop ~/.gemini/config/skills/
```

選配需求: 抓線上網站或錄 story video 需要 Playwright + Chromium (`npm install -D playwright && npx playwright install chromium`, 在執行 skill 的專案內安裝); story video 另需 `ffmpeg` 在 PATH 上.

## 正確用法

### 1. 叫用

Claude Code 用 `/unslop` (名稱衝突時用 `/unslop:unslop`), Codex 用 `$unslop`, 或直接用自然語言描述, agent 會自己挑到這個 skill.

### 2. 依輸入來源準備原稿

| 手上有 | 做法 |
| ------ | ---- |
| 專案原始碼 | 找出渲染該畫面的檔案, 能跑就跑起來截圖, 以使用者看到的畫面判斷 |
| 只有網址 | `node <skill>/scripts/snapshot.mjs <url> unslop/original.html` (可加 `--dark` / `--light`), 同時產出 `original.live.png`; 會先捲動整頁以抓到 lazy-load 內容, 但 script 畫的動態背景不會保留 |
| 只有截圖 | 先手刻成一個靜態 HTML 當原稿 |

後兩種的產出是 `一份重新設計過的 code`, 交給擁有該網站的人實作, 不會直接改到線上網站.

### 3. 一次執行的流程

1. Audit: 對照 `signs.md` 列出每個 slop, 一行一個並寫出位置.
2. 判斷程度: `Heavy slop` 跑完整流程; `Light slop` 走 polish 模式, 只修 audit 找到的, 沿用既有品牌; `No slop` 直說並最多建議 2-3 個附確切數值的調整.
3. 寫 5 行 direction (Product, Character, World, Moments, Signature) 後 `停下來等使用者確認`. 只有使用者事先說 "just do it" 才跳過.
4. 依序跑 9 個 pass: Spacing -> Color -> Icons -> Copy -> Type -> Layout -> Details -> Visuals -> Motion.
5. 跑 checklist, 再 render 比對 audit 清單, 最後回報 direction, 每個 pass 改了什麼, 以及寫過的 image prompt.

效率約束: 一次執行最多 render 3 次; 要生成圖時, direction 一確認就送出, 生成期間繼續做前面的 pass.

### 4. Hero 視覺與付費生成

- Dashboard 與 app 幾乎不需要生成圖, moment 來自產品本身 (細心繪製的 chart, 用 code 做的實體物件如付款卡).
- Landing hero 首選一張來自產品世界的寬幅 scene, 文字放在刻意留白處, product preview 從底部升起.
- 會使用已連線的 image service (fal, Replicate, OpenAI images, Higgsfield); 沒有就提議接成 MCP server. `每次付費生成前都會秀出 prompt 並詢問`. 沒有工具時, 用 code 做得到的先做, 其餘寫好 prompt 並留下正確尺寸的空位.

### 5. Story video (選用)

一支方形, 約 10 秒的影片, 畫面一步步自我 unslop, slop 上有紅色 redline. 只有使用者要分享時才做 (會多花幾分鐘):

1. 建 `unslop-story/story.html`: 原稿的獨立副本, 原始樣式是 step 0.
2. 3-5 個 step 以 `<meta name="unslop-steps" content="Spacing, Color and icons, Type and copy, Details">` 宣告; 每個 step 寫成疊加的 CSS (`.s1 .card { padding: 24px }`), 只改值不改結構, 讓瀏覽器能 tween.
3. 無法 tween 的處理: 漸層放獨立圖層淡出; 換字用 `<span class="swap" data-at="3">` (全片最多 8 個); 消失用 `data-out`, 出現用 `data-in`; pill 圓角寫成高度一半而非 `9999px`.
4. Redline 標記: `data-gap` (每個 stack 4-6 個), `data-radius` (2-4 個), `data-slop` (選用, 最多 5 個).
5. 鏡頭: `<meta name="unslop-frame" content="x y width">` 裁到變化最大的區塊.
6. 輸出: `node <skill>/scripts/story.mjs unslop-story/story.html unslop-story.mp4` (`--wide` 為 16:9, `--play` 先在瀏覽器循環預覽), 交付前用 ffmpeg 抽幾格檢查.

story file 只是展示用; 正式成果仍以專案原本的 stack 交付, 不含 `.s*` class 與 story 屬性. 依授權, 每支影片帶 `unslop by @heydamir` 字樣.

### 6. 不會做的事

- 不改產品名, 不改功能, 內容與 section 順序; 不發明功能, 宣稱或資料 (無法查證的宣稱會詢問, 不實就移除).
- 不以新裝飾取代舊裝飾 (glass, grain, noise, neon, bento grid 都只是換一種潮流).
- 不因為「想不一樣」就改動已經好的畫面; 也不讓每個產品都落到同一種外觀.

## 提示詞範例

- `/unslop this dashboard`
- `/unslop the hero of https://example.com`
- `/unslop my landing page and make the story video so I can post it`
- "這個 settings 頁看起來很 AI, 只修 audit 找到的問題, 保留現有品牌色."
- "Audit this screen and propose a direction, but don't edit yet."

## 適用情境

- 用 AI 快速做出的 dashboard / landing page 需要在上線前收斂成有設計感的版本, 且不想動到功能與資訊架構.
- 只拿到客戶的網址或截圖, 要產出一份可交付的 redesign code 與前後對照.
- 想把「先 audit, 再定方向, 停下確認, 再分層修」的審查節奏套用到其他設計任務. 與 [pixelfx](pixelfx.md), [line-art-animation](line-art-animation.md) 搭配時, unslop 負責讓 chrome 安靜, 後兩者負責提供有材質感的 `moment`.
