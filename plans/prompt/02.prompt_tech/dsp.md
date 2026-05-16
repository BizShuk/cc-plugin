https://www.promptingguide.ai/techniques/dsp

## Summary: Directional Stimulus Prompting (DSP)

**What It Is:**
DSP is a prompting technique proposed by Li et al. (2023) that better guides LLMs to generate desired outputs, particularly summaries.

**How It Works:**
A tuneable policy language model generates "stimulus/hint" signals that guide a larger, frozen black-box LLM. The technique leverages reinforcement learning to optimize the guidance process.

**Key Characteristics:**
- The policy LM can be small and optimized specifically for generating effective hints
- Uses RL to optimize LLM outputs
- Compares favorably against standard prompting approaches

**Current Status:**
The page indicates that full examples are "coming soon," so no practical implementation examples are currently available in this guide.

**Related Context:**
The technique is categorized under advanced prompting techniques alongside methods like Chain-of-Thought, Active-Prompt, and Program-Aided Language Models.