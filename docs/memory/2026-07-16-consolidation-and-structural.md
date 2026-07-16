# 一致性清理與結構性整合完成紀錄 (Consolidation and Structural Completion Record)

Status: `completed`

Source plan: `2026-07-08-consolidation-and-structural-plan.md`（原始內容保留於 git history）

## 結論 (Conclusion)

CLI、plugin metadata、文件單一來源與部署可移植性已完成整合。`plans/` 不再保留已完成設計；產品候選 TUI 經審核後轉入 backlog，Agent 設計則以已實作規格收斂。

## 完成項目 (Completed Work)

### CLI 基礎設施

- [x] Viper 預設值集中於 `config/config.go`，移除空白 embedded JSON。
- [x] 升級 `gosdk`，由 `gosdk/log` 安裝 `slog` handler。
- [x] 終止錯誤、retention 警告與 distill 完成事件使用結構化 logging。
- [x] Topology 自 `model/` 移至 `pkg/topology/`。
- [x] 建立 `cc-plugin topology verify/unlinked/query/backlinks/index/rewrite`。
- [x] 建立真正的 Cobra root，使 `cc-plugin distill` 與 `cc-plugin topology` 合約成立。

### Plugin 生態

- [x] 新增 `scripts/pluginmeta`，掃描 skill、agent、README、frontmatter、manifest 與 marketplace。
- [x] `plugin.json` 的 `skills`／`agents` 保持空陣列，使用標準目錄 auto-discovery。
- [x] `anti-sabotage` 轉為標準 `SKILL.md`；修正 `session-retro`、`summarize-sh` kebab-case。
- [x] 修正缺少觸發詞或 folded description 的 metadata。
- [x] 八個本地 plugin 均有 README，marketplace 補登 `god`。

### 文件單一來源

- [x] config、logging、topology、skill metadata 與 feature agent 設計移入 `docs/specs/`。
- [x] Claude usage TUI 經掃描上限、隱私、token 定義與 E2E 測試審核後移入 `docs/backlog/`。
- [x] `README.md`、`CLAUDE.md`、`README.todo` 與實際目錄／CLI 同步。
- [x] 重複或失效 plans 移除，歷史決策由本文件與 git history 保存。

### 部署可移植性

- [x] `.gitmodules` 只保留兩個實際 gitlink。
- [x] PM2 設定改用 PATH 與 `__dirname`，cron job 設 `autorestart: false`。
- [x] `run.sh` 使用 `set -euo pipefail`、`REPO_ROOT`、備份守衛與可重入 helper。
- [x] VS Code launch 對齊 `topology verify`，settings 移除重複與失效項目。
- [x] plugin-scoped MCP/LSP、JSON、PM2 與 shell syntax 均完成驗證。

## 驗證結果 (Verification Result)

```text
go test ./... -count=1                 PASS
go vet ./...                           PASS
go build ./...                         PASS
cc-plugin topology verify              OK
go run ./scripts/pluginmeta            OK
markdownlint affected experiment files PASS (MD013 disabled per skill)
jq plugin/MCP/LSP/VS Code JSON          PASS
node ecosystem.config.js               PASS
bash -n run.sh                          PASS
run.sh temporary-HOME idempotency       PASS
git submodule status                    PASS
git diff --check                        PASS
```

## 延後項目 (Deferred)

- Claude Code 使用統計 TUI：等待 `docs/backlog/2026-07-16-claude-usage-tui.md` 的實作前 gate 完成。
