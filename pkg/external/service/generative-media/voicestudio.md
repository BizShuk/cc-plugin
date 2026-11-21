# voicestudio

來源: [debpalash/VoiceStudio](https://github.com/debpalash/VoiceStudio) (AGPL-3.0, ~113k stars, 最新 release v0.5.8 @ 2026-11-12) · [website](https://voicestudio.sh) · [docs](https://github.com/debpalash/VoiceStudio/blob/main/docs/README.md)

開源, 完全本機執行的 ElevenLabs 替代品: voice cloning (語音複製), voice design (用文字描述生成聲音), video dubbing (影片配音翻譯), dictation (聽寫), transcription (逐字稿) 與 audiobook 製作, 號稱支援 646 種語言. 預設 TTS 引擎是 k2-fsa/OmniVoice, 模型由 app 內的 Model Manager 按需下載. 免費版不需帳號也不需雲端; 付費的 Pro (商用授權, 多裝置 GPU 派工) 與 Cloud (hosted API, 尚未開放) 為選配.

這個 repo `本身是 app`, 但 `附帶兩個 Agent skill`: `skills/voicestudio/SKILL.md` (操作音訊) 與 `skills/voicestudio-maintainer/SKILL.md` (維護 repo). 列在這裡作為外部參考, 用於讓 Claude 透過本機 REST API 或 MCP server 產生語音, 轉錄音訊, 自動化配音與旁白.

## 功能

| 功能 | 說明 |
| ---- | ---- |
| Voice cloning | 給一段乾淨的 reference 錄音, 以該聲音唸出任意文字 |
| Voice design | 用文字描述 (性別, 年齡, 音高, 風格, 口音) 生成新聲音並存成 profile |
| Video dubbing | 翻譯並配音影片, 語音對齊原片時間軸; 支援多講者, 保留背景音 |
| Dictation & transcription | 浮動聽寫 widget, 本機轉錄, speaker diarization, vocal isolation |
| Narrate & batch | 故事 / audiobook, batch queue, watch folder, 可重用 recipe |
| Integrations | 本機 OpenAI-compatible API, WebSocket 即時轉錄, MCP server, CLI |
| AI watermark | 生成音訊加浮水印 |

引擎 (依硬體與已安裝模型而定):

- TTS: OmniVoice (預設, 另有 GGUF / subprocess 版), CosyVoice 3, IndexTTS 2.5, VoxCPM2, MOSS-TTS, GPT-SoVITS, KittenTTS, MLX-Audio, sherpa-onnx 等 20 餘種.
- ASR: WhisperX (預設), Faster-Whisper, MLX Whisper, whisper.cpp, Parakeet TDT, Moonshine, FunASR, sherpa-onnx (即時聽寫), 或外接 OpenAI-compatible server.

## 技術

- 桌面: Electron (由 0.5.3 以前的 Tauri 版遷移) + `frontend/`; 後端: Python FastAPI (`backend/`, Alembic migration), 預設聽 `:3900`.
- 加速: Apple Silicon 走 Metal (MPS), NVIDIA CUDA, AMD ROCm (opt-in), GGUF 引擎可走 Vulkan; 沒 GPU 也能用 CPU, 只是慢.
- 平台: macOS 13.3+ (Apple Silicon; Intel 只能接遠端 backend), Windows x64, Linux x64, Docker (GHCR), Android 10+ / iOS 17+ (preview, 下載模型後離線).

## Agent 介面

| 介面 | Endpoint | 用途 |
| ---- | -------- | ---- |
| Health / schema | `GET :3900/health`, `GET :3900/openapi.json` | 先確認 backend 在線, 再以 OpenAPI 探索該版本的 native endpoint |
| OpenAI-compatible TTS | `POST :3900/v1/audio/speech` | `model=tts-1`, `voice=<profile_id>`; `alloy` 等名稱只是相容 alias |
| OpenAI-compatible ASR | `POST :3900/v1/audio/transcriptions` | `response_format` 可選 `json` / `text` / `verbose_json` / `srt` / `vtt` |
| Voices | `GET :3900/v1/audio/voices` | 列出已存 profile 與引擎, 不要自己捏 profile ID |
| WebSocket 即時轉錄 | `ws://:3900/v1/audio/transcriptions/stream` | 送 WebM/Opus 或 `?pcm=1&sr=16000`, 回 partial / final |
| MCP (Streamable HTTP) | `http://localhost:3900/mcp/` | backend 開著就有, 不需另起服務; 保留結尾斜線最相容 |
| MCP (stdio) | `python -m backend.mcp_shim` | 只支援 stdio 的 client |
| Capabilities | `GET :3900/.well-known/voicestudio-speech` | 回傳 `voicestudio.speech.v1` 能力文件 |

MCP tools: `generate_speech`, `clone_voice`, `describe_voice`, `design_voice`, `transcribe`, `list_voices`, `list_personalities`, `list_languages`, `check_health`.

MCP 相關環境變數 (設在 backend 端):

- `OMNIVOICE_MCP_OUTPUT_MODE`: `resources` (預設, 回 `wav_base64` 塞進 context) / `files` / `both`. 給 agent 用建議 `files`, 避免 base64 音訊吃掉 token.
- `OMNIVOICE_MCP_BASE_PATH`: 檔案讀寫的 security boundary; 未設時 `audio_path` / `ref_audio_path` 參數一律拒絕.
- `OMNIVOICE_MCP_TIMEOUT_S`: 固定 timeout; 未設時跟隨 backend 的生成 / 轉錄 timeout.

Claude Code 接法 (`.mcp.json`, 桌面 app 的 Integrations 頁也能匯出):

```json
{
  "mcpServers": {
    "voicestudio": { "type": "http", "url": "http://127.0.0.1:3900/mcp/" }
  }
}
```

安裝它自帶的 skill: `npx skills add debpalash/VoiceStudio`.

## 安裝

- macOS: `brew tap debpalash/tap && brew install --cask debpalash/tap/voicestudio`, 或從 [releases](https://github.com/debpalash/VoiceStudio/releases/latest) 下載 DMG.
- Docker (CPU):

  ```sh
  export OMNIVOICE_API_KEY="$(openssl rand -hex 32)"
  docker run -d --name omnivoice -p 127.0.0.1:3900:3900 \
    -e OMNIVOICE_API_KEY -e OMNIVOICE_DATA_DIR=/app/omnivoice_data \
    -v omnivoice-data:/app/omnivoice_data ghcr.io/debpalash/voicestudio:stable
  ```

- 從原始碼: `git clone` 後 `bun run dev` (會一併裝 Python 依賴並啟動 Electron).

## 適用情境

- 讓 Claude 以固定聲音產生旁白, podcast 片段或 audiobook, 全程不送雲端.
- 影片翻譯配音 (多講者, 對齊時間軸), 或批次把會議錄音轉成 `srt` 字幕.
- 既有使用 OpenAI audio API 的 script 改指向 `:3900`, 換成本機模型.
- 與 [openwhispr](../../app/openwhispr.md) 比較: OpenWhispr 專注聽寫與會議筆記; VoiceStudio 重心在 TTS, cloning 與 dubbing, 聽寫只是附帶功能.

## 注意

- 首次 `generate_speech` 若尚未下載 OmniVoice 模型, 會在 backend 內下載約 2.3 GB, 容易 timeout; 先在 app 的 Model Catalogue 裝好模型.
- AGPL-3.0: 若把它改造後以網路服務提供給他人, 須公開修改後的原始碼. 模型權重與生成音訊各有授權, 以 Model Manager 顯示為準.
- 只能 clone 取得本人授權的聲音; 官方明文禁止未授權模仿, 詐騙等用途.
- loopback 連線預設免認證; 對外開放時必須設 `OMNIVOICE_API_KEY` 並走 HTTPS/WSS.
- 原生聽寫控制 port `3902` (Tauri 時代的 `/v1/dictation/*`) 已移除, 舊教學的該段不再適用.
