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
  "keywords": [
    "keyword1",
    "keyword2"
  ],
  "agents": [],
  "skills": []
}
```

#### Key Fields
- `name`: Unique name identifier for the plugin (lowercase, hyphens/numbers allowed).
- `version`: SemVer formatted version string (e.g., `1.0.0`).
- `description`: Plain text description summarizing the plugin's capabilities.
- `author`: Nested object containing author contact information.
- `homepage` / `repository`: URL links to project pages or git repositories.
- `keywords`: Array of tags for discoverability.
- `skills` / `agents`: Array of paths to specific skill/agent configurations if explicit loading is required; leave empty for auto-discovery by folder convention (i.e. `skills/` and `agents/` directories).

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
      "name": "external-plugin",
      "source": {
        "source": "github",
        "repo": "username/repository"
      }
    }
  ]
}
```

#### Key Fields
- `name`: Name identifier for the marketplace workspace.
- `owner`: Nested object containing owner name and email.
- `plugins`: Array of plugin registrations.
  - `name`: The name of the registered plugin.
  - `source`: The source location of the plugin. Can be a relative path or a github source object.

## Common Mistakes
- Putting the manifest in the root directory directly (e.g. `plugin.json` instead of `.claude-plugin/plugin.json`).
- Using uppercase letters or special characters in the `name` field.
- Forgetting to include contact info under the `author` field.
- Forgetting to update the `marketplace.json` when adding a new local plugin folder.
