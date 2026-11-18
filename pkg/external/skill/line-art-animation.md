# line-art-animation

來源: [JDhruv14/line-art-animation](https://github.com/JDhruv14/line-art-animation) (MIT)

一個 Claude skill, 搭配一個零依賴的 TypeScript library, 用來產生 ink-on-paper (墨線手繪) 風格的動畫 SVG 元件. 人物會走路, 騎車, 對滑鼠互動做出反應; 場景有日夜循環. 所有圖形都是即時以 SVG 繪製, 任何尺寸都清晰.

## 內容

| 部分 | 用途 |
| ---- | ---- |
| `SKILL.md`, `references/` | 工作流程與風格規則, 分為 style, motion, scenes, verify, matching (對照參考影片), video 幾份指南 |
| `lib/ink.ts` | 繪圖 primitives: limb 是沿關節曲線畫的兩道開放筆觸, 這是手繪感的主要來源 |
| `lib/gait.ts` | walk/run solver: 用解算而非 keyframe 產生腿部動作. 著地腳以移動速度後滑, 擺動腳走 minimum-jerk 路徑, 髖部高度隨步伐起伏 |
| `lib/motion.ts` | spring, pendulum, wind |
| `lib/daynight.ts` | 依時段變化的色盤 |
| `examples/` | `walker` (約 150 行的單一人物) 與 `village` (完整日式街景, 約十幾種互動) |
| `scripts/` | 為影片做 deterministic 逐格擷取, 以及量測參考影片的工具 |

## Skill 工作方式

1. 把 `lib/` 複製進使用者的專案.
2. 先畫一個人物, 請使用者核可, 再建構整個場景.
3. 交付前逐格 (frame by frame) 檢查動態.
4. 選配: 對照參考影片重現步態, 或渲染影片 (例如 30 秒 1080p).

## 安裝與試跑

```sh
git clone https://github.com/JDhruv14/line-art-animation ~/.claude/skills/line-art-animation
```

試跑範例需要 Node 22.18 以上:

```sh
npm install
npm run dev        # http://localhost:5173
```

URL flags: `?start=0.9` 設定起始時段 (0 是午夜, 0.5 是正午), `?tau=0.5` 固定時段, `?day=30` 讓一整天只跑 30 秒, `?debug` 手動逐步推進時間.

## 適用情境

- 想要有手繪質感, 可互動的 hero 或插圖動畫, 不用影片檔.
- 想把參考影片轉成可循環播放的 line-art component.
- 可當作 [De-AI spec](../../../docs/specs/2026-11-18-de-ai.md) 中 `moment` 的一種實作: 一個來自產品世界, 有材質感的視覺元素.
