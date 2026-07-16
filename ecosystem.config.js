module.exports = {
    apps: [
        // Ollama
        {
            name: "Ollama",
            script: "ollama",
            namespace: "Agent",
            args: ["serve"],
            instances: 1
        }, // Agent Memory (Agent)
        {
            name: "Agent Memory",
            script: "agentmemory",
            namespace: "Agent",
            instances: 1
        },
        // agy-cc-plugin (planner)
        {
            name: "agy-cc-plugin-system",
            script: "agy",
            args: [
                "--add-dir",
                __dirname,
                "-p",
                "run /system-planner for current workspace"
            ],
            namespace: "planner",
            instances: 1,
            cron: "10 0-9 * * *",
            autorestart: false,
            watch: false
        }
    ]
};
