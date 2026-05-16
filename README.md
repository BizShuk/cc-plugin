# README

```bash
# install claude code
curl -fsSL https://claude.ai/install.sh | bash
```

## MCP List

`claude mcp add --scope [local|project|user] <mcp_name> <mcp_command>`

- [Playwright](https://playwright.dev/docs/getting-started-mcp), `claude mcp add --scope project playwright npx @playwright/mcp@latest`
- [Chrome DevTools](https://github.com/ChromeDevTools/chrome-devtools-mcp), `claude mcp add --scope project  chrome-devtools npx @chrome-devtools/mcp@latest`

    ```sh
    /plugin marketplace add ChromeDevTools/chrome-devtools-mcp
    /plugin install chrome-devtools-mcp@chrome-devtools-plugins
    ```

## skills

- [Scrapling 技能 (Scrapling Skill)](https://github.com/D4Vinci/Scrapling/blob/main/agent-skill/Scrapling-Skill/SKILL.md)
