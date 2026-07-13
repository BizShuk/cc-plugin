---
name: auto-evolving
description: >
    Initialize and run the Auto-Evolving Knowledge Engine to discover system blind spots, synthesize solutions across four philosophy frameworks (Taoism, Socratic Dialectics, Systems Theory, First Principles), evaluate value, and evolve a project's core foundations (axioms, principles, hypotheses) with confidence decay.
    Triggers on: "auto-evolving", "knowledge engine", "evolution cycle", "confidence decay", "promote hypothesis", "auto-evolving skill".
version: "1.0.0"
allowed-tools: Read, Write, Bash, Glob
user-invocable: true
disable-model-invocation: false
effort: high
context: workspace
metadata:
    type: methodology
    platforms: [macos, linux]
---


# Auto-Evolving Knowledge Engine

This skill guides Codex through initializing, running, and managing the Auto-Evolving Knowledge Engine in any project workspace.

## Core Operations

### 1. Initialize the Workspace

When the user requests to setup the engine or initialize the folder structure, copy all directories and templates from the bundled `template/` directory inside this skill to the root of the target workspace.

The structure copied should look like this:
```
<workspace>/
├── README.md
├── CLAUDE.md
├── AGENTS.md -> CLAUDE.md
├── README.todo
├── engine/
│   ├── system_prompt.md
│   ├── config.yaml
│   └── philosophy.md
├── foundation/
│   ├── axioms/
│   │   └── README.md
│   ├── principles/
│   │   └── README.md
│   ├── hypotheses/
│   │   └── README.md
│   ├── retired/
│   │   └── README.md
│   └── confidence.yaml
├── cycles/
│   └── README.md
├── branches/
│   ├── harmony/
│   ├── dialectic/
│   ├── velocity/
│   ├── incremental/
│   └── scoreboard.md
├── scripts/
│   ├── kb_manager.py
│   ├── run_cycle.sh
│   ├── decay.sh
│   └── promote.sh
└── tmp/
```

Make sure the helper shell scripts under `scripts/` are marked as executable:
```bash
chmod +x scripts/kb_manager.py scripts/run_cycle.sh scripts/decay.sh scripts/promote.sh
```

### 2. Execute a Cycle

To run an evolution cycle:
1. Locate or read `<workspace>/CLAUDE.md` and `<workspace>/README.md` to anchor your context.
2. Run `./scripts/run_cycle.sh` to initialize a new markdown cycle file under `cycles/YYYY-MM/YYYY-MM-DD-NNN.md`.
3. Follow the rules defined in `engine/system_prompt.md` to:
   - **DISCOVERY (探)**: Identify exactly one critical blind spot, unexplored possibility, or missing knowledge in the workspace.
   - **SYNTHESIS (合)**: Propose solutions based on the four philosophical branches (Harmony, Dialectics, Velocity, Incremental) and save them to the respective `branches/` folders.
   - **EVALUATION (衡)**: Score the proposals against four axes: System Value, Business Value, Epistemic Value, and Antifragility.
   - **FOUNDATION UPDATE (沉)**: Save approved knowledge to `hypotheses/` and register it in `confidence.yaml` with initial confidence `0.5`.

### 3. Manage Confidence Decay

Every cycle (or regular interval), confidence for items that are not verified/cited should decay.
Run:
```bash
./scripts/decay.sh
```
This triggers `scripts/kb_manager.py decay` which automatically updates all confidence ratings in `foundation/confidence.yaml`. Any entry dropping below the demotion threshold (`0.3` by default) should be flagged for demotion or retired to the `foundation/retired/` archive directory.

### 4. Promote and Demote Knowledge Items

When hypotheses are verified through implementation or tests, they are promoted.
Run:
```bash
./scripts/promote.sh <item_id> [target_layer]
```
For example, to promote hypothesis `HY-001` to a principle:
```bash
./scripts/promote.sh HY-001 principle
```
This moves the file from `foundation/hypotheses/HY-001-<topic>.md` to `foundation/principles/PR-001-<topic>.md` and updates the entry in `foundation/confidence.yaml`.
