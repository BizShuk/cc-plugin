---
name: service-a
type: service
zone: payments
tags: [billing, checkout]
aliases: [svc-a]
sources:
    - type: repo
      ref: github.com/example/service-a
---

# Service A

面向顧客的結帳與發票服務。

## Checkout

kind: method

接受購物車、建立訂單並觸發收款。

References:

- calls [[service-b#Validate]] — 下單前驗證
- writes-to [[billing-db#Orders]]

## Invoice Issuing

kind: method

收款完成後開立發票並寄送。

References:

- depends-on [[service-b#Webhook Dispatch]] — 等待收款事件
- writes-to [[billing-db#Invoices]]
- uses [[#Billing Cycle]]

## Billing Cycle

kind: concept

月結週期規則：每月 1 日結算、寬限 7 天。

References:

- uses [[service-b#Validate]] — 結算前重新驗證未付訂單

## External Sources

- [計費規則文件](https://example.com/billing-rules)

## Backlinks

<!-- auto-generated: do not hand-edit -->

- calls ← [[service-b#Health Probe]]
