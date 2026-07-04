"""Exclusion engine: built-in defaults + per-repo .gitignore overlay."""

import re
from pathlib import Path
from typing import Optional


# --- Layer 1: Built-in defaults (23 rules, generalized from wallet analysis) ---

RE_GEN_DIR = re.compile(r'(^|/)kitex_gen/|(^|/)rpc_gen/|(^|/)thrift_gen/')
RE_CSV_GENERIC = re.compile(r'(^|/)data/.*\.csv$')

CSV_REPO_RULES = [
    ("pns_monorepo", re.compile(r'app/retry_script/(uids|x\d+|songs)/')),
    ("wallet-recharge", re.compile(r'asyncjobs/order_repair/data/orders\.csv')),
    ("wallet-migration-trigger-v2", re.compile(r'config/data/.*test_accounts.*\.csv(\.bak)?$')),
    ("wallet_faas", re.compile(r'top_recharger/biz/html/latest\.csv')),
    ("wallet-region-manager", re.compile(r'conf/testing_account/.*test_accounts.*\.csv')),
    ("wallet-account", re.compile(r'conf/testing_account/.*test_accounts.*\.csv')),
]

RE_PNS_TRUNCATED = re.compile(r'^\.\.\./')
RE_ALLOWLIST = re.compile(r'(^|/)allowlist/')
RE_VENDOR = re.compile(r'(^|/)vendor/')
RE_TEST_MOCK = re.compile(r'(^|/)mock.*\.go$|(^|/).*_test\.go$')
RE_OVERPASS = re.compile(r'overpass_auto_generated|invoker\.go$|virtual_common_idl')
RE_PNS_RETRY = re.compile(r'app/retry_script/(acl|bot|check_0819|clean_)')
RE_DEPS = re.compile(r'(^|/)go\.sum$|(^|/)go\.mod$|(^|/)BUILD\.bazel$|(^|/)go\.work\.sum$|(^|/)go_deps\.bzl$')
RE_KITEX_TRUNCATED = re.compile(r'core_tiktok/(k-[^/]+\.go|[^/]*_core\.go)$')

BW_LIST_RULES = [
    ("wallet-penalty-script", re.compile(r'ban_list\.go$')),
    ("pns_monorepo", re.compile(r'whitelist/source\.go$')),
]

RE_CLAUDE_SKILLS = re.compile(r'(^|/)\.claude/skills/|(^|/)\.ttadk/|(^|/)\.cursor/skills/')
RE_DOC_EXPORT = re.compile(r'(^|/)doc_export/')
RE_DOCS = re.compile(r'(^|/)openspec/|(^|/)docs/superpowers/')
RE_OUTPUT = re.compile(r'(^|/)output/')
RE_COVERAGE = re.compile(r'(^|/)coverage\.out$')
RE_BAZEL_DEPS = re.compile(r'(^|/)deps\.bzl$')
RE_LOG_FILES = re.compile(r'\.log\.\d{4}-\d{2}-\d{2}|(^|/)access\.log$')
RE_DATA_TXT = re.compile(r'(^|/)prod_item_ids\.txt$')
RE_PROTO_GEN = re.compile(r'(^|/)gen/pb/|(^|/)proto_gen/|\.pb\.go$')
RE_SKILLS_DIR = re.compile(r'(^|/)skills/')
RE_RPC_CLIENT = re.compile(r'(^|/)rpc_client/')
RE_THRIFT_RENAME = re.compile(r'\{.*thrift_gen.*=>.*thrift_gen.*\}')

# Fund-specific generated files
RE_FUND_GEN = re.compile(
    r'(^|/)fundsniffer/.*/(rules_v2_library\.go|rules_v2\.go|rules\.go|'
    r'table\.json|output\.yaml|config\.yaml|convert_rules_library\.go|'
    r'convert_rules\.go|filters\.go)$'
)


def is_excluded_default(file_path: str, repo_name: str = "") -> tuple[bool, Optional[str]]:
    """Check built-in exclusion rules. Returns (excluded, reason)."""
    if RE_GEN_DIR.search(file_path):
        return True, "*_gen/"
    if RE_CSV_GENERIC.search(file_path):
        return True, "CSV"
    for rn, pattern in CSV_REPO_RULES:
        if repo_name == rn and pattern.search(file_path):
            return True, "CSV"

    is_pns = (repo_name == "pns_monorepo")
    if is_pns and RE_PNS_TRUNCATED.match(file_path):
        return True, "pns_truncated"
    if RE_ALLOWLIST.search(file_path):
        return True, "allowlist"
    if RE_VENDOR.search(file_path):
        return True, "vendor"
    if RE_TEST_MOCK.search(file_path):
        return True, "test_mock"
    if RE_OVERPASS.search(file_path):
        return True, "overpass"
    if is_pns and RE_PNS_RETRY.search(file_path):
        return True, "pns_retry"
    if RE_DEPS.search(file_path):
        return True, "deps"
    if not is_pns and RE_KITEX_TRUNCATED.search(file_path):
        return True, "kitex_truncated"
    for rn, pattern in BW_LIST_RULES:
        if repo_name == rn and pattern.search(file_path):
            return True, "bw_list"
    if RE_CLAUDE_SKILLS.search(file_path):
        return True, "claude_skills"
    if RE_DOC_EXPORT.search(file_path):
        return True, "doc_export"
    if RE_DOCS.search(file_path):
        return True, "docs"
    if RE_OUTPUT.search(file_path):
        return True, "output"
    if RE_COVERAGE.search(file_path):
        return True, "coverage"
    if RE_BAZEL_DEPS.search(file_path):
        return True, "deps"
    if RE_LOG_FILES.search(file_path):
        return True, "log_files"
    if RE_DATA_TXT.search(file_path):
        return True, "data_txt"
    if RE_PROTO_GEN.search(file_path):
        return True, "proto_gen"
    if RE_SKILLS_DIR.search(file_path):
        return True, "skills"
    if RE_RPC_CLIENT.search(file_path):
        return True, "rpc_client"
    if RE_THRIFT_RENAME.search(file_path):
        return True, "thrift_rename"
    if RE_FUND_GEN.search(file_path):
        return True, "fund_gen"

    return False, None


# --- Layer 2: Per-repo .gitignore ---

class GitignoreOverlay:
    """Parses a .gitignore file and matches paths against it."""

    def __init__(self, repo_path: Path):
        self.spec = None
        gf = repo_path / ".gitignore"
        if not gf.exists():
            return
        try:
            from pathspec import PathSpec
            with open(gf) as f:
                self.spec = PathSpec.from_lines("gitwildmatch", f)
        except ImportError:
            self.spec = None  # fallback handled in is_excluded

    def matches(self, file_path: str) -> bool:
        if self.spec is None:
            return self._check_ignore_fallback(file_path)
        return self.spec.match_file(file_path)

    def _check_ignore_fallback(self, file_path: str) -> bool:
        """Fallback using git check-ignore (slower but always correct)."""
        import subprocess
        r = subprocess.run(
            ["git", "check-ignore", "-q", file_path],
            capture_output=True, timeout=5
        )
        return r.returncode == 0


def is_excluded(file_path: str, repo_name: str = "",
                gitignore: Optional[GitignoreOverlay] = None) -> tuple[bool, Optional[str]]:
    """Full exclusion check: defaults first, then .gitignore overlay."""
    excluded, reason = is_excluded_default(file_path, repo_name)
    if excluded:
        return True, reason
    if gitignore and gitignore.matches(file_path):
        return True, ".gitignore"
    return False, None
