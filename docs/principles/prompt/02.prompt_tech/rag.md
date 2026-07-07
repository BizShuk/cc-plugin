https://www.promptingguide.ai/techniques/rag

## Summary: Retrieval-Augmented Generation (RAG)

**What is RAG?**

RAG is a technique introduced by Meta AI researchers that combines information retrieval with text generation to address knowledge-intensive tasks. It helps language models access external knowledge sources, improving factual consistency and reducing hallucination.

**How It Works:**

The system takes an input and retrieves relevant documents from a knowledge source (like Wikipedia). These documents are concatenated with the original prompt and fed to a text generator. This approach uses both "parametric memory" (a seq2seq model) and "non-parametric memory" (dense vector index) to produce outputs.

**Key Benefits:**

RAG enables models to access up-to-date information without retraining, making it valuable when facts evolve over time. It performs well on benchmarks like Natural Questions, WebQuestions, and CuratedTrec, and generates more factual, specific, and diverse responses.

**Use Case Example:**

The page references a notebook tutorial for building a RAG system to generate short, concise machine learning paper titles using open-source LLMs.