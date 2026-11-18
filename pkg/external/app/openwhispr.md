# openwhispr

來源: [OpenWhispr/openwhispr](https://github.com/OpenWhispr/openwhispr) (MIT, ~11.8k stars) · [website](https://openwhispr.com) · [docs](https://docs.openwhispr.com)

開源免費的語音轉文字 (voice-to-text) 桌面 app, 定位是 WisprFlow 與 Granola 的替代品: 按 hotkey 說話, 文字直接貼到游標位置; 另有 AI agent, 會議逐字稿與筆記. 主打 privacy-first: 可完全離線用本機模型 (音訊不離開裝置), 也可走雲端 BYOK (Bring Your Own Key, 自備 API key). 無 telemetry. 支援 macOS / Windows / Linux, 另有 Expo 手機版 (`openwhispr-mobile/`).

這個 repo `不是 Claude skill`: 沒有 `SKILL.md`, 是一個完整的 Electron app. 列在這裡是作為外部參考, 用於語音輸入工具選型, 或經由它的 MCP server 讓 AI assistant 存取筆記與逐字稿.

## 功能

| 功能                  | 說明                                                                                                                                        |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| Voice dictation       | 全域 hotkey, 在任何 app 中聽寫並自動貼上                                                                                                    |
| Dictation translation | 另一個 hotkey, 用 A 語言說, 貼上 B 語言                                                                                                     |
| Voice Assistant       | hotkey 直接把語音當指令送給 AI (GPT-5 / Claude / Gemini / Groq / OpenRouter / 本機模型); 選取文字可就地改寫, 可選擇附上螢幕截圖當 context   |
| Meeting transcription | 自動偵測 Zoom / Teams / FaceTime 通話, 即時 speaker diarization (分辨誰在說話) 與 voice fingerprint, 接 Google / Microsoft / Apple Calendar |
| Notes                 | 資料夾, semantic search, 雲端同步, AI actions; 登入後可用 team spaces 與網頁分享                                                            |
| Audio import          | 拖入音訊 / 影片, 批次上傳, 或貼 YouTube URL 轉錄                                                                                            |
| Public API & MCP      | 以 API 管理筆記與逐字稿, 或透過 [MCP server](https://docs.openwhispr.com/integrations/mcp) 接 AI assistant                                  |
| Enterprise            | SSO / SCIM, 集中管理 Amazon Bedrock 或 Azure OpenAI 存取                                                                                    |

## 技術

- Electron 41 + React 19 + TypeScript + Tailwind v4 + shadcn/ui, 本機儲存用 `better-sqlite3`.
- 本機 ASR: `whisper.cpp` (GPU 加速 Metal / CUDA / Vulkan), NVIDIA Parakeet 與 Cohere Transcribe 經 `sherpa-onnx`; 本機 LLM 用 `llama.cpp`.
- `examples/custom-asr-shim/`: 讓 self-hosted 的非 OpenAI-compatible ASR API 接進來.

## 安裝

- 使用者: 從 [releases](https://github.com/OpenWhispr/openwhispr/releases/latest) 下載 `.dmg` / `.exe` / `.AppImage` / `.deb` / `.rpm`, 或 `brew install --cask openwhispr`.
- 從原始碼: Node.js 24+, `npm install && npm run dev`.

## 適用情境

- 想在 Claude Code / terminal 用語音輸入 prompt, 且不想把音訊送雲端 (本機 Whisper / Parakeet).
- 會議錄音轉逐字稿並分辨講者, 再交給 AI 摘要.
- 透過 MCP 讓 Claude 讀寫 OpenWhispr 的筆記與逐字稿.

## 注意

- Intel Mac 沒有即時 speaker identification 與 voice fingerprint (ONNX Runtime 1.24 起停止提供 macOS x86_64 binary), semantic search 退回 keyword.
- 雲端同步與 team spaces 走 OpenWhispr Cloud (Neon Postgres); 隱私敏感內容應維持本機模式 (推論, 未逐項查證同步範圍).
