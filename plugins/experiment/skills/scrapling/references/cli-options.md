# Scrapling CLI — 完整選項與範例 (Full CLI Options)

`scrapling extract` 指令群的全部選項。SKILL.md 只保留指令選擇規則，細節在此。

```bash
Usage: scrapling extract [OPTIONS] COMMAND [ARGS]...

Commands:
  get             Perform a GET request and save the content to a file.
  post            Perform a POST request and save the content to a file.
  put             Perform a PUT request and save the content to a file.
  delete          Perform a DELETE request and save the content to a file.
  fetch           Use a browser to fetch content with browser automation and flexible options.
  stealthy-fetch  Use a stealthy browser to fetch content with advanced stealth features.
```

輸出格式由副檔名決定：`.md` 轉 Markdown（文件閱讀最佳）、`.html` 保留原始 HTML、`.txt` 存純文字。

## Key options (requests)

以下選項為 4 個 HTTP request 指令 (`get` / `post` / `put` / `delete`) 共用：

| Option                                     | Input type | Description                                                                                                                                    |
|:-------------------------------------------|:----------:|:-----------------------------------------------------------------------------------------------------------------------------------------------|
| -H, --headers                              |    TEXT    | HTTP headers in format "Key: Value" (can be used multiple times)                                                                               |
| --cookies                                  |    TEXT    | Cookies string in format "name1=value1; name2=value2"                                                                                          |
| --timeout                                  |  INTEGER   | Request timeout in seconds (default: 30)                                                                                                       |
| --proxy                                    |    TEXT    | Proxy URL in format "http://username:password@host:port"                                                                                       |
| -s, --css-selector                         |    TEXT    | CSS selector to extract specific content from the page. It returns all matches.                                                                |
| -p, --params                               |    TEXT    | Query parameters in format "key=value" (can be used multiple times)                                                                            |
| --follow-redirects / --no-follow-redirects |    None    | Whether to follow redirects (default: "safe", rejects redirects to internal/private IPs)                                                       |
| --verify / --no-verify                     |    None    | Whether to verify SSL certificates (default: True)                                                                                             |
| --impersonate                              |    TEXT    | Browser to impersonate. Can be a single browser (e.g., Chrome) or a comma-separated list for random selection (e.g., Chrome, Firefox, Safari). |
| --stealthy-headers / --no-stealthy-headers |    None    | Use stealthy browser headers (default: True)                                                                                                   |
| --ai-targeted                              |    None    | Extract only main content and sanitize hidden elements for AI consumption (default: False)                                                     |

`post` 與 `put` 專屬：

| Option     | Input type | Description                                                                             |
|:-----------|:----------:|:----------------------------------------------------------------------------------------|
| -d, --data |    TEXT    | Form data to include in the request body (as string, ex: "param1=value1&param2=value2") |
| -j, --json |    TEXT    | JSON data to include in the request body (as string)                                    |

範例：

```bash
# Basic download
scrapling extract get "https://news.site.com" news.md

# Download with custom timeout
scrapling extract get "https://example.com" content.txt --timeout 60

# Extract only specific content using CSS selectors
scrapling extract get "https://blog.example.com" articles.md --css-selector "article"

# Send a request with cookies
scrapling extract get "https://scrapling.requestcatcher.com" content.md --cookies "session=abc123; user=john"

# Add user agent
scrapling extract get "https://api.site.com" data.json -H "User-Agent: MyBot 1.0"

# Add multiple headers
scrapling extract get "https://site.com" page.html -H "Accept: text/html" -H "Accept-Language: en-US"
```

## Key options (browsers)

`fetch` 與 `stealthy-fetch` 共用：

| Option                                   | Input type | Description                                                                                                                                              |
|:-----------------------------------------|:----------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------|
| --headless / --no-headless               |    None    | Run browser in headless mode (default: True)                                                                                                             |
| --disable-resources / --enable-resources |    None    | Drop unnecessary resources for speed boost (default: False)                                                                                              |
| --network-idle / --no-network-idle       |    None    | Wait for network idle (default: False)                                                                                                                   |
| --real-chrome / --no-real-chrome         |    None    | If you have a Chrome browser installed on your device, enable this, and the Fetcher will launch an instance of your browser and use it. (default: False) |
| --timeout                                |  INTEGER   | Timeout in milliseconds (default: 30000)                                                                                                                 |
| --wait                                   |  INTEGER   | Additional wait time in milliseconds after page load (default: 0)                                                                                        |
| -s, --css-selector                       |    TEXT    | CSS selector to extract specific content from the page. It returns all matches.                                                                          |
| --wait-selector                          |    TEXT    | CSS selector to wait for before proceeding                                                                                                               |
| --proxy                                  |    TEXT    | Proxy URL in format "http://username:password@host:port"                                                                                                 |
| -H, --extra-headers                      |    TEXT    | Extra headers in format "Key: Value" (can be used multiple times)                                                                                        |
| --dns-over-https / --no-dns-over-https   |    None    | Route DNS through Cloudflare's DoH to prevent DNS leaks when using proxies (default: False)                                                              |
| --block-ads / --no-block-ads             |    None    | Block requests to ~3,500 known ad and tracker domains (default: False)                                                                                   |
| --ai-targeted                            |    None    | Extract only main content and sanitize hidden elements for AI consumption (default: False). Also enables ad blocking automatically.                      |

`fetch` 專屬：

| Option   | Input type | Description                                                 |
|:---------|:----------:|:------------------------------------------------------------|
| --locale |    TEXT    | Specify user locale. Defaults to the system default locale. |

`stealthy-fetch` 專屬：

| Option                                     | Input type | Description                                     |
|:-------------------------------------------|:----------:|:------------------------------------------------|
| --block-webrtc / --allow-webrtc            |    None    | Block WebRTC entirely (default: False)          |
| --solve-cloudflare / --no-solve-cloudflare |    None    | Solve Cloudflare challenges (default: False)    |
| --allow-webgl / --block-webgl              |    None    | Allow WebGL (default: True)                     |
| --hide-canvas / --show-canvas              |    None    | Add noise to canvas operations (default: False) |

範例：

```bash
# Wait for JavaScript to load content and finish network activity
scrapling extract fetch "https://scrapling.requestcatcher.com/" content.md --network-idle

# Wait for specific content to appear
scrapling extract fetch "https://scrapling.requestcatcher.com/" data.txt --wait-selector ".content-loaded"

# Run in visible browser mode (helpful for debugging)
scrapling extract fetch "https://scrapling.requestcatcher.com/" page.html --no-headless --disable-resources

# Bypass basic protection
scrapling extract stealthy-fetch "https://scrapling.requestcatcher.com" content.md

# Solve Cloudflare challenges
scrapling extract stealthy-fetch "https://nopecha.com/demo/cloudflare" data.txt --solve-cloudflare --css-selector "#padded_content a"

# Use a proxy for anonymity.
scrapling extract stealthy-fetch "https://site.com" content.md --proxy "http://proxy-server:8080"
```

## Docker

使用者沒有 Python 或不想裝 Python 時，可改用 Docker image。此路徑`只能`跑 CLI 指令，不能寫 Python：

```bash
docker pull pyd4vinci/scrapling
# 或
docker pull ghcr.io/d4vinci/scrapling:latest
```
