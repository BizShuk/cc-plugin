---
name: happy
description: Test skill that echoes a greeting; use to verify skill invocation works.
argument-hint: "[name] [time-of-day]"
disable-model-invocation: false
user-invocable: true
allowed-tools: Read
model: haiku
effort: medium
---

# Happy Hello Skill

Output the exact phrase "happy hello" when invoked, nothing else.

Then report the arguments received: `$1` is the name, `$2` is the time of day.
`$ARGUMENTS` holds the full argument string when you need it verbatim.
