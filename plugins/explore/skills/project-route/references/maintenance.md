# Maintaining the index

The routing index lives at `~/projects/.project_index/projects.json`, with
`router.py` beside it. Routing itself never needs these commands — they are for
when the index is wrong, not when it is merely old.

## Rebuilding

The router auto-rebuilds whenever the index is older than `stale_after_days`
(default 3). Force a refresh only when you want it now — you just added a
project, edited tags, or want to bypass the window:

```bash
python3 ~/projects/.project_index/router.py rebuild
```

`rebuild` preserves the hand-curated `tags` and `categories`; the auto-scan
never overwrites them. New projects start with an empty tag list.

## Tags

A project with no tags is invisible to tag-based routing. List the untagged
ones and fill them in:

```bash
python3 -c "import json;d=json.load(open('$HOME/projects/.project_index/projects.json'));print([k for k,v in d['projects'].items() if not v['tags']])"
```

## Categories

A top-level directory becomes a category automatically when it has no
`README.md` / `CLAUDE.md` of its own and holds 2+ project-like children. A
container that does carry its own umbrella docs (for example `game/`,
`collections/`) is not auto-detected and must be listed by hand in the index's
`categories` array — otherwise its children are indexed as subprojects rather
than as routable projects.

## Moving a project

Moving a project into a category does not change its key or its tags; only its
`path` and `category` change. Re-run `rebuild` after any such move.
