---
name: fusion-methods
description: >
  Use when hitting a conceptual singularity where knowledge compresses to
  a point and needs expansion, when merging conflicting architectural
  concerns, or when two incompatible abstractions must coexist. Triggers
  on: "design conflict", "how to combine", "orthogonal concerns",
  "first principles", "adapter pattern", "lifecycle", "trade-off".
version: "1.0.0"
metadata:
  type: pattern
  tier: philosophy
---

# fusion-methods

當知識與概念壓縮至極限的`奇異點 (Singularity)` 狀態時，透過四種結構化交織法促成突破性擴展。

## 四大融合思維工具

### 1. 正交性 (Orthogonality)

讓重疊的概念垂直交叉，互不干涉，撐起全新維度。

- 系統應用：分離 data plane 與 control plane
- 辨識信號：兩個關注點總是一起改變 → 可能是耦合而非正交
- 驗證方法：改 A 是否迫使改 B？若否 → 正交

### 2. 同構性 (Isomorphism / 第一性原理)

剝除術語表象，找出不同概念底層相同的物理或數學形狀。

- 系統應用：`快取`、`快照`、`備份` 本質上都是「凍結狀態的副本」
- 辨識信號：不同團隊用不同術語描述同一件事
- 驗證方法：抽掉術語後，運算結構是否一致？

### 3. 引入催化劑 (Catalyst / 介面層)

在互斥的概念之間插入輕量級協議，使其平滑互動。

- 系統應用：在`安全性`與`可觀測性`之間加入 Auth 層
- 辨識信號：兩個元件直接互動時產生矛盾
- 驗證方法：移除催化劑後，兩端是否立即衝突？

### 4. 辯證循環 (Dialectical Cycle)

將對立的概念拉長至時間軸，形成生生不息的生命週期。

- 系統應用：Build → Deploy → Monitor → Destroy → Rebuild
- 辨識信號：A 與 B 看似互斥，但在不同時間點各自必要
- 驗證方法：是否能畫出 A → B → A 的循環？

## 速查決策

| 衝突類型           | 首選融合法   | 範例                           |
| :----------------- | :----------- | :----------------------------- |
| 同時存在互不影響   | 正交性       | 水平擴展 vs 垂直擴展           |
| 表面不同本質相同   | 同構性       | Cache vs Snapshot vs Replica   |
| 直接對立無法共存   | 催化劑       | 安全 vs 可觀測 → Auth 介面     |
| 時序上交替出現     | 辯證循環     | 建設 vs 破壞 → CI/CD 生命週期  |
