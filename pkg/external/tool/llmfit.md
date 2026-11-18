# llmfit

來源: [AlexsJones/llmfit](https://github.com/AlexsJones/llmfit) (MIT, ~39k stars, Rust, 最新 release `v1.1.20` 2026-11-03) · [website](https://llmfit.axjns.dev) · [API.md](https://github.com/AlexsJones/llmfit/blob/main/API.md)

一個 terminal 工具, 回答「我這台機器能順跑哪些 open-source LLM」. 偵測 CPU, RAM, GPU / VRAM 與加速 backend (CUDA / Metal / ROCm / SYCL / Ascend), 再把內建 catalog 中數百個模型逐一以 quality, speed, fit, context 四個維度打分, 並自動挑出放得下的最高品質 quantization. 預設開 TUI (Terminal UI, 終端互動介面), 也有 CLI, Web dashboard, REST API 與 MCP server.

這個 repo `不是 Claude skill`: 根目錄沒有 `SKILL.md`, 是一個 Rust workspace. 內附的 [`skills/llmfit-advisor`](https://github.com/AlexsJones/llmfit/tree/main/skills/llmfit-advisor) 是給 OpenClaw agent 用的 skill. 列在這裡作為外部參考, 用於本機模型選型, 或經由 MCP / `--json` 讓 agent 查詢「這台機器能跑什麼」.

## 功能

| 功能                    | 說明                                                                                                                         |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Hardware detection      | `nvidia-smi` (多 GPU 合計 VRAM), `rocm-smi`, Intel Arc sysfs, Apple Silicon unified memory, `npu-smi`                        |
| Dynamic quantization    | 從 Q8_0 往下走到 Q2_K, 挑放得下的最高品質; 全 context 放不下再試半 context. 支援 GGUF / AWQ / GPTQ / EXL2 / MLX              |
| Scoring                 | Quality / Speed / Fit / Context 各 0-100, 依 use case (General, Coding, Reasoning, Chat, Multimodal, Embedding) 加權成總分    |
| Speed estimation        | 以記憶體頻寬估 tok/s: `(bandwidth_GB_s / model_size_GB) × 0.55`; 約 80 張已知 GPU 查表, 未知則用 per-backend 常數           |
| MoE                     | 以 active parameters 估 VRAM 與速度 (例: Mixtral 8x7B 由 23.9 GB 降到 ~6.6 GB)                                               |
| Hardware simulation     | `--memory` / `--ram` / `--cpu-cores` / `--profile <name\|file.json>` 模擬別台或未發售的硬體                                  |
| Plan                    | `llmfit plan <model> --context 8192 --target-tps 25` 反推跑某模型需要的硬體                                                  |
| Runtime providers       | 偵測 Ollama, llama.cpp, MLX, Docker Model Runner, LM Studio 已裝的模型, TUI 按 `d` 直接下載                                  |
| Benchmark               | `llmfit bench` 量實際 tok/s 與 TTFT, 實測值取代估計值; `--share` 從 TUI 直接把結果送 PR 回 community 資料庫                  |
| Cluster mode (preview)  | 讀 Ray Dashboard 產生 `cluster.toml`, 判斷模型能否跨多節點切分給 vLLM; 只讀不排程                                             |
| Storage                 | `llmfit storage --keep 3` 估算保留 N 個模型需要的 SSD 容量                                                                    |

## Agent 介面

- CLI: 任何子命令加 `--json`, 例如 `llmfit recommend --json --use-case coding --limit 3`, `llmfit --json system`, `llmfit info "<model>"`.
- REST: `llmfit serve` (預設 `127.0.0.1:8787`), 端點 `/health`, `/api/v1/system`, `/api/v1/models`, 並附 Web dashboard; 可改用 `--unix-socket <path>`.
- MCP (stdio): `llmfit serve --mcp`, 提供 6 個 tool: `get_system_specs`, `recommend_models`, `search_models`, `plan_hardware`, `get_runtimes`, `get_installed_models`.

Claude Code 接法 (`.mcp.json`, 推論, 未實測):

```json
{
  "mcpServers": {
    "llmfit": { "command": "llmfit", "args": ["serve", "--mcp"] }
  }
}
```

## 安裝

- macOS / Linux: `brew install AlexsJones/llmfit/llmfit`, 或 `curl -fsSL https://llmfit.axjns.dev/install.sh | sh` (加 `-s -- --local` 裝到 `~/.local/bin`).
- Python: `uv tool install -U llmfit` 或 `uvx llmfit`.
- Windows: `winget install --exact --id AlexsJones.llmfit` / `scoop install llmfit`.
- Container: `docker run -d -p 8787:8787 ghcr.io/alexsjones/llmfit serve`.
- 從原始碼: `cargo build --release`; 要內嵌 dashboard 需先在 `llmfit-web/` 跑 `npm ci && npm run build`.

## 適用情境

- 在 Apple Silicon 開發機上挑 MLX / GGUF 模型給 pm2 常駐的本機 LLM server 前, 先確認 fit 與預估 tok/s, 避免把記憶體吃爆.
- 採購或升級硬體前, 用 `--profile` / `plan` 模擬「換成 64GB / 某張 GPU 能跑到哪個模型」.
- 讓 agent 透過 MCP 自己查本機能力, 再決定要叫本機模型還是雲端 API.

## 注意

- 速度與 fit 是`估計值`, 依頻寬模型推算; 要可靠數字請跑 `llmfit bench` 實測, 或看 `llmfit info` 的 `estimate_basis` 確認假設.
- 模型 catalog 在編譯時內嵌 (`llmfit-core/data/hf_models.json`), 新模型要等新 release, 或依 [custom models](https://github.com/AlexsJones/llmfit/blob/main/docs/custom-models.md) 自行加入.
- `serve --host 0.0.0.0` 會把硬體資訊與 API 暴露到 LAN, 文件未提到任何驗證機制 (推論); 對外請放在 reverse proxy 後或只綁 loopback.
- 只有在使用下載, provider 查詢或 leaderboard 等功能時才會連外; 本身無 telemetry.
- Cluster mode 仍是 preview, 只對 mock Ray Dashboard 測試過.
