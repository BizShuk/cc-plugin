# waves (dli/waves)

- Source: <https://github.com/dli/waves>
- Demo: <http://david.li/waves> / video <http://www.youtube.com/watch?v=IrUehq6vJss>
- Author: David Li, License: MIT
- Stats (2026-11-20 via GitHub API): 1117 stars, 151 forks, last commit 2014-09-05, last push 2018-05-10
- Size: 5 個檔案, 約 1700 行純 JavaScript + GLSL, 無 build step, 無外部 library

## 它做什麼 (What it does)

瀏覽器內即時的 `FFT ocean`(以頻域頻譜加上逆 FFT 產生海面高度場的方法, Tessendorf 風格) 海浪模擬.
使用者可在 3D 場景中直接拖曳三個控制項:

- `WIND` 箭頭: 風向與風速 (m/s), 改變頻譜形狀
- `SIZE` 滑桿: 模擬區塊的實體尺寸 (m), 改變波長分佈
- `CHOPPINESS` 滑桿: 水平位移強度, 讓浪峰變尖 (旁邊有 2D canvas 畫出示意剖面)

### 管線 (Pipeline), 全部在 GPU 上

1. `initial spectrum` - fragment shader 依風速與尺寸算出 512x512 的初始頻譜振幅, 只在參數改變時重算
2. `phase` - 每幀依 dispersion relation 推進相位, ping-pong 兩張 texture; 相位與頻譜分開存, 拖動參數時浪不會跳
3. `spectrum` - 合成當前時刻的複數頻譜, 同時產生 X/Z 方向的 choppy displacement
4. `FFT` - Stockham formulation 的 GPU FFT, 水平 + 垂直共 `2 * log2(512) = 18` 個 pass, 輸出 displacement map
5. `normal map` - 從 displacement 差分出法向量
6. `ocean render` - 256x256 網格 (2000m 見方) 在 vertex shader 讀 displacement 位移; fragment shader 做 Schlick Fresnel 混合天空色與海水色, 再做 exponential tone mapping

UI 本身是 `CSS 3D transform` 的 DOM 元素, 每幀把 WebGL camera 的 view matrix 轉成 CSS matrix, 讓 HTML 控制項看起來「躺」在海面上; 滑鼠點擊以 unproject 射線打到 y=0 平面做 hit test.

## 效能與規模 (Performance and Scale)

數字由原始碼常數 (`RESOLUTION=512`, `GEOMETRY_RESOLUTION=256`, `GEOMETRY_SIZE=2000`) 算出; 未實際量測 FPS.

| 元素 (Element) | 規模 | 每幀成本 |
| --- | --- | --- |
| Initial spectrum | 512x512 texel | 1 個 fullscreen pass, `只在參數改變時`執行 |
| Phase | 512x512, ping-pong 2 張 | 1 pass |
| Spectrum (含 choppy X/Z) | 512x512 | 1 pass |
| GPU FFT | 512 點, 水平 9 + 垂直 9 stage | 18 pass |
| Normal map | 512x512 | 1 pass |
| Ocean mesh | 65,536 vertices, 130,050 triangles, 390,150 indices | 1 個 `drawElements` |
| CSS 3D UI | 約 20 個 DOM 元素 (箭頭 2 + 滑桿 2x3 + 靜態標籤) | 1 次 camera `transform` 更新 |

- `Draw calls`: 每幀 22 個 (參數改變時 23); 其中 21 個 offscreen pass 各跑 262,144 個 fragment, 合計約 `5.5M` 次 fragment shader 呼叫 / 幀, 成本集中在 FFT
- `GPU memory`: 8 張 `RGBA float` 512x512 texture, 每張 4 MiB, 合計約 `32 MiB`; vertex buffer 約 1.3 MB, index buffer 約 0.78 MB
- `空間解析度`: 網格間距 2000m / 255 ≈ `7.8m`; displacement map 是 512 但網格只有 256, 幾何只取到一半細節, 剩下靠 per-fragment normal map 補
- `尺度是放大而非鋪磚`: 模擬區塊 (`SIZE`, 預設 250m) 被整片拉伸到 2000m 網格上, displacement 乘上 `2000 / size` (預設 8 倍); 調 `SIZE` 等於調視覺尺度, 沒有重複鋪排
- `CPU 端`: 每幀只有 uniform 更新與 1 次 CSS matrix 轉換; 每次 mousemove 做 1 次 4x4 反矩陣 unproject, 可忽略
- `擴充上限`: 網格加密受 `UNSIGNED_SHORT` index 上限卡住 (需 `OES_element_index_uint`); FFT 解析度提高到 1024 會變 20 pass 且 texture 記憶體 x4 (約 128 MiB)

## 最棒的地方 (What is great)

- `極小但完整`: 1700 行涵蓋頻譜, GPU FFT, 渲染, 互動, 零依賴, 打開 `index.html` 就能跑, 是讀懂 FFT ocean 的最佳教材之一
- `物理頻譜而非 Phillips`: 初始頻譜的常數 (`KM=370`, `CM=0.23`, `Omega=0.84`) 與公式結構對應 Elfouhaily et al. 1997 的 unified directional spectrum (`推論`: 程式碼未引用出處, 由公式比對得出), 比常見的 Phillips spectrum 更寫實
- `GPU FFT on WebGL1`: 只靠 fragment shader + float texture 完成 Stockham FFT, 在沒有 compute shader 的年代示範了 GPGPU 做法
- `相位獨立儲存`: 改風速或尺寸時只重算振幅, 相位延續, 海面平滑過渡而不是重新洗牌
- `CSS 3D UI 與 WebGL 同步`: 控制項存在於場景空間, 是很少見且優雅的互動設計; 參數直接在海面上操作, 回饋直觀
- `效能友善`: 初始頻譜與部分 uniform 以 `changed` flag 惰性更新, framebuffer 預先建好避免切換 attachment

## 最糟的地方 (What is worst)

- `已停止維護`: 2014 年後無 commit; open issues 中有 3 則是垃圾廣告, 無人處理
- `Shader 運算子優先序 bug`: phase shader 的 `k * k / KM * KM` 實際等於 `k * k`, 與 initial spectrum 使用的 `square(k / KM)` 不一致, 短波 (capillary) 的 dispersion 被嚴重高估
- `頻譜公式疑似錯誤`: `exp(-square(sqrt(k / kp) - 1.0) / 2.0 * square(sigma))` 是乘上 `sigma^2` 而非除以 `2 * sigma^2`, peak enhancement 幾乎恆為常數 (`推論`: 依 Elfouhaily 原式比對)
- `場景有限`: 單一 2000m 網格, 無 tiling, 無 LOD, 遠處直接出現邊界; camera 的 azimuth / elevation 都被 clamp 住來掩飾這點
- `光照簡化`: 天空是單一常數色, 無環境反射, 無泡沫 (foam), 無次表面散射 (subsurface scattering), 太陽只有 diffuse 項
- `相容性老舊`: 依賴 WebGL1 + `OES_texture_float` + `OES_texture_float_linear`, 多數手機不支援 float linear filtering 就直接顯示錯誤頁; 只處理滑鼠事件, 沒有 touch
- `工程結構`: 全域變數與全域常數 (無 module), shader 以字串陣列 join 內嵌, 無測試; 字型走 `http://` 在 HTTPS 頁面會被擋
- `參數寫死`: 解析度 512, 網格 256, 顏色, 曝光皆為 `shared.js` 常數, 沒有任何 runtime 設定; 256x256 網格剛好撐滿 `UNSIGNED_SHORT` index 上限 65536, 無法再加密

## 可借鏡 (Takeaways)

- 想做 FFT ocean: 讀 `simulation.js` 的 pass 順序即可得到完整管線; 移植到 WebGL2 / WebGPU 時改用 compute shader 與 `RGBA32F` storage texture
- 要修正上述兩處 shader 公式後再拿來做物理正確性比較
- 「把 HTML UI 放進 3D 場景」的 CSS matrix 同步手法可單獨抽出重用
