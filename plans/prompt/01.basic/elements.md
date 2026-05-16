# Elements of a Prompt

Source: https://www.promptingguide.ai/introduction/elements

## Summary

Effective prompts consist of four key elements:

### 1. Instruction
A specific task or action you want the model to perform. This is the core directive that guides the model's response.

**Example:**
```
Translate the following English text to French
```

### 2. Context
External information or additional background that helps steer the model toward better, more relevant responses. This can include examples or domain-specific knowledge.

**Example:**
```
You are a medical billing specialist. Translate the following...
```

### 3. Input Data
The specific question, text, or data you want the model to process or analyze.

**Example:**
```
Text: The patient was seen for a routine check-up on May 15th.
```

### 4. Output Indicator
The expected type or format of the response (e.g., labeling, a specific structure).

**Example:**
```
Sentiment: [positive/negative/neutral]
```

## Combined Example

A complete prompt demonstrating all four elements:

```
Classify the text into neutral, negative, or positive.

Text: I think the pizza was okay.

Sentiment:
```

In this example:
- **Instruction**: "Classify the text into neutral, negative, or positive"
- **Input Data**: "I think the pizza was okay."
- **Output Indicator**: "Sentiment:"

## Notes

- Not all four elements are required for every prompt
- The format depends on your specific task
- Context can include examples that help the model better understand expected output
- More context generally leads to better, more accurate responses