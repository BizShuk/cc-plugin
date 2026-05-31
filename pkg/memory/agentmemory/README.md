<https://github.com/rohitg00/agentmemory#how-it-works>

npx skills add rohitg00/agentmemory -y -a claude-code -a antigravity -a hermes-agent

│ REST API <http://localhost:3111> │
│ Viewer <http://localhost:3113> │
│ Streams ws://localhost:3112 │
│ Engine ws://localhost:49134 │
│ iii console (install: curl -fsSL <https://install.iii.dev/iii/main/install.sh> | sh)

> iii-engine on PATH is v0.16.1 but agentmemory v0.9.24 hard-pins v0.11.2. Engine API drift causes runtime failures (e.g. state::list-not-found on v0.13.0). Downgrade with: `curl -fsSL https://github.com/iii-hq/iii/releases/download/iii/v0.11.2/iii-aarch64-apple-darwin.tar.gz | tar -xz -C ~/.local/bin`. Or set AGENTMEMORY_III_VERSION=0.16.1 to override at your own risk.

`AGENTMEMORY_III_VERSION=0.16.1 agentmemory`

echo "export AGENTMEMORY_III_VERSION=0.16.1" >> ~/.bash_plugin
pm2 start agentmemory
