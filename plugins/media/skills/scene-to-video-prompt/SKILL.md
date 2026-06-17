---
name: scene-to-video-prompt
description: >
    Use when the user wants to convert a storyboard script into individual video prompts for AI video generators like Kling, Veo, Sora, or Seedance. Triggers on: "convert script to video prompts", "generate video prompt", "分鏡提示詞", "影片提示詞", "場景轉提示詞", "分鏡腳本".
version: "1.0.0"
allowed-tools: []
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 場景轉影片提示詞技能 (Scene to Video Prompt Skill)

此技能可將一份結構化的分鏡腳本，自動拆解並翻譯成符合 AI 影片生成模型規格的英文提示詞，特別對齊 Kling 等多分鏡 (multi-shot) Storyboard 介面。

## 核心處理原則 (Core Principles)

### 1. 分鏡切分原則 (Storyboard Breakdown Rules)

每當腳本中出現以下任一變化時，即應切分為新的鏡頭 (shot)：
- 場景變更 (Change of location)：例如從室內到室外。
- 景別變更 (Change of camera scale)：例如從遠景切到特寫。
- 時間跳躍 (Time jump)：例如從白天到黑夜。
- 主動作改變 (Major action shift)：例如從奔跑中突然坐下。

### 2. 提示詞六要素 (Six Elements of a Prompt)

為每個鏡頭組裝提示詞時，必須包含以下要素（並使用英文輸出，以獲得模型最佳理解度）：
- 主體 (Subject)：主角外觀的具體特徵描述。
- 動作 (Action)：該鏡頭主體所做的單一、清晰動作。
- 場景與環境 (Setting/Environment)：背景、周圍環境與材質細節。
- 鏡頭語言 (Camera/Shot type)：指定景別（如 close-up、medium shot）與運鏡（如 pan left、dolly in）。
- 光線與氛圍 (Lighting/Mood)：如 golden hour、neon reflections、moody lighting。
- 風格 (Style)：如 cinematic、shot on 35mm film、3D animation。

### 3. 一致性鎖定 (Consistency Locking)

跨鏡頭描述同一個角色或場景時，關鍵名詞與外觀細節描述（如髮型、服裝顏色、材質）必須逐字複製貼上，不可使用代名詞（如 he/she）或改變措辭，以防模型在生成時讓角色變臉。

---

## 輸出格式範本 (Output Format Template)

請依以下結構輸出分鏡提示詞：

```markdown
# 影片提示詞清單 (Video Prompt List)

## 一致性區塊 (Consistency Block)
- 角色特徵描述 (Character Consistency)：[請在此列出跨鏡頭必須複製的主角外觀英文描述]

## 分鏡清單 (Shot Prompts)

### 鏡頭 1 (Shot 1)
- 景別與運鏡 (Camera & Movement)：[例如 Medium shot, pan right]
- 英文提示詞 (Prompt)：[主體 + 動作 + 場景 + 鏡頭 + 光線 + 風格 的英文組裝]
- 轉場建議 (Transition)：[例如 None / Cut]

### 鏡頭 2 (Shot 2)
- 景別與運鏡 (Camera & Movement)：[例如 Close-up, static shot]
- 英文提示詞 (Prompt)：[複製的角色特徵描述 + 新的動作與細節]
- 轉場建議 (Transition)：[例如 Fade out]

---

## 模型特定優化建議 (Model-specific Tips)

- Kling 3.0：支援 3 到 12 個分鏡，若需要 4K，可在風格要素中加入 `extreme detail, 4K resolution`。
- Seedance 2.0：對角色一致性支援度極佳，確保一致性區塊的描述精確。
- Veo 3.1：若需產生同步對白，可在 Prompt 後方以括號註明對白文字。
```
