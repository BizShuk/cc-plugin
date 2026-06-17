---
name: character-setting
description: >
    Use when the user wants to define, design, or maintain visual consistency for a character across multiple generated video clips. Triggers on: "character setting", "character consistency", "角色設定", "設定角色", "角色一致性", "建立角色".
version: "1.0.0"
allowed-tools: []
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 角色視覺一致性設定技能 (Character Visual Consistency Setting Skill)

此技能提供了一套系統化的角色設定框架，用以定義、設計並維持角色在 AI 影片生成過程中的視覺一致性 (visual consistency)。

## 核心處理原則 (Core Principles)

### 1. 角色特徵維度 (Character Feature Dimensions)

定義一個角色時，應從以下五個關鍵維度進行具體描述，並將其轉化為明確的英文特徵：
- 面部特徵 (Face)：包含年齡、眼色、五官輪廓、有無雀斑、鬍鬚或疤痕。
  - 範例：`a 25-year-old Caucasian woman with blue eyes and sharp jawline`
- 髮型與髮色 (Hair)：包含長度、樣式、顏色與質地。
  - 範例：`short black bob haircut with straight bangs`
- 服裝細節 (Clothing)：包含款式、顏色、材質與花紋。避免籠統的服裝詞彙。
  - 範例：`wearing a beige cotton trench coat over a white crewneck t-shirt`
- 體型與特徵 (Body)：包含身高、體態。
  - 範例：`slender build, athletic stance`
- 固定配飾 (Accessories)：如眼鏡、帽子、耳環或項鍊。
  - 範例：`thin round gold-rimmed glasses`

### 2. 模型提示詞撰寫規範 (Writing Rules for Prompts)

- 避免使用主觀形容詞：如 `beautiful`、`handsome`，因為每個模型對這些詞的理解落差極大，且容易在不同鏡頭中產生大幅度變形。
- 採用具體特徵：以 `emerald green eyes, symmetrical facial features` 代替 `beautiful eyes`。
- 善用錨點物件：給角色一個非常獨特且容易辨識的錨點配件，例如 `red canvas backpack` 或 `bright yellow scarf`，這能強烈引導模型識別為同一個角色。

---

## 輸出格式範本 (Output Format Template)

請依以下結構輸出角色設定：

```markdown
# 角色設定：[角色名稱]

## 角色特徵表 (Character Attributes)
- 性別與年齡：[例如 女性，約 30 歲]
- 面部：[面部細節英文描述]
- 髮型：[髮型細節英文描述]
- 服裝：[服裝細節英文描述]
- 配飾：[配飾細節英文描述]

## 一致性提示詞區塊 (Consistency Prompt Block)
以下為已完成優化的英文一致性提示詞區塊，可直接複製至 scene-to-video-prompt 的主體描述中：

`[角色外觀英文描述段落，例如：A 30-year-old Asian woman with long wavy dark brown hair, soft facial features, wearing thin round gold-rimmed glasses, a beige linen shirt, and a bright yellow scarf]`

## 模型套用建議 (Model Integration Tips)
- Seedance 2.0：此模型對服裝和配件的細節捕捉非常敏感，若有自訂特徵，可將此 Prompt 放入主體描述。
- Kling 3.0：如果角色變形，建議簡化服裝花紋，改用單一顏色與純色布料描述。
```
