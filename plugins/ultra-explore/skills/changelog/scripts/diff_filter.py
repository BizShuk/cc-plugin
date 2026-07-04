"""Stage 3: git log -p per week, filter to kept files, write _diffs/YYYY-Www.diff."""

import re
from pathlib import Path
from utils import git_log_p, read_json, write_text
from exclusion import is_excluded_default, GitignoreOverlay

# Match git patch file headers: --- a/path/to/file or +++ b/path/to/file
PATCH_FILE_RE = re.compile(r'^[-+]{3} [ab]/(.+)$')
# Match `diff --git a/path b/path` and capture the path
DIFF_GIT_RE = re.compile(r'^diff --git a/(.+?) b/.*$')


def _is_file_excluded(file_path: str, repo_name: str, gitignore) -> bool:
    exc, _ = is_excluded_default(file_path, repo_name)
    if not exc and gitignore and gitignore.matches(file_path):
        exc = True
    return exc


def filter_diff(repo: Path, since: str, until: str,
                gitignore, repo_name: str) -> str:
    """Run git log -p for date range, return only hunks of kept files.

    Strategy: track current file via `diff --git` headers, then evaluate
    `--- a/...` and `+++ b/...` to set the per-file excluded state. Lines
    belonging to an excluded file (including the `diff --git` header) are
    skipped. The excluded state is reset on every new `diff --git` line.
    """
    text = git_log_p(repo, since, until)
    lines = text.split('\n')
    out = []
    current_file = None
    current_excluded = False

    for line in lines:
        # Detect a new file section start
        m_git = DIFF_GIT_RE.match(line)
        if m_git:
            current_file = m_git.group(1)
            current_excluded = _is_file_excluded(current_file, repo_name, gitignore)
            if current_excluded:
                # Skip the entire file section (header, ---, +++, hunks)
                current_file = None
                current_excluded = True
                continue
            out.append(line)
            continue

        # Detect --- a/path or +++ b/path (refines current_file; both refer
        # to the same file in a normal patch). Keep behavior consistent with
        # diff --git decision.
        m_patch = PATCH_FILE_RE.match(line)
        if m_patch:
            # If we have not yet excluded, check this path. They should match
            # the diff --git path, but stay defensive.
            if current_excluded:
                continue
            out.append(line)
            continue

        # Hunk / body lines
        if current_excluded:
            continue
        out.append(line)

    return '\n'.join(out)


def generate_diffs(repo: Path, output_dir: Path,
                   no_gitignore: bool = False) -> list[Path]:
    """Read stats.json, generate filtered diffs per week."""
    stats = read_json(output_dir / "stats.json")
    repo_name = stats["repo"]
    gitignore = None if no_gitignore else GitignoreOverlay(repo)
    diffs_dir = output_dir / "_diffs"
    diffs_dir.mkdir(parents=True, exist_ok=True)
    paths = []

    for bucket in stats["weekly_buckets"]:
        week = bucket["week"]
        since = bucket["start"]
        until = bucket["end"]
        filtered = filter_diff(repo, since, until, gitignore, repo_name)
        out = diffs_dir / f"{week}.diff"
        write_text(out, filtered)
        paths.append(out)

    return paths
