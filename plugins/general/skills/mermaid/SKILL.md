---
name: mermaid
description: >
    Use when creating, choosing, or writing the syntax for any Mermaid.js diagram —
    flowcharts, sequence / state / class / ER diagrams, gantt, kanban, mindmap,
    timeline, journey, requirement, C4, packet, block, architecture, treemap,
    treeView, pie / xy / quadrant / sankey / radar / venn charts, git graphs and
    ishikawa. Triggers on: "draw a diagram", "mermaid", "flowchart", "sequence
    diagram", "visualize this architecture", "which chart type should I use",
    "render a graph", "畫一張圖", "用 mermaid".
version: "1.0.0"
allowed-tools: Read, Glob, Bash
disable-model-invocation: false
user-invocable: true
effort: low
context: fork
metadata:
    type: reference
    platforms: [macos, linux]
    homepage: [https://mermaid.](https://mermaid.ai/open-source/intro/)ai/open-source/syntax
---

# mermaid

Pick the right Mermaid.js diagram, then open its reference file for the exact
keyword, syntax notes and a copy-paste example.

## How to use

1. Match the user's intent to a row in the selection tables below.
2. Read the linked `references/NN-*.md` for that type — each file holds the
   keyword, a one-line 說明, and a ready-to-render example.
3. Copy the fenced `mermaid` block, adapt the labels/data, and emit it.
4. If the output target is a **terminal** (no markdown preview / browser),
   pipe the diagram through `mermaid-ascii` to render it as ASCII art
   (see Terminal rendering below).

All examples are valid against current Mermaid.js; `-beta` keywords are the
official spelling and must be kept (they are not placeholders).

## Terminal rendering

When the user is working in a terminal context (SSH, CLI output, no browser),
use [`mermaid-ascii`](https://github.com/AlexanderGrooff/mermaid-ascii) to
render diagrams as Unicode/ASCII art instead of emitting a raw fenced block.

### Prerequisites

If `mermaid-ascii` is not on `$PATH`, install it first:

```bash
go install github.com/AlexanderGrooff/mermaid-ascii@master
```

### Usage

```bash
# From stdin (pipe the mermaid syntax)
printf 'graph LR\n  A --> B --> C' | mermaid-ascii -f-

# From a file
mermaid-ascii -f diagram.mmd
```

### Key flags

| Flag | Description |
| ---- | ----------- |
| `-f-` | Read mermaid from stdin |
| `-f FILE` | Read from file |
| `-a` | Pure ASCII (no Unicode box-drawing) |
| `-x N` | Horizontal spacing between nodes (default 5) |
| `-y N` | Vertical spacing between nodes (default 5) |
| `-p N` | Padding between text and border (default 1) |

### Supported diagram types

`mermaid-ascii` supports `graph` / `flowchart` (TD, TB, LR) and
`sequenceDiagram`. Other Mermaid types (state, class, ER, gantt, etc.) are
**not supported** — fall back to a fenced mermaid block for those.

### Example output

```
┌───┐     ┌───┐     ┌───┐
│   │     │   │     │   │
│ A ├────►│ B ├────►│ C │
│   │     │   │     │   │
└───┘     └───┘     └───┘
```

## Selection guide

### 一、流程、架構與軟體設計 (Process & Architecture)

| Use it for                                       | Type             | Keyword               | Reference                                                           |
| ------------------------------------------------ | ---------------- | --------------------- | ------------------------------------------------------------------- |
| Process / decision / branching flow              | Flowchart        | `flowchart` / `graph` | [01-flowchart.md](references/01-flowchart.md)                       |
| Time-ordered messages between actors/services    | Sequence Diagram | `sequenceDiagram`     | [02-sequence-diagram.md](references/02-sequence-diagram.md)         |
| Finite state machine, lifecycle transitions      | State Diagram    | `stateDiagram-v2`     | [03-state-diagram.md](references/03-state-diagram.md)               |
| OOP classes, attributes, inheritance/composition | Class Diagram    | `classDiagram`        | [04-class-diagram.md](references/04-class-diagram.md)               |
| Database tables and cardinality                  | ER Diagram       | `erDiagram`           | [05-er-diagram.md](references/05-er-diagram.md)                     |
| Cloud infra / deployment topology                | Architecture     | `architecture-beta`   | [06-architecture-diagram.md](references/06-architecture-diagram.md) |
| Stacked high-level system blocks                 | Block Diagram    | `block-beta`          | [07-block-diagram.md](references/07-block-diagram.md)               |
| C4 software context / containers                 | C4 Diagram       | `C4Context`           | [08-c4-diagram.md](references/08-c4-diagram.md)                     |
| Binary packet header field layout                | Packet Diagram   | `packet`              | [09-packet-diagram.md](references/09-packet-diagram.md)             |

### 二、專案管理、規劃與組織 (Planning & Organization)

| Use it for                                 | Type                | Keyword              | Reference                                                         |
| ------------------------------------------ | ------------------- | -------------------- | ----------------------------------------------------------------- |
| Project schedule, task spans, dependencies | Gantt Chart         | `gantt`              | [10-gantt-chart.md](references/10-gantt-chart.md)                 |
| Agile board by task status                 | Kanban              | `kanban`             | [11-kanban.md](references/11-kanban.md)                           |
| Brainstorm radiating from a core idea      | Mindmap             | `mindmap`            | [12-mindmap.md](references/12-mindmap.md)                         |
| Linear sequence of dated events            | Timeline            | `timeline`           | [13-timeline.md](references/13-timeline.md)                       |
| User experience, pain points, scores       | User Journey        | `journey`            | [14-user-journey.md](references/14-user-journey.md)               |
| Requirements traced to test cases          | Requirement Diagram | `requirementDiagram` | [15-requirement-diagram.md](references/15-requirement-diagram.md) |
| UI / command / event / read-model flow     | Event Modeling      | `eventmodeling`      | [16-event-modeling.md](references/16-event-modeling.md)           |
| Nested rectangles sized by weight          | Treemap             | `treemap-beta`       | [17-treemap.md](references/17-treemap.md)                         |
| Plain-text file/dir hierarchy              | TreeView            | `treeView-beta`      | [18-treeview.md](references/18-treeview.md)                       |

### 三、數據、圖表與策略分析 (Data & Analytics)

| Use it for                            | Type           | Keyword         | Reference                                               |
| ------------------------------------- | -------------- | --------------- | ------------------------------------------------------- |
| Share of a whole (percentages)        | Pie Chart      | `pie`           | [19-pie-chart.md](references/19-pie-chart.md)           |
| Bars + lines on x/y axes              | XY Chart       | `xychart-beta`  | [20-xy-chart.md](references/20-xy-chart.md)             |
| 2x2 prioritization / strategy matrix  | Quadrant Chart | `quadrantChart` | [21-quadrant-chart.md](references/21-quadrant-chart.md) |
| Flow / volume between nodes           | Sankey Diagram | `sankey-beta`   | [22-sankey-diagram.md](references/22-sankey-diagram.md) |
| Multi-axis skill/attribute comparison | Radar Chart    | `radar-beta`    | [23-radar-chart.md](references/23-radar-chart.md)       |
| Business value-chain vs evolution map | Wardley Map    | `wardley-beta`  | [24-wardley-map.md](references/24-wardley-map.md)       |
| Set overlap / intersection            | Venn Diagram   | `venn-beta`     | [25-venn-diagram.md](references/25-venn-diagram.md)     |

### 四、版本控制與品質診斷 (Technical & Quality)

| Use it for                             | Type      | Keyword    | Reference                                                   |
| -------------------------------------- | --------- | ---------- | ----------------------------------------------------------- |
| Git commits, branches, merges          | Git Graph | `gitGraph` | [26-git-graph.md](references/26-git-graph.md)               |
| Cause-and-effect / root-cause analysis | Ishikawa  | `ishikawa` | [27-ishikawa-diagram.md](references/27-ishikawa-diagram.md) |

Full index: [references/README.md](references/README.md).

## Common mistakes

- Dropping the `-beta` suffix from beta keywords (`architecture-beta`,
  `treeView-beta`, `sankey-beta`, etc.) — the diagram will not render.
- Treeview/treemap/mindmap hierarchy is defined by **indentation only**; mixed
  or inconsistent indent levels silently break the nesting.
- Forgetting `stateDiagram-v2` (v1 is legacy) or `classDiagram` casing.
- Wrapping the example's `title`/labels in the wrong quote style — keep the
  quoting shown in each reference file.
- Using `mermaid-ascii` with unsupported diagram types (anything other than
  `graph`/`flowchart` and `sequenceDiagram`) — it will error out; use a
  fenced mermaid block instead.
- Putting the entire mermaid on one semicolon-separated line when piping to
  `mermaid-ascii` — it requires newline-separated syntax.
