# Big Tech Product Team — Role-Based System Prompts (with Skills)

_Modeled on Meta / Google / Amazon / TikTok-style cross-functional product teams._

`How to use:` Copy a block into the `system prompt` field of each agent. Each prompt is self-contained. To make an agent reply in a specific language, append: `Always respond in Traditional Chinese, keeping technical terms in English.`

`Team pattern:` Product Manager + Engineering Manager = `orchestrators (協調者)`; the specialists are `workers`. Designer -> Eng -> QA -> SRE roughly forms a `pipeline (流水線)`.

`Skill legend:` Each role lists `Core technical skills (硬技能)` and `Cross-functional skills (軟技能)`. At senior level, the soft skills usually matter more.

---

## 角色定位與檔案連結 (Role Definitions and Links)

| #   | 角色 (Role)                                        | 一句話定位 (One-sentence Description) | 檔案連結 (File Link)                                                                                                                |
| --- | -------------------------------------------------- | ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| 1   | 產品經理 (Product Manager / PM)                    | 定義「為什麼做、做什麼」              | [product_manager.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/product_manager.md)                                 |
| 2   | 工程主管 (Engineering Manager)                     | 拆解技術任務、做最終技術決策          | [engineering_manager.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/engineering_manager.md)                         |
| 3   | 後端工程師 (Backend Engineer)                      | 設計可擴展、容錯的伺服器邏輯          | [backend_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/backend_engineer.md)                               |
| 4   | 前端工程師 (Frontend Engineer)                     | 介面效能與無障礙 (a11y)               | [frontend_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/frontend_engineer.md)                             |
| 5   | 機器學習工程師 (Machine Learning Engineer / MLE)   | 把模型推到生產環境                    | [machine_learning_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/machine_learning_engineer.md)             |
| 6   | 資料科學家 (Data Scientist)                        | A/B 實驗與因果分析                    | [data_scientist.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/data_scientist.md)                                   |
| 7   | 資料工程師 (Data Engineer)                         | 可靠的資料管線 (pipeline)             | [data_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/data_engineer.md)                                     |
| 8   | 產品設計師 (Product / UX Designer)                 | 以使用者需求設計流程                  | [product_designer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/product_designer.md)                               |
| 9   | 網站可靠性工程師 (Site Reliability Engineer / SRE) | SLO、監控、故障處理                   | [site_reliability_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/site_reliability_engineer.md)             |
| 10  | 資安工程師 (Security Engineer)                     | 威脅建模與防禦                        | [security_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/security_engineer.md)                             |
| 11  | 推薦系統工程師 (Recommendation Systems Engineer)   | 排序與個人化 (TikTok/Meta 招牌)       | [recommendation_systems_engineer.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/recommendation_systems_engineer.md) |
| 12  | 成長分析師 (Growth / Product Analyst)              | 找出漏斗 (funnel) 缺口                | [growth_analyst.md](file:///Users/shuk/projects/cc-plugin/pkg/agent_team/roles/growth_analyst.md)                                   |

---

_Tip: For a multi-agent setup, keep shared rules (output language, company terminology, available tools) at the project level, and let each prompt carry only its role-specific content. Use the skill lists above as a quality rubric — if an agent's output ignores a listed skill (e.g. the Data Scientist omits confidence intervals), that is your signal to tighten its prompt._
