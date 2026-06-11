package model

import (
	"fmt"
	"path/filepath"
	"regexp"
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
