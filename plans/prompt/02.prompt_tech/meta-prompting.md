https://www.promptingguide.ai/techniques/meta-prompting

## Summary: Meta-Prompting Technique

**What It Is:**
Meta Prompting is an advanced prompting technique that emphasizes the structural and syntactical patterns of problems rather than their specific content details. It creates a more abstract, structured approach to interacting with LLMs.

**Key Characteristics:**
- Structure-oriented: prioritizes format and pattern over content
- Syntax-focused: uses syntax as a guiding template
- Uses abstract examples as frameworks
- Draws from type theory for categorizing prompt components
- Applicable across various domains

**How to Use It:**
Instead of providing detailed content examples (like few-shot prompting), meta prompting focuses on:
- The format and structure of problems
- The logical arrangement of components
- Abstract frameworks that illustrate patterns

**Advantages Over Few-Shot Prompting:**
- Requires fewer tokens (more efficient)
- Provides fair comparison between models
- Functions similarly to zero-shot prompting
- Minimizes influence of specific examples

**Applications:**
- Complex reasoning tasks
- Mathematical problem-solving
- Coding challenges
- Theoretical queries

**Important Consideration:**
Meta prompting assumes the LLM has inherent knowledge about the task. Performance may decrease with highly unique or novel tasks, similar to zero-shot prompting limitations.

**Example Context:**
The page references an example from Zhang et al. (2024) comparing meta prompting to few-shot prompting on the MATH benchmark (visible in an accompanying image showing structural differences between the two approaches).