---
name: apple-notes
description: "基於終端機與 Python 的蘋果備忘錄管理工具 (A command-line and Python tool for Apple Notes)，支援透過指令操作 MacOS Notes.app。"
version: 1.0.0
homepage: https://github.com/RhetTbull/macnotesapp
author: Hermes Agent
license: MIT
platforms: [macos]
metadata:
    hermes:
        tags: [Notes, tasks, todo, macOS, Apple]
prerequisites:
    commands: [macnotesapp]
---

# 蘋果備忘錄 (apple-notes / macnotesapp)

使用 `notes` 或 `macnotesapp` 指令可以在終端機 (terminal) 直接與 Apple MacOS Notes.app 進行互動。此工具也提供了 Python API，方便在腳本中調用。

## 初始設定 (Setup)

推薦的安裝方式是使用 `uv`，或者在 Apple Silicon (M1 等) 裝置上使用 Homebrew：

- **透過 uv 安裝 (Install via uv)**:

    ```bash
    uv tool install --python 3.13 macnotesapp
    ```

- **免安裝直接執行 (Run without installing)**:

    ```bash
    uvx --python 3.13 macnotesapp
    ```

- **透過 Homebrew 安裝 (Install via Homebrew)** (僅支援 Apple Silicon):

    ```bash
    brew tap RhetTbull/macnotesapp
    brew install macnotesapp
    ```

## 核心指令與操作 (Core Commands)

主要的指令為 `notes`，以下為支援的操作。

> 預設使用 Markdown (Use Markdown by default)：新增或編輯筆記時，預設帶上 `-m, --markdown`，內文以 Markdown 撰寫 (例如 `notes add -m`、`notes edit -m`)；檢視筆記時亦可加上 `-m` 以 Markdown 輸出。

### 全域選項 (Global Options)

- `-v, --version`：顯示版本資訊。
- `-h, --help`：顯示說明訊息 (各子指令亦可加上 `--help` 取得專屬說明)。

### 查詢與檢視 (List and View)

- `notes list TEXT`：列出筆記 (支援依據帳號或文字內容過濾)。
    - `-a, --account ACCOUNT`：限定特定帳號；可重複使用以涵蓋多個帳號。
- `notes cat NOTE_NAME`：將一個或多個筆記印到標準輸出 (STDOUT)。
    - `-p, --plaintext`：以純文字輸出筆記。
    - `-m, --markdown`：以 Markdown 輸出筆記。
    - `-h, --html`：以 HTML 輸出筆記。
    - `-j, --json`：以 JSON 輸出筆記 (內文預設為 HTML；若同時指定 `--plaintext` 或 `--markdown`，則內文改用該格式)。
- `notes accounts`：印出備忘錄的帳號資訊 (Notes accounts)。
    - `-j, --json`：以 JSON 格式輸出。
- `notes dump`：傾印 (dump) 所有的筆記或特定筆記以供除錯。
    - `-s, --selected`：僅傾印目前選取的筆記。
    - `-B, --no-body`：不傾印筆記內文 (body)。

### 新增與編輯筆記 (Add and Edit Notes)

- `notes add NOTE`：新增一則筆記。如果有多行，第一行將作為標題，其餘作為內文。
    - `-s, --show`：新增後在 Notes.app 中顯示該筆記。
    - `-F, --file FILENAME`：從檔案新增筆記。
    - `-u, --url URL`：從網址下載內容並轉換成易讀的版本存入筆記。
    - `-h, --html`：內文以 HTML 格式處理。
    - `-m, --markdown`：內文以 Markdown 格式處理。
    - `-p, --plaintext`：內文以純文字處理 (預設值，可透過 `notes config` 變更)。
    - `-e, --edit`：新增前先以預設編輯器編輯筆記內容。
    - `-a, --account ACCOUNT`：將筆記新增到指定帳號。
    - `-f, --folder FOLDER`：將筆記新增到指定資料夾。
    - `cat file.txt | notes add`：亦支援從標準輸入 (STDIN) 新增筆記。
- `notes edit NOTE_NAME`：編輯現有筆記的內文 (body)。
    - `-b, --body TEXT`：直接設定內文，不開啟編輯器。
    - `-h, --html`：將內文視為 HTML。
    - `-m, --markdown`：將內文視為 Markdown。
    - `-a, --account ACCOUNT`：指定搜尋的帳號。
- `notes rename OLD_NAME NEW_NAME`：重新命名筆記標題。
    - `-a, --account ACCOUNT`：指定搜尋的帳號。

### 刪除與移動筆記 (Delete and Move)

- `notes delete NOTE_NAME`：刪除特定筆記。
    - `-y, --yes`：略過確認提示。
    - `-a, --account ACCOUNT`：指定搜尋的帳號。
- `notes move NOTE_NAME`：將筆記移動至不同的資料夾。
    - `-f, --folder TEXT`：目標資料夾 (必填)。
    - `-a, --account ACCOUNT`：指定搜尋的帳號。

### 資料夾管理 (Folder Management)

- `notes mkdir FOLDER_NAME`：建立新的資料夾。
    - `-a, --account ACCOUNT`：指定要建立資料夾的帳號。
- `notes rmdir FOLDER_NAME`：刪除指定的資料夾。
    - `-y, --yes`：略過確認提示。
    - `-a, --account ACCOUNT`：指定要刪除資料夾的帳號。

### 其他配置 (Config)

- `notes config`：設定預設的帳號、編輯器等配置。
- `notes help <command>`：印出指定指令的說明。

## 已知限制 (Known Limitations)

- 不支援被密碼鎖定的筆記 (Unlocked 狀態的密碼筆記可讀取，但 Locked 狀態則無法)。
- 無法寫入標籤 (Tags)，手動加入的 `#tagname` 會被視為純文字。讀取時，內文的標籤會被過濾掉。
- 目前不支援處理附件 (Attachments) (相關內容會被忽略)。
- 目前僅能存取最上層資料夾 (top-level folders) 內的筆記。

## Python API 整合 (Python Usage)

此工具提供了一套 Python 函式庫，讓你可以編寫程式來控制 Notes.app：

```python
from macnotesapp import NotesApp

notesapp = NotesApp()
# 取得所有的筆記
notes = notesapp.notes()

# 新增筆記到預設帳號與資料夾
new_note = notesapp.make_note(
    name="New Note",
    body="This is a new note created with #macnotesapp"
)
```
