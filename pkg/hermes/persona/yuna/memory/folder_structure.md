---
trigger: always_on
description: Overall AI Rules
globs: "**/*"
---

# Overall AI Rules

Apply these mental models systmatically when analyzing problems and generating solutions.

## Language

Use `Traditional Chinese` to reply message or describe the information/details. For name/term, should use local language with English and round brackets.

ex1: 中正紀念堂 (Chiang Kai-shek Memorial Hall)

中正紀念堂 is in Taiwan Taipei, here is using Traditional Chinese, so use it by default, but still need to attach with English as it's a name/term

ex2: Catedral de Santa Eulalia de Barcelona(Barcelona Cathedral)

## folder structure/file content

### folder structure

This rule defines the standard directory structure centered in the user's home directory (`~`).

User Home Directory Tree Overview

```text
<User Home Directory>
├── .gemini/                # Gemini-specific brain and context
│   ├── GEMINI.md           # Global rule for gemini cli and antigravity
│   └── antigravity/        # Antigravity agent specific state
│       ├── skills/         # Antigravity agent skills
│       └── global_workflows/ # Antigravity global workflow
├── .claude/                # Claude agent specific State
├── .local/                 # Local tools & env (Go, NPM, Homebrew)
├── .hermes/                # Agent brain & workspace context
├── bin/                    # Global CLI tools
└── projects/               # Parallel project workspace
    └── ...                 # Other micro-service style projects (Parallel microservices)
```

#### System Directories

##### ~/.local/

The `.local/` folder stores user-specific data, configurations, and environment states.

```text
.local/
├── bin/                    # 二進位檔 (Binaries): Local-specific execution scripts and tools.
├── go/                     # Go 工作區 (Go Workspace): Core Go development workspace.
├── homebrew/               # Homebrew (Homebrew): Local installation and managed packages.
├── npm-global/             # NPM 全域 (NPM Global): Shared Node.js tools across projects.
├── share/                  # 共用資料 (Shared Data): Documents, icons, and shared application resources.
└── state/                  # 狀態與日誌 (State & Logs): Application runtime logs and cache states.
```

#### Projects Directory (~/projects/\*)

```text
~/projects/<project_name>/

├── .agents/                # Antigravity 專案工作流 (Antigravity project workflow) (Symlink to .agent)
│   ├── skills/             # Project-specific skills
│   ├── rules/              # Project-specific rules
│   ├── workflows/          # Project-specific workflows
│   └── settings.json       # Project-specific agent settings
├── .gemini/                # Symlink to .agent
├── .gitignore              # Git ignore file (symlinked to .geminiignore)
├── GEMINI.md               # Project-specific AI context
├── README.md               # Project overview (Must include domain folder structure)
├── README.todo             # Task tracking and open items
├── specs/                  # Phase-based specifications and plans
└── [Domain Folders]/       # Physical project structure (e.g., cmd/, svc/, model/, etc.)
```
