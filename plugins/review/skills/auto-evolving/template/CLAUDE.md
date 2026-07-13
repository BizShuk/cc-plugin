# 自演化知識引擎 — 技術脈絡

## 項目結構

```tree
research/
├── engine/                 # 引擎核心：System Prompt + 配置 + 哲學參考
│   ├── system_prompt.md    # System Prompt 正本
│   ├── config.yaml         # 執行參數（衰減係數、閾值、週期間隔）
│   └── philosophy.md       # 四大哲學脈參考
│
├── foundation/             # 知識基礎：三層沉澱 + 淘汰歸檔
│   ├── axioms/             # 公理層 — AX-<NNN>-<topic>.md
│   ├── principles/         # 原則層 — PR-<NNN>-<topic>.md
│   ├── hypotheses/         # 假說層 — HY-<NNN>-<topic>.md
│   ├── retired/            # 淘汰層 — RT-<NNN>-<topic>.md
│   └── confidence.yaml     # 信心度登記簿
│
├── cycles/                 # 執行週期記錄，按月歸檔
│   └── YYYY-MM/            # 月份子目錄
│       └── YYYY-MM-DD-NNN.md
│
├── branches/               # 四策略分支演化追蹤
│   ├── harmony/            # 策略 A — 和諧融入
│   ├── dialectic/          # 策略 B — 辯證衝突
│   ├── velocity/           # 策略 C — 速度優先
│   ├── incremental/        # 策略 D — 漸進驗證
│   └── scoreboard.md       # 策略表現記分板
│
├── plans/                  # 進行中計畫
├── docs/
│   ├── memory/             # 歷史操作與決策 retrospective
│   ├── backlog/            # 待辦想法
│   └── specs/              # 既有設計與規格
├── scripts/                # 自動化腳本
└── tmp/                    # 暫存資料
```

## 關鍵決策

1. 知識分三層存儲（公理/原則/假說）+ 淘汰歸檔，不刪除任何知識
2. 信心度衰減機制：每 7 個週期未被引用或驗證，信心度 -= 0.05
3. 四策略分支獨立演化，不強求統一，讓實踐裁判
4. 評估採四軸（系統/業務/認知/反脆弱），總分 40 分制
5. 重大變更必須暫停通知人類批准
6. 週期記錄按月份歸檔，避免單一目錄過度膨脹

## 慣例

- 知識檔案命名：`<前綴>-<三位序號>-<topic>.md`（如 `AX-001-idempotent-ops.md`）
- 週期檔案命名：`YYYY-MM-DD-<NNN>.md`（NNN 為當日序號）
- `confidence.yaml` 是所有知識條目的唯一信心度來源
- 哲學推論僅在合成與評估階段啟動，資料收集保持中性
- 讀取目標項目的 `CLAUDE.md` 和 `README.md` 是每次喚醒的必要前置動作
