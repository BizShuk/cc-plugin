# Apple Reminders (Apple 提醒事項)

這是一個名為 `apple-reminders` (Apple 提醒事項) 的 `skill` (技能)，主要透過 `remindctl` (提醒事項命令列工具) 在 `macOS` 系統上管理 `Apple Reminders` (Apple 提醒事項)。

其主要功能包括：

1. `View Reminders` (檢視提醒事項)：可以查看今天、明天、本週、過期或指定日期的提醒。
2. `Manage Lists` (管理列表)：列出所有列表、查看特定列表，以及建立或刪除列表。
3. `Create Reminders` (建立提醒事項)：支援設定標題、指定列表以及設定截止日期 (例如 `today`、`tomorrow` 或特定時間格式)。
4. `Complete/Delete` (完成或刪除)：透過 `ID` (識別碼) 來標記完成或刪除提醒事項。
5. `Output Formats` (輸出格式)：支援 `JSON` (JSON 格式)、`Plain` (純文字) 或 `Quiet` (僅顯示數量) 等輸出方式。

此 `skill` (技能) 適合用於需要將待辦事項同步到 `iPhone` 或 `iPad` 的場景，但不建議用於排定系統內部的 `alert` (警示) 或管理 `Calendar` (行事曆) 活動。

---

Source: <https://github.com/openclaw/openclaw/blob/main/skills/apple-reminders/SKILL.md>
