# Agent Memory System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 建置一個跨筆電/伺服器的分層 agent 記憶系統，核心交付物是每日凌晨執行的蒸餾器 (cron distiller)，把暫存記憶蒸餾成 long-term (agentmemory) 與驗證 100% 真實事實 (mempalace)。

**Architecture:** 蒸餾器核心邏輯只依賴內部 Protocol（SourceReader / LLMClient / LongTermStore / TruthStore / Pruner），以 in-memory fake 做 TDD；真實外部工具 (claude-mem、gbrain、agentmemory、mempalace) 由設定驅動的薄 adapter 接上，介面在 Phase 0 探查後寫入設定。同步用 syncthing 單向 folder，排程用 cron。

**Tech Stack:** Python 3.11+、uv（套件/虛擬環境）、pytest（測試）、sqlite3（狀態庫，標準庫）、ollama HTTP API（抽取用 LLM）、syncthing（同步）、cron（排程）。

---

## Assumptions & Decisions（可改，實作前請確認）

這些是 spec 未鎖定、我為了讓計畫可落地所做的決策。若不同意請於開工前提出：

- 蒸餾器語言用 `Python 3.11+`，以 `uv` 管理，測試用 `pytest`。
- 新套件位置：`pkg/memory/distiller/`（src layout：`pkg/memory/distiller/src/distiller/`）。
- 抽取 LLM 預設走本地 `ollama`（model 由設定指定，與 gbrain 既有 `bge-m3` 同生態，local-first）。
- 執行期資料：狀態庫 `~/.distiller/state.db`、日誌 `~/.distiller/logs/`。
- 排程：伺服器 cron 每日 `03:00` 跑 `distiller run`（run 內含蒸餾＋30 天保留清理）。
- 同步：syncthing 兩個單向 folder（claude-mem 筆電→伺服器 send-only；agentmemory 伺服器→筆電 send-only）。
- 真實外部工具 (claude-mem SQLite schema、agentmemory/mempalace 寫入指令) 的確切介面未知，於 Phase 0 探查並寫入 `config.toml`，adapter 以設定驅動，不在程式碼硬寫。

---

## File Structure

```tree
pkg/memory/distiller/
├── pyproject.toml                      # uv 專案、entry point: distiller
├── INTERFACES.md                       # Phase 0 探查輸出（各工具確切 CLI/schema）
├── config.sample.toml                  # 設定範本（路徑、指令模板、ollama model）
├── src/distiller/
│   ├── __init__.py
│   ├── models.py                       # Observation / Candidate / Memory / Fact / RunReport
│   ├── interfaces.py                   # Protocol：SourceReader / LLMClient / LongTermStore / TruthStore / Pruner
│   ├── fingerprint.py                  # 正規化 + sha256 指紋（去重/冪等）
│   ├── verify.py                       # 升級 mempalace 的驗證政策
│   ├── state.py                        # StateStore：cursor / seen(corroboration) / distilled
│   ├── pipeline.py                     # Distiller 編排器
│   ├── retention.py                    # 30 天保留清理
│   ├── llm_ollama.py                   # OllamaLLM 抽取 client
│   ├── config.py                       # 設定載入 + 組裝真實 adapter
│   ├── cli.py                          # entry point：run / retain
│   └── adapters/
│       ├── __init__.py
│       ├── gbrain_working.py           # SourceReader：讀 gbrain/working *.md（檔案系統）
│       ├── claude_mem.py               # SourceReader：讀 claude-mem SQLite（設定驅動）
│       └── command_writer.py           # CommandWriter + AgentMemoryStore + MempalaceStore（設定驅動 subprocess）
└── tests/
    ├── conftest.py                     # fakes：FakeSourceReader / FakeLLM / FakeLongTerm / FakeTruth / FakePruner
    ├── test_fingerprint.py
    ├── test_verify.py
    ├── test_state.py
    ├── test_pipeline.py
    ├── test_retention.py
    ├── test_gbrain_working.py
    ├── test_claude_mem.py
    ├── test_command_writer.py
    └── test_llm_ollama.py
```

伺服器/筆電佈署相關變更：

```bash
run.sh                                  # 新增：安裝 distiller、建 ~/.distiller、cron 安裝
pkg/memory/distiller/scripts/install_cron.sh   # 安裝每日 03:00 cron
pkg/memory/distiller/scripts/syncthing.md      # syncthing 兩 folder 設定程序
pkg/memory/hermes-continuity/                  # hermes 連續性：gbrain/working 讀寫 helper + 接線說明
```

---

## Phase 0：探查外部工具介面 (Discovery)

### Task 0: 安裝並記錄各工具確切介面

真實 adapter 依賴這裡的結果，先做。

**Files:**

- Create: `pkg/memory/distiller/INTERFACES.md`

- [x] **Step 1: 在伺服器安裝/確認各工具可用**

Run:

```bash
claude-mem --version || npx claude-mem --version
gbrain doctor --json | jq '.checks[] | {name, status}'
hermes doctor
```

Expected: 三者各自回報版本/健康狀態（非 command not found）。

- [x] **Step 2: 探查 claude-mem 的 SQLite schema**

Run:

```bash
find "$HOME/.claude-mem" -name "*.db" -o -name "*.sqlite*"
# 對找到的 DB：
sqlite3 "$HOME/.claude-mem/<db>" ".tables"
sqlite3 "$HOME/.claude-mem/<db>" ".schema <observations_table>"
```

記錄到 INTERFACES.md：DB 路徑、觀察紀錄表名、id 欄、時間戳欄（epoch?）、文字內容欄。

- [x] **Step 3: 探查 agentmemory / mempalace 的寫入介面**

Run:

```bash
# agentmemory（rohitg00）：確認是 Python 套件還是 CLI
python -c "import agentmemory; print(agentmemory.__file__)" 2>/dev/null || pipx list | grep -i agentmemory
agentmemory --help 2>/dev/null || true
# mempalace：確認 MCP/CLI 寫入方式
mempalace --help 2>/dev/null || true
```

記錄到 INTERFACES.md：各自「新增一筆記憶」的確切指令或 Python API；若為 CLI，記下可接受 stdin JSON 的指令模板。

- [x] **Step 4: 探查 gbrain/working 與 hermes 接線點**

Run:

```bash
ls -la "$HOME/.gbrain"
ls -la "$HOME/.hermes" && sed -n '1,40p' "$HOME/.hermes/config.yaml"
```

記錄到 INTERFACES.md：gbrain/working 目錄實體路徑、hermes 可掛載的擴充點（AGENTS.md / hooks / pre-message script）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/INTERFACES.md
git commit -m "docs(distiller): record external tool interfaces from discovery"
```

---

## Phase 1：專案骨架 (Scaffold)

### Task 1: 建立 uv 專案與測試骨架

**Files:**

- Create: `pkg/memory/distiller/pyproject.toml`
- Create: `pkg/memory/distiller/src/distiller/__init__.py`
- Create: `pkg/memory/distiller/tests/conftest.py`

- [ ] **Step 1: 建立 pyproject.toml**

```toml
[project]
name = "distiller"
version = "0.1.0"
description = "Agent memory distiller: temp memory -> long-term + verified facts"
requires-python = ">=3.11"
dependencies = []

[project.scripts]
distiller = "distiller.cli:main"

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[tool.hatch.build.targets.wheel]
packages = ["src/distiller"]

[dependency-groups]
dev = ["pytest>=8"]
```

- [ ] **Step 2: 建立套件 init**

`src/distiller/__init__.py`:

```python
"""Agent memory distiller."""
__version__ = "0.1.0"
```

- [ ] **Step 3: 建立測試 fakes（conftest）**

`tests/conftest.py`:

```python
import pytest
from distiller.models import Observation, Candidate, Memory, Fact


class FakeReader:
    def __init__(self, name, observations):
        self.name = name
        self._obs = observations
    def read_since(self, last_ts):
        return [o for o in self._obs if o.timestamp > last_ts]


class FakeLLM:
    def __init__(self, candidates):
        self._candidates = candidates
        self.calls = []
    def extract(self, observations):
        self.calls.append(list(observations))
        return list(self._candidates)


class FakeLongTerm:
    def __init__(self):
        self.upserts = []
    def upsert(self, memory: Memory):
        self.upserts.append(memory)


class FakeTruth:
    def __init__(self):
        self.upserts = []
    def upsert(self, fact: Fact):
        self.upserts.append(fact)


class FakePruner:
    def __init__(self):
        self.pruned = []
    def prune(self, source, source_id):
        self.pruned.append((source, source_id))


@pytest.fixture
def make_obs():
    def _make(source="claude-mem", source_id="1", timestamp=100, text="hello"):
        return Observation(source=source, source_id=source_id, timestamp=timestamp, text=text)
    return _make
```

- [ ] **Step 4: 安裝依賴並確認 pytest 可跑（暫無測試）**

Run:

```bash
cd pkg/memory/distiller && uv sync && uv run pytest -q
```

Expected: pytest 啟動，回報 `no tests ran`（exit 5 可接受）。conftest import 會因 models 尚未建立而失敗 — 這是預期，下一個 Task 建立 models 後即解。

- [ ] **Step 5: Commit**

```bash
git add pkg/memory/distiller/pyproject.toml pkg/memory/distiller/src pkg/memory/distiller/tests
git commit -m "chore(distiller): scaffold uv project and test fakes"
```

---

## Phase 2：領域模型 (Domain Models)

### Task 2: 定義 models 與 interfaces

**Files:**

- Create: `pkg/memory/distiller/src/distiller/models.py`
- Create: `pkg/memory/distiller/src/distiller/interfaces.py`

- [ ] **Step 1: 寫 models**

`src/distiller/models.py`:

```python
from dataclasses import dataclass, field


@dataclass(frozen=True)
class Observation:
    """SourceReader 讀出的一筆原始觀察。"""
    source: str
    source_id: str
    timestamp: int            # epoch seconds
    text: str
    metadata: dict = field(default_factory=dict)


@dataclass(frozen=True)
class Candidate:
    """LLM 抽取出的候選記憶。"""
    text: str                                   # verbatim
    entities: tuple[str, ...]
    kind: str                                   # "fact" | "experience" | "preference" | "inference"
    source_refs: tuple[tuple[str, str], ...]    # ((source, source_id), ...)
    first_person: bool = False
    confirmed_by_human: bool = False


@dataclass(frozen=True)
class Memory:
    """寫入 agentmemory 的 long-term 記憶。"""
    fingerprint: str
    text: str
    entities: tuple[str, ...]
    kind: str
    created_at: int


@dataclass(frozen=True)
class Fact:
    """寫入 mempalace 的驗證 100% 真實事實。"""
    fingerprint: str
    text: str
    entities: tuple[str, ...]
    evidence: tuple[tuple[str, str], ...]
    created_at: int


@dataclass
class RunReport:
    sources_read: int = 0
    observations: int = 0
    candidates: int = 0
    long_term_written: int = 0
    facts_written: int = 0
```

- [ ] **Step 2: 寫 interfaces（Protocol）**

`src/distiller/interfaces.py`:

```python
from typing import Iterable, Protocol
from .models import Observation, Candidate, Memory, Fact


class SourceReader(Protocol):
    name: str
    def read_since(self, last_ts: int) -> Iterable[Observation]: ...


class LLMClient(Protocol):
    def extract(self, observations: list[Observation]) -> list[Candidate]: ...


class LongTermStore(Protocol):
    def upsert(self, memory: Memory) -> None: ...


class TruthStore(Protocol):
    def upsert(self, fact: Fact) -> None: ...


class Pruner(Protocol):
    def prune(self, source: str, source_id: str) -> None: ...
```

- [ ] **Step 3: 確認 conftest 可 import（模型就緒）**

Run: `cd pkg/memory/distiller && uv run pytest -q`
Expected: collection 不再 import 失敗；回報 `no tests ran`。

- [ ] **Step 4: Commit**

```bash
git add pkg/memory/distiller/src/distiller/models.py pkg/memory/distiller/src/distiller/interfaces.py
git commit -m "feat(distiller): add domain models and store protocols"
```

---

## Phase 3：指紋去重 (Fingerprint)

### Task 3: fingerprint 模組

**Files:**

- Create: `pkg/memory/distiller/src/distiller/fingerprint.py`
- Test: `pkg/memory/distiller/tests/test_fingerprint.py`

- [ ] **Step 1: 寫失敗測試**

`tests/test_fingerprint.py`:

```python
from distiller.fingerprint import normalize, fingerprint


def test_normalize_collapses_whitespace_and_case():
    assert normalize("  Hello   World  ") == "hello world"


def test_fingerprint_stable_regardless_of_entity_order():
    a = fingerprint("Same text", ("alice", "bob"))
    b = fingerprint("same   TEXT", ("bob", "alice"))
    assert a == b


def test_fingerprint_differs_on_different_text():
    assert fingerprint("text one", ()) != fingerprint("text two", ())
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_fingerprint.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.fingerprint`）。

- [ ] **Step 3: 實作 fingerprint**

`src/distiller/fingerprint.py`:

```python
import hashlib
import re


def normalize(text: str) -> str:
    return re.sub(r"\s+", " ", text.strip().lower())


def fingerprint(text: str, entities: tuple[str, ...]) -> str:
    h = hashlib.sha256()
    h.update(normalize(text).encode())
    h.update(b"|")
    h.update("|".join(sorted(entities)).encode())
    return h.hexdigest()
```

- [ ] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_fingerprint.py -q`
Expected: PASS（3 passed）。

- [ ] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/fingerprint.py pkg/memory/distiller/tests/test_fingerprint.py
git commit -m "feat(distiller): add content fingerprint for dedup/idempotency"
```

---

## Phase 4：驗證政策 (Verification Policy)

### Task 4: verify 模組

實作 spec「驗證政策」：升級進 mempalace 需滿足任一條件。

**Files:**

- Create: `pkg/memory/distiller/src/distiller/verify.py`
- Test: `pkg/memory/distiller/tests/test_verify.py`

- [ ] **Step 1: 寫失敗測試**

`tests/test_verify.py`:

```python
from distiller.models import Candidate
from distiller.verify import qualifies_for_truth


def _cand(**kw):
    base = dict(text="t", entities=(), kind="fact", source_refs=(("s", "1"),),
                first_person=False, confirmed_by_human=False)
    base.update(kw)
    return Candidate(**base)


def test_inference_never_qualifies_even_with_corroboration():
    assert qualifies_for_truth(_cand(kind="inference"), corroboration=5) is False


def test_human_confirmed_qualifies():
    assert qualifies_for_truth(_cand(confirmed_by_human=True), corroboration=1) is True


def test_first_person_life_fact_qualifies():
    assert qualifies_for_truth(_cand(first_person=True, kind="experience"), corroboration=1) is True


def test_two_source_corroboration_qualifies():
    assert qualifies_for_truth(_cand(), corroboration=2) is True


def test_single_unconfirmed_does_not_qualify():
    assert qualifies_for_truth(_cand(), corroboration=1) is False
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_verify.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.verify`）。

- [ ] **Step 3: 實作 verify**

`src/distiller/verify.py`:

```python
from .models import Candidate


def qualifies_for_truth(candidate: Candidate, corroboration: int) -> bool:
    """是否可升級寫入 mempalace（100% 真實）。"""
    if candidate.kind == "inference":
        return False
    if candidate.confirmed_by_human:
        return True
    if candidate.first_person and candidate.kind in ("fact", "experience"):
        return True
    if corroboration >= 2:
        return True
    return False
```

- [ ] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_verify.py -q`
Expected: PASS（5 passed）。

- [ ] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/verify.py pkg/memory/distiller/tests/test_verify.py
git commit -m "feat(distiller): add mempalace verification policy"
```

---

## Phase 5：狀態庫 (State Store)

### Task 5: StateStore（cursor / corroboration / distilled）

**Files:**

- Create: `pkg/memory/distiller/src/distiller/state.py`
- Test: `pkg/memory/distiller/tests/test_state.py`

- [x] **Step 1: 寫失敗測試**

`tests/test_state.py`:

```python
from distiller.state import StateStore


def test_cursor_roundtrip(tmp_path):
    s = StateStore(tmp_path / "state.db")
    assert s.get_cursor("claude-mem") == 0
    s.set_cursor("claude-mem", 1234)
    assert s.get_cursor("claude-mem") == 1234
    s.close()


def test_record_seen_counts_distinct_sources(tmp_path):
    s = StateStore(tmp_path / "state.db")
    assert s.record_seen("fp1", "claude-mem") == 1
    assert s.record_seen("fp1", "claude-mem") == 1   # same source, no increase
    assert s.record_seen("fp1", "gbrain-working") == 2
    s.close()


def test_distilled_marking_and_due_for_prune(tmp_path):
    s = StateStore(tmp_path / "state.db")
    s.mark_distilled("gbrain-working", "a.md", at=1000)
    s.mark_distilled("gbrain-working", "b.md", at=5000)
    assert s.already_distilled("gbrain-working", "a.md") is True
    due = s.due_for_prune(before_ts=3000)
    assert due == [("gbrain-working", "a.md")]
    s.drop_distilled("gbrain-working", "a.md")
    assert s.already_distilled("gbrain-working", "a.md") is False
    s.close()
```

- [x] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_state.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.state`）。

- [x] **Step 3: 實作 StateStore**

`src/distiller/state.py`:

```python
import sqlite3
from pathlib import Path


class StateStore:
    def __init__(self, path):
        self.path = str(path)
        Path(self.path).parent.mkdir(parents=True, exist_ok=True)
        self.con = sqlite3.connect(self.path)
        self._init_schema()

    def _init_schema(self):
        self.con.executescript(
            """
            CREATE TABLE IF NOT EXISTS cursor (
                source TEXT PRIMARY KEY,
                last_ts INTEGER NOT NULL
            );
            CREATE TABLE IF NOT EXISTS seen (
                fingerprint TEXT NOT NULL,
                source TEXT NOT NULL,
                first_seen INTEGER NOT NULL,
                PRIMARY KEY (fingerprint, source)
            );
            CREATE TABLE IF NOT EXISTS distilled (
                source TEXT NOT NULL,
                source_id TEXT NOT NULL,
                distilled_at INTEGER NOT NULL,
                PRIMARY KEY (source, source_id)
            );
            """
        )
        self.con.commit()

    def get_cursor(self, source: str) -> int:
        row = self.con.execute(
            "SELECT last_ts FROM cursor WHERE source = ?", (source,)
        ).fetchone()
        return int(row[0]) if row else 0

    def set_cursor(self, source: str, ts: int) -> None:
        self.con.execute(
            "INSERT INTO cursor(source, last_ts) VALUES(?, ?) "
            "ON CONFLICT(source) DO UPDATE SET last_ts = excluded.last_ts",
            (source, int(ts)),
        )
        self.con.commit()

    def record_seen(self, fingerprint: str, source: str) -> int:
        """記錄此指紋曾由 source 產生，回傳目前佐證來源數（distinct sources）。"""
        self.con.execute(
            "INSERT OR IGNORE INTO seen(fingerprint, source, first_seen) "
            "VALUES(?, ?, strftime('%s','now'))",
            (fingerprint, source),
        )
        self.con.commit()
        return int(
            self.con.execute(
                "SELECT COUNT(*) FROM seen WHERE fingerprint = ?", (fingerprint,)
            ).fetchone()[0]
        )

    def already_distilled(self, source: str, source_id: str) -> bool:
        return (
            self.con.execute(
                "SELECT 1 FROM distilled WHERE source = ? AND source_id = ?",
                (source, source_id),
            ).fetchone()
            is not None
        )

    def mark_distilled(self, source: str, source_id: str, at: int) -> None:
        self.con.execute(
            "INSERT OR REPLACE INTO distilled(source, source_id, distilled_at) "
            "VALUES(?, ?, ?)",
            (source, source_id, int(at)),
        )
        self.con.commit()

    def due_for_prune(self, before_ts: int) -> list[tuple[str, str]]:
        return [
            (r[0], r[1])
            for r in self.con.execute(
                "SELECT source, source_id FROM distilled WHERE distilled_at < ? "
                "ORDER BY source, source_id",
                (int(before_ts),),
            ).fetchall()
        ]

    def drop_distilled(self, source: str, source_id: str) -> None:
        self.con.execute(
            "DELETE FROM distilled WHERE source = ? AND source_id = ?",
            (source, source_id),
        )
        self.con.commit()

    def close(self) -> None:
        self.con.close()
```

- [x] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_state.py -q`
Expected: PASS（3 passed）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/state.py pkg/memory/distiller/tests/test_state.py
git commit -m "feat(distiller): add sqlite state store (cursor/corroboration/distilled)"
```

---

## Phase 6：編排管線 (Pipeline)

### Task 6: Distiller 編排器

**Files:**

- Create: `pkg/memory/distiller/src/distiller/pipeline.py`
- Test: `pkg/memory/distiller/tests/test_pipeline.py`

- [ ] **Step 1: 寫失敗測試**

`tests/test_pipeline.py`:

```python
from distiller.models import Candidate
from distiller.pipeline import Distiller
from distiller.state import StateStore
from tests.conftest import FakeReader, FakeLLM, FakeLongTerm, FakeTruth


def _cand(text, **kw):
    base = dict(text=text, entities=("alice",), kind="fact",
                source_refs=(("claude-mem", "1"),),
                first_person=False, confirmed_by_human=False)
    base.update(kw)
    return Candidate(**base)


def _clock(value=1_000_000):
    return lambda: value


def test_run_writes_long_term_for_every_candidate(tmp_path, make_obs):
    reader = FakeReader("claude-mem", [make_obs(timestamp=10)])
    llm = FakeLLM([_cand("alice likes tea")])
    lt, truth = FakeLongTerm(), FakeTruth()
    state = StateStore(tmp_path / "s.db")
    report = Distiller([reader], llm, lt, truth, state, now=_clock()).run()
    assert len(lt.upserts) == 1
    assert report.long_term_written == 1
    state.close()


def test_run_promotes_only_verified_to_truth(tmp_path, make_obs):
    reader = FakeReader("claude-mem", [make_obs(timestamp=10)])
    llm = FakeLLM([
        _cand("unconfirmed claim"),                       # 不升級
        _cand("alice was born in Taipei", confirmed_by_human=True),  # 升級
    ])
    lt, truth = FakeLongTerm(), FakeTruth()
    state = StateStore(tmp_path / "s.db")
    Distiller([reader], llm, lt, truth, state, now=_clock()).run()
    assert len(lt.upserts) == 2
    assert len(truth.upserts) == 1
    assert truth.upserts[0].text == "alice was born in Taipei"
    state.close()


def test_run_advances_cursor_and_marks_distilled(tmp_path, make_obs):
    reader = FakeReader("claude-mem", [make_obs(source_id="42", timestamp=77)])
    llm = FakeLLM([_cand("x")])
    state = StateStore(tmp_path / "s.db")
    Distiller([reader], llm, FakeLongTerm(), FakeTruth(), state, now=_clock()).run()
    assert state.get_cursor("claude-mem") == 77
    assert state.already_distilled("claude-mem", "42") is True
    state.close()


def test_run_is_incremental(tmp_path, make_obs):
    state = StateStore(tmp_path / "s.db")
    state.set_cursor("claude-mem", 100)
    reader = FakeReader("claude-mem", [make_obs(timestamp=50), make_obs(timestamp=150, source_id="2")])
    llm = FakeLLM([_cand("only new")])
    Distiller([reader], llm, FakeLongTerm(), FakeTruth(), state, now=_clock()).run()
    # 只有 timestamp>100 的觀察被送入 LLM
    assert len(llm.calls) == 1
    assert [o.timestamp for o in llm.calls[0]] == [150]
    state.close()
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_pipeline.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.pipeline`）。

- [ ] **Step 3: 實作 Distiller**

`src/distiller/pipeline.py`:

```python
import time

from .fingerprint import fingerprint
from .models import Fact, Memory, RunReport
from .verify import qualifies_for_truth


class Distiller:
    def __init__(self, readers, llm, long_term, truth, state, now=time.time):
        self.readers = readers
        self.llm = llm
        self.long_term = long_term
        self.truth = truth
        self.state = state
        self.now = now

    def run(self) -> RunReport:
        report = RunReport()
        for reader in self.readers:
            last = self.state.get_cursor(reader.name)
            observations = list(reader.read_since(last))
            if not observations:
                continue
            report.sources_read += 1
            report.observations += len(observations)

            candidates = self.llm.extract(observations)
            report.candidates += len(candidates)

            for c in candidates:
                fp = fingerprint(c.text, c.entities)
                corroboration = 1
                for src, _ in c.source_refs:
                    corroboration = self.state.record_seen(fp, src)
                now_ts = int(self.now())
                self.long_term.upsert(
                    Memory(fp, c.text, c.entities, c.kind, now_ts)
                )
                report.long_term_written += 1
                if qualifies_for_truth(c, corroboration):
                    self.truth.upsert(
                        Fact(fp, c.text, c.entities, c.source_refs, now_ts)
                    )
                    report.facts_written += 1

            for o in observations:
                if not self.state.already_distilled(o.source, o.source_id):
                    self.state.mark_distilled(o.source, o.source_id, int(self.now()))
            self.state.set_cursor(reader.name, max(o.timestamp for o in observations))
        return report
```

- [ ] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_pipeline.py -q`
Expected: PASS（4 passed）。

- [ ] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/pipeline.py pkg/memory/distiller/tests/test_pipeline.py
git commit -m "feat(distiller): add distillation pipeline orchestrator"
```

---

## Phase 7：保留清理 (Retention)

### Task 7: retention.sweep（30 天）

**Files:**

- Create: `pkg/memory/distiller/src/distiller/retention.py`
- Test: `pkg/memory/distiller/tests/test_retention.py`

- [x] **Step 1: 寫失敗測試**

`tests/test_retention.py`:

```python
from distiller.retention import sweep
from distiller.state import StateStore
from tests.conftest import FakePruner

DAY = 86400


def test_sweep_prunes_only_items_older_than_max_age(tmp_path):
    state = StateStore(tmp_path / "s.db")
    now = 100 * DAY
    state.mark_distilled("gbrain-working", "old.md", at=now - 31 * DAY)
    state.mark_distilled("gbrain-working", "fresh.md", at=now - 5 * DAY)
    pruner = FakePruner()
    pruned = sweep(state, {"gbrain-working": pruner}, now=now, max_age_days=30)
    assert pruned == 1
    assert pruner.pruned == [("gbrain-working", "old.md")]
    assert state.already_distilled("gbrain-working", "old.md") is False
    assert state.already_distilled("gbrain-working", "fresh.md") is True
    state.close()


def test_sweep_skips_sources_without_pruner_but_still_drops_state(tmp_path):
    state = StateStore(tmp_path / "s.db")
    now = 100 * DAY
    state.mark_distilled("claude-mem", "x", at=now - 40 * DAY)
    pruned = sweep(state, {}, now=now, max_age_days=30)
    assert pruned == 1
    assert state.already_distilled("claude-mem", "x") is False
    state.close()
```

- [x] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_retention.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.retention`）。

- [x] **Step 3: 實作 retention**

`src/distiller/retention.py`:

```python
DAY_SECONDS = 86400


def sweep(state, pruners: dict, now: int, max_age_days: int = 30) -> int:
    """清理蒸餾後超過 max_age_days 的來源項目，回傳清理筆數。"""
    cutoff = int(now) - max_age_days * DAY_SECONDS
    count = 0
    for source, source_id in state.due_for_prune(cutoff):
        pruner = pruners.get(source)
        if pruner is not None:
            pruner.prune(source, source_id)
        state.drop_distilled(source, source_id)
        count += 1
    return count
```

- [x] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_retention.py -q`
Expected: PASS（2 passed）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/retention.py pkg/memory/distiller/tests/test_retention.py
git commit -m "feat(distiller): add 30-day retention sweep"
```

---

## Phase 8：gbrain/working 來源讀取器 + pruner

### Task 8: GbrainWorkingReader（檔案系統，可完整實作）

**Files:**

- Create: `pkg/memory/distiller/src/distiller/adapters/__init__.py`
- Create: `pkg/memory/distiller/src/distiller/adapters/gbrain_working.py`
- Test: `pkg/memory/distiller/tests/test_gbrain_working.py`

- [x] **Step 1: 寫失敗測試**

`tests/test_gbrain_working.py`:

```python
import os
from distiller.adapters.gbrain_working import GbrainWorkingReader


def test_reads_markdown_files_as_observations(tmp_path):
    (tmp_path / "alice.md").write_text("alice context")
    sub = tmp_path / "topics"
    sub.mkdir()
    (sub / "trip.md").write_text("trip notes")
    reader = GbrainWorkingReader(tmp_path)
    obs = sorted(reader.read_since(0), key=lambda o: o.source_id)
    assert reader.name == "gbrain-working"
    assert [o.source_id for o in obs] == ["alice.md", "topics/trip.md"]
    assert obs[0].text == "alice context"


def test_read_since_filters_by_mtime(tmp_path):
    p = tmp_path / "old.md"
    p.write_text("old")
    os.utime(p, (1000, 1000))
    reader = GbrainWorkingReader(tmp_path)
    assert list(reader.read_since(2000)) == []
    assert len(list(reader.read_since(500))) == 1


def test_prune_deletes_file(tmp_path):
    p = tmp_path / "gone.md"
    p.write_text("bye")
    reader = GbrainWorkingReader(tmp_path)
    reader.prune("gbrain-working", "gone.md")
    assert not p.exists()
```

- [x] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_gbrain_working.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.adapters.gbrain_working`）。

- [x] **Step 3: 實作 reader（兼任 pruner）**

`src/distiller/adapters/__init__.py`:

```python

```

（空檔，標記為 package）

`src/distiller/adapters/gbrain_working.py`:

```python
from pathlib import Path

from ..models import Observation


class GbrainWorkingReader:
    """讀 gbrain/working 目錄下的 *.md，每檔一筆 Observation；亦可作為 Pruner。"""
    name = "gbrain-working"

    def __init__(self, root):
        self.root = Path(root)

    def read_since(self, last_ts: int):
        if not self.root.exists():
            return
        for path in sorted(self.root.rglob("*.md")):
            ts = int(path.stat().st_mtime)
            if ts > last_ts:
                rel = str(path.relative_to(self.root))
                yield Observation(
                    source=self.name,
                    source_id=rel,
                    timestamp=ts,
                    text=path.read_text(),
                )

    def prune(self, source: str, source_id: str) -> None:
        target = self.root / source_id
        if target.exists():
            target.unlink()
```

- [x] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_gbrain_working.py -q`
Expected: PASS（3 passed）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/adapters pkg/memory/distiller/tests/test_gbrain_working.py
git commit -m "feat(distiller): add gbrain/working source reader and pruner"
```

---

## Phase 9：claude-mem 來源讀取器（設定驅動 SQLite）

### Task 9: ClaudeMemReader

欄位名由 Phase 0 探查後寫入設定，避免硬寫捏造 schema。

**Files:**

- Create: `pkg/memory/distiller/src/distiller/adapters/claude_mem.py`
- Test: `pkg/memory/distiller/tests/test_claude_mem.py`

- [x] **Step 1: 寫失敗測試（用臨時 sqlite 模擬 claude-mem 表）**

`tests/test_claude_mem.py`:

```python
import sqlite3
from distiller.adapters.claude_mem import ClaudeMemReader


def _seed(db_path):
    con = sqlite3.connect(db_path)
    con.execute("CREATE TABLE observations (id INTEGER PRIMARY KEY, ts INTEGER, content TEXT)")
    con.executemany(
        "INSERT INTO observations(id, ts, content) VALUES(?,?,?)",
        [(1, 100, "first"), (2, 200, "second")],
    )
    con.commit()
    con.close()


def test_reads_rows_after_cursor(tmp_path):
    db = tmp_path / "claude-mem.db"
    _seed(db)
    reader = ClaudeMemReader(
        db_path=str(db), table="observations",
        id_col="id", ts_col="ts", text_col="content",
    )
    obs = list(reader.read_since(100))
    assert reader.name == "claude-mem"
    assert [o.source_id for o in obs] == ["2"]
    assert obs[0].text == "second"
    assert obs[0].timestamp == 200


def test_reads_all_when_cursor_zero(tmp_path):
    db = tmp_path / "claude-mem.db"
    _seed(db)
    reader = ClaudeMemReader(str(db), "observations", "id", "ts", "content")
    assert len(list(reader.read_since(0))) == 2
```

- [x] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_claude_mem.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.adapters.claude_mem`）。

- [x] **Step 3: 實作 reader**

`src/distiller/adapters/claude_mem.py`:

```python
import sqlite3

from ..models import Observation


class ClaudeMemReader:
    """讀 claude-mem 的 SQLite 觀察紀錄。表/欄名由設定提供（Phase 0 探查）。"""
    name = "claude-mem"

    def __init__(self, db_path, table, id_col, ts_col, text_col):
        self.db_path = db_path
        self.table = table
        self.id_col = id_col
        self.ts_col = ts_col
        self.text_col = text_col

    def read_since(self, last_ts: int):
        con = sqlite3.connect(self.db_path)
        try:
            query = (
                f"SELECT {self.id_col}, {self.ts_col}, {self.text_col} "
                f"FROM {self.table} WHERE {self.ts_col} > ? ORDER BY {self.ts_col}"
            )
            for sid, ts, text in con.execute(query, (last_ts,)):
                yield Observation(
                    source=self.name,
                    source_id=str(sid),
                    timestamp=int(ts),
                    text=text or "",
                )
        finally:
            con.close()
```

注意：表/欄名來自受信任的設定檔（非外部輸入），故以 f-string 組裝可接受。

- [x] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_claude_mem.py -q`
Expected: PASS（2 passed）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/adapters/claude_mem.py pkg/memory/distiller/tests/test_claude_mem.py
git commit -m "feat(distiller): add config-driven claude-mem sqlite reader"
```

---

## Phase 10：agentmemory / mempalace 寫入器（設定驅動 subprocess）

### Task 10: CommandWriter + AgentMemoryStore + MempalaceStore

實際寫入指令由 Phase 0 探查後填入設定的指令模板，透過 stdin 傳 JSON。

**Files:**

- Create: `pkg/memory/distiller/src/distiller/adapters/command_writer.py`
- Test: `pkg/memory/distiller/tests/test_command_writer.py`

- [x] **Step 1: 寫失敗測試**

`tests/test_command_writer.py`:

```python
import json
from distiller.models import Memory, Fact
from distiller.adapters.command_writer import CommandWriter, AgentMemoryStore, MempalaceStore


class RecordingWriter:
    def __init__(self):
        self.sent = []
    def send(self, record):
        self.sent.append(record)


def test_agentmemory_store_serializes_memory():
    w = RecordingWriter()
    store = AgentMemoryStore(w)
    store.upsert(Memory("fp1", "alice likes tea", ("alice",), "preference", 1000))
    assert w.sent == [{
        "fingerprint": "fp1", "text": "alice likes tea",
        "entities": ["alice"], "kind": "preference", "created_at": 1000,
    }]


def test_mempalace_store_serializes_fact():
    w = RecordingWriter()
    store = MempalaceStore(w)
    store.upsert(Fact("fp2", "born in Taipei", ("alice",), (("claude-mem", "9"),), 2000))
    assert w.sent == [{
        "fingerprint": "fp2", "text": "born in Taipei", "entities": ["alice"],
        "evidence": [["claude-mem", "9"]], "created_at": 2000,
    }]


def test_command_writer_invokes_subprocess_with_stdin_json(monkeypatch):
    captured = {}
    def fake_run(cmd, input, text, check):
        captured["cmd"] = cmd
        captured["input"] = input
        captured["check"] = check
    monkeypatch.setattr("distiller.adapters.command_writer.subprocess.run", fake_run)
    CommandWriter(["agentmemory", "ingest", "-"]).send({"a": 1})
    assert captured["cmd"] == ["agentmemory", "ingest", "-"]
    assert json.loads(captured["input"]) == {"a": 1}
    assert captured["check"] is True
```

- [x] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_command_writer.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.adapters.command_writer`）。

- [x] **Step 3: 實作 writers**

`src/distiller/adapters/command_writer.py`:

```python
import json
import subprocess

from ..models import Memory, Fact


class CommandWriter:
    """以設定的指令模板執行 subprocess，將 record 以 JSON 經 stdin 傳入。"""
    def __init__(self, cmd: list[str]):
        self.cmd = cmd

    def send(self, record: dict) -> None:
        subprocess.run(self.cmd, input=json.dumps(record), text=True, check=True)


class AgentMemoryStore:
    """LongTermStore：寫入 agentmemory。"""
    def __init__(self, writer):
        self.writer = writer

    def upsert(self, memory: Memory) -> None:
        self.writer.send({
            "fingerprint": memory.fingerprint,
            "text": memory.text,
            "entities": list(memory.entities),
            "kind": memory.kind,
            "created_at": memory.created_at,
        })


class MempalaceStore:
    """TruthStore：寫入 mempalace（只有 distiller 走此路徑）。"""
    def __init__(self, writer):
        self.writer = writer

    def upsert(self, fact: Fact) -> None:
        self.writer.send({
            "fingerprint": fact.fingerprint,
            "text": fact.text,
            "entities": list(fact.entities),
            "evidence": [list(e) for e in fact.evidence],
            "created_at": fact.created_at,
        })
```

- [x] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_command_writer.py -q`
Expected: PASS（3 passed）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/adapters/command_writer.py pkg/memory/distiller/tests/test_command_writer.py
git commit -m "feat(distiller): add config-driven agentmemory/mempalace writers"
```

---

## Phase 11：ollama 抽取 client

### Task 11: OllamaLLM

**Files:**

- Create: `pkg/memory/distiller/src/distiller/llm_ollama.py`
- Test: `pkg/memory/distiller/tests/test_llm_ollama.py`

- [x] **Step 1: 寫失敗測試（monkeypatch urlopen）**

`tests/test_llm_ollama.py`:

```python
import io
import json
import distiller.llm_ollama as mod
from distiller.llm_ollama import OllamaLLM
from distiller.models import Observation


def _fake_response(payload):
    class _Ctx:
        def __enter__(self_inner):
            return io.BytesIO(json.dumps(payload).encode())
        def __exit__(self_inner, *a):
            return False
    return _Ctx()


def test_extract_parses_candidates(monkeypatch):
    model_reply = {"message": {"content": json.dumps({"candidates": [
        {"text": "alice likes tea", "entities": ["alice"], "kind": "preference",
         "first_person": False, "confirmed_by_human": True},
    ]})}}
    monkeypatch.setattr(mod.urllib.request, "urlopen",
                        lambda req, timeout=None: _fake_response(model_reply))
    obs = [Observation("claude-mem", "1", 100, "user: alice likes tea")]
    cands = OllamaLLM(model="qwen2.5").extract(obs)
    assert len(cands) == 1
    assert cands[0].text == "alice likes tea"
    assert cands[0].confirmed_by_human is True
    assert cands[0].source_refs == (("claude-mem", "1"),)


def test_extract_empty_when_no_candidates(monkeypatch):
    model_reply = {"message": {"content": json.dumps({"candidates": []})}}
    monkeypatch.setattr(mod.urllib.request, "urlopen",
                        lambda req, timeout=None: _fake_response(model_reply))
    obs = [Observation("claude-mem", "1", 100, "chit chat")]
    assert OllamaLLM(model="qwen2.5").extract(obs) == []
```

- [x] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_llm_ollama.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.llm_ollama`）。

- [x] **Step 3: 實作 OllamaLLM**

`src/distiller/llm_ollama.py`:

```python
import json
import urllib.request

from .models import Observation, Candidate

EXTRACT_SYSTEM = (
    "You extract durable, reusable memories from agent/chat observations. "
    'Return ONLY a JSON object: {"candidates": [...]}. Each candidate has: '
    "text (verbatim statement worth remembering), "
    "entities (list of canonical names: people/projects/topics), "
    'kind (one of "fact","experience","preference","inference"), '
    "first_person (true if it is the human's own first-person life fact/experience), "
    "confirmed_by_human (true only if the human explicitly confirmed it). "
    "Omit chit-chat and transient operational noise."
)


class OllamaLLM:
    def __init__(self, model: str, host: str = "http://localhost:11434", timeout: int = 120):
        self.model = model
        self.host = host.rstrip("/")
        self.timeout = timeout

    def extract(self, observations: list[Observation]) -> list[Candidate]:
        joined = "\n\n".join(
            f"[{o.source}:{o.source_id}] {o.text}" for o in observations
        )
        payload = {
            "model": self.model,
            "messages": [
                {"role": "system", "content": EXTRACT_SYSTEM},
                {"role": "user", "content": joined},
            ],
            "format": "json",
            "stream": False,
        }
        req = urllib.request.Request(
            self.host + "/api/chat",
            data=json.dumps(payload).encode(),
            headers={"Content-Type": "application/json"},
        )
        with urllib.request.urlopen(req, timeout=self.timeout) as resp:
            body = json.loads(resp.read())
        content = json.loads(body["message"]["content"])
        refs = tuple((o.source, o.source_id) for o in observations)
        results = []
        for c in content.get("candidates", []):
            results.append(
                Candidate(
                    text=c["text"],
                    entities=tuple(c.get("entities", [])),
                    kind=c.get("kind", "fact"),
                    source_refs=refs,
                    first_person=bool(c.get("first_person", False)),
                    confirmed_by_human=bool(c.get("confirmed_by_human", False)),
                )
            )
        return results
```

- [x] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_llm_ollama.py -q`
Expected: PASS（2 passed）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/llm_ollama.py pkg/memory/distiller/tests/test_llm_ollama.py
git commit -m "feat(distiller): add ollama extraction client"
```

---

## Phase 12：設定載入 + CLI 進入點

### Task 12: config + cli（組裝真實 adapter）

**Files:**

- Create: `pkg/memory/distiller/src/distiller/config.py`
- Create: `pkg/memory/distiller/src/distiller/cli.py`
- Create: `pkg/memory/distiller/config.sample.toml`

- [x] **Step 1: 寫設定範本**

`config.sample.toml`（值由 Phase 0 探查後替換 `<...>`）:

```toml
[state]
db_path = "~/.distiller/state.db"

[retention]
max_age_days = 30

[llm]
model = "qwen2.5"
host = "http://localhost:11434"

[sources.claude_mem]
db_path = "~/.claude-mem/<db-file>"
table = "<observations_table>"
id_col = "<id_col>"
ts_col = "<ts_col>"
text_col = "<text_col>"

[sources.gbrain_working]
root = "~/.gbrain/working"

[stores.agentmemory]
cmd = ["<agentmemory-ingest-command>", "-"]

[stores.mempalace]
cmd = ["<mempalace-ingest-command>", "-"]
```

- [x] **Step 2: 實作 config（載入 + 組裝）**

`src/distiller/config.py`:

```python
import os
import tomllib
from pathlib import Path

from .adapters.claude_mem import ClaudeMemReader
from .adapters.command_writer import CommandWriter, AgentMemoryStore, MempalaceStore
from .adapters.gbrain_working import GbrainWorkingReader
from .llm_ollama import OllamaLLM
from .state import StateStore


def _expand(p: str) -> str:
    return os.path.expanduser(p)


class AppConfig:
    def __init__(self, data: dict):
        self.data = data

    @classmethod
    def load(cls, path: str) -> "AppConfig":
        with open(_expand(path), "rb") as fh:
            return cls(tomllib.load(fh))

    def state(self) -> StateStore:
        return StateStore(_expand(self.data["state"]["db_path"]))

    def max_age_days(self) -> int:
        return int(self.data.get("retention", {}).get("max_age_days", 30))

    def llm(self) -> OllamaLLM:
        cfg = self.data["llm"]
        return OllamaLLM(model=cfg["model"], host=cfg.get("host", "http://localhost:11434"))

    def readers(self):
        readers = []
        cm = self.data["sources"].get("claude_mem")
        if cm:
            readers.append(ClaudeMemReader(
                db_path=_expand(cm["db_path"]), table=cm["table"],
                id_col=cm["id_col"], ts_col=cm["ts_col"], text_col=cm["text_col"],
            ))
        gb = self.data["sources"].get("gbrain_working")
        if gb:
            readers.append(GbrainWorkingReader(_expand(gb["root"])))
        return readers

    def pruners(self) -> dict:
        pruners = {}
        gb = self.data["sources"].get("gbrain_working")
        if gb:
            pruners["gbrain-working"] = GbrainWorkingReader(_expand(gb["root"]))
        return pruners

    def long_term(self) -> AgentMemoryStore:
        return AgentMemoryStore(CommandWriter(self.data["stores"]["agentmemory"]["cmd"]))

    def truth(self) -> MempalaceStore:
        return MempalaceStore(CommandWriter(self.data["stores"]["mempalace"]["cmd"]))
```

- [x] **Step 3: 實作 cli**

`src/distiller/cli.py`:

```python
import argparse
import sys
import time

from .config import AppConfig
from .pipeline import Distiller
from .retention import sweep

DEFAULT_CONFIG = "~/.distiller/config.toml"


def main(argv=None) -> int:
    parser = argparse.ArgumentParser(prog="distiller")
    parser.add_argument("command", choices=["run", "retain"])
    parser.add_argument("--config", default=DEFAULT_CONFIG)
    parser.add_argument("--no-retain", action="store_true",
                        help="run 時不接著做保留清理")
    args = parser.parse_args(argv)

    cfg = AppConfig.load(args.config)
    state = cfg.state()
    try:
        if args.command == "run":
            report = Distiller(
                cfg.readers(), cfg.llm(), cfg.long_term(), cfg.truth(), state
            ).run()
            print(
                f"[distiller] sources={report.sources_read} obs={report.observations} "
                f"cand={report.candidates} long_term={report.long_term_written} "
                f"facts={report.facts_written}"
            )
            if not args.no_retain:
                pruned = sweep(state, cfg.pruners(), now=int(time.time()),
                               max_age_days=cfg.max_age_days())
                print(f"[distiller] pruned={pruned}")
        elif args.command == "retain":
            pruned = sweep(state, cfg.pruners(), now=int(time.time()),
                           max_age_days=cfg.max_age_days())
            print(f"[distiller] pruned={pruned}")
    finally:
        state.close()
    return 0


if __name__ == "__main__":
    sys.exit(main())
```

- [x] **Step 4: 全套件測試 + CLI 煙霧測試（--help）**

Run:

```bash
cd pkg/memory/distiller && uv run pytest -q && uv run distiller --help
```

Expected: 全部測試 PASS；`distiller --help` 印出 usage（含 run/retain）。

- [x] **Step 5: Commit**

```bash
git add pkg/memory/distiller/src/distiller/config.py pkg/memory/distiller/src/distiller/cli.py pkg/memory/distiller/config.sample.toml
git commit -m "feat(distiller): add config loader and CLI entrypoint"
```

---

## Phase 13：佈署接線（run.sh + cron）

### Task 13: 安裝 distiller、建執行期目錄、安裝每日 cron

**Files:**

- Create: `pkg/memory/distiller/scripts/install_cron.sh`
- Modify: `run.sh`（新增 distiller 區塊）

- [x] **Step 1: 寫 cron 安裝腳本**

`pkg/memory/distiller/scripts/install_cron.sh`:

```bash
#!/bin/bash
set -euo pipefail

DISTILLER_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CRON_LINE="0 3 * * * cd $DISTILLER_DIR && /usr/bin/env uv run distiller run >> $HOME/.distiller/logs/run.log 2>&1"

# 移除舊的 distiller 排程後重新加入（冪等）
( crontab -l 2>/dev/null | grep -v 'distiller run' ; echo "$CRON_LINE" ) | crontab -
echo "Installed daily 03:00 distiller cron:"
crontab -l | grep 'distiller run'
```

- [x] **Step 2: 在 run.sh 新增 distiller 區塊**

在 `run.sh` 末端（CCStatusline 區塊之後）新增：

```bash
# Distiller (agent memory)
mkdir -p "$HOME/.distiller/logs"
if [ ! -f "$HOME/.distiller/config.toml" ]; then
    cp "$(pwd)/pkg/memory/distiller/config.sample.toml" "$HOME/.distiller/config.toml"
    echo "[run.sh] Created ~/.distiller/config.toml — 請依 INTERFACES.md 填入 <...> 佔位值"
fi
ln -sf "$HOME/.distiller/config.toml" "$(pwd)/config/"
( cd "$(pwd)/pkg/memory/distiller" && uv sync )
```

- [x] **Step 3: 驗證安裝流程（伺服器）**

Run:

```bash
bash run.sh
test -f "$HOME/.distiller/config.toml" && echo "config OK"
chmod +x pkg/memory/distiller/scripts/install_cron.sh
bash pkg/memory/distiller/scripts/install_cron.sh
```

Expected: 印出 `config OK`，並列出已安裝的 `0 3 * * * ... distiller run` cron。

- [x] **Step 4: 端到端 dry-run（填好 config 後）**

Run:

```bash
cd pkg/memory/distiller && uv run distiller run --no-retain --config ~/.distiller/config.toml
```

Expected: 印出 `[distiller] sources=... obs=... cand=... long_term=... facts=...`，無例外。

- [x] **Step 5: Commit**

```bash
git add run.sh pkg/memory/distiller/scripts/install_cron.sh
git commit -m "chore(distiller): wire install into run.sh and add daily cron installer"
```

---

## Phase 14：同步拓樸 (syncthing)

### Task 14: 設定兩個單向 folder 並驗證

syncthing 的 folder 設定透過其 Web GUI 或 REST API 完成；本任務提供程序與驗證指令。

**Files:**

- Create: `pkg/memory/distiller/scripts/syncthing.md`

- [ ] **Step 1: 撰寫 syncthing 設定程序文件**

`pkg/memory/distiller/scripts/syncthing.md` 內容需涵蓋：

```markdown
# Syncthing 同步拓樸

兩個單向 folder（避免雙向 live-DB 損毀）：

## Folder A：claude-mem 筆電 → 伺服器

- 筆電端：新增 folder 指向 `~/.claude-mem`，Folder Type = `Send Only`
- 伺服器端：接收到 `~/.distiller/incoming/claude-mem`（供 distiller 讀為 claude-mem 快照）
- 在筆電 folder 加 `.stignore`：忽略鎖檔/暫存（`*-wal`, `*-shm`），降低同步進行中讀到半寫狀態的機率

## Folder B：agentmemory 伺服器 → 筆電

- 伺服器端：新增 folder 指向 agentmemory 資料目錄，Folder Type = `Send Only`
- 筆電端：Folder Type = `Receive Only`（唯讀查詢）
```

（伺服器端 claude-mem 快照路徑 `~/.distiller/incoming/claude-mem/<db>` 需回填到 `~/.distiller/config.toml` 的 `[sources.claude_mem].db_path`。）

- [ ] **Step 2: 確認 syncthing 已安裝並列出裝置**

Run:

```bash
syncthing --version
curl -s -H "X-API-Key: $(syncthing cli config gui apikey get 2>/dev/null)" http://localhost:8384/rest/system/connections | jq '.connections | keys'
```

Expected: 印出 syncthing 版本與已連線裝置 id（含對端機器）。

- [ ] **Step 3: 依文件在兩端建立 Folder A、Folder B（Web GUI <http://localhost:8384）>**

完成後驗證 folder 狀態：

```bash
curl -s -H "X-API-Key: <apikey>" "http://localhost:8384/rest/db/status?folder=<folderA-id>" | jq '{state, needFiles, globalFiles}'
```

Expected: `state` 為 `idle`，`needFiles` 為 0（已同步）。

- [ ] **Step 4: 驗證資料確實落地**

Run（伺服器）：

```bash
ls -la "$HOME/.distiller/incoming/claude-mem/"
```

Expected: 看到從筆電同步來的 claude-mem DB 檔。

- [ ] **Step 5: Commit**

```bash
git add pkg/memory/distiller/scripts/syncthing.md
git commit -m "docs(distiller): document syncthing one-way sync topology"
```

---

## Phase 15：hermes 對話連續性接線

### Task 15: gbrain/working 讀寫 helper + hermes 接線說明

hermes 每則訊息：載入該聯絡人/主題的 working note ＋ agentmemory 查詢結果以重建脈絡；回覆後把新交流追加回 working note。

**Files:**

- Create: `pkg/memory/distiller/src/distiller/continuity.py`
- Test: `pkg/memory/distiller/tests/test_continuity.py`
- Create: `pkg/memory/hermes-continuity/README.md`

- [ ] **Step 1: 寫失敗測試**

`tests/test_continuity.py`:

```python
from distiller.continuity import WorkingNotes


def test_append_creates_and_appends(tmp_path):
    notes = WorkingNotes(tmp_path)
    notes.append("alice", "user: hi")
    notes.append("alice", "assistant: hello")
    text = notes.load("alice")
    assert "user: hi" in text
    assert "assistant: hello" in text
    assert (tmp_path / "alice.md").exists()


def test_load_missing_returns_empty(tmp_path):
    assert WorkingNotes(tmp_path).load("nobody") == ""


def test_contact_is_slugified(tmp_path):
    notes = WorkingNotes(tmp_path)
    notes.append("Alice Wang", "x")
    assert (tmp_path / "alice-wang.md").exists()
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_continuity.py -q`
Expected: FAIL（`ModuleNotFoundError: distiller.continuity`）。

- [ ] **Step 3: 實作 WorkingNotes**

`src/distiller/continuity.py`:

```python
import re
from pathlib import Path


def slugify(name: str) -> str:
    s = re.sub(r"[^\w\s-]", "", name.strip().lower())
    return re.sub(r"[\s_]+", "-", s)


class WorkingNotes:
    """hermes 的 gbrain/working running notes 讀寫（每聯絡人/主題一檔）。"""
    def __init__(self, root):
        self.root = Path(root)
        self.root.mkdir(parents=True, exist_ok=True)

    def _path(self, contact: str) -> Path:
        return self.root / f"{slugify(contact)}.md"

    def load(self, contact: str) -> str:
        p = self._path(contact)
        return p.read_text() if p.exists() else ""

    def append(self, contact: str, line: str) -> None:
        p = self._path(contact)
        with p.open("a") as fh:
            fh.write(line.rstrip("\n") + "\n")
```

- [ ] **Step 4: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_continuity.py -q`
Expected: PASS（3 passed）。

- [ ] **Step 5: 撰寫 hermes 接線說明**

`pkg/memory/hermes-continuity/README.md` 需涵蓋：

```markdown
# Hermes 對話連續性接線

依 Phase 0 (INTERFACES.md) 確認的 hermes 擴充點，於每則訊息：

1. Pre-message：以聯絡人/主題為 key
    - 讀 `WorkingNotes.load(contact)` 取回 running notes
    - 查 agentmemory（long-term）取回已蒸餾事實
    - 兩者注入 hermes 的對話上下文（AGENTS.md 指示 + pre-message script）
2. Post-message：`WorkingNotes.append(contact, "<role>: <text>")` 追加新交流
3. gbrain/working 不同步、每日蒸餾後超過 30 天清理（由 distiller 負責）

接線方式擇一（依 INTERFACES.md 結果）：

- hermes hooks（若支援 pre/post message hook）
- 或 wrapper script 包住 hermes 入口
- 或在 `~/.hermes/AGENTS.md` 指示 agent 主動讀寫 working note 路徑
```

- [ ] **Step 6: Commit**

```bash
git add pkg/memory/distiller/src/distiller/continuity.py pkg/memory/distiller/tests/test_continuity.py pkg/memory/hermes-continuity/README.md
git commit -m "feat(distiller): add gbrain/working continuity helper and hermes wiring doc"
```

---

## Phase 16：端到端驗證 (Validation)

### Task 16: 對照 spec 的 Validation 節做整體驗收

**Files:**

- Create: `pkg/memory/distiller/tests/test_e2e.py`

- [ ] **Step 1: 寫端到端測試（fake LLM，真實 reader/state/writer 用 fake 收集）**

`tests/test_e2e.py`:

```python
from distiller.models import Candidate
from distiller.pipeline import Distiller
from distiller.state import StateStore
from distiller.adapters.gbrain_working import GbrainWorkingReader
from tests.conftest import FakeLLM, FakeLongTerm, FakeTruth


def test_e2e_distill_then_idempotent_rerun(tmp_path):
    working = tmp_path / "working"
    working.mkdir()
    (working / "alice.md").write_text("user: I was born in Taipei")
    reader = GbrainWorkingReader(working)
    cand = Candidate(
        text="born in Taipei", entities=("alice",), kind="fact",
        source_refs=(("gbrain-working", "alice.md"),),
        first_person=True, confirmed_by_human=False,
    )
    llm = FakeLLM([cand])
    lt, truth = FakeLongTerm(), FakeTruth()
    state = StateStore(tmp_path / "s.db")

    r1 = Distiller([reader], llm, lt, truth, state).run()
    assert r1.long_term_written == 1
    assert r1.facts_written == 1            # first_person fact -> mempalace

    # 重跑：cursor 已前進，無新觀察 -> 不重複寫入（冪等）
    r2 = Distiller([reader], llm, lt, truth, state).run()
    assert r2.observations == 0
    assert len(lt.upserts) == 1
    assert len(truth.upserts) == 1
    state.close()
```

- [ ] **Step 2: 跑測試確認通過**

Run: `cd pkg/memory/distiller && uv run pytest tests/test_e2e.py -q`
Expected: PASS（1 passed）。

- [ ] **Step 3: 全套件綠燈**

Run: `cd pkg/memory/distiller && uv run pytest -q`
Expected: 全部 PASS。

- [ ] **Step 4: 對照 spec Validation 節手動驗收（伺服器，真實工具）**

依 spec [驗證] 節逐項確認並記錄結果：

```bash
# 端到端：產生暫存記憶 -> 同步 -> 跑 distiller -> 確認 agentmemory/mempalace 出現對應記憶
cd pkg/memory/distiller && uv run distiller run --config ~/.distiller/config.toml
# 寫入隔離：確認沒有 distiller 以外的路徑寫 mempalace（檢視 INTERFACES.md 列出的寫入點僅 distiller 使用）
# 連續性：模擬 hermes 不連續對話，確認 working note + agentmemory 能重建脈絡
```

Expected：agentmemory 出現 long-term 記憶、驗證事實出現在 mempalace、重跑冪等。

- [ ] **Step 5: Commit**

```bash
git add pkg/memory/distiller/tests/test_e2e.py
git commit -m "test(distiller): add end-to-end distill + idempotency test"
```

---

## Self-Review

**1. Spec coverage（spec 各節 → 對應任務）**

| Spec 節                                    | 對應任務                                                                                                      |
| ------------------------------------------ | ------------------------------------------------------------------------------------------------------------- |
| 元件角色：claude-mem 來源                  | Task 9                                                                                                        |
| 元件角色：gbrain/working 連結層            | Task 8、Task 15                                                                                               |
| 元件角色：agentmemory long-term            | Task 10                                                                                                       |
| 元件角色：mempalace 真實（write-isolated） | Task 10（只有 distiller 走 TruthStore）、Task 16 驗收寫入隔離                                                 |
| 蒸餾器契約：排程每日凌晨                   | Task 13（cron 0 3 \* \* \*）                                                                                  |
| 蒸餾器契約：輸入/輸出                      | Task 6、Task 9、Task 8、Task 10                                                                               |
| 蒸餾器契約：驗證政策                       | Task 4                                                                                                        |
| 蒸餾器契約：冪等性                         | Task 3、Task 5、Task 16                                                                                       |
| 蒸餾器契約：增量（cursor）                 | Task 5、Task 6                                                                                                |
| 保留：30 天                                | Task 7、Task 12（run 內含）、Task 13                                                                          |
| 同步拓樸：syncthing 單向                   | Task 14                                                                                                       |
| 查詢層                                     | agentmemory(Task 10)＋claude-mem 既有；筆電副本由 Task 14 同步                                                |
| hermes 連續性                              | Task 15                                                                                                       |
| 嵌入一致性（bge-m3）                       | Task 0 探查 + config.sample.toml 之 llm/embedding 設定（agentmemory/mempalace 向量後端一致，於 Phase 0 確認） |
| 已排除 memsearch / codegraph               | 計畫未納入（符合 spec）                                                                                       |

**2. Placeholder scan：** 程式碼步驟皆含完整可跑程式碼；外部工具未知處集中於 `config.sample.toml` 的 `<...>` 與 Phase 0 探查，並以設定驅動，非程式碼內 TODO。

**3. Type consistency：** `Observation/Candidate/Memory/Fact/RunReport`（models.py）跨 Task 6/8/9/10/11/16 一致；`SourceReader.read_since`、`LLMClient.extract`、`LongTermStore.upsert`、`TruthStore.upsert`、`Pruner.prune` 簽章跨 interfaces 與各 adapter 一致；`StateStore` 方法（get_cursor/set_cursor/record_seen/already_distilled/mark_distilled/due_for_prune/drop_distilled）於 Task 5 定義、Task 6/7 使用一致。

---

## 未決依賴（執行前提醒）

- Phase 0 必須先完成，Phase 9/10/14/15 的設定值與接線方式才有確切來源。
- 嵌入模型一致性：agentmemory / mempalace 的向量後端若無法設為 ollama bge-m3，需於 Phase 0 記錄替代方案。
- 伺服器 OS 若為 macOS，cron 可改用 launchd（Task 13 之 install_cron.sh 需對應調整）。
