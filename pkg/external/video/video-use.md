# video-use

| 欄位 | 值 |
| --- | --- |
| Source | https://github.com/browser-use/video-use |
| Type | `skill` |
| Author / License | Browser Use / `MIT` |
| Stats | 31503 stars, 3653 forks, last commit 2026-10-08 (`main`), no GitHub releases (2026-11-24 via GitHub API) |
| Noted | 2026-11-24 |
| Brief | 以對話驅動影片剪輯的 Agent Skill, 透過詞級音訊逐字稿與按需視覺條讓 Coding Agent 直接規劃 EDL 並以 FFmpeg 渲染成片 |

## 摘要 (Summary)

`video-use` 是一套專為 Coding Agent (如 Claude Code, Codex) 設計的影片自動剪輯技能, 不靠抽大量影片幀進視覺模型, 而是透過 "逐字逐字音訊轉錄 (`takes_packed.md`) + 切點按需生成的視覺/波形複合圖 (`timeline_view`)" 讓 LLM 精準下刀並產生 EDL, 再由本地 Python helpers 驅動 FFmpeg 完成調色, 30ms 音訊淡入淡出, 動畫疊加, 字幕壓製與 `-14 LUFS` 響度標準化. 最大優勢是大幅壓低長片分析的 Token 成本, 並內建 12 條防靜默錯誤的硬規則與 3 輪自我驗證機制; 最大限制是強制依賴雲端 ElevenLabs Scribe API 進行語音轉錄 (無本地 ASR 選項), 且極度以人聲對白為核心, 對無對白純畫面素材較吃力. 適合口播, 教學, 訪談與產品宣傳片快速自動化成片; 不適合完全離線環境, 零 API 預算或大量純 B-roll 視覺藝術剪輯. 建議 `採用` (若需同時保留互動式 GUI 時間軸微調, 可搭配同類的 [`openshot-qt`](openshot-qt.md) 對照選型; 本專案主打零 GUI 對話式剪輯, 與 `openshot-qt` 的全功能時間軸和 JSON-RPC CLI 互補).

## 它是什麼 (What is it)

`video-use` 解決通用 LLM 無法直接 "看懂" 數萬幀原始影片並進行幀級精準剪輯的痛點. 它仿效 `browser-use` 提供結構化 DOM 給網頁代理的做法, 為影片剪輯建立兩層資訊介面與一條固定順序的 FFmpeg 管線:

| 階段 / 元件 | 對應腳本與產物 | 核心機制 |
| --- | --- | --- |
| 1. 轉錄與壓縮 (Layer 1) | `helpers/transcribe.py`, `helpers/transcribe_batch.py`, `helpers/pack_transcripts.py` -> `takes_packed.md` | 抽取 16kHz 單聲道 WAV 送往 ElevenLabs Scribe (`scribe_v1`) 取得詞級時間戳, 講者辨識與音效事件, 再依靜音 `>= 0.5s` 或換人折疊成純文字稿 |
| 2. 視覺鑽取 (Layer 2) | `helpers/timeline_view.py` -> PNG 複合圖 | 僅在停頓模糊或驗證切點時生成 "橫向縮圖條 + RMS 波形 + 單字標註 + 靜音陰影" 單張圖片 |
| 3. 策略與並行動畫 | `SKILL.md`, `skills/manim-video/` -> `edl.json`, `animations/slot_<id>/` | 先與使用者確認剪輯策略, 再以並行 sub-agent 透過 HyperFrames, Remotion, Manim 或 PIL 製作動畫疊層 |
| 4. 渲染與自我驗證 | `helpers/grade.py`, `helpers/render.py` -> `edit/final.mp4` | 分段調色與 30ms `afade` -> `-c copy` 無損拼接 -> `setpts` 疊加動畫 -> 最後燒錄 `master.srt` -> 兩段式 `loudnorm`, 再於切點執行 `timeline_view` 檢查與 critic 審查 |

使用方式: 將整個 repo symlink 至 `~/.claude/skills/video-use` (或對應 Agent 的 skills 目錄), 安裝 Python 依賴與 `ffmpeg`, 在 `.env` 填入 `ELEVENLABS_API_KEY` 後, 直接於原始影片資料夾開啟 Coding Agent 下達剪輯指令即可, 所有產物會寫入該資料夾下的 `edit/`.

## 視角評估 (Perspectives)

| 視角 | 優勢 (Advantage) | 劣勢 (Disadvantage) | 結論 |
| --- | --- | --- | --- |
| Token 成本與剪輯架構 | 以 `takes_packed.md` 取代全片抽幀, 切點才看圖; 分段拼接避免重複編碼並防呆字幕順序與音訊爆音 | 高度依賴語音對白, 純畫面 B-roll 難以只靠文本剪輯; 固定採用 CPU `libx264` 渲染無硬體編碼加速 | 口播/教學類剪輯的 Token 效率極高, 工程防呆規則極為成熟 |
| 外部依賴與隱私費用 | 影片畫面留在本機處理, 僅上傳 16kHz 音軌; 靜音軌 (`< -60 dBFS`) 先攔截且轉錄結果本地快取 | 強制綁定付費 ElevenLabs Scribe API (音訊會外送雲端且無本地 Whisper 選項); 動畫引擎需額外依賴 | 影像資料安全但音訊需走雲端付費 API, 不適合完全離線需求 |
| 安全邊界與系統權限 | 強制將所有輸出隔離於 `<videos_dir>/edit/`, 動刀前必經策略確認; `ffmpeg` 呼叫無 `shell=True` | 依賴廣泛 `Bash` 與平行 `Agent` 權限, 並會在執行期觸發 `npx --yes` 等遠端套件下載與安裝 | 專案目錄隔離謹慎, 但執行層授權極大, 適合受控或本機個人工作區 |
| 跨 Agent 可攜性與維護 | 標準 `SKILL.md` + 獨立 Python helpers, 支援 symlink 跨 Agent 部署且隨附完整 `manim-video` 子技能 | 平行動畫工作流點名依賴 Claude Code `Agent` tool; 無 GitHub Release/Tag 且測試覆蓋率低 | 跨 Agent 遷移容易, 但最佳多代理體驗仍以 Claude Code 為準 |

### Token 成本與剪輯架構

- 證據: `README.md:71-93`, `helpers/pack_transcripts.py:1-16, 38-163`, `helpers/timeline_view.py:1-63`, `helpers/render.py:42-57, 263-314, 527-706`, `SKILL.md:18-33, 91-104`.
- 優勢: `pack_transcripts.py` 將原始 Scribe JSON 折疊為極精簡的階段對白文字 (README 估算為原始 JSON 約 1/10 token, 相較全片抽幀可避免數千萬 token 雜訊, 未實測); `SKILL.md:18-33` 寫死 12 條硬規則解決 LLM 剪輯常見靜默失敗: 詞邊界吸附與 30-200ms 緩衝, 30ms `afade` 消爆音, 分段套用 `grade.py` 後再 `-c copy` 無損拼接, `setpts=PTS-STARTPTS+T/TB` 對齊動畫疊層, 字幕永遠置於濾鏡鏈最後一層且 `MarginV=90` 避開短影音底部介面遮擋 (`render.py:44-57`), 以及預設 `-14 LUFS` 兩段式響度正規化; 成片前強制在切點重跑 `timeline_view.py` 自我檢查 (上限 3 輪).
- 劣勢: 核心假設是 "Audio is primary", 面對缺乏旁白的風景 Montage, 純操作螢幕錄影或默片剪輯時, 文字稿無法反映畫面動作變化; `render.py:308, 698` 硬編碼 `-c:v libx264` 軟體編碼 (CRF 20/22/28), 未利用 macOS VideoToolbox 或 GPU 編碼器, 長片渲染時間會較長 (推論, 未實測).
- 結論: 將影片剪輯拆解為高訊噪比文本 + 按需視覺採樣, 是目前極少見且具工程驗證深度的 Agent 剪輯架構.

### 外部依賴與隱私費用

- 證據: `pyproject.toml:6-17`, `install.md:14-69, 92-100`, `helpers/transcribe.py:33-49, 84-114, 158-174`, `SKILL.md:62-68, 331`, `.env.example:1`.
- 優勢: 核心 Python 套件極少 (`requests`, `librosa`, `matplotlib`, `pillow`, `numpy`, 注意 `timeline_view.py:89` 註解其實連 `librosa` 都改以 `wave` + `numpy` 手刻避免硬依賴); 原始影音檔案不離開本機, 僅以 `ffmpeg` 萃取 16kHz 單聲道 PCM WAV 上傳轉錄, 且上傳前透過 `peak_dbfs < -60.0` 自動阻擋誤選靜音音軌 (如 OBS 雙軌錄音選錯軌, `transcribe.py:165-173`); 轉錄結果快取於 `edit/transcripts/` 防止重複付費.
- 劣勢: 轉錄完全依賴外部付費 API (`https://api.elevenlabs.io/v1/speech-to-text`, `model_id: "scribe_v1"`), `SKILL.md:331` 更明確反對在本地跑 Whisper (認為慢且會正規化掉口癖詞), 這意味著所有對話音訊必須送往 ElevenLabs 伺服器且無離線替代方案; 若使用動畫生成, 還須另備 Node.js 22+ (`HyperFrames`), `Remotion`, 或 `manim` 與數 GB 等級的 LaTeX (`texlive-full` / `mactex`, `skills/manim-video/SKILL.md:27`).
- 結論: 兼顧畫面隱私與高準確度詞級轉錄, 但付費外連 ElevenLabs 是無法繞過的硬門檻.

### 安全邊界與系統權限

- 證據: `SKILL.md:31-33, 62, 67, 81, 214-216`, `install.md:24-30, 68`, `helpers/render.py:62-66, 393`.
- 優勢: 嚴格區隔工具目錄與資料目錄 (Hard Rule 12 禁止在 `video-use/` 內寫入任何編輯產物, 連 `.env` 也不得寫進 `<videos_dir>`); 強制確認關卡 (Hard Rule 11 要求先核對自然語言策略再動剪輯); `install.md:68` 明確限制不得嘗試猜測 `sudo` 密碼; `helpers/*.py` 呼叫外部指令皆採 list 形式傳遞給 `subprocess.run`, 未開啟 shell 注入風險.
- 劣勢: 技能正常運作仰賴極高的工具權限: 允許 Agent 執行任意 `Bash` 指令, 在首次使用時自動執行 `npx --yes hyperframes` 或 `npx create-video@latest` 下載並執行 npm 套件 (`SKILL.md:67, 214-216`), 以及讀寫外部目錄檔案與呼叫 LLM/Scribe 外部 API; 同時允許將自訂 raw `ffmpeg` filter 字串直接寫入 `edl.json` 並執行 (`render.py:68-85`).
- 結論: 輸出路徑與使用者確認流程設計良好, 惟自動化安裝與任意指令執行能力仍要求在可信任的本機工作環境下使用.

### 跨 Agent 可攜性與維護

- 證據: `SKILL.md:1-4, 31, 70, 81, 257`, `skills/manim-video/SKILL.md:1-5`, `install.md:1-4, 70-90`, `tests/` (3 支測試), GitHub API (`pushed_at: 2026-10-08`, `open_issues_count: 163`, 2026-11-24).
- 優勢: 頂層 `SKILL.md` 採用標準 YAML frontmatter (`name`, `description`), 所有 helper scripts 與 `SKILL.md` 透過相對路徑綁定, 只要 symlink 整個 repo 資料夾即可同時在 Claude Code, Codex, Hermes, Openclaw 等系統運作; 隨附極完整的 `skills/manim-video/` 子技能 (含 15 份動畫與視覺設計規範文件); 社群熱度高 (31.5k stars, 3.6k forks).
- 劣勢: `SKILL.md:31, 81, 257` 明確引用 Claude Code 專有的 `Agent` tool 與 `Bash` 進行平行動畫子代理調度, 切換至其他不支援平行 sub-agent tool 的模型/介面時會退化為循序執行或須改寫提示; GitHub 倉庫無任何 Release 或 Tag (`404` via GitHub API), 單元測試僅有 `test_render_captions.py`, `test_render_fps.py`, `test_render_orientation.py` 三支, `pushed_at` 停在 `2026-10-08` 且累積 163 筆 open issues.
- 結論: 核心工具鏈與文件結構跨平台可攜性佳, 但完整多代理自動化體驗最貼合 Claude Code 生態, 且需留意上游近期更新節奏與未發佈正式版的現況.

## 細節 (Details)

- **本機現況**: 本機 (`~/.claude/skills/`, `~/.codex/skills/`, `~/Developer/video-use`) 尚未安裝 `video-use`; 本機已具備 `ffmpeg` 與 `ffprobe` (`/Users/shuk/.local/homebrew/bin/`).
- **安裝指令**:
  ```bash
  git clone https://github.com/browser-use/video-use ~/Developer/video-use
  ln -sfn ~/Developer/video-use ~/.claude/skills/video-use   # Claude Code
  cd ~/Developer/video-use && uv sync                        # 或 pip install -e .
  cp .env.example .env                                       # 填入 ELEVENLABS_API_KEY
  ```
- **核心 Helper 指令一覽**:
  - `python helpers/transcribe_batch.py <videos_dir> --workers 4`: 批次轉錄整個資料夾 (`--audio-track 1` 可指定 OBS 麥克風軌).
  - `python helpers/pack_transcripts.py --edit-dir <videos_dir>/edit`: 產生 `takes_packed.md`.
  - `python helpers/timeline_view.py <video> <start> <end>`: 產生含縮圖與波形的決策/驗收圖.
  - `python helpers/render.py <edl.json> -o <out.mp4> --build-subtitles [--preview|--draft]`: 渲染成片 (`--preview` 為 1080p CRF 22; `--draft` 為 720p CRF 28).
- **注意事項**:
  - `SKILL.md:188` 說明字幕 `MarginV=35`, 但實際 `helpers/render.py:56` 已調高為 `MarginV=90` 以避開直式短影音平台底部 UI (文件與實作數字略有差異, 以 `render.py` 為準).
  - 專案另有雲端代管版 `video-use.com` Waitlist 與 `Browser Use Box` 遠端執行宣傳 (`README.md:11, 42`), 但本開源 repo 本身純為本機 CLI/Skill.

## 相關參考 (References)

- 同類影片編輯筆記: [`openshot-qt.md`](openshot-qt.md) (桌面 GUI + `--cli` JSON-RPC 剪輯器, 與本 skill 純對話/FFmpeg EDL 管線定位互補)
- 相關 Agent Skill 筆記: [`agent-desktop-macos.md`](../skill/agent-desktop-macos.md) · [`web-design-guidelines.md`](../skill/web-design-guidelines.md)
- 專案首頁與動畫整合: [video-use.com](https://video-use.com) · [HyperFrames](https://github.com/heygen-com/hyperframes) · [Remotion](https://www.remotion.dev/) · [Manim](https://www.manim.community/)
