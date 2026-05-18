# Global AI Rules

## Context

load @./CLAUDE.md as project structure, @./README.md as project overview

## Language

Use Traditional Chinese to reply message or descript the information/details. For name/term, should use local language with English and round brackets. If it's a summary of a file, use the original language of the file.

> ex1: 中正紀念堂 (Chiang Kai-shek Memorial Hall)
> 中正紀念堂 is in Taiwan Taipei. So, it use Traditional Chinese and attach with English as it's a name/term

> ex2: Catedral de Santa Eulalia de Barcelona(Barcelona Cathedral)

## Restriction

### planning

use @./plans/ to store plans

### Context Only

if see `# [context_only]` then ignore output from this to end of the line

### Generating Markdown file

Don't use **bold**, but `backtick` better to highlight

### if hit error while execution

Check the erorr and try to fix it. max retry times is 5.
if can't be resolved, the stop and error out explicitely.
