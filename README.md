# README

## MCP List

`claude mcp add --scope [local|project|user] <mcp_name> <mcp_command>`

- `claude mcp add playwright npx @playwright/mcp@latest`, [Playwright](https://playwright.dev/docs/getting-started-mcp)
- `claude mcp add chrome-devtools npx @chrome-devtools/mcp@latest`, [Chrome DevTools](https://github.com/ChromeDevTools/chrome-devtools-mcp)

    ```sh
    /plugin marketplace add ChromeDevTools/chrome-devtools-mcp
    /plugin install chrome-devtools-mcp@chrome-devtools-plugins
    ```
