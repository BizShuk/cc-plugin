# Prompt Engineering Examples

Source: https://www.promptingguide.ai/introduction/examples

## Key Points Summary

This guide covers practical prompt engineering techniques across multiple domains, emphasizing specificity in instructions, strategic use of examples, and proper prompt structure.

---

## 1. Text Summarization

**Technique**: Use instruction-based prompts to condense information.

**Example**:
```
Prompt: "Explain the above in one sentence"
```

**Key Insight**: Direct instructions to distill lengthy text into concise summaries.

---

## 2. Information Extraction

**Technique**: Prompt models to pull specific details from given context.

**Example**:
Extracting "ChatGPT" from a paragraph about AI use in research papers by asking targeted questions.

**Key Insight**: Clear, specific questions yield precise extractions.

---

## 3. Question Answering

**Technique**: Combine context, instructions, and output indicators in a structured prompt.

**Example**:
A passage about Teplizumab combined with a targeted question and an instruction like "Answer in one word."

**Key Insight**: Structure prompts with context + question + format instruction for reliable answers.

---

## 4. Text Classification

**Technique**: Provide explicit instructions and examples to improve classification consistency.

**Example**:
Specifying desired output format (e.g., "neutral" instead of "Neutral") to ensure consistency.

**Key Insight**: Ambiguous labels lead to inconsistent outputs; explicit formatting instructions resolve this.

---

## 5. Role Prompting / Conversation

**Technique**: Instruct the model on identity and tone to shape the response style.

**Example**:
- Technical-scientific response: "Explain quantum computing like a physics professor"
- Beginner-friendly response: "Explain it like you're talking to a 10-year-old"

**Key Insight**: Assigning a role influences vocabulary, depth, and communication style.

---

## 6. Code Generation

**Technique**: Describe the desired outcome and let the model generate the appropriate code.

**Example**:
- JavaScript: "Write a JavaScript function that asks for a user's name and displays a greeting"
- SQL: "Given a database schema with tables X, Y, Z, write a MySQL query to..."

**Key Insight**: Clear descriptions of input/output requirements produce accurate code.

---

## 7. Reasoning (Chain-of-Thought)

**Technique**: Break complex problems into stages rather than asking for direct answers.

**Example**:
Instead of: "What is 15 * 17?"
Use: "Solve step by step. First multiply 15 * 10, then..."

**Key Insight**: Complex tasks require step-by-step reasoning prompts for better results than simple direct questions.

---

## General Principles

1. **Specificity**: Clear, detailed instructions outperform vague ones
2. **Structure**: Combine context + instruction + format indicators
3. **Examples**: Strategic use of examples improves consistency
4. **Role Assignment**: Defining persona shapes response style
5. **Incremental Decomposition**: Break complex tasks into stages