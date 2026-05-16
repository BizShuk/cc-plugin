https://www.promptingguide.ai/techniques/reflexion

## What is Reflexion?

Reflexion is a framework that helps language model agents learn from mistakes through verbal feedback. Instead of using traditional reinforcement learning with numeric rewards, it provides linguistic self-reflection that agents use to improve future decisions.

## How It Works

The framework uses three components:
1. **Actor** – Generates actions and receives observations from the environment
2. **Evaluator** – Scores the actor's outputs with reward signals
3. **Self-Reflection** – Creates verbal feedback based on past mistakes and stores experiences in memory

The cycle repeats: generate trajectory → evaluate → reflect → improve.

## Key Benefits
- Agents learn through trial and error without fine-tuning
- Provides nuanced, specific feedback rather than simple scores
- Creates interpretable memory of past experiences
- Works well for decision-making, reasoning, and coding tasks

## Example Applications
- **AlfWorld tasks**: Completed 130/134 decision-making challenges
- **HotPotQA**: Improved reasoning over multiple documents
- **Programming**: Achieved state-of-the-art results on HumanEval and MBPP benchmarks

## Limitations
- Depends on the agent's ability to accurately self-evaluate
- Memory constraints with large tasks
- Challenges with non-deterministic code functions

The framework is particularly useful when traditional reinforcement learning is impractical but nuanced feedback and explicit memory are needed for improvement.