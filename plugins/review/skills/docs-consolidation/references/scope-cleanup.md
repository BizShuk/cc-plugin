# docs-consolidation — 範疇清理 (Scope Cleanup)

對象是 `README.md` 與 `CLAUDE.md` `本身的內容`，不是歷史文件的數量。
判準表見 `SKILL.md` 的`範疇判準`章節；本檔是逐階段的指令與樣板。

## Phase S0 — Audit

`稽核先於改寫`。每一筆發現寫成 `doc 說 X → 實際 Y`，附檔名與行號，
經使用者確認後才進入後續階段。未經確認不得動手改寫。

### Step S0.0 — Token 基線 (Token Baseline)

改寫前先記錄正典文件的 token 數；S6 以`同一指令`重測，回報前後差異：

```bash
python3 -c 'import sys,tiktoken; e=tiktoken.get_encoding("o200k_base");
[print(len(e.encode(open(f).read())), f) for f in sys.argv[1:]]' README.md CLAUDE.md
```

無 `tiktoken` 時改用 `wc -c` 除以 4 估算，並在報告標註`估算值`。
基線只做`量測`，不做承諾 ——「不預告最終行數」對 token 數同樣適用。

### Step S0.1 — 掃描七類徵狀

```bash
# 1) 外部 repo：自承範疇的章節，以及對方的實作細節
grep -n "外部 repo\|external repo\|本節描述\|不對應本 repo" README.md CLAUDE.md

# 2) 歷史敘述
grep -n "已移除\|已刪除\|已解體\|已併回\|已下沉\|已脫離\|不再\|原 \`\|曾經" README.md CLAUDE.md

# 3) 易腐計數（行數、檔案數、容量、版本數）
grep -nE '[0-9]+ 行|[0-9]+ 個(檔案|模組|module|package)|[0-9]+ ?(MiB|KiB|MB|KB)|→[0-9]+' README.md CLAUDE.md

# 4) 可執行斷言（文件裡的 pass/fail 指令）
grep -n "必須為空\|必須非空\|只該有\|不該出現\|grep -\|go list -deps\|test ! -" README.md CLAUDE.md

# 5) 機器專屬路徑
grep -n "/Users/\|/home/\|~/projects/" README.md CLAUDE.md

# 6) 重複：兩檔都有的章節標題與結構樹
grep -c '^```tree\|^```text' README.md CLAUDE.md
diff <(grep '^## ' README.md) <(grep '^## ' CLAUDE.md)

# 7) 過長章節（細節下放候選）：逐章節統計行數
for f in README.md CLAUDE.md; do
  awk -v f="$f" '/^## /{if(s)print f": "c" 行  "s; s=$0; c=0; next}{c++}
                 END{if(s)print f": "c" 行  "s}' "$f"
done
```

第 7 類的判定：`使用方式`／`開發指南`類章節超過約 25 行即為`細節下放`候選，
目的地與保留規則見 [content-ownership.md](content-ownership.md) 的`細節下放`章節。

### Step S0.2 — 逐筆驗證文件宣稱

文件說的`每一件可驗證的事`都要實跑，不得目視採信。至少涵蓋：

```bash
# 結構樹宣稱 vs 實際目錄
ls <每個樹狀圖列出的目錄>

# 模組／套件計數宣稱 vs 實際
go list -m -f '{{.Dir}}' | wc -l        # Go workspace 為例

# 檔案級引用是否還存在
ls <文件提到的每個具體檔案路徑>

# 連結目標是否可解析
grep -oE '\]\([^)h][^)]*\)' README.md CLAUDE.md \
  | sed 's/.*](//; s/)$//; s/#.*//' | sort -u \
  | while read -r f; do [ -z "$f" ] || [ -e "$f" ] || echo "BROKEN: $f"; done
```

### Step S0.3 — 處置表

| 位置 | 內容 | 失效原因 | 目的地 | 實測 |
| ---- | ---- | -------- | ------ | ---- |

`實測`欄記錄支持該判定的指令輸出。沒有實測的列不得進入後續階段。

## Phase S1 — Automate

`可執行的斷言先落地，才有資格從文件刪掉`。順序反過來會產生無人把關的空窗期。

### Step S1.1 — 先跑一次文件裡的斷言

**文件裡的斷言預設為錯。** 沒人執行的斷言不只是沒用，它會`腐爛成錯誤資訊`，
而且會被下一個人原樣搬進測試，把錯誤固化。

兩類常見的錯誤斷言：

| 錯誤形態 | 例 | 正確寫法 |
| -------- | -- | -------- |
| 把`遞移閉包`當成`直接依賴` | 「`X/sub` 的 deps 不該出現 `core`」，但 `X` 自己就 import `core` | 斷言`白名單`：閉包內不得出現清單以外的項目 |
| 用文字比對代替語意檢查 | `grep 'vendor-name' core/*.go` 必須為空 | 斷言 `import`；註解與測試字串是合法的 |

先跑，紅的先修正`斷言本身`，再進入 S1.2。

### Step S1.2 — 落成測試或腳本

| 性質 | 去處 |
| ---- | ---- |
| 有 pass/fail 語意的不變式 | 測試（納入預設測試指令，如 `go test ./...`） |
| 沒有 pass/fail、只是流程 | `scripts/`（如跨 module 的 build 迴圈） |

驗證工具`必須唯讀`。會產出檔案的建置指令要導向暫存目錄並在結束時清理；
否則稽核工具本身會弄髒 repo（例：`go build ./...` 在單一 main package 時
會把執行檔寫進工作目錄）。

### Step S1.3 — 注入違規，證明測試會紅

守護`既有`不變式的測試天生全綠，**沒證明過會紅的 guard test 等於沒寫**。

```
1. 在受測邊界注入一筆違規（多一個 import、多一個依賴）
2. 跑測試 → 必須 FAIL，且訊息要`指名違規者`
3. 還原，確認 diff 為空
4. 重跑 → PASS
```

subtest 名稱不得含 `/`（會被 `-run` 當成巢狀分隔符），否則 `-run` 會靜默匹配
不到任何測試並印出 `no tests to run` —— 那看起來像成功。

## Phase S2 — Verify destination

**刪除前必須實測目的地已有該內容。** 這一步決定這是`刪除`還是`遺失`。

```bash
# 例：確認里程碑歷史已在 CHANGELOG，才刪 README 的里程碑表
grep -n "M1\|M2\|里程碑" docs/CHANGELOG.md
```

| 查核結果 | 動作 |
| -------- | ---- |
| 目的地`已有`同等內容 | 直接刪除來源，不搬移 |
| 目的地`沒有` | 先寫進目的地，確認寫入成功，才刪來源 |
| 目的地是`外部 repo` | 只留一行指向連結；不把內容複製到對方 repo（那是對方的工作） |

## Phase S3 — Cut

- 編輯以 `anchor 文字`定位，`不用行號`：前面的刪除會讓後面所有行號失效。
- 逐段刪除，每段刪完立即跑該段對應的驗證指令。
- 刪`歷史敘述`時保留句子的`現況`部分，只砍「曾經如何」的子句。
- 刪`重複章節`前，先把只在該處出現的內容併入 owner 檔，避免連帶遺失。

## Phase S4 — Sweep

**刪完重讀全文，不只看改動處。** 分段刪除看不見這三類殘留：

| 殘留 | 例 |
| ---- | -- |
| 懸空引用 | 表格被刪後，別處仍寫著 `(M2)`、`見上表` |
| 交叉連結 | 指向已刪章節的錨點連結 |
| 孤立列 | 表格只剩一列、章節只剩標題 |

## Phase S5 — Lint

```bash
# 連結可解析
grep -oE '\]\([^)h][^)]*\)' README.md CLAUDE.md \
  | sed 's/.*](//; s/)$//; s/#.*//' | sort -u \
  | while read -r f; do [ -z "$f" ] || [ -e "$f" ] || echo "BROKEN: $f"; done

# 機器路徑、外部 repo 細節、歷史敘述、易腐計數皆已清空
grep -n "/Users/\|/home/" README.md CLAUDE.md
grep -n "已移除\|已解體\|已併回\|已下沉" README.md CLAUDE.md

# 測試與腳本全綠，且 repo 未被弄髒
<專案的預設測試指令>
git status --short
```

`git status --short` 只該出現本次`預期修改`的檔案。出現其他檔案代表驗證工具
有寫入副作用，回 S1.2 修工具，不是接受這個結果。

## Phase S6 — Report

回報格式：

```
範疇清理 — <專案>

處置：
- 外部 repo 細節：<n> 行 → 刪除，保留 <k> 條消費端契約
- 歷史敘述：<n> 處 → 已在 docs/CHANGELOG.md（刪除，非搬移）
- 易腐計數：<n> 處 → 刪除
- 可執行斷言：<n> 條 → <測試檔>（其中 <m> 條原本就是錯的，已改寫）
- 重複章節：<清單> → 由 <owner 檔> 單一擁有
- 細節下放：<n> 行 → docs/cli.md、docs/development.md（正典檔留 quick start + 指標）
- 機器路徑：<n> 處 → 相對路徑

行數：README.md <a> → <b>；CLAUDE.md <c> → <d>
Token：README.md <a> → <b>（-<x>%）；CLAUDE.md <c> → <d>（-<y>%）
  （S0.0 同一指令重測；含下放後的 docs/cli.md、docs/development.md 新增量）
驗收：測試全綠 / 連結全解析 / working tree 只含預期檔案
未處理：<清單與理由>
```

**不預告最終行數。** 刪完才知道 —— 事先承諾的數字只能靠刪真契約來達成。
估錯了就照實回報實際值與剩下的是什麼，不要硬湊。
