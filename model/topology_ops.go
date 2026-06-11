package model

import (
	"fmt"
	"path/filepath"
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
	forward := map[string]bool{}
	for _, ed := range t.Edges() {
		forward[edgeKey(ed)] = true
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
