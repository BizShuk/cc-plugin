https://www.promptingguide.ai/agents/context-engineering-deep-dive

# Context Engineering Deep Dive: Building a Deep Research Agent

## Overview

Building effective AI agents requires substantial tuning of system prompts, tool definitions, agent architecture, and input/output specifications. This guide explores practical context engineering through developing a deep research agent.

## Agent Architecture Design

### The Original Design Problem

The initial architecture connected web search directly to the deep research agent, placing excessive burden on a single agent responsible for:
- Managing tasks
- Saving information to memory
- Executing web searches
- Generating final reports

This caused context to grow too long, agents forgetting to execute searches, missed task updates, and unreliable behavior.

### The Improved Multi-Agent Architecture

The solution separated concerns by introducing a dedicated search worker agent with the following benefits:
- Separation of concerns
- Improved reliability through focused responsibilities
- Model selection flexibility

**Architecture Pattern:**
- Parent agent handles planning and orchestration
- Search worker focuses exclusively on web searches

## System Prompt Engineering

### Key Components

**1. High-Level Agent Definition**
Clear role definition establishes the agent's purpose for planning and executing search tasks to generate deep research reports.

**2. General Instructions**
Explicit workflow instructions specify:
- User provides a query
- Query gets converted into a search plan with multiple tasks
- Executed searches maintained in a spreadsheet
- Final report generated

**3. Essential Context**
Including current date information is crucial for research agents. Without temporal context, agents often search for outdated information since LLMs typically have knowledge cutoffs behind the current date.

**4. Tool Definitions**
The biggest performance improvements often come from clearly explaining tool usage in the system prompt, not just defining tool parameters. Detailed instructions should specify when and how to use each tool.

## Key Techniques

### Explicit Status Definitions

Without explicit status definitions, agents use inconsistent values like "pending" vs "to-do" or "completed" vs "done". Be explicit about allowed values to eliminate ambiguity.

**Example:**
```json
{
  "status": ["pending", "in_progress", "completed", "failed"]
}
```

### Flexible vs Rigid Approaches

**Flexible:** Instruct agents to use tools in the order that makes most sense.

**Rigid:** For production consistency, use instructions requiring all tasks to be executed without skipping.

## Context Engineering Iteration Process

Context engineering requires iterative improvement:
1. Test with diverse queries
2. Identify issues (missed tasks, wrong status values)
3. Add specific instructions to address problems
4. Repeat the cycle

Even after multiple iterations, opportunities remain for:
- Search task metadata augmentation
- Enhanced search planning
- Date range specification

## Advanced Considerations

### Sub-Agent Communication

Keep sub-agent inputs minimal and focused. The search worker needs only the search query text, not full context or task metadata.

### Context Length Management

Strategies include:
- Using separate agents to isolate context
- Implementing memory management tools
- Summarizing long outputs
- Clearing task lists between queries

### Error Handling

System prompts should include failure scenario instructions:
- Retries on failure
- Mark status as "failed" with reasons
- Notify users of critical errors
- Never proceed silently when operations fail

## Key Takeaways

Context engineering transforms unreliable prototypes into robust, production-ready systems through:

- Clear role definitions
- Explicit tool instructions
- Essential context provision
- Iterative improvement

**Success requires:**
1. Significant iteration time
2. Careful architectural decisions
3. Explicit instructions eliminating assumptions
4. Continuous refinement based on observed behavior
5. Balancing flexibility with control