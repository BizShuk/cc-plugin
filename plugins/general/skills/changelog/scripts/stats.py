"""Stage 2: exclusion + aggregation → stats.json."""

from collections import defaultdict
from datetime import datetime
from pathlib import Path
from utils import read_jsonl, write_json, iso_week, iso_week_range
from exclusion import is_excluded_default, GitignoreOverlay
from schema import StatsJSON, WeeklyBucket, CommitterStat, Top10Commit


def percentile(sorted_values: list[int], p: float) -> float:
    """Linear interpolation percentile (numpy default method)."""
    n = len(sorted_values)
    if n == 0:
        return 0
    k = (n - 1) * p / 100.0
    f = int(k)
    c = k - f
    if f + 1 < n:
        return sorted_values[f] + c * (sorted_values[f + 1] - sorted_values[f])
    return float(sorted_values[f])


def compute_stats(repo: Path, output_dir: Path,
                  no_gitignore: bool = False) -> Path:
    """Read _raw/commits.jsonl, apply exclusion, compute stats, write stats.json."""
    repo_name = repo.name
    raw_path = output_dir / "_raw" / "commits.jsonl"
    commits = read_jsonl(raw_path)

    gitignore = None if no_gitignore else GitignoreOverlay(repo)

    all_commits = []
    total_raw = 0
    total_excluded = 0
    excluded_cat = defaultdict(int)
    # Per-author tracking
    author_lines = defaultdict(int)
    author_weekly = defaultdict(lambda: defaultdict(int))

    for c in commits:
        raw_lines = 0
        kept_lines = 0
        kept_files = []
        for f in c['files']:
            change = f['adds'] + f['dels']
            raw_lines += change
            exc, reason = is_excluded_default(f['file'], repo_name)
            if not exc and gitignore and gitignore.matches(f['file']):
                exc, reason = True, ".gitignore"
            if exc:
                excluded_cat[reason] += change
                total_excluded += change
            else:
                kept_lines += change
                kept_files.append(f)

        total_raw += raw_lines
        c['raw_lines'] = raw_lines
        c['kept_lines'] = kept_lines
        c['kept_files'] = kept_files
        week = iso_week(c['date'])
        c['week'] = week
        all_commits.append(c)
        author_lines[c['author']] += kept_lines
        author_weekly[c['author']][week] += kept_lines

    total_kept = total_raw - total_excluded
    kept_sorted = sorted([c['kept_lines'] for c in all_commits])
    n = len(kept_sorted)
    ps = {}
    for p in [10, 25, 50, 75, 90, 95, 99]:
        ps[f'P{p}'] = round(percentile(kept_sorted, p))

    # Distribution
    buckets = [(0, 0, "0"), (1, 10, "1-10"), (11, 50, "11-50"),
               (51, 100, "51-100"), (101, 500, "101-500"),
               (501, 1000, "501-1,000"), (1001, 5000, "1,001-5,000"),
               (5001, float('inf'), "5,001+")]
    dist = [{"range": label, "count": sum(1 for v in kept_sorted if lo <= v <= hi),
             "pct": round(100.0 * sum(1 for v in kept_sorted if lo <= v <= hi) / max(n, 1), 1)}
            for lo, hi, label in buckets]

    # Top 10
    top10 = sorted(all_commits, key=lambda c: -c['kept_lines'])[:10]
    top10_out = [{
        "repo": repo_name, "hash": c['hash'], "author": c['author'],
        "date": c['date'], "subject": c['subject'],
        "raw_lines": c['raw_lines'], "kept_lines": c['kept_lines'],
        "files": c['kept_files'],
    } for c in top10]

    # Committers
    total_author_lines = sum(author_lines.values())
    committers = [{
        "author": a, "kept_lines": l,
        "pct": round(100.0 * l / max(total_author_lines, 1), 1),
        "weekly": dict(author_weekly[a]),
    } for a, l in sorted(author_lines.items(), key=lambda x: -x[1])]

    # Weekly buckets
    weeks = defaultdict(list)
    for c in all_commits:
        weeks[c['week']].append(c)
    weekly_buckets = []
    for week in sorted(weeks.keys()):
        wc = weeks[week]
        start, end = iso_week_range(week)
        weekly_buckets.append({
            "week": week, "start": start, "end": end,
            "commit_count": len(wc),
            "kept_lines": sum(c['kept_lines'] for c in wc),
            "commits": [{
                "hash": c['hash'], "author": c['author'],
                "date": c['date'], "subject": c['subject'],
                "kept_lines": c['kept_lines'],
            } for c in wc],
            "diff_file": f"_diffs/{week}.diff",
        })

    # Category breakdown
    cat_breakdown = sorted(excluded_cat.items(), key=lambda x: -x[1])
    excluded_categories = [{
        "category": cat, "lines": lines,
        "pct": round(100.0 * lines / max(total_excluded, 1), 1),
    } for cat, lines in cat_breakdown]

    result: StatsJSON = {
        "repo": repo_name,
        "generated": datetime.now().isoformat(),
        "summary": {
            "total_commits": n,
            "total_raw_lines": total_raw,
            "total_excluded_lines": total_excluded,
            "total_kept_lines": total_kept,
            "excluded_pct": round(100.0 * total_excluded / max(total_raw, 1), 1),
        },
        "percentiles": {"Min": min(kept_sorted) if kept_sorted else 0,
                        **ps,
                        "Max": max(kept_sorted) if kept_sorted else 0,
                        "Average": round(sum(kept_sorted) / max(n, 1), 1)},
        "distribution": dist,
        "excluded_by_category": excluded_categories,
        "top10_commits": top10_out,
        "committers": committers,
        "weekly_buckets": weekly_buckets,
    }

    out = output_dir / "stats.json"
    write_json(out, result)
    return out
