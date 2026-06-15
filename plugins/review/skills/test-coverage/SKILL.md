---
name: test-coverage
description: >
    Review tests for the gaps that matter — untested core logic, missing error
    and edge-case paths, assertion-free or tautological tests, and brittle
    coupling to internals. Prioritizes by business risk, not raw coverage
    percentage. Use when auditing test quality or completeness. Triggers on:
    "review test coverage", "what's untested", "are these tests good",
    "測試覆蓋", "missing tests", "test gaps".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: review
---

# Test Coverage Review

Coverage percentage measures lines executed, not behavior verified. This skill
finds the gaps that carry risk and the tests that only look like tests, then
ranks what to add by business impact.

## Procedure

1. Map the code: which packages hold core business logic vs glue. Run the
   coverage tool if cheap (`go test ./... -cover`) and read the report.
2. For the highest-risk code, check whether its error paths, edge cases, and
   branches are actually asserted — not merely executed.
3. Inspect existing tests for the quality smells below; a high-coverage file can
   still be untested in the ways that matter.
4. Output a risk-ranked gap list, plus any existing tests to strengthen.

## What to look for

| Gap / smell      | Why it matters                                           |
| ---------------- | -------------------------------------------------------- |
| Untested core    | The logic that earns or loses money has no test           |
| Missing error path | Only the happy path is covered                          |
| Edge cases       | Empty, nil, boundary, concurrency, timezone untested      |
| No assertion     | Test runs code but asserts nothing (false coverage)       |
| Tautology        | Assertion restates the implementation; always passes      |
| Over-mocking     | Mocks so heavy the test proves only the mock works         |
| Brittle coupling | Asserts on internals; breaks on safe refactors            |

## Output

```
Test coverage review (ranked by risk)
1. [high] cmd/distill.go fingerprint dedup — no test for hash collision path
2. [high] model/store.go retain() — error path untested
3. [med]  ollama_test.go — runs Extract() but asserts nothing on output
4. strengthen: state_test.go couples to a private field; assert behavior instead
```

Rank by what failure would cost, not by line count. A blunt "add a test for X"
beats a coverage number with no context.

## Related

- `[[business-improvement]]` to know which paths carry the most business risk
- `[[consistency]]` when tests encode rules that contradict the code
