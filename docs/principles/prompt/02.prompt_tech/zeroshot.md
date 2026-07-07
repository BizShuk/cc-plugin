https://www.promptingguide.ai/techniques/zeroshot

## Summary: Zero-Shot Prompting

**What It Is:**
Zero-shot prompting is a technique where you give an LLM instructions to perform a task without providing examples or demonstrations. The model relies on its training to understand and execute the request.

**How to Use It:**
Simply instruct the model directly on what you want it to do, without any examples included in your prompt.

**Example:**
The page provides this illustration:
> "Classify the text into neutral, negative or positive. Text: I think the vacation is okay. Sentiment:"
>
> Output: Neutral

The model understands "sentiment" without examples—that's the zero-shot capability in action.

**Background:**
- Models like GPT-3.5 Turbo, GPT-4, and Claude 3 are trained on large amounts of data, making them capable of zero-shot performance
- Instruction tuning (Wei et al., 2022) and RLHF (reinforcement learning from human feedback) improve this capability

**When It Doesn't Work:**
If zero-shot prompting fails, the recommendation is to add examples to your prompt, which leads to few-shot prompting—a technique covered in the guide's next section.