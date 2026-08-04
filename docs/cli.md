# CLI 參考 (CLI Reference)

`cc-plugin` 的完整指令參考; quick start 見 [`README.md`](../README.md) 的使用方式章節。

## 記憶蒸餾 (Memory Distillation)

```bash
# 執行完整蒸餾管道（讀取 → 提取 → 寫入 → 清理）
cc-plugin distill

# 僅提取記憶（從 stdin 讀取 JSON 觀察值）
cc-plugin extract < observations.json

# 清理狀態（重置遊標、已見、已蒸餾紀錄）
cc-plugin reset
```

## 資料匯出 (Data Export)

```bash
# 匯出 mempalace 類別清單
cc-plugin export mempalace

# 匯出 mempalace 完整 Markdown 結構
cc-plugin export mempalace --data -o ./export

# 匯出 gbrain 觀察值（增量）
cc-plugin export gbrain

# 匯出 claude-mem 觀察值（全量）
cc-plugin export claudemem --all
```

`claudemem` 增量匯出使用獨立遊標; 從舊版本升級後, 第一次增量匯出會安全地完整
匯出既有資料一次。設計細節見
[`docs/specs/2026-07-16-claudemem-export-reliability-design.md`](specs/2026-07-16-claudemem-export-reliability-design.md)。

## Topology 知識圖譜 (Topology Graph)

```bash
# 驗證 topology-builder 參考圖譜
cc-plugin topology verify

# 查詢 entity 邊
cc-plugin topology query service-a

# 重算 backlinks 與 _index.md
cc-plugin topology rewrite --root <topology-root>
```

## 環境初始化 (Environment Initialization)

```bash
# 初始化軟連結與設定同步
chmod +x scripts/run.sh && ./scripts/run.sh

# 安裝技能至 AI Agents
npx skills add .
```

環境細節（前置需求、部署、設定同步範圍）見 [`docs/development.md`](development.md)。
