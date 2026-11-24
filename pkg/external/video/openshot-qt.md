# openshot-qt

| 欄位 | 值 |
| --- | --- |
| Source | https://github.com/OpenShot/openshot-qt |
| Type | `project` |
| Author / License | OpenShot Studios, LLC (Jonathan Thomas 等) / `GPL-3.0-or-later` |
| Stats | 6832 stars, 857 forks, last commit 2026-11-24 (`develop`), latest release `v4.1.1` (2026-11-02) (2026-11-24 via GitHub API) |
| Noted | 2026-11-24 |
| Brief | 跨平台開源非線性影片剪輯器 (NLE), 以 Python/Qt 介面結合 C++ `libopenshot` 引擎, 並內建 Agent-first CLI 與 ComfyUI 工作流整合 |

## 摘要 (Summary)

`openshot-qt` 是 OpenShot Video Editor 的桌面端與自動化主專案, 透過 Python Qt 前端驅動底層 C++ `libopenshot` 與 FFmpeg 剪輯引擎. 最大優勢在於完全開源免費 (`GPLv3`), 並在 `develop` 分支 (`4.2.0-dev1`) 新增了以 JSON-RPC 運作的 Agent-first CLI (`--cli`) 與 25 份本機 ComfyUI AI 工作流模板, 讓外部 AI Agent 能直接建立與剪輯 `.osp` 專案; 最大限制是最新穩定釋出 `v4.1.1` 的 macOS 安裝包仍僅提供 `x86_64` (Apple Silicon 尚未有獨立預編譯包), 且開箱預設開啟 GA4 匿名指標與版號檢查連線. 適合需要 "本機離線剪輯 + 可由 AI/腳本驅動時間軸" 或尋找無浮水印開源 NLE 的情境; 不適合追求原生 M 系列極致即時 4K 調色效能的專業後期剪輯. 建議 `採用` (若需 Agent-first CLI 則建議 `觀望` 待 `4.2.0` 正式版或自行從 `develop` 建置; 相較於純靠對白文字稿與 FFmpeg 批次組裝的 [`video-use`](video-use.md), 本專案更側重完整 NLE 時間軸, GUI 與可程式化多軌合成).

## 它是什麼 (What is it)

`openshot-qt` 解決一般桌面剪輯器昂貴, 閉源或難以程式化控制的問題, 提供跨平台圖形介面與可腳本化的剪輯工作階段. 它本身負責 UI, 專案管理 (`.osp` UTF-8 JSON), 外掛工作流與自動化協定, 實際的影音解碼, 合成, 轉場與匯出交由 C++ 函式庫 `libopenshot` 與 `libopenshot-audio` 處理.

| 層級 | 元件與路徑 | 職責 |
| --- | --- | --- |
| 桌面啟動與 Qt 抽象層 | `src/launch.py`, `src/qt_api.py` | 統一相容 `PyQt5` / `PyQt6` / `PySide6`, 初始化使用者目錄 (`~/.openshot_qt`) 與主視窗 |
| 專案與媒體服務 | `src/classes/project_data.py`, `src/classes/proxy_service.py` | 讀寫 UTF-8 JSON 格式 `.osp` 專案, 管理代理轉碼 (`optimized/`) 與波形快取 |
| Agent-first CLI 與協定 | `src/cli.py`, `src/classes/agent/host.py`, `src/classes/agent/service.py` | 透過 `QLocalServer` (Unix socket) 與 JSON-RPC 2.0 提供無頭 (`headless`) 或連入現有 GUI 的即時剪輯控制 |
| 本機 AI 與 3D 擴充 | `src/classes/comfy_client.py`, `src/comfyui/*.json`, `src/blender/` | 對接本機或自架 ComfyUI WebSocket/HTTP API (25 份模板), 本機 ONNX (YOLO26 / EfficientSAM / Cutie) 與 Blender 5.0+ 3D 動畫 |
| 影音核心引擎 (外部相依) | `OpenShot/libopenshot`, `OpenShot/libopenshot-audio`, FFmpeg 7.1 | 實際時間軸幀合成, OpenCV 追蹤, Vulkan/硬體編解碼與最終渲染匯出 |

使用方式: 終端使用者可安裝預編譯版後直接開啟 GUI; 開發者或 Agent 可透過 `openshot-qt --cli session start` / `openshot-cli` 啟動無頭工作階段, 再依序呼叫 `file import`, `clip add`, `edit apply` 與 `export start`.

## 視角評估 (Perspectives)

| 視角 | 優勢 (Advantage) | 劣勢 (Disadvantage) | 結論 |
| --- | --- | --- | --- |
| 平台支援與相依架構 | 支援 Linux/macOS/Windows, 跨 Qt5/Qt6 三種 binding, 提供完整預編譯包 | 從原始碼建置需編譯 C++ `libopenshot` 與多項系統依賴; `v4.1.1` macOS 釋出包僅 `x86_64` | 優先使用官方預編譯包; Apple Silicon 暫依賴 Rosetta 2 或自行編譯 |
| 隱私與雲端邊界 | 剪輯運算, `.osp`, ComfyUI (預設 `127.0.0.1:8188`) 與 ONNX 模型全在本機運行 | 預設 `send_metrics: true` 會傳送 GA4 事件, 啟動時查詢 `openshot.org` 版號並可觸發 Sentry | 素材不離機, 注重絕對離線應於設定關閉 `send_metrics` |
| 整合與自動化能力 | 內建 Agent-first CLI (JSON-RPC + 共享 Undo + 256 指令批次 + 聯絡表預覽), `.osp` 為純 JSON | 無原生 MCP Server; `--cli` 完整功能位於 `develop` (`4.2.0-dev1`), 尚未進入穩定版 `v4.1.1` | 具備極佳程式化與 AI Agent 驅動潛力, 優於多數桌面 NLE |
| 效能, 資源與授權 | `GPL-3.0-or-later` 無付費牆; 支援 Vulkan / VideoToolbox / NVENC / VA-API 等硬體匯出與代理媒體 | `hw-decoder` 預設關閉且屬實驗性質; 預設記憶體快取上限 `768 MB`; AI/3D 依賴外部工具 | 免費功能完整度高, 但重度剪輯效能與即時性仍遜於商用原生引擎 |

### 平台支援與相依架構

- 證據: `README.md:96-106`, `.github/workflows/ci.yml:50-121`, `.gitlab-ci.yml:19-255`, `installer/build-mac-dmg.sh:5-32`, GitHub Release `v4.1.1` assets (`OpenShot-v4.1.1-x86_64.dmg`: 308,547,400 bytes / ~294 MiB; `OpenShot-v4.1.1-x86_64.AppImage`: 258,921,664 bytes / ~247 MiB; `OpenShot-v4.1.1-x86_64.exe`: 248,497,059 bytes / ~237 MiB; `OpenShot-v4.1.1-x86.exe`: 243,824,048 bytes / ~233 MiB, 2026-11-24 via GitHub API).
- 優勢: 官方對 Linux (`AppImage`), macOS (`DMG`, 含 codesign 與 notarization), Windows (`x64`/`x86` `.exe` 與 `MSIX` 打包腳本) 皆有完整 CI 建置管線; Python 層透過 `src/qt_api.py` 同時支援 `PyQt5`, `PyQt6`, `PySide6`, 並在 Ubuntu, Windows, macOS 上以矩陣 CI 驗證相容性 (`ci.yml:92-121`).
- 劣勢: Python 程式碼無法獨立執行, 必須搭配對應版本的 C++ `libopenshot` (`MINIMUM_LIBOPENSHOT_VERSION = "1.2.0"`, `src/classes/info.py:34`), `libopenshot-audio`, FFmpeg 7.1, OpenCV 4.13 與 Protobuf (`ci.yml:53-70`); 雖然 `installer/build-mac-dmg.sh:29-31` 與 `.gitlab-ci.yml:98` 已出現 `arm64-support` 相關分支檢查, 但 `v4.1.1` 公開釋出資產中僅有 `x86_64` 版 `.dmg`, Apple Silicon Mac 原生支援尚未落地於正式發行版.
- 結論: 三大桌面平台安裝門檻低, 惟 Mac ARM 使用者目前需接受 Rosetta 2 轉譯或等待後續釋出.

### 隱私與雲端邊界

- 證據: `src/settings/_default.settings:610-639`, `src/classes/metrics.py:81-225`, `src/classes/sentry.py:51-116`, `src/classes/version.py:34-70`, `src/classes/app.py:172-174`, `src/classes/comfy_client.py:61-108`, `src/tests/test_video_widget_transform.py:854-859`.
- 優勢: 核心剪輯, 轉檔與專案檔儲存完全在本機運作; 檢查 `src/classes/` 與 `src/windows/` 未發現將使用者影音素材上傳至 `cloud` 或 OpenShot 伺服器的連線邏輯; AI 生成整合採用使用者自建的 ComfyUI (預設 URL 為 `http://127.0.0.1:8188`, `_default.settings:610-615`), 物件偵測/遮罩模型 (YOLO26, EfficientSAM, Cutie) 亦下載至本機 `~/.openshot_qt/yolo/` 離線推論; 若將 `send_metrics` 設為 `false`, `sentry.disable_tracing()` 與 `metrics.py:179` 會停止傳送錯誤與統計.
- 劣勢: `_default.settings:634-639` 預設 `"send_metrics": true`, 會透過 `GA4_ENDPOINT` (`https://www.google-analytics.com/mp/collect`) 傳送畫面瀏覽, UI 事件與系統/Qt/Python 版本參數 (`metrics.py:101-115`) 並附帶 `unique_install_id` (`metrics.py:213`); 程式亦會向 `https://www.openshot.org/version/json/` 取得版本與動態錯誤採樣率 (`version.py:42-66`), 並於打包 `sentry_sdk` 時初始化硬編碼之 Sentry DSN (`sentry.py:109`).
- 結論: 媒體資料主權完全保留於本機, 但開箱預設有遙測與版本檢查外連行為, 需手動關閉 `send_metrics` 才算徹底靜默.

### 整合與自動化能力

- 證據: `src/launch.py:56-61`, `src/cli.py:1-188`, `src/classes/agent/protocol.py:33-155`, `src/classes/agent/host.py:38-174`, `src/classes/agent/service.py:252-550`, `src/classes/project_data.py:76-106`, `src/comfyui/*.json`.
- 優勢: `develop` 分支導入設計相當完整的 `Agent-first CLI`, 無論單獨啟動無頭 session (`--cli session start --project ...`, Linux 無顯示器自動切換 `QT_QPA_PLATFORM=offscreen`) 或連線到已開啟的 GUI 視窗 (`session list` -> `session show`), 皆透過 `QLocalServer` + JSON-RPC 2.0 + `secrets.token_hex(32)` 驗證進行互動 (`host.py:44-94`), 且 CLI 與 GUI 共享 Undo/Redo 歷史 (`cli.py:56-57`); 提供 `timeline overview`, `timeline view --mode waveforms`, `file strip` (產生時間軸聯絡圖給多模態模型檢視), `edit bulk --dry-run` 與支援 `$ref` 鏈結的批次操作 (`edit apply`, 單批上限 256 展開指令, `cli.py:180-186`); `.osp` 為純 UTF-8 JSON 可直接由程式讀寫; 另附 25 份 ComfyUI JSON 模板 (含 SAM2 遮罩追蹤, Whisper SRT 字幕, RIFE 2x 補幀, Wan/SVD 影生影等).
- 劣勢: 專案內沒有原生 MCP Server 實作, 須由 AI Agent 透過 shell 呼叫 `--cli` 指令集或自行包一層 MCP wrapper; 此 Agent CLI 架構主要於 `4.2.0-dev1` 開發中 (`cli.py` 與 `agent/*.py` 版權年份皆標注 `2026`), 在最新穩定版 `v4.1.1` 中的可用程度未經安裝實測 (推論穩定版尚未包含完整 `--cli` 命令樹, 需使用 Daily Build 或 `develop` 分支).
- 結論: 在開源桌面影片剪輯器中少見地直接為 AI Coding/Media Agent 設計一等公民介面, 非常適合作為自動化剪輯管線後端.

### 效能, 資源與授權

- 證據: `COPYING:1-20`, `README.md:34, 148-159`, `src/settings/_default.settings:495-608, 687-820`, `src/classes/export_settings.py:34-101`, `src/classes/proxy_service.py:55-90`, `src/classes/info.py:176`.
- 優勢: 採用 `GPL-3.0-or-later` 授權, 無功能付費牆, 解析度限制或匯出浮水印; 預覽快取支援 `CacheMemory` (預設 `768 MB`, 最多 `900` frames, 預讀 `24-48` frames, 數值取自 `_default.settings:498-556`, 未實測) 與 `CacheDisk` 縮放快取; 內建 `ProxyService` 產生輕量化優化代理檔 (`optimized/`); GPU 管線支援 `gpu-mode: auto` (Vulkan) 並於 `preferred_export_codec` 依序優先選用 `h264_vulkan`, `h264_videotoolbox`, `h264_nvenc`, `h264_qsv`, `h264_amf`, `h264_vaapi` (`export_settings.py:96-99`), 失敗時退回 `libx264`.
- 劣勢: 預設硬體解碼 `hw-decoder` 為 `"0"` (`No acceleration`, `_default.settings:710-717`) 且上限預設 `1950x1100` (`_default.settings:796-807`, 未實測), `README.md:34` 亦註明硬體編解碼仍屬實驗性 (Experimental); 3D 動畫字幕需要額外安裝 Blender `5.0+` (`info.py:176`), 25 種 ComfyUI AI 特效亦須自備外部 ComfyUI 實例與 GPU 算力; 相比 C++ 原生專業 NLE (如 DaVinci Resolve / Kdenlive), Python/Qt 架構在高軌數與重度特效預覽時的記憶體與幀率表現可能受限 (推論, 本機未安裝實測).
- 結論: 授權自由且功能廣度佳, 日常剪輯與自動化輸出夠用, 但進階硬體加速與重度剪輯體感需自行調校與實測.

## 細節 (Details)

- **本機現況**: 本機 (`/Applications/` 與 `$PATH`) 目前未安裝 `openshot-qt`.
- **預編譯安裝**: 從 [OpenShot Download](https://www.openshot.org/download/) 或 [GitHub Releases (`v4.1.1`)](https://github.com/OpenShot/openshot-qt/releases/tag/v4.1.1) 下載 `.dmg` (macOS x86_64, ~294 MiB), `.AppImage` (Linux x86_64, ~247 MiB) 或 `.exe` (Windows x64/x86, ~233-237 MiB). 若要體驗 `4.2.0-dev1` 的 `--cli` 與 ComfyUI 工作流, 需從官方網站下載 **Daily Builds**.
- **從原始碼啟動** (需先行建置並透過 `PYTHONPATH` 掛載 `libopenshot` Python bindings):
  ```sh
  PYTHONPATH=/path/to/libopenshot/build/bindings/python python3 src/launch.py
  # 或僅啟動 Agent-first CLI (查詢 help 不需要 Qt 與 libopenshot)
  python3 src/launch.py --cli help
  ```
- **關鍵設定與目錄** (`~/.openshot_qt/`):
  - `openshot.settings`: 建議將 `send_metrics` 設為 `false` 以停用 GA4 與 Sentry.
  - `sessions/`: Agent CLI 探索檔案 (`<session_id>.json`, 權限 `0o600`), 記錄 Unix socket endpoint 與臨時 token.
  - `optimized/` / `yolo/` / `comfyui-output/`: 分別存放代理剪輯檔, 下載的 ONNX 分割模型與 ComfyUI 生成輸出.

## 相關參考 (References)

- 同類影片編輯筆記: [`video-use.md`](video-use.md) (純 Coding Agent 對話與 FFmpeg/Scribe 剪輯管線, 與本專案 GUI + JSON-RPC CLI 定位互補)
- 其他桌面應用筆記: [`openwhispr.md`](../app/openwhispr.md) (語音轉文字與會議筆記桌面應用)
- 核心影音引擎: [OpenShot/libopenshot](https://github.com/OpenShot/libopenshot) 與 [OpenShot/libopenshot-audio](https://github.com/OpenShot/libopenshot-audio)
- 官方文件與開發指南: [OpenShot User Guide](https://www.openshot.org/user-guide/) · [Become a Developer](https://github.com/OpenShot/openshot-qt/wiki/Become-a-Developer)
