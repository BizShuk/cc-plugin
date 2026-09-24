## Language

- Reply in Traditional Chinese by default.
- Keep names and terms in their local language, followed by English in brackets, e.g. 中正紀念堂 (Chiang Kai-shek Memorial Hall).
§
## Where Things Live

- Home (`~`) is organized by owner: each agent keeps its own state in a dot folder (`~/.hermes`, `~/.claude`, `~/.gemini`), tooling and runtime state live in `~/.local`, and global CLI tools live in `~/bin`.
- Every project lives in `~/projects/`, and projects sit side by side as independent services.
- A project describes itself: `README.md` explains what it is and how it is laid out, `README.todo` tracks open work, and specs and plans have their own folder.
- Code is grouped by domain (`cmd/`, `svc/`, `model/`), not by file type.
- One source of truth: share agent settings and ignore files through symlinks instead of copies.
- Before creating a file, find the folder that owns that concern; if none exists, create one that follows the same pattern.
§
## Apple Apps

- Use Apple Reminders app to store to-do list
- Use Apple Calendar app to store events/schedule and share the event with <biz.shuk@gmail.com>(Apple account)
