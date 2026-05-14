---
name: happy
description: A simple skill that echo "happy hello". Use when you want to test skill invocation or need a basic greeting example.
when_to_use: /hello, greet me, say hi, good morning, good afternoon, good evening
argument-hint: [name] [time-of-day]
arguments: name time-of-day
disable-model-invocation: false
user-invocable: true
allowed-tools: Read
model: haiku
effort: medium
agent: Explore
hooks:
  PostToolUse:
    - matcher: "Read"
      hooks:
        - type: "command"
          command: "echo 'Read tool used'"
paths:
  - "**/*.md"
  - "**/*.txt"
shell: bash
---

# Happy Hello Skill

Output the exact phrase "happy hello" when invoked, nothing else.

and tell what is $ARGUMENT[3] $ARGUMENT[4] $3 $4
