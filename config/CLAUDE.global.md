# Global AI Rules

If any conflict and can't determine, AskUserQuestions and wait for the answer.

## Output Style

- answer conclusion first
- concise response
- if any relations/associations, try to below
    - Indented List (tree graph)
    - Minimalist Relationship Expression
    - Mermaid
    - Markdown Table

## Context

Load @./CLAUDE.md as project structure, if any updates from structure, need to update too
Load @./README.md as project overview, if any business scope change, need to update too

## Language

Use Traditional Chinese to reply message or describe the information/details. For name/term, should use local language with English and round brackets. If it's a summary of a file, use the original language of the file.

> ex1: 中正紀念堂 (Chiang Kai-shek Memorial Hall)
> 中正紀念堂 is in Taiwan Taipei. So, it use Traditional Chinese and attach with English as it's a name/term

> ex2: Catedral de Santa Eulalia de Barcelona(Barcelona Cathedral)

## Restriction

### Planning

use @./plans/ to store plans

### Context Only

if see `# [context_only]` then ignore output from this to end of the line

### Generating Markdown file

Don't use **bold**, but `backtick` better to highlight

### if hit error while execution

Check the error and try to fix it. max retry times is 5.
if can't be resolved, then stop and error out explicitly.

## Convention

### Command Line

- `monitor` sub command is used to overall monitoring for the command
