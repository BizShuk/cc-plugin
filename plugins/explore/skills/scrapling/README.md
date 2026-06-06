# 網頁爬蟲框架 (Scrapling)

`Scrapling` 是一個自適應網頁爬蟲與資料擷取框架 (Adaptive Web Scraping Framework)，支援單次請求、動態渲染、反爬蟲繞過 (Anti-bot Bypass)，以及分散式高並發爬蟲。

本說明文件介紹如何安裝並使用此工具。

## 核心特色 (Key Features)

- `自適應解析器`：當網頁結構改變時，自動重新定位元素。
- `反爬蟲繞過`：無需外部 API 或額外設定，開箱即用繞過 Cloudflare Turnstile 等反爬蟲系統。
- `隱身瀏覽器` (Stealthy Headless Browsing)：提供指紋偽裝與真實瀏覽器模擬。
- `爬蟲蜘蛛框架` (Spider Framework)：支援高並發、暫停與續傳、代理伺服器輪替 (Proxy Rotation) 等功能。

## 安裝與設定 (Setup)

需要 `Python 3.10+` 環境。

### 1. 使用 Pip 安裝

在您的 Python 虛擬環境中執行：

```bash
pip install "scrapling[all]>=0.4.8"
```

下載瀏覽器依賴：

```bash
scrapling install --force
```

### 2. 使用 Docker 映像檔 (Docker Image)

不需安裝 Python 環境即可直接使用命令列工具：

```bash
docker pull pyd4vinci/scrapling
```

## 命令列介面使用 (CLI Usage)

使用 `scrapling extract` 指令群可直接擷取網頁內容，支援輸出為 Markdown (`.md`)、HTML (`.html`) 或純文字 (`.txt`)。

### 常用指令

- `get`：適用於簡單靜態網頁、部落格或新聞文章。

    ```bash
    scrapling extract get "https://example.com" page.md --ai-targeted
    ```

- `fetch`：適用於現代單頁式應用程式 (SPA) 或動態 JavaScript 渲染的網站。

    ```bash
    scrapling extract fetch "https://example.com" page.md --network-idle
    ```

- `stealthy-fetch`：適用於受 Cloudflare 或防爬蟲保護的網站。

    ```bash
    scrapling extract stealthy-fetch "https://example.com" page.md --solve-cloudflare
    ```

> [!IMPORTANT]
> 命令列爬取時，強烈建議使用 `--ai-targeted` 參數以過濾廣告與無用元素，並保護避免提示詞注入 (Prompt Injection)。

## Python 程式開發 (Code Usage)

### 1. 基本擷取 (Fetcher & Session)

```python
from scrapling.fetchers import Fetcher, StealthyFetcher

# 靜態 GET 請求（偽裝 Chrome TLS 指紋）
page = Fetcher.get('https://example.com')
title = page.css('h1::text').get()

# 隱身瀏覽器擷取（繞過 Cloudflare）
page = StealthyFetcher.fetch('https://nopecha.com/demo/cloudflare')
links = page.css('#padded_content a').getall()
```

### 2. 爬蟲蜘蛛 (Spiders)

繼承 `Spider` 即可實作具備並發請求、自動遵守 `robots.txt`、以及支援以 `Ctrl+C` 暫停與續傳功能的爬蟲：

```python
from scrapling.spiders import Spider, Response

class MySpider(Spider):
    name = "myspider"
    start_urls = ["https://quotes.toscrape.com/"]
    concurrent_requests = 5
    robots_txt_obey = True

    async def parse(self, response: Response):
        for quote in response.css('.quote'):
            yield {
                "text": quote.css('.text::text').get(),
                "author": quote.css('.author::text').get(),
            }
```

執行並儲存進度以支援續傳：

```python
MySpider(crawldir="./crawl_data").start()
```

## 開發守則 (Guardrails)

- 遵守 `robots.txt` 與網站服務條款 (ToS)。
- 在大規模爬取時設定合理的下載延遲 (`download_delay`)，避免對伺服器造成負擔。
- 未經授權請勿嘗試繞過付費牆 (Paywalls) 或登入驗證。
