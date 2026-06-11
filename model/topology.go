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
