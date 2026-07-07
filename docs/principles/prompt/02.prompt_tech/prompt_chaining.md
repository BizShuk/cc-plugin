https://www.promptingguide.ai/techniques/prompt_chaining

## Summary: Prompt Chaining

**What It Is:**
Prompt chaining breaks complex tasks into subtasks, using each LLM response as input for the next prompt in sequence.

**Key Benefits:**
- Improves reliability and performance on complex tasks
- Increases transparency and controllability
- Makes debugging easier by allowing analysis at each stage
- Enhances personalization for conversational applications

**How to Use It:**
Split a task into logical subtasks, then chain prompts so the output of one feeds into the next.

**Example - Document QA:**
1. **Prompt 1** extracts relevant quotes from a document based on the user's question
2. **Prompt 2** takes those quotes plus the original document and generates a helpful, accurate response

This two-step approach outperforms asking a single detailed prompt to answer the question directly.

**Best For:**
Complex tasks where responses need multiple transformations, building conversational assistants, and improving user experience through staged processing.