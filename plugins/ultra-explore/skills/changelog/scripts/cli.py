"""Click CLI: repo-changelog command."""

import click
from pathlib import Path
from discover import discover
from collect import collect
from stats import compute_stats
from diff_filter import generate_diffs
from changelog_render import write_skeleton


def _output_dir(repo_name: str, output: str | None, base_only: bool = False) -> Path:
    if output:
        base = Path(output).expanduser().resolve()
        if base_only:
            return base
        return base / repo_name
    return Path.home() / "projects" / "product" / "projects" / repo_name


@click.group()
def main():
    """repo-changelog: Git history analysis & CHANGELOG generation."""
    pass


@main.command()
@click.argument("root", type=click.Path(exists=True, file_okay=False))
@click.option("--max-depth", default=6, help="Max depth for discovery")
def discover_cmd(root, max_depth):
    """List all nested .git repos under ROOT."""
    repos = discover(Path(root), max_depth)
    for r in repos:
        click.echo(str(r))


@main.command()
@click.argument("repo", type=click.Path(exists=True, file_okay=False))
@click.option("--output", help="Output directory")
def collect_cmd(repo, output):
    """Run git log --numstat → _raw/commits.jsonl."""
    repo_path = Path(repo).resolve()
    out_dir = _output_dir(repo_path.name, output)
    out = collect(repo_path, out_dir)
    click.echo(f"Wrote {out}")


@main.command()
@click.argument("repo", type=click.Path(exists=True, file_okay=False))
@click.option("--output", help="Output directory")
@click.option("--no-gitignore", is_flag=True, help="Skip per-repo .gitignore overlay")
def stats_cmd(repo, output, no_gitignore):
    """Compute stats → stats.json."""
    repo_path = Path(repo).resolve()
    out_dir = _output_dir(repo_path.name, output)
    out = compute_stats(repo_path, out_dir, no_gitignore=no_gitignore)
    click.echo(f"Wrote {out}")


@main.command()
@click.argument("repo", type=click.Path(exists=True, file_okay=False))
@click.option("--output", help="Output directory")
@click.option("--no-gitignore", is_flag=True, help="Skip per-repo .gitignore overlay")
def diff_cmd(repo, output, no_gitignore):
    """Generate filtered diffs per week → _diffs/."""
    repo_path = Path(repo).resolve()
    out_dir = _output_dir(repo_path.name, output)
    paths = generate_diffs(repo_path, out_dir, no_gitignore=no_gitignore)
    click.echo(f"Generated {len(paths)} diff files in {out_dir / '_diffs'}")


@main.command()
@click.argument("repo", type=click.Path(exists=True, file_okay=False))
@click.option("--output", help="Output directory")
@click.option("--no-gitignore", is_flag=True, help="Skip per-repo .gitignore overlay")
def run_cmd(repo, output, no_gitignore):
    """Run all stages (collect + stats + diff) for a single repo."""
    repo_path = Path(repo).resolve()
    out_dir = _output_dir(repo_path.name, output)
    click.echo(f"[1/3] collect: {repo_path.name}")
    collect(repo_path, out_dir)
    click.echo(f"[2/3] stats")
    compute_stats(repo_path, out_dir, no_gitignore=no_gitignore)
    click.echo(f"[3/3] diff")
    generate_diffs(repo_path, out_dir, no_gitignore=no_gitignore)
    # Also write skeleton
    write_skeleton(out_dir / "stats.json", out_dir / "CHANGELOG.md")
    click.echo(f"Done → {out_dir}")


@main.command()
@click.argument("root", type=click.Path(exists=True, file_okay=False))
@click.option("--output", help="Base output directory")
@click.option("--no-gitignore", is_flag=True, help="Skip per-repo .gitignore overlay")
@click.option("--max-depth", default=6, help="Max depth for discovery")
def run_all_cmd(root, output, no_gitignore, max_depth):
    """Run all stages for all repos under ROOT."""
    import sys
    repos = discover(Path(root), max_depth)
    click.echo(f"Found {len(repos)} repos")
    base_dir = _output_dir("", output, base_only=True)
    errors = []
    for i, repo in enumerate(repos, 1):
        click.echo(f"[{i}/{len(repos)}] {repo.name}")
        try:
            out_dir = base_dir / repo.name
            collect(repo, out_dir)
            compute_stats(repo, out_dir, no_gitignore=no_gitignore)
            generate_diffs(repo, out_dir, no_gitignore=no_gitignore)
            write_skeleton(out_dir / "stats.json", out_dir / "CHANGELOG.md")
        except Exception as e:
            click.echo(f"  ERROR: {e}", err=True)
            errors.append({"repo": str(repo), "error": str(e)})
    if errors:
        click.echo(f"\n{len(errors)} errors, see above")
        sys.exit(1)
    click.echo(f"\nDone — {len(repos)} repos processed")


if __name__ == "__main__":
    main()
