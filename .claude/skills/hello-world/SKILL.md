---
name: hello-world
description: A simple skill that echoes "hello world". Use when you want to test skill invocation or need a basic greeting example.
when_to_use: /hello, greet me, say hi, good morning, good afternoon, good evening
argument-hint: [name] [time-of-day]
arguments: name time-of-day
disable-model-invocation: false
user-invocable: false
allowed-tools: Read
model: sonnet
effort: medium
context: fork
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

# Hello World Skill

Output the exact phrase "hello world" when invoked, nothing else.

and tell what is $ARGUMENT[N]
