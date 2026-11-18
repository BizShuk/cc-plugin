# claude-spend

來源: [writetoaniketparihar-collab/claude-spend](https://github.com/writetoaniketparihar-collab/claude-spend) (npm 宣告 MIT, ~440 stars, 最新 npm `1.0.6`, 最後更新 2026-03-01)

讀 `~/.claude/` 下的 Claude Code session 檔, 在 localhost 開一個 Web dashboard, 回答「我的 Token 都花在哪些 prompt 上」. 零設定, 只支援 Claude Code.

## 功能

- 依 conversation / day / model 拆解 Token 用量.
- 標出最貴的 prompt 與使用模式 (usage pattern).

## 使用

```bash
npx claude-spend                # 開 http://localhost:3456
claude-spend --port 8080        # 改 port
claude-spend --no-open          # 不自動開瀏覽器
```

## 適用情境

- 想看「單一 prompt / 對話」層級的成本排行, 找出吃 Token 的元兇; 這是 ccusage 表格較難直觀呈現的視角.

## 注意

- 本機實測 dashboard 在 Safari 打不開, 請用 Chrome 等瀏覽器.
- 專案已約 8 個月未更新, 新版 Claude Code 紀錄格式若變動可能失效 (推論).
- 全程本機運作, README 聲明不送出任何資料.
