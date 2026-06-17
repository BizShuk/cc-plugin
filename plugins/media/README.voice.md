/
I tried to use the AI voice model. Could you let me know how to prong or configuration in order to get the proper voices. How many of the ways or aspect can help me to do specialization for the voice?

8:23 PM

Architected comprehensive framework for voice model specialization dimensions
Architected comprehensive framework for voice model specialization dimensions
先用一個生活化的比喻把整件事框起來:設定 AI 語音模型,就像在「指導一位配音員 (voice actor) 錄音」。你要決定三件事——這位配音員「天生是誰」、「你要他怎麼演」、以及「錄音棚的技術規格」。把這三層拆開,你就會發現專屬化 (specialization) 的著力點其實非常多。

下面我分成「控制方式」和「可調面向」兩部分,最後給你一個實作計畫。

一、先搞懂:控制語音有四種「介面方式」
不同模型 (ElevenLabs、OpenAI TTS、Google、MiniMax、Azure⋯⋯) 吃的「指令格式」不一樣,先認得它們長什麼樣:

文字描述式 (prompt-based / descriptive)
你直接用一段話描述你要的聲音,例如「一位溫暖、四十歲左右的男聲,語速偏慢、帶點疲憊感」。新一代模型 (像 OpenAI 的 gpt-4o-mini-tts) 支援這種「用講的」方式。
比喻:像跟配音員說「我要那種深夜電台 DJ 的感覺」。
參數式 (parameter-based)
用一格一格的滑桿/數值控制,例如 stability(穩定度)、similarity(相似度)、speed(語速)、style exaggeration(風格誇張度)。ElevenLabs 是典型。
比喻:像錄音師調混音台上的旋鈕。
參考音訊式 (reference audio / voice cloning)
你丟一段聲音樣本進去,模型去模仿那個音色。
比喻:像跟配音員說「請學這個人講話」。
標記語言式 (markup / SSML)
在文字裡插入標籤來控制細節,例如 <break time="500ms"/> 控制停頓、<emphasis> 控制重音、<phoneme> 控制發音。Azure、Google 大量使用。
比喻:像在劇本上用筆標註「這裡停一下」「這個字要重唸」。
多數成熟的平台會「混用」這幾種——先選音色 → 用參數微調 → 再用 SSML 精修。

二、可以做專屬化的「面向」:大約 12 個
我把它們分三層,對應前面那個比喻。

A. 身份層(這位配音員「是誰」)

音色 / 聲音身份 (timbre / voice identity):選預設聲音,或用 voice cloning 複製特定人聲。這是最根本的一層。
語言與口音 (language & accent):同一個音色可能有美式、英式;中文有台灣腔、大陸腔的差別。
性別與年齡感 (gender & age):有些模型可直接指定。
B. 表現層(你要他「怎麼演」)——這層最影響「自然度」4. 情緒 / 語氣 (emotion / tone):開心、悲傷、興奮、冷靜、嚴肅。5. 語速 (speed / rate):快慢。6. 音高 (pitch):高低。7. 音量 (volume):整體響度。8. 停頓與節奏 (pauses & pacing):用 SSML 的 <break> 製造呼吸感。9. 重音 / 強調 (emphasis / stress):把關鍵字唸重。10. 說話風格 (style):旁白 (narration)、對話 (conversational)、播報 (newscasting) 是不同模式。
比喻:同一位演員,演旁白和演吵架的語氣完全不同。

C. 技術層(錄音棚的規格) 11. 穩定度 vs 多樣性 (stability vs variability):調高穩定度 → 每次唸都很一致但較單調;調低 → 更有情感起伏但偶爾「演過頭」。12. 發音修正 (pronunciation):用 <phoneme> 或自訂字典,處理人名、品牌、縮寫、專有名詞。13. 取樣率 / 音檔格式 (sample rate / format)、種子 (seed):seed 固定後可重現同一個輸出,方便除錯。

三、給你的實作計畫(由淺入深)
Step 1 — 先確認你用哪個模型
這是最關鍵的一步,因為「prompt 怎麼寫」完全取決於平台。描述式模型要寫自然語言、參數式模型要調滑桿、SSML 模型要寫標籤——搞錯方向會白忙。

Step 2 — 鎖定音色,先求穩
選一個接近的預設聲音(或 clone 一段乾淨的參考音訊,建議 30 秒以上、無背景雜音)。先不要動太多參數,聽聽基準長怎樣。

Step 3 — 一次只改一個變數
像實驗一樣:固定其他條件,只調語速,聽差別;再只調情緒。一次改太多會分不清是哪個參數造成的效果。

Step 4 — 用 SSML 做最後精修
基礎對了之後,再針對個別句子加停頓、重音、修發音。

Step 5 — 固定 seed + 存設定
找到滿意組合後記錄下所有參數(含 seed),這樣可以重現,也方便日後微調。

要我給你具體可貼上去用的範例(prompt 寫法或 SSML 範本),我需要知道你用的是哪一個語音模型/平台。你是用 OpenAI、ElevenLabs、Google/Azure,還是其他的?告訴我我就能直接針對它寫出可用的設定。

Compare current voices models on the market including VOXCPM

8:29 PM

Thought for 4s
Thought for 4s
I'll research the current landscape, including VoxCPM, before comparing. Let me search.

Synthesizing open-weights tier specifics alongside VoxCPM comparisons
Synthesizing open-weights tier specifics alongside VoxCPM comparisons
Good detail on VoxCPM. Now let me check the commercial and other open-source leaders.

Synthesized comprehensive voice model landscape with structured comparison framework
Synthesized comprehensive voice model landscape with structured comparison framework
先講結論:2026 年的語音模型市場已經沒有「單一最強」這回事了。語音生成已不再是「一家獨大」的問題——ElevenLabs、OpenAI TTS、Hume、Cartesia、PlayHT、Gemini TTS,以及開源權重 (open-weights) 陣營各自在不同場景勝出。選擇取決於你要的是「音質」「速度」「情緒」「便宜」還是「能自己部署」。
SurePrompts

下面我用一個生活化的比喻幫你定位:挑語音模型就像挑餐廳——米其林星級(音質最好但貴)、快餐連鎖(快又穩)、自家廚房(開源、自己煮、免費但要會用)。

一、按「強項」分群
音質與聲音克隆 (voice quality & cloning) 王者
ElevenLabs 仍是成熟的全能選手,介面精緻、有快速的 Flash v2.5 模型。不過它已不是唯一頂尖選項:在 Artificial Analysis 語音競技場 (Speech Arena) 的 ELO 排名上,Inworld Realtime TTS 1.5 Max 與 Google Gemini 3.1 Flash TTS 在更低價格下,音質超越了 ElevenLabs 的 Eleven v3。另外 Fish Audio 在 TTS-Arena2 排名第一,支援 80+ 語言與 50+ 種情緒控制 (emotion controls)。
Befreed + 2

可用文字「指揮」聲音 (instructable / steerable)
這是 OpenAI 的強項。除了選 alloy、echo、nova 等預設聲音外,你可以在請求裡直接下指令控制角色、語氣與表達,例如「像冷靜的博物館導覽員那樣講」「聽起來有點惱怒」「用快速新聞快報的方式唸」——這種自然語言可控性是 2026 年所有閉源供應商中最強的。
SurePrompts

最便宜
OpenAI 的 tts-1 與 tts-1-hd 都是每百萬字元 15 美元 (US$15/M chars),六個內建聲音、支援 57+ 種語言;但缺點是沒有 SSML、沒有聲音克隆、沒有逐字時間戳,也不能調整語速。相對地 ElevenLabs 走訂閱制 (subscription),約 US$0.30/1K 字元。
TokenMix
TokenMix

低延遲 / 即時語音代理 (real-time voice agent)
Cartesia 把首位元組時間 (time-to-first-byte) 壓到 100 毫秒以下,專為即時應用設計。這類場景(客服、語音助理)講求的是「秒回」,音質可以稍微讓步。
SurePrompts

情緒表現 (emotion):Hume 是這一塊的代表。

企業級雲端:Google Cloud TTS、Azure、Amazon Polly,穩定、可規模化、有免費額度。

二、VoxCPM 專門拆解(你特別問到的)
VoxCPM 屬於「自家廚房」這一派——開源、可自部署、免費,而且是這個陣營裡很特別的一個。

它的核心創新是 tokenizer-free(免分詞器) 架構。一般 TTS 會先把語音切成離散 token(像把音樂剁成 MIDI 音符再重組),VoxCPM 則直接在連續語音空間建模,類似在原始連續空間裡作曲,而不是切成 MIDI 再重建,因此保留了語音的自然流暢與表現力。架構上 它建立在 MiniCPM-4 骨幹上,透過階層式語言建模 (hierarchical language modeling) 與 FSQ 約束達成語意-聲學的隱性解耦 (semantic-acoustic decoupling),同時提升表現力與穩定性。
Voxcpm
huggingface

效能與規格:

VoxCPM-0.5B 用 180 萬小時雙語語料訓練,在開源系統中達到 state-of-the-art 的 zero-shot TTS 表現。
arXiv
在 SEED-TTS-EVAL 上,英文 WER 1.85%、中文 CER 0.93%,超越 IndexTTS2 與 CosyVoice2;說話者相似度 (speaker similarity) 英文 72.9%、中文 77.2%。這對你(台灣、中英混用情境)特別有意義——它的中英文本來就是強項。
arXiv
支援串流合成 (streaming synthesis),在消費級 RTX 4090 上即時因子 (Real-Time Factor, RTF) 低至 0.17,可做即時應用。(RTF 0.17 = 生成 1 秒語音只需 0.17 秒運算,綽綽有餘)
Openbmb
採 Apache 2.0 授權,程式碼與權重全開,可從 GitHub / HuggingFace 取得自行部署。
Voxcpm
最新版本:VoxCPM2 是最新主要版本——2B 參數、用超過 200 萬小時多語語料訓練,現支援 30 種語言、Voice Design(聲音設計)、可控聲音克隆 (Controllable Voice Cloning),並輸出 48kHz 錄音室等級音質。
GitHub

一個值得注意的特性:它會根據輸入文字的語意「自動推斷」適當的語氣與風格,而不需要你下指令。這是雙面刃——省事,但比起 OpenAI 那種「用指令精準指揮」的可控性低。
Openbmb

風險提醒:VoxCPM 強大的 zero-shot 聲音克隆能力可能被濫用於製作仿冒、詐騙或散播假訊息的深偽 (deepfake) 語音。自部署時要自己負責使用倫理。
Hugging Face

順帶一提,開源陣營還有對手:Mistral AI 於 2026 年 3 月 26 日發表 Voxtral TTS,號稱首個前沿品質的開放權重 TTS,4B 參數、16GB 以上 VRAM 單卡可跑,內建 20 個預設聲音加 zero-shot 克隆。其他常見的還有 Fish Audio(也可自架)、Kokoro、XTTS-v2、CosyVoice、IndexTTS。
Digital Applied Team

三、橫向比較表
模型 類型 中英文 最大強項 大致定價 適合誰
VoxCPM / VoxCPM2 開源 極強 自然度、可自部署、免費、中文佳 免費(自負算力) 想自己掌控、注重隱私、中英內容
ElevenLabs 閉源 好 音質、聲音克隆、生態成熟 ~$0.30/1K 字元 有聲書、長內容、品質至上
OpenAI TTS 閉源 好 文字指令可控、便宜、易整合 $15/M 字元 已用 GPT 生態、求簡單
Fish Audio 閉源/可自架 強 情緒控制多、多語、價優 約 $15/M 多語、要細膩情緒
Cartesia Sonic 閉源 好 超低延遲 API 計費 即時語音代理
Voxtral (Mistral) 開源 好 開放權重、品質高 免費(自負算力) 想自部署但要高品質
Google / Azure / Polly 閉源雲 好 穩定、規模化、免費額度 多階梯 企業級大量部署
四、給你的選擇計畫
考量你在台灣、中英混用、又有自架基礎建設的能力,我的建議路線:

先定義你的「主場景」:有聲書/旁白 → 重音質(ElevenLabs / VoxCPM2);即時對話 → 重延遲(Cartesia);要省錢且能寫程式控制 → OpenAI;要隱私自主 → VoxCPM2 / Voxtral。
既然你能自部署,優先試 VoxCPM2:中英文本就是它的強項,Apache 2.0 無授權包袱,RTX 4090 等級就能即時跑。先用預設設定產一段基準音,再調克隆參考音訊。
拿 1~2 段你的真實腳本當「對照測試 (A/B test)」:同一段文字分別丟進 VoxCPM2、ElevenLabs、OpenAI,盲聽比較自然度與中文發音。市場排名是參考,你的耳朵和你的內容才是準。
算總成本:閉源是「按量付費、零維運」;開源是「免費授權、但要付 GPU 與維運時間」。用量大時開源划算,用量小或不想顧機器時閉源反而便宜。
想要的話,我可以給你一份 VoxCPM2 的本機部署 + 中英克隆的實作步驟,或一份「同一腳本跨模型盲測」的測試腳本範本。你想先往哪個方向走?

給份 voxcpm 部署指南 然後告訴我他需要的資源硬體 Mac mini是否可以部署？

8:35 PM

驗證VoxCPM硬體需求與Mac相容性以提供部署指南。
驗證VoxCPM硬體需求與Mac相容性以提供部署指南。
要給你準確的部署步驟和硬體需求,我先去抓官方倉庫的實際資訊,特別是 Mac(Apple Silicon)的支援情況。

整合部署資訊、評估Mac相容性、規劃分層指南。
整合部署資訊、評估Mac相容性、規劃分層指南。
抓到官方倉庫的最新資料了。先直接回答你最關心的問題:

Mac mini 能不能部署?能,但要看是哪一台
Apple Silicon(M 系列)Mac mini → 可以。 官方的 Web demo 明確支援 Mac:app.py 的 --device 參數可選 auto、cpu、mps、cuda,在 Apple Silicon Mac 上,auto 會自動使用 MPS(Apple 的 GPU 加速後端)。也就是說它會吃你 M 晶片的 GPU。
GitHub

Intel Mac mini → 技術上可以,但只能跑 CPU,會很慢,實用性低,建議走下面的 VoxCPM.cpp 路線。

關鍵不是「能不能跑」,而是記憶體夠不夠。Apple Silicon 是「統一記憶體 (unified memory)」——CPU 和 GPU 共用同一池 RAM,就像一間「廚房和餐廳共用的食材櫃」。模型要佔走一大塊,系統還要留一塊,所以 RAM 容量是真正的門檻。

硬體 / 資源需求(依版本)
VoxCPM 有三個版本,需求差很多。下表的 VRAM 在 Mac 上就等於要佔用的統一記憶體:

版本 參數量 記憶體需求 音質 語言 建議 Mac mini 配置
VoxCPM2(最新) 2B ~8 GB
GitHub
48kHz 錄音室級
GitHub
30 種
GitHub
16GB 以上(24GB 更穩)
VoxCPM1.5(穩定) 0.6B ~6 GB
GitHub
44.1kHz 中、英 16GB
VoxCPM-0.5B(舊版) 0.5B ~5 GB
GitHub
16kHz 中、英 8~16GB
白話建議:

8GB Mac mini:很吃緊。模型加系統幾乎把記憶體吃滿,容易卡頓或被 swap 拖慢。勉強跑 0.5B 版。
16GB Mac mini:VoxCPM2 的甜蜜點,可以順跑。
24GB 以上:最舒服,還能同時做別的事。
速度方面要有心理準備:VoxCPM2 在 NVIDIA RTX 4090 上的即時因子 (RTF) 約 0.3。Mac 的 MPS 會比 4090 慢不少,所以比較適合「批次離線生成」(batch),不一定能做到電話客服那種即時對話。如果你要追求 Mac 上的速度,官方生態裡有專門優化:VoxCPMANE 提供 Apple Neural Engine 後端、VoxCPM.cpp 提供 GGML/GGUF 的 CPU/Vulkan 推論、VoxCPM-ONNX 提供 ONNX 的 CPU 推論。
GitHub
GitHub

部署指南(Apple Silicon Mac mini,最簡路線)
前置需求:Python ≥ 3.10(且 < 3.13)、PyTorch ≥ 2.5.0。(官方寫的 CUDA ≥ 12.0 是給 NVIDIA 顯卡用的,Mac 走 MPS 不需要 CUDA。)
GitHub

Step 1 — 裝環境
建議用 conda 或 venv 開一個乾淨環境,避免污染系統 Python(像「給這個專案一個獨立的工作檯」):

bash
conda create -n voxcpm python=3.11 -y
conda activate voxcpm
pip install voxcpm
Step 2 — 寫一支最小測試腳本
存成 test.py。注意我加了 device="mps" 讓它走 Mac GPU:

python
from voxcpm import VoxCPM
import soundfile as sf

model = VoxCPM.from_pretrained(
"openbmb/VoxCPM2",
load_denoiser=False,
)

wav = model.generate(
text="哈囉,這是 VoxCPM2 在我的 Mac mini 上跑出來的聲音。",
cfg_value=2.0,
inference_timesteps=10,
)
sf.write("demo.wav", wav, model.tts_model.sample_rate)
print("完成:demo.wav")
第一次執行會自動從 HuggingFace 下載權重(2B 模型檔不小,要等一下、確保網路順)。

Step 3 — 試「聲音設計」(Voice Design,免參考音訊)
這是 VoxCPM2 的亮點:用一句話描述就生出新聲音。格式是把描述放在文字最前面的括號裡,例如 (描述)要合成的文字:
GitHub

python
wav = model.generate(
text="(一位年輕女性,溫柔甜美的聲線)歡迎使用 VoxCPM2!",
cfg_value=2.0,
inference_timesteps=10,
)
sf.write("design.wav", wav, model.tts_model.sample_rate)
Step 4 — 試「聲音克隆」(Voice Cloning)
丟一段參考音訊,它複製音色,還能加風格指令:

python
wav = model.generate(
text="(語速稍快、開朗的語氣)這是帶風格控制的克隆聲音。",
reference_wav_path="path/to/voice.wav",
cfg_value=2.0,
inference_timesteps=10,
)
sf.write("clone.wav", wav, model.tts_model.sample_rate)
Step 5 — 開圖形介面玩(可選)
不想寫程式可以開官方 Web UI:

bash
python app.py --device auto # Apple Silicon 會自動用 MPS

# 然後瀏覽器開 <http://localhost:8808>

給你的行動計畫
先確認你的 Mac mini 規格:是不是 M 系列(Apple Silicon)?統一記憶體幾 GB?這決定一切。
記憶體 < 16GB → 先跑 VoxCPM-0.5B 版(把上面 "openbmb/VoxCPM2" 換成 "openbmb/VoxCPM-0.5B"),確認流程通了再考慮升級。
16GB 以上 → 直接上 VoxCPM2,享受 48kHz 與 30 語言。
若覺得 MPS 太慢:改用 VoxCPMANE(Apple Neural Engine 優化)或 VoxCPM.cpp(GGUF 量化版,記憶體佔用更低),這兩個是專為 Apple 硬體調過的。
想要專屬聲音:之後可以做微調——只要 5~10 分鐘的音訊,用 LoRA 就能適配特定說話者或領域。(不過微調建議在有 NVIDIA GPU 的機器或雲端跑,Mac 訓練會吃力。)
GitHub
一個務實提醒:克隆/聲音設計有隨機性,官方建議同一段可生成 1~3 次來挑出滿意的聲音或風格。
GitHub

要不要我幫你把這份指南整理成一個可以存下來照做的 Markdown 檔?或者你先告訴我你的 Mac mini 是哪一代、幾 GB 記憶體,我幫你判斷該選哪個版本、跑起來大概多快。

Create the markdown step by step guideline for me

8:43 PM

Crafting markdown guide directly in outputs directory
Crafting markdown guide directly in outputs directory
I'll create a clean Markdown guide you can save and follow. Let me check for a relevant skill first.

Orchestrated comprehensive Traditional Chinese markdown guide creation
Orchestrated comprehensive Traditional Chinese markdown guide creation
No markdown skill needed — it's plain text. I'll write the guide now in Traditional Chinese.

Created a file, read a file
Created a file, read a file
指南做好了,可以下載儲存。內容涵蓋十個段落:版本選擇、硬體需求、安裝、第一次測試、Voice Design、Voice Cloning、圖形介面、Mac 加速方案、微調,以及一份從頭到尾的執行順序。

兩個重點再幫你濃縮:Apple Silicon Mac mini 走 MPS 可以部署,16GB 統一記憶體是 VoxCPM2 的甜蜜點;若記憶體吃緊就先用 0.5B 版,或改用為 Apple 硬體優化的 VoxCPMANE / VoxCPM.cpp。

如果你之後告訴我 Mac mini 的實際型號和記憶體,我可以幫你把指南裡的版本選擇和預期速度標註得更精準。

Voxcpm 部署指南
Document · MD
幫我建立10種聲音樣板

8:47 PM

設計十種聲音角色並規劃實踐方案。
設計十種聲音角色並規劃實踐方案。
幫你設計 10 種涵蓋常見用途的聲音樣板,每一種都用 VoxCPM 的 Voice Design 格式(把描述放進括號),並附建議參數和一支「一次生成全部 10 種」的 Python 批次腳本。我直接整理成檔案,方便你跟前面那份部署指南一起存。

Created a file, read a file
Created a file, read a file
10 種樣板做好了。涵蓋的用途分別是:溫暖旁白、活潑客服、新聞主播、療癒助眠、Podcast 主持、童趣故事、廣告促銷、長者智者、企業簡報、電影預告。

每一種都附了三樣東西:可直接貼上的 Voice Design 描述字串、建議的 cfg_value 與 inference_timesteps 參數,以及一個生活化比喻幫你抓感覺。檔案最後還有一支 make_voices.py 批次腳本,跑一次就會在 voice_samples/ 產出全部 10 個試聽檔。

兩個使用重點:描述越具體越好(性別、年齡、音色、語速、情緒都能寫進去),而且 Voice Design 有隨機性,同一段多生幾次挑最好的是正常做法。

聽完後如果想把某個樣板再調得更貼近你要的(例如指定台灣口音、特定年齡感),把你選中的那個告訴我,我幫你把描述詞和參數再精修一版。

Voxcpm 10種聲音樣板
Document · MD
如果我要模仿聲音的話要怎麼prompt

8:52 PM

釐清聲音克隆需參考音訊，規劃提示寫法與倫理考量。
釐清聲音克隆需參考音訊，規劃提示寫法與倫理考量。
關鍵觀念要先講清楚:「模仿一個真實的聲音」和「用文字描述造聲音」是兩件事。

前面那 10 種樣板用的是 Voice Design——純文字描述,等於「無中生有」造一個聲音。但你說的「模仿聲音」是要複製某個已經存在的真實聲線,這在 VoxCPM 叫 Voice Cloning(聲音克隆),它一定需要一段參考音訊 (reference audio),光靠文字 prompt 是做不到的。

比喻:Voice Design 像跟配音員說「給我一個溫柔女聲」;Voice Cloning 則是「請聽這段錄音,學這個人講話」。文字 prompt 在克隆裡只是「附加指令」,真正的靈魂是那段參考音訊。

克隆時的「prompt」其實分成兩個部分
第一個是參考音訊(決定音色),第二個才是文字 prompt(裡面又分「要唸的內容」和「括號裡的風格控制」)。VoxCPM 有三種模式,prompt 寫法不同:

模式一:基本克隆(只給音訊)
最單純,只複製音色:

python
wav = model.generate(
text="這是用我的參考音訊克隆出來的聲音。",
reference_wav_path="path/to/voice.wav",
)
模式二:可控克隆(音訊 + 風格 prompt)
複製音色的同時,用括號下風格指令,調整語速、情緒、語氣——寫法跟 Voice Design 一樣:

python
wav = model.generate(
text="(語速稍快、開朗的語氣)這是帶風格控制的克隆聲音。",
reference_wav_path="path/to/voice.wav",
cfg_value=2.0,
inference_timesteps=10,
)
這就是你問的「prompt 怎麼寫」的核心:把風格描述放進文字最前面的括號裡,例如:

(平靜、緩慢)
(興奮、語速快)
(悲傷、低沉)
(溫柔、像在說悄悄話)
模式三:極致克隆(音訊 + 逐字稿,最像)
想要最高相似度,就同時給「參考音訊」和「它對應的逐字稿 (transcript)」。為了最像,官方建議把同一段音訊同時傳給 prompt_wav_path 和 reference_wav_path:

python
wav = model.generate(
text="這是極致克隆的示範。",
prompt_wav_path="path/to/voice.wav",
prompt_text="參考音訊裡那個人實際講的逐字稿。", # 一字不差很重要
reference_wav_path="path/to/voice.wav", # 選用,提升相似度
)
重點:prompt_text 要精確對應參考音訊的內容,模型才能「對齊」音訊和文字。

真正決定像不像的,是參考音訊品質(80% 成敗在這)
文字 prompt 只是微調,音訊才是關鍵。準備參考音訊時:

乾淨無雜音:沒有背景音樂、回音、其他人聲。像錄 podcast 那種品質。
長度足夠:建議 10~30 秒的清晰說話,太短資訊不夠。
語氣自然:用平穩、有代表性的說話片段,別挑情緒極端或喊叫的。
單一說話者:整段只有要克隆的那個人。
比喻:這像請人臨摹一幅畫——你給的原稿越清晰、光線越好,臨摹得越像;原稿模糊還有指紋油漬,再厲害也學不準。

行動計畫
錄/找一段乾淨的參考音訊(10~30 秒,WAV 格式最保險)。
先跑模式一,聽基本音色像不像。
不夠像 → 升級模式三(補上精確逐字稿),相似度通常大幅提升。
要調語氣 → 用模式二,在括號加風格指令。
同一段多生 1~3 次挑最好的(克隆有隨機性)。
要長期穩定的專屬聲音 → 用 5~10 分鐘音訊做 LoRA 微調,比即時克隆更穩。
一個必要提醒:克隆能做出非常逼真的聲音,請只克隆你自己的聲音,或已取得對方同意的聲音。官方明文禁止用於假冒、詐騙或散播假訊息,也建議標示為 AI 生成。

要的話,我可以把這套克隆 prompt 寫法整理進前面那份部署指南,變成獨立一節讓你一起存。需要嗎?

寫出來

8:54 PM

決定撰寫獨立的聲音克隆指南文件。
決定撰寫獨立的聲音克隆指南文件。
好的,幫你整理成一份獨立的克隆指南檔,可以跟前面兩份一起存。

Created a file, read a file
Created a file, read a file
寫好了,可以下載。這份克隆指南涵蓋六個部分:克隆與 Voice Design 的差別、克隆 prompt 的兩個組成、三種克隆模式的程式碼、參考音訊的準備要點、一鍵測試三種模式的完整腳本,以及使用倫理。

最後幫你抓三個重點:模仿真實聲音一定要有參考音訊,光靠文字做不到;像不像 80% 取決於音訊品質(乾淨、10~30 秒、單一說話者);文字 prompt 的角色只是用括號做風格微調(語速、情緒、語氣)。

現在你手上有三份可以一起用的檔案了:部署指南、10 種聲音樣板、聲音克隆指南。等你實際在 Mac mini 上跑起來,如果遇到報錯或克隆效果不理想,把情況告訴我,我幫你排查。

Voxcpm 聲音克隆指南
Document · MD

Usage limit reached ∙ Resets 1:20 AM ∙ limits shared with Claude Code
Keep working

Claude is AI and can make mistakes. Please double-check responses.
