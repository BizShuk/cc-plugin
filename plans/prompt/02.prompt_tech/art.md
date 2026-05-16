https://www.promptingguide.ai/techniques/art

## ART (Automatic Reasoning and Tool-use) Summary

**What is ART?**
ART combines Chain-of-Thought prompting with tool use in an interleaved manner. Proposed by Paranjape et al. (2023), it uses a frozen LLM to automatically generate intermediate reasoning steps as a program.

**How It Works:**
1. Given a new task, it selects demonstrations of multi-step reasoning and tool use from a task library
2. At test time, it pauses generation whenever external tools are called
3. Integrates tool output before resuming generation
4. Generalizes from demonstrations to decompose new tasks in a zero-shot fashion
5. Extensible—humans can fix mistakes or add new tools by updating the task and tool libraries

**Performance:**
- Substantially improves over few-shot prompting and automatic CoT on unseen tasks
- Outperforms on BigBench and MMLU benchmarks
- Exceeds hand-crafted CoT prompts when human feedback is incorporated

**Key Advantage:** Unlike approaches requiring hand-crafted task-specific demonstrations and carefully scripted interleaving, ART automatically generates intermediate reasoning steps while remaining extensible through simple library updates.