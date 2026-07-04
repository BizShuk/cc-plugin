"""Find all nested .git folders under a root directory."""

from pathlib import Path


def discover(root: Path, max_depth: int = 6) -> list[Path]:
    """Find all .git directories (not submodules) under root.
    Returns list of parent directories (the actual repos)."""
    repos = set()
    for gitdir in root.rglob(".git"):
        if gitdir.is_file():
            continue  # skip submodule pointers (git submodule uses .git file)
        # Check depth: count path parts relative to root
        rel = gitdir.relative_to(root)
        if len(rel.parts) > max_depth:
            continue
        repos.add(gitdir.parent.resolve())
    return sorted(repos)
