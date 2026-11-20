# 網頁水面實作比較 (Web Water Implementations Comparison)

比較 [waves.md](waves.md), [webgl-water.md](webgl-water.md), [3D-Ocean.md](3D-Ocean.md) 三種做法,
並納入兩個沿 `waves` 路線發展的新一代 repo, 判斷哪一種在 `高效能 + 大量物件` 下擴充能力最好.

- 判斷依據: 2026-11-20 shallow clone 五個 repo 後讀原始碼與 GitHub API; `未實測 FPS`, 成本皆由原始碼常數推導
- 「物件」指漂浮或與水互動的場景物體 (船, 浮標, 碎片), 不是水面網格本身

## 結論 (Verdict)

- `架構路線勝出: 頻譜 + FFT (waves 系)`. 水面成本 `O(R² log R)` 與物件數 N, 海域面積都無關; displacement map 是週期函數, 天然可鋪磚 (tiling); 物件只需「取樣高度」這個單向耦合, 可整批放到 GPU
- `程式碼本身勝出: 都不是原始三個`. `dli/waves` 是 WebGL1 + 全域變數 + 無 scene graph, 要加任何物件都得自寫 renderer; 實際起點應選 [Spiri0/Threejs-WebGPU-IFFT-Ocean](https://github.com/Spiri0/Threejs-WebGPU-IFFT-Ocean), 再補 GPU buoyancy
- `webgl-water` 的價值在`局部雙向互動` (物體推水), 但每個 pixel 的成本隨 N 線性成長, 只能當「相機附近的局部 wake 層」疊加
- `3D-Ocean` 只有 three.js scene graph 這一點可取; 水面是平的, 沒有可擴充的物理核心

## 三種做法的本質 (Approaches)

| 面向 | waves (dli) | webgl-water (evanw) | 3D-Ocean (Seth-arc) |
| --- | --- | --- | --- |
| 水面模型 | 頻譜 + GPU FFT (Tessendorf) | 波動方程 heightfield | 無模擬, normal map 捲動 |
| 有無真實浪高 | 有 (含 choppy 水平位移) | 有 (小振幅) | `無`, 幾何永遠是平面 |
| 物件耦合 | `無物件` | 雙向, 但寫死 1 顆球 | 無 (宣稱的 ripple 是死碼) |
| 場景範圍 | 單一 2000m 網格, 無 tiling / LOD | 固定 2x2 泳池 | 10km 平面, far=500 裁切 |
| 框架 | 原生 WebGL1 | 自家 `lightgl.js` | three.js r158 |
| 每幀 draw call | 22 | 9 | 約 3 |
| 主要成本 | FFT 18 pass (約 5.5M fragment) | 螢幕解析度 x 解析 ray tracing | fill rate + mirror 重繪 |

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

## 新一代 waves 系 repo (Modern Successors)

| Repo | Stars / 最後 push | 框架 | 擴充相關能力 | 缺口 |
| --- | --- | --- | --- | --- |
| [jbouny/fft-ocean](https://github.com/jbouny/fft-ocean) | 354 / 2015-11-12 | three.js (舊版, WebGL) | `dli/waves` 的 three.js 移植 | 無 buoyancy, 已停更 |
| [Popov72/OceanDemo](https://github.com/Popov72/OceanDemo) | 152 / 2024-07-23 | Babylon.js 5.28 WebGPU | 3 個 cascade (`lengthScale` 250 / 17 / 5m), WGSL compute FFT, `Buoyancy.addMesh()` 浮力 API | 浮力在 CPU |
| [Spiri0/Threejs-WebGPU-IFFT-Ocean](https://github.com/Spiri0/Threejs-WebGPU-IFFT-Ocean) | 143 / 2025-12-12 | three.js ^0.179 WebGPU (TSL) | 512 IFFT, 3 cascade, `quadtree LOD` 水面區塊, 7 個 Web Worker 建網格, entity-manager 架構 | `無 buoyancy` (原始碼 grep 無結果) |

`Popov72/OceanDemo` 的浮力路徑 (讀 `wavesGenerator.ts` 與 `buoyancy.ts` 確認):

1. 每次 `update()` 都對 cascade 0 的 displacement 做非同步 `readPixels`, 結果存成 `Uint16Array` (half float)
2. CPU 端每個物件取 3 個 frame 點, 每點做 4 次 displacement 反推迭代求水面高度, 再算傾角
3. 成本: CPU `O(N x 3 x 4)`, 每幀固定一次 GPU->CPU 讀回; 只讀最大 cascade, 5m / 17m 的短浪不影響浮力; 讀回非同步, 浮力資料延遲 1 幀以上 (`推論`)

N 到數百仍可行 (`推論`), 但 readback 是管線同步點, 而且物件在 CPU 端逐一更新矩陣, N 再大就卡在 JS.

## 最佳擴充架構 (Recommended Architecture)

以 `Spiri0/Threejs-WebGPU-IFFT-Ocean` 為底, 依層疊加:

1. `水面核心`: 保留其 3-cascade IFFT + quadtree LOD (已具備)
2. `GPU buoyancy`: 用 compute shader 對 N 個物件的取樣點讀 displacement (同 OceanDemo 的 4 次反推迭代), 結果直接寫進 `InstancedMesh` 的 instance matrix storage buffer; 全程不讀回 CPU, 成本 `O(N)` 且 N 可到上萬 (`推論`)
3. `需要 CPU 知道位置的少數物件` (玩家船, 碰撞判定): 才走 OceanDemo 式非同步 readback, 只讀取樣點而非整張 texture
4. `局部雙向互動`: 在相機附近放一塊 `webgl-water` 式 heightfield, 物件以 instanced splat 推水, 結果疊加到 FFT displacement 上 (類似 Unity Crest 的 dynamic waves 做法); 修正 `webgl-water` 的兩個缺陷: step 數依 `dt` 固定步長, 擾動正負成對以守恆水量
5. `反射`: 不用 3D-Ocean 的 mirror 重繪 (每物件成本 x2), 改用 Sky 產生的 PMREM 環境貼圖 + screen-space reflection

| 能力 | 取自 |
| --- | --- |
| 頻譜模型, 相位與振幅分離 | `dli/waves` |
| 多 cascade, compute FFT, 浮力演算法 | `Popov72/OceanDemo` |
| three.js WebGPU, LOD, Worker 架構 | `Spiri0/Threejs-WebGPU-IFFT-Ocean` |
| 物體排水, caustics 面積比 | `evanw/webgl-water` |
| Sky 大氣散射 (改接 PMREM) | `Seth-arc/3D-Ocean` / three.js 範例 |

## 未驗證事項 (Unverified)

- 五個 repo 都沒有在本機跑過; FPS, N 的實際上限都是由原始碼結構推得
- `Spiri0/Threejs-WebGPU-IFFT-Ocean` 的 three.js 版本相依 (`^0.179`) 與 iOS 支援只看了 package.json 與論壇描述, 未在裝置上測
- 「GPU buoyancy 寫 instance matrix」在 three.js TSL 的可行寫法未實作驗證
