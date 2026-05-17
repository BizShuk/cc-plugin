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

主要的指令為 `notes`，以下為支援的操作：

### 查詢與檢視 (List and View)

- `notes list`：列出筆記 (支援依據帳號或文字內容過濾)。
- `notes cat`：將一個或多個筆記印到標準輸出 (STDOUT)。
- `notes accounts`：印出備忘錄的帳號資訊 (Notes accounts)。
- `notes dump`：傾印 (dump) 所有的筆記或特定筆記以供除錯。

### 新增與編輯筆記 (Add and Edit Notes)

- `notes add "筆記內容"`：新增一則筆記。如果有多行，第一行將作為標題。
- `notes add --edit` (或 `-e`)：開啟預設編輯器撰寫新筆記。
- `cat file.txt | notes add`：從標準輸入 (STDIN) 新增筆記。
- `notes add --url URL`：從網址下載內容並轉換成易讀的版本存入筆記。
- `notes edit`：編輯現有筆記的內文 (body)。
- `notes rename`：重新命名筆記標題。

### 刪除與移動筆記 (Delete and Move)

- `notes delete`：刪除特定筆記。
- `notes move`：將筆記移動至不同的資料夾。

### 資料夾管理 (Folder Management)

- `notes mkdir`：建立新的資料夾。
- `notes rmdir`：刪除指定的資料夾。

### 其他配置 (Config)

- `notes config`：設定預設的帳號、編輯器等配置。

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
