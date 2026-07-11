---
name: creating-plugin-metadata
description: Use when authoring or initializing a plugin.json manifest file for a Claude Code workspace plugin.
---

# Creating Plugin Metadata

## Overview
A guide for authoring the `.claude-plugin/plugin.json` manifest file to define custom plugins, skills, and agents.

## When to Use
- Initializing a new workspace plugin configuration.
- Defining plugin name, description, author, repository, and keywords.
- Structuring sub-components (skills, agents) for discovery.

When NOT to use:
- Generating general Go project configuration or system settings.

## Manifest Schema

The plugin metadata manifest must be placed at `.claude-plugin/plugin.json` relative to the plugin root directory.

### Structure

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

### Key Fields

- `name`: Unique name identifier for the plugin (lowercase, hyphens/numbers allowed).
- `version`: SemVer formatted version string (e.g., `1.0.0`).
- `description`: Plain text description summarizing the plugin's capabilities.
- `author`: Nested object containing author contact information.
- `homepage` / `repository`: URL links to project pages or git repositories.
- `keywords`: Array of tags for discoverability.
- `skills` / `agents`: Array of paths to specific skill/agent configurations if explicit loading is required; leave empty for auto-discovery by folder convention (i.e. `skills/` and `agents/` directories).

## Common Mistakes
- Putting the manifest in the root directory directly (e.g. `plugin.json`). It must live under the `.claude-plugin/` subdirectory.
- Using uppercase letters or special characters in the `name` field.
- Forgetting to include contact info under the `author` field.
