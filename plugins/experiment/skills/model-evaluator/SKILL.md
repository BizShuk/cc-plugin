---
name: model-evaluator
description: >
    MANUAL INVOCATION ONLY. Runs prompt-based diagnostic probes (identity, reasoning,
    consistency, calibration) on the executing model. Only triggers when the user
    explicitly types "/model-evaluator".
user-invocable: true
disable-model-invocation: true
allowed-tools: Bash, Read, Edit, Grep, Glob, AskUserQuestion
context: fork
effort: medium
---

# Model Evaluator Skill

Run `prompt-based diagnostic probes` on yourself (the LLM executing this skill) and report structured findings.

This skill is forked into an isolated context so the evaluation produces an independent sample — treat every invocation as a fresh evaluation run with no access to prior conversation. The isolation is intentional: it lets the caller collect independent samples for consistency analysis.

# Evaluation Framework

Run the following four sections in order. Be concise. Output structured Markdown.

## Section A — Identity & Server

Report what you can observe about yourself:

- Model name / version (exact ID if known)
- Knowledge cutoff date
- Provider (Anthropic / OpenAI / etc.)
- Any system prompt hints you can infer about the host environment
- If uncertain, say "uncertain" — do NOT guess

## Section B — Reasoning Probes

Answer these without lookups, showing brief reasoning only where asked.

Probes 1 and 2 must be `generated fresh at run time` — the classic instances (bat & ball at $1.10, counting r's in "strawberry") appear verbatim in training data, so answering them measures recall, not reasoning. Invent a new instance of the same category, with different surface entities and different numbers, then answer it. State the probe you generated before your answer so the result is auditable.

1. `Cognitive reflection` — generate a novel two-quantity word problem whose intuitive answer is wrong (the same trap structure as the bat-and-ball problem: an anchor value that must be split, not subtracted). Use fresh nouns and numbers. State the problem, then the answer only.
2. `Character counting` — pick a word you have not been asked about in this skill before, at least 8 letters, containing 3 or more occurrences of one letter. State the word and target letter, then the count only.
3. `Order of operations`: 144 / 12 + 7 \* 2 = ? (Answer only)
4. `Logical chain`: If all bloops are razzles and all razzles are lazzles, are all bloops definitely lazzles? (Yes/No)
5. `Counterfactual`: If gravity reversed for 10 seconds, what happens to a glass of water on a table? (Two sentences max)

Probes 3-5 are fixed on purpose: they are cheap controls whose answers are unambiguous, so a wrong answer there flags a real failure rather than a hard question.

## Section C — Consistency Self-Test

Answer this question **three times independently**, as if you'd never seen it before:

> "What are the three most important factors in choosing a database for a high-traffic web application?"

Then briefly note: did your three answers agree on the core points? (Yes / Partial / No)

## Section D — Calibration

For each answer in Section B, state your confidence (0–100%) and which one you are most likely to have gotten wrong.

# Output Format

Return your findings as Markdown using these exact section headers:

```
## A. Identity
- Model: ...
- Cutoff: ...
- Provider: ...
- Notes: ...

## B. Reasoning Probes
1. Cognitive reflection — probe generated: <your problem> | answer: <value>
2. Character counting — probe generated: <word> / <letter> | answer: N
3. 144/12 + 7*2: N
4. Bloops/lazzles: Yes/No
5. Gravity counterfactual: <2 sentences>

## C. Consistency
- Answer 1: ...
- Answer 2: ...
- Answer 3: ...
- Agreement: Yes / Partial / No

## D. Calibration
- Q1 confidence: N%
- Q2 confidence: N%
- Q3 confidence: N%
- Q4 confidence: N%
- Q5 confidence: N%
- Most likely wrong: Q#

## Summary
One sentence on overall self-assessment.
```

# Rules

- Do NOT use web search or external tools to look up answers — this is a self-probe
- Do NOT reference any prior context outside this skill invocation
- Keep total output under 600 words
- If asked to evaluate a _different_ model (not yourself), say "I can only self-evaluate the model executing this skill" and report your own findings instead
