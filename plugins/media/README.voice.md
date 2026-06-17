# 語音模型與 VoxCPM 技術指南 (Voice Models & VoxCPM Technical Guide)

這份文件介紹了 AI 語音模型專屬化 (Voice Specialization) 的多個維度，比較了當前市場上的主流語音模型，並針對開源語音模型 `VoxCPM` 進行深入剖析與聲音克隆 (Voice Cloning) 指南。

> [!NOTE]
> 本文件旨在協助開發者理解如何透過參數調校、提示詞與參考音訊，讓 AI 語音達到最理想的呈現效果。

---

## 結論與核心建議 (Key Takeaway)

- `多元的專屬化面向`：AI 語音的調整可分為身份層 (Identity Layer)、表現層 (Expression Layer) 與技術層 (Technical Layer)，著力點眾多。
- `市場模型分工明確`：ElevenLabs 著重音質與成熟生態；OpenAI 擁有最強的文字指令可控性；Cartesia 專攻超低延遲；而 `VoxCPM` 則是開源自主部署且中英雙語表現優異的代表。
- `克隆成敗取決於音訊`：聲音克隆相似度有 80% 決定於參考音訊 (Reference Audio) 的乾淨度與品質，文字提示詞 (Style Prompt) 僅用於風格微調。

---

## 一、語音專屬化的四大介面方式 (Four Interface Modes for Specialization)

調整語音模型時，根據不同模型所支援的介面，有以下四種常見的控制方式：

1. `文字描述式 (Prompt-based / Descriptive)`：直接使用自然語言描述目標聲音。例如：「一位溫暖、四十歲左右的男聲，語速偏慢、帶點疲憊感」。適用於 OpenAI 等支援自然語言指令的模型。
2. `參數式 (Parameter-based)`：透過數值或滑桿控制聲音屬性，例如 stability (穩定度)、similarity (相似度)、speed (語速) 等。ElevenLabs 是典型代表。
3. `參考音訊式 (Reference Audio / Voice Cloning)`：提供一段聲音樣本，讓模型模仿該音色與說話風格。
4. `標記語言式 (Markup / SSML)`：在文字中插入特定標籤（例如 `<break time="500ms"/>` 或 `<emphasis>`）以精確控制停頓、重音與發音。Google 與 Azure TTS 大量使用。

---

## 二、語音專屬化的十二個面向 (Twelve Dimensions of Voice Specialization)

語音專屬化可以拆分為三個層次，包含共 12 個可調面向：

### A. 身份層 (Identity Layer) — 配音員是誰

1. `音色與聲音身份 (Timbre / Voice Identity)`：選用預設音色，或透過聲音克隆複製特定人聲。
2. `語言與口音 (Language & Accent)`：調整同音色在不同語言下的腔調，例如美式、英式英文，或台灣腔、大陸腔中文。
3. `性別與年齡感 (Gender & Age)`：指定語音呈現的性別與年齡特徵。

### B. 表現層 (Expression Layer) — 怎麼演繹

4. `情緒與語氣 (Emotion & Tone)`：例如開心、悲傷、興奮、冷靜、嚴肅等。
5. `語速 (Speed / Rate)`：控制發音的快慢。
6. `音高 (Pitch)`：調整聲音的高低音調。
7. `音量 (Volume)`：控制整體的響度。
8. `停頓與節奏 (Pauses & Pacing)`：利用標記語言控制停頓，製造呼吸感與說話節奏。
9. `重音與強調 (Emphasis / Stress)`：加強特定字詞的發音強度。
10. `說話風格 (Style)`：例如旁白 (Narration)、對話 (Conversational)、新聞播報 (Newscasting) 等。

### C. 技術層 (Technical Layer) — 錄音規格與演算法

11. `穩定度與多樣性 (Stability vs Variability)`：調高穩定度會使每次輸出更一致但較單調；調低則情感起伏更大，但可能偶爾失真。
12. `發音修正 (Pronunciation)`：透過發音標籤 (Phoneme) 或自訂字典，修正縮寫、人名或專有名詞的唸法。
13. `取樣率、格式與種子值 (Sample Rate, Format & Seed)`：固定種子值 (Seed) 可重現相同的語音輸出，方便除錯與微調。

---

## 三、當前主流語音模型橫向比較 (Comparison of Major Voice Models)

| 模型名稱            | 類型          | 中英文支援 | 最大強項                             | 大致計費方式        | 適用場景                         |
| ------------------- | ------------- | ---------- | ------------------------------------ | ------------------- | -------------------------------- |
| `VoxCPM / VoxCPM2`  | 開源          | 極強       | 自然度高、可本地部署、中英雙語佳     | 免費 (需自備算力)   | 隱私自主性要求高、中英混用內容   |
| `ElevenLabs`        | 閉源          | 優秀       | 極致音質、聲音克隆、成熟的開發生態   | 約 $0.30 / 1K 字元  | 有聲書、高規格旁白、品質至上者   |
| `OpenAI TTS`        | 閉源          | 良好       | 文字指令可控性高、價格實惠、易於整合 | $15 / 每百萬字元    | 已在 GPT 生態系、追求簡單快速者  |
| `Fish Audio`        | 閉源 / 可自架 | 強         | 豐富的情緒控制、多國語言、價格具優勢 | 約 $15 / 每百萬字元 | 需要多語系與細膩情緒控制的場景   |
| `Cartesia Sonic`    | 閉源          | 良好       | 超低延遲 (TTFT < 100ms)              | 按 API 使用量計費   | 即時語音代理、對話機器人         |
| `Voxtral (Mistral)` | 開源          | 良好       | 開放權重、品質高                     | 免費 (需自備算力)   | 需自行部署但要求商業級品質者     |
| `Google / Azure`    | 閉源雲端      | 優秀       | 服務穩定、企業級規模化、提供免費額度 | 階梯式計費          | 企業級大量語音合成、需 SSML 精修 |

---

## 四、開源模型 VoxCPM 深度解析 (In-depth Analysis of VoxCPM)

`VoxCPM` 屬於開源、可自主部署的語音模型，在保留隱私與降低長期授權成本上具有極大優勢。

### 核心架構與創新

- `免分詞器架構 (Tokenizer-free Architecture)`：傳統 Text-to-Speech (TTS) 模型會將語音切成離散的語音 token，而 VoxCPM 直接在連續的語音空間建模，保留了語音原本的自然流暢度。
- `語意與聲學解耦 (Semantic-Acoustic Decoupling)`：建立在 MiniCPM-4 骨幹上，透過階層式語言建模 (Hierarchical Language Modeling) 與 FSQ 約束，提升了聲音的表現力與發音穩定性。
- `自動語意推斷`：模型會自動根據輸入文字的上下文語意推斷出合適的語氣與風格，減少手動下指令的繁瑣，但相對的可控度較低。

### 效能表現與版本

- `VoxCPM-0.5B`：採用 180 萬小時雙語語料訓練，提供 Zero-shot 語音生成。在 SEED-TTS-EVAL 評測中，英文 WER 為 1.85%，中文 CER 為 0.93%。在 NVIDIA RTX 4090 上即時因子 (Real-Time Factor, RTF) 低至 0.17。
- `VoxCPM2`：為最新主要版本，擴展至 2B 參數量，使用超過 200 萬小時多語語料訓練，支援 30 種語言，具備 `聲音設計 (Voice Design)` 與 `可控聲音克隆 (Controllable Voice Cloning)` 功能，最高可輸出 48kHz 錄音室等級的音訊。

> [!CAUTION]
> `安全與倫理風險`：VoxCPM 強大的 Zero-shot 聲音克隆能力可能會被濫用於製作 Deepfake 語音、詐騙或散播假訊息。使用者在部署與使用時必須自行承擔倫理與法律責任，建議清楚標明為 AI 生成。

---

## 五、聲音克隆提示詞指南 (Voice Cloning & Prompting Guide)

在 `VoxCPM` 中，模仿一個真實的聲音必須使用 `聲音克隆 (Voice Cloning)`，這需要提供一段參考音訊 (Reference Audio)，而非單純使用文字提示詞 (Prompt)。

### 聲音克隆的三種模式

1. `基本克隆 (Basic Cloning)`：僅複製參考音訊的音色。

    ```python
    wav = model.generate(
        text="這是使用我的參考音訊克隆出來的聲音。",
        reference_wav_path="path/to/voice.wav"
    )
    ```

2. `可控克隆 (Controllable Cloning)`：複製音色，並在文字開頭加上括號風格提示詞，微調語速或情緒。

    ```python
    wav = model.generate(
        text="(語速稍快、開朗的語氣)這是帶風格控制的克隆聲音。",
        reference_wav_path="path/to/voice.wav",
        cfg_value=2.0,
        inference_timesteps=10
    )
    ```

    可用的風格提示詞範例：`(平靜、緩慢)`、`(興奮、語速快)`、`(溫柔、像在說悄悄話)`。

3. `極致克隆 (Ultimate Cloning)`：同時提供參考音訊與其對應的精確逐字稿 (Transcript)，以達到最高的相似度與對齊效果。官方建議將同一段音訊同時傳給 `prompt_wav_path` 與 `reference_wav_path`。

    ```python
    wav = model.generate(
        text="這是極致克隆的示範。",
        prompt_wav_path="path/to/voice.wav",
        prompt_text="參考音訊中人物實際說出的逐字稿，必須一字不差。",
        reference_wav_path="path/to/voice.wav"
    )
    ```

### 參考音訊的四大準備要點

聲音克隆的成敗有 80% 取決於參考音訊的品質，請遵循以下原則：

- `乾淨無雜音`：不可有背景音樂、環境噪音、回音或其他人的聲音。
- `長度適中`：建議提供 10 至 30 秒的清晰說話片段。
- `語氣自然`：選擇語調平穩、具有代表性的日常說話內容，避免極端的情緒波動。
- `單一說話者`：整段音訊中只能有目標模仿對象一個人的聲音。

---

## 六、實作與測試步驟 (Execution Plan)

1. `部署環境`：請參考 [VocCPM_setup.md](./VocCPM_setup.md) 進行本機或 Mac mini 部署。
2. `測試樣板`：完成部署後，可參考 [voice/templates/README.md](./voice/templates/README.md) 使用 10 種預設聲音樣板進行測試。
3. `對照測試 (A/B Testing)`：使用同一段測試腳本，分別使用 `Voice Design` 與不同引導強度 (`cfg_value`) 生成，盲聽比較並挑選出最適合的聲音配置。
