https://www.promptingguide.ai/techniques/activeprompt

## Summary: Active-Prompt

**What is Active-Prompt:**
Active-Prompt is a prompting technique proposed by Diao et al. (2023) that adapts large language models to task-specific example prompts with human-designed chain-of-thought (CoT) reasoning.

**The Problem It Solves:**
Traditional CoT methods rely on fixed human-annotated exemplars. These fixed examples may not be most effective for different tasks.

**How It Works:**
1. Query the LLM with or without few CoT examples
2. Generate k possible answers for training questions
3. Calculate uncertainty metric using answer disagreement
4. Select most uncertain questions for human annotation
5. Use newly annotated exemplars to infer answers

**Key Benefit:**
By dynamically selecting and annotating the most uncertain examples, Active-Prompt tailors exemplars to specific tasks rather than relying on one-size-fits-all human annotations.