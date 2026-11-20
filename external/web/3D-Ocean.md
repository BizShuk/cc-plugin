# 3D-Ocean (Seth-arc/3D-Ocean)

- Source: <https://github.com/Seth-arc/3D-Ocean>
- Demo: 無 (repo 未設 homepage, 需自行 `python -m http.server` 開 `simulation.html`)
- Author: Sethu Nguna, License: `未宣告` (GitHub 無 LICENSE 檔; README 只寫 "open source and available for educational and personal use", 法律上等同 all rights reserved)
- Stats (2026-11-20 via GitHub API): 1 star, 0 forks, 2 commits 皆在 2025-12-01 ("Initial commit", "First upload"), 之後無更新
- Size: 3 個文字檔 + 1 張 1.1 MB 截圖; `simulation.html` 778 行 (約 230 行 CSS/HTML + 約 545 行 JS), `index.html` 156 行純 landing page; 依賴 Three.js `r158` (unpkg CDN importmap), 無 build step

## 它做什麼 (What it does)

以 Three.js 官方範例的 `Water` (平面鏡反射 + 捲動 normal map 的水面材質) 與 `Sky` (Preetham 大氣散射天空盒) 組出一片開放海面, 外加:

- `第一人稱飛行`: WASD / 方向鍵水平移動, Space / Shift 升降, pointer lock 滑鼠轉向; `C` 切換到 `OrbitControls`
- `控制面板`: 5 個 slider (Wave Height / Wave Speed / Wave Length / Wind Strength / Surface Detail) 加 FPS 顯示, glassmorphism 樣式, 可摺疊
- `互動漣漪 (宣稱)`: 相機高速貼近水面或穿越水面時產生 ripple, 以 vertex shader 注入位移
- `自適應畫質 (宣稱)`: 依 60 幀平均 FPS 調整 `pixelRatio`

場景設定 (`turbidity=10`, `rayleigh=2`, `mieCoefficient=0.005`, `mieDirectionalG=0.8`, `distortionScale=3.7`, `waterColor=0x001e0f`, sun `elevation=2` / `azimuth=180`) 與 three.js `webgl_shaders_ocean` 範例的預設值相同; README 的 Acknowledgments 也承認 Water 與 Sky 來自官方範例.

### 每幀流程 (Frame Loop)

1. 更新 `time` uniform (normal map 捲動的唯一驅動力)
2. 第一人稱移動, 算速度, 視條件 push ripple (最多 10 個)
3. 把 sun light 與 fill light 移到相機相對位置
4. FPS 統計與 `pixelRatio` 調整
5. `renderer.render`: Water 的 `onBeforeRender` 先用 mirror camera 把場景 (扣掉水面) 畫進 512x512 render target, 再畫主畫面

## 效能與規模 (Performance and Scale)

數字由原始碼常數與 three r158 `Water.js` / `Sky.js` 原始碼算出; 未實際量測 FPS.

| 元素 (Element) | 規模 | 每幀成本 |
| --- | --- | --- |
| Water mesh | `PlaneGeometry(10000, 10000)`, 未指定 segments → `4 vertices, 2 triangles` | 1 draw call; fragment shader 取 normal map 4 次 + mirror 1 次 |
| Mirror reflection | 512x512 render target (`textureWidth/Height`) | 整個場景再畫一次 (實際只有 Sky 1 draw) |
| Sky | `BoxGeometry` scale 10000, 12 triangles, vertex shader 強制 `z = w` 貼在 far plane | 主畫面 + mirror 各 1 draw; 大氣散射是 per-fragment 計算 |
| Shadow map | 2048x2048, `PCFSoftShadowMap` | 有開啟但場景`沒有任何 castShadow / receiveShadow 物件`, 純浪費 |
| Lights | 1 directional (sun, castShadow) + 1 ambient + 1 directional fill | 參與 Water 的 lights uniform |
| Ripple 陣列 | 上限 10 個 vec3 + float | uniform 更新, 但 shader 端從未讀取 (見下節) |
| UI | 5 個 range slider + FPS 文字 + 返回按鈕 | 每 30 幀改 1 次 DOM 文字 |

- `Draw calls`: 每幀約 3 個 (主畫面 Sky + Water, mirror 的 Sky), 幾何成本趨近零; 成本集中在 `fill rate` (像素填充率): dpr 2 的 1080p 螢幕約 830 萬像素, Water 與 Sky 都是滿版 fragment shader, 加上 `antialias: true` 的 MSAA
- `GPU memory` (`推論`, 依 render target 尺寸估算): shadow map 2048² RGBA8 + depth 約 `32 MiB`, mirror target 約 2 MiB, normal map 貼圖 (`waternormals.jpg`, 從 threejs.org 載入) 若為 1024² 約 4 MiB (含 mipmap 約 5.3 MiB)
- `海面尺度`: 平面 10 km 見方, 但相機 `far = 500`, 500 單位以外的水面被裁掉, 由 far plane 上的 Sky 補位; 遠方海平線其實是裁切邊界而不是真正的地平線
- `海浪尺度`: 沒有任何幾何位移, 浪完全是 normal map 4 個 octave 的 UV 捲動 (`uv/103`, `uv/107`, `uv/8907`, `uv/1091` 四種尺度), 所以`海面輪廓永遠是平的`, 只有光影在動
- `CPU 端`: 每幀 `updateLighting` / `updateSun` / `updateMovement` 會 `new` 或 `clone` 約 6 個 `Vector3`, 產生少量 GC 壓力; 其餘可忽略
- `擴充上限`: 因為海浪不是模擬出來的, 規模與畫質只取決於螢幕解析度與 mirror target 大小, 與海域面積無關; 要真正有浪就必須換掉整個水面方案 (加 segments + Gerstner / FFT displacement)

## 最棒的地方 (What is great)

- `入門門檻極低`: 單一 HTML + CDN importmap, 開 local server 就能跑, 是「Three.js Water + Sky 範例加上第一人稱操控」的最小可運作版本
- `視覺效價比高`: 3 個 draw call 就有反射, Fresnel, 夕陽大氣散射與 ACES tone mapping, 在中低階 GPU 上也容易達到 60 FPS (`推論`: 依 draw call 與 shader 複雜度判斷, 未實測)
- `操控完整`: 第一人稱 pointer lock + Orbit 雙模式切換, 鍵盤狀態以 key map 管理, 移動以 `deltaTime` 縮放, 與 frame rate 無關
- `UI 打磨`: glassmorphism 面板, 摺疊動畫, 內嵌 FPS; landing page 與 simulation 分離, 適合做 demo 展示
- `README 齊全`: 操作說明, 參數範圍, 本地啟動指令, 已知問題與 roadmap 都有寫

## 最糟的地方 (What is worst)

- `互動漣漪是死碼`: 第二個 `shader.replace` 要找的字串 `vec4 worldPosition = modelMatrix * vec4( transformed, 1.0 );` 不存在於 r158 `Water.js` vertex shader (它用的是 `mirrorCoord = modelMatrix * vec4( position, 1.0 )`), 替換靜默失敗, `calculateUserWave` 從未被呼叫; 就算替換成功, 水面只有 4 個頂點 (位於 ±5000 角落), 半徑 50 / 30 的漣漪也打不到任何頂點. README 主打的 "Interactive Water Physics" 實際不存在
- `5 個 slider 有 4 個名不副實`:
  - `Wave Height` 改的是 `size` uniform, 在 shader 中是 normal map 的 UV 縮放 (頻率), 不是浪高
  - `Wave Speed` 只更新數字顯示, handler 註解寫 "handled in animation loop" 但 loop 裡沒有任何讀取
  - `Wave Length` 與 `Surface Detail` 都寫入同一個 `distortionScale`, 互相覆蓋, 後拉的那個生效
  - `Wind Strength` 寫入 `alpha`, 但 Water 材質未設 `transparent` 且 renderer `alpha: false`, 畫面沒有任何變化
- `自適應畫質有缺陷`: 降級時 `pixelRatio` 被 `Math.min(dpr * q, 1)` 直接壓到 1 (dpr 2 螢幕一次掉 4 倍像素), 回升時上限又回到 2, 可能在兩者間來回震盪 (`推論`); 每次 `setPixelRatio` 都會重配 drawing buffer 造成卡頓; README 宣稱的「shadow map size and texture resolution dynamically scaled」在程式碼中不存在
- `無效的開銷與殘留`: 開了 2048² 軟陰影但沒有物件投影或受影; `PMREMGenerator` 建立後從未使用 (官方範例用它從 Sky 生成環境貼圖); body 的 CSS 漸層背景被不透明 canvas 完全蓋住; `outputEncoding` / `sRGBEncoding` 在 r152 後已 deprecated
- `沒有真正的海浪`: 無 displacement, 無 FFT 或 Gerstner, 無泡沫, 無水下視角; 標題的 "photorealistic ocean simulation" 實際是 shading-only 的鏡面水面
- `工程結構`: 全部 JS 擠在單一 HTML, 全域狀態, 縮排混亂, 無測試; 執行期依賴 unpkg 與 threejs.org 的外部資源, 離線或 CDN 失效就無法顯示 normal map; 無 touch 支援, 手機無法操作
- `授權不明`: 沒有 LICENSE, 不能直接拿來商用或改作

## 可借鏡 (Takeaways)

- 想快速做「漂亮但不需要物理」的水面背景: Three.js `Water` + `Sky` 即可, 這個 repo 是最小組裝範例; 但記得把 `PMREMGenerator` 接上 `scene.environment` 並關掉無用的 shadow map
- 對 three.js 內建 shader 做 `String.replace` 注入時, 必須驗證替換是否命中 (比對替換前後字串), 否則會像本專案一樣靜默失效; 更穩的做法是 `onBeforeCompile`
- 需要真正的浪高與互動, 參考 [waves.md](waves.md) 的 GPU FFT 管線, 並把水面換成高 segments 的網格
