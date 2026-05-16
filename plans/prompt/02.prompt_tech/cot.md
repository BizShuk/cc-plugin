https://www.promptingguide.ai/techniques/cot

## Summary of Chain-of-Thought (CoT) Prompting

Chain-of-Thought prompting enables complex reasoning by showing intermediate steps before final answers. It works by combining few-shot examples that demonstrate the reasoning process, not just the answer.

**Key points:**

1. **What it is:** Introduced in 2022, CoT prompting enables "complex reasoning capabilities through intermediate reasoning steps."

2. **How to use it:** Provide a few examples that show your reasoning process, then ask the model to answer a new question using similar reasoning.

3. **Zero-shot CoT:** Simply adding "Let's think step by step" to prompts can improve reasoning without any examples.

4. **Auto-CoT:** An automatic approach that clusters questions and generates reasoning chains using LLMs to eliminate manual example crafting.

5. **Important note:** The authors state this is "an emergent ability that arises with sufficiently large language models."

**Example shown:** Instead of just answering "True/False" for whether odd numbers sum to an even number, the prompt shows how to identify the odd numbers, add them, then determine the answer.