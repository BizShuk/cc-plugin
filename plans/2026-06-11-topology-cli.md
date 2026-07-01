# Topology CLI Implementation Plan

> For agentic workers: REQUIRED SUB-SKILL — use superpowers:subagent-driven-development
> (recommended) or superpowers:executing-plans to implement this plan task-by-task.
> Steps use checkbox (`- [ ]`) syntax for tracking.

`Goal:` 在本 repo 新增 `cc-plugin topology` Cobra 子命令群，對 `topology-builder`
技能產出的知識圖譜（`<root>/<zone>/<entity>.md`，Obsidian wikilink 邊）做機械化操作：
完整性驗證、邊查詢、Unlinked 報告、Backlinks 重算、`_index.md` 重生成。

`Architecture:` 解析與圖運算放在 `model/`（領域層，純函數、可單測），命令薄殼放在
`cmd/topology/`（仿 `cmd/export/` 子包模式），由 `cmd/root.go` 註冊。檔案格式即
資料庫：所有運算讀寫 markdown 檔，無額外狀態儲存。

`Tech Stack:` Go 1.26 / cobra / `gopkg.in/yaml.v3`（frontmatter）/ 標準庫 regexp。
不引入 markdown parser——技能格式行導向，行掃描 + regex 即足夠（YAGNI）。

`Format contract`（來源：`plugins/general/skills/topology-builder/SKILL.md`）：

- entity 檔：YAML frontmatter（`name/type/zone/tags/aliases/sources`）+ 二級標題維度章節
- 維度標題下第一個非空行：`kind: concept|method|state|interface`
- 正向邊（僅在維度章節）：`- <relation> [[entity#Section]] — note`（note 可省略，
  分隔符是「空格 + em-dash + 空格」）；同檔引用 `[[#Section]]`
- `## External Sources` 與 `## Backlinks` 是固定章節、不是維度
- backlink 行：`- <relation> ← [[source#Section]]`，只能由全圖正向邊推導
- 參考樣本：`plugins/general/skills/topology-builder/references/`

---

## File Structure

| 檔案 | 動作 | 職責 |
| :--- | :--- | :--- |
| `model/topology.go` | Create | 型別（`TopoEntity/TopoDimension/TopoEdge/Topology`）、frontmatter 與邊解析、`LoadTopology` |
| `model/topology_ops.go` | Create | 圖運算：`Verify`、`Unlinked`、`BacklinksFor`、`RenderBacklinksSection`、`RenderIndex` |
| `model/topology_test.go` | Create | fixture 建構 + 解析/載入測試 |
| `model/topology_ops_test.go` | Create | 圖運算測試 |
| `cmd/topology/topology.go` | Create | 父命令 `TopologyCmd()`、`--root` persistent flag、共用 `loadFromFlags` |
| `cmd/topology/verify.go` | Create | `verify` 子命令 |
| `cmd/topology/unlinked.go` | Create | `unlinked` 子命令 |
| `cmd/topology/query.go` | Create | `query <entity>` 子命令（`--in/--out/--depth`） |
| `cmd/topology/rewrite.go` | Create | `backlinks` 與 `index` 子命令（`--write`） |
| `cmd/topology/topology_test.go` | Create | 命令層整合測試 |
| `cmd/root.go` | Modify (`init()` 區塊, 約 L23-34) | 註冊 `topology.TopologyCmd()` |
| `go.mod` | Modify | `gopkg.in/yaml.v3` 升為直接依賴 |
| `CLAUDE.md` | Modify | 結構樹、模組對應、測試清單 |
| `plugins/general/skills/topology-builder/SKILL.md` | Modify | Verification 章節補 CLI 用法 |

---

### Task 1: 依賴與骨架

`Files:`

- Modify: `go.mod`

- [ ] Step 1: 確認 yaml.v3 目前為間接依賴

Run: `grep yaml ../go.mod`
Expected: `gopkg.in/yaml.v3 ... // indirect`

- [ ] Step 2: 升為直接依賴

Run: `cd /Users/shuk/projects/cc-plugin && go get gopkg.in/yaml.v3 && go mod tidy`
Expected: exit 0；`go.mod` require 區塊出現 `gopkg.in/yaml.v3`（無 indirect 註記，
在 model 程式碼 import 後 `go mod tidy` 才會固定，Task 2 結束時再跑一次即可）

- [ ] Step 3: Commit

```bash
git add go.mod go.sum
git commit -m "chore: promote gopkg.in/yaml.v3 to direct dependency"
```

### Task 2: model/topology.go — 型別與解析

`Files:`

- Create: `model/topology.go`
- Test: `model/topology_test.go`

- [ ] Step 1: 寫 fixture helper 與失敗測試

`model/topology_test.go`：

```go
package model

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
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./model/ -run TestLoadTopology -v`
Expected: FAIL — `undefined: LoadTopology`

- [ ] Step 3: 實作 `model/topology.go`

```go
package model

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// TopoKinds 是維度 kind 標註的合法取值（見 topology-builder skill）。
var TopoKinds = map[string]bool{
	"concept": true, "method": true, "state": true, "interface": true,
}

// TopoRelations 是正向邊的合法關係動詞。
var TopoRelations = map[string]bool{
	"calls": true, "uses": true, "reads-from": true, "writes-to": true,
	"publishes-to": true, "subscribes-to": true, "depends-on": true,
	"mentions": true, "owned-by": true,
}

// TopoEdge 是一條有向邊：FromEntity#FromDim -Relation-> ToEntity#ToDim。
type TopoEdge struct {
	FromEntity string
	FromDim    string
	Relation   string
	ToEntity   string
	ToDim      string
	Note       string
}

// TopoDimension 是 entity 檔內一個 `## ` 維度章節。
type TopoDimension struct {
	Name  string
	Kind  string
	Edges []TopoEdge
}

// TopoEntity 是一個 entity 檔（frontmatter + 維度 + 既有 backlink 行）。
type TopoEntity struct {
	Name       string          `yaml:"name"`
	Type       string          `yaml:"type"`
	Zone       string          `yaml:"zone"`
	Aliases    []string        `yaml:"aliases"`
	Path       string          `yaml:"-"`
	Dimensions []TopoDimension `yaml:"-"`
	Backlinks  []TopoEdge      `yaml:"-"`
}

// Topology 是整個圖：以檔名（不含 .md）為鍵。
type Topology struct {
	Root     string
	Entities map[string]*TopoEntity
	Findings []string // 載入期問題：跨 zone 重名等
}

var (
	// `- calls [[service-b#Validate]] — note`，note 分隔符為「空格 em-dash 空格」
	topoEdgeRe     = regexp.MustCompile(`^- ([a-z-]+) \[\[([^\]#]*)(?:#([^\]]+))?\]\](?: — (.*))?$`)
	topoBacklinkRe = regexp.MustCompile(`^- ([a-z-]+) ← \[\[([^\]#]*)(?:#([^\]]+))?\]\]$`)
)

// LoadTopology 讀取 <root>/<zone>/*.md（略過 _index.md）建圖。
func LoadTopology(root string) (*Topology, error) {
	topo := &Topology{Root: root, Entities: map[string]*TopoEntity{}}
	zones, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read topology root: %w", err)
	}
	for _, z := range zones {
		if !z.IsDir() {
			continue
		}
		files, err := filepath.Glob(filepath.Join(root, z.Name(), "*.md"))
		if err != nil {
			return nil, fmt.Errorf("glob zone %s: %w", z.Name(), err)
		}
		for _, f := range files {
			name := strings.TrimSuffix(filepath.Base(f), ".md")
			if name == "_index" {
				continue
			}
			if prev, ok := topo.Entities[name]; ok {
				topo.Findings = append(topo.Findings,
					fmt.Sprintf("duplicate filename: %s (%s, %s)", name, prev.Path, f))
				continue
			}
			e, err := parseEntityFile(f)
			if err != nil {
				return nil, fmt.Errorf("parse %s: %w", f, err)
			}
			topo.Entities[name] = e
		}
	}
	return topo, nil
}

// Names 回傳排序後的 entity 檔名，供確定性迭代。
func (t *Topology) Names() []string {
	names := make([]string, 0, len(t.Entities))
	for n := range t.Entities {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Edges 回傳全圖正向邊（不含 Backlinks 區段）。
func (t *Topology) Edges() []TopoEdge {
	var out []TopoEdge
	for _, n := range t.Names() {
		for _, d := range t.Entities[n].Dimensions {
			out = append(out, d.Edges...)
		}
	}
	return out
}

func parseEntityFile(path string) (*TopoEntity, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	body := string(raw)
	e := &TopoEntity{Path: path}
	if strings.HasPrefix(body, "---\n") {
		parts := strings.SplitN(body[4:], "\n---\n", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unterminated frontmatter")
		}
		if err := yaml.Unmarshal([]byte(parts[0]), e); err != nil {
			return nil, fmt.Errorf("frontmatter: %w", err)
		}
		body = parts[1]
	}
	name := strings.TrimSuffix(filepath.Base(path), ".md")
	var dim *TopoDimension
	inBacklinks := false
	sawFirstLine := false
	for _, line := range strings.Split(body, "\n") {
		switch {
		case strings.HasPrefix(line, "## "):
			h := strings.TrimPrefix(line, "## ")
			inBacklinks = h == "Backlinks"
			if h == "Backlinks" || h == "External Sources" {
				dim = nil
				continue
			}
			e.Dimensions = append(e.Dimensions, TopoDimension{Name: h})
			dim = &e.Dimensions[len(e.Dimensions)-1]
			sawFirstLine = false
		case inBacklinks:
			if m := topoBacklinkRe.FindStringSubmatch(line); m != nil {
				e.Backlinks = append(e.Backlinks, TopoEdge{
					Relation: m[1], FromEntity: m[2], FromDim: m[3], ToEntity: name,
				})
			}
		case dim != nil:
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !sawFirstLine {
				sawFirstLine = true
				if strings.HasPrefix(trimmed, "kind: ") {
					dim.Kind = strings.TrimPrefix(trimmed, "kind: ")
					continue
				}
			}
			if m := topoEdgeRe.FindStringSubmatch(line); m != nil {
				to := m[2]
				if to == "" {
					to = name // 同檔引用 [[#Section]]
				}
				dim.Edges = append(dim.Edges, TopoEdge{
					FromEntity: name, FromDim: dim.Name,
					Relation: m[1], ToEntity: to, ToDim: m[3], Note: m[4],
				})
			}
		}
	}
	return e, nil
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./model/ -run TestLoadTopology -v && go mod tidy`
Expected: PASS

- [ ] Step 5: Commit

```bash
git add model/topology.go model/topology_test.go go.mod go.sum
git commit -m "feat(topology): add graph model and entity file parser"
```

### Task 3: Verify — 完整性檢查

`Files:`

- Create: `model/topology_ops.go`
- Test: `model/topology_ops_test.go`

- [ ] Step 1: 寫失敗測試

`model/topology_ops_test.go`：

```go
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
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./model/ -run TestVerify -v`
Expected: FAIL — `undefined: (*Topology).Verify`

- [ ] Step 3: 實作 `model/topology_ops.go` 的 Verify

```go
package model

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Verify 回傳所有機械可查的規則違反；空切片代表通過。
func (t *Topology) Verify() []string {
	findings := append([]string{}, t.Findings...)
	for _, n := range t.Names() {
		e := t.Entities[n]
		if e.Name != n {
			findings = append(findings,
				fmt.Sprintf("%s: frontmatter name %q != filename", n, e.Name))
		}
		zone := filepath.Base(filepath.Dir(e.Path))
		if e.Zone != zone {
			findings = append(findings,
				fmt.Sprintf("%s: frontmatter zone %q != folder %q", n, e.Zone, zone))
		}
		for _, d := range e.Dimensions {
			if !TopoKinds[d.Kind] {
				findings = append(findings,
					fmt.Sprintf("%s#%s: missing or invalid kind %q", n, d.Name, d.Kind))
			}
		}
	}
	for _, ed := range t.Edges() {
		if !TopoRelations[ed.Relation] {
			findings = append(findings,
				fmt.Sprintf("%s#%s: unknown relation %q", ed.FromEntity, ed.FromDim, ed.Relation))
		}
		target, ok := t.Entities[ed.ToEntity]
		if !ok {
			findings = append(findings,
				fmt.Sprintf("%s#%s: missing entity: [[%s]]", ed.FromEntity, ed.FromDim, ed.ToEntity))
			continue
		}
		if ed.ToDim != "" && !hasDim(target, ed.ToDim) {
			findings = append(findings,
				fmt.Sprintf("%s#%s: missing heading: [[%s#%s]]",
					ed.FromEntity, ed.FromDim, ed.ToEntity, ed.ToDim))
		}
	}
	forward := map[string]bool{}
	for _, ed := range t.Edges() {
		forward[edgeKey(ed)] = true
	}
	for _, n := range t.Names() {
		for _, b := range t.Entities[n].Backlinks {
			if !forward[edgeKey(b)] {
				findings = append(findings,
					fmt.Sprintf("%s: backlink without forward edge: %s ← [[%s#%s]]",
						n, b.Relation, b.FromEntity, b.FromDim))
			}
		}
	}
	return findings
}

func edgeKey(e TopoEdge) string {
	return e.Relation + "|" + e.FromEntity + "|" + e.FromDim + "|" + e.ToEntity
}

func hasDim(e *TopoEntity, name string) bool {
	for _, d := range e.Dimensions {
		if d.Name == name {
			return true
		}
	}
	return false
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./model/ -run TestVerify -v`
Expected: PASS（兩個測試都過）

- [ ] Step 5: Commit

```bash
git add model/topology_ops.go model/topology_ops_test.go
git commit -m "feat(topology): add integrity verification"
```

### Task 4: Unlinked 計算

`Files:`

- Modify: `model/topology_ops.go`
- Test: `model/topology_ops_test.go`

- [ ] Step 1: 寫失敗測試（附加到 `model/topology_ops_test.go`）

```go
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
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./model/ -run TestUnlinked -v`
Expected: FAIL — `undefined: (*Topology).Unlinked`

- [ ] Step 3: 實作（附加到 `model/topology_ops.go`）

```go
// Unlinked 回傳無入邊與無出邊的 entity 清單。
// 只計跨實體正向邊：同檔引用與 Backlinks 區段一律不計。
func (t *Topology) Unlinked() (noInbound, noOutbound []string) {
	in, out := map[string]bool{}, map[string]bool{}
	for _, ed := range t.Edges() {
		if ed.ToEntity == ed.FromEntity {
			continue
		}
		if _, ok := t.Entities[ed.ToEntity]; !ok {
			continue
		}
		out[ed.FromEntity] = true
		in[ed.ToEntity] = true
	}
	for _, n := range t.Names() {
		if !in[n] {
			noInbound = append(noInbound, n)
		}
		if !out[n] {
			noOutbound = append(noOutbound, n)
		}
	}
	return noInbound, noOutbound
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./model/ -run TestUnlinked -v`
Expected: PASS

- [ ] Step 5: Commit

```bash
git add model/topology_ops.go model/topology_ops_test.go
git commit -m "feat(topology): add unlinked entity report"
```

### Task 5: Backlinks 重算與區段改寫

`Files:`

- Modify: `model/topology_ops.go`
- Test: `model/topology_ops_test.go`

- [ ] Step 1: 寫失敗測試（附加）

```go
func TestBacklinksFor(t *testing.T) {
	topo, err := LoadTopology(writeFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	got := topo.BacklinksFor("billing-db")
	want := []string{
		"- reads-from ← [[service-b#Validate]]",
		"- writes-to ← [[service-a#Checkout]]",
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRenderBacklinksSection(t *testing.T) {
	content := "# X\n\n## Dim\n\nkind: method\n\n## Backlinks\n\n" +
		"<!-- auto-generated: do not hand-edit -->\n\n- stale ← [[old#Gone]]\n"
	out := RenderBacklinksSection(content, []string{"- calls ← [[a#B]]"})
	if !strings.Contains(out, "- calls ← [[a#B]]") {
		t.Errorf("missing new backlink:\n%s", out)
	}
	if strings.Contains(out, "stale") {
		t.Errorf("stale entry survived:\n%s", out)
	}
	if !strings.Contains(out, "## Dim") {
		t.Errorf("earlier section lost:\n%s", out)
	}
	// 冪等：重跑結果不變
	if again := RenderBacklinksSection(out, []string{"- calls ← [[a#B]]"}); again != out {
		t.Errorf("not idempotent")
	}
}
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./model/ -run 'TestBacklinksFor|TestRenderBacklinksSection' -v`
Expected: FAIL — undefined

- [ ] Step 3: 實作（附加到 `model/topology_ops.go`）

```go
const topoBacklinkMarker = "<!-- auto-generated: do not hand-edit -->"

// BacklinksFor 由全圖正向邊重算 name 的 backlink 行（排序去重）。
func (t *Topology) BacklinksFor(name string) []string {
	seen := map[string]bool{}
	var lines []string
	for _, ed := range t.Edges() {
		if ed.ToEntity != name || ed.FromEntity == name {
			continue
		}
		l := fmt.Sprintf("- %s ← [[%s#%s]]", ed.Relation, ed.FromEntity, ed.FromDim)
		if !seen[l] {
			seen[l] = true
			lines = append(lines, l)
		}
	}
	sort.Strings(lines)
	return lines
}

// RenderBacklinksSection 以 lines 重寫 content 的 ## Backlinks 區段；
// 區段不存在時附加到檔尾。
func RenderBacklinksSection(content string, lines []string) string {
	section := "## Backlinks\n\n" + topoBacklinkMarker + "\n"
	if len(lines) > 0 {
		section += "\n" + strings.Join(lines, "\n") + "\n"
	}
	idx := strings.Index(content, "## Backlinks")
	if idx == -1 {
		return strings.TrimRight(content, "\n") + "\n\n" + section
	}
	rest := content[idx:]
	if next := strings.Index(rest[1:], "\n## "); next != -1 {
		return content[:idx] + section + rest[next+1:]
	}
	return content[:idx] + section
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./model/ -run 'TestBacklinksFor|TestRenderBacklinksSection' -v`
Expected: PASS

- [ ] Step 5: Commit

```bash
git add model/topology_ops.go model/topology_ops_test.go
git commit -m "feat(topology): recompute backlinks from forward edges"
```

### Task 6: `_index.md` 生成

`Files:`

- Modify: `model/topology_ops.go`
- Test: `model/topology_ops_test.go`

- [ ] Step 1: 寫失敗測試（附加）

```go
func TestRenderIndex(t *testing.T) {
	topo, err := LoadTopology(writeFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	existing := "# Old\n\n## Frontier\n\n- `payment-gateway` — 外部收款 API\n"
	out := topo.RenderIndex(existing)
	for _, want := range []string{
		"| [[service-a]] | payments | service | 2 |",
		"subgraph core",
		"payment-gateway", // Frontier 保留
		"## Unlinked",
		"- [[billing-db]]",
		"- 無 (None)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("index missing %q in:\n%s", want, out)
		}
	}
	if !strings.HasSuffix(strings.TrimRight(out, "\n"), "- [[billing-db]]") {
		t.Errorf("Unlinked 必須是檔尾章節:\n%s", out)
	}
}
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./model/ -run TestRenderIndex -v`
Expected: FAIL — undefined

- [ ] Step 3: 實作（附加到 `model/topology_ops.go`，需新增 import `regexp`）

```go
// RenderIndex 產生 _index.md 內容；existing 為現有檔案內容（無檔給空字串），
// 其 ## Frontier 區段原文保留。Unlinked 依 skill 規範必為檔尾章節。
func (t *Topology) RenderIndex(existing string) string {
	var b strings.Builder
	b.WriteString("# Topology Index\n\n## Entity Registry\n\n")
	b.WriteString("| Entity | Zone | Type | Dimensions |\n| :--- | :--- | :--- | :--- |\n")
	names := t.Names()
	sort.SliceStable(names, func(i, j int) bool {
		zi, zj := t.Entities[names[i]].Zone, t.Entities[names[j]].Zone
		if zi != zj {
			return zi < zj
		}
		return names[i] < names[j]
	})
	for _, n := range names {
		e := t.Entities[n]
		fmt.Fprintf(&b, "| [[%s]] | %s | %s | %d |\n", n, e.Zone, e.Type, len(e.Dimensions))
	}
	b.WriteString("\n## Overview Diagram\n\n```mermaid\nflowchart LR\n")
	id := map[string]string{}
	zones := map[string][]string{}
	for _, n := range names {
		id[n] = fmt.Sprintf("n%d", len(id))
		zones[t.Entities[n].Zone] = append(zones[t.Entities[n].Zone], n)
	}
	zoneNames := make([]string, 0, len(zones))
	for z := range zones {
		zoneNames = append(zoneNames, z)
	}
	sort.Strings(zoneNames)
	for _, z := range zoneNames {
		fmt.Fprintf(&b, "    subgraph %s\n", z)
		for _, n := range zones[z] {
			fmt.Fprintf(&b, "        %s[%s]\n", id[n], n)
		}
		b.WriteString("    end\n")
	}
	seen := map[string]bool{}
	for _, ed := range t.Edges() {
		if ed.FromEntity == ed.ToEntity {
			continue
		}
		if _, ok := t.Entities[ed.ToEntity]; !ok {
			continue
		}
		key := ed.FromEntity + ">" + ed.ToEntity
		if seen[key] {
			continue
		}
		seen[key] = true
		fmt.Fprintf(&b, "    %s --> %s\n", id[ed.FromEntity], id[ed.ToEntity])
	}
	b.WriteString("```\n\n")
	b.WriteString(frontierSection(existing))
	noIn, noOut := t.Unlinked()
	b.WriteString("\n## Unlinked\n\n無入邊 (no inbound)：\n\n")
	b.WriteString(wikiList(noIn))
	b.WriteString("\n無出邊 (no outbound)：\n\n")
	b.WriteString(wikiList(noOut))
	return b.String()
}

func wikiList(names []string) string {
	if len(names) == 0 {
		return "- 無 (None)\n"
	}
	var b strings.Builder
	for _, n := range names {
		fmt.Fprintf(&b, "- [[%s]]\n", n)
	}
	return b.String()
}

func frontierSection(existing string) string {
	re := regexp.MustCompile(`(?s)## Frontier\n.*?(\n## |\z)`)
	if m := re.FindString(existing); m != "" {
		m = strings.TrimSuffix(m, "\n## ")
		return strings.TrimRight(m, "\n") + "\n"
	}
	return "## Frontier\n\n- 無 (None)\n"
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./model/ -run TestRenderIndex -v`
Expected: PASS

- [ ] Step 5: Commit

```bash
git add model/topology_ops.go model/topology_ops_test.go
git commit -m "feat(topology): render _index.md with registry, mermaid, unlinked"
```

### Task 7: cmd/topology — 父命令、verify、unlinked

`Files:`

- Create: `cmd/topology/topology.go`
- Create: `cmd/topology/verify.go`
- Create: `cmd/topology/unlinked.go`
- Modify: `cmd/root.go:23-34`（`init()`）
- Test: `cmd/topology/topology_test.go`

- [ ] Step 1: 寫失敗測試

`cmd/topology/topology_test.go`：

```go
package topology

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 與 model fixture 同構的最小拓撲（套件邊界外無法共用 _test helper）。
func writeFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"payments/service-a.md": `---
name: service-a
type: service
zone: payments
---

# Service A

## Checkout

kind: method

References:

- calls [[service-b#Validate]]
`,
		"core/service-b.md": `---
name: service-b
type: service
zone: core
---

# Service B

## Validate

kind: method
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

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := TopologyCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestVerifyCommandClean(t *testing.T) {
	out, err := run(t, "verify", "--root", writeFixture(t))
	if err != nil {
		t.Fatalf("verify: %v\n%s", err, out)
	}
	if !strings.Contains(out, "OK") {
		t.Errorf("want OK, got:\n%s", out)
	}
}

func TestUnlinkedCommand(t *testing.T) {
	out, err := run(t, "unlinked", "--root", writeFixture(t))
	if err != nil {
		t.Fatalf("unlinked: %v", err)
	}
	if !strings.Contains(out, "no outbound") || !strings.Contains(out, "service-b") {
		t.Errorf("want service-b in no outbound, got:\n%s", out)
	}
}
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./cmd/topology/ -v`
Expected: FAIL — `undefined: TopologyCmd`

- [ ] Step 3: 實作三個檔案

`cmd/topology/topology.go`：

```go
package topology

import (
	"github.com/bizshuk/cc-plugin/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const defaultRoot = "~/projects/product/topologies"

// TopologyCmd returns the top-level topology Cobra command.
func TopologyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topology",
		Short: "Operate on topology-builder knowledge graphs",
	}
	cmd.PersistentFlags().String("root", defaultRoot, "topology root directory")
	cmd.AddCommand(VerifyCmd())
	cmd.AddCommand(UnlinkedCmd())
	return cmd
}

func loadFromFlags(cmd *cobra.Command) (*model.Topology, error) {
	root, _ := cmd.Flags().GetString("root")
	if expanded, err := homedir.Expand(root); err == nil {
		root = expanded
	}
	return model.LoadTopology(root)
}
```

`cmd/topology/verify.go`：

```go
package topology

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VerifyCmd returns the verify subcommand.
func VerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Check graph integrity: links, kinds, backlinks, duplicates",
		RunE: func(cmd *cobra.Command, args []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("load topology: %w", err)
			}
			findings := topo.Verify()
			for _, f := range findings {
				fmt.Fprintln(cmd.OutOrStdout(), f)
			}
			if len(findings) > 0 {
				return fmt.Errorf("%d finding(s)", len(findings))
			}
			fmt.Fprintln(cmd.OutOrStdout(), "OK")
			return nil
		},
	}
}
```

`cmd/topology/unlinked.go`：

```go
package topology

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UnlinkedCmd returns the unlinked subcommand.
func UnlinkedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlinked",
		Short: "List entities with no inbound or no outbound edges",
		RunE: func(cmd *cobra.Command, args []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("load topology: %w", err)
			}
			noIn, noOut := topo.Unlinked()
			w := cmd.OutOrStdout()
			fmt.Fprintln(w, "no inbound:")
			printNames(w, noIn)
			fmt.Fprintln(w, "no outbound:")
			printNames(w, noOut)
			return nil
		},
	}
}

func printNames(w interface{ Write([]byte) (int, error) }, names []string) {
	if len(names) == 0 {
		fmt.Fprintln(w, "  (none)")
		return
	}
	for _, n := range names {
		fmt.Fprintf(w, "  %s\n", n)
	}
}
```

`cmd/root.go` 修改（import 加 `"github.com/bizshuk/cc-plugin/cmd/topology"`，
`init()` 內 `RootCmd.AddCommand(export.ExportCmd())` 之後加一行）：

```go
	RootCmd.AddCommand(topology.TopologyCmd())
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./cmd/... -v && go build ./...`
Expected: PASS；build 成功

- [ ] Step 5: Commit

```bash
git add cmd/topology/ cmd/root.go
git commit -m "feat(topology): add topology verify and unlinked subcommands"
```

### Task 8: query 子命令

`Files:`

- Create: `cmd/topology/query.go`
- Test: `cmd/topology/topology_test.go`

- [ ] Step 1: 寫失敗測試（附加）

```go
func TestQueryCommand(t *testing.T) {
	out, err := run(t, "query", "service-a", "--root", writeFixture(t))
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	want := "service-a#Checkout -calls-> service-b#Validate"
	if !strings.Contains(out, want) {
		t.Errorf("want %q, got:\n%s", want, out)
	}
}

func TestQueryUnknownEntity(t *testing.T) {
	_, err := run(t, "query", "nobody", "--root", writeFixture(t))
	if err == nil || !strings.Contains(err.Error(), "unknown entity") {
		t.Errorf("want unknown entity error, got %v", err)
	}
}
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./cmd/topology/ -run TestQuery -v`
Expected: FAIL — `unknown command "query"`

- [ ] Step 3: 實作 `cmd/topology/query.go`，並在 `topology.go` 的
`TopologyCmd()` 加 `cmd.AddCommand(QueryCmd())`

```go
package topology

import (
	"fmt"
	"io"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
)

// QueryCmd returns the query subcommand.
func QueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query <entity>",
		Short: "List inbound and outbound edges of an entity (depth 1-2)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("load topology: %w", err)
			}
			name := args[0]
			if _, ok := topo.Entities[name]; !ok {
				return fmt.Errorf("unknown entity: %s", name)
			}
			onlyIn, _ := cmd.Flags().GetBool("in")
			onlyOut, _ := cmd.Flags().GetBool("out")
			depth, _ := cmd.Flags().GetInt("depth")
			if depth < 1 || depth > 2 {
				return fmt.Errorf("depth must be 1 or 2")
			}
			showIn := !onlyOut || onlyIn
			showOut := !onlyIn || onlyOut
			printEdges(cmd.OutOrStdout(), topo, name, showOut, showIn, depth)
			return nil
		},
	}
	cmd.Flags().Bool("in", false, "only inbound edges")
	cmd.Flags().Bool("out", false, "only outbound edges")
	cmd.Flags().Int("depth", 1, "traversal depth (1 or 2)")
	return cmd
}

func printEdges(w io.Writer, topo *model.Topology, name string, showOut, showIn bool, depth int) {
	printed := map[string]bool{}
	visited := map[string]bool{name: true}
	frontier := []string{name}
	for level := 1; level <= depth; level++ {
		var next []string
		for _, n := range frontier {
			for _, ed := range topo.Edges() {
				var neighbor string
				switch {
				case showOut && ed.FromEntity == n && ed.ToEntity != n:
					neighbor = ed.ToEntity
				case showIn && ed.ToEntity == n && ed.FromEntity != n:
					neighbor = ed.FromEntity
				default:
					continue
				}
				line := fmt.Sprintf("%s#%s -%s-> %s#%s",
					ed.FromEntity, ed.FromDim, ed.Relation, ed.ToEntity, ed.ToDim)
				if !printed[line] {
					printed[line] = true
					fmt.Fprintln(w, line)
				}
				next = append(next, neighbor)
			}
		}
		frontier = nil
		for _, n := range next {
			if !visited[n] {
				visited[n] = true
				frontier = append(frontier, n)
			}
		}
	}
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./cmd/topology/ -v`
Expected: PASS（全部）

- [ ] Step 5: Commit

```bash
git add cmd/topology/query.go cmd/topology/topology.go cmd/topology/topology_test.go
git commit -m "feat(topology): add edge query subcommand"
```

### Task 9: backlinks 與 index 子命令（--write）

`Files:`

- Create: `cmd/topology/rewrite.go`
- Test: `cmd/topology/topology_test.go`

- [ ] Step 1: 寫失敗測試（附加）

```go
func TestBacklinksCommandWrite(t *testing.T) {
	root := writeFixture(t)
	out, err := run(t, "backlinks", "--write", "--root", root)
	if err != nil {
		t.Fatalf("backlinks: %v\n%s", err, out)
	}
	raw, err := os.ReadFile(filepath.Join(root, "core", "service-b.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "- calls ← [[service-a#Checkout]]") {
		t.Errorf("service-b backlinks not rebuilt:\n%s", raw)
	}
}

func TestIndexCommandWrite(t *testing.T) {
	root := writeFixture(t)
	if _, err := run(t, "index", "--write", "--root", root); err != nil {
		t.Fatalf("index: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "_index.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Entity Registry", "## Unlinked", "[[service-a]]"} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("_index.md missing %q", want)
		}
	}
}
```

- [ ] Step 2: 跑測試確認失敗

Run: `go test ./cmd/topology/ -run 'TestBacklinksCommandWrite|TestIndexCommandWrite' -v`
Expected: FAIL — unknown command

- [ ] Step 3: 實作 `cmd/topology/rewrite.go`，並在 `TopologyCmd()` 加
`cmd.AddCommand(BacklinksCmd())` 與 `cmd.AddCommand(IndexCmd())`

```go
package topology

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/cobra"
)

// BacklinksCmd returns the backlinks subcommand.
func BacklinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlinks",
		Short: "Recompute ## Backlinks sections from forward edges",
		RunE: func(cmd *cobra.Command, args []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("load topology: %w", err)
			}
			write, _ := cmd.Flags().GetBool("write")
			for _, n := range topo.Names() {
				e := topo.Entities[n]
				raw, err := os.ReadFile(e.Path)
				if err != nil {
					return fmt.Errorf("read %s: %w", e.Path, err)
				}
				updated := model.RenderBacklinksSection(string(raw), topo.BacklinksFor(n))
				if updated == string(raw) {
					continue
				}
				if !write {
					fmt.Fprintf(cmd.OutOrStdout(), "would update: %s\n", e.Path)
					continue
				}
				if err := os.WriteFile(e.Path, []byte(updated), 0o644); err != nil {
					return fmt.Errorf("write %s: %w", e.Path, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "updated: %s\n", e.Path)
			}
			return nil
		},
	}
	cmd.Flags().Bool("write", false, "write changes to files")
	return cmd
}

// IndexCmd returns the index subcommand.
func IndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Regenerate _index.md (registry, mermaid, unlinked; frontier preserved)",
		RunE: func(cmd *cobra.Command, args []string) error {
			topo, err := loadFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("load topology: %w", err)
			}
			idxPath := filepath.Join(topo.Root, "_index.md")
			existing := ""
			if raw, err := os.ReadFile(idxPath); err == nil {
				existing = string(raw)
			}
			out := topo.RenderIndex(existing)
			write, _ := cmd.Flags().GetBool("write")
			if !write {
				fmt.Fprint(cmd.OutOrStdout(), out)
				return nil
			}
			if err := os.WriteFile(idxPath, []byte(out), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", idxPath, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "updated: %s\n", idxPath)
			return nil
		},
	}
	cmd.Flags().Bool("write", false, "write _index.md")
	return cmd
}
```

- [ ] Step 4: 跑測試確認通過

Run: `go test ./cmd/topology/ -v`
Expected: PASS（全部）

- [ ] Step 5: Commit

```bash
git add cmd/topology/rewrite.go cmd/topology/topology.go cmd/topology/topology_test.go
git commit -m "feat(topology): add backlinks rebuild and index regeneration"
```

### Task 10: 全量驗證與文件同步

`Files:`

- Modify: `CLAUDE.md`（結構樹、模組對應表、測試清單）
- Modify: `plugins/general/skills/topology-builder/SKILL.md`（Verification 章節）

- [ ] Step 1: 全套件測試與真實樣本煙霧測試

Run:

```bash
go test ./... -count=1
go build -o /tmp/cc-plugin-test main.go
/tmp/cc-plugin-test topology verify --root plugins/general/skills/topology-builder/references
/tmp/cc-plugin-test topology unlinked --root plugins/general/skills/topology-builder/references
```

Expected: 測試全過；verify 輸出 `OK`；unlinked 輸出 no inbound `(none)`、
no outbound `billing-db`（與 references/_index.md 的 Unlinked 章節一致）

- [ ] Step 2: 更新 `CLAUDE.md`

- 結構樹 `cmd/` 下加一行：`│   ├── topology/             # topology 子命令群 — 知識圖譜驗證/查詢/重建`
- 模組對應表加一列：`| 拓撲圖譜操作      | cmd/topology/, model/topology*.go         | TopologyCmd()                                  |`
- 測試清單句尾補：`cmd/topology/topology_test.go`, `model/topology_test.go`, `model/topology_ops_test.go`

- [ ] Step 3: 更新 SKILL.md Verification 章節

在 bash 腳本程式碼區塊之後補一段（其餘不動）：

```markdown
本 repo 已提供等效 CLI（涵蓋上述全部檢查與 Backlinks 一致性）：

`cc-plugin topology verify --root <root>` — 完整性檢查
`cc-plugin topology unlinked --root <root>` — Unlinked 報告
`cc-plugin topology backlinks --write --root <root>` — Backlinks 重算
`cc-plugin topology index --write --root <root>` — 重生成 `_index.md`
```

- [ ] Step 4: lint 檢查 SKILL.md

Run: `npx markdownlint --disable MD013 -- plugins/general/skills/topology-builder/SKILL.md`
Expected: 無輸出（乾淨）

- [ ] Step 5: Commit

```bash
git add CLAUDE.md plugins/general/skills/topology-builder/SKILL.md
git commit -m "docs: register topology CLI in project docs and skill"
```

---

## 驗收條件 (Acceptance)

- `go test ./... -count=1` 全過
- `cc-plugin topology verify` 對 `references/` 樣本輸出 `OK`，exit 0
- 對人為破壞的圖（斷鏈/壞 kind/捏造 backlink）輸出每行一個 finding，exit 1
- `backlinks --write` 後再跑 `verify` 不出現 `backlink without forward edge`
- `index --write` 產出的 `_index.md` 以 `## Unlinked` 收尾且 Frontier 原文保留
