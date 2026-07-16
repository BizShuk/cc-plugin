package topology

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFixture 建立一個 3 entities / 2 zones 的合法小型拓撲，
// 內容與 plugins/general/skills/topology-builder/references/ 樣本同構。
func writeFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"payments/service-a.md": `---
name: service-a
type: service
zone: payments
aliases: [svc-a]
---

# Service A

## Checkout

kind: method

References:

- calls [[service-b#Validate]] — 下單前驗證
- writes-to [[billing-db#Orders]]

## Billing Cycle

kind: concept

References:

- uses [[#Checkout]]

## Backlinks

<!-- auto-generated: do not hand-edit -->

- calls ← [[service-b#Health Probe]]
`,
		"core/service-b.md": `---
name: service-b
type: service
zone: core
---

# Service B

## Health Probe

kind: method

References:

- calls [[service-a#Checkout]]

## Validate

kind: method

References:

- reads-from [[billing-db#Orders]]

## Backlinks

<!-- auto-generated: do not hand-edit -->

- calls ← [[service-a#Checkout]]
`,
		"payments/billing-db.md": `---
name: billing-db
type: datastore
zone: payments
---

# Billing DB

## Orders

kind: concept

## Backlinks

<!-- auto-generated: do not hand-edit -->

- writes-to ← [[service-a#Checkout]]
- reads-from ← [[service-b#Validate]]
`,
	}
	for rel, content := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestLoadTopology(t *testing.T) {
	topo, err := LoadTopology(writeFixture(t))
	if err != nil {
		t.Fatalf("LoadTopology: %v", err)
	}
	if got := len(topo.Entities); got != 3 {
		t.Fatalf("entities = %d, want 3", got)
	}
	a := topo.Entities["service-a"]
	if a == nil || a.Zone != "payments" || a.Type != "service" {
		t.Fatalf("service-a frontmatter = %+v", a)
	}
	if got := len(a.Dimensions); got != 2 {
		t.Fatalf("service-a dimensions = %d, want 2 (Backlinks 不是維度)", got)
	}
	checkout := a.Dimensions[0]
	if checkout.Name != "Checkout" || checkout.Kind != "method" {
		t.Fatalf("dim[0] = %+v", checkout)
	}
	if got := len(checkout.Edges); got != 2 {
		t.Fatalf("checkout edges = %d, want 2", got)
	}
	e0 := checkout.Edges[0]
	if e0.Relation != "calls" || e0.ToEntity != "service-b" ||
		e0.ToDim != "Validate" || e0.Note != "下單前驗證" {
		t.Fatalf("edge[0] = %+v", e0)
	}
	// 同檔引用 [[#Checkout]] 的 ToEntity 解析為自身
	self := a.Dimensions[1].Edges[0]
	if self.ToEntity != "service-a" || self.ToDim != "Checkout" {
		t.Fatalf("self edge = %+v", self)
	}
	if got := len(a.Backlinks); got != 1 || a.Backlinks[0].FromEntity != "service-b" {
		t.Fatalf("backlinks = %+v", a.Backlinks)
	}
}
