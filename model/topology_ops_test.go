package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyCleanFixture(t *testing.T) {
	topo, err := LoadTopology(writeFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	if findings := topo.Verify(); len(findings) != 0 {
		t.Fatalf("want clean, got findings: %v", findings)
	}
}

func TestVerifyFindings(t *testing.T) {
	root := writeFixture(t)

	// rogue.md — 注入：壞 kind（Loose Dim 無 kind）、斷鏈、捏造 backlink、未知關係、缺失 heading
	bad := `---
name: rogue
type: module
zone: core
---

# Rogue

## Loose Dim

References:

- calls [[ghost#Nowhere]]

## Bad Edges

kind: concept

References:

- summons [[service-a#Checkout]]
- calls [[service-a#No Such Section]]

## Backlinks

<!-- auto-generated: do not hand-edit -->

- uses ← [[service-a#Checkout]]
`
	if err := os.WriteFile(filepath.Join(root, "core", "rogue.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}

	// mismatch.md — frontmatter name != filename (name: other-name, file: mismatch.md)
	mismatch := `---
name: other-name
type: service
zone: core
---

# Mismatch

## Info

kind: concept
`
	if err := os.WriteFile(filepath.Join(root, "core", "mismatch.md"), []byte(mismatch), 0o644); err != nil {
		t.Fatal(err)
	}

	// wrongzone.md — frontmatter zone != folder (zone: payments, file lives in core/)
	wrongzone := `---
name: wrongzone
type: service
zone: payments
---

# Wrong Zone

## Info

kind: concept
`
	if err := os.WriteFile(filepath.Join(root, "core", "wrongzone.md"), []byte(wrongzone), 0o644); err != nil {
		t.Fatal(err)
	}

	topo, err := LoadTopology(root)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(topo.Verify(), "\n")
	for _, want := range []string{
		// original three
		"rogue#Loose Dim: missing or invalid kind",
		"missing entity: [[ghost]]",
		"backlink without forward edge",
		// four new rules
		`mismatch: frontmatter name "other-name" != filename`,
		`wrongzone: frontmatter zone "payments" != folder "core"`,
		`unknown relation "summons"`,
		`missing heading: [[service-a#No Such Section]]`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("findings missing %q in:\n%s", want, got)
		}
	}
}

func TestUnlinked(t *testing.T) {
	topo, err := LoadTopology(writeFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	noIn, noOut := topo.Unlinked()
	if len(noIn) != 0 {
		t.Errorf("noInbound = %v, want empty (三個 entity 都有入邊)", noIn)
	}
	if len(noOut) != 1 || noOut[0] != "billing-db" {
		t.Errorf("noOutbound = %v, want [billing-db]", noOut)
	}
}
