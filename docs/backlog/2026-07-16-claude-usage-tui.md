# Claude Code 使用統計 TUI (Claude Code Usage TUI)

Status: `deferred`

## 結論 (Conclusion)

保留產品價值，但本次一致性清理不進入實作。原計畫已確立 `bubbletea + bubbles + lipgloss` 與 master-detail 互動，但掃描邊界、token 定義、隱私輸出與端對端驗證尚未達到可落地標準。

## 已接受方向 (Accepted Direction)

- 串流讀取 `~/.claude/projects/*.jsonl` 與 `~/.claude/history.jsonl`，不保留完整 record。
- 資料掃描層與 TUI 狀態／render 分離；先交付 `--dump` 再接 UI。
- Session、Prompt、Token、Project 四種檢視採 master-detail 導覽。
- 新程式應置於 `cmd/tui/` 子套件，由 `cmd/root.go` 註冊，遵循目前 `export`、`memory`、`topology` 結構。

## 進入實作前的必要決策 (Required Gates)

- [ ] 定義全域掃描上限：總檔案數、總 bytes、單檔 bytes、逾時與取消行為；`MaxFileBytes=0` 不得代表無限制。
- [ ] 修正超長 JSONL 行策略：`bufio.Scanner` 遇到超過 token 上限會停止該檔，不能宣稱只 skip 單筆；需改用可恢復 reader 或明確 fail-file。
- [ ] 定義 `TotalTokens` 是否包含 cache creation/read，避免不同畫面重複加總或與供應商帳單口徑不一致。
- [ ] `--dump` 預設不得輸出 prompt 原文、完整本機路徑或 pasted contents；建立明確的 privacy contract。
- [ ] 補資料層 benchmark 與記憶體上限，使用至少一個大型 synthetic JSONL fixture。
- [ ] 除 parser 測試外，加入 loading/error/resize/navigation 的 `Update` 狀態測試與非互動 smoke test。
- [ ] 明訂損壞檔、部分掃描與 warning 時的 exit code，避免把不完整統計當成功結果。

## 驗收草案 (Draft Acceptance)

```text
go test ./cmd/tui -count=1
go test ./... -count=1
go vet ./...
go run . tui --dump --projects-dir <fixture> --history-path <fixture>
```

完成上述 gate 後，再由 backlog 晉升為新的 `plans/YYYY-MM-DD-claude-usage-tui.md`。
