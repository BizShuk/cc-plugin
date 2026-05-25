# LLMFit

這是一個名為 `llmfit` 的命令列工具，旨在協助使用者根據其系統的硬體配置（包括 `RAM`、`CPU` 與 `GPU`）找到最適合運行的 `LLM` (大型語言模型)。

```bash
uvx llmfit
```

其主要功能包括：

1. `Hardware Detection` (硬體檢測)：自動檢測系統規格，並建議適合的 `Quantization` (量化) 版本。
2. `Model Scoring` (模型評分)：根據品質、速度、配適度與 `Context` (上下文) 維度為模型評分。
3. `Multiple Runtimes` (多種執行環境)：支援 `Ollama`、`llama.cpp`、`MLX`、`LM Studio` 等多種本地模型供應商。
4. `Interactive TUI` (互動式終端介面)：提供直觀的介面列出模型，並顯示預估的 `tok/s` (每秒標記數) 與記憶體使用量。
5. `Community Leaderboard` (社群排行榜)：查看其他使用相同硬體的使用者所回報的真實效能數據。
6. `Hardware Simulation` (硬體模擬)：允許使用者模擬不同的硬體環境來評估模型效能。

此工具適合想要在本地運行 `LLM` 但不確定硬體是否足以負擔的使用者，能有效避免下載到無法順暢運行的模型。

---

Source: <https://github.com/AlexsJones/llmfit>
