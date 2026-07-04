#!/usr/bin/env python3
"""kb_history — deterministic git-history pipeline for the ultra-explore plugin.

Single-file, stdlib-only (no pip install). Redeveloped from the general
plugin's `changelog` skill (repo-changelog) with generalized exclusion rules.

Per repo it writes (default output: ~/projects/product/projects/<repo_name>/):
    _raw/commits.jsonl   hash/author/date/subject + numstat per commit
    stats.json           summary, percentiles, committers, weekly buckets
    _diffs/<week>.diff   per-ISO-week diff filtered to kept files
    CHANGELOG.md         skeleton with <!-- LLM: ... --> weekly placeholders

Usage:
    python3 kb_history.py run <repo> [--output DIR] [--no-gitignore]
    python3 kb_history.py discover <root> [--max-depth N]
"""

import argparse
import json
import re
import subprocess
import sys
from collections import defaultdict
from datetime import datetime, timedelta
from pathlib import Path

# ---------------------------------------------------------------- exclusion

# Generic noise rules (kept lines/diffs exclude these; .gitignore overlay on top)
EXCLUDE_RULES = [
    ("vendored",  re.compile(r'(^|/)(vendor|node_modules|third_party)/')),
    ("generated", re.compile(r'(^|/)(gen|proto_gen|kitex_gen|rpc_gen|thrift_gen|dist|build|out|output|target)/')),
    ("generated", re.compile(r'\.(pb|gen)\.(go|py|ts|js)$|\.min\.(js|css)$')),
    ("deps-lock", re.compile(r'(^|/)(go\.sum|go\.mod|go\.work\.sum|package-lock\.json|yarn\.lock|pnpm-lock\.yaml|poetry\.lock|Cargo\.lock|Gemfile\.lock|deps\.bzl|go_deps\.bzl|BUILD\.bazel)$')),
    ("test-mock", re.compile(r'(^|/)(mocks?|fakes?|fixtures|testdata)/|_test\.go$|\.spec\.(ts|js|tsx|jsx)$|_test\.py$')),
    ("agent-cfg", re.compile(r'(^|/)(\.claude|\.cursor|\.ttadk|\.gemini|\.agent)/')),
    ("data-file", re.compile(r'\.(csv|jsonl|parquet|sqlite|db|log)$|\.log\.\d{4}-\d{2}-\d{2}')),
    ("coverage",  re.compile(r'(^|/)coverage\.(out|xml|html)$|(^|/)\.coverage$')),
    ("binary",    re.compile(r'\.(png|jpe?g|gif|pdf|zip|tar|gz|ico|woff2?|ttf|mp[34])$')),
]


def exclusion_reason(file_path: str):
    for reason, rx in EXCLUDE_RULES:
        if rx.search(file_path):
            return reason
    return None


class GitignoreOverlay:
    """Batch .gitignore matcher via `git check-ignore --stdin` (no deps)."""

    def __init__(self, repo: Path, enabled: bool = True):
        self.repo = repo
        self.enabled = enabled
        self._cache = {}

    def prime(self, paths):
        """Resolve many paths in one subprocess call."""
        todo = [p for p in set(paths) if p not in self._cache]
        if not (self.enabled and todo):
            return
        r = subprocess.run(
            ["git", "-C", str(self.repo), "check-ignore", "--stdin"],
            input="\n".join(todo), capture_output=True, text=True, errors="replace", timeout=60)
        ignored = set(r.stdout.splitlines())
        for p in todo:
            self._cache[p] = p in ignored

    def matches(self, file_path: str) -> bool:
        if not self.enabled:
            return False
        if file_path not in self._cache:
            self.prime([file_path])
        return self._cache.get(file_path, False)


def is_excluded(file_path: str, gitignore: GitignoreOverlay):
    reason = exclusion_reason(file_path)
    if reason:
        return True, reason
    if gitignore.matches(file_path):
        return True, ".gitignore"
    return False, None

# -------------------------------------------------------------------- utils


def iso_week(date_str: str) -> str:
    iso = datetime.strptime(date_str, "%Y-%m-%d").isocalendar()
    return f"{iso.year}-W{iso.week:02d}"


def iso_week_range(week: str):
    year_str, _, wk_str = week.partition("-W")
    year, wk = int(year_str), int(wk_str)
    jan4 = datetime(year, 1, 4)
    iso = jan4.isocalendar()
    monday = jan4 - timedelta(days=iso.weekday - 1) + timedelta(weeks=wk - iso.week)
    return monday.strftime("%Y-%m-%d"), (monday + timedelta(days=6)).strftime("%Y-%m-%d")


def write_text(path: Path, text: str):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text)


def percentile(sorted_values, p):
    n = len(sorted_values)
    if n == 0:
        return 0
    k = (n - 1) * p / 100.0
    f = int(k)
    c = k - f
    if f + 1 < n:
        return sorted_values[f] + c * (sorted_values[f + 1] - sorted_values[f])
    return float(sorted_values[f])

# ---------------------------------------------------------- stage 1 collect

COMMIT_RE = re.compile(r'^Commit: (\S+) \| (.+?) \| (\d{4}-\d{2}-\d{2}) \| (.*)$')


def collect(repo: Path, out_dir: Path) -> list:
    """git log --numstat --first-parent → _raw/commits.jsonl (oldest first)."""
    cmd = ["git", "-C", str(repo), "log", "--numstat", "--first-parent",
           "--reverse", "--pretty=format:Commit: %h | %an | %ad | %s",
           "--date=short"]
    text = subprocess.run(cmd, capture_output=True, text=True, errors="replace", timeout=600).stdout

    commits, current = [], None
    for line in text.split("\n"):
        m = COMMIT_RE.match(line)
        if m:
            if current:
                commits.append(current)
            current = {"hash": m.group(1), "author": m.group(2),
                       "date": m.group(3), "subject": m.group(4), "files": []}
            continue
        if current and line.strip():
            parts = line.split("\t")
            if len(parts) >= 3:
                try:
                    adds = int(parts[0].replace("-", "0"))
                    dels = int(parts[1].replace("-", "0"))
                except ValueError:
                    continue
                current["files"].append({"adds": adds, "dels": dels, "file": parts[2]})
    if current:
        commits.append(current)

    out = out_dir / "_raw" / "commits.jsonl"
    out.parent.mkdir(parents=True, exist_ok=True)
    with open(out, "w") as f:
        for c in commits:
            f.write(json.dumps(c, ensure_ascii=False) + "\n")
    return commits

# ------------------------------------------------------------ stage 2 stats


def compute_stats(repo: Path, out_dir: Path, commits: list,
                  gitignore: GitignoreOverlay) -> dict:
    gitignore.prime([f["file"] for c in commits for f in c["files"]])

    total_raw = total_excluded = 0
    excluded_cat = defaultdict(int)
    author_lines = defaultdict(int)
    author_weekly = defaultdict(lambda: defaultdict(int))

    for c in commits:
        raw = kept = 0
        kept_files = []
        for f in c["files"]:
            change = f["adds"] + f["dels"]
            raw += change
            exc, reason = is_excluded(f["file"], gitignore)
            if exc:
                excluded_cat[reason] += change
                total_excluded += change
            else:
                kept += change
                kept_files.append(f)
        total_raw += raw
        c["raw_lines"], c["kept_lines"], c["kept_files"] = raw, kept, kept_files
        c["week"] = iso_week(c["date"])
        author_lines[c["author"]] += kept
        author_weekly[c["author"]][c["week"]] += kept

    kept_sorted = sorted(c["kept_lines"] for c in commits)
    n = len(kept_sorted)
    dist_buckets = [(0, 0, "0"), (1, 10, "1-10"), (11, 50, "11-50"),
                    (51, 100, "51-100"), (101, 500, "101-500"),
                    (501, 1000, "501-1,000"), (1001, 5000, "1,001-5,000"),
                    (5001, float("inf"), "5,001+")]

    top10 = sorted(commits, key=lambda c: -c["kept_lines"])[:10]
    weeks = defaultdict(list)
    for c in commits:
        weeks[c["week"]].append(c)

    total_author = sum(author_lines.values())
    stats = {
        "repo": repo.name,
        "generated": datetime.now().isoformat(),
        "summary": {
            "total_commits": n,
            "total_raw_lines": total_raw,
            "total_excluded_lines": total_excluded,
            "total_kept_lines": total_raw - total_excluded,
            "excluded_pct": round(100.0 * total_excluded / max(total_raw, 1), 1),
        },
        "percentiles": {
            "Min": kept_sorted[0] if kept_sorted else 0,
            **{f"P{p}": round(percentile(kept_sorted, p)) for p in (10, 25, 50, 75, 90, 95, 99)},
            "Max": kept_sorted[-1] if kept_sorted else 0,
            "Average": round(sum(kept_sorted) / max(n, 1), 1),
        },
        "distribution": [
            {"range": label,
             "count": sum(1 for v in kept_sorted if lo <= v <= hi),
             "pct": round(100.0 * sum(1 for v in kept_sorted if lo <= v <= hi) / max(n, 1), 1)}
            for lo, hi, label in dist_buckets],
        "excluded_by_category": [
            {"category": cat, "lines": lines,
             "pct": round(100.0 * lines / max(total_excluded, 1), 1)}
            for cat, lines in sorted(excluded_cat.items(), key=lambda x: -x[1])],
        "top10_commits": [
            {"repo": repo.name, "hash": c["hash"], "author": c["author"],
             "date": c["date"], "subject": c["subject"],
             "raw_lines": c["raw_lines"], "kept_lines": c["kept_lines"],
             "files": c["kept_files"]} for c in top10],
        "committers": [
            {"author": a, "kept_lines": l,
             "pct": round(100.0 * l / max(total_author, 1), 1),
             "weekly": dict(author_weekly[a])}
            for a, l in sorted(author_lines.items(), key=lambda x: -x[1])],
        "weekly_buckets": [
            {"week": w, "start": iso_week_range(w)[0], "end": iso_week_range(w)[1],
             "commit_count": len(wc),
             "kept_lines": sum(c["kept_lines"] for c in wc),
             "commits": [{"hash": c["hash"], "author": c["author"], "date": c["date"],
                          "subject": c["subject"], "kept_lines": c["kept_lines"]}
                         for c in wc],
             "diff_file": f"_diffs/{w}.diff"}
            for w, wc in sorted(weeks.items())],
    }
    out = out_dir / "stats.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(stats, indent=2, ensure_ascii=False))
    return stats

# ------------------------------------------------------------ stage 3 diffs

DIFF_GIT_RE = re.compile(r'^diff --git a/(.+?) b/.*$')
PATCH_FILE_RE = re.compile(r'^[-+]{3} [ab]/(.+)$')


def filter_diff(repo: Path, since: str, until: str, gitignore: GitignoreOverlay) -> str:
    cmd = ["git", "-C", str(repo), "log", "-p", "--first-parent",
           f"--since={since}", f"--until={until}"]
    text = subprocess.run(cmd, capture_output=True, text=True, errors="replace", timeout=600).stdout
    out, excluded = [], False
    paths = [m.group(1) for line in text.split("\n") if (m := DIFF_GIT_RE.match(line))]
    gitignore.prime(paths)
    for line in text.split("\n"):
        m = DIFF_GIT_RE.match(line)
        if m:
            excluded = is_excluded(m.group(1), gitignore)[0]
            if not excluded:
                out.append(line)
            continue
        if excluded:
            continue
        out.append(line)
    return "\n".join(out)


def generate_diffs(repo: Path, out_dir: Path, stats: dict, gitignore: GitignoreOverlay):
    for bucket in stats["weekly_buckets"]:
        filtered = filter_diff(repo, bucket["start"], bucket["end"], gitignore)
        write_text(out_dir / "_diffs" / f"{bucket['week']}.diff", filtered)

# --------------------------------------------------------- stage 4 skeleton

LLM_PLACEHOLDER = ("<!-- LLM: 3-5 sentence narrative summarizing the feature changes "
                   "visible in the diff for this week. Write in past tense, mention "
                   "key files/modules changed, and focus on what changed from a user/API "
                   "perspective, not implementation details. -->")


def write_skeleton(stats: dict, output_path: Path):
    s = stats["summary"]
    L = [f"# CHANGELOG — {stats['repo']}", "",
         f"> Generated: {stats['generated'][:10]} | Commits: {s['total_commits']:,} | "
         f"Contributors: {len(stats['committers'])}",
         f"> Excluded: {s['excluded_pct']}% of raw line changes "
         f"({s['total_excluded_lines']:,} / {s['total_raw_lines']:,})", "",
         "## Committer Activity", "",
         "| Author | Total Kept Lines | % |", "|--------|-----------------|---|"]
    for cm in stats["committers"]:
        L.append(f"| {cm['author']} | {cm['kept_lines']:,} | {cm['pct']}% |")
    L += ["", "## Top 10 Commits (by kept lines)", "",
          "| # | Hash | Author | Date | Lines | Subject |",
          "|---|------|--------|------|-------|---------|"]
    for i, c in enumerate(stats["top10_commits"], 1):
        L.append(f"| {i} | `{c['hash']}` | {c['author']} | {c['date']} | "
                 f"{c['kept_lines']:,} | {c['subject'][:80]} |")
    L += ["", "---", ""]
    for b in stats["weekly_buckets"]:
        L += [f"## Week of {b['start']} ({b['week']}) — "
              f"{b['commit_count']} commits, {b['kept_lines']:,} kept lines", "",
              LLM_PLACEHOLDER, "", "### Commits", "",
              "| Hash | Author | Date | Lines | Subject |",
              "|------|--------|------|-------|---------|"]
        for c in b["commits"]:
            L.append(f"| `{c['hash']}` | {c['author']} | {c['date']} | "
                     f"{c['kept_lines']:,} | {c['subject']} |")
        L += ["", "---", ""]
    write_text(output_path, "\n".join(L))

# ----------------------------------------------------------------- commands


def discover(root: Path, max_depth: int = 6):
    """Nested git repos under root (skips submodule .git files)."""
    repos = set()
    for gitdir in root.rglob(".git"):
        if gitdir.is_file():
            continue
        if len(gitdir.relative_to(root).parts) > max_depth:
            continue
        repos.add(gitdir.parent.resolve())
    return sorted(repos)


def cmd_run(repo: Path, output: str | None, no_gitignore: bool):
    out_dir = (Path(output).expanduser().resolve() if output
               else Path.home() / "projects" / "product" / "projects" / repo.name)
    gitignore = GitignoreOverlay(repo, enabled=not no_gitignore)
    print(f"[1/4] collect: {repo.name}")
    commits = collect(repo, out_dir)
    print(f"[2/4] stats ({len(commits)} commits)")
    stats = compute_stats(repo, out_dir, commits, gitignore)
    print(f"[3/4] diffs ({len(stats['weekly_buckets'])} weeks)")
    generate_diffs(repo, out_dir, stats, gitignore)
    print("[4/4] skeleton")
    write_skeleton(stats, out_dir / "CHANGELOG.md")
    print(f"Done → {out_dir}")


def main():
    p = argparse.ArgumentParser(prog="kb_history.py", description=__doc__)
    sub = p.add_subparsers(dest="cmd", required=True)

    r = sub.add_parser("run", help="run collect+stats+diffs+skeleton for one repo")
    r.add_argument("repo", type=Path)
    r.add_argument("--output", help="output dir (default ~/projects/product/projects/<repo>)")
    r.add_argument("--no-gitignore", action="store_true")

    d = sub.add_parser("discover", help="list nested git repos under a root")
    d.add_argument("root", type=Path)
    d.add_argument("--max-depth", type=int, default=6)

    a = p.parse_args()
    if a.cmd == "run":
        if not (a.repo / ".git").exists():
            sys.exit(f"not a git repo: {a.repo}")
        cmd_run(a.repo.resolve(), a.output, a.no_gitignore)
    elif a.cmd == "discover":
        for r_ in discover(a.root.resolve(), a.max_depth):
            print(r_)


if __name__ == "__main__":
    main()
