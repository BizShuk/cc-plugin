# Scrapling Spiders — 大規模爬取 (Crawling Framework)

單頁抓取用 CLI 或 `python-api.md` 的 fetcher；需要`跟連結、併發、斷點續爬`時才用 spider。

## 基本 Spider

併發請求、多 session、可暫停續爬：

```python
from scrapling.spiders import Spider, Request, Response

class QuotesSpider(Spider):
    name = "quotes"
    start_urls = ["https://quotes.toscrape.com/"]
    concurrent_requests = 10
    robots_txt_obey = True  # Respect robots.txt rules

    async def parse(self, response: Response):
        for quote in response.css('.quote'):
            yield {
                "text": quote.css('.text::text').get(),
                "author": quote.css('.author::text').get(),
            }

        next_page = response.css('.next a')
        if next_page:
            yield response.follow(next_page[0].attrib['href'])

result = QuotesSpider().start()
print(f"Scraped {len(result.items)} quotes")
result.items.to_json("quotes.json")
```

## 單一 Spider 內混用多種 session

受保護的頁面走 stealth session，其餘走快速 session：

```python
from scrapling.spiders import Spider, Request, Response
from scrapling.fetchers import FetcherSession, AsyncStealthySession

class MultiSessionSpider(Spider):
    name = "multi"
    start_urls = ["https://example.com/"]

    def configure_sessions(self, manager):
        manager.add("fast", FetcherSession(impersonate="chrome"))
        manager.add("stealth", AsyncStealthySession(headless=True), lazy=True)

    async def parse(self, response: Response):
        for link in response.css('a::attr(href)').getall():
            # Route protected pages through the stealth session
            if "protected" in link:
                yield Request(link, sid="stealth")
            else:
                yield Request(link, sid="fast", callback=self.parse)  # explicit callback
```

## 暫停與續爬 (Checkpoints)

```python
QuotesSpider(crawldir="./crawl_data").start()
```

按 Ctrl+C 會優雅暫停，進度自動存檔。之後以`相同的` `crawldir` 重新啟動即從中斷處續跑。

## 開發模式 (Development Mode)

調整 `parse()` 邏輯期間，在 spider class 上設 `development_mode = True`：第一次執行把 response 快取到磁碟，後續執行改播放快取，可無限次重跑而不再打目標站。快取預設在 `.scrapling_cache/{spider.name}/`，可用 `development_cache_dir` 覆寫。上線的 spider 不得留著此旗標。

## 規則式爬取 (CrawlSpider)

需要「跟隨符合 regex 的連結」時用 `CrawlSpider`，不要自己寫連結抽取迴圈：

```python
from scrapling.spiders import CrawlSpider, CrawlRule, LinkExtractor

class BlogCrawler(CrawlSpider):
    name = "blog"
    start_urls = ["https://example.com"]

    def rules(self):
        return [
            CrawlRule(LinkExtractor(allow=r"/posts/"), callback=self.parse_post),
            CrawlRule(LinkExtractor(allow=r"/page/\d+/")),  # follow pagination, no callback
        ]

    async def parse_post(self, response):
        yield {"title": response.css("h1::text").get()}
```

## Sitemap 驅動爬取

`SitemapSpider` 使用相同的 `rules()` API：抓取 `sitemap_urls`、遞迴展開 sitemap index，並將每個 URL 交給你的 rules 分派。直接把 `robots.txt` 的 URL 放進 `sitemap_urls`，spider 會自動抽出其中每個 `Sitemap:` 指令。

`LinkExtractor` 的 allow / deny / restrict_css / canonicalize 完整參數見[上游 spiders 文件](https://github.com/D4Vinci/Scrapling/tree/main/docs)。
