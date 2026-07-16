package topology

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"core/service-a.md": `---
name: service-a
type: service
zone: core
---

# Service A

## Call

kind: method

- calls [[service-b#Handle]]

## Backlinks

<!-- auto-generated: do not hand-edit -->

- uses ← [[service-b#Handle]]
`,
		"core/service-b.md": `---
name: service-b
type: service
zone: core
---

# Service B

## Handle

kind: method

- uses [[service-a#Call]]

## Backlinks

<!-- auto-generated: do not hand-edit -->

- calls ← [[service-a#Call]]
`,
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := TopologyCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs(args)
	err := cmd.Execute()
	return output.String(), err
}

func TestVerifyCommand(t *testing.T) {
	output, err := run(t, "verify", "--root", writeFixture(t))
	if err != nil {
		t.Fatalf("verify: %v\n%s", err, output)
	}
	if output != "OK\n" {
		t.Fatalf("output = %q, want OK", output)
	}
}

func TestQueryCommand(t *testing.T) {
	output, err := run(t, "query", "service-a", "--root", writeFixture(t))
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if !strings.Contains(output, "service-a#Call -calls-> service-b#Handle") {
		t.Fatalf("missing outbound edge:\n%s", output)
	}
}

func TestRewriteCommand(t *testing.T) {
	root := writeFixture(t)
	path := filepath.Join(root, "core", "service-b.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.ReplaceAll(string(raw), "- calls ← [[service-a#Call]]\n", "")
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := run(t, "rewrite", "--root", root); err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	rewritten, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rewritten), "- calls ← [[service-a#Call]]") {
		t.Fatalf("backlink was not restored:\n%s", rewritten)
	}
	if _, err := os.Stat(filepath.Join(root, "_index.md")); err != nil {
		t.Fatalf("index not written: %v", err)
	}
}
