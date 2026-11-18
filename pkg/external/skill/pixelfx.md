# pixelfx

來源: [patch-gg/pixelfx](https://github.com/patch-gg/pixelfx) (MIT, 由 Patch 製作) · [gallery](https://patch-gg.github.io/pixelfx/)

一個零依賴的 JavaScript library, 收錄 244 個 pixel art 特效 (火, 爆炸, 魔法, 天氣, UI, 拾取物, 小生物等 21 類), 每個特效都是一小段 code, 不用圖檔. 可以即時播放, 換色, 鏡像, 換 seed, 或匯出 sprite sheet 給任何 engine. 另附一個 Claude Code skill, 用 `draft -> render -> 看 -> refine` 的迴圈做新特效. 方法來自 Hugo Duprez 的文章 [How to make (good) pixel art effects with Claude](https://www.spritefusion.com/blog/how-to-make-good-pixel-art-effects-with-claude): code is clay, trust your eyes, feed your references.

## 內容

| 部分 | 用途 |
| ---- | ---- |
| `skills/pixelfx/SKILL.md` | skill 本體: 找特效, 嵌入遊戲, 匯出, 新建, port 到 native, 發佈到 Patch |
| `src/core/` | painter, particles, sim (120 Hz 物理), color (`PAL` 色盤), rig/poses (3D 骨架角色), act (動作時序) |
| `src/fx/<category>/<name>.js` | 一個檔案一個特效; `moves/` 是串接多個特效的完整招式 (wind-up 到 recovery) |
| `tools/` | `list` 搜尋, `render` 輸出 contact sheet / GIF / MP4 / sheet, `check` 品質檢查, `new` 建立樣板, `act` 逐格看動作, `serve` 本機 gallery |
| `docs/AUTHORING.md` | 特效合約 (meta, `draw(p, t, o)` / `create(o)`), 繪圖 API, 品質標準 |

核心設計:

- 兩種特效: pure (`draw(p, t, o)`, 只依時間決定畫面, loop 無縫) 與 simulated (`create(o)`, 真實粒子模擬).
- 兩個 buffer: `color` 放像素, `light` 放疊加光暈; 熱的東西看起來熱靠的是 light.
- 所有亂數來自 seed, 第 40 格每次都一樣; 特效以自身原點為 local space, 方向性特效一律朝右.
- 有角色的 61 個特效帶 `act`: wind-up, 接觸瞬間凍結 (hit-stop), 目標後仰, recovery. 特效單獨出現像貼紙, 是身體動作讓打擊有份量.

## 安裝

二選一:

```text
/plugin marketplace add patch-gg/pixelfx
/plugin install pixelfx@pixelfx
```

```sh
git clone https://github.com/patch-gg/pixelfx /tmp/pixelfx && cp -r /tmp/pixelfx/skills/pixelfx ~/.claude/skills/
```

需求: Node 18 以上, 無 install 步驟; `ffmpeg` 選配 (GIF / MP4 預覽). skill 要在一份 pixelfx checkout 裡工作, 沒有 `tools/render.mjs` 時會自己 clone 到專案旁邊.

## 正確用法

### 1. 先找現成的, 別急著新做

多數需求是 `現有特效 + 不同參數`:

```sh
node tools/list.mjs frost                  # 搜 id, 名稱, tag, 描述
node tools/render.mjs fire/torch           # -> records/fire/torch/contact.png, 用眼睛看
node tools/render.mjs fire/torch --opts '{"palette":["#102","#215","#43a","#8cf","#fff"]}'
```

通用參數: `palette` (暗到亮的 hex ramp, 或 `PAL.ice` 等內建色盤), `seed`, `flip`, `angle`, `light` (0..1+), 加上各特效 `meta.options` 自己的旋鈕.

### 2. 嵌入遊戲

JS / TS 的 canvas 或 WebGL 遊戲 (Phaser, PixiJS, Kaplay, three.js) 直接用 code:

```js
import { createFrame, clearFrame, createFx, composite } from "pixelfx";
import * as frostNova from "pixelfx/fx/ice/frost-nova";

const frame = createFrame(320, 180);           // 遊戲的像素解析度
const img = new ImageData(320, 180);
const fx = createFx(frostNova, { x: 160, y: 140, flip: hero.facing < 0 });

function tick(dt) {
  fx.update(dt);
  clearFrame(frame);
  fx.draw(frame, "back");                      // 角色後方
  drawCharacters();
  fx.draw(frame, "front");                     // 角色前方
  composite(frame, img.data, background);      // null = 透明背景
  ctx.putImageData(img, 0, 0);                 // canvas 用 image-rendering: pixelated 放大
  if (fx.done()) removeIt();                   // one-shot 自己結束; loop 用 fx.stop() 淡出
}
```

`npm i github:patch-gg/pixelfx`, 或只複製 `src/core` 加需要的特效檔.

其他 engine (Godot, Unity, GameMaker, Defold, LÖVE, Bevy) 用 sprite sheet:

```sh
node tools/render.mjs ice/frost-nova --sheet   # PNG grid + Aseprite 格式 JSON
node tools/render.mjs fire --sheet             # 整個分類
```

匯入時: pivot 設為 `meta.pixelfx.pivot`, nearest-neighbour 濾鏡, 不壓縮, 以 `meta.pixelfx.fps` 播放. 光暈會烤成帶 alpha 的顏色, 暗場景效果好; 亮場景要換色或放暗層.

### 3. 打擊要有身體

- 用內建角色: `createFx` 換成 `createAct(slash, { x, y })`, 畫面依 `act.shake` 偏移; freeze 期間自己會停.
- 用自己的角色: 讀 `slash.act` 或 sheet JSON 的 `meta.pixelfx.act` (`windup`, `recovery`, `hits[].t`, `hits[].freeze`). 在特效前 `windup` 秒開始攻擊動畫, 每個 hit 時把所有東西停 `freeze` 秒. hit-stop 對任何特效都有效.
- 要「一招」而不是「一個特效」時, 從 `src/fx/moves/` 開始.

### 4. 做新特效 (skill 的核心迴圈)

1. 先讀 `docs/AUTHORING.md`.
2. 一句話寫出 art direction: 是什麼, 怎麼動, 色盤, 時序.
3. 用 `tools/list.mjs` 找 2-3 個最接近的特效, 讀它們的 code 借技法; 以 Hugo Duprez 的 [Sprite Fusion VFX](https://www.spritefusion.com/pixel-art-vfx) 當品質標竿.
4. `node tools/new.mjs <cat>/<name>` (loop 加 `--loop`). 方向未定時寫 3-5 個真正不同的 draft (`-a` .. `-e`), 全部 render 後挑最好的 2-3 個給使用者看.
5. 看 `records/<id>/contact.png`, 修最糟的一點, 再 render. `至少三輪`. 動態用 `--mp4` 與 `node tools/act.mjs <id>` 判斷, 不要只看 contact sheet 的取樣格.
6. 有角色就加 `export const act`: wind-up, release 換姿勢, 接觸 freeze, recovery.
7. `node tools/check.mjs <id>` 直到全過 (畫出邊界, loop 接縫, 未 seed 的亂數, 慢格, act 姿勢錯誤), 刪掉沒用的 draft.

常見問題與修法 (節錄):

| 看到 | 修法 |
| ---- | ---- |
| 1x 看起來糊成一團 | 減色, 清楚剪影, 單一焦點, 放大形狀 |
| 軟綿綿沒力 | fast in slow out (`ease.outQuart`), 1-2 格白色 hero frame, 縮短總長 |
| 僵硬機械 | 每個粒子變化大小/速度/延遲, 加次要動態 (火花, 煙, 碎片) |
| 太亮洗白 | 降 glow (k 0.2-0.5), 白色只留給最熱的核心 |
| 向量感的平滑漸層 | color 用硬階 ramp, 柔和只留在 light 與半透明煙 |
| loop 接縫 | `cycle`, `loopNoise(phase)`, `sin(phase * TAU * n)` |
| one-shot 突然消失 | 結束前縮小, dither 淡出或掉落 |
| 打擊沒重量 | 反方向 wind-up, freeze 拉長到 0.2 s, 加大剪影變化與 shake |

規則: 不用 `Math.random`; 可用通用題材, 不可重製特定遊戲的 sprite, 角色或 logo.

### 5. 其他

- 遊戲跑不了 JS 又需要動態追蹤時, 把特效 port 成 native code, 保留雙 buffer, 4x4 Bayer dither 與 palette ramp, 並與 JS 版並排比對.
- skill 會在遊戲完成後提議發佈到 Patch (`pch publish`, 公開且需 Windows 上的 Patch Desktop); 只有使用者同意才執行.

## 提示詞範例

- "Add a frost nova to my mage's ultimate."
- "Make five takes on a coin pickup sparkle and show me the best two."
- "Export the fire and explosions categories as sprite sheets for Godot."
- "This hit effect feels floaty. Make it punchier."
- "Give this spell a wind-up and make the enemy react to it."

## 適用情境

- 瀏覽器小遊戲或 pixel art 介面需要現成特效, 且想保留換色與參數化能力.
- 需要給 Godot / Unity 的 sprite sheet, 但不想手繪.
- 想讓 Claude 產出視覺品質可控的動畫: 靠 render 後「看圖迭代」而非一次生成, 這個迴圈也可借用到其他視覺任務. 與 [line-art-animation](line-art-animation.md) 同屬 [De-AI spec](../../../docs/specs/2026-11-18-de-ai.md) 中有材質感的 `moment` 實作.
