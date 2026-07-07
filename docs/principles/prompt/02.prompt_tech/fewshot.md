https://www.promptingguide.ai/techniques/fewshot

## Summary: Few-Shot Prompting

**What It Is:**
Few-shot prompting enables in-context learning by providing examples (demonstrations) in your prompt to guide the model's output. It helps when zero-shot capabilities aren't enough.

**How to Use It:**
- Include 1-10 example demonstrations before your actual query
- The format matters more than correct labels—even random labels improve performance
- Matching your examples to the true label distribution helps
- Keep the format consistent across examples

**Key Examples:**
1. Teaching new words: "A 'whatpu' is a small, furry animal..." provides one example, then asks the model to use a new word in a sentence.
2. Sentiment classification with scrambled labels still produces correct results.

**Limitation:**
Few-shot prompting struggles with complex reasoning tasks, like multi-step math problems. When demonstrations don't help, consider chain-of-thought prompting or fine-tuning.