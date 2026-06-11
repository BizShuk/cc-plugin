---
name: service-b
type: service
zone: core
tags: [validation, events]
aliases: [svc-b, validator]
sources:
    - type: repo
      ref: github.com/example/service-b
---

# Service B

核心驗證與事件分發服務。

## Health Probe

kind: method

定期以合成交易探測下游服務可用性。

References:

- calls [[service-a#Checkout]] — 合成交易探測

## Validate

kind: method

驗證訂單資料與付款狀態；結果供結帳與結算流程使用。

References:

- reads-from [[billing-db#Orders]]

## Webhook Dispatch

kind: interface

對外發布收款與退款事件的 webhook 介面；訂閱者自行註冊回呼網址。

## Backlinks

<!-- auto-generated: do not hand-edit -->

- calls ← [[service-a#Checkout]]
- depends-on ← [[service-a#Invoice Issuing]]
- uses ← [[service-a#Billing Cycle]]
