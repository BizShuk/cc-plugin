https://www.promptingguide.ai/techniques/pal

## Summary: PAL (Program-Aided Language Models)

**What is PAL?**
PAL is a prompting technique where an LLM generates **program code** as intermediate reasoning steps instead of free-form text. The generated code is then executed by a programmatic runtime (like Python) to produce the final answer.

**How to Use It:**
1. Provide exemplars showing how to translate questions into code
2. Let the LLM generate Python code based on your question
3. Execute the generated code to get the answer

**Key Example:**
For the question "Today is 27 February 2023. I was born exactly 25 years ago. What is the date I was born in MM/DD/YYYY?"

The LLM generates:
```python
today = datetime(2023, 2, 27)
born = today - relativedelta(years=25)
born.strftime('%m/%d/%Y')
```
After executing: **`02/27/1998`**

**Why PAL Works:**
By offloading calculations to a Python interpreter, PAL avoids arithmetic errors that LLMs commonly make. The LLM handles reasoning and code generation, while Python handles precise computation.