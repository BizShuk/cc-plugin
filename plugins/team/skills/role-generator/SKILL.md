---
name: role-generator
description: >
    Use when the user wants to generate, write, or optimize a system prompt or role definition for an AI agent. Triggers on: "role generator", "role prompt", "system prompt design", "角色提示", "角色設定", "產生角色", "系統提示設計".
version: "1.0.0"
allowed-tools: []
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 角色提示生成技能 (Role Prompt Generator Skill)

此技能協助使用者依據五大核心原則設計高品質的 `系統提示 (System Prompt)`，並可融合大型科技公司的文化特色（如 Amazon、Meta、TikTok、Google）以提升模型的人設定錨效果。

## 核心設計結構 (Core Prompt Structure)

每個產出的角色 `系統提示 (System Prompt)` 必須包含以下五個核心區塊：

1. `身分與專長 (Identity & Expertise)`：定義角色是誰、資歷以及擅長領域。
2. `職責範圍 (Scope of Responsibility)`：明確劃分可以做與不可以做的事（邊界限制）。
3. `核心技能與思考方式 (Skills & Methodology)`：該角色擁有的硬技能與軟技能，以及面對問題時的邏輯（如 Working Backwards）。
4. `輸出格式 (Output Format)`：定義該角色的產出結構（如 JSON、Markdown、程式碼註解等）。
5. `限制與護欄 (Constraints & Guardrails)`：定義錯誤處理方式或不確定時應採取的措施。

---

## 輸出格式範本 (Output Format Template)

請依以下結構輸出角色的 `系統提示 (System Prompt)`：

```markdown
# [角色英文名稱] (中文角色名稱)

[一段描述身分與大廠文化特點的引導語，例如：You are a Senior Software Engineer at a large tech company...]

## 職責範圍 (Scope)
- [可以做的事項 1]
- [可以做的事項 2]
- [明確指出不可以做的事項]

## 具備技能 (Skills)
- 核心技術技能 (Core Technical Skills)：[列出該角色專屬的工具或知識庫]
- 跨職能技能 (Cross-Functional Skills)：[列出溝通、協作與決策的軟實力]

## 思考方式 (How you think)
- [思考路徑與決策邏輯 1]
- [思考路徑與決策邏輯 2]

## 輸出格式 (Output format)
- [定義輸出的首要元素與後續結構]

## 品質標準與護欄 (Quality bar & Guardrails)
- [定義判定合格的標準]
- [遇到不確定性時的暫停與發問機制]
```
