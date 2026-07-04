# Feature: Topology 模組抽離與 CLI 化 (Topology Extraction & CLI)

## 問題陳述 (Problem Statement)

`model/` 套件同時承載了核心記憶狀態管理 (`StateStore`) 以及 Markdown 知識圖譜解析的 `Topology` 邏輯，違反單一職責原則：

- 領域職責過載：`model/topology.go` (185 行) 與 `model/topology_ops.go` (219 行) 屬於 Markdown 知識圖譜工具，與記憶蒸餾領域無關
- 孤立程式碼：`topology.go` 與 `topology_ops.go` 目前在核心 CLI 指令 (`cmd/`) 中完全沒有任何實質調用，僅在 `topology_test.go` 與 `topology_ops_test.go` 中被引用
- 缺乏 CLI 接入點：拓撲驗證與邊界檢查功能無法透過 `cc-plugin` 命令直接執行

## 受影響檔案 (Affected Files)

| 原始檔案路徑 | 目標檔案路徑 | 調整說明 |
| :--- | :--- | :--- |
| `model/topology.go` | `pkg/topology/topology.go` | 移動並變更 package 為 `topology` |
| `model/topology_ops.go` | `pkg/topology/topology_ops.go` | 同上 |
| `model/topology_test.go` | `pkg/topology/topology_test.go` | 對應單元測試遷移 |
| `model/topology_ops_test.go` | `pkg/topology/topology_ops_test.go` | 對應單元測試遷移 |
| — | `cmd/topology.go` | 新增 `cc-plugin topology` 命令，實作 `verify`、`rewrite` 與 `unlinked` 子命令 |

## 實作步驟 (Implementation Steps)

1. 在 `pkg/topology/` 下建立目錄，將 `model/topology.go` 與 `model/topology_ops.go` 移入
2. 將檔案的 package 宣告改為 `topology`，並更新所有 import 路徑
3. 將對應的 `*_test.go` 測試檔一併遷移
4. 新增 `cmd/topology.go`，實作 Cobra 子指令：
   - `cc-plugin topology verify --root <path>` — 驗證 Markdown 知識圖譜邊界
   - `cc-plugin topology rewrite --root <path>` — 重寫 backlink
   - `cc-plugin topology unlinked --root <path>` — 找出未連結的節點
5. 在 `cmd/root.go` 註冊 `topology` 子命令

## 驗證方式 (Verification)

- 執行 `go test ./pkg/topology/...` 確認解析與邊界驗證功能正常
- 執行 `go build -o cc-plugin main.go` 後，手動執行 `cc-plugin topology verify --root plugins/general/skills/topology-builder/references` 驗證輸出
- 確認所有 `cmd/*.go` 與 `model/` 中已無 `model/topology` 的 import 殘留

## 來源 (Source Plans)

- [`architecture-system-modularization.md`](architecture-system-modularization.md) §1 診斷 1-2, §3, §4, §6 Phase 2 + Phase 4