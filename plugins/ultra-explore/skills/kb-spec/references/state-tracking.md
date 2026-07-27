# kb-spec — 狀態追蹤格式 (State Tracking Formats)

`run-id` 格式：`<yyyy-mm-dd>-<skill縮寫>-<slug>`，如 `2026-07-04-repo-service-a`。

## manifest.json（開跑時寫一次，之後唯讀）

```json
{
  "run": "2026-07-04-repo-service-a",
  "skill": "kb-ingest-repo",
  "source": "github.com/org/service-a",
  "created": "2026-07-04",
  "total_items": 1240,
  "batch_size": 50,
  "batches": [
    { "id": 1, "items": ["cmd/main.go", "model/order.go"] }
  ]
}
```

## progress.json（每批更新一次）

```json
{
  "run": "2026-07-04-repo-service-a",
  "status": "in-progress",
  "batches_done": 3,
  "batches_total": 25,
  "items_done": 150,
  "items_skipped": 4,
  "outputs": ["_inbox/2026-07-04-order-lifecycle.md"],
  "last_batch_finished": "2026-07-04"
}
```

- `status` 取值：`in-progress` | `done` | `failed` | `paused`
- 跳過的 item 必須在 `log.md` 寫一行原因（如 `SKIP vendor/lib.go — 排除清單`）

## STATUS.md（儀表板，每批同步更新）

```markdown
# KB Pipeline Status

| Run                        | Skill          | Progress | Status      | Updated    |
| -------------------------- | -------------- | -------- | ----------- | ---------- |
| 2026-07-04-repo-service-a  | kb-ingest-repo | 3/25     | in-progress | 2026-07-04 |
```

## 分批建議 (Batch Sizing)

| 來源             | 每批大小       |
| ---------------- | -------------- |
| repo 原始碼      | 50 檔          |
| git log 歷史     | 300 commits    |
| 文件 (web/檔案)  | 20 份          |
| 對話紀錄         | 200 則訊息     |
| schema 物件      | 30 表/topic    |
| distill captures | 20 個 captures |
