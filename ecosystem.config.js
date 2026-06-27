module.exports = {
    apps: [
        // Agent Memory (Agent)
        {
            name: "Agent Memory",
            script: "/Users/shuk/.local/nvm/versions/node/v24.11.1/bin/agentmemory",
            namespace: "Agent",
            instances: 1
        },
        // agy-cc-plugin (planner)
        {
            name: "agy-cc-plugin",
            script: "agy",
            args: [
                "--add-dir",
                "/Users/shuk/projects/cc-plugin",
                "-p",
                "'run /system-planner for current workspace, and output under <workspace>/plans/'"
            ],
            namespace: "planner",
            cwd: "/Users/shuk/projects/cc-plugin",
            instances: 1,
            cron: "10 0-9 * * *"
        }
    ]
};
