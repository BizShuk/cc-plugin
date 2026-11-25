module.exports = {
    apps: [
        // Local Ollama LLM server (`ollama serve`), registered paused.
        {
            name: "Ollama Server - Always",
            script: "ollama",
            namespace: "Agent",
            args: ["serve"],
            instances: 1,
            optional: true
        },
        // Hermes agent gateway on the shared ~/.venv, registered paused.
        {
            name: "Hermes Gateway - Always",
            script: "~/.hermes/scripts/gateway.sh",
            namespace: "Agent",
            instances: 1,
            optional: true,
            env: {
                HERMES_HOME: "~/.hermes",
                VENV_PATH: "~/.venv",
                HERMES_SUPERVISED_CHILD: "1"
            }
        },
        // Hermes WebUI (foreground bootstrap on :8787) on the shared ~/.venv, registered paused.
        {
            name: "Hermes WebUI - Always",
            script: "~/.hermes/scripts/webui.sh",
            namespace: "Agent",
            instances: 1,
            optional: true,
            env: {
                HERMES_HOME: "~/.hermes",
                VENV_PATH: "~/.venv"
            }
        },
        // Antigravity CLI Remote Control foreground server (`mac-server`), clearing GEMINI_API_KEY so OAuth auth is used; registered paused and restarted daily at 00:00.
        {
            name: "Agy Remote Control - Always",
            script: "agy",
            args: ["remote-control", "serve"],
            namespace: "Agent",
            instances: 1,
            cron_restart: "0 0 * * *",
            optional: true
        }
    ]
};
