# MarkItDown Python API

CLI 之外，`markitdown` 也可以當函式庫用 — 需要在既有 Python pipeline 裡轉檔、
或要對回傳結果做後處理時使用。

```python
from markitdown import MarkItDown

md = MarkItDown()
result = md.convert("report.pdf")                # local file
result = md.convert_url("https://example.com")   # URL
print(result.text_content)
```

- `convert(path)` 吃本機檔案路徑，副檔名決定 converter。
- `convert_url(url)` 抓取靜態 HTML；JS 渲染的頁面拿不到內容。
- 回傳物件的 `text_content` 是 Markdown 字串。

啟用 plugin（預設關閉，等同 CLI 的 `--use-plugins`）：

```python
md = MarkItDown(enable_plugins=True)
```

其餘可用的 extras 與 plugin 清單見 SKILL.md 的 Installation 一節。
