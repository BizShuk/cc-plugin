---
name: billing-db
type: datastore
zone: payments
tags: [postgres]
aliases: [billing-database]
sources:
    - type: database
      ref: postgres://billing
---

# Billing DB

訂單與發票的主要儲存；典型的「無出邊」實體（datastore 只被讀寫）。

## Orders

kind: concept

訂單表群：`orders`、`order_items`。

## Invoices

kind: concept

發票表群：`invoices`、`credit_notes`。

## Retention Policy

kind: concept

已結案資料保留 24 個月後歸檔。

## Backlinks

<!-- auto-generated: do not hand-edit -->

- writes-to ← [[service-a#Checkout]]
- writes-to ← [[service-a#Invoice Issuing]]
- reads-from ← [[service-b#Validate]]
