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
        }
    ]
};
