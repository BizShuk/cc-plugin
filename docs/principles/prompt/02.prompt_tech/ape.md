https://www.promptingguide.ai/techniques/ape

## Summary: Automatic Prompt Engineer (APE)

**What is APE:**
APE is a framework for automatic instruction generation and selection proposed by Zhou et al. (2022). It frames instruction generation as natural language synthesis addressed as black-box optimization using LLMs to generate and search over candidate solutions.

**How to Use:**
1. Provide a LLM with output demonstrations for a specific task
2. Generate multiple instruction candidates from those demonstrations
3. Execute the candidates using your target model
4. Select the best instruction based on evaluation scores

**Example:**
APE discovered a better zero-shot Chain-of-Thought prompt than the human-engineered "Let's think step by step." The discovered prompt: "Let's work this out in a step by step way to be sure we have the right answer." This elicits chain-of-thought reasoning and improved performance on MultiArith and GSM8K benchmarks.

**Related Techniques:**
- Prompt-OIRL: Uses offline inverse reinforcement learning for query-dependent prompts
- OPRO: Uses LLMs to optimize prompts (e.g., "Take a deep breath" improves math performance)
- AutoPrompt: Gradient-guided search for automatic prompt creation
- Prefix Tuning and Prompt Tuning: Learnable continuous prompts via backpropagation