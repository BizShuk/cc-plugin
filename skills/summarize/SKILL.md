---
name: summarize
description: >
    Use when onboarding to a new project, after major refactors, or when README.md
    and CLAUDE.md are missing, outdated, or empty. Explores the entire workspace
    and writes structured documentation into two files at the repo root.
allowed-tools: Read, Bash, Glob, Write, Edit
disable-model-invocation: false
user-invocable: true
effort: high
context: fork
---

# summarize

Explore and conclude the whole workspace, then write structured documentation
into `README.md` and `CLAUDE.md` at the repository root.

## Overview

This skill performs a full-workspace scan and produces two canonical docs:

| File        | Purpose                                                                              |
| ----------- | ------------------------------------------------------------------------------------ |
| `README.md` | 商業邏輯 (Business logic): what the app does, how to use it, improvement suggestions |
| `CLAUDE.md` | 技術脈絡 (Technical context): project structure, key decisions, how to build/run     |

## When to Use

- Project has no `README.md` or `CLAUDE.md`
- Existing docs are stale after significant changes
- Onboarding to an unfamiliar codebase
- After a major refactor or migration
- User explicitly asks to summarize or document the workspace

## Procedure

### Step 1 — Discover project layout

```bash
# Get top-level structure
find . -maxdepth 1 -not -path './.git' -not -name '.' | sort

# Get deeper structure (exclude noise)
find . -maxdepth 3 -type f \
  -not -path '*/.git/*' \
  -not -path '*/node_modules/*' \
  -not -path '*/vendor/*' \
  -not -path '*/__pycache__/*' \
  -not -path '*/.venv/*' \
  -not -path '*/dist/*' \
  -not -path '*/build/*' \
  | head -200
```

### Step 2 — Identify key files

Read the following files if they exist (in order of priority):

1. `package.json` / `go.mod` / `pyproject.toml` / `Cargo.toml` — language & deps
2. `Makefile` / `Dockerfile` / `docker-compose.yml` — build & run
3. `*.config.*` / `.env.example` — configuration shape
4. Entry points: `main.*`, `index.*`, `app.*`, `cmd/`
5. Existing `README.md` and `CLAUDE.md` — preserve anything still valid

### Step 3 — Read critical source files

Skim the top 5-10 most important source files to understand:

- Core domain models / types
- Main entry point logic
- API routes or CLI commands
- Key business rules or algorithms

Do NOT read every file. Focus on high-signal files only.

### Step 4 — Write `README.md`

Write or overwrite `${workspace}/README.md` with:

```markdown
# <Project Name>

<1-2 sentence elevator pitch>

## 功能概述 (Features)

- Feature 1: brief description
- Feature 2: brief description

## 快速開始 (Quick Start)

### 前置需求 (Prerequisites)

- Requirement 1
- Requirement 2

### 安裝 (Installation)

<actual install commands from the project>

### 執行 (Run)

<actual run commands from the project>

## 使用方式 (Usage)

<key CLI commands, API endpoints, or UI flows>

## 架構簡介 (Architecture Overview)

<high-level description of how the system works, 3-5 sentences>

## 改善建議 (Improvement Suggestions)

Based on codebase analysis:

- [ ] Suggestion 1: rationale
- [ ] Suggestion 2: rationale
- [ ] Suggestion 3: rationale
```

**Rules for README.md:**

- Write in Traditional Chinese section headers with English in parentheses
- Use actual commands found in the project, not placeholders
- Improvement suggestions must be specific and actionable, based on real findings
- Minimum 3 suggestions, max 7
- Suggestions should cover: code quality, architecture, DX, testing, documentation

### Step 5 — Write `CLAUDE.md`

Write or overwrite `${workspace}/CLAUDE.md` with:

```markdown
# <Project Name> — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

<actual directory tree, 2-3 levels deep>

## 技術棧 (Tech Stack)

- Language: <detected>
- Framework: <detected>
- Build tool: <detected>
- Key dependencies: <top 5-8 deps>

## 關鍵決策 (Key Decisions)

- Decision 1: why this approach was chosen (inferred from code patterns)
- Decision 2: ...

## 核心邏輯 (Core Logic)

### <Component/Module 1>

Brief description of what it does and how it connects to other parts.

### <Component/Module 2>

...

## 開發指南 (Development Guide)

### 建置 (Build)

<exact build commands>

### 測試 (Test)

<exact test commands, or note if no tests exist>

### 部署 (Deploy)

<deployment method if detectable, or "未偵測到部署設定 (No deployment config detected)">

## 慣例 (Conventions)

- Naming: <detected patterns>
- Error handling: <detected patterns>
- Logging: <detected patterns>
```

**Rules for CLAUDE.md:**

- Write in Traditional Chinese section headers with English in parentheses
- Project structure must be the actual tree, not a template
- Key decisions should be inferred from code patterns (e.g., "uses dependency injection via constructor" not guesses)
- If something is not detectable, say so explicitly — do not fabricate

### Step 6 — Project basic setup (symbolic links)

After writing the docs, run the setup script to create symbolic links so
multiple AI agents share the same configuration:

```bash
bash "$(dirname "$0")/setup-links.sh" "${workspace}"
```

The script creates:

| Symlink         | Target       | Purpose             |
| --------------- | ------------ | ------------------- |
| `GEMINI.md`     | `CLAUDE.md`  | Gemini CLI 讀取     |
| `AGENTS.md`     | `CLAUDE.md`  | 通用 agent context  |
| `.geminiignore` | `.gitignore` | Gemini CLI 忽略檔案 |

**Safety:** skips if the link already exists or if the target is a regular
file (logs a `WARN`). See `setup-links.sh` for details.

### Step 7 — Summary report

After writing both files and setting up symlinks, output a brief summary:

```text
✅ summarize 完成

README.md: <line count> 行, <N> 項改善建議
CLAUDE.md: <line count> 行, <N> 個核心模組

Symlinks:
- GEMINI.md -> CLAUDE.md ✅ (created | already exists | skipped)
- AGENTS.md -> CLAUDE.md ✅ (created | already exists | skipped)
- .geminiignore -> .gitignore ✅ (created | already exists | skipped)

主要發現:
- <finding 1>
- <finding 2>
- <finding 3>
```

## Failure Modes

| Situation                                    | Action                                                                  |
| -------------------------------------------- | ----------------------------------------------------------------------- |
| Workspace is empty                           | Write minimal stubs, note "空專案 (Empty project)"                      |
| Cannot detect language/framework             | Note "未偵測到 (Not detected)" in relevant sections                     |
| Existing README/CLAUDE have valuable content | Merge — preserve valid sections, update stale ones                      |
| Too many files to scan                       | Focus on top-level + entry points, note "僅掃描部分檔案 (Partial scan)" |

## Important

- Never fabricate information. If you cannot determine something, say so.
- Preserve any existing content that is still accurate.
- All section headers use Traditional Chinese with English in parentheses.
- Commands must be real commands found in the project, not placeholders.
