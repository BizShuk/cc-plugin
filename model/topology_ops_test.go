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
	// 注入三類錯誤：壞 kind、斷鏈、捏造 backlink
	bad := `---
name: rogue
type: module
zone: core
---

# Rogue

## Loose Dim

References:

- calls [[ghost#Nowhere]]

## Backlinks

<!-- auto-generated: do not hand-edit -->

- uses ← [[service-a#Checkout]]
`
	if err := os.WriteFile(filepath.Join(root, "core", "rogue.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	topo, err := LoadTopology(root)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(topo.Verify(), "\n")
	for _, want := range []string{
		"rogue#Loose Dim: missing or invalid kind",
		"missing entity: [[ghost]]",
		"backlink without forward edge",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("findings missing %q in:\n%s", want, got)
		}
	}
}
