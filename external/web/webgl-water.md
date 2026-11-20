# webgl-water (evanw/webgl-water)

- Source: <https://github.com/evanw/webgl-water>
- Demo: <http://madebyevan.com/webgl-water/> / video <http://youtube.com/watch?v=R0O_9bp3EKQ>
- Caustics 技術說明: <https://medium.com/@evanwallace/rendering-realtime-caustics-in-webgl-2a99a29a0b2c>
- Author: Evan Wallace (esbuild / Figma 共同創辦人), License: 檔頭宣告 MIT, 但 repo 沒有 LICENSE 檔 (GitHub API `license: null`)
- Stats (2026-11-20 via GitHub API): 1286 stars, 270 forks, 15 open issues (含 2 個 PR), last commit 2016-01-07, last push 2022-12-31
- Size: 應用程式約 840 行 (`main.js` 281 + `renderer.js` 353 + `water.js` 172 + `cubemap.js` 33), 另附作者自己的 `lightgl.js` (2068 行) 與 float linear polyfill (130 行); 無 build step; 7 張 JPEG 共 232 KB

## 它做什麼 (What it does)

瀏覽器內即時的「泳池 + 球」水面互動 demo. 一個 2x2 見方的磁磚泳池, 水面以 `heightfield`(每個格點只存高度與垂直速度的 2D 網格) 波動方程模擬, 一顆球可被拖動或受重力落水並排開水面.

互動:

- 在水面上拖曳 -> 產生漣漪
- 拖曳背景 -> 旋轉 camera; 拖曳球 -> 移動球
- `SPACE` 暫停, `G` 切換球的重力, 按住 `L` -> 光源方向改成目前視線方向
- 支援 touch 與 `devicePixelRatio`

### 管線 (Pipeline), 每幀依序執行

1. `moveSphere` - fullscreen pass, 依球的舊/新位置對水面「加回舊體積, 扣掉新體積」, 做出排水效果
2. `stepSimulation` x2 - fullscreen pass, 每格速度往四鄰平均高度拉 (`v += (avg - h) * 2`), 乘 `0.995` 阻尼, 再把速度積分進高度; ping-pong 兩張 float texture
3. `updateNormals` - 從高度差分出法向量, 存回同一張 texture 的 `ba` channel (texel 內容 = 高度, 速度, normal.x, normal.z)
4. `updateCaustics` - 把 200x200 水面網格沿折射光線投影到池底, fragment shader 用 `dFdx/dFdy` 比較投影前後三角形面積, 「面積變小 = 光被聚焦 = 變亮」, 寫進 1024x1024 caustic texture; 同一個 pass 順便算球的 blob shadow 與池緣陰影
5. 主畫面: 池壁 (cube 拿掉頂面) -> 水面 x2 (一次 cull front 畫水上視角, 一次 cull back 畫水下視角) -> 球

水面 fragment shader 不做 reflection / refraction render pass, 而是在 shader 內對 `解析場景`(只有軸對齊 box 與一顆球, 都有封閉解交點公式) 直接 `ray trace`: 反射與折射各打一條射線, 命中球就算球色, 命中池壁就取磁磚 + caustic, 打出泳池就取 cubemap 天空 + 太陽高光, 最後以 Fresnel 混合. Ambient occlusion 也是對球與池壁的解析公式.

## 效能與規模 (Performance and Scale)

數字由原始碼常數算出 (`Texture(256, 256)`, `Mesh.plane({detail: 200})`, `Mesh.sphere({detail: 10})`, `Texture(1024, 1024)`); 未實際量測 FPS.

| 元素 (Element) | 規模 | 每幀成本 |
| --- | --- | --- |
| Water simulation grid | 256x256 = 65,536 格, float32 RGBA ping-pong 2 張 | 4 個 fullscreen pass (`moveSphere` 1 + `stepSimulation` 2 + `updateNormals` 1) = 262,144 次 fragment 呼叫; update pass 每格讀 5 個 texel |
| Water surface mesh | 201x201 = 40,401 vertices, 80,000 triangles, 240,000 indices (`Uint16`) | 每幀畫 3 次 (caustics 1 + 水上/水下 2), 約 121K 次 vertex shader 呼叫 |
| Caustic map | 1024x1024 RGBA8 | 1 pass, 每幀都重算 (水面靜止也算) |
| Sphere mesh | 402 vertices, 800 triangles (八面體細分 n=10) | 1 draw call |
| Pool (cube) | 5 面 10 triangles | 1 draw call |
| 滑鼠漣漪 | 半徑 0.03 的 cosine bump | 每個 `mousemove` 額外 1 pass |

- `Draw calls`: 每幀 9 個 (offscreen 5 + 主畫面 4), 拖曳時每次 mousemove +1
- `GPU memory`: 模擬 2 x 1 MiB (float32; 降級 half float 時 2 x 0.5 MiB) + caustic 4 MiB + 磁磚 512x512 + cubemap 6 x 256x256, 合計約 `7 MiB`; 水面 vertex buffer 約 485 KB, index buffer 約 480 KB
- `真正的瓶頸是螢幕解析度`: 水面 shader 每個 pixel 要 6 次高度 fetch (1 次 + 5 次「peaked」迭代) 再加 2 條射線各含球/box 求交與最多 3 次 texture fetch; canvas 尺寸乘上 `devicePixelRatio`, 所以 Retina 螢幕的 fragment 數是 DPR=1 的 4 倍; 模擬本身 (256²) 相對便宜
- `波速與衰減` (由公式推導): 更新式等價於 `c²Δt²/Δx² = 0.5`, 剛好落在 2D 顯式波動方程的穩定上限; 波速約 0.71 格/step, 每幀 2 step, 在 60 Hz 下約 3 秒橫越整個池; 阻尼 `0.995^138 ≈ 0.5`, 振幅約 1.15 秒衰減一半
- `CPU 端`: 只有球的重力/浮力/黏滯積分與 uniform 更新; 每次滑鼠事件做 1 次 ray-plane 或 ray-sphere 求交, 可忽略
- `擴充上限`: 水面網格 201² 已用掉 `Uint16` index 上限 65,536 的 62%, `detail` 最多只能到 255 (256² = 65,536 vertices), 再密就要 `OES_element_index_uint`; 模擬格數提高只線性增加 fullscreen pass 成本, 但波速以「格/step」計, 格數加倍會讓波看起來慢一倍, 必須同步調 step 數

## 最棒的地方 (What is great)

- `Shader 內解析 ray tracing`: 場景刻意限制成 box + sphere, 換來在 fragment shader 內直接算出正確的反射, 折射, 二次命中與 ambient occlusion, 不需任何額外 render pass 或 cubemap 重繪; 是「限制場景換取品質」的經典示範
- `Caustics 面積比技巧`: 用 `dFdx/dFdy` 比較折射前後三角形面積來估光強, 一個 pass 就得到可信的焦散紋, 作者另寫了專文說明, 後來被大量教學引用
- `雙向耦合`: 球的移動以「加舊體積, 扣新體積」推開水面, 球的浮力依沒入比例計算, 水下另有黏滯減速; 互動感遠超只能點漣漪的 demo
- `水上/水下同一網格`: 以 face culling 切換兩支 shader, 從水下仰看時有全反射與水面 Fresnel, 不需另建幾何
- `極小且零依賴`: 約 840 行應用程式碼 + 自家 `lightgl.js`, 開 `index.html` 即跑; 相容處理周到 (float -> half float 降級, linear filtering 能力實測 polyfill, `OES_standard_derivatives` 缺席時退化成固定亮度, touch, DPR)

## 最糟的地方 (What is worst)

- `模擬速度綁定 FPS`: `stepSimulation()` 每幀固定呼叫 2 次, 不看 `dt`(只有球的物理用 `dt`); 120 Hz 螢幕上的水波速度是 60 Hz 的 2 倍 (`推論`: 由程式碼直接推得, 未在高刷新率螢幕實測)
- `水量不守恆`: 滑鼠漣漪的 `strength` 固定 `+0.01`, 每次 mousemove 都往池裡加水; 初始 20 顆雨滴才有正負交替. 長時間拖曳後水位上升, 對應 issue #10 「Water rises」與 #19
- `行動裝置永遠不平靜`: issue #11 回報 iOS / Android 上水面無法回到平滑; 程式碼在 float RTT 不可用時降級成 `HALF_FLOAT_OES`, `推論`: fp16 精度不足讓 `0.995` 阻尼與小振幅抵消失效
- `Caustics 不累加`: 整個 `lightgl.js` 沒有啟用 `BLEND`, caustic pass 中重疊的折射三角形互相覆蓋而不是相加 (`推論`: 多個波峰聚焦到同一點時亮度會被低估)
- `場景寫死`: 泳池固定 `[-1, 1]`, 只有一顆球, 解析 ray tracing 無法擴充到任意 mesh; 池壁為 `CLAMP_TO_EDGE` 反射邊界, 無吸收邊界 (issue #3 要求), 無法變成碗形 (issue #8)
- `已停止維護`: 2016 年後無 commit, 15 個 open issue 與 2 個 PR 皆無回應; issue #15 回報 demo 站 HTTPS 憑證錯誤 (2026-11-20 以 `curl -I https://madebyevan.com/webgl-water/` 已回 200)
- `工程結構老舊`: 全域變數, shader 以 `\` 續行字串內嵌, 依賴 `lightgl.js` 模擬 fixed-function matrix stack (`gl.matrixMode`, `gl.rotate`); 無 LICENSE 檔, 無測試; `main.js` 的 `frame` 變數算了沒用, cubemap 的 `yneg` 直接重用 `ypos` 圖
- `浪費的固定成本`: caustic 與 normal 每幀全量重算, 水面靜止時也一樣; 沒有 dirty flag

## 可借鏡 (Takeaways)

- 要做小範圍互動水面 (泳池, 水盆, UI 特效): 直接沿用「heightfield 4 pass + 解析 ray tracing」架構, 但把 step 數改成依 `dt` 累積的固定步長迴圈, 並讓滑鼠漣漪正負成對以守恆水量
- Caustics 面積比技巧可獨立抽出, 移植時記得開 additive blending (`gl.blendFunc(gl.ONE, gl.ONE)`) 讓聚焦正確累加
- 與 [waves.md](waves.md) 互補: `waves` 是頻譜 + FFT 的大範圍海洋 (不可互動, 無邊界), `webgl-water` 是波動方程的小範圍水體 (可互動, 有邊界與物體耦合)
