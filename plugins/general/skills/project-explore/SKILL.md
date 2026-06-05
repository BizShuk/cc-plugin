---
name: project-explore
description: >
    Use when onboarding to a new project, after major refactors, or when README.md
    and CLAUDE.md are missing, outdated, or empty. Explores the entire workspace
    and writes structured documentation into two files at the repo root.
    README.md covers functional requirements (business domains, domain flow, entities).
    CLAUDE.md covers non-functional requirements (structure, tech decisions, build/deploy).
allowed-tools: Read, Bash, Glob, Write, Edit
disable-model-invocation: false
user-invocable: true
effort: high
context: fork
---

# project-explore

Explore and conclude the whole workspace, then write structured documentation
into `README.md` and `CLAUDE.md` at the repository root.

## Overview

This skill performs a full-workspace scan and produces two canonical docs:

| File        | Focus                                                                                                      |
| ----------- | ---------------------------------------------------------------------------------------------------------- |
| `README.md` | `功能性需求 (Functional Requirements)`: business domains, domain flow, entities, use cases                 |
| `CLAUDE.md` | `非功能性需求 (Non-functional Requirements)`: project structure, tech decisions, build/deploy, conventions |

### Design Philosophy

A project may have many handlers, services, or modules, but they belong to only
a few `業務領域 (Business Domains)`. The README should be organized by domain,
not by file or handler.

`Example:` A data service may have 20+ handlers, but they belong to ~4 domains:

- `資料複製 (Data Replication)` — pure data sync across regions
- `用戶遷移 (User Migration)` — user region change workflows
- `資料修復 (Data Fix)` — ad-hoc data correction utilities
- `稽核日誌 (Audit Log)` — tracking and compliance

The README groups these by domain with flow descriptions, not by listing every handler.

## When to Use

- Project has no `README.md` or `CLAUDE.md`
- Existing docs are stale after significant changes
- Onboarding to an unfamiliar codebase
- After a major refactor or migration
- User explicitly asks to summarize or document the workspace

## Procedure

### Step 1 — Discover project layout

Use the `Glob` tool to discover files. Apply the following exclusion patterns:

```
Excluded directories and files:
  .git, .svn, .hg, .DS_Store, Thumbs.db
  archive/, *.bak*
  .geminiignore, .gitlab, .pre-commit-config.yaml
  *.code-workspace, .golangci.yml, go.sum
  .specify, .gemini, .agent, .serena, .ttadk, .coco
  .devops, .settings
  .classpath, .project
  target/, out/, dist/, output/
  .mvn/, node_modules/, __pycache__/, .venv/
  __debug_bin*
  gen/**, kitex_gen/**, thrift_gen/**
  .playwright-mcp
  .vscode/, .claude/
```

```bash
# Step 1a: top-level structure
Glob("*", maxDepth=1)

# Step 1b: deeper structure (exclude noise)
Glob("**/*", maxDepth=5, exclude=[
  "**/.git/**", "**/.svn/**", "**/.hg/**",
  "**/.DS_Store", "**/Thumbs.db",
  "**/archive/**", "**/*.bak*",
  "**/.geminiignore", "**/.gitlab/**",
  "**/.pre-commit-config.yaml",
  "**/*.code-workspace", "**/.golangci.yml", "**/go.sum",
  "**/.specify/**", "**/.gemini/**", "**/.agent/**",
  "**/.serena/**", "**/.ttadk/**", "**/.coco/**",
  "**/.devops/**", "**/.settings/**",
  "**/.classpath", "**/.project",
  "**/target/**", "**/out/**", "**/dist/**", "**/output/**",
  "**/.mvn/**", "**/node_modules/**",
  "**/__pycache__/**", "**/.venv/**",
  "**/__debug_bin*",
  "**/gen/**", "**/kitex_gen/**", "**/thrift_gen/**",
  "**/.playwright-mcp/**"
])
```

If `Glob` tool is not available, fall back to bash:

```bash
find . -maxdepth 3 -type f \
  -not -path '*/.git/*' -not -path '*/.svn/*' -not -path '*/.hg/*' \
  -not -name '.DS_Store' -not -name 'Thumbs.db' \
  -not -path '*/archive/*' -not -name '*.bak*' \
  -not -path '*/.gemini/*' -not -path '*/.agent/*' \
  -not -path '*/.serena/*' -not -path '*/.ttadk/*' -not -path '*/.coco/*' \
  -not -path '*/.specify/*' -not -path '*/.gitlab/*' \
  -not -path '*/.devops/*' -not -path '*/.settings/*' \
  -not -path '*/target/*' -not -path '*/out/*' \
  -not -path '*/dist/*' -not -path '*/output/*' \
  -not -path '*/.mvn/*' -not -path '*/node_modules/*' \
  -not -path '*/__pycache__/*' -not -path '*/.venv/*' \
  -not -name '__debug_bin*' \
  -not -path '*/gen/*' -not -path '*/kitex_gen/*' -not -path '*/thrift_gen/*' \
  -not -path '*/.playwright-mcp/*' \
  -not -name '*.code-workspace' -not -name '.golangci.yml' -not -name 'go.sum' \
  | sort | head -200
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

### Step 4 — Identify business domains

This is the critical analysis step. Group all handlers, services, and modules
into `業務領域 (Business Domains)`:

1. List all handler/controller/route files
2. Identify common themes and purposes
3. Group them into 3-7 domains
4. For each domain, trace the data flow: `entry point → service → repository → external`

### Step 5 — Write `README.md`

Write or overwrite `${workspace}/README.md` with:

```markdown
# <Project Name>

<1-2 sentence elevator pitch: what business problem this solves>

## 業務領域 (Business Domains)

### <Domain 1 Name>

<2-3 sentences: what this domain does, why it exists, when it is triggered>

`領域流程 (Domain Flow):`

1. <Step 1: entry point / trigger>
2. <Step 2: core processing>
3. <Step 3: outcome / side effects>

`核心實體 (Key Entities):` <Entity A>, <Entity B>, <Entity C>

`相關處理器 (Related Handlers):` <HandlerX>, <HandlerY>

---

### <Domain 2 Name>

<same structure as above>

---

## 領域關聯 (Domain Relationships)

<Describe how domains interact with each other. Which domain's output is
another domain's input? Are there shared entities?>

## 使用方式 (Usage)

<key CLI commands, API endpoints, or UI flows — organized by domain>

## 改善建議 (Improvement Suggestions)

Based on codebase analysis:

- [ ] Suggestion 1: rationale
- [ ] Suggestion 2: rationale
- [ ] Suggestion 3: rationale
```

`Rules for README.md:`

- Write in Traditional Chinese section headers with English in parentheses
- Organize by `業務領域 (Business Domain)`, NOT by file or handler
- Each domain must have: description, flow, key entities, related handlers
- Domain flow should trace the real code path, not be abstract
- Use actual function/handler names found in the project
- Improvement suggestions must be specific and actionable, based on real findings
- Minimum 3 suggestions, max 7
- Suggestions should cover: domain boundaries, missing use cases, data flow gaps

### Step 6 — Write `CLAUDE.md`

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

## 模組對應 (Module Mapping)

Map each business domain (from README) to its technical implementation:

| 業務領域 (Domain) | 套件/模組 (Package/Module) | 進入點 (Entry Point) |
| ----------------- | -------------------------- | -------------------- |
| <Domain 1>        | `pkg/xxx`, `handler/yyy`   | `HandleXxx()`        |
| <Domain 2>        | `pkg/aaa`, `handler/bbb`   | `HandleAaa()`        |

## 開發指南 (Development Guide)

### 前置需求 (Prerequisites)

- Requirement 1
- Requirement 2

### 安裝 (Installation)

<actual install commands from the project>

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
- Testing: <detected patterns>
```

`Rules for CLAUDE.md:`

- Write in Traditional Chinese section headers with English in parentheses
- Project structure must be the actual tree, not a template
- Must include `模組對應 (Module Mapping)` table linking domains to code locations
- Key decisions should be inferred from code patterns (e.g., "uses dependency injection via constructor" not guesses)
- Commands must be real commands found in the project, not placeholders
- If something is not detectable, say so explicitly — do not fabricate

### Step 7 — Project basic setup (symbolic links)

After writing the docs, run the setup script to create symbolic links so
multiple AI agents share the same configuration:

```bash
bash "$(dirname "$0")/setup-links.sh" "${workspace}"
```

The script creates:

| Symlink         | Target       | Purpose             |
| --------------- | ------------ | ------------------- |
| `AGENTS.md`     | `CLAUDE.md`  | 通用 agent context  |
| `.geminiignore` | `.gitignore` | Gemini CLI 忽略檔案 |

`Safety:` skips if the link already exists or if the target is a regular
file (logs a `WARN`). See `setup-links.sh` for details.

### Step 8 — Summary report

After writing both files and setting up symlinks, output a brief summary:

```text
✅ summarize 完成

README.md: <line count> 行, <N> 個業務領域, <N> 項改善建議
CLAUDE.md: <line count> 行, <N> 個核心模組

Symlinks:
- AGENTS.md -> CLAUDE.md ✅ (created | already exists | skipped)
- .geminiignore -> .gitignore ✅ (created | already exists | skipped)

業務領域摘要:
- <Domain 1>: <1-sentence summary>
- <Domain 2>: <1-sentence summary>
- <Domain 3>: <1-sentence summary>
```

## Failure Modes

| Situation                                    | Action                                                                           |
| -------------------------------------------- | -------------------------------------------------------------------------------- |
| Workspace is empty                           | Write minimal stubs, note "空專案 (Empty project)"                               |
| Cannot detect language/framework             | Note "未偵測到 (Not detected)" in relevant sections                              |
| Existing README/CLAUDE have valuable content | Merge — preserve valid sections, update stale ones                               |
| Too many files to scan                       | Focus on top-level + entry points, note "僅掃描部分檔案 (Partial scan)"          |
| Too many files to scan                       | Focus on top-level + entry points, note "僅掃描部分檔案 (Partial scan)"          |
| Cannot identify clear domains                | Group by package/directory and note "領域邊界不明確 (Domain boundaries unclear)" |

## Important

- Never fabricate information. If you cannot determine something, say so.
- Preserve any existing content that is still accurate.
- All section headers use Traditional Chinese with English in parentheses.
- Commands must be real commands found in the project, not placeholders.
- README focuses on `WHAT` the system does (functional). CLAUDE.md focuses on `HOW` it works (technical).
- Organize README by business domain, not by file structure.
