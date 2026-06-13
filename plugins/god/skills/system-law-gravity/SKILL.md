---
name: system-law-gravity
description: >
  Use when designing load balancing, traffic routing, rate limiting, or
  distributing workload across multiple instances. Triggers on:
  "load balancing", "traffic distribution", "rate limiting", "auto-scaling",
  "request routing", "capacity planning", "orchestrator role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 宇宙基石
  law: 重力 (Gravity)
  role: 指揮家 (Orchestrator)
---

# law-gravity — 重力法則 × 指揮家

> 法則領域：宇宙基石 (Cosmic Foundation)

## 法則定義

`重力法則 (Law of Gravity)` ✖ `指揮家 (Orchestrator)`

平衡流量引力，執行負載平衡。

## 系統對應

重力在系統中是流量的自然拉力，指揮家負責對抗失衡：

- `負載平衡 (Load Balancing)`：讓請求均勻分散到健康節點
- `流量整形 (Traffic Shaping)`：限制突發流量，保護下游
- `自動擴縮 (Auto-Scaling)`：依流量重力動態調整資源
- `熔斷 (Circuit Breaking)`：當某節點過載時自動切斷

## 架構問題

當重力法則被違反時，會出現：

- 熱點節點崩潰，其他節點閒置（流量引力失衡）
- 雪崩效應（過載向下傳播）
- 容量不足或浪費（靜態分配無法適應動態流量）

## 設計原則

1. 任何單點都應能被替換或擴展，不能成為必要的「重力中心」
2. 在系統邊界（入口）就做流量整形，不讓壓力傳入
3. 設計降級策略：過載時優雅降級，不是崩潰
