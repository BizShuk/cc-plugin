---
name: markdownlint
description: >
    Use when generating or modifying Markdown files under plugins/experiment/skills directory. Triggers on Markdown file creation or modification in the experiment skills workspace.
version: "1.0.0"
allowed-tools: Read, write_file
user-invocable: false
disable-model-invocation: false
effort: low
context: fork
metadata:
    type: reference
    platforms: [macos, linux]
---

# Markdownlint

## Overview

This skill provides a subset of `markdownlint` (v0.40.0) formatting rules and repository-specific conventions to ensure consistent, clean, and valid Markdown files inside the codebase.

## When to Use

- Triggered when creating new Markdown files (`.md` extension) under `plugins/experiment/skills/`.
- Triggered when editing or updating existing Markdown files.

When NOT to use:

- Do not apply for files other than Markdown format.

## Usage

Use `npx markdownlint` to check and format Markdown files. This avoids installing `markdownlint-cli` globally.

If the command is not found, install `markdownlint-cli` locally first:

```bash
# Check a specific file
npx markdownlint path/to/file.md

# Check all Markdown files in a directory
npx markdownlint plugins/experiment/skills/**/*.md

# Automatically fix basic issues
npx markdownlint --fix path/to/file.md
```

## Quick Reference Table

The following table summarizes core rules from `markdownlint` (v0.40.0) and custom constraints:

| Rule ID     | Alias                              | Key Requirement                                                                    | Status |
| :---------- | :--------------------------------- | :--------------------------------------------------------------------------------- | :----- |
| `MD001`     | `heading-increment`                | Heading levels must increment by exactly one level at a time.                      | `on`   |
| `MD003`     | `heading-style`                    | Consistent heading style (always use ATX style `#`).                               | `on`   |
| `MD004`     | `ul-style`                         | Unordered list style                                                               | `on`   |
| `MD005`     | `list-indent`                      | Consistent indentation for list items at the same level.                           | `on`   |
| `MD007`     | `ul-indent`                        | Unordered list indentation (2 spaces for sublists).                                | `on`   |
| `MD009`     | `no-trailing-spaces`               | Zero trailing whitespace at the end of lines.                                      | `on`   |
| `MD010`     | `no-hard-tabs`                     | Indentation must use spaces, not hard tabs.                                        | `on`   |
| `MD011`     | `no-reversed-links`                | Reversed link syntax                                                               | `on`   |
| `MD012`     | `no-multiple-blanks`               | Maximum of 1 consecutive empty blank line.                                         | `on`   |
| `MD013`     | `line-length`                      | Line length                                                                        | `off`  |
| `MD014`     | `commands-show-output`             | Dollar signs used before commands without showing output                           | `on`   |
| `MD018`     | `no-missing-space-atx`             | No space after hash on atx style heading                                           | `on`   |
| `MD019`     | `no-multiple-space-atx`            | Multiple spaces after hash on atx style heading                                    | `on`   |
| `MD020`     | `no-missing-space-closed-atx`      | No space inside hashes on closed atx style heading                                 | `on`   |
| `MD021`     | `no-multiple-space-closed-atx`     | Multiple spaces inside hashes on closed atx style heading                          | `on`   |
| `MD022`     | `blanks-around-headings`           | Headings must be surrounded by blank lines (1 line above, 1 line below).           | `on`   |
| `MD023`     | `heading-start-left`               | Headings must start at the beginning of the line (no indentation).                 | `on`   |
| `MD024`     | `no-duplicate-heading`             | Multiple headings with the same content                                            | `on`   |
| `MD025`     | `single-title`                     | Exactly one top-level H1 heading (`# Title`) at the beginning.                     | `on`   |
| `MD026`     | `no-trailing-punctuation`          | Trailing punctuation in heading                                                    | `on`   |
| `MD027`     | `no-multiple-space-blockquote`     | Multiple spaces after blockquote symbol                                            | `on`   |
| `MD028`     | `no-blanks-blockquote`             | Blank line inside blockquote                                                       | `on`   |
| `MD029`     | `ol-prefix`                        | Ordered list item prefix                                                           | `on`   |
| `MD030`     | `list-marker-space`                | Spaces after list markers                                                          | `on`   |
| `MD031`     | `blanks-around-fences`             | Code blocks must be surrounded by blank lines.                                     | `on`   |
| `MD032`     | `blanks-around-lists`              | Lists must be surrounded by blank lines.                                           | `on`   |
| `MD033`     | `no-inline-html`                   | Inline HTML                                                                        | `on`   |
| `MD034`     | `no-bare-urls`                     | Bare URL used                                                                      | `on`   |
| `MD035`     | `hr-style`                         | Horizontal rule style                                                              | `on`   |
| `MD036`     | `no-emphasis-as-heading`           | Emphasis used instead of a heading                                                 | `on`   |
| `MD037`     | `no-space-in-emphasis`             | Spaces inside emphasis markers                                                     | `on`   |
| `MD038`     | `no-space-in-code`                 | Spaces inside code span elements                                                   | `on`   |
| `MD039`     | `no-space-in-links`                | Spaces inside link text                                                            | `on`   |
| `MD040`     | `fenced-code-language`             | Fenced code blocks must have a language specified (e.g., `bash`, `go`, `text`).    | `on`   |
| `MD041`     | `first-line-heading`               | First line in a file must be a top-level heading.                                  | `on`   |
| `MD042`     | `no-empty-links`                   | No empty links                                                                     | `on`   |
| `MD043`     | `required-headings`                | Required heading structure                                                         | `on`   |
| `MD044`     | `proper-names`                     | Proper names should have the correct capitalization                                | `on`   |
| `MD045`     | `no-alt-text`                      | Images should have alternate text (alt text)                                       | `on`   |
| `MD046`     | `code-block-style`                 | Code block style                                                                   | `on`   |
| `MD047`     | `single-trailing-newline`          | Files must end with a single trailing newline.                                     | `on`   |
| `MD048`     | `code-fence-style`                 | Code fence style                                                                   | `on`   |
| `MD049`     | `emphasis-style`                   | Emphasis style                                                                     | `on`   |
| `MD050`     | `strong-style`                     | Strong style                                                                       | `on`   |
| `MD051`     | `link-fragments`                   | Link fragments should be valid                                                     | `on`   |
| `MD052`     | `reference-links-images`           | Reference links and images should use a label that is defined                      | `on`   |
| `MD053`     | `link-image-reference-definitions` | Link and image reference definitions should be needed                              | `on`   |
| `MD054`     | `link-image-style`                 | Link and image style                                                               | `on`   |
| `MD055`     | `table-pipe-style`                 | Table pipe style                                                                   | `on`   |
| `MD056`     | `table-column-count`               | Table column count                                                                 | `on`   |
| `MD058`     | `blanks-around-tables`             | Tables should be surrounded by blank lines                                         | `on`   |
| `MD059`     | `descriptive-link-text`            | Link text should be descriptive                                                    | `on`   |
| `MD060`     | `table-column-style`               | Table column style                                                                 | `on`   |
| `CUSTOM-01` | `no-bold-formatting`               | Never use double asterisk `**` for bold formatting. Use backticks `` ` `` instead. | `on`   |
| `CUSTOM-02` | `relative-local-links`             | Local file links must use relative paths, not absolute paths or `file://` URLs.    | `on`   |

## Core Pattern

Here is a comparison demonstrating how to correctly format a Markdown file.

### Incorrect Formatting

```markdown
# My Skill

Here is a paragraph with trailing space.
We skipped level 3 heading.

#### Incorrect Subheading

We also used **bold formatting** which is prohibited.

- Unordered list
    - Misaligned indentation (3 spaces)
```

### Correct Formatting

```markdown
# My Skill

Here is a paragraph without trailing space.
We increment heading properly.

## Correct Subheading

We also used `backtick formatting` to highlight text.

- Unordered list
    - Correct indentation (2 spaces)
```

## Common Mistakes

- Skipping heading levels: Going from `#` straight to `###` is invalid. Always go `#` -> `##` -> `###`.
- Trailing spaces: Unintentional space at the end of lines. Always trim trailing spaces.
- Missing blank lines: Forgetting blank lines before/after lists, code blocks, or headings.
- Bold text: Using `**bold**` instead of `backticks` for emphasizing terms.
- Absolute local links: Use relative paths for local files, such as `[README](../../README.md)`, instead of `/Users/...` or `file://` URLs.
