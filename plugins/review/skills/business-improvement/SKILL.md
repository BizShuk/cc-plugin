---
name: business-improvement
description: >
    Review a feature, flow, or codebase for business-value improvements — gaps,
    friction, missed revenue or retention levers, risky edge cases, and
    redundant work — then propose ranked, concrete changes tied to outcomes.
    Use when asked how to improve the product or business logic. Triggers on:
    "business improvement", "how to improve this", "業務改善", "增加價值",
    "what's missing", "review the business value".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: xhigh
context: fork
metadata:
    type: review
---

# Business Improvement Review

Find where the product underserves its own business goal, then propose changes
ranked by value-over-effort. This is a review, not a redesign: every suggestion
must cite the code or flow it improves and the outcome it moves.

## Procedure

1. Establish the goal — read `README.md` and `CLAUDE.md` for the stated business
   domain; if absent, infer it from entry points and name it back explicitly.
2. Trace the core flow end to end (acquisition → activation → core action →
   retention). Note where the user does work the system could do.
3. Score each finding on the matrix below; drop anything you cannot tie to an
   outcome.
4. Output a ranked list. Lead with the single highest value-over-effort item.

## What to look for

| Lens          | Improvement signal                                          |
| ------------- | ----------------------------------------------------------- |
| Friction      | Extra steps, manual input, or waits the system could remove |
| Gap           | A promised value the flow never actually delivers           |
| Leverage      | A small change unlocking revenue, retention, or activation  |
| Risk          | An edge case that silently loses data, money, or trust      |
| Redundancy    | Duplicated effort or features that cancel each other out    |
| Measurability | An outcome that is happening but never measured             |

## Output

```text
Business improvement — <goal in one line>
1. [high value / low effort] <change> — improves <outcome>, at <file:line/flow>
2. [med / med]               <change> — improves <outcome>, at <...>
   Risk: <what breaks if ignored>
```

Rank strictly by value-over-effort. Do not list more than the top items that
genuinely matter; a short ranked list beats an exhaustive one.

## Related

- `[[consistency]]` if an improvement would contradict an existing rule
- Use the `business-extract` skill first when the business goal is unclear
