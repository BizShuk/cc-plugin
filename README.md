# My Plugin

A sample Claude Code plugin demonstrating all available plugin components.

## Installation

```bash
# Local development
claude --plugin-dir ./my-plugin

# Or install from marketplace when published
/plugin install my-plugin
```

## Components

### Skills

- `/my-plugin:hello` - Greet users with a personalized message
- `/my-plugin:code-review` - Review code for best practices

### Agents

- `security-reviewer` - Specialized agent for security-focused code review

### Hooks

- `PostToolUse` - Runs lint:fix after Write/Edit operations
- `PostAgentMessage` - Logs after agent completion

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
my-plugin/
├── .claude-plugin/
│   └── plugin.json       # Plugin manifest
├── agents/
│   └── security-reviewer.json
├── hooks/
│   └── hooks.json
├── monitors/
│   └── monitors.json
├── skills/
│   ├── hello/
│   │   └── SKILL.md
│   └── code-review/
│       └── SKILL.md
├── bin/                  # Executables added to PATH
├── settings.json         # Default settings
├── .lsp.json           # LSP server config
└── .mcp.json           # MCP server config
```

## License

MIT
