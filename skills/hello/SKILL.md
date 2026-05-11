---
name: hello
description: Greet the user with a friendly personalized message. Use when saying hello, greeting someone new, or starting a conversation.
when_to_use: /hello, greet me, say hi, good morning, good afternoon, good evening
argument-hint: [name]
arguments: name
disable-model-invocation: false
user-invocable: true
allowed-tools: Read
model: inherit
effort: medium
context: null
agent: null
hooks: null
paths: null
shell: bash
---

# Hello Skill

Greet the user by name and ask how you can help them today.

## Greeting

Say hello to **$name** warmly and enthusiastically. Ask what they'd like to work on today.

## Example

If the user says `/hello Alice`, respond with:
"Hello Alice! Great to see you. What can I help you with today?"
