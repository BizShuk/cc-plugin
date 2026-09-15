# README.todo 格式規範

```md
# TODO

- [ ] xxx

## <feature>

- [ ] yyy

## Archive
```

## 規則 (Rules)

- 待辦是`唯一`的去處：`README.md` 與 `CLAUDE.md` 不寫 `- [ ]` 項目，
  `bootstrap` 與 `refresh` 分析出的改善建議一律落在這裡
- 每筆待辦具體可執行，說明理由；依主題分組在 `## <feature>` 之下
- 完成的項目勾選為 `- [x]` 並移入 `## Archive`，由 `consolidate` 模式搬進
  `docs/CHANGELOG.md`；`## Archive` 標題一律保留
