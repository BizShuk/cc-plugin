# Hermes Agent

## Installation

```bash
# Copy the plugin to the appropriate location
curl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh | bash

# Verify the installation
hermes doctor

# Configure model
hermes model
```

##

```bash
# 📁 Your files:
#    Config:    /Users/shuk/.hermes/config.yaml
#    API Keys:  /Users/shuk/.hermes/.env
#    Data:      /Users/shuk/.hermes/cron/, sessions/, logs/
#    Code:      /Users/shuk/.hermes/hermes-agent

# 🚀 Commands:
hermes                  # Start chatting
hermes setup            # Configure API keys & settings
hermes config           # View/edit configuration
hermes config edit      # Open config in editor
hermes gateway install  # Install gateway service (messaging + cron)
hermes update           # Update to latest version
```
