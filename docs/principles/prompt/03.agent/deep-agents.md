https://www.promptingguide.ai/agents/deep-agents

# Deep Agents

## Core Concept

Deep Agents are advanced AI systems designed to solve complex, multi-step problems through strategic planning, memory management, and intelligent delegation. They represent a shift from shallow agents that break down on long tasks to robust systems for deep research and agentic coding.

## Key Components

### 1. Planning

Deep Agents maintain structured task plans that can be updated, retried, and recovered from - essentially a "living to-do list" guiding the agent toward long-term goals. This differs from ad-hoc reasoning within single context windows.

### 2. Orchestrator & Sub-agent Architecture

Instead of relying on one massive agent with long context, deep agents use:

- An orchestrator managing specialized sub-agents
- Specialized agents for search, coding, knowledge base retrieval, analysis, verification, and writing
- Clean separation of concerns for efficient context management
- Claude Code exemplifies this approach for coding tasks

### 3. Context Retrieval and Agentic Search

- External memory storage (files, vectors, databases)
- Reference intermediate work without context overload
- Hybrid memory combining agentic search + semantic search
- References ReasoningBank and Agentic Context Engineering research

### 4. Context Engineering

- Explicit, detailed, intentional instructions
- Clear definitions for when to plan, use sub-agents, name files, collaborate with humans
- Structured outputs and system prompt optimization
- Tool definition optimization

### 5. Verification

- Critical for reliability and production readiness
- LLM-as-a-Judge for automated verification
- Addresses hallucination, sycophancy, prompt injection issues
- Systematic evaluation pipelines for building verifiers

## Examples and Applications

- Customer support systems for educational platforms
- Agentic RAG (Retrieval Augmented Generation) systems
- Deep research tasks
- Agentic coding workflows

## Related Techniques

The broader guide covers related patterns including ReAct and Reflexion as prompting techniques that can be incorporated into deep agent systems. These techniques enhance the reasoning and self-reflection capabilities of deep agents.