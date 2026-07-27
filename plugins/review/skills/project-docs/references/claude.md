# CLAUDE.md 模板與規則

## 模板 (Template)

```markdown
# <Project Name> — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

<實際目錄樹，2-3 層深>

## 技術棧 (Tech Stack)

- Language: <detected>
- Framework: <detected>
- Build tool: <detected>
- Key dependencies: <top 5-8 deps>

## 關鍵決策 (Key Decisions)

- Decision 1：為何選擇此做法（從程式碼模式推斷）
- Decision 2：...

## 模組對應 (Module Mapping)

把每個業務領域（從 README）對應到技術實作：

| 業務領域 (Domain) | 套件/模組 (Package/Module) | 進入點 (Entry Point) |
| ----------------- | -------------------------- | -------------------- |
| <Domain 1>        | `pkg/xxx`, `handler/yyy`   | `HandleXxx()`        |
| <Domain 2>        | `pkg/aaa`, `handler/bbb`   | `HandleAaa()`        |

## 開發指南 (Development Guide)

### 前置需求 (Prerequisites)

- Requirement 1
- Requirement 2

### 安裝 (Installation)

<專案實際的 install commands>

### 建置 (Build)

<精確的 build commands>

### 測試 (Test)

<精確的 test commands，或註明無測試>

### 部署 (Deploy)

<可偵測的部署方式，或「未偵測到部署設定 (No deployment config detected)」>

## 慣例 (Conventions)

- Naming: <detected patterns>
- Error handling: <detected patterns>
- Logging: <detected patterns>
- Testing: <detected patterns>
```
