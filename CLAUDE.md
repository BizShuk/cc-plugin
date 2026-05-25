# CC Plugin

A global Claude Code plugin demonstrating settings, hooks, monitors, skills, and agents.

## Installation

```bash
# Local development installation
claude --plugin-dir .

# Or install from marketplace when published
# /plugin install cc-plugin
```

## Components

### Skills

- `apple-calendar` - Manage Apple Calendar events
- `apple-notes` - Manage Apple Notes
- `apple-reminders` - Manage Apple Reminders
- `golang-code-quality` - Review Go code for SOLID principles and idioms
- `golang-mvc` - Guide on Go MVC patterns
- `summarize` - Onboard projects and create README.md/CLAUDE.md
- `superpower` - Check and invoke relevant skills

### Agents

- `golang-refactor` - Specialized agent for Go code refactoring
- `feature` - Guide and implement features

### Hooks

- `PostToolUse` - Automatically format files using specific tools (e.g., `go fmt` for Go files)

### Monitors

- `error-log` - Watch `./logs/error.log` for errors
- `access-log` - Watch `./logs/access.log` for access events

### LSP

- Go language server (`gopls`) for Go code intelligence

### MCP

- Filesystem MCP server for file operations

## Development

```bash
# Reload plugins without restarting
/reload-plugins
```

## Structure

```
.
├── .claude-plugin/
│   └── plugin.json       # Plugin manifest
├── agents/               # Custom Agents
│   ├── feature.md
│   └── golang-refactor.md
├── hooks/                # Hooks configuration & scripts
│   ├── hooks.json
│   └── post-tool.sh
├── monitors/             # Monitors configuration
│   └── monitors.json
├── skills/               # Custom Skills
│   ├── apple-*/
│   ├── golang-*/
│   └── ...
├── pkg/                  # Package configuration templates & resources
├── config/               # Symlinked configurations (local and home)
├── settings.json         # Default settings
├── .lsp.json             # LSP server config
├── .mcp.json             # MCP server config
├── run.sh                # Setup script for Unix/macOS
└── run.ps1               # Setup script for Windows
```

## License

MIT
