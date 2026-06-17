/
I try to use ai video model but I don’t know how to make the results better. Guide me to general prompt, and how to separate scene on order to generate property.

8:25 PM

Synthesized prompting principles and scene-structuring strategies for video generation
Synthesized prompting principles and scene-structuring strategies for video generation
AI 影片生成的成果好不好，關鍵在於你怎麼「描述」和怎麼「拆解」。我先講清楚一個 prompt（提示詞）該有哪些骨架，再教你怎麼分鏡（scene breakdown），最後給你一套可以照做的流程。

一個好的影片 Prompt 該包含什麼
把寫 prompt 想像成「請一位攝影師幫你拍片」。你不會只說「拍一隻貓」，而是會告訴他：拍什麼、怎麼拍、什麼氛圍。一個完整的 prompt 通常有這六個要素：

1. 主體（Subject）——畫面裡的主角是誰／是什麼，加上具體外觀。
   不要寫「a woman」，要寫「a young woman with short black hair, wearing a beige trench coat」。越具體，模型越不會「自由發揮」。

2. 動作（Action）——主體在做什麼。
   影片跟圖片最大的差別就是「會動」，所以動作一定要寫清楚：「slowly turning her head to look at the camera」。

3. 場景與環境（Setting / Environment）——在哪裡、周圍有什麼。
   例如「in a rainy Tokyo street at night, neon signs reflecting on wet pavement」。

4. 鏡頭語言（Camera / Shot type）——這是新手最常漏掉、卻最能提升質感的部分。
   就像導演會說「給我一個特寫」，你要指定：

景別：close-up（特寫）、medium shot（中景）、wide shot（遠景）
運鏡：dolly in（推近）、pan left（左搖）、tracking shot（跟拍）、static shot（固定鏡頭）
鏡頭：shallow depth of field（淺景深，背景模糊）5. 光線與氛圍（Lighting / Mood）——決定整支片「好不好看」的靈魂。
例如「golden hour, warm soft light」（黃金時刻、暖色柔光）或「moody, low-key lighting」（陰暗低調光）。

6. 風格（Style）——寫實、動畫、電影感、復古膠片等。
   例如「cinematic, shot on 35mm film」或「3D Pixar style animation」。

把上面組起來，一個高品質 prompt 長這樣：

Medium shot, a young woman with short black hair in a beige trench coat slowly turning to look at the camera, rainy Tokyo street at night, neon reflections on wet ground, shallow depth of field, cinematic, moody lighting.

生活化比喻：寫 prompt 就像「點一杯客製化飲料」。你只說「珍奶」，店員只能給你標準款；但你說「半糖、少冰、加椰果、用燕麥奶」，拿到的才會接近你心裡想的那杯。模型也是一樣——你給的細節越精準，它猜的空間越小，成果就越穩定。

為什麼要「分鏡」以及怎麼分
目前大多數 AI 影片模型一次只能生成 5～10 秒的片段，而且一個片段裡如果塞太多動作或換場景，模型會「打架」——畫面崩壞、主體變形、邏輯混亂。

所以正確做法是：一個鏡頭只做一件事，再把多個鏡頭接起來。這就是分鏡（storyboard / shot breakdown）。

生活化比喻：這就像拍 Vlog。你不會用一鏡到底拍完整天，而是「早餐一段、出門一段、咖啡廳一段」，最後剪在一起。AI 影片也要這樣切。

分鏡的拆解原則：

每當發生以下情況，就應該「換一個鏡頭」：

場景變了（從室內到室外）
景別變了（從遠景切到特寫）
時間跳躍了（白天變晚上）
主要動作換了（走路→坐下）
舉例，假設你想做「一個人走進咖啡廳點咖啡」，要拆成三個獨立 prompt：

Shot 1（建立鏡頭）：Wide shot, a man pushing open a glass cafe door, warm interior light, morning.
Shot 2（中景動作）：Medium shot, the same man walking to the counter, looking at the menu board.
Shot 3（特寫細節）：Close-up, barista's hands pouring espresso into a white cup, steam rising.
維持一致性（consistency）的小技巧：跨鏡頭時，主角的外觀描述（髮型、服裝、顏色）要「逐字複製貼上」，不要每次換句話講。模型沒有記憶，每個片段都是重新理解，描述一變它就會畫成另一個人。

給你的執行計畫
先寫文字腳本：用一兩句話寫出你要的整段故事，例如「女孩在雨中的東京街頭回頭看鏡頭」。
切分鏡：按上面的「換鏡頭原則」，列出 3～5 個 shot，每個只裝一個動作。
每個 shot 套六要素模板：主體＋動作＋場景＋鏡頭＋光線＋風格，逐句填。記得跨鏡頭時主角描述複製貼上。
單鏡頭測試：先只生成第一個鏡頭，跑 2～3 次看哪個最好。不要一次全做，省時間也省額度。
微調再放大：哪裡崩了就只改那個要素（例如動作太複雜就簡化），確定好了再生成下一個鏡頭。
後製接起來：所有片段都滿意後，用剪輯軟體（甚至手機 CapCut）依順序接合、加轉場與配樂。
一個建議：先從「固定鏡頭＋單一簡單動作」開始練（例如 static shot, a cup of coffee with steam rising）。這種最不容易崩，能讓你快速建立對模型「脾氣」的手感，再慢慢加運鏡和複雜度。

你目前是用哪一個影片模型？不同模型（像 Sora、Kling、Runway、Veo）對 prompt 的吃法略有差異，如果你告訴我，我可以給你更貼該模型的寫法。

Help me do the comparison first

8:27 PM

Scrutinizing pricing details and capabilities comprehensively
Scrutinizing pricing details and capabilities comprehensively
我先查一下目前最新的狀況，因為這些模型版本和能力更新得很快。

Synthesizing comparison framework with table and recommendations
Synthesizing comparison framework with table and recommendations
我再確認一下各家的價格和免費方案，這對剛起步的你最重要。

Synthesized video AI tools into structured comparison framework
Synthesized video AI tools into structured comparison framework
我整理了目前（2026 年 6 月）的主流模型。先給你一個快速比較表，再用生活化的方式幫你挑。

主流 AI 影片模型比較
模型 最強項 解析度／長度 免費方案 適合誰
Kling 3.0（快手） 高 CP 值、運鏡流暢 原生 4K@60fps／約 10 秒 最佛心，每天 66 點、免綁卡 新手、預算有限
Seedance 2.0（字節跳動） 角色一致性最強、12 檔多模態輸入 1080p／約 15 秒 有限 連續劇情、品牌角色
Sora 2（OpenAI） 物理擬真最好、單片最長 高畫質／最長 25 秒 含在 ChatGPT 訂閱中 寫實、複雜場景
Veo 3.1（Google） 電影級畫質＋原生對白語音 4K／含同步音訊 透過 Gemini 有限免費 需要講話、配音的片
Runway Gen-4.5 創作平台工具整合 高畫質 訂閱制 想邊生成邊剪輯
Wan 2.2 / 2.6（開源） 免費、可本地跑 720p～1080p 完全免費（要自備 GPU） 有顯卡、想無限生成
幾個重點補充：Kling 3.0 是主流模型中免費方案最大方的，每天 66 點且不需信用卡。Seedance 2.0 的角色一致性被廣泛認為是業界最強，對連續劇情或品牌角色比單純畫質更重要，而且它是唯一深度整合 TikTok CapCut 剪輯軟體的模型。另外，截至目前，六大模型中有四個原生支援同步音訊：Veo 3.1、Seedance 2.0、Sora 2、Kling 3.0，Veo 3.1 是唯一能生成真正同步對白語音（dialogue）、而不只是音效的模型。
Atlas Cloud + 4

用生活化的方式幫你選
把選模型想成「買車」——沒有最好的車，只有最適合你用途的車：

Kling = 國民神車（如 Toyota）：好上手、省錢、品質穩，新手第一台就選它。每天免費額度夠你大量試錯練手感。
Seedance = 全家用的 MPV：如果你要做「同一個角色出現在多個鏡頭」的連續故事（像短劇、品牌吉祥物），它記人臉的能力最強，角色不會每換一個鏡頭就變臉。
Sora = 越野性能車：要拍複雜的物理互動（水花、碰撞、多人動作），它最不容易穿幫，單片還能拍到 25 秒。
Veo = 附專業音響的房車：你的片需要人物「開口講話」或要電影感配音，它是唯一能做到嘴型對上對白的。
Wan = 自己組裝的車：完全免費，但你得有一張夠力的顯卡，技術門檻較高（以你的技術背景倒是可行）。
給你的選擇與起步計畫
第一步：先確認你的用途。 這決定一切。你可以先回答自己：你的片需要人物講話嗎？是同一角色的連續故事，還是各自獨立的畫面？要不要 4K？

第二步：新手先從 Kling 3.0 起步。 理由是它每天有免費額度、不用綁卡、4K 畫質、運鏡又自然，是練「寫 prompt → 看結果 → 調整」這個循環成本最低的選擇。在台灣可直接用網頁版註冊。

第三步：用同一個 prompt 跨模型測試。 等你抓到手感後，拿你上一則我們講的「六要素 prompt」，在 Kling、Veo（透過 Gemini）、Sora（若你有 ChatGPT 訂閱）各跑一次同樣的內容，直接比結果，挑出最合你題材的那台。

第四步：依結果定位主力＋備援。 通常會變成「一台主力跑大量、一台特殊需求備援」。例如平常用 Kling 量產，需要角色一致時切 Seedance。

提醒一點：價格和版本每幾週就在變，2026 年的 AI 影片市場破碎、競爭激烈，而且每個月都在變便宜，所以真正要訂閱前，記得去官方頁面確認最新價格。
Get AI Perks

你大概是想做哪種片？（短劇／產品廣告／社群短影音／純藝術）告訴我，我可以直接幫你鎖定最適合的那一兩台，省下你比較的時間。

那個它有拆分鏡的來做提示詞的介面

8:32 PM

識別了用戶詢問場景分割提示詞介面的需求。
識別了用戶詢問場景分割提示詞介面的需求。
你是指「內建分鏡腳本（storyboard）介面、可以一鏡一鏡排好再生成」的工具吧？我查一下哪幾家有這功能。

整理了多個影片生成工具的分鏡介面功能。
整理了多個影片生成工具的分鏡介面功能。
有的，這種「一鏡一鏡分開填提示詞、再一次生成」的功能叫 Multi-Shot Storyboard（多鏡分鏡）。目前分兩類：模型內建的，和專門的腳本拆鏡平台。

模型內建分鏡介面的
Kling 3.0——最推薦你先用這個。 它的招牌功能就是這個。你可以定義 3 到 12 個鏡頭，每個鏡頭各自填提示詞、運鏡角度、轉場，然後一次批次生成成一段連貫的敘事。而且跨鏡頭會自動維持視覺一致性。這正是你要的介面，又有免費額度，新手起步首選。
LushBinary
Medium

Seedance 2.0——進階首選。 它也有多鏡敘事，做法類似 Kling，但多了跨鏡頭的「統一音訊連續性」——背景音樂、環境音、對白會自然地跨場景銜接。如果你的片有講話或配樂，它接得比較順。
LushBinary

Veo 3.1——透過 Google Flow 平台。 Flow 裡有 Storyboard 和 Recut 等工具可以排場景，但要注意：Veo 單鏡較短，約 8 秒上限，比較像「排好幾個短鏡再拼」。
AIMLAPI

Sora 2 的 Storyboard——已經停用了，要注意。 OpenAI 已宣布 Sora 的網頁與 App 版在 2026 年 4 月 26 日關閉，API 服務則在 9 月 24 日結束。所以它的分鏡功能現在官方管道已經用不到了，只剩有 API 合作的第三方平台還能接。先別把它當主力。
Noviai

專門的「腳本→分鏡」平台
如果你想要更完整的「貼一段腳本，它自動幫你拆成一格一格分鏡」的體驗，有專門的工具：LTX Studio 和 M Studio 這類平台，是把 AI 影片生成連接到分鏡腳本、配音、審稿、輸出的完整流程。你可以開一個瀏覽器分頁、貼上腳本，M Studio 就把它變成一格一格的分鏡，附帶 AI 生成的畫面、動態腳本（animatic）、音訊和可交付的成品輸出。
M Studio
M Studio

生活化比喻：模型內建的分鏡（像 Kling）就像「相機本身有連拍模式」——夠用、方便；而 LTX / M Studio 這種平台像「整間附導演、剪接、錄音的工作室」——你給劇本，它幫你從頭排到尾。新手用前者就很夠，等你要做長片才需要後者。

給你的起步計畫
先用 Kling 3.0 的 Multi-Shot Storyboard 練手。 它有免費額度又內建這個介面，最低成本就能體驗「拆鏡→各填提示詞→一次生成」的完整流程。
把我們第一則講的「六要素」填進每一格。 分鏡介面裡每個 shot 的輸入框，就是放你的「主體＋動作＋場景＋鏡頭＋光線＋風格」的地方。記得跨鏡頭時主角描述要複製貼上維持一致。
先做 3 鏡的短序列測試。 不要一開始就排滿 12 鏡。先用 3 個鏡頭（建立鏡頭→中景動作→特寫）跑通，確認一致性 OK，再加長。
要做長片或要配音時，再升級到 Seedance 或 LTX/M Studio。 等基本功穩了，需要跨鏡頭音訊連貫或自動拆腳本，再換工具。
你想先做幾秒、幾個鏡頭的東西？如果你給我一句你想拍的內容，我可以直接幫你示範把它拆成 Kling 分鏡介面可以一格一格貼上去的提示詞。

幫我根據場景腳本環境鏡頭 等等 產生一個 transcription skill

8:38 PM

解析需求並規劃建立結構化影片提示轉換技能。
解析需求並規劃建立結構化影片提示轉換技能。
我來幫你做一個可重複使用的 skill，它能把你的場景腳本，自動轉成分鏡＋六要素的影片提示詞。先看一下建立 skill 的規範。

評估主觀輸出技能，略過嚴格基準測試。
評估主觀輸出技能，略過嚴格基準測試。
我看一下 skill 撰寫指南的細節部分。

設計場景轉視頻提示的技能架構與輸出格式。
設計場景轉視頻提示的技能架構與輸出格式。
規範看懂了。你要的是一個能把「場景腳本＋環境＋鏡頭」這些素材，自動轉成分鏡＋六要素影片提示詞的 skill。這類產出偏創意性，所以我會直接幫你寫好、打包成可安裝的 .skill 檔，不用跑一整套評測流程。

我先建立 skill 的資料夾與主檔。

診斷腳本寫入權限問題並尋求可寫入位置。
診斷腳本寫入權限問題並尋求可寫入位置。
主檔和模型小抄都寫好了。現在打包成可安裝的 .skill 檔。

組裝技能檔案並準備中文說明與使用指南。
組裝技能檔案並準備中文說明與使用指南。
打包完成。檔案已驗證通過。

Presented file
Presented file
幫你做好了，這是一個叫 scene-to-video-prompt 的 skill（技能包）。

這個 skill 會做什麼
它就是把我們前幾則討論的東西，固化成一個 Claude 隨時能調用的「自動轉換器」。你只要丟一段場景描述（中文也行，像「一個女生在下雨的東京街頭回頭看鏡頭」），它就會自動幫你：

拆分鏡（scene breakdown）——依「換場景／換景別／換時間／換動作」原則，把腳本切成 3～6 個獨立鏡頭。
每鏡填六要素——主體＋動作＋場景＋鏡頭＋光線＋風格，組成一句完整的英文提示詞（模型對英文最聽話）。
鎖一致性區塊（consistency block）——把角色、關鍵物件、場景的固定描述抽出來，逐字複製進每個鏡頭，避免主角換鏡就變臉。
輸出成可直接貼上的分鏡清單——格式專門對齊 Kling 的 Multi-Shot Storyboard 介面，每個 shot 有提示詞、運鏡、轉場欄位。
裡面還附了一份 model-tips.md，當你指定 Kling／Seedance／Veo／Sora 時，它會自動套用各家的長度、4K、音訊、分鏡支援等差異來調整輸出。

怎麼用（安裝計畫）
下載上面那個 .skill 檔。
到 Claude 的 Settings → Capabilities → Skills（網頁或桌面版）把它上傳安裝。
之後在對話裡只要說「幫我把這個場景做成影片提示詞」或直接貼腳本，它就會自動觸發、產出分鏡清單。
把每個 shot 的內容貼進 Kling 的 Multi-Shot Storyboard 一格一格生成。
生活化比喻：原本你每次都要自己「翻譯」想法成攝影師聽得懂的話；這個 skill 等於請了一位固定的副導，你講白話，他自動寫成規格化的拍攝指令。

要不要我直接拿一個你實際想拍的場景，當場跑一次給你看輸出長什麼樣？如果輸出的格式或預設風格你想調整，我也可以改了再重新打包一版給你。

scene-to-video-prompt
Skill
再幫我做一個腳本生產優化skill 讓我可以從簡單的prompt生成有劇情般的腳本

8:51 PM

設計了提示轉故事腳本的技能架構與輸出格式。
設計了提示轉故事腳本的技能架構與輸出格式。
好主意——這正好是上一個 skill 的「上游」。整條線會變成：一句簡單想法 →（這個新 skill）有劇情的腳本 →（上一個 skill）分鏡提示詞。我來做這個「腳本生成優化」skill。

識別並修復YAML語法衝突問題。
識別並修復YAML語法衝突問題。
YAML 因為描述裡有個冒號出錯了，我修掉它。

Edited a file, ran a command
Edited a file, ran a command
驗證通過、打包完成。

Presented file
Presented file
做好了，這個 skill 叫 prompt-to-story-script。

這個 skill 會做什麼
它解決一個核心問題：一句平淡的 prompt 只是「一張會動的圖」，不是故事。 故事需要「改變」——有人想要某個東西，發生了一個轉折，最後落在跟開頭不同的情緒上。就算只有 15 秒也能有這個小弧線。

你丟一句簡單想法（像「一個女生在下雨的東京」），它會幫你：

定四個錨點——目標長度、類型／情緒（genre/tone）、主角＋他「想要什麼」、以及那個關鍵「轉折（the shift）」。沒講的它會自己推測並標註假設，不會一直反問你。
搭故事弧線（story arc）——短片預設用「鋪陳→轉折→收尾（Setup → Turn → Payoff）」三拍；較長的片會展開成五到七拍。
改寫成「畫面化」的腳本——影片沒有旁白，所以它會把「她很寂寞」這種內心戲，翻成「她看著對街共撐一把傘的情侶，把衣領拉緊」這種鏡頭看得見的動作。
輸出結構化腳本——含 logline（一句話故事）、故事設定、逐拍腳本，最後給一個「punch-up tip（加分提示）」告訴你怎麼改最能提升感染力。
附的 story-structures.md 還收了不同長度（15／30／60 秒）和不同類型（情感、懸疑、喜劇、史詩、產品廣告、MV）的節奏模板，做長片或特定類型時它會自動參考。

兩個 skill 怎麼串起來用
這就是重點——你現在有一條完整生產線：

一句想法 → 跑 prompt-to-story-script → 有劇情的腳本 → 跑 scene-to-video-prompt → 分鏡提示詞 → 貼進 Kling 生成

生活化比喻：第一個 skill 是「編劇」，把你的靈感寫成有起承轉合的劇本；第二個是「副導＋攝影指導」，把劇本拆成一個個可拍的鏡頭指令。兩個接在一起，你只要出一句話，就能一路走到能生成的提示詞。

安裝計畫
下載上面的 .skill 檔（連同上一個 scene-to-video-prompt 一起）。
到 Claude 的 Settings → Capabilities → Skills 各自上傳安裝。
之後直接說「幫我把這個想法變成有劇情的腳本」就會觸發第一個；接著說「把這份腳本做成影片提示詞」就接到第二個。
要不要我現在用你的某個真實想法，把兩個 skill 串起來、一路從「一句話」跑到「可貼進 Kling 的分鏡提示詞」給你看完整效果？如果哪個 skill 的預設（節奏、語言、輸出格式）想調整，我改完可以再幫你重新打包。

prompt-to-story-script
Skill

Usage limit reached ∙ Resets 1:20 AM ∙ limits shared with Claude Code
Keep working
寫一個角色設定的skill

Claude is AI and can make mistakes. Please double-check responses.
