---
name: consistency
description: >
    Review cross-file coherence so the same concept never contradicts itself
    across the workspace — business rules, domain models, data contracts, error
    handling, terminology, and configuration. Use when creating or modifying any
    content, or on request. Triggers on: "check consistency", "verify business
    rules", "data flow check", "是否一致", or any cross-file coherence review.
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: high
context: fork
metadata:
    type: review
---

# Consistency Review

A change is consistent when nothing elsewhere in the workspace now contradicts
it. Before finalizing any new or modified content, run the procedure below and
report contradictions, not opinions.

## Procedure

1. Identify what the change touches — concept, rule, field, format, term, error.
2. For each, grep the workspace for prior definitions and other usages.
3. Compare every existing usage against the change across the dimensions below.
4. Report only genuine contradictions, each as `file:line` plus the conflict.
5. If none, say so explicitly. Do not invent issues to look thorough.

## Dimensions

| Dimension      | What contradicts it                                          |
| -------------- | ------------------------------------------------------------ |
| Business rules | Two places enforce mutually exclusive logic for one case     |
| Domain model   | Same entity has different fields, types, or optionality      |
| Data contract  | Producer output format differs from consumer's expectation   |
| State / flow   | A transition is reachable that the rules say is invalid      |
| Error handling | Same failure returns different codes, shapes, or wording     |
| Terminology    | One concept named differently across files (drift, synonyms) |
| Configuration  | A default, path, or limit set differently in two places      |

## Relationship rule

When two things reference each other (entity links, caller/callee contracts,
docs ↔ code), the reference must hold in `both` directions. A one-way reference
is an inconsistency.

## Output

```text
Consistency review
- [conflict] new rule X in a.go:42 contradicts b.go:88 (X handled inversely)
- [drift]    "tenant" (a.go) vs "account" (c.go) for the same concept
- [ok]       data contract a→b matches; error codes aligned
```

State the severity inline (`conflict` blocks, `drift` is cleanup, `ok` is
confirmation). Never report a contradiction without both sides cited.

## Related

- `[[naming-convention]]` for terminology drift at the identifier level
- `[[doc-sync]]` for docs-vs-code contradictions specifically
