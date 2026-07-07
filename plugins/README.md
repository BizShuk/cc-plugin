# 插件目錄說明 (Plugins Directory Overview)

本目錄收納本專案的所有 `插件 (Plugins)`。除了本地開發的插件外，也透過 `.claude-plugin/marketplace.json` 與 `.gitmodules` 整合了多個 `外部插件 (External Plugins)` 與 `Git 子模組 (Git Submodules)`。

## 外部插件與子模組對照表 (External Plugins & Submodules Mapping)

以下為本專案所引入的外部插件與對應的 `Git 子模組` 狀態：

| 插件名稱 (Plugin Name)          | 來源倉庫 (Source Repository)                                                                          | 本地路徑 (Local Path)                                                                            | 類型 (Type)                           | 說明 (Description)                                             |
| ------------------------------- | ----------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ | ------------------------------------- | -------------------------------------------------------------- |
| `awesome-claude-code-subagents` | [VoltAgent/awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents) | N/A                                                                                              | `Marketplace 外部插件`                | 收錄優秀的 Claude Code 子代理 (Subagents) 資源。               |
| `gosdk`                         | [bizshuk/gosdk](https://github.com/bizshuk/gosdk)                                                     | [plugins/gosdk](file:///Users/shuk/projects/cc-plugin/plugins/gosdk)                             | `Marketplace 外部插件` / `Git 子模組` | Go 語言開發工具包，提供共用函式庫與設定管理。                  |
| `inf`                           | [bizshuk/inf](https://github.com/bizshuk/inf)                                                         | N/A                                                                                              | `Marketplace 外部插件`                | LGTM 觀測後端與運維基礎設施。                                  |
| `superpowers`                   | [obra/superpowers](https://github.com/obra/superpowers)                                               | [plugins/superpower](file:///Users/shuk/projects/cc-plugin/plugins/superpower)                   | `Marketplace 外部插件` / `Git 子模組` | 進階 AI Agent 工具與框架。                                     |
| `understand-anything`           | [Egonex-AI/Understand-Anything](https://github.com/Egonex-AI/Understand-Anything)                     | [plugins/understand-anything](file:///Users/shuk/projects/cc-plugin/plugins/understand-anything) | `Marketplace 外部插件` / `Git 子模組` | AI 輔助程式碼庫理解與知識圖譜生成工具。                        |
| `last30days`                    | [mvanhorn/last30days-skill](https://github.com/mvanhorn/last30days-skill)                             | [plugins/last30days-skill](file:///Users/shuk/projects/cc-plugin/plugins/last30days-skill)       | `Marketplace 外部插件` / `Git 子模組` | 趨勢與研究插件，支援 Reddit、X (Twitter)、YouTube 等趨勢研究。 |
| `ui-ux-pro-max-skill`           | [nextlevelbuilder/ui-ux-pro-max-skill](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill)       | [plugins/ui-ux-pro-max-skill](file:///Users/shuk/projects/cc-plugin/plugins/ui-ux-pro-max-skill) | `Marketplace 外部插件` / `Git 子模組` | UI/UX 設計與評估技能包。                                       |

## 其他內嵌 Git 子模組 (Other Embedded Git Submodules)

除了頂層插件外，某些本地插件的特定技能或套件包內部也以 `Git 子模組` 形式依賴了外部開源專案：

- `plugins/explore/skills/scrapling/Scrapling`
    - 來源：[D4Vinci/Scrapling](https://github.com/D4Vinci/Scrapling.git)
    - 用途：具備防機器人偵測繞過與 JavaScript 渲染的高效爬蟲工具。
- `plugins/explore/skills/summarize.sh/summarize`
    - 來源：[steipete/summarize](https://github.com/steipete/summarize.git)
    - 用途：文字與內容摘要工具。
- `plugins/media/seedance-2.0`
    - 來源：[Emily2040/seedance-2.0](https://github.com/Emily2040/seedance-2.0)
    - 用途：影片生成與多媒體創作相關資源。
- `pkg/system-prompts/CL4R1T4S`
    - 來源：[elder-plinius/CL4R1T4S](https://github.com/elder-plinius/CL4R1T4S.git)
    - 用途：系統提示詞 (System Prompts) 資源包。
- `pkg/tools/career-ops`
    - 來源：[santifer/career-ops](https://github.com/santifer/career-ops)
    - 用途：開發輔助工具集。

## 初始化與更新說明 (Initialization & Update)

若要初始化與更新所有子模組，請於專案根目錄下執行：

```bash
git submodule update --init --recursive
```
