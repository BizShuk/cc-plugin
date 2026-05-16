https://www.promptingguide.ai/techniques/react

## Summary: ReAct (Reasoning + Acting) Prompting

**What is ReAct?**
ReAct is a framework where LLMs generate both reasoning traces and task-specific actions in an interleaved manner. It combines "acting" and "reasoning" to allow models to learn new tasks and make decisions.

**Key Innovation:**
Unlike chain-of-thought prompting, which lacks access to external information, ReAct interfaces with external sources like knowledge bases to incorporate additional information into reasoning. This helps reduce fact hallucination.

**How It Works:**
The model follows a Thought → Action → Observation cycle:
- **Thought**: Reasoning about what to do next
- **Action**: Performing an action (e.g., search, lookup)
- **Observation**: Receiving feedback from the environment

**Prompt Structure:**
ReAct uses few-shot exemplars with multiple thought-action-observation steps. Example: "Question: What is the elevation range for the area that the eastern sector of the Colorado orogeny extends into?"

**Results:**
ReAct outperforms several state-of-the-art baselines on knowledge-intensive tasks (HotPotQA, Fever) and decision-making tasks (ALFWorld, WebShop). The best approach combines ReAct with chain-of-thought for both internal knowledge and external information.

**Practical Implementation:**
LangChain provides built-in support for ReAct agents using tools like Google search and calculators. Example query: "Who is Olivia Wilde's boyfriend? What is his current age raised to the 0.23 power?"