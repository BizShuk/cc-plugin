https://www.promptingguide.ai/techniques/tot

## Summary: Tree of Thoughts (ToT)

**What is ToT:**
Tree of Thoughts is a framework that generalizes chain-of-thought prompting by maintaining a tree structure where "thoughts represent coherent language sequences that serve as intermediate steps toward solving a problem." It enables language models to self-evaluate progress and explore solutions using search algorithms like breadth-first and depth-first search.

**How to Use It:**
To implement ToT, you must define:
- Number of candidates to explore
- Number of thought steps

For Game of 24, for example, the framework uses 3 steps with b=5 candidates kept per step. The model evaluates each thought as "sure/maybe/impossible" with regard to reaching the goal.

**Example Prompt (Hulbert's simplified version):**
> "Imagine three different experts are answering this question. All experts will write down 1 step of their thinking, then share it with the group."

This approach outperformed traditional prompting methods significantly in the original research.

**Key Distinctions:**
- Yao et al. use generic DFS/BFS/beam search
- Long's version uses a reinforcement learning-trained "ToT Controller" that can evolve and learn

**Related Concepts:**
PanelGPT extends this by using panel discussions among multiple LLMs to evaluate thoughts collaboratively.