---
name: system-law-life
description: >
  Use when designing health checks, auto-restart, self-healing, graceful
  shutdown, or system resilience mechanisms. Triggers on: "health check",
  "auto-restart", "self-healing", "graceful shutdown", "watchdog",
  "liveness probe", "resilience", "fault tolerance", "creator role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 系統意志
  law: 生命 (Life)
  role: 造物主 (Creator)
---

# law-life — 生命法則 × 造物主

> 法則領域：系統意志 (System Will)

## 法則定義

`生命法則 (Law of Life)` ✖ `造物主 (Creator)`

建立保活機制與系統韌性。

## 系統對應

造物主賦予系統生命力，使其能自我維持：

- `健康檢查 (Health Check / Liveness Probe)`：持續確認服務是否存活
- `自動重啟 (Auto-Restart)`：崩潰後自動恢復，不需人工介入
- `優雅關閉 (Graceful Shutdown)`：讓服務在停止前完成進行中的工作
- `斷路器 (Circuit Breaker)`：偵測故障並隔離，防止級聯崩潰

## 架構問題

當生命法則被忽視時，系統是「脆性」的：

- 單次崩潰需要人工介入才能恢復
- 關閉時丟失進行中的請求
- 依賴的下游故障會拖垮整個系統

## 設計原則

1. 設計系統時預設它**會**失敗，然後設計恢復機制
2. 保活機制本身不能成為單點故障
3. 區分 liveness（存活）與 readiness（就緒）：存活但未就緒的服務不應接受流量
