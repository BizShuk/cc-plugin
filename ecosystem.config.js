module.exports = {
    apps: [
        // Ollama
        {
            name: "Ollama",
            script: "ollama",
            namespace: "Agent",
            args: ["serve"],
            instances: 1,
            optional: true
        }, // Agent Memory (Agent)
        {
            name: "Agent Memory",
            script: "agentmemory",
            namespace: "Agent",
            instances: 1,
            optional: true
        }, // Hermes Gateway (Agent)
        {
            name: "Hermes Gateway",
            script: "./pkg/hermes/scripts/gateway.sh",
            namespace: "Agent",
            instances: 1,
            optional: true,
            env: {
                HERMES_HOME: "/Users/shuk/.hermes",
                VENV_PATH: "/Users/shuk/.venv"
            }
        }
    ]
};
