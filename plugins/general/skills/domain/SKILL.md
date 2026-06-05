---
name: domain
description: >
  Loads project domain knowledge for the current working directory. Reads
  README.md (business domain + project architecture) at the repo root and merges
  with CLAUDE.md if present. Returns a structured context block for the calling
  agent. Invoke at the start of any feature or refactor session to establish
  project context. Can be invoked by the user or by a calling agent.
allowed-tools: Read, Bash, Glob
disable-model-invocation: false
user-invocable: true
effort: medium
context: fork
---

# domain

Load and surface the business and technical context for the current project.

## Input Sources (checked in order)

1. `README.md` (repo root) — primary domain knowledge (business domain + architecture)
2. `CLAUDE.md` (repo root) — supplementary conventions, if present

## Procedure

### Step 1 — Load primary domain knowledge

Check if `README.md` exists at the repo root.

- **If yes:** read it in full.
- **If no:** emit the following warning and continue with Step 2:

```
WARN: README.md not found at repo root.
  Run `/domain-init` to generate domain knowledge from this project.
  Proceeding with minimal context inferred from project structure.
```

Then run:
```bash
find . -maxdepth 3 -type d -not -path '*/.git/*' -not -path '*/vendor/*' | sort
```

List the top-level directories and infer their purpose from names (e.g., `handler/`,
`service/`, `model/`, `cmd/`, `internal/`). Include this as the "Architecture" section
of the output.

### Step 2 — Merge supplementary conventions

Check if `CLAUDE.md` exists at the repo root.

- **If yes:** read it and incorporate any conventions, constraints, or known patterns
  into the "Conventions" section of the output.
- **If no:** skip silently.

### Step 3 — Emit structured output

Output the following block so the calling agent can parse it:

```
### DOMAIN CONTEXT — <project name>

**Business Domain:** <1-3 sentences from the Overview section of README.md,
or inferred from project structure if README is missing>

**Key Entities:** <bullet list of core domain objects>

**Architecture:** <layer summary — packages, entry points, external services>

**Conventions:** <naming, error handling, testing approach from CLAUDE.md if present>
```

## Failure Modes

- Neither `README.md` nor `CLAUDE.md` exists → warn, infer from directory structure,
  suggest `domain-init`, continue. Never block the calling agent.
- `find` command returns no results → output "Architecture: unable to infer (empty project)"
  and continue.
