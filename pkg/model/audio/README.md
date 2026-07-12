# 語音模型清單 (Audio Model List)

針對中英雙語 (Code-Switching) 場景推薦的語音模型：

## 文字轉語音 (Text to Speech)

| 分類                       | 模型          | 特色                                                       |
| :------------------------- | :------------ | :--------------------------------------------------------- |
| `品質首選 (The Best)`      | `CosyVoice 2` | 情感過渡自然，支援零樣本聲音複製 (Zero-Shot Voice Cloning) |
|                            | `Fish Speech` | 發音地道，適合科技專有名詞與專業術語                       |
| `高效輕量 (The Efficient)` | `F5-TTS`      | 顯存占用低 (Low VRAM Footprint)，無吞字漏字問題            |
|                            | `MeloTTS`     | CPU 執行友善，即時推理 (CPU Real-Time Inference) 成本低    |

## 語音轉文字 (Speech to Text)

| 分類                       | 模型                     | 特色                                                            |
| :------------------------- | :----------------------- | :-------------------------------------------------------------- |
| `品質首選 (The Best)`      | `Whisper Large V3`       | 跨語言語境理解 (Cross-Lingual Contextual Awareness)，抗噪能力強 |
|                            | `SenseVoice`             | 支援富文本識別 (Rich Text Features)，包含情緒與環境音偵測       |
| `高效輕量 (The Efficient)` | `SenseVoice-Small`       | 推理速度極快，適合即時對話系統                                  |
|                            | `Whisper Turbo / Distil` | 兼顧精準度與伺服器成本，吞吐量提升                              |

## 實戰組合 (Workflow Strategy)

- `即時助理 (Real-Time Agent)`：`SenseVoice-Small` (STT) + `CosyVoice 2` (TTS) ── 主打低延遲 (Low Latency) 互動。
- `影音後製 (Offline Pipeline)`：`Whisper Large V3` (STT) + `Fish Speech` (TTS) ── 主打最高精準度與自然聽感。
