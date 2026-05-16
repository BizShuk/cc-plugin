# LLM Settings Guide

Source: https://www.promptingguide.ai/introduction/settings

## Key Parameters

### Temperature
Controls randomness in outputs. Lower values (0-0.3) produce deterministic, consistent responses ideal for factual QA. Higher values (0.7-1.0) generate creative, varied outputs for creative tasks.

**Examples:**
- Fact-based Q&A: temperature = 0
- Creative writing: temperature = 0.9

### Top P (Nucleus Sampling)
Controls which token pool the model draws from. Low Top P (0-0.5) restricts to high-probability tokens for focused answers. High Top P (0.9-1.0) includes more diverse, lower-probability tokens.

**Important:** Use either temperature OR Top P, not both.

### Max Tokens
Sets an upper bound on token generation. Prevents excessive responses and controls costs.

### Stop Sequences
Strings that halt generation when encountered. Useful for controlling response structure.

### Frequency Penalty
Penalizes repeated tokens proportionally to their frequency. Reduces word repetition but may affect coherence.

### Presence Penalty
Applies equal penalty to all repeated tokens. Prevents phrase repetition while encouraging topic diversity.

**Important:** Choose either frequency OR presence penalty, not both.

## Best Practices

1. Adjust temperature OR Top P (not both) - they control similar randomness
2. Adjust frequency OR presence penalty (not both) - they reduce repetition differently
3. Start with defaults and tune based on task requirements
4. Use Max Tokens and Stop Sequences for structural control