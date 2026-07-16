# Feature Agent 與專門 Agent 邊界規格 (Feature and Specialized Agent Boundaries)

Status: `completed`

## 結論 (Conclusion)

`feature` 是新功能的端到端 composition agent，不取代 refactor、dead-code、performance、debugging 或 security 專門流程。共同工程規則由 skills 提供，agent 只保留角色、階段、委派邊界與輸出契約。

## 角色邊界 (Role Boundaries)

| 請求 | 路由 |
| :--- | :--- |
| 新功能規劃與實作 | `plugins/general/agents/feature.md` |
| Go 重構／死碼／效能審查 | 對應的 Go 專門 agent/skill |
| 一般品質與文件審查 | `plugins/review/agents/review-coordinator.md` |
| 邏輯錯誤與安全漏洞 | debugging、code-review 或 security-review 專門流程 |

`feature` 僅處理單一 repository；跨 repo 功能以明確介面契約交接，分別在各 repo 啟動獨立工作。

## Feature Agent 契約 (Feature Agent Contract)

```text
Understand → Clarify → Plan → Implement → Verify
```

- `Understand`：載入專案文件與適用 skills，辨識受影響套件與外部依賴。
- `Clarify`：一次整理 security、testability、observability、performance、maintainability 與 compatibility 決策。
- `Plan`：列出檔案、介面、遷移、測試與回滾方式。
- `Implement`：依專案分層落地，測試與實作同行。
- `Verify`：執行完整測試與靜態分析，回報已滿足及延後的 NFR。

## 委派與輸出 (Delegation and Output)

- Agent 依賴 skill 注入 Go naming、MVC 與 code-quality 規則，不把規則複製進每個任務。
- 遇到 refactor、bug、dead code 或 performance review，明確拒絕錯誤路由並指向專門流程。
- 完成摘要必含：變更檔案、測試結果、NFR 狀態、延後事項與跨 repo contract。
- `plugins/general/.claude-plugin/plugin.json` 透過 `agents/` 目錄自動探索，不維護重複 agents 清單。

## 實作證據 (Implementation Evidence)

- [x] `plugins/general/agents/feature.md` 已包含五階段流程、NFR 排序與拒絕模板。
- [x] `plugins/review/agents/review-coordinator.md` 已承擔唯讀跨維度審查。
- [x] general/review plugin README 已說明對應 agent。
- [x] plugin metadata 檢查會驗證 agent 檔案與 README，但保持 manifest auto-discovery。
