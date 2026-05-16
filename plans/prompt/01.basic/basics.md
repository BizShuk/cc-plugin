# Basics of Prompting

Source: https://www.promptingguide.ai/introduction/basics

## Key Points

### Prompt Structure
A prompt can include an instruction, question, context, inputs, or examples. The more clearly you state what you want, the better the model's output.

### Prompt Formatting
- **Standard form**: `<Instruction>` or `<Question>?`
- **QA style**:
  ```
  Q: <Question>?
  A:
  ```

### Zero-Shot Prompting
Direct requests without any examples. The model uses its existing knowledge to respond.

Example:
```
User: What is prompt engineering?
Model: Prompt engineering is the process of designing and optimizing prompts...
```

### Few-Shot Prompting
Providing a few demonstrations so the model can learn in-context and follow the same pattern.

Example:
```
This is awesome! // Positive
This is bad! // Negative
Wow that movie was rad! // Positive
What a horrible show! // Negative
```

## Practical Examples

**Basic prompting:**
- Bad: `"The sky is"` → `"blue."`
- Good: `"Complete the sentence: The sky is"` → `"blue during the day and dark at night."`

**Using clear instructions:**
- Include specific instructions about the format, tone, or structure you want
- Add context when helpful for better responses

**Few-shot for style matching:**
```
Input: This is amazing! // Positive
Input: Terrible experience // Negative
Input: I love this product! //
Model output: Positive
```

## Summary
Clear, well-structured prompts produce better results. Use instructions for direction, context for relevance, and examples (few-shot) to guide output format and style when zero-shot is insufficient.