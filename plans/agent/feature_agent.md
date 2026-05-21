# Feature Agent Design Plan

## Context

The goal is a **general-purpose feature implementation agent** that works across any project.
It brings domain context, architectural conventions, and a structured NFR checklist to every
feature task — so the model doesn't drift into pure functional implementation and miss
security, observability, or testability concerns.

The agent is **user-level** (`~/.claude/agents/feature.md`) so it's available in every
project without per-repo configuration. Its skills live in the **cc-plugin plugin**
(`skills/<name>/SKILL.md`) so they're distributed and versioned alongside the other
Golang skills. Per-project configuration lives in `workspace/` (which is gitignored or
committed per team preference).

---

## Files to Create

| File | Location |
|------|----------|
| Feature agent | `agents/feature.md` |
| Domain skill | `skills/domain/SKILL.md` (not ready — will patch later) |
| Golang MVC skill | `skills/golang-mvc/SKILL.md` |

---

## 1. Feature Agent — `~/.claude/agents/feature.md`

### Frontmatter

```yaml
---
name: feature
description: >
  General-purpose feature implementation agent. Use when building, planning, or
  scoping any new feature end-to-end: "implement X", "add support for Y", "build
  the Z endpoint/service/flow". Does NOT trigger for refactoring, dead-code removal,
  or performance review — use golang-refactor for those.
tools: Read, Edit, Write, Bash, Grep, Glob, AskUserQuestion, TodoWrite
model: inherit
permissionMode: acceptEdits
skills: golang-code-quality, golang-dev, golang-mvc, golang-naming
mcpServers:
hooks:
memory: local
background: false
effort: xhigh
isolation: worktree
color: cyan
initialPrompt:
---
```

### Body Structure (5 parts)

**Part 1 — Role/Identity**
Senior software engineer. Correctness before cleverness. Smallest vertical slice first.
Makes implicit NFRs explicit before writing any code.

**Part 2 — Feature Description (Embedded Prompt)**
The feature description is the **user's invocation message** (`@feature implement X...`),
NOT `initialPrompt`. The body uses marker comments so the agent treats the invocation
message as authoritative input:

```
<!-- FEATURE_DESCRIPTION_START -->
Treat the user's opening message as the feature requirement.
If no feature is described, ask before proceeding.
<!-- FEATURE_DESCRIPTION_END -->
```

`initialPrompt` is reserved for static bootstrapping (e.g., "load domain context first").

**Part 3 — NFRs (ranked by criticality)**

| Priority | NFR | Blocking? |
|---|---|---|
| 1 | Security (auth/z, input validation, no secrets in logs) | Yes — blocking |
| 2 | Testability (interface injection, table-driven tests) | Yes — required before merge |
| 3 | Observability (structured logs, metrics, trace spans) | Yes — required before prod |
| 4 | Performance (hot paths, N+1 queries, pooling) | No — validate post-deploy |
| 5 | Maintainability (naming, no duplication, constants) | No — ongoing |
| 6 | Backward compatibility (migration plan if schema changes) | Context-dependent |

Agent batches all NFR ambiguities into a single `AskUserQuestion` call in Phase 2.

**Part 4 — Workflow Phases**

1. **Understand** — read feature description; for Go projects invoke `golang-mvc` skill; identify affected packages/models/APIs; confirm summary with user
2. **Clarify** — batch ALL questions (feature + NFR ambiguities) into one `AskUserQuestion`; document decisions
3. **Plan** — write implementation plan (new files, modified packages, interface definitions, migration steps, test plan); get user approval before writing code
4. **Implement** — follow `golang-mvc` layer rules; build top-down (models → repos/services → handlers → wiring); write tests alongside each layer; run `go build ./... && go vet ./...` after each package
5. **Verify** — run full test suite; walk NFR checklist section-by-section; output done summary (files changed, NFRs met or deferred with justification)

`TodoWrite` tracks progress across all phases from the start.

**Part 5 — Skill Routing**

| Scenario | Action |
|---|---|
| Go project | Always invoke `golang-mvc` in Phase 1 |
| Non-Go project | Use language/framework conventions |
| Feature touches DB schema | Flag in Phase 2: clarify migration plan required |
| Feature spans multiple repos | Scope to one repo; instruct user to invoke `@feature` a second time in the other repo |

---

## 2. Domain Skill — `skills/domain/SKILL.md`

**Decision: skill, not direct file read.** A skill gracefully degrades when
`workspace/README.md` doesn't exist, merges multiple sources (`workspace/README.md` +
`CLAUDE.md`), and reads `workspace/domains.json` to return the external skill registry.
A raw file read from the agent body would fail silently on a fresh project.

```yaml
---
name: domain
description: >
  Loads project domain knowledge. Reads workspace/README.md (business domain +
  architecture) and merges with CLAUDE.md if present. Returns structured context
  block and external skill registry from workspace/domains.json. Invoke at the
  start of any feature or refactor session. Can be invoked by user or model.
allowed-tools: Read, Bash, Glob
disable-model-invocation: false
user-invocable: true
effort: medium
context: fork
---
```

**Degradation behavior:**
- `workspace/README.md` missing → warn, scan top-level packages from filesystem, suggest running `domain-init`
- `CLAUDE.md` missing → skip silently
- Entry in `workspace/domains.json` points to a non-existent path → emit `WARN:` block per missing skill (not a hard failure)

**Output format** (structured so calling agent can parse it):
```
### DOMAIN CONTEXT — <project>
Business Domain: ...
Key Entities: ...
Architecture: ...
Conventions: ...
External Skill Registry: <contents of domains.json or "none">
```

---

## 3. Golang MVC Skill — `skills/golang-mvc/SKILL.md`

Actively enforces and applies MVC layer conventions. Modifies and refactors code
that violates layer boundaries. Distinct from `golang-code-quality` which focuses
on SOLID and code-level patterns.

```yaml
---
name: golang-mvc
description: >
  Go MVC architecture conventions for feature implementation and refactoring.
  Actively enforces and applies layer conventions — modifies and refactors code
  that violates MVC boundaries.
allowed-tools: Bash, Read, Edit, Grep, Glob, AskUserQuest
disable-model-invocation: true
user-invocable: false
effort: medium
context: fork
---
```

**Key content sections:**

- **Layer map** — handler / service / repository / model / config — with import rules (what each may/must-not import)
- **Feature checklist** — model first → repo/service → handler → validation → wiring in `main.go`/`bootstrap.go`
- **Interface placement rule** — interfaces defined at the consumer (handler), not the implementor (repo). Critical for testability.
- **Constructor injection rules** — ≤4 deps: plain params; 5+ deps: `Options` struct
- **Error handling** — wrap at every layer boundary; sentinel errors in owning package; HTTP mapping only at handler; log once
- **Context conventions** — always first param; never stored in struct; `WithTimeout` wraps every DB/external call
- **Test patterns** — handler tests mock interfaces + `httptest`; repo tests use real DB or `sqlmock`; table-driven tests for validation

---

## 4. Domain Init Skill — REMOVED

> Removed from the project. Domain knowledge bootstrapping will be handled
> differently when the `domain` skill is ready.

---

## 5. External Domain Skill Loading Strategy

**Primary: `workspace/domains.json` registry**

```json
{
  "skills": [
    { "name": "payment-domain", "path": "~/.claude/skills/payment-domain" },
    { "name": "logistics-domain", "path": "../acme-core/skills/logistics" }
  ]
}
```

- Committed to the project repo — every team member sees the same set
- Consistent with `workspace/README.md` convention (`workspace/` = project meta dir)
- `domain` skill reads this and returns registry contents to the calling `feature` agent
- Agent invokes each listed skill in Phase 1 before Clarify

**Secondary: symlinks in `~/.claude/skills/`**

For skills needed across many projects without per-project `domains.json` entries:
```bash
ln -s ~/projects/acme/skills/payment ~/.claude/skills/acme-payment
```

Then reference by name in `domains.json`. This gives: declarative (JSON) + globally available (symlink).

**Why not arguments or settings.json:**
- `--domain` flag requires the user to remember it every invocation (error-prone)
- `.claude/settings.json` mixes permissions/env with domain routing (wrong abstraction)

---

## Improvements Over Original Spec

1. **`initialPrompt` is not the injection point.** It's static — baked into the file. The feature description should be the user's invocation message (`@feature implement X`). The body uses marker comments to make the contract explicit. `initialPrompt` is reserved for static bootstrapping.

2. **NFRs must be ranked, not flat.** An unranked list treats security and changelog entries equally. Ranking by criticality (Security → Testability → Observability → Performance → Maintainability) tells the agent which are blockers vs. nice-to-haves, preventing over-engineering on simple features.

3. **Domain skill needs graceful degradation.** Binary present/missing logic blocks the agent on fresh projects. Degraded mode (filesystem scan + warning + suggestion to run `domain-init`) makes the feature agent immediately usable on any project.

4. **`golang-mvc` must be separate from `golang-code-quality`.** Review (backward-looking, Edit/Write) and generation guidance (forward-looking, read-only) have conflicting tool sets and mental models. Merging them would pollute new-code workflows with review-mode instructions.

5. **Parallel project scaling is explicit, not implicit.** The agent scopes itself to one repo. For cross-repo features, it instructs the user to invoke `@feature` a second time in the other repo (`isolation: worktree` already makes this safe). Self-orchestrating multi-repo parallelism from one agent would blow out context.

6. **`domain-init` needs idempotency.** Naive regeneration overwrites hand-written notes, so teams stop re-running it. The `<!-- MANUAL -->` marker pattern protects hand-written sections and outputs a diff on re-run, making it safe to invoke as a post-refactor hook.

7. **Missing external skill paths must warn explicitly.** If a path in `domains.json` doesn't exist on the current machine, the `domain` skill emits a `WARN:` block per missing entry and continues. Silent failure mid-Phase-4 is much harder to debug.

---

## Verification

1. **Domain skill**: create a project without `workspace/README.md`; invoke `/domain` — should warn and list packages, not fail
2. **Domain init**: run `/domain-init` on a Go project → check `workspace/README.md` is generated; run again → check hand-written `<!-- MANUAL -->` sections are preserved
3. **Golang MVC**: invoke `@feature implement a new user registration endpoint` in a Go project — agent should load both `domain` and `golang-mvc` in Phase 1, produce a layered plan (model → repo → handler), and ask NFR questions before writing code
4. **External skills**: add a `workspace/domains.json` with a non-existent path → `domain` skill should emit `WARN:` and continue
5. **Parallel projects**: invoke `@feature` in two separate terminal sessions on two different repos — both run in separate worktrees without collision (`isolation: worktree`)
