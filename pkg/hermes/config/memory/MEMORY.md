§
## Language

- Reply in Traditional Chinese by default.
- Keep names and terms in their local language, followed by English in brackets, e.g. 中正紀念堂 (Chiang Kai-shek Memorial Hall).
§
## Where Things Live

- Home (`~`) is organized by owner: each agent keeps its own state in a dot folder (`~/.hermes`, `~/.claude`, `~/.gemini`), tooling and runtime state live in `~/.local`, and global CLI tools live in `~/bin`.
- Every project lives in `~/projects/`. Projects sit side by side as independent services, OR under a category folder. A category is itself a thin container repo (README+CLAUDE.md+.gitmodules); subprojects are git submodules each with its own Unified Interface (README.md, CLAUDE.md, AGENTS.md symlink, README.todo).
- 14 categories as of 2026-09-24 (one-line each — read README.md/CLAUDE.md for detail):
  - `ai/` — agent frameworks: agentSDK, cc-plugin, conversation_agent, customer_service, m-agent (contains picoclaw 3rd-level)
  - `collections/` — personal receipt Go CLI; vision model (currently `minimax`); pm2 cron; safe-by-default dry-run
  - `data/` — external-data fetch/normalize/store: datahub, msghub, open_data, realtime_dataflow, stock, ticketmaster, trifecta, vid-note, yfin
  - `env_setup/` — macOS/Ubuntu bootstrap shell framework; owns `~/bin` symlink, idempotent `.bash_plugin`; domains in `scripts/{system,network,cleanup,io,dump,install,backup,uninstall}`
  - `game/` — design KB + prototypes; `game-architect/` is the KB (mechanisms/lineage/art_pipeline); lower-case+underscores naming, no bold, Mermaid edge text in double-quotes
  - `iphone/` — iOS apps + Apple platform skills; push, template, MinimalBrowser, castlan, iphone_sync, md-viewer, surf_eye, tally; all target iPhone 17 Pro Max simulator only
  - `platform/` — bizshuk backend services: ads, feedback, gosdk, identity, inf, n8n, payment, rnet, superset, vscode-shuk
  - `playground/` — disposable experiments: archive, ltv
  - `product/` — end-user products: bizshuk.github.io, foreclosed, fun, ip-incubation, night_market, outfit, payment_gateway, surf_*, tools-web
  - `research/` — domain notes: ab, business-strategy, coding, company, daily-feedback, healthcare-system, linux, mens-outfit, property, pty, relationship, seo, ubuntu-server-tuning, work; plus `supermemory-research-report-2026-07-16.md`
  - `social/` — ONE Instagram+Threads app (Next.js+Vite+drizzle); not a category
  - `tools/` — CLI/dev utils: auth, autop, dux, go-dependency-analysis, gx, img, macemailapp, macnotesapp, mactrans, mdserver, pm2, port, proxy, sandbox, sbackup, sessiond, skills-cli, trans, video-utils, voice, vscoed-plugin, ytdl
  - `web/` — automation, config, static `index.html`; `tmp/` empty scratch
- Conventions: each category repo holds `.gitmodules/.gitignore/.claude-plugin/plugin.json/.vscode/settings.json`. Submodule pointer bumps at category layer; subproject edits in subproject. Multi-agent skill dirs (`.claude/.agents/.grok/.hermes/.pi/skills/`) are aliases to one `skills/` source.
- Code grouped by domain (`cmd/`, `svc/`, `model/`), not by file type. Share settings/ignores through symlinks (one source of truth).
§
## Memory file layout warning
- `~/.hermes/memories/USER.md` is a SYMLINK → `~/projects/ai/cc-plugin/config/CLAUDE.global.md` (git-tracked). NEVER edit via memory tool — pollutes skills installs and cc-plugin tree.
- `~/.hermes/memories/MEMORY.md` is currently the same APFS inode (hardlink) as `~/projects/ai/cc-plugin/pkg/hermes/config/memory/MEMORY.md`. Edits here show as a working-tree `M` in that cc-plugin subdir. Per user 2026-09-24: modify freely but DO NOT commit or push; user will handle the cc-plugin side manually.
- For one-off env facts outside the user profile, use this file (`memory target=memory`).
§
## Apple Apps

- Use Apple Reminders app to store to-do list
- Use Apple Calendar app to store events/schedule and share the event with <biz.shuk@gmail.com>(Apple account)
