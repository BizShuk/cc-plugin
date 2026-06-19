# base plugin

`base` 是 **must-install 預設插件**。提供每次 Claude Code 工作階段都會執行的底層 hooks — 屬於「裝了就別關」這類的低噪音基礎設施。

目前內容：1 個 hook，於 agent loop 結束時發出終端機鈴聲。

## 為什麼要裝

長任務或分心時，agent 結束的通知很容易錯過。一聲 BEL 比盯著終端機輪詢更省注意力。

## Hooks

| Event          | 行為                  | 目的                                       |
| -------------- | --------------------- | ------------------------------------------ |
| `Stop`         | `printf '\a'` → stderr | 正常完成時通知                              |
| `StopFailure`  | `printf '\a'` → stderr | 失敗/中斷時通知（語意與 `Stop` 相同）        |

兩個事件共用同一支 `stop-bell.sh`，行為無差異。

## 檔案結構

```text
plugins/base/
├── .claude-plugin/
│   └── plugin.json      # manifest
├── hooks/
│   ├── hooks.json       # Stop / StopFailure 註冊
│   └── stop-bell.sh     # 終端機 bell 實作
└── README.md
```

## 安裝

`base` 已在 `.claude-plugin/marketplace.json` 註冊為第一個 plugin，視為預設：

```bash
claude plugin install base
```

或在 `~/.claude/settings.json` 強制啟用：

```json
{
    "enabledPlugins": {
        "base": true
    }
}
```

## 停用方式

**不推薦** — 失去 agent 結束通知。

如必要，編輯 `~/.claude/settings.json`：

```json
{
    "enabledPlugins": {
        "base": false
    }
}
```

或暫時關閉單一 hook：把 `plugins/base/hooks/hooks.json` 對應 event 改為空陣列 `[]`。

## 疑難排解

| 現象                | 原因                              | 解法                                     |
| ------------------- | --------------------------------- | ---------------------------------------- |
| 沒聽到鈴聲          | 終端機關閉了 bell (`visible-bell`) | 終端機偏好設定開啟 audible bell           |
| 終端機顯示亂碼      | 不支援 BEL                       | 確認 `$TERM` 有效；腳本會 fallback 到 `tput bel` |
| Hook 沒執行         | Plugin 未啟用                      | 檢查 `enabledPlugins` 設定               |

## 設計原則

- **零依賴**：僅使用 `bash` 內建 + 可選 `tput`
- **冪等**：每次觸發只輸出一個 BEL
- **靜默失敗**：若 `tput` 不可用，fallback 到 `printf '\a'`，不輸出錯誤
- **無 stdout 污染**：bell 寫到 stderr，避免干擾 hook 輸出 JSON 解析
