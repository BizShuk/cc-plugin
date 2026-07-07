# Prompt Engineering Tips

Source: <https://www.promptingguide.ai/introduction/tips>

## Key Points Summary

### 1. Start Simple

Treat prompt design as an iterative process. Begin with simple prompts and gradually add complexity as needed. Break larger tasks into smaller subtasks to avoid overwhelming the process.

### 2. Use Clear Instructions

Commands like "Write," "Classify," "Summarize," "Translate," and "Order" help direct the model effectively. Place instructions at the beginning of prompts and use separators (like ###) to distinguish instructions from context.

**Example:**

```
### Instruction ###
Translate the text below to Spanish:
Text: 'hello!'
```

Output: "¡Hola!"

### 3. Be Specific

The more descriptive and detailed your prompt, the better the results. Focus on relevant details that contribute to the task. Providing examples within prompts is highly effective for achieving desired output formats.

### 4. Avoid Impreciseness

Be direct rather than overly clever. Instead of vague instructions like "keep it short," specify exact parameters such as "Use 2-3 sentences."

### 5. Frame Positively

Tell the model what TO do rather than what NOT to do. This promotes clearer, more actionable responses.

**Example:** Rather than saying "DO NOT ASK FOR INTERESTS," specify what the agent SHOULD do instead.
