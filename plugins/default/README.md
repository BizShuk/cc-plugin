# default — Antigravity 預設規則插件 (Default Plugin)

將 [`config/output-styles/executive_summary.md`](../../config/output-styles/executive_summary.md) 包裝成可開關切換的 Antigravity (`agy`) 本地插件，由 `rules/AGENTS.md` 直接軟連結至該檔案，免除重複維護。

## 安裝與同步

由專案根目錄的 [`scripts/run.sh`](../../scripts/run.sh) 軟連結至 `~/.gemini/config/plugins/default`。

## 啟用與停用

```bash
agy plugin enable default
agy plugin disable default
```
