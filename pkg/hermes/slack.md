# slack channel

<https://hermes-agent.nousresearch.com/docs/user-guide/messaging/slack>

1. generate manifest json file and copy it

```bash
hermes slack manifest --write
```

2. Go to <https://api.slack.com/apps> → Create New App → From an app manifest => paste the manifest

3. Install App to Workspace

    ```text
    In the sidebar, go to Settings → Install App
    Click Install to Workspace
    Review the permissions and click Allow
    After authorization, you'll see a Bot User OAuth Token starting with xoxb-
    Copy this token — this is your SLACK_BOT_TOKEN
    ```

4. find member id

    ```text
    In Slack, click on the user's name or avatar
    Click View full profile
    Click the ⋮ (more) button
    Select Copy member ID
    ```

5. config `$HOME/.config/hermes/.env`

```env
# Required
SLACK_BOT_TOKEN=xoxb-your-bot-token-here
SLACK_APP_TOKEN=xapp-your-app-token-here
SLACK_ALLOWED_USERS=U01ABC2DEF3              # Comma-separated Member IDs

# Optional
SLACK_HOME_CHANNEL=C01234567890              # Default channel for cron/scheduled messages
SLACK_HOME_CHANNEL_NAME=general              # Human-readable name for the home channel (optional)
```

6. /invite @Hermes Agent
7. configuration

    ```yaml
    platforms:
    slack:
        # Controls how multi-part responses are threaded
        # "off"   — never thread replies to the original message
        # "first" — first chunk threads to user's message (default)
        # "all"   — all chunks thread to user's message
        reply_to_mode: "first"

        extra:
        # Whether to reply in a thread (default: true).
        # When false, channel messages get direct channel replies instead
        # of threads. Messages inside existing threads still reply in-thread.
        reply_in_thread: true

        # Also post thread replies to the main channel
        # (Slack's "Also send to channel" feature).
        # Only the first chunk of the first reply is broadcast.
        reply_broadcast: false
    ```
