# Threejs-WebGPU-IFFT-Ocean (Spiri0/Threejs-WebGPU-IFFT-Ocean)

- Source: <https://github.com/Spiri0/Threejs-WebGPU-IFFT-Ocean>
- Demo: <https://spiri0.github.io/Threejs-WebGPU-IFFT-Ocean/index.html> (2026-11-20 以 `curl -I` 確認回 200)
- License: MIT (有 LICENSE 檔)
- Stats (2026-11-20 via GitHub API): 143 stars, 8 forks, 1 open issue (#7 "How would you add buoyancy to this?", 2024-11-14 開, 無回應), created 2024-02-08, last push 2025-12-12, commit 數 >= 100
- Size: 應用程式約 4,000 行 JS / WGSL (`src/` + `resources/shader/`), 最大的是 `quadtree.js` 293 行, `ocean.js` 284 行, ocean fragment shader 262 行; 另把 `node_modules/three` (0.179.0, 32 MB) 直接 commit 進 repo; `resources/` 8 MB
- Stack: three.js `WebGPURenderer` + TSL `wgslFn` (把手寫 WGSL 包成 node), 無 build step (import map 直接指向 `node_modules/three/build/three.webgpu.js`)

## 它做什麼 (What it does)

在瀏覽器以 WebGPU 跑 `頻譜 + IFFT` (Inverse Fast Fourier Transform, 先在頻率空間描述每個波長的振幅, 每幀反轉換回空間得到整張高度圖) 的無限大海面, 並用 quadtree LOD 鋪出 1,000 km 見方的海.

- `波譜模型`: JONSWAP + TMA 水深修正 + 方向分佈 (cosine-2s, spread blend, swell), 兩組波譜疊加 (預設風向 0° 與 240°, fetch 100 km 與 300 km); 20 個參數在 lil-gui 即時調整, 改動時重算 initial spectrum
- `三層 cascade`: length scale `250 m / 17 m / 5 m`, 每層 512x512; 以 wavenumber 上下界切分頻帶 (`boundaryHigh = 2π/L(i+1) * 6`), 避免同一波長在兩層重複出現
- `泡沫 (foam)`: 由 displacement 的 Jacobian 行列式 (表面被擠壓的程度) 判斷浪頭, 並以 `turbulence` buffer 跨幀累積與衰減
- `水面幾何`: quadtree 依相機距離切出 chunk, 每個 chunk 是 36x36 格的平面網格, 由 7 個 Web Worker 產生; vertex shader 做 CDLOD 式 morphing (相鄰 LOD 交界處把奇數頂點滑向偶數頂點) 消除接縫與跳變
- `著色`: 法向量由 derivative texture 在 fragment 逐 pixel 算 (不依賴網格密度), Fresnel 混合 cubemap 反射與固定深海色, Blinn 式高光, 500 m 到 4,000 m 漸入霧色
- `天空`: Preetham 式大氣散射 sky box, 每次更新以 `CubeCamera` 渲染成 256 cubemap 給水面反射
- `操控`: 鍵盤飛行 (`W/A/S/D` 速度與翻滾, 方向鍵俯仰與轉向, `Space` 加速), 相機起點高度 5 m

`推論`: 波譜函式命名 (`JonswapAlpha`, `spreadPower`, `cosine2s`, `shortWavesFade`), 預設 length scale 250/17/5 與 turbulence 公式都與 Unity 專案 gasgiant/FFT-Ocean 相同, 應是其 WebGPU 移植.

### 每幀管線 (Per-frame Pipeline)

1. 每個 cascade: `timeSpectrum` (把 h0 依色散關係推進到目前時間, 拆成 4 組複數場) -> 4 次 IFFT -> `texturesMerger` (寫 displacement / derivative / jacobian 三張 storage texture)
2. 每次 IFFT: `init` 1 + horizontal butterfly 8 + vertical butterfly 9 + `permute` 1 = 19 個 dispatch
3. `CubeCamera.update` 重繪 sky 6 面
4. Quadtree 依相機位置重算 leaf, 新 chunk 丟給 worker pool, 舊 chunk 等 pool 閒置後移除
5. 主畫面 render: 每個 chunk 一個 `Mesh`

## 效能與規模 (Performance and Scale)

數字由原始碼常數推導 (`IFFT_RESOLUTION: 512`, `LENGTH_SCALES: [250, 17, 5]`, `QT_OCEAN_MIN_CELL_RESOLUTION: 36`, `QT_OCEAN_MIN_LOD_RADIUS: 15`, `QT_OCEAN_MIN_NUM_LAYERS: 15`, `OCEAN_SIZE: 500000`); chunk 數由 scratchpad 內重現 `quadtree._Insert` 邏輯的腳本算出; `未實際量測 FPS`.

| 元素 (Element) | 規模 | 每幀成本 |
| --- | --- | --- |
| IFFT cascades | 3 層 x 512² = 786,432 格; texel 尺寸 0.49 m / 3.3 cm / 0.98 cm | `234 個 compute dispatch` (每層 1 + 4x19 + 1 = 78), 每個 262,144 threads, 合計約 6,100 萬次 invocation |
| Butterfly passes | 每層 68 次, 共 204 次 | 每 thread 讀 2 個 vec4 + 1 個 butterfly vec4, 並對自己的 vec4 做 read-modify-write; 不計 cache 約 21 MB/pass, 4.3 GB/frame (`推論`: 上界估算, 實際受 cache 影響) |
| Compute GPU memory | 每層約 31.7 MiB (spectrum 4 + waveData 4 + 4 組複數場 8 + pingpong 4 + turbulence 1 + 三張 texture 含 mip 約 10.7), 3 層約 `95 MiB` | butterfly 表只有 72 KiB, 啟動時算一次 |
| Ocean chunks | 相機高 5 m 時 `748 個 leaf`; 50 m: 640; 500 m: 460; 5,000 m: 316 | 每個 chunk 一個 `Mesh` = 1 draw call, 上限 748 (視錐裁切後實際數量未量測) |
| Chunk geometry | 37² = 1,369 vertices, 2,592 triangles, 7,776 `Uint32` indices, 約 80 KB/chunk | 5 m 高度時總計 102 萬 vertices, 194 萬 triangles, 約 57 MiB geometry |
| LOD 範圍 | 最細 chunk 3.8 m 寬 (頂點間距約 10.6 cm), 最粗 250 km 寬 (間距約 6.9 km); quadtree 根 1,000 km 見方 | 相機移動跨越 LOD 邊界時才重建, 由 worker 產生後在主執行緒建 `BufferGeometry` |
| Vertex shader | 每頂點 3 cascade x 4 次 `textureLoad` (手寫 bilinear) = 12 次 + morph 迴圈最多 18 次 | 5 m 高度時約 1,230 萬次 texture load (裁切前) |
| Fragment shader | 每 pixel 6 次 `textureSample` (derivative + jacobian x 3) + 1 次 cubemap | `setPixelRatio(devicePixelRatio)` + MSAA, Retina 螢幕 fragment 數為 DPR 1 的 4 倍 |
| Sky cubemap | 6 面 x 256², half float, 約 4 MiB 含 mip | 每次更新重繪 6 面 (sky 只在 GUI 調整時才會變) |
| Worker pool | 寫死 7 個 module worker (作者註解: 平板 8 核扣主執行緒) | 結果以 `postMessage` 結構化複製回傳, 沒有用 transfer list |

- `可見細節的實際距離`: vertex 位移依 `min(LOD_SCALE * L / dist, 1)` 衰減 (`LOD_SCALE = 3.7`), 250 m 層在 925 m 內全幅, 17 m 層 63 m, 5 m 層 18.5 m; 霧在 500 m 開始, 4,000 m 全霧. 尺寸 >= 3,906 m 的 chunk 共 300 個 (佔 748 的 40%), 基本上只畫霧色 (`推論`: 由 chunk 尺寸與距離推得)
- `時間步進`: 主迴圈固定 60 Hz step, 每幀最多 1 次 (`accumulator` 上限 200 ms); 波的相位用 `performance.now()` 所以波速是實時的, 但 120 Hz 螢幕上每兩幀才更新一次 displacement (`推論`: 會出現 60 Hz 的位移階梯感, 未實測)
- `擴充上限`: compute 成本與 `IFFT_RESOLUTION²` 和 `log2(N)` 成正比, 1024 會讓每層 dispatch 從 78 變 86, 每個 dispatch threads 變 4 倍; `QT_OCEAN_MIN_CELL_RESOLUTION` 只能用偶數 (morph 公式要求), 從 36 改 72 時 geometry 與 vertex 成本約 4 倍, draw call 數不變
- `物件數`: 本身沒有任何浮力或物體互動. 要加 N 個浮體時, 水面成本不隨 N 增加, 新增成本只在取樣 displacement 的部分; 既有架構沒有 CPU readback, 相關設計見 [comparison.md](comparison.md)

## 最棒的地方 (What is great)

- `物理正確的波譜`: JONSWAP + TMA + 方向分佈 + 雙波譜疊加, 參數有物理單位 (風速, fetch, 水深), GUI 調整後即時重算; 是這批 repo 中波浪最寫實的參數化
- `三層 cascade 頻帶切分`: 大浪, 中浪, 細紋分層處理且頻帶不重疊, 近距離有 1 cm 級細節, 遠處有 250 m 長浪
- `Storage buffer IFFT`: README 記錄 2025-03-26 從 storage texture 改 storage buffer, 從 4x256 提升到 3x512 且效能明顯改善; 運算全在 GPU, 主執行緒只負責發 dispatch
- `Jacobian 泡沫 + 時間累積`: 浪頭泡沫跟著實際擠壓量出現並慢慢消散, 不是靠高度閾值
- `CDLOD quadtree + morphing`: 1,000 km 的海面只用數百個 chunk, 相機附近頂點間距 10 cm, LOD 交界以頂點 morph 平滑過渡, 不需要 skirt
- `融入 three.js scene graph`: 用 `WebGPURenderer` 與 TSL node, 可以直接在同一個 scene 放 GLTF 模型, 燈光與後製; 這是 [waves.md](waves.md) 與 [webgl-water.md](webgl-water.md) 都做不到的
- `零 build 即可部署`: import map + commit 進 repo 的 three.js, GitHub Pages workflow 直接發佈; MIT 授權且 2025-12 仍有更新 (跟進 r179 的 workgroup / dispatch size API)

## 最糟的地方 (What is worst)

- `沒有浮力或物體互動`: issue #7 從 2024-11 開著無回應; 作者 README 只指向另一個 voxelizer repo; 沒有任何 displacement readback 或取樣 API
- `IFFT dispatch 數量過多`: 每個 butterfly stage 一個 dispatch, 4 組複數場各自跑一次 IFFT, 共 234 個 dispatch/frame; ping/pong 交錯存在同一個 vec4, 每次寫入都要先讀回整個 vec4. `推論`: 改成 workgroup shared memory 一次做完一個方向, 並把 4 組場打包成 2 個 vec4 一起轉換, dispatch 數可降到每層 4 個左右
- `幾何管線浪費`: 每個 chunk 在本地座標其實是同一張 36x36 網格, 卻由 worker 每次重算世界座標並各自成為獨立 `Mesh`; CPU 算的 `normal` attribute 沒有被 shader 使用 (每 chunk 16 KB), 每頂點的 `width` / `lod` 是 chunk 常數; `推論`: 改成單一 `InstancedMesh` + per-instance offset/scale 可把 748 個 draw call 收成 1 個, worker 也可移除
- `GPU 記憶體洩漏`: 整個 `src/` 沒有任何 `dispose()`; 退役 chunk 只從 group 移除, chunk pool 從未被填回, 每個新 chunk 都重新配置 `BufferGeometry`. `推論`: 依 three.js 語意, 移出 scene 不會釋放 GPU buffer, 長時間飛行會持續累積
- `遠景成本無效`: 約 40% chunk 位於全霧區, 仍要跑完整 vertex shader 的 12 次 texture load; 最大 cascade 只有 250 m, 作者在 shader 註解承認遠處 mipmap 耗盡會閃爍, 所以用霧遮掉; 實際可見海面約 4 km
- `反射方向是 hack`: fragment 先算 `reflect()`, 隨即以 `halfVec` 覆蓋並對調 x/y 分量再取 cubemap, 不是物理反射方向; 材質標 `transparent = true` 但 alpha 恆為 1, 748 個 mesh 每幀走透明排序
- `每幀固定成本`: sky cubemap 6 面每次更新都重繪 (sky 靜態); `shadowMap` (PCFSoft) 開啟但場景沒有投影燈; `PostProcessing` / `pass` / `depthPass` 建了但從未使用
- `未使用的資產與死碼`: 載入 `uv_grid_opengl.jpg` (550 KB), `simplex-noise.png` (394 KB) 與 `cube/sky` 6 張 JPEG (420 KB) 但 shader 沒有讀取; `MilkyWay` (2 MB), `uv` (3.2 MB), `pisa` cubemap 與 `noise.exr` (1 MB) 完全沒被引用; fragment shader 留有 6 個未呼叫的函式
- `平台限制`: 只支援 WebGPU, 不支援時直接 throw, 沒有 WebGL fallback; 操控只有鍵盤 (`keyCode`, 已 deprecated), 沒有 touch 與滑鼠, 行動裝置無法操作; iOS 實機表現未知
- `工程衛生`: 32 MB `node_modules` 進版控; `vite.config.js` 引用 `vite` 與 `vite-plugin-top-level-await`, 但 `package.json` 只有 `three`; README 標題寫 r174, 實際是 0.179; quadtree 收到未定義的 `QT_OCEAN_MIN_CELL_SIZE`; 7 個 worker 寫死; 沒有測試; README 說下載後應改回 `SharedArrayBuffer`, 但程式碼現況是 `ArrayBuffer` 且無 transfer

## 可借鏡 (Takeaways)

- 要做 three.js 生態內的大範圍海洋: 直接沿用其 `波譜 + 3 cascade IFFT` 與 Jacobian foam, 這部分品質最高
- 移植時優先改三件事: (1) IFFT 改 shared-memory 單 dispatch 每方向, (2) chunk 改單一 `InstancedMesh` 並移除 worker, (3) sky cubemap 只在參數變動時重繪
- 浮力與大量物件: 在 displacement texture 上加一個 compute pass, 對 N 個物件取樣後寫入 instance matrix storage buffer, 避免 CPU readback; 具體分層設計見 [comparison.md](comparison.md)
