---
name: prompt-to-story-script
description: >
    Use when the user wants to generate or optimize a video script or story from a simple idea or prompt. Triggers on: "generate script", "write a story script", "優化腳本", "產生劇本", "故事腳本", "腳本創作".
version: "1.0.0"
allowed-tools: []
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 提示詞轉故事腳本技能 (Prompt to Story Script Skill)

此技能可將一句簡單的點子或提示詞，擴展並優化為具備情節張力、故事弧線與視覺化畫面描述的短片腳本。

## 核心處理原則 (Core Principles)

### 1. 定位故事錨點 (Establish Story Anchors)

開始寫作前，必須明確以下四個元素：
- 目標長度 (Target Duration)：例如 15 秒、30 秒或 60 秒。
- 類型與情緒 (Genre & Tone)：例如科幻、懸疑、溫馨、商業廣告。
- 主角與渴望 (Protagonist & Desire)：主角是誰？他想要達到什麼目的？
- 關鍵轉折 (The Shift)：故事中發生的核心轉折事件或情緒變化。

### 2. 故事弧線設計 (Story Arc Design)

依據片長採用不同的節奏結構：
- 15 秒短片：採用三拍結構（鋪陳 -> 轉折 -> 收尾 / Setup -> Turn -> Payoff）。
- 30 秒至 60 秒中長片：展開至五拍（鋪陳 -> 觸發 -> 衝突 -> 轉折 -> 結局）。

### 3. 畫面化動作改寫 (Visual Action Translation)

影片是視覺媒介，避免使用抽象的內心戲或情緒形容詞，應將其改寫為鏡頭看得見的畫面與動作：
- 錯誤範例：`他感到非常寂寞`
- 正確範例：`他獨自坐在咖啡廳角落，看著對街共撐一把傘的情侶，默默地將大衣領子拉緊`

---

## 輸出格式範本 (Output Format Template)

請依以下結構輸出腳本：

```markdown
# 故事標題：[請填寫標題]

## 故事設定 (Setup Details)
- 片長：[例如 15 秒]
- 類型：[例如 溫馨]
- 故事簡介 (Logline)：[一句話概括故事]

## 故事大綱與轉折 (Story Arc)
- 起 (Setup)：[開場畫面與狀態]
- 承/轉 (Turn)：[發生了什麼改變]
- 合 (Payoff)：[最後的結局與情緒落點]

## 逐鏡腳本 (Script Breakdown)

### 鏡頭 1：[起]
- 畫面描述 (Visual)：[視覺細節]
- 角色動作 (Action)：[主角的具體動作]
- 音效與音樂 (Audio)：[背景音或配樂建議]

### 鏡頭 2：[承/轉]
- 畫面描述 (Visual)：[改變發生的視覺信號]
- 角色動作 (Action)：[主角對轉折的反應]
- 音效與音樂 (Audio)：[音效變化]

### 鏡頭 3：[合]
- 畫面描述 (Visual)：[最終的情緒落點畫面]
- 角色動作 (Action)：[結尾動作]
- 音效與音樂 (Audio)：[音樂收尾]

## 加分提示 (Punch-up Tips)
- [給予一到兩點關於如何提升這段影片張力或表現力的具體建議]
```
