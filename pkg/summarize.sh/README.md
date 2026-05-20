<https://summarize.sh/docs/config.html>

```bash

mdkir -p ~/.summarize/
cp config.sample.json ~/.summarize/config.json
ln -s ~/.summarize ./config/

npx @steipete/summarize
# read ./.env, ~/.summarize/config.json



# 1. yt-dlp：影片下載神器
# 主要功能： 專門用來下載網路影片與音訊。它支援 YouTube、Vimeo、Bilibili 等上千個網站。
# 生活化比喻： 就像一個萬能的「影片抓取器」，只要給它網址，它就能把最高畫質的影片或純音檔拔下來。
# yt-dlp 下載與轉檔操作. 來源： OSTechNix
brew install ffmpeg
# 2. ffmpeg：多媒體處理萬能瑞士刀
# 主要功能： 負責影片和音訊的轉檔、剪輯、壓縮與解碼（Codec）。
# 重要關聯： yt-dlp 下載影片時，必須依賴 ffmpeg 在幕後把下載下來的影像流與聲音流「縫合」成一個完整的 MP4 或 MP3 檔案。
# FFmpeg 幕後影音處理運作原理. 來源： Medium
brew install yt-dlp
# 3. whisper-cpp：AI 語音轉文字工具
# 主要功能： 這是由 OpenAI 開源的 Whisper 語音辨識模型（Speech-to-Text），經過 C/C++ 最佳化後的版本。它可以在你的電腦上完全離線運行。
# 生活化比喻： 就像一個「聽寫祕書」，你把剛剛下載的音檔丟給它，它就能自動聽懂裡面的人在說什麼，並精準輸出成逐字稿或帶有時間軸的 SRT 字幕檔。
brew install whisper-cpp
# 4. tesseract（選填）：圖片文字辨識（OCR）
# 主要功能： 光學字元辨識（Optical Character Recognition, OCR）。能把圖片或影片截圖中的「死文字」辨識並提取出來。
# 在此處的用途： 你的備註寫著 --slides-ocr，這代表當你下載某些線上課程、演講或簡報（Slides）影片時，這個工具可以自動幫你「看」畫面上投影片寫了什麼字，並把它轉成讀得懂的文本。
brew install tesseract
```
