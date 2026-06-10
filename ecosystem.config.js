module.exports = {
    apps: [
        {
            name: "Agent Memory",
            script: "agentmemory",
            error_file: "~/.config/agentmemory/daemon.err",
            out_file: "~/.config/agentmemory/daemon.out",
            autorestart: true,
            env: {
                NODE_ENV: "development"
            },
            env_production: {
                NODE_ENV: "production"
            }
        }
    ]
};
