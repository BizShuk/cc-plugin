# 多媒體影片生成與劇本插件 (Media Generation and Script Plugin)

此插件提供了一套完整的影片生成 (video generation) 與劇本創作 (script writing) 技能。其工作流可協助使用者將簡單的想法擴展為有故事劇情的腳本，再自動拆解為多分鏡 (multi-shot) 的影片提示詞，並能為特定角色建立一致性外觀描述。

## 技能清單 (Skills List)

### 1. 提示詞轉故事腳本 (prompt-to-story-script)

- 路徑：`./skills/prompt-to-story-script`
- 用途：將簡單的一句話靈感或提示，優化並擴展為具備完整起承轉合與畫面化動作的故事腳本。

### 2. 場景轉影片提示詞 (scene-to-video-prompt)

- 路徑：`./skills/scene-to-video-prompt`
- 用途：將分鏡腳本拆分成多個獨立的分鏡 (storyboard)，並套用主體、動作、場景、鏡頭、光線、風格六要素，產生符合 Kling 或 Veo 等多鏡介面的英文提示詞。

### 3. 角色一致性設定 (character-setting)

- 路徑：`./skills/character-setting`
- 用途：建立角色視覺特徵的描述範本，以確保在不同鏡頭中角色的面部、服裝與髮型等細節能維持一致。

## 工作流 (Workflow)

一整個影片創作流程可以透過以下三個步驟串聯完成：

```
簡單想法 (Idea)
    │
    ▼ (prompt-to-story-script)
故事腳本 (Story Script)
    │
    ▼ (scene-to-video-prompt) + (character-setting 角色描述)
多分鏡影片提示詞 (Multi-shot Prompts)
```

1. 透過 `prompt-to-story-script` 將您的點子寫成故事。
2. 使用 `character-setting` 為故事主角設定具體且一致的外觀特徵。
3. 透過 `scene-to-video-prompt` 將劇本與角色設定結合成最終可直接貼入影片生成模型的提示詞清單。
