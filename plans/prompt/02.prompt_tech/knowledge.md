https://www.promptingguide.ai/techniques/knowledge

## Summary: Generated Knowledge Prompting

**What It Is:**
Generated Knowledge Prompting is a technique where an LLM first generates relevant knowledge before making a prediction. Introduced by Liu et al. in 2022, this approach helps improve accuracy on commonsense reasoning tasks.

**How to Use It:**
1. Generate multiple knowledge statements by providing example input-knowledge pairs, then asking for knowledge related to your target question
2. Integrate the generated knowledge with your question in a QA format
3. Request an explanation and answer based on the provided knowledge

**Key Example:**
A prompt asking "Part of golf is trying to get a higher point total than others. Yes or No?" initially produced an incorrect "Yes" response. When the model was first asked to generate knowledge about golf, it produced accurate facts explaining that golf involves minimizing strokes, not maximizing points. Providing this knowledge in the follow-up prompt led to correct answers—though confidence levels varied between different generated knowledge statements.

This technique addresses LLM limitations on knowledge-intensive tasks by having the model access "generated" information rather than relying solely on internal parameters.