---
name: dependency-hygiene
description: >
    Review a project's dependencies for hygiene — unused or duplicate packages,
    unpinned or wildcard versions, stale majors, heavy deps used trivially, and
    direct-vs-transitive confusion. Language-aware across go.mod, package.json,
    requirements/pyproject, and lockfiles. Use when auditing dependencies.
    Triggers on: "review dependencies", "dependency audit", "unused packages",
    "依賴審查", "are these deps needed", "check go.mod".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: false
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: review
---

# Dependency Hygiene Review

Fewer, current, justified dependencies. This skill reads the manifest, checks
each dependency against real usage, and reports what to drop, pin, or upgrade —
it does not run installs or modify lockfiles unless asked.

## Procedure

1. Locate the manifest(s): `go.mod`, `package.json`, `requirements.txt`,
   `pyproject.toml`, plus their lockfiles.
2. For each direct dependency, grep the source for an actual import or call.
3. Apply the checks below; prefer the language's own tooling when available
   (`go mod tidy -diff`, `npm ls`, `pip check`) and report its output.
4. Output a per-dependency verdict: keep, drop, pin, or upgrade.

## Checks

| Check             | Smell                                                      |
| ----------------- | ---------------------------------------------------------- |
| Unused            | Declared but never imported or called                      |
| Duplicate         | Two deps covering the same need; one stdlib-replaceable    |
| Unpinned          | Wildcard / range / missing version where a pin is expected |
| Stale             | Major versions behind, or unmaintained upstream            |
| Heavy-for-trivial | A large dep pulled in for one small function               |
| Mis-scoped        | A runtime dep that should be dev-only, or vice versa       |
| Phantom direct    | Used transitively but not declared as a direct dependency  |

## Output

```
Dependency hygiene review (go.mod)
- drop:    github.com/x/unused — no imports found
- pin:     example.com/lib @ latest → pin to v1.4.2
- upgrade: zap v1.21 → v1.27 (3 majors? no, minors; safe)
- keep:    cobra, viper, gorm — all used, current
```

State the source of truth (which manifest) and never claim a dep is unused
without showing the empty grep. Recommend; do not auto-remove unless asked.

## Related

- `[[consistency]]` if two modules pin the same lib to different versions
- `[[doc-sync]]` when the documented tech stack lists deps no longer present
