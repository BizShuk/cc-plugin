"""TypedDict definitions for the stats.json schema."""

from typing import TypedDict, NotRequired


class FileStat(TypedDict):
    adds: int
    dels: int
    file: str


class CommitRecord(TypedDict):
    hash: str
    author: str
    date: str
    subject: str
    raw_lines: int
    kept_lines: int
    kept_files: list[FileStat]


class Top10Commit(TypedDict):
    repo: str
    hash: str
    author: str
    date: str
    subject: str
    raw_lines: int
    kept_lines: int
    files: list[FileStat]


class WeeklyCommit(TypedDict):
    hash: str
    author: str
    date: str
    subject: str
    kept_lines: int


class WeeklyBucket(TypedDict):
    week: str
    start: str
    end: str
    commit_count: int
    kept_lines: int
    commits: list[WeeklyCommit]
    diff_file: str


class CommitterStat(TypedDict):
    author: str
    kept_lines: int
    pct: float
    weekly: dict[str, int]  # "2026-W25": 1240


class DistributionBucket(TypedDict):
    range: str
    count: int
    pct: float


class ExcludedCategory(TypedDict):
    category: str
    lines: int
    pct: float


class StatsJSON(TypedDict):
    repo: str
    generated: str
    summary: dict  # total_commits, total_raw_lines, total_excluded_lines, total_kept_lines, excluded_pct
    percentiles: dict[str, int]
    distribution: list[DistributionBucket]
    excluded_by_category: list[ExcludedCategory]
    top10_commits: list[Top10Commit]
    committers: list[CommitterStat]
    weekly_buckets: list[WeeklyBucket]


EMPTY_STATS: StatsJSON = {
    "repo": "",
    "generated": "",
    "summary": {},
    "percentiles": {},
    "distribution": [],
    "excluded_by_category": [],
    "top10_commits": [],
    "committers": [],
    "weekly_buckets": [],
}
