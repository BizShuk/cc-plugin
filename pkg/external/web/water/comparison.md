# 網頁水面實作比較 (Web Water Implementations Comparison)

比較 [waves.md](waves.md), [webgl-water.md](webgl-water.md), [3D-Ocean.md](3D-Ocean.md),
[Threejs-WebGPU-IFFT-Ocean.md](Threejs-WebGPU-IFFT-Ocean.md) 四種做法, 並納入其他沿 `waves` 路線發展的 repo,
判斷哪一種在 `高效能 + 大量物件` 下擴充能力最好.

- 判斷依據: 2026-11-20 shallow clone 五個 repo 後讀原始碼與 GitHub API; `未實測 FPS`, 成本皆由原始碼常數推導
- 「物件」指漂浮或與水互動的場景物體 (船, 浮標, 碎片), 不是水面網格本身

## 結論 (Verdict)

- `架構路線勝出: 頻譜 + FFT (waves 系)`. 水面成本 `O(R² log R)` 與物件數 N, 海域面積都無關; displacement map 是週期函數, 天然可鋪磚 (tiling); 物件只需「取樣高度」這個單向耦合, 可整批放到 GPU
- `程式碼本身勝出: Threejs-WebGPU-IFFT-Ocean, 但需先修三處`. `dli/waves` 是 WebGL1 + 全域變數 + 無 scene graph, 要加任何物件都得自寫 renderer; [Spiri0/Threejs-WebGPU-IFFT-Ocean](https://github.com/Spiri0/Threejs-WebGPU-IFFT-Ocean) 是唯一同時具備 compute IFFT, 多 cascade, LOD 與 three.js scene graph 的 repo. 但它自己的水面已用掉 `748 draw call` 與約 `95 MiB` GPU memory, 且有 chunk 記憶體洩漏; 不先修, 加物件的預算會被水面吃掉
- `webgl-water` 的價值在`局部雙向互動` (物體推水), 但每個 pixel 的成本隨 N 線性成長, 只能當「相機附近的局部 wake 層」疊加
- `3D-Ocean` 只有 three.js scene graph 這一點可取; 水面是平的, 沒有可擴充的物理核心

## 四種做法的本質 (Approaches)

| 面向 | waves (dli) | webgl-water (evanw) | 3D-Ocean (Seth-arc) | Threejs-WebGPU-IFFT-Ocean (Spiri0) |
| --- | --- | --- | --- | --- |
| 水面模型 | 頻譜 + GPU FFT (Tessendorf) | 波動方程 heightfield | 無模擬, normal map 捲動 | JONSWAP + TMA 頻譜, 3 cascade IFFT |
| 有無真實浪高 | 有 (含 choppy 水平位移) | 有 (小振幅) | `無`, 幾何永遠是平面 | 有, 250 m 長浪到 1 cm 細節 + Jacobian 泡沫 |
| 物件耦合 | `無物件` | 雙向, 但寫死 1 顆球 | 無 (宣稱的 ripple 是死碼) | `無` (issue #7 詢問浮力, 無回應) |
| 場景範圍 | 單一 2000m 網格, 無 tiling / LOD | 固定 2x2 泳池 | 10km 平面, far=500 裁切 | 1,000 km quadtree LOD, 實際可見約 4 km (霧) |
| 框架 | 原生 WebGL1 | 自家 `lightgl.js` | three.js r158 (WebGL) | three.js 0.179 WebGPU + TSL, 無 WebGL fallback |
| 每幀 draw call | 22 | 9 | 約 3 | 234 compute dispatch + 約 748 draw (相機 5 m, 視錐裁切前) |
| GPU memory | 約 32 MiB | 約 7 MiB | 約 40 MiB (多為無用 shadow map) | 約 95 MiB cascade + 57 MiB chunk geometry |
| 主要成本 | FFT 18 pass (約 5.5M fragment) | 螢幕解析度 x 解析 ray tracing | fill rate + mirror 重繪 | IFFT dispatch 數 + chunk draw call |

## 物件數 N 增加時的擴充性 (Scaling with N Objects)

### webgl-water: `O(pixels x N)`, 最差

原始碼中球是`單一 uniform` (`sphereCenter`, `sphereRadius`), 散佈在 5 個 shader. 擴到 N 個物件時:

- 水面 fragment shader 的反射與折射各 1 條射線, 每條都要對每個物件求交 -> 每 pixel `2N` 次求交; ambient occlusion 也是逐物件相乘
- `moveSphere` 每個物件 1 個 fullscreen pass (65,536 fragment), 或改成 shader 內迴圈 N 次
- caustic pass 的球影同樣逐物件計算
- 解析 ray tracing 只對 sphere / box 有封閉解, 任意 mesh 無法加入

結論: 成本與螢幕像素和 N 相乘, N 到數十就會掉幀 (`推論`, 依 shader 結構推得). 但「加舊體積, 扣新體積」的排水手法本身可改成 `instanced splat` (每物件畫一個小 quad 到 heightfield), 成本變成 `O(N x 物件足跡)`, 這部分值得保留.

### 3D-Ocean: 渲染成本 `2x`, 物理成本 `0` (因為沒有物理)

- three.js scene graph 加物件最容易, `InstancedMesh` 可把 N 個同型物件壓成 1 個 draw call
- 但 `Water` 的 mirror reflection 會把整個場景再畫一次, 每個物件的頂點成本都 `x2`
- 水面沒有高度, 物件無從「浮」起來; 要有浮力就得換掉整個水面方案

### waves 系 (FFT): 水面 `O(1)`, 物件 `O(N)` 且可全 GPU

- 水面成本只看 FFT 解析度 R, 與 N 和海域大小無關
- 物件浮力 = 在物件位置取樣 displacement map, 屬單向耦合 (水推物件, 物件不推水), 不需要每物件額外 render pass
- 原始 `dli/waves` 把 displacement map 設成 `CLAMP_TO_EDGE`, 只畫單一網格; FFT 輸出本身是週期函數, 改成 `REPEAT` 即可無限鋪磚

### Threejs-WebGPU-IFFT-Ocean: 水面 `O(1)` 但基線偏高, 物件 `O(N)` 待補

它解掉 `dli/waves` 的範圍問題 (quadtree LOD), 物件也能直接放進 three.js scene. 但擴充前有三個基線成本:

- `水面本身就是「大量物件」`: 每個 chunk 是獨立 `Mesh`, 相機 5 m 時 748 個 draw call; 再加 N 個物件時, draw call 預算多半已被水面用掉. 所有 chunk 的本地網格都是同一張 36x36, `推論`: 改成單一 `InstancedMesh` 可收成 1 個 draw call, worker 也能移除
- `記憶體會一直長`: `src/` 沒有任何 `dispose()`, 退役 chunk 的 `BufferGeometry` 留在 GPU 上. 長時間飛行加上 N 個物件會更快耗盡記憶體
- `compute 管線偏碎`: 3 層 x 78 = 234 個 dispatch/frame; `推論`: 用 workgroup shared memory 每方向 1 個 dispatch, 可降到每層約 4 個

浮力的接點已經有了: displacement 是 storage texture, 加一個 compute pass 對 N 個取樣點讀值即可, 不需要 CPU readback. 但它的 vertex shader 依距離用 `LOD_SCALE = 3.7` 衰減位移 (5 m 層 18.5 m 外就沒了), 浮力取樣要讀`未衰減`的原值, 否則遠方物件會與畫面上的浪對不上

## 新一代 waves 系 repo (Modern Successors)

| Repo | Stars / 最後 push | 框架 | 擴充相關能力 | 缺口 |
| --- | --- | --- | --- | --- |
| [jbouny/fft-ocean](https://github.com/jbouny/fft-ocean) | 354 / 2015-11-12 | three.js (舊版, WebGL) | `dli/waves` 的 three.js 移植 | 無 buoyancy, 已停更 |
| [Popov72/OceanDemo](https://github.com/Popov72/OceanDemo) | 152 / 2024-07-23 | Babylon.js 5.28 WebGPU | 3 個 cascade (`lengthScale` 250 / 17 / 5m), WGSL compute FFT, `Buoyancy.addMesh()` 浮力 API | 浮力在 CPU |
| [Spiri0/Threejs-WebGPU-IFFT-Ocean](https://github.com/Spiri0/Threejs-WebGPU-IFFT-Ocean) (詳見 [筆記](Threejs-WebGPU-IFFT-Ocean.md)) | 143 / 2025-12-12 | three.js 0.179 WebGPU (TSL) | 512 IFFT, 3 cascade, CDLOD quadtree + morphing, 7 個 Web Worker 建網格, Jacobian 泡沫; `推論` 是 Unity `gasgiant/FFT-Ocean` 的移植 | `無 buoyancy`, 無 `dispose()`, 748 個獨立 chunk `Mesh`, 只能鍵盤操作 |

`Popov72/OceanDemo` 的浮力路徑 (讀 `wavesGenerator.ts` 與 `buoyancy.ts` 確認):

1. 每次 `update()` 都對 cascade 0 的 displacement 做非同步 `readPixels`, 結果存成 `Uint16Array` (half float)
2. CPU 端每個物件取 3 個 frame 點, 每點做 4 次 displacement 反推迭代求水面高度, 再算傾角
3. 成本: CPU `O(N x 3 x 4)`, 每幀固定一次 GPU->CPU 讀回; 只讀最大 cascade, 5m / 17m 的短浪不影響浮力; 讀回非同步, 浮力資料延遲 1 幀以上 (`推論`)

N 到數百仍可行 (`推論`), 但 readback 是管線同步點, 而且物件在 CPU 端逐一更新矩陣, N 再大就卡在 JS.

## 最佳擴充架構 (Recommended Architecture)

以 `Spiri0/Threejs-WebGPU-IFFT-Ocean` 為底, 依層疊加:

1. `水面核心`: 保留其 3-cascade IFFT + quadtree LOD (已具備), 加物件前先修三處: chunk 改單一 `InstancedMesh` 並移除 worker, 退役 chunk 呼叫 `dispose()`, IFFT 改 shared-memory 每方向 1 dispatch
2. `GPU buoyancy`: 用 compute shader 對 N 個物件的取樣點讀`未經 LOD 衰減`的 displacement (同 OceanDemo 的 4 次反推迭代), 結果直接寫進 `InstancedMesh` 的 instance matrix storage buffer; 全程不讀回 CPU, 成本 `O(N)` 且 N 可到上萬 (`推論`)
3. `需要 CPU 知道位置的少數物件` (玩家船, 碰撞判定): 才走 OceanDemo 式非同步 readback, 只讀取樣點而非整張 texture
4. `局部雙向互動`: 在相機附近放一塊 `webgl-water` 式 heightfield, 物件以 instanced splat 推水, 結果疊加到 FFT displacement 上 (類似 Unity Crest 的 dynamic waves 做法); 修正 `webgl-water` 的兩個缺陷: step 數依 `dt` 固定步長, 擾動正負成對以守恆水量
5. `反射`: 不用 3D-Ocean 的 mirror 重繪 (每物件成本 x2), 改用 Sky 產生的 PMREM 環境貼圖 + screen-space reflection

| 能力 | 取自 |
| --- | --- |
| 頻譜模型, 相位與振幅分離 | `dli/waves` |
| 多 cascade, compute FFT, 浮力演算法 | `Popov72/OceanDemo` |
| three.js WebGPU, JONSWAP 頻譜, CDLOD morphing, Jacobian 泡沫 | `Spiri0/Threejs-WebGPU-IFFT-Ocean` |
| 物體排水, caustics 面積比 | `evanw/webgl-water` |
| Sky 大氣散射 (改接 PMREM) | `Seth-arc/3D-Ocean` / three.js 範例 |

## 效能面總結 (Performance Conclusion)

各面向數字取自四份筆記與本文前述分析, 全部由原始碼常數推導, `未實測 FPS`. `Spiri0` 欄的 draw call 與幾何量以相機高 5 m, 視錐裁切前計.

### 🔹 模擬運算 (Simulation Compute)

| 項目 | waves | webgl-water | 3D-Ocean | Spiri0 IFFT-Ocean |
| --- | --- | --- | --- | --- |
| 模擬格點 | 512² | 256² | `無模擬` | 3 x 512² = 786,432 |
| 每幀 pass / dispatch | 21 offscreen pass (initial spectrum 只在參數變動時) | 4 fullscreen + 1 caustic (1024²) | 0 (只更新 `time` uniform) | 234 compute dispatch |
| 每幀 shader invocation | 約 5.5M fragment | 約 0.26M fragment + caustic | 0 | 約 6,100 萬 |
| 解析度加倍的代價 | 1024: 20 pass, texture 記憶體 x4 | 波速以格/step 計, 必須同步調 step 數 | 不適用 | 1024: 每層 78 -> 86 dispatch, threads x4 |

### 🔹 GPU 記憶體 (GPU Memory)

| 項目 | waves | webgl-water | 3D-Ocean | Spiri0 IFFT-Ocean |
| --- | --- | --- | --- | --- |
| 模擬 / 貼圖 | 32 MiB (8 張 RGBA float 512²) | 約 7 MiB | 約 40 MiB, 其中 32 MiB 是無用 shadow map | 約 95 MiB (3 cascade) |
| 幾何 | 約 2 MB | 約 1 MB | 可忽略 (4 vertices) | 約 57 MiB (748 chunk) |
| 長時間穩定性 | 固定配置 | 固定配置 | 固定配置 | `持續成長` (無 `dispose()`) |

### 🔹 繪製與 CPU (Draw Calls and CPU)

| 項目 | waves | webgl-water | 3D-Ocean | Spiri0 IFFT-Ocean |
| --- | --- | --- | --- | --- |
| 每幀 draw call | 22 | 9 (拖曳時每次 mousemove +1) | 約 3 | 約 748 + 234 dispatch |
| 水面幾何 | 65,536 vertices, 130K triangles | 40,401 vertices, 80K triangles, 每幀畫 3 次 | 4 vertices, 2 triangles | 102 萬 vertices, 194 萬 triangles |
| CPU 端 | uniform 更新 + 1 次 CSS matrix | 球的物理積分 | 每幀 new 數個 `Vector3`, 少量 GC | quadtree 重算 + 7 worker 以結構化複製回傳 |

### 🔹 像素填充 (Fill Rate per Pixel)

| 項目 | waves | webgl-water | 3D-Ocean | Spiri0 IFFT-Ocean |
| --- | --- | --- | --- | --- |
| 水面 fragment 成本 | normal map + Fresnel + tone mapping, 輕 | 6 次高度 fetch + 2 條解析射線, `最重` | normal map 4 次 + mirror 1 次 + 大氣散射 sky | 6 次 `textureSample` + 1 次 cubemap |
| 額外全畫面成本 | 無 | 無 | mirror 512² 重畫場景, MSAA | MSAA, 748 mesh 走透明排序 |
| 高 DPR 螢幕 | 筆記未記錄 | fragment 數 x DPR² | 自適應降級會在 DPR 1 與 2 間震盪 | fragment 數 x DPR² |

### 🔹 擴充性 (Scalability)

| 項目 | waves | webgl-water | 3D-Ocean | Spiri0 IFFT-Ocean |
| --- | --- | --- | --- | --- |
| N 個物件 | 無物件; 加入後 `O(N)`, 可全 GPU | `O(pixels x N)`, 數十個即吃緊 (`推論`) | 渲染 x2 (mirror), 沒有浮力可算 | 水面不隨 N 增加; 物件 `O(N)` 待實作, draw call 預算已被水面佔用 |
| 海域面積 | 單一 2000 m, `Uint16` index 上限 | 固定泳池 | 無關 (水面是平的) | 1,000 km LOD, 實際可見約 4 km |
| 幀率獨立 | 筆記未記錄 | `綁 FPS` (每幀固定 2 step) | 移動依 `deltaTime` | 固定 60 Hz step, 120 Hz 螢幕會有位移階梯感 (`推論`) |
| 平台 | WebGL1 + float linear, 多數手機不支援 | WebGL1, 有 half float 降級; 手機水面停不下來 | three.js WebGL, 無 touch | 只支援 WebGPU, 只能用鍵盤 |

### 🔹 綜合排名 (Overall Ranking)

`1` 最好, `4` 最差; 排名以「高效能 + 大量物件」為目標.

| 面向 | waves | webgl-water | 3D-Ocean | Spiri0 IFFT-Ocean |
| --- | --- | --- | --- | --- |
| 模擬效率 (每單位浪況的成本) | 2 | 3 | 4 (沒有浪) | 1 |
| 現況 GPU / draw call 負擔 | 2 | 1 | 3 | 4 |
| 物件數擴充 | 2 | 4 | 3 | 1 |
| 海域擴充 | 3 | 4 | 2 (平面) | 1 |
| 長時間穩定性 | 1 | 3 (水位上升) | 2 | 4 (記憶體洩漏) |
| 作為擴充起點 | 3 (WebGL1, 無 scene graph) | 4 | 2 (只有 scene graph) | `1` |

結論: 現在最省資源的是 `webgl-water` 與 `waves`, 但兩者的上限都被寫死 (泳池, 單一網格, `Uint16` index). 物件與海域都放大時, 只有 `Spiri0 IFFT-Ocean` 的成本不隨兩者成長; 它目前效能最差的兩項 (748 draw call 與記憶體洩漏) 都屬實作問題, 用 `InstancedMesh` 與 `dispose()` 可以解掉 (`推論`), 演算法本身不需要換.

## 未驗證事項 (Unverified)

- 五個 repo 都沒有在本機跑過; FPS, N 的實際上限都是由原始碼結構推得
- `Spiri0/Threejs-WebGPU-IFFT-Ocean` 的 748 chunk 是視錐裁切前的數字, 實際 draw call 與記憶體洩漏速度未量測; iOS 支援只看了論壇描述
- 「GPU buoyancy 寫 instance matrix」在 three.js TSL 的可行寫法未實作驗證
