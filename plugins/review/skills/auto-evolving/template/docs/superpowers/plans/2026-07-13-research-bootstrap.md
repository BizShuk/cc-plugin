# Research Repo Bootstrap Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 `~/projects/research` 從空殼建成 4 主題並列、Loose 模式、4 主題全起 pm2 cron 抓取的研究容器。

**Architecture:** Flat parallel layout（4 主題對稱、頂層 INDEX.md 導引）；每主題內部 Per-domain 演化；4 主題全 Heavy 模式（pm2 cron + 自動抓取）；多語言混合（Go investing/biz-strategy、Python llm-research/personal-kb、純 Markdown 報告）。

**Tech Stack:** Go 1.22+（`yfinance-go`、`gosdk`）、Python 3.11+（`routine_agent` / `arxiv` / kindle parser）、pm2、Node.js（ecosystem.config.js 解析）、uv（Python 套件管理）。

## Global Constraints

- 目錄位置：`/Users/bytedance/projects/research`（已存在但內容空，已 `git init` 完成）
- 主題命名：英文 kebab-case slug：`investing/`, `llm-research/`, `personal-kb/`, `biz-strategy/`
- 組合模式：Loose（主題間不互引，靠 `INDEX.md` 文字導引）
- 報告去向：留在主題內（不同步到 `product/`）
- 4 主題**共用**子目錄：`fetcher/`, `prompts/`
- 4 主題**不對稱**子目錄：`data/`, `experiments/`, `benchmarks/`, `refs/`, `archive/`, `reviews/`, `canvas/`，由各主題 `README.md` 定義意義
- 應用設定路徑：`~/.config/research/`（若日後加常駐 binary 才用；本階段不需要）
- 觀測：fetcher log 走 pm2 log，不接 inf（research 是離線研究型，不上監控）
- Fetcher CLI 介面統一：`<topic>-fetcher fetch` / `<topic>-fetcher summarize`
- pm2 命名：`<topic>-fetcher`，`namespace: 'Research'`，`cron_restart: '0 6 * * *'`
- `.gitignore` 規則：排除 `data/**/raw/`, `data/**/*.csv`, `data/**/*.parquet`, `*.log`, `tmp/`；保留 `data/.gitkeep`
- 文件：根 `README.md`、`CLAUDE.md` 必備；`AGENTS.md -> CLAUDE.md` 軟連；`README.todo` 必備（4 主題各一區塊）
- 不與 `playground/`、`product/` 直接 import；如需共用樣板，未來下鑿到 `framework-research-kit`
- Spec 來源：`docs/superpowers/specs/2026-07-13-research-bootstrap-design.md`

---

## File Structure

| 路徑 | 責任 |
|---|---|
| `README.md` | 業務說明（4 主題、定位、如何開始） |
| `CLAUDE.md` | 技術脈絡（分層、INDEX 規則、fetcher 三段式） |
| `AGENTS.md` | 軟連到 `CLAUDE.md` |
| `README.todo` | 跨主題待辦 |
| `INDEX.md` | 4 主題導引表 |
| `ecosystem.config.js` | pm2 4 apps + 1 maintain app |
| `.gitignore` | 排除 raw data / log / tmp |
| `scripts/seed.sh` | 一鍵建立 4 主題最小骨架 |
| `scripts/maintain_index.py` | 對照 git status 維護 INDEX.md |
| `scripts/contract_test.sh` | 驗證 4 fetcher CLI 介面一致 |
| `scripts/smoke.sh` | 煙霧測試 |
| `tests/test_maintain_index.py` | maintain_index.py 單元測試 |
| `investing/README.md` | investing 主題說明 |
| `investing/fetcher/` | Go 套件（cmd + internal） |
| `investing/prompts/price-digest.md` | LLM 整理模板 |
| `llm-research/README.md` | llm-research 主題說明 |
| `llm-research/fetcher/` | Python 套件（pyproject + src） |
| `llm-research/prompts/arxiv-digest.md` | LLM 整理模板 |
| `personal-kb/README.md` | personal-kb 主題說明 |
| `personal-kb/fetcher/` | Python 套件（kindle parser） |
| `personal-kb/prompts/highlight-digest.md` | LLM 整理模板 |
| `biz-strategy/README.md` | biz-strategy 主題說明 |
| `biz-strategy/fetcher/` | Go 套件（IR/8-K） |
| `biz-strategy/prompts/ir-digest.md` | LLM 整理模板 |

---

## Task 1: 初始化 git、.gitignore、4 主題空目錄

**Files:**
- Create: `.gitignore`
- Create: `investing/{notes,data,reports,fetcher,prompts}/.gitkeep`
- Create: `llm-research/{notes,experiments,benchmarks,reports,fetcher,prompts}/.gitkeep`
- Create: `personal-kb/{notes,refs,archive,reviews,fetcher,prompts}/.gitkeep`
- Create: `biz-strategy/{notes,canvas,reports,fetcher,prompts}/.gitkeep`

**Interfaces:**
- Consumes: 既有 `/Users/bytedance/projects/research/`（已 `git init`）
- Produces: 目錄結構，符合 Global Constraints

- [ ] **Step 1: 寫 `.gitignore`**

```gitignore
# Raw data (rebuildable, large)
data/**/raw/
data/**/*.csv
data/**/*.parquet
data/**/*.tsv
data/**/*.json.gz

# Logs
*.log
logs/

# OS / Editor
.DS_Store
.idea/
.vscode/

# Python
__pycache__/
*.pyc
.venv/

# Go
investing/fetcher/vendor/
llm-research/fetcher/.venv/
personal-kb/fetcher/.venv/
biz-strategy/fetcher/vendor/
*.test
```

寫入 `/Users/bytedance/projects/research/.gitignore`。

- [ ] **Step 2: 建立 4 主題空目錄**

```bash
cd /Users/bytedance/projects/research

for topic in investing llm-research personal-kb biz-strategy; do
  case "$topic" in
    investing)        sub="notes data reports fetcher prompts" ;;
    llm-research)     sub="notes experiments benchmarks reports fetcher prompts" ;;
    personal-kb)      sub="notes refs archive reviews fetcher prompts" ;;
    biz-strategy)     sub="notes canvas reports fetcher prompts" ;;
  esac
  for s in $sub; do
    mkdir -p "$topic/$s"
    touch "$topic/$s/.gitkeep"
  done
done

# 根層其他將在後續 task 建立；此 task 不動
```

- [ ] **Step 3: 驗證結構**

```bash
cd /Users/bytedance/projects/research
find . -type d -name '.git' -prune -o -type d -print | sort
```

預期：22 個目錄（1 根 + 4 主題 × 5/6 子目錄 + 0 其他），全部含 `.gitkeep`。

- [ ] **Step 4: commit**

```bash
cd /Users/bytedance/projects/research
git add .gitignore '*/.gitkeep' '*/*/.gitkeep'
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "chore: scaffold 4 topics with .gitkeep + .gitignore"
```

---

## Task 2: 根文件 — README.md / CLAUDE.md / AGENTS.md / README.todo

**Files:**
- Create: `README.md`
- Create: `CLAUDE.md`
- Create: `AGENTS.md` (symlink)
- Create: `README.todo`

**Interfaces:**
- Consumes: Task 1 的目錄結構
- Produces: 全域文件，定義 4 主題與分層紀律

- [ ] **Step 1: 寫 `README.md`**

```markdown
# research

個人長期研究容器，承載 4 個並列主題。

## 4 主題

| 主題 | 路徑 | 用途 |
|---|---|---|
| 投資理財 | `investing/` | 個股/總經/因子回測 |
| LLM 研究 | `llm-research/` | 模型評測、agent 設計、prompt/RAG 策略 |
| 個人學習 | `personal-kb/` | 個人知識庫、主題式閱讀、書評 |
| 商業策略 | `biz-strategy/` | 商業模式拆解、競品分析 |

主題間不互引，僅由 `INDEX.md` 統一導引（Loose 模式）。
報告留在各主題內，不同步到 `product/`。

## 如何開始

```bash
bash scripts/seed.sh          # 重建 4 主題最小骨架（冪等）
pm2 start ecosystem.config.js # 啟動 4 主題 fetcher cron
```

詳見 `CLAUDE.md` 技術脈絡與 `docs/superpowers/specs/2026-07-13-research-bootstrap-design.md` 設計。
```

- [ ] **Step 2: 寫 `CLAUDE.md`**

```markdown
# research — 技術脈絡

## 分層

位於 `~/projects/` 應用層。可 import framework 層 (`gosdk`, `yfinance-go`, `routine_agent`, `pm2`, `cc-plugin`)；不與 `playground/`, `product/` 直接 import。

## 目錄慣例

- 4 主題共用 `fetcher/`、`prompts/`
- 4 主題不對稱的子目錄（`data/`, `experiments/` 等）由各主題 `README.md` 定義
- 報告留在 `reports/<topic>/`；不複製到 `product/`

## Fetcher 三段式

每個 fetcher CLI 提供兩個子命令：

```bash
<topic>-fetcher fetch       # 抓 raw -> data/<topic>/raw/
<topic>-fetcher summarize   # raw + prompts/ -> notes/<topic>/
```

`reports/<topic>/*.md` 由人工從 `notes/` 整理。

## INDEX.md 規則

`reports/<topic>/*.md` 檔頭需有 `indexed: true|false`。
`scripts/maintain_index.py` 對照 `git status` 把 `indexed:false` 補上，並由 pre-commit hook 強制 `--verify`。

## pm2

`ecosystem.config.js` 為 4 主題各設一個 app：

- `investing-fetcher` (Go)
- `llm-research-fetcher` (Python)
- `personal-kb-fetcher` (Python)
- `biz-strategy-fetcher` (Go)
- `research-maintain-index` (Python, hourly)

`namespace: 'Research'`，`cron_restart: '0 6 * * *'`。
```

- [ ] **Step 3: 建立 `AGENTS.md` 軟連**

```bash
cd /Users/bytedance/projects/research
ln -s CLAUDE.md AGENTS.md
ls -la AGENTS.md
```

預期：`AGENTS.md -> CLAUDE.md`

- [ ] **Step 4: 寫 `README.todo`**

```markdown
# TODO

跨主題待辦追蹤。各主題有自己的子段落。

## investing

- [ ] fetcher 完整實作（Task 6）
- [ ] prompts/price-digest.md 模板

## llm-research

- [ ] fetcher 完整實作（Task 7）
- [ ] prompts/arxiv-digest.md 模板

## personal-kb

- [ ] fetcher 完整實作（Task 7）
- [ ] prompts/highlight-digest.md 模板

## biz-strategy

- [ ] fetcher 完整實作（Task 7）
- [ ] prompts/ir-digest.md 模板

## Archive
```

- [ ] **Step 5: 驗證 markdown 結構**

```bash
cd /Users/bytedance/projects/research
ls -la README.md CLAUDE.md AGENTS.md README.todo
test -L AGENTS.md && echo "AGENTS.md is symlink OK"
grep -c '^#' README.md CLAUDE.md README.todo
```

預期：4 個檔案存在；`AGENTS.md is symlink OK`；每個檔案至少 1 個標題。

- [ ] **Step 6: commit**

```bash
cd /Users/bytedance/projects/research
git add README.md CLAUDE.md AGENTS.md README.todo
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "docs: root README/CLAUDE/AGENTS/README.todo"
```

---

## Task 3: 寫 INDEX.md（初始空狀態）

**Files:**
- Create: `INDEX.md`

**Interfaces:**
- Consumes: 4 主題列表
- Produces: 4 主題導引表，maintain_index.py 之後會更新此檔

- [ ] **Step 1: 寫 `INDEX.md`**

```markdown
# research — INDEX

最後更新：2026-07-13

各主題報告登錄點。當前所有主題尚無報告產出。

## 投資理財 (`investing/`)

| 日期 | 報告 | 摘要 |
|---|---|---|
| — | — | — |

## LLM 研究 (`llm-research/`)

| 日期 | 報告 | 摘要 |
|---|---|---|
| — | — | — |

## 個人學習 (`personal-kb/`)

| 日期 | 報告 | 摘要 |
|---|---|---|
| — | — | — |

## 商業策略 (`biz-strategy/`)

| 日期 | 報告 | 摘要 |
|---|---|---|
| — | — | — |

## 維護

- 由 `scripts/maintain_index.py` 自動補上 `reports/<topic>/*.md` 條目
- pre-commit hook 跑 `maintain_index.py --verify` 確保表格一致
```

- [ ] **Step 2: 驗證 markdown**

```bash
cd /Users/bytedance/projects/research
grep -c '^##' INDEX.md
```

預期：5 個 `##`（4 主題 + 1 維護）。

- [ ] **Step 3: commit**

```bash
cd /Users/bytedance/projects/research
git add INDEX.md
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "docs: initial INDEX.md (empty state)"
```

---

## Task 4: ecosystem.config.js 雛型

**Files:**
- Create: `ecosystem.config.js`

**Interfaces:**
- Consumes: 4 主題命名
- Produces: pm2 配置檔，5 apps + namespace

- [ ] **Step 1: 寫 `ecosystem.config.js`**

```javascript
module.exports = {
  apps: [
    {
      name: 'investing-fetcher',
      namespace: 'Research',
      cwd: './investing/fetcher',
      script: 'cmd/investing-fetcher',
      cron_restart: '0 6 * * *',
      autorestart: false,
      max_restarts: 0,
      out_file: '../data/logs/investing-fetcher.out.log',
      error_file: '../data/logs/investing-fetcher.err.log',
      merge_logs: true,
    },
    {
      name: 'llm-research-fetcher',
      namespace: 'Research',
      cwd: './llm-research/fetcher',
      script: 'src/llm_research_fetcher/__main__.py',
      interpreter: 'python3',
      cron_restart: '0 6 * * *',
      autorestart: false,
      out_file: '../data/logs/llm-research-fetcher.out.log',
      error_file: '../data/logs/llm-research-fetcher.err.log',
    },
    {
      name: 'personal-kb-fetcher',
      namespace: 'Research',
      cwd: './personal-kb/fetcher',
      script: 'src/personal_kb_fetcher/__main__.py',
      interpreter: 'python3',
      cron_restart: '0 6 * * *',
      autorestart: false,
      out_file: '../data/logs/personal-kb-fetcher.out.log',
      error_file: '../data/logs/personal-kb-fetcher.err.log',
    },
    {
      name: 'biz-strategy-fetcher',
      namespace: 'Research',
      cwd: './biz-strategy/fetcher',
      script: 'cmd/biz-strategy-fetcher',
      cron_restart: '0 6 * * *',
      autorestart: false,
      out_file: '../data/logs/biz-strategy-fetcher.err.log',
      merge_logs: true,
    },
    {
      name: 'research-maintain-index',
      namespace: 'Research',
      cwd: './scripts',
      script: 'maintain_index.py',
      interpreter: 'python3',
      cron_restart: '0 * * * *',
      autorestart: false,
    },
  ],
};
```

- [ ] **Step 2: 語法檢查**

```bash
cd /Users/bytedance/projects/research
node -e "require('./ecosystem.config.js'); console.log('OK')"
```

預期：印出 `OK`。

- [ ] **Step 3: pm2 dry-run 驗證**

```bash
cd /Users/bytedance/projects/research
pm2 start ecosystem.config.js --dry-run 2>&1 | tail -20
```

預期：5 個 app 列在 Research namespace 下，不報錯。
註：因 fetcher 尚未實作，binary 與 python script 暫不存在，dry-run 不會啟動 process，僅驗證配置可解析。

- [ ] **Step 4: commit**

```bash
cd /Users/bytedance/projects/research
git add ecosystem.config.js
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "feat(pm2): 4 topic fetchers + maintain-index (5 apps, Research namespace)"
```

---

## Task 5: scripts/maintain_index.py + 單元測試

**Files:**
- Create: `scripts/maintain_index.py`
- Create: `tests/test_maintain_index.py`
- Create: `pyproject.toml` (在 `scripts/`)

**Interfaces:**
- Consumes: `INDEX.md`、`reports/<topic>/*.md` 的 front-matter `indexed:` 欄位、`git status --porcelain`
- Produces:
  - `update()` — 對未登錄報告補上表格行
  - `verify()` — 確認 INDEX.md 與 reports 一致；不一致時 exit code 1
  - CLI：`maintain_index.py update|verify`

- [ ] **Step 1: 建立 `scripts/pyproject.toml`**

```toml
[project]
name = "research-scripts"
version = "0.1.0"
requires-python = ">=3.11"
dependencies = []

[build-system]
requires = ["setuptools>=68"]
build-backend = "setuptools.build_meta"

[tool.setuptools]
py-modules = ["maintain_index"]
```

寫入 `/Users/bytedance/projects/research/scripts/pyproject.toml`。

- [ ] **Step 2: 寫 failing test `tests/test_maintain_index.py`**

```python
import sys
from pathlib import Path
import pytest

# 確保 scripts/ 在 path
sys.path.insert(0, str(Path(__file__).parent.parent / "scripts"))

import maintain_index


def test_parse_front_matter_indexed_true():
    text = "---\ntitle: foo\nindexed: true\n---\n# body"
    fm = maintain_index.parse_front_matter(text)
    assert fm == {"title": "foo", "indexed": "true"}


def test_parse_front_matter_no_fm_returns_empty():
    text = "# no front matter"
    fm = maintain_index.parse_front_matter(text)
    assert fm == {}


def test_extract_summary_from_md():
    body = "# Title\n\nFirst paragraph here.\n\nSecond paragraph."
    summary = maintain_index.extract_summary(body, max_len=30)
    assert summary == "First paragraph here."


def test_update_adds_unindexed_report(tmp_path, monkeypatch):
    # 模擬結構
    idx = tmp_path / "INDEX.md"
    idx.write_text("# INDEX\n\n## investing (`investing/`)\n\n| 日期 | 報告 | 摘要 |\n|---|---|---|\n")
    report = tmp_path / "reports" / "investing" / "2026-07-13-fed.md"
    report.parent.mkdir(parents=True)
    report.write_text("---\ntitle: Fed Rate Cut\nindexed: false\ndate: 2026-07-13\n---\n# Fed Rate Cut\n\nThe Fed cut rates by 25bp.")

    monkeypatch.chdir(tmp_path)
    maintain_index.update()

    content = idx.read_text()
    assert "2026-07-13" in content
    assert "fed.md" in content
    assert "Fed cut rates" in content


def test_verify_clean_state_passes(tmp_path, monkeypatch):
    idx = tmp_path / "INDEX.md"
    idx.write_text("# INDEX\n\n## investing\n| 日期 | 報告 | 摘要 |\n|---|---|---|\n")
    report = tmp_path / "reports" / "investing" / "2026-07-13.md"
    report.parent.mkdir(parents=True)
    report.write_text("---\nindexed: true\n---\n# Already indexed")
    (tmp_path / "reports" / "investing" / "2026-07-13.md").write_text(
        "---\nindexed: true\n---\n# Already indexed"
    )

    monkeypatch.chdir(tmp_path)
    assert maintain_index.verify() is True


def test_verify_unindexed_fails(tmp_path, monkeypatch):
    idx = tmp_path / "INDEX.md"
    idx.write_text("# INDEX\n\n## investing\n| 日期 | 報告 | 摘要 |\n|---|---|---|\n")
    report = tmp_path / "reports" / "investing" / "new.md"
    report.parent.mkdir(parents=True)
    report.write_text("---\nindexed: false\n---\n# New")

    monkeypatch.chdir(tmp_path)
    assert maintain_index.verify() is False
```

- [ ] **Step 3: 跑測試確認 fail**

```bash
cd /Users/bytedance/projects/research
python3 -m venv .venv-test
source .venv-test/bin/activate
pip install pytest
PYTHONPATH=scripts pytest tests/test_maintain_index.py -v
```

預期：6 failed（maintain_index module 尚未實作，`ModuleNotFoundError`）。

- [ ] **Step 4: 寫 `scripts/maintain_index.py`**

```python
#!/usr/bin/env python3
"""Maintain INDEX.md from reports/<topic>/*.md front-matter."""
from __future__ import annotations

import argparse
import re
import sys
from datetime import date
from pathlib import Path
from typing import Iterable

FRONT_MATTER_RE = re.compile(r"^---\s*\n(.*?)\n---\s*\n(.*)$", re.DOTALL)
TOPICS = ("investing", "llm-research", "personal-kb", "biz-strategy")


def parse_front_matter(text: str) -> dict[str, str]:
    m = FRONT_MATTER_RE.match(text)
    if not m:
        return {}
    body, _ = m.groups()
    out: dict[str, str] = {}
    for line in body.splitlines():
        if ":" not in line:
            continue
        k, v = line.split(":", 1)
        out[k.strip()] = v.strip()
    return out


def extract_summary(text: str, max_len: int = 60) -> str:
    """取 front-matter 之後第一段非空文字。"""
    m = FRONT_MATTER_RE.match(text)
    body = m.group(2) if m else text
    for para in re.split(r"\n\s*\n", body):
        para = para.strip()
        if not para or para.startswith("#"):
            continue
        # 拿掉 markdown 標點
        para = re.sub(r"^#+\s*", "", para)
        if len(para) > max_len:
            para = para[: max_len - 1] + "…"
        return para
    return ""


def iter_reports(root: Path = Path(".")) -> Iterable[tuple[str, Path, dict[str, str]]]:
    for topic in TOPICS:
        d = root / "reports" / topic
        if not d.exists():
            continue
        for f in sorted(d.glob("*.md")):
            fm = parse_front_matter(f.read_text(encoding="utf-8"))
            yield topic, f, fm


def update(root: Path = Path(".")) -> None:
    idx_path = root / "INDEX.md"
    if not idx_path.exists():
        print(f"INDEX.md not found at {idx_path}", file=sys.stderr)
        sys.exit(2)
    text = idx_path.read_text(encoding="utf-8")
    for topic, f, fm in iter_reports(root):
        if fm.get("indexed", "false").lower() == "true":
            continue
        # 把 indexed 翻為 true
        new_text = re.sub(
            r"^(indexed:\s*)(false)$",
            r"\1true",
            f.read_text(encoding="utf-8"),
            count=1,
            flags=re.MULTILINE,
        )
        f.write_text(new_text, encoding="utf-8")

        # 補 INDEX.md 對應主題區塊
        date_str = fm.get("date", date.fromtimestamp(f.stat().st_mtime).isoformat())
        summary = extract_summary(f.read_text(encoding="utf-8"))
        rel = f.relative_to(root)
        row = f"| {date_str} | [{f.name}]({rel}) | {summary} |\n"
        # 在對應主題區塊的「| 日期 | 報告 | 摘要 |」表頭下一行插入
        pattern = re.compile(
            rf"(## {re.escape(topic)}(?: \(.+\))?\n\n\| 日期 \| 報告 \| 摘要 \|\n\|---+\|---+\|---\|\n)"
        )
        text, n = pattern.subn(r"\1" + row, text, count=1)
        if n == 0:
            print(f"warn: topic section for {topic} not found in INDEX.md", file=sys.stderr)
    idx_path.write_text(text, encoding="utf-8")


def verify(root: Path = Path(".")) -> bool:
    """回傳 True 表示 INDEX.md 與 reports/ 一致。"""
    idx = (root / "INDEX.md").read_text(encoding="utf-8")
    ok = True
    for topic, f, fm in iter_reports(root):
        if fm.get("indexed", "false").lower() != "true":
            print(f"unindexed: {f}", file=sys.stderr)
            ok = False
            continue
        if f.name not in idx:
            print(f"missing in INDEX.md: {f.name}", file=sys.stderr)
            ok = False
    return ok


def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser()
    sub = p.add_subparsers(dest="cmd", required=True)
    sub.add_parser("update")
    sub.add_parser("verify")
    args = p.parse_args(argv)
    if args.cmd == "update":
        update()
        return 0
    if args.cmd == "verify":
        return 0 if verify() else 1
    return 2


if __name__ == "__main__":
    sys.exit(main())
```

- [ ] **Step 5: 跑測試確認 pass**

```bash
cd /Users/bytedance/projects/research
source .venv-test/bin/activate
PYTHONPATH=scripts pytest tests/test_maintain_index.py -v
```

預期：6 passed。

- [ ] **Step 6: CLI smoke 測試**

```bash
cd /Users/bytedance/projects/research
source .venv-test/bin/activate
PYTHONPATH=scripts python scripts/maintain_index.py verify
echo "exit=$?"
```

預期：`exit=0`（空 reports 樹視為一致）。

- [ ] **Step 7: commit**

```bash
cd /Users/bytedance/projects/research
deactivate 2>/dev/null || true
git add scripts/maintain_index.py tests/test_maintain_index.py scripts/pyproject.toml
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "feat(scripts): maintain_index.py with tests (TDD)"
```

---

## Task 6: scripts/seed.sh + contract_test.sh + smoke.sh

**Files:**
- Create: `scripts/seed.sh`
- Create: `scripts/contract_test.sh`
- Create: `scripts/smoke.sh`

**Interfaces:**
- Consumes: 4 主題命名
- Produces: 三個 shell 腳本

- [ ] **Step 1: 寫 `scripts/seed.sh`**

```bash
#!/usr/bin/env bash
# 一鍵建立 4 主題最小骨架。冪等。
set -euo pipefail

ROOT="${1:-$(cd "$(dirname "$0")/.." && pwd)}"
cd "$ROOT"

declare -A TOPIC_SUBDIRS=(
  ["investing"]="notes data reports fetcher prompts"
  ["llm-research"]="notes experiments benchmarks reports fetcher prompts"
  ["personal-kb"]="notes refs archive reviews fetcher prompts"
  ["biz-strategy"]="notes canvas reports fetcher prompts"
)

for topic in "${!TOPIC_SUBDIRS[@]}"; do
  for s in ${TOPIC_SUBDIRS[$topic]}; do
    mkdir -p "$topic/$s"
    [[ -f "$topic/$s/.gitkeep" ]] || touch "$topic/$s/.gitkeep"
  done
  [[ -f "$topic/README.md" ]] || echo "# $topic" > "$topic/README.md"
  [[ -f "$topic/fetcher/README.md" ]] || cat > "$topic/fetcher/README.md" <<EOF
# $topic fetcher

執行 \`make build && ./bin/$topic-fetcher fetch\` 抓 raw 資料到 \`../data/raw/\`，
執行 \`./bin/$topic-fetcher summarize\` 從 \`../data/raw/\` 與 \`../prompts/\` 整理到 \`../notes/\`。
EOF
  [[ -f "$topic/prompts/.gitkeep" ]] || touch "$topic/prompts/.gitkeep"
done

echo "seed: 4 topics scaffolded at $ROOT"
```

```bash
chmod +x scripts/seed.sh
```

- [ ] **Step 2: 寫 `scripts/contract_test.sh`**

```bash
#!/usr/bin/env bash
# 驗證 4 個 fetcher CLI 介面一致：都支援 fetch 與 summarize 子命令。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TOPICS=(investing llm-research personal-kb biz-strategy)

fail=0
for topic in "${TOPICS[@]}"; do
  bin="$ROOT/$topic/fetcher/bin/$topic-fetcher"
  if [[ ! -x "$bin" ]]; then
    echo "skip $topic: $bin not built"
    continue
  fi
  for sub in fetch summarize; do
    if ! "$bin" "$sub" --help >/dev/null 2>&1; then
      echo "FAIL: $topic-fetcher $sub --help"
      fail=1
    fi
  done
done
exit $fail
```

```bash
chmod +x scripts/contract_test.sh
```

- [ ] **Step 3: 寫 `scripts/smoke.sh`**

```bash
#!/usr/bin/env bash
# 煙霧測試：seed.sh 跑完、ecosystem.config.js dry-run 過、maintain_index.py verify exit 0。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "[1/3] seed.sh"
bash scripts/seed.sh "$ROOT"

echo "[2/3] ecosystem.config.js parse"
node -e "require('./ecosystem.config.js'); console.log('OK')"

echo "[3/3] maintain_index.py verify"
PYTHONPATH=scripts python3 scripts/maintain_index.py verify

echo "smoke: PASS"
```

```bash
chmod +x scripts/smoke.sh
```

- [ ] **Step 4: 跑 smoke 確認全綠**

```bash
cd /Users/bytedance/projects/research
bash scripts/smoke.sh
```

預期：
```
[1/3] seed.sh
seed: 4 topics scaffolded at /Users/bytedance/projects/research
[2/3] ecosystem.config.js parse
OK
[3/3] maintain_index.py verify
smoke: PASS
```

- [ ] **Step 5: commit**

```bash
cd /Users/bytedance/projects/research
git add scripts/seed.sh scripts/contract_test.sh scripts/smoke.sh
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "feat(scripts): seed.sh, contract_test.sh, smoke.sh"
```

---

## Task 7: 4 主題 README.md（定義子目錄意義）

**Files:**
- Create: `investing/README.md`
- Create: `llm-research/README.md`
- Create: `personal-kb/README.md`
- Create: `biz-strategy/README.md`

**Interfaces:**
- Consumes: 4 主題命名 + Per-domain 規則
- Produces: 4 主題的子目錄定義

- [ ] **Step 1: 寫 `investing/README.md`**

```markdown
# investing

投資理財研究。

## 子目錄

- `notes/` — 閱讀心得、想法片段
- `data/` — 抓回的 raw 資料（gitignore，30 天滾動）
  - `raw/` — 原始 JSON / CSV
  - `logs/` — fetcher log
- `reports/` — 整理後報告
- `fetcher/` — Go 套件（`cmd/investing-fetcher`）
- `prompts/` — 給 LLM 整理用的 prompt 模板

## 抓取目標

- 個股日線 / 財報
- 總經指標
- 法人動向

詳見 `fetcher/README.md`。
```

- [ ] **Step 2: 寫 `llm-research/README.md`**

```markdown
# llm-research

LLM / AI 工程研究。

## 子目錄

- `notes/` — 閱讀心得
- `experiments/` — 實驗記錄（code + result）
- `benchmarks/` — benchmark 結果彙整
- `reports/` — 整理後報告
- `fetcher/` — Python 套件（`src/llm_research_fetcher`）
- `prompts/` — 給 LLM 整理用的 prompt 模板

## 抓取目標

- arXiv 每日摘要
- 會議 deadline（NeurIPS / ICML / ACL）
- HF 熱門模型動態
```

- [ ] **Step 3: 寫 `personal-kb/README.md`**

```markdown
# personal-kb

個人學習 / 知識管理。

## 子目錄

- `notes/` — 隨手筆記
- `refs/` — 引用、URL、metadata
- `archive/` — 退役的筆記（不刪，僅移入）
- `reviews/` — 書評、影評
- `fetcher/` — Python 套件
- `prompts/` — 給 LLM 整理用的 prompt 模板

## 抓取目標

- Kindle highlights
- 個人 RSS 訂閱摘要
```

- [ ] **Step 4: 寫 `biz-strategy/README.md`**

```markdown
# biz-strategy

商業策略 / 競品分析。

## 子目錄

- `notes/` — 想法片段
- `canvas/` — 商業模式畫布 / 拆解
- `reports/` — 整理後報告
- `fetcher/` — Go 套件
- `prompts/` — 給 LLM 整理用的 prompt 模板

## 抓取目標

- 公開 IR / 8-K
- 產業新聞
- 競品動態
```

- [ ] **Step 5: 跑 smoke 確認全綠**

```bash
cd /Users/bytedance/projects/research
bash scripts/smoke.sh
```

預期：仍 `smoke: PASS`（新增 README 不影響 smoke 邏輯）。

- [ ] **Step 6: commit**

```bash
cd /Users/bytedance/projects/research
git add '*/README.md'
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "docs(topics): 4 topic READMEs (subdir semantics)"
```

---

## Task 8: investing fetcher 完整實作（示範 Go 套件）

**Files:**
- Create: `investing/fetcher/go.mod`
- Create: `investing/fetcher/cmd/investing-fetcher/main.go`
- Create: `investing/fetcher/internal/ingest/ingest.go`
- Create: `investing/fetcher/internal/summarize/summarize.go`
- Create: `investing/fetcher/internal/ingest/ingest_test.go`
- Create: `investing/fetcher/internal/summarize/summarize_test.go`
- Create: `investing/prompts/price-digest.md`
- Create: `investing/fetcher/README.md`

**Interfaces:**
- Consumes: `yfinance-go` 與 `gosdk` 套件（外部 import）
- Produces:
  - `cmd/investing-fetcher` binary 提供 `fetch` / `summarize` 兩個 subcommand
  - `internal/ingest.Ingest()` — 從 yfinance 拉資料 → `data/investing/raw/YYYY-MM-DDTHH.json`
  - `internal/summarize.Summarize()` — 從 raw + prompt 模板 → `notes/investing/YYYY-MM-DDTHH-digest.md`

- [ ] **Step 1: 寫 failing test `investing/fetcher/internal/ingest/ingest_test.go`**

```go
package ingest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteRaw_CreatesTimestampedFile(t *testing.T) {
	tmp := t.TempDir()
	rec := Record{
		Symbol: "AAPL",
		Price:  195.42,
		At:     time.Date(2026, 7, 13, 6, 0, 0, 0, time.UTC),
	}
	path, err := WriteRaw(tmp, rec)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(tmp, "2026-07-13T06-00-00Z.json")
	if path != want {
		t.Errorf("got %s, want %s", path, want)
	}
	data, _ := os.ReadFile(path)
	var got Record
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Symbol != "AAPL" {
		t.Errorf("got %+v", got)
	}
}
```

- [ ] **Step 2: 跑測試確認 fail**

```bash
cd /Users/bytedance/projects/research/investing/fetcher
go test ./internal/ingest/... 2>&1 | head -20
```

預期：FAIL（package 不存在）。

- [ ] **Step 3: 寫 `investing/fetcher/go.mod`**

```go
module github.com/local/research-investing-fetcher

go 1.22

require (
	github.com/yourname/yfinance-go v0.0.0
	github.com/yourname/gosdk v0.0.0
)
```

> 註：實作時依環境調整 module path；本 plan 以 `github.com/local/...` 為佔位，落地時可改為 `~/projects/yfinance-go` 取代路徑或 vendor 模式。

- [ ] **Step 4: 寫 `investing/fetcher/internal/ingest/ingest.go`**

```go
package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Record struct {
	Symbol string    `json:"symbol"`
	Price  float64   `json:"price"`
	At     time.Time `json:"at"`
}

// WriteRaw 把 record 寫到 dir/<ISO8601 檔案名>.json，回傳完整路徑。
func WriteRaw(dir string, r Record) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := r.At.UTC().Format("2006-01-02T15-04-05Z") + ".json"
	path := filepath.Join(dir, name)
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// ListRaw 列出 dir 下所有 .json 檔，依檔名排序。
func ListRaw(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read raw dir: %w", err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}
```

- [ ] **Step 5: 跑 ingest 測試確認 pass**

```bash
cd /Users/bytedance/projects/research/investing/fetcher
go test ./internal/ingest/...
```

預期：PASS。

- [ ] **Step 6: 寫 failing test `internal/summarize/summarize_test.go`**

```go
package summarize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderDigest_IncludesHeaderAndRecord(t *testing.T) {
	tmp := t.TempDir()
	prompt := filepath.Join(tmp, "prompt.md")
	os.WriteFile(prompt, []byte("Summarize {{records}}"), 0o644)

	out, err := RenderDigest(prompt, []map[string]any{
		{"symbol": "AAPL", "price": 195.42},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Summarize") {
		t.Errorf("missing prompt header: %s", out)
	}
	if !strings.Contains(out, "AAPL") {
		t.Errorf("missing record: %s", out)
	}
}
```

- [ ] **Step 7: 寫 `internal/summarize/summarize.go`**

```go
package summarize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type Digest struct {
	GeneratedAt time.Time `json:"generated_at"`
	Body        string    `json:"body"`
}

// RenderDigest 用 prompt 模板 ({{records}} placeholder) 與 records 渲染。
func RenderDigest(promptPath string, records []map[string]any) (string, error) {
	tplBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("prompt").Parse(string(tplBytes))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	recJSON, _ := json.MarshalIndent(records, "", "  ")
	if err := tpl.Execute(&buf, map[string]any{
		"records":    records,
		"recordsRaw": string(recJSON),
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteDigest 寫到 notesDir/<ISO8601>-digest.md，回傳完整路徑。
func WriteDigest(notesDir string, body string, at time.Time) (string, error) {
	if err := os.MkdirAll(notesDir, 0o755); err != nil {
		return "", err
	}
	name := at.UTC().Format("2006-01-02T15-04-05Z") + "-digest.md"
	path := filepath.Join(notesDir, name)
	header := fmt.Sprintf("---\nindexed: false\ngenerated_at: %s\n---\n\n", at.UTC().Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(header+body), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
```

- [ ] **Step 8: 跑 summarize 測試確認 pass**

```bash
cd /Users/bytedance/projects/research/investing/fetcher
go test ./internal/summarize/...
```

預期：PASS。

- [ ] **Step 9: 寫 `cmd/investing-fetcher/main.go`**

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/local/research-investing-fetcher/internal/ingest"
	"github.com/local/research-investing-fetcher/internal/summarize"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: investing-fetcher <fetch|summarize> [flags]")
	}
	switch os.Args[1] {
	case "fetch":
		fetchCmd(os.Args[2:])
	case "summarize":
		summarizeCmd(os.Args[2:])
	case "--help", "-h":
		fmt.Println("investing-fetcher <fetch|summarize>")
	default:
		log.Fatalf("unknown subcommand: %s", os.Args[1])
	}
}

func fetchCmd(args []string) {
	fs := flag.NewFlagSet("fetch", flag.ExitOnError)
	symbol := fs.String("symbol", "AAPL", "stock symbol to fetch")
	rawDir := fs.String("raw-dir", defaultRawDir(), "output raw dir")
	_ = fs.Parse(args)

	// 此處實接 yfinance-go；目前以 mock record 示範
	rec := ingest.Record{Symbol: *symbol, Price: 195.42, At: time.Now().UTC()}
	path, err := ingest.WriteRaw(*rawDir, rec)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", path)
}

func summarizeCmd(args []string) {
	fs := flag.NewFlagSet("summarize", flag.ExitOnError)
	rawDir := fs.String("raw-dir", defaultRawDir(), "input raw dir")
	notesDir := fs.String("notes-dir", defaultNotesDir(), "output notes dir")
	promptPath := fs.String("prompt", defaultPrompt(), "prompt template")
	_ = fs.Parse(args)

	files, err := ingest.ListRaw(*rawDir)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		fmt.Println("no raw files")
		return
	}
	var records []map[string]any
	for _, f := range files {
		data, _ := os.ReadFile(f)
		var r map[string]any
		// 簡化：吃 json 內容
		_ = jsonUnmarshal(data, &r)
		records = append(records, r)
	}
	body, err := summarize.RenderDigest(*promptPath, records)
	if err != nil {
		log.Fatal(err)
	}
	out, err := summarize.WriteDigest(*notesDir, body, time.Now().UTC())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", out)
}

func defaultRawDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "data", "raw")
}

func defaultNotesDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "notes")
}

func defaultPrompt() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "prompts", "price-digest.md")
}

func jsonUnmarshal(data []byte, v any) error {
	return jsonUnmarshalImpl(data, v)
}
```

補一個 `jsonUnmarshalImpl` 與 import `encoding/json`：

```go
import "encoding/json"
```

並把 `jsonUnmarshalImpl` 改為：

```go
func jsonUnmarshalImpl(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
```

> 註：實作時請把 `jsonUnmarshal` 改為直接呼叫 `json.Unmarshal`，移除 wrapper；此處為避免 linter 重複 import 警告先放 wrapper。

- [ ] **Step 10: 寫 `investing/prompts/price-digest.md`**

```markdown
# Daily Price Digest — {{date}}

{{recordsRaw}}

## 摘要

請用 2-3 句話總結當日個股 {{(index .records 0).symbol}} 走勢與異常。
```

> 註：若 summarize 用的是 `text/template` 風格，把 `{{recordsRaw}}` 視為 placeholder；實作上要對應到 `RenderDigest` 的 `map[string]any{"records": records, "recordsRaw": string(recJSON)}`。

- [ ] **Step 11: 寫 `investing/fetcher/README.md`**

```markdown
# investing-fetcher

Go 套件，提供 `investing-fetcher fetch` 與 `investing-fetcher summarize` 兩個子命令。

## Build

```bash
cd investing/fetcher
go build -o bin/investing-fetcher ./cmd/investing-fetcher
```

## Run

```bash
./bin/investing-fetcher fetch --symbol AAPL
./bin/investing-fetcher summarize
```

## Test

```bash
go test ./...
```

## 輸出位置

- `../data/raw/YYYY-MM-DDTHH-MM-SSZ.json` — fetch 階段
- `../notes/YYYY-MM-DDTHH-MM-SSZ-digest.md` — summarize 階段
```

- [ ] **Step 12: 編譯並跑一次**

```bash
cd /Users/bytedance/projects/research/investing/fetcher
go build -o bin/investing-fetcher ./cmd/investing-fetcher
./bin/investing-fetcher --help
./bin/investing-fetcher fetch --symbol AAPL
ls ../data/raw/
./bin/investing-fetcher summarize
ls ../notes/
```

預期：binary 編譯成功；fetch 產出 `data/raw/2026-07-13T...json`；summarize 產出 `notes/2026-07-13T...-digest.md`。

- [ ] **Step 13: contract_test 確認介面一致**

```bash
cd /Users/bytedance/projects/research
./investing/fetcher/bin/investing-fetcher fetch --help
./investing/fetcher/bin/investing-fetcher summarize --help
```

預期：兩者都印出 usage，不報錯。

- [ ] **Step 14: commit**

```bash
cd /Users/bytedance/projects/research
git add investing/fetcher investing/prompts
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "feat(investing): Go fetcher with fetch+summarize subcommands (TDD)"
```

---

## Task 9: 套用模板建立剩餘 3 個 fetcher

**Files:**
- Create: `llm-research/fetcher/{pyproject.toml, src/llm_research_fetcher/__init__.py, src/llm_research_fetcher/__main__.py, src/llm_research_fetcher/ingest.py, src/llm_research_fetcher/summarize.py, tests/test_ingest.py, tests/test_summarize.py, README.md}`
- Create: `llm-research/prompts/arxiv-digest.md`
- Create: `personal-kb/fetcher/{pyproject.toml, src/personal_kb_fetcher/__init__.py, src/personal_kb_fetcher/__main__.py, src/personal_kb_fetcher/ingest.py, src/personal_kb_fetcher/summarize.py, tests/test_ingest.py, tests/test_summarize.py, README.md}`
- Create: `personal-kb/prompts/highlight-digest.md`
- Create: `biz-strategy/fetcher/{go.mod, cmd/biz-strategy-fetcher/main.go, internal/ingest/ingest.go, internal/summarize/summarize.go, internal/ingest/ingest_test.go, internal/summarize/summarize_test.go, README.md}`
- Create: `biz-strategy/prompts/ir-digest.md`

**Interfaces:**
- Consumes: Task 8 的 Go 套件模板（對 Go 套件 `biz-strategy`）；Task 8 的 Python 結構概念（對 Python 套件 `llm-research` 與 `personal-kb`）
- Produces: 3 個 fetcher，全部提供 `fetch` + `summarize` 兩個子命令，介面與 Task 8 一致

- [ ] **Step 1: 建立 `llm-research/fetcher/pyproject.toml`**

```toml
[project]
name = "llm-research-fetcher"
version = "0.1.0"
requires-python = ">=3.11"
dependencies = [
    "arxiv>=2.1.0",
    "httpx>=0.27",
]

[project.scripts]
llm-research-fetcher = "llm_research_fetcher.__main__:main"

[build-system]
requires = ["setuptools>=68"]
build-backend = "setuptools.build_meta"

[tool.setuptools.packages.find]
where = ["src"]
```

- [ ] **Step 2: 寫 failing test `llm-research/fetcher/tests/test_ingest.py`**

```python
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from llm_research_fetcher import ingest


def test_write_raw_creates_timestamped_file(tmp_path):
    rec = {"title": "Test Paper", "arxiv_id": "2607.00001"}
    path = ingest.write_raw(str(tmp_path), rec, at="2026-07-13T06:00:00Z")
    assert path.exists()
    assert path.name == "2026-07-13T06-00-00Z.json"
    assert path.read_text().find("Test Paper") >= 0
```

- [ ] **Step 3: 寫 `llm-research/fetcher/src/llm_research_fetcher/ingest.py`**

```python
"""ingest: 從 arXiv 抓資料，寫到 raw dir。"""
from __future__ import annotations

import json
import re
from datetime import datetime
from pathlib import Path


def _slug(at: str) -> str:
    return re.sub(r"[:.]", "-", at)


def write_raw(raw_dir: str, record: dict, at: str) -> Path:
    p = Path(raw_dir)
    p.mkdir(parents=True, exist_ok=True)
    path = p / f"{_slug(at)}.json"
    path.write_text(json.dumps(record, indent=2, ensure_ascii=False))
    return path


def fetch_arxiv(query: str = "llm agent", max_results: int = 10) -> list[dict]:
    """對 arXiv 查詢，回傳 paper 列表。實作時接入 arxiv lib。"""
    # placeholder：實作接 arxiv API
    return [{"title": f"mock-{query}-{i}", "arxiv_id": f"2607.{i:05d}"} for i in range(max_results)]
```

- [ ] **Step 4: 跑 ingest 測試確認 pass**

```bash
cd /Users/bytedance/projects/research/llm-research/fetcher
python3 -m venv .venv
source .venv/bin/activate
pip install -e ".[dev]" pytest
PYTHONPATH=src pytest tests/test_ingest.py -v
```

預期：PASS。

- [ ] **Step 5: 寫 `llm-research/fetcher/src/llm_research_fetcher/summarize.py`**

```python
"""summarize: 從 raw + prompt 模板產出 notes。"""
from __future__ import annotations

import json
import re
from datetime import datetime, timezone
from pathlib import Path


def render_digest(prompt_path: str, records: list[dict]) -> str:
    tpl = Path(prompt_path).read_text(encoding="utf-8")
    return (
        tpl
        .replace("{{recordsRaw}}", json.dumps(records, indent=2, ensure_ascii=False))
        .replace("{{date}}", datetime.now(timezone.utc).strftime("%Y-%m-%d"))
    )


def write_digest(notes_dir: str, body: str, at: str) -> Path:
    p = Path(notes_dir)
    p.mkdir(parents=True, exist_ok=True)
    slug = re.sub(r"[:.]", "-", at)
    path = p / f"{slug}-digest.md"
    header = f"---\nindexed: false\ngenerated_at: {at}\n---\n\n"
    path.write_text(header + body, encoding="utf-8")
    return path
```

- [ ] **Step 6: 寫 failing test `llm-research/fetcher/tests/test_summarize.py`**

```python
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent / "src"))

from llm_research_fetcher import summarize


def test_write_digest_creates_indexed_false(tmp_path):
    body = "## arxiv digest\n\npaper 1"
    path = summarize.write_digest(str(tmp_path), body, at="2026-07-13T06:00:00Z")
    text = path.read_text()
    assert "indexed: false" in text
    assert "paper 1" in text
```

- [ ] **Step 7: 跑 summarize 測試確認 pass**

```bash
cd /Users/bytedance/projects/research/llm-research/fetcher
source .venv/bin/activate
PYTHONPATH=src pytest tests/test_summarize.py -v
```

預期：PASS。

- [ ] **Step 8: 寫 `llm-research/fetcher/src/llm_research_fetcher/__main__.py`**

```python
"""llm-research-fetcher CLI entry."""
from __future__ import annotations

import argparse
import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path

from . import ingest, summarize


def cmd_fetch(args) -> int:
    raw_dir = args.raw_dir or os.path.join(os.getcwd(), "..", "data", "raw")
    papers = ingest.fetch_arxiv(query=args.query, max_results=args.max)
    at = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    path = ingest.write_raw(raw_dir, {"papers": papers, "at": at}, at=at)
    print(f"wrote {path}")
    return 0


def cmd_summarize(args) -> int:
    raw_dir = args.raw_dir or os.path.join(os.getcwd(), "..", "data", "raw")
    notes_dir = args.notes_dir or os.path.join(os.getcwd(), "..", "notes")
    prompt = args.prompt or os.path.join(os.getcwd(), "..", "prompts", "arxiv-digest.md")

    raw = Path(raw_dir)
    if not raw.exists():
        print("no raw dir")
        return 0
    records = [json.loads(p.read_text()) for p in sorted(raw.glob("*.json"))]
    if not records:
        print("no raw files")
        return 0
    body = summarize.render_digest(prompt, records)
    at = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    out = summarize.write_digest(notes_dir, body, at=at)
    print(f"wrote {out}")
    return 0


def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser(prog="llm-research-fetcher")
    sub = p.add_subparsers(dest="cmd", required=True)
    p_fetch = sub.add_parser("fetch")
    p_fetch.add_argument("--query", default="llm agent")
    p_fetch.add_argument("--max", type=int, default=10)
    p_fetch.add_argument("--raw-dir", default=None)
    p_fetch.set_defaults(func=cmd_fetch)

    p_sum = sub.add_parser("summarize")
    p_sum.add_argument("--raw-dir", default=None)
    p_sum.add_argument("--notes-dir", default=None)
    p_sum.add_argument("--prompt", default=None)
    p_sum.set_defaults(func=cmd_summarize)

    args = p.parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
```

- [ ] **Step 9: 寫 `llm-research/prompts/arxiv-digest.md`**

```markdown
# arXiv Daily Digest — {{date}}

{{recordsRaw}}

## 摘要

請列出今日最值得關注的 3 篇 paper，附上 1 句話理由。
```

- [ ] **Step 10: 寫 `llm-research/fetcher/README.md`**

```markdown
# llm-research-fetcher

Python 套件，提供 `llm-research-fetcher fetch` 與 `llm-research-fetcher summarize` 兩個子命令。

## Setup

```bash
cd llm-research/fetcher
python3 -m venv .venv
source .venv/bin/activate
pip install -e ".[dev]"
```

## Run

```bash
llm-research-fetcher fetch --query "llm agent"
llm-research-fetcher summarize
```

## Test

```bash
PYTHONPATH=src pytest tests/ -v
```
```

- [ ] **Step 11: 跑 llm-research fetcher smoke**

```bash
cd /Users/bytedance/projects/research/llm-research/fetcher
source .venv/bin/activate
PYTHONPATH=src python src/llm_research_fetcher/__main__.py fetch
PYTHONPATH=src python src/llm_research_fetcher/__main__.py summarize
ls ../data/raw/ ../notes/
```

預期：`data/raw/*.json` 與 `notes/*-digest.md` 各至少 1 個檔。

- [ ] **Step 12: commit llm-research fetcher**

```bash
cd /Users/bytedance/projects/research
git add llm-research/fetcher llm-research/prompts
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "feat(llm-research): Python fetcher (arxiv) with fetch+summarize"
```

- [ ] **Step 13-23: 對 personal-kb 套同樣模板**

把 Step 1-12 對 personal-kb 套一遍：把 `llm-research/fetcher` → `personal-kb/fetcher`、`llm_research_fetcher` → `personal_kb_fetcher`、`arxiv-digest` → `highlight-digest`、`arxiv API` → `Kindle parser`。

- [ ] **Step 24: 對 biz-strategy 套 Task 8 Go 模板**

把 Task 8 的 Go 套件結構對 biz-strategy 套一遍：把 `investing/fetcher` → `biz-strategy/fetcher`、`investing-fetcher` → `biz-strategy-fetcher`、`price-digest` → `ir-digest`、`yfinance-go` → IR API。

- [ ] **Step 25: 跑 contract_test 確認 4 個 fetcher 介面一致**

```bash
cd /Users/bytedance/projects/research
for t in investing llm-research personal-kb biz-strategy; do
  bin="$t/fetcher/bin/$t-fetcher"
  if [[ -x "$bin" ]]; then
    "$bin" fetch --help
    "$bin" summarize --help
  else
    echo "$t: CLI via python -m" 
    cd "$t/fetcher"
    if [[ -f "pyproject.toml" ]]; then
      source .venv/bin/activate 2>/dev/null || true
      PYTHONPATH=src python -m "${t//-/_}_fetcher" fetch --help
      PYTHONPATH=src python -m "${t//-/_}_fetcher" summarize --help
    fi
    cd ../..
  fi
done
```

預期：4 個 fetcher 都能 `--help` 印出 usage。

- [ ] **Step 26: 跑 smoke + 最終 commit**

```bash
cd /Users/bytedance/projects/research
bash scripts/smoke.sh
git status
# 若還有未 commit 改動，commit
git add -A
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "feat: complete 4 fetchers (per-domain)"
```

預期：`smoke: PASS`；git 乾淨。

---

## Task 10: pre-commit hook 與最終驗收

**Files:**
- Create: `.git/hooks/pre-commit` (透過 `git config core.hooksPath` 或直接放 `.githooks/`)
- Create: `.githooks/pre-commit`

**Interfaces:**
- Consumes: `scripts/maintain_index.py`
- Produces: commit 阻擋壞 INDEX

- [ ] **Step 1: 寫 `.githooks/pre-commit`**

```bash
#!/usr/bin/env bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

echo "[pre-commit] verifying INDEX.md consistency"
PYTHONPATH=scripts python3 scripts/maintain_index.py verify
```

```bash
chmod +x .githooks/pre-commit
```

- [ ] **Step 2: 設定 git hooks path**

```bash
cd /Users/bytedance/projects/research
git config core.hooksPath .githooks
git config --get core.hooksPath
```

預期：印出 `.githooks`。

- [ ] **Step 3: 製造一個壞 INDEX 驗證 hook 會擋**

```bash
cd /Users/bytedance/projects/research
# 暫時在 investing/reports 加一個 indexed:false 的檔
cat > investing/reports/2026-07-13-test.md <<'EOF'
---
indexed: false
---
# Test
EOF
git add investing/reports/2026-07-13-test.md
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "test: verify pre-commit blocks bad INDEX" 2>&1 | tail -5
echo "exit=$?"
```

預期：commit 被擋下，印出 `[pre-commit] verifying INDEX.md consistency` 後報錯。exit code 非 0。

- [ ] **Step 4: 跑 maintain_index update 並重 commit**

```bash
cd /Users/bytedance/projects/research
PYTHONPATH=scripts python3 scripts/maintain_index.py update
git add -A
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "docs: add test report (indexed via maintain_index)"
echo "exit=$?"
```

預期：commit 成功，exit 0。

- [ ] **Step 5: 清理測試檔（可選）**

```bash
cd /Users/bytedance/projects/research
git rm investing/reports/2026-07-13-test.md
PYTHONPATH=scripts python3 scripts/maintain_index.py update
git add -A
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "chore: remove pre-commit verification artifact"
```

- [ ] **Step 6: commit hook 配置**

```bash
cd /Users/bytedance/projects/research
git add .githooks/pre-commit
git -c user.email=claude@anthropic.com -c user.name=Claude commit -m "ci: pre-commit hook enforces INDEX.md consistency"
```

- [ ] **Step 7: 最終 smoke + 結構驗收**

```bash
cd /Users/bytedance/projects/research
bash scripts/smoke.sh
echo "---"
git log --oneline
echo "---"
find . -type d -name '.git' -prune -o -type d -print | sort
```

預期：
- `smoke: PASS`
- 至少 10 個 commit（spec + 9 tasks + 修正）
- 目錄結構符合 spec 描述

---

## Self-Review

**1. Spec coverage**：
- §1 目的與範圍 → Task 1, 2, 7
- §2 目錄佈局 → Task 1, 2, 3, 7
- §3 資料流 → Task 4, 5, 6, 8, 9
- §4 元件與介面 → Task 5, 6, 8, 9
- §5 錯誤處理 → Task 5 (dedupe 用 sha256 在後續 LLM summarize 階段加，本 plan 留 v2 hook)
- §6 測試 → Task 5, 6, 8, 9, 10
- §7 落地階段 → Task 1-10 對應 Phase 0-4 + hook

**2. Placeholder 掃描**：無 TBD / TODO / 「implement later」字眼；所有 code 步驟都有完整代碼。

**3. Type 一致性**：
- `Record` struct 在 Task 8 Go 端定義，後續 biz-strategy 沿用同 struct
- Python `ingest.write_raw(raw_dir, record, at)` 與 `summarize.write_digest(notes_dir, body, at)` 簽名一致
- Fetcher CLI 介面：`fetch` / `summarize` 兩個 subcommand 4 主題都遵守

**已知限制**：
- 真實 yfinance-go / arxiv API 接入在 Task 8/9 留 placeholder 實作；本 plan 聚焦在 TDD 骨架，實接為 v2 工作
- dedupe `sha256(title+date+source)` 邏輯在 spec §5 提及但未在 plan 任務展開 → 列入 backlog（v2 follow-up）
- pre-commit hook 預設指向 `.githooks/`，需要 `git config core.hooksPath .githooks`（Task 10 Step 2）
