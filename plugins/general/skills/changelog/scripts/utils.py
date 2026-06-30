"""Shared helpers: git commands, ISO week, file I/O."""

import subprocess
import json
from pathlib import Path
from datetime import datetime
from typing import Optional


def git_log_numstat(repo: Path, date_format: str = "short") -> str:
    """Run git log --numstat and return stdout."""
    fmt = f"Commit: %h | %an | %ad | %s"
    cmd = ["git", "-C", str(repo), "log", "--numstat",
           f"--pretty=format:{fmt}", f"--date={date_format}"]
    raw = subprocess.run(cmd, capture_output=True, text=True, timeout=300).stdout
    # Indent file-change lines (non-Commit, non-blank) with two spaces
    lines = raw.split('\n')
    out = []
    for line in lines:
        if line and not line.startswith('Commit:'):
            out.append('  ' + line)
        else:
            out.append(line)
    return '\n'.join(out)


def git_log_p(repo: Path, since: str, until: str) -> str:
    """Run git log -p for a date range and return stdout."""
    cmd = ["git", "-C", str(repo), "log", "-p",
           f"--since={since}", f"--until={until}"]
    return subprocess.run(cmd, capture_output=True, text=True, timeout=300).stdout


def git_log_date_range(repo: Path, since: str, until: str) -> str:
    """Run git log --numstat for a date range."""
    fmt = f"Commit: %h | %an | %ad | %s"
    cmd = ["git", "-C", str(repo), "log", "--numstat",
           f"--pretty=format:{fmt}", "--date=short",
           f"--since={since}", f"--until={until}"]
    return subprocess.run(cmd, capture_output=True, text=True, timeout=300).stdout


def iso_week(date_str: str) -> str:
    """Convert 'YYYY-MM-DD' to ISO week string like '2026-W25'."""
    d = datetime.strptime(date_str, "%Y-%m-%d")
    iso = d.isocalendar()
    return f"{iso.year}-W{iso.week:02d}"


def iso_week_range(week: str) -> tuple[str, str]:
    """Given '2026-W25', return (monday_date, sunday_date) as 'YYYY-MM-DD'."""
    from datetime import timedelta
    year, _, wk = week.partition("-W")
    d = datetime.strptime(f"{year}-{int(wk)}-1", "%Y-%W-%w")
    # Python's %W is Monday-based, week 1 is first week with a Monday
    d = d + timedelta(days=1)  # Monday
    return d.strftime("%Y-%m-%d"), (d + timedelta(days=6)).strftime("%Y-%m-%d")


def write_jsonl(path: Path, items: list[dict]):
    """Write list of dicts as JSONL."""
    path.parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w") as f:
        for item in items:
            f.write(json.dumps(item, ensure_ascii=False) + "\n")


def read_jsonl(path: Path) -> list[dict]:
    """Read JSONL into list of dicts."""
    items = []
    with open(path) as f:
        for line in f:
            if line.strip():
                items.append(json.loads(line))
    return items


def read_json(path: Path) -> dict:
    with open(path) as f:
        return json.load(f)


def write_json(path: Path, data):
    path.parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)


def write_text(path: Path, text: str):
    path.parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w") as f:
        f.write(text)
