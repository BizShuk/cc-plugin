"""Stage 1: git log --numstat → _raw/commits.jsonl."""

import re
from pathlib import Path
from utils import git_log_numstat, write_jsonl


COMMIT_RE = re.compile(r'^Commit: (\S+) \| (.+?) \| (\d{4}-\d{2}-\d{2}) \| (.+)$')


def parse_git_log(text: str) -> list[dict]:
    """Parse git log --numstat output into list of commit dicts."""
    commits = []
    current = None
    files = []

    for line in text.split('\n'):
        if not line.strip():
            if current:
                current['files'] = files
                commits.append(current)
                current = None
                files = []
            continue

        m = COMMIT_RE.match(line)
        if m:
            if current:
                current['files'] = files
                commits.append(current)
                files = []
            current = {
                'hash': m.group(1),
                'author': m.group(2),
                'date': m.group(3),
                'subject': m.group(4),
                'files': [],
            }
            continue

        if current:
            parts = line.split('\t')
            if len(parts) >= 3:
                try:
                    adds = int(parts[0].replace('-', '0'))
                    dels = int(parts[1].replace('-', '0'))
                except ValueError:
                    continue
                files.append({'adds': adds, 'dels': dels, 'file': parts[2]})

    if current:
        current['files'] = files
        commits.append(current)

    return commits


def collect(repo: Path, output_dir: Path) -> Path:
    """Run git log --numstat on a repo, write _raw/commits.jsonl."""
    text = git_log_numstat(repo)
    commits = parse_git_log(text)
    out = output_dir / "_raw" / "commits.jsonl"
    write_jsonl(out, commits)
    return out
