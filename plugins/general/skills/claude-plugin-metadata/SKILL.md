---
name: claude-plugin-metadata
description: >
    Use when authoring, initializing, or updating plugin.json and marketplace.json manifests for a Claude Code workspace plugin or individual skills. Triggers on: "create plugin", "update plugin", "plugin.json", "marketplace.json", "register skill".
version: "1.0.0"
allowed-tools: Read, Write, Bash
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: reference
    platforms: [macos, linux]
---

# Claude Plugin Metadata

## Overview

A guide for authoring the `.claude-plugin/plugin.json` and `.claude-plugin/marketplace.json` manifest files to define and register custom plugins, skills, and agents.

## When to Use

- Initializing a new workspace plugin configuration.
- Defining plugin name, description, author, repository, and keywords.
- Registering or updating plugins in the workspace-wide marketplace.json.
- Structuring sub-components (skills, agents) for discovery.

When NOT to use:

- Generating general Go project configuration or system settings.

## Manifest Schemas

### 1. Plugin Metadata Manifest (plugin.json)

The plugin metadata manifest must be placed at `.claude-plugin/plugin.json` relative to the plugin root directory.

#### Structure

```json
{
    "name": "plugin-name",
    "version": "1.0.0",
    "description": "Short description of what the plugin does",
    "author": {
        "name": "Author Name",
        "email": "email@example.com",
        "url": "https://github.com/username"
    },
    "homepage": "https://github.com/username/repository",
    "repository": "https://github.com/username/repository",
    "license": "MIT",
    "keywords": ["keyword1", "keyword2"]
}
```

Default `skills/` and `agents/` are auto-discovered. Omit those fields unless the next section applies.

#### Key Fields

- `name`: Unique name identifier for the plugin (lowercase, hyphens/numbers allowed).
- `version`: SemVer formatted version string (e.g., `1.0.0`).
- `description`: Plain text description summarizing the plugin's capabilities.
- `author`: Nested object containing author contact information.
- `homepage` / `repository`: URL links to project pages or git repositories.
- `keywords`: Array of tags for discoverability.
- `skills`: Only unexpected locations or external sources. See below.
- `agents`: Omit for default `agents/`. List only if agents live outside that directory.

#### `skills` field

The default `skills/` directory is always scanned. The field **adds** extra locations; it does not replace `skills/`.

| Situation                                        | Put in `skills`?            |
| ------------------------------------------------ | --------------------------- |
| Default path `./skills/<name>/SKILL.md`          | No — auto-discovered        |
| Git submodule checked out under `skills/`        | No — still the default path |
| Unexpected location (e.g. `./lib/extra-skills/`) | Yes                         |
| External source (`owner/repo-skill`)             | Yes                         |

Same rule on a marketplace plugin entry: do not list default-path or submodule skills there either.

```json
{
    "name": "plugin-name",
    "version": "1.0.0",
    "skills": ["./lib/extra-skills/", "owner/repo-skill"]
}
```

### 2. Workspace Marketplace Registry (marketplace.json)

The workspace plugin marketplace manifest must be placed at `.claude-plugin/marketplace.json` relative to the workspace root directory.

#### Structure

```json
{
    "name": "workspace-marketplace",
    "owner": {
        "name": "Developer Name",
        "email": "developer@example.com"
    },
    "plugins": [
        {
            "name": "plugin-name",
            "source": "./plugins/plugin-name"
        },
        {
            "name": "vendored-plugin",
            "source": "owner/repo_name"
        }
    ]
}
```

#### Key Fields

- `name`: Name identifier for the marketplace workspace.
- `owner`: Nested object containing owner name and email.
- `plugins`: Array of plugin registrations.
    - `name`: The name of the registered plugin.
    - `source`: Where to fetch the plugin. See below.

#### `source` field

| Plugin kind                                | `source`                                        |
| ------------------------------------------ | ----------------------------------------------- |
| First-party, lives in this repo            | `./plugins/<name>`                              |
| Git submodule that is itself a GitHub repo | `owner/repo_name` — not the local checkout path |
| External GitHub, not vendored here         | `owner/repo_name`                               |

A local submodule folder is a checkout, not the marketplace source. If `.gitmodules` points at `https://github.com/owner/repo_name`, write `"source": "owner/repo_name"`. Do not write `"./plugins/<name>"` just because the folder exists.

## Common Mistakes

- Putting the manifest in the root directory directly (e.g. `plugin.json` instead of `.claude-plugin/plugin.json`).
- Using uppercase letters or special characters in the `name` field.
- Forgetting to include contact info under the `author` field.
- Forgetting to update the `marketplace.json` when adding a new local plugin folder.
- Listing default-path skills (`./skills`, `./skills/<name>`) or a submodule under `skills/` in the `skills` field.
- Using a local path (`./plugins/<name>`) for a plugin that is a git submodule of a GitHub repo. Use `owner/repo_name`.
