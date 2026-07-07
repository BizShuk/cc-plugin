https://www.promptingguide.ai/techniques/consistency

## Summary: Self-Consistency Prompting Technique

### What It Is

Self-Consistency is an advanced prompting technique proposed by Wang et al. (2022) that aims to improve reasoning quality by sampling multiple diverse reasoning paths and selecting the most consistent answer.

### How It Works

The technique:
- Replaces greedy decoding used in chain-of-thought prompting
- Samples multiple reasoning paths using few-shot CoT
- Uses the generations to identify the most common/frequent answer
- Boosts performance on arithmetic and commonsense reasoning tasks

### Example

A flawed single-output approach gives an incorrect answer:
> "When I was 6 my sister was half my age. Now I'm 70 how old is my sister?" → **35** (wrong)

With self-consistency, multiple reasoning paths are generated, producing different answers. The majority/consistent answer becomes the final result:
- Output 1: 67
- Output 2: 67
- Output 3: 35

The answer "67" appears most frequently, so it is selected as the final answer.