# Firecrawl — Auth (Path D) and REST API (Path E)

## Path D: Account Authorization Or API Key

Use this when the human still needs to sign up, sign in, authorize access, or
obtain an API key.

If you ran the install command with `--browser`, the human was already prompted
to sign in. Check if the key is available before running this flow. If you
already have a valid `FIRECRAWL_API_KEY`, skip this path.

If you're the human reading this in the browser, create an account or sign in at
<https://www.firecrawl.dev/signin?view=signup&source=agent-suggested>.

If you're an agent and need the human to authorize an API key:

`Step 1 — Generate auth parameters:`

```bash
SESSION_ID=$(openssl rand -hex 32)
CODE_VERIFIER=$(openssl rand -base64 32 | tr '+/' '-_' | tr -d '=\n' | head -c 43)
CODE_CHALLENGE=$(printf '%s' "$CODE_VERIFIER" | openssl dgst -sha256 -binary | openssl base64 -A | tr '+/' '-_' | tr -d '=')
```

`Step 2 — Ask the human to open this URL:`

```text
https://www.firecrawl.dev/cli-auth?code_challenge=$CODE_CHALLENGE&source=coding-agent#session_id=$SESSION_ID
```

If they already have a Firecrawl account, they'll sign in and authorize. If not,
they'll create one first and then authorize. The API key comes back to you
automatically after they click "Authorize."

`Step 3 — Poll for the API key:`

```bash
POST https://www.firecrawl.dev/api/auth/cli/status
Content-Type: application/json

{"session_id": "$SESSION_ID", "code_verifier": "$CODE_VERIFIER"}
```

Poll every 3 seconds. Responses:

- `{"status": "pending"}` — keep polling
- `{"status": "complete", "apiKey": "fc-...", "teamName": "..."}` — done

`Step 4 — Save the key and continue:`

```bash
echo "FIRECRAWL_API_KEY=fc-..." >> .env
```

## Path E: Use Firecrawl Without Installing Anything

Use this when you don't want to install a CLI or skills package. This works for
both use cases:

- `Live web work` — an agent calling the API directly for search, scrape, or interact during a session
- `Building with Firecrawl` — integrating the REST API into app code

You still need an API key. Two ways to get one:

- `Human pastes it in` — if you already have a key, just set `FIRECRAWL_API_KEY=fc-...` in your environment or pass it directly
- `Automated flow` — do Path D to walk the human through browser auth and receive the key automatically

`Base URL:` `https://api.firecrawl.dev/v2`

`Auth header:` `Authorization: Bearer fc-YOUR_API_KEY`

### Available endpoints

- `POST /search` — discover pages by query, returns results with optional full-page content
- `POST /scrape` — extract clean markdown from a single URL
- `POST /interact` — browser actions on live pages (clicks, forms, navigation)
- `POST /support/ask` — diagnose a failing Firecrawl call. Pass `{ question, jobId? }`; returns a prose `answer` plus machine-readable `fixParameters` to retry with. Auto-scoped to your team via the bearer key
- `POST /support/docs-search` — answer "how do I…" questions from Firecrawl's official docs. Pass `{ question }`; returns the answer plus citations to the docs pages used

### Documentation and references

The API docs are the source of truth for request/response schemas, parameters,
and SDKs:

- `API reference:` <https://docs.firecrawl.dev>
- `Skills repo` (for agent integration patterns): <https://github.com/firecrawl/skills>

## Session-specific auth

Use this API key for this onboarding session (or set the `FIRECRAWL_API_KEY`
environment variable):

```dotenv
FIRECRAWL_API_KEY=fc-YOUR_API_KEY
```
