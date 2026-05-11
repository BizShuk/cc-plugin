---
name: hello
description: Greet the user with a friendly personalized message. Use when saying hello, greeting someone new, or starting a conversation.
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

# Hello Skill

Greet the user by name and ask how you can help them today.

## Greeting Formats

Based on the time of day, adjust your greeting appropriately:

| $1 (time-of-day) | Greeting Style   |
| ---------------- | ---------------- |
| morning          | "Good morning"   |
| afternoon        | "Good afternoon" |
| evening          | "Good evening"   |
| night            | "Good night"     |

## Argument Examples

### Single argument ($name)

```bash
/hello Alice
```

Response: **Hello Alice!** Nice to meet you.

### Multiple arguments ($name and $time-of-day)

```bash
/hello Alice morning
```

Response: **Good morning, Alice!** Hope you had a great start to your day.

### Using indexed arguments ($ARGUMENTS[0], $ARGUMENTS[1])

```bash
/hello Bob afternoon
```

Same as above but using indexed placeholders:

- `$ARGUMENTS[0]` → Bob
- `$ARGUMENTS[1]` → afternoon
- Shorthand: `$0` → Bob, `$1` → afternoon

### Full arguments string ($ARGUMENTS)

```bash
/hello Charlie
```

When $ARGUMENTS is used, it captures all input:

- `$ARGUMENTS` → Charlie

### Real-world example

```bash
/hello David evening
```

When $ARGUMENTS is used, it captures all input:

- `$ARGUMENTS` → Charlie

### Real-world example

```bash
/hello David evening
```

Output: **Good evening, David!** I hope your day is going well. What would you like to work on?
