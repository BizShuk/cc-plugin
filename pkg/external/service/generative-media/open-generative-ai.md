# open-generative-ai

來源: [Anil-matcha/Open-Generative-AI](https://github.com/Anil-matcha/Open-Generative-AI) (MIT, 約 35k stars) · [hosted 版](https://muapi.ai/open-generative-ai)

一個可 self-host 的 AI 圖片 / 影片生成 studio (Next.js web + Electron desktop), 本質是 [MuAPI](https://muapi.ai) 這個 model gateway 的前端. 把數百個模型 (Flux, Nano Banana, Midjourney, Kling, Sora, Veo, Seedance, Wan...) 包成 14 個 studio 分頁, 用同一套 UI 操作. 它`不是 Claude skill`: repo 裡沒有 `SKILL.md`; 給 agent 用的部分在另一個 repo [muapi-skills](https://github.com/SamurAIGPT/muapi-skills).

## 內容

| 部分                                                            | 用途                                                                                                                                                   |
| --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `app/`                                                          | Next.js App Router; `/studio` 是主頁, `/api/*` 把請求 proxy 到 MuAPI                                                                                   |
| `packages/studio/`                                              | 共用 React component library, 每個 studio 一個 `*Studio.jsx`; `models.js` 是模型定義的單一來源, 各家另有 `<vendor>Models.js` / `<vendor>Parameters.js` |
| `packages/studio/src/muapi.js`                                  | API client: submit 後 poll 直到完成                                                                                                                    |
| `electron/`                                                     | desktop app; `lib/localInference*` 管 sd.cpp, `lib/wan2gp*` 管 Wan2GP 本機或 MCP 推論                                                                  |
| `packages/Vibe-Workflow`, `Open-Poe-AI`, `Open-AI-Design-Agent` | git submodule, 分別提供 Workflow builder, Agent, Design Agent 三個 studio                                                                              |

Studio 分頁: Image, Video, Audio, AI Clipping, Vibe Motion, Lip Sync, Recast (body swap), Cinema, Marketing, Workflow, Agent, Design Agent, Apps, AI Influencer, MCP & CLI.

## 運作方式

雲端生成 (預設):

1. 使用者在 UI 填入 MuAPI access key, 存在瀏覽器 `localStorage`.
2. `POST /api/v1/{model-endpoint}` 送出 prompt 與參數, header 為 `x-api-key`.
3. `GET /api/v1/predictions/{request_id}/result` 輪詢到 `completed`.
4. 參考圖先 `POST /api/v1/upload_file` 取得 hosted URL; 多圖模型一次送 `images_list` (最多 14 張).

本機推論 (只有 desktop app):

| Engine          | 說明                                                            | 適用                                                            |
| --------------- | --------------------------------------------------------------- | --------------------------------------------------------------- |
| `sd.cpp` (內建) | stable-diffusion.cpp, Apple Silicon 走 Metal                    | 圖片: SD 1.5, SDXL, Z-Image; 8 GB Mac 只建議 SD 1.5             |
| `Wan2GP` (自備) | 本機跑 `wgp.py`, 或連另一台 GPU 機器的 MCP server (`:7866/mcp`) | 影片 (Wan 2.1/2.2, Hunyuan, LTX-2) 與大型圖片模型; 需 CUDA/ROCm |

Wan2GP 的 MCP 模式讓 Mac 只跑 UI, 推論丟給 LAN 上的 GPU 機器.

## 值得借鏡的設計

- 模型目錄集中在一份資料檔, UI 依模型能力動態顯示 aspect ratio, 解析度, 時長控制; 有無參考圖自動切換 t2i/i2i 或 t2v/i2v 模型集.
- Cinema Studio 把相機, 鏡頭, 焦段, 光圈選項轉成 prompt modifier, 是「結構化選項 -> prompt」的好例子.
- Pending job 存在本機, 重新整理後自動續 poll.

## 注意事項

- `No content filters` 是行銷說法: 過濾與否取決於 MuAPI 與底層模型供應商, 這個 repo 只是前端.
- `Self-hosted, data stays local` 只對本機推論成立; 雲端生成的 prompt 與上傳圖都會送到 MuAPI. `推論`: 雲端生成依 MuAPI credit 計費, README 的 `Free` 指的是程式碼.
- README 大量篇幅在推廣 MuAPI White Label 付費方案.
- desktop 安裝檔未簽章 (macOS 需 `xattr -cr`, Windows 會跳 SmartScreen).

## 安裝與試跑

```sh
git clone --recurse-submodules https://github.com/Anil-matcha/Open-Generative-AI.git
cd Open-Generative-AI
npm run setup          # 必要: init submodule + install + build workspace packages
npm run dev            # web, http://localhost:3000
npm run electron:dev   # desktop, 含本機推論
docker compose up --build   # http://localhost:3001
```

需求: Node 18 以上, MuAPI access key (只用本機模型可略過).

## 適用情境

- 想要一個介面比較多家圖片 / 影片模型的輸出, 而不想逐家註冊.
- 想參考多模型 generation UI 的結構: 模型目錄, 參數控制, 非同步 job 輪詢.
- 想在 LAN 內用自己的 GPU 跑 Wan2GP, 由 Mac 操作.
- 在 Claude Code 內直接生成媒體, 改看 [muapi-skills](https://github.com/SamurAIGPT/muapi-skills) 或 MuAPI 的 MCP endpoint; 只要生圖的話, 現有的 `imagine` skill 走自家 proxy 即可.
