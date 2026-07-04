"""Render non-LLM skeleton of CHANGELOG.md from stats.json."""

from pathlib import Path
from utils import read_json, write_text


def render_skeleton(stats_path: Path) -> str:
    """Read stats.json and return the CHANGELOG.md skeleton (without LLM narratives)."""
    d = read_json(stats_path)
    s = d["summary"]
    lines = []

    # Header
    lines.append(f"# CHANGELOG — {d['repo']}")
    lines.append("")
    lines.append(f"> Generated: {d['generated'][:10]} | Commits: {s['total_commits']:,} | "
                 f"Contributors: {len(d['committers'])}")
    lines.append(f"> Excluded: {s['excluded_pct']}% of raw line changes "
                 f"({s['total_excluded_lines']:,} / {s['total_raw_lines']:,})")
    lines.append("")

    # Committer Activity
    lines.append("## Committer Activity")
    lines.append("")
    lines.append("| Author | Total Kept Lines | % |")
    lines.append("|--------|-----------------|---|")
    for cm in d["committers"]:
        lines.append(f"| {cm['author']} | {cm['kept_lines']:,} | {cm['pct']}% |")
    lines.append("")

    # Weekly breakdown (collapsed)
    all_weeks = sorted(set(w for cm in d["committers"] for w in cm["weekly"]))
    lines.append("<details><summary>Weekly breakdown</summary>")
    lines.append("")
    header = "| Author | " + " | ".join(all_weeks) + " |"
    lines.append(header)
    lines.append("|" + "|".join(["--------"] * (len(all_weeks) + 1)) + "|")
    for cm in d["committers"]:
        vals = " | ".join(str(cm["weekly"].get(w, 0)) for w in all_weeks)
        lines.append(f"| {cm['author']} | {vals} |")
    lines.append("")
    lines.append("</details>")
    lines.append("")

    # Top 10 Commits
    lines.append("## Top 10 Commits (by kept lines)")
    lines.append("")
    lines.append("| # | Hash | Author | Date | Lines | Subject |")
    lines.append("|---|------|--------|------|-------|---------|")
    for i, c in enumerate(d["top10_commits"], 1):
        lines.append(f"| {i} | `{c['hash']}` | {c['author']} | {c['date']} | "
                     f"{c['kept_lines']:,} | {c['subject'][:80]} |")
    lines.append("")

    # Top 10 file lists (collapsed)
    lines.append("<details><summary>Top 10 file lists</summary>")
    lines.append("")
    for i, c in enumerate(d["top10_commits"], 1):
        lines.append(f"**#{i}** `{c['hash']}` — {c['subject'][:80]}")
        lines.append("")
        for f in sorted(c["files"], key=lambda x: -(x["adds"] + x["dels"])):
            lines.append(f"- `{f['file']}` (+{f['adds']} -{f['dels']})")
        lines.append("")
    lines.append("</details>")
    lines.append("")

    # Weekly sections (LLM placeholder)
    lines.append("---")
    lines.append("")
    for bucket in d["weekly_buckets"]:
        lines.append(f"## Week of {bucket['start']} ({bucket['week']}) — "
                     f"{bucket['commit_count']} commits, {bucket['kept_lines']:,} kept lines")
        lines.append("")
        lines.append("<!-- LLM: 3-5 sentence narrative summarizing the feature changes "
                     "visible in the diff for this week. Write in past tense, mention "
                     "key files/modules changed, and focus on what changed from a user/API "
                     "perspective, not implementation details. -->")
        lines.append("")
        lines.append("### Commits")
        lines.append("")
        lines.append("| Hash | Author | Date | Lines | Subject |")
        lines.append("|------|--------|------|-------|---------|")
        for c in bucket["commits"]:
            lines.append(f"| `{c['hash']}` | {c['author']} | {c['date']} | "
                         f"{c['kept_lines']:,} | {c['subject']} |")
        lines.append("")
        lines.append("---")
        lines.append("")

    return "\n".join(lines)


def write_skeleton(stats_path: Path, output_path: Path) -> Path:
    """Write CHANGELOG skeleton to output_path."""
    content = render_skeleton(stats_path)
    write_text(output_path, content)
    return output_path
