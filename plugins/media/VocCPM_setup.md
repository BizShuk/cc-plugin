# VoxCPM 部署指南(Mac mini / 通用版)

> 一份從零開始、可照做的 VoxCPM 文字轉語音 (Text-to-Speech, TTS) 本機部署手冊。
> 重點:**Apple Silicon Mac mini 可以部署**,關鍵在「統一記憶體 (unified memory)」夠不夠。

---

## 0. 先搞懂:我該裝哪個版本?

VoxCPM 有三個版本,需求差很多。下表的「記憶體需求」在 Mac 上等於要佔用的統一記憶體。

| 版本              | 參數量 | 記憶體需求 | 音質           | 語言   | 狀態 |
| ----------------- | ------ | ---------- | -------------- | ------ | ---- |
| **VoxCPM2**(推薦) | 2B     | ~8 GB      | 48kHz 錄音室級 | 30 種  | 最新 |
| **VoxCPM1.5**     | 0.6B   | ~6 GB      | 44.1kHz        | 中、英 | 穩定 |
| **VoxCPM-0.5B**   | 0.5B   | ~5 GB      | 16kHz          | 中、英 | 舊版 |

**白話選擇法:**

- 記憶體 8GB → 先跑 `VoxCPM-0.5B`(會吃緊,僅供試跑)
- 記憶體 16GB → `VoxCPM2` 的甜蜜點,順跑
- 記憶體 24GB↑ → `VoxCPM2` 最舒服

> **比喻**:統一記憶體像「廚房和餐廳共用的食材櫃」。模型要佔走一大塊,系統還要留一塊,櫃子太小就會卡。

---

## 1. 硬體與環境需求

**硬體**

- ✅ **Apple Silicon(M 系列)Mac mini**:走 MPS(Apple GPU 加速),官方支援。
- ⚠️ **Intel Mac mini**:只能跑 CPU,很慢,建議改用 `VoxCPM.cpp` 路線(見第 7 節)。
- 💡 速度提醒:VoxCPM2 在 NVIDIA RTX 4090 上的即時因子 (Real-Time Factor, RTF) 約 0.3;Mac 的 MPS 會慢一些,**比較適合批次離線生成,不一定能做即時對話**。

**軟體**

- Python ≥ 3.10 且 < 3.13
- PyTorch ≥ 2.5.0
- (官方標示的 CUDA ≥ 12.0 只給 NVIDIA 顯卡;**Mac 走 MPS 不需要 CUDA**)

---

## 2. 安裝(Apple Silicon Mac mini 最簡路線)

**Step 2-1｜建立乾淨的 Python 環境**

像「給這個專案一個獨立工作檯」,避免污染系統 Python。

```bash
conda create -n voxcpm python=3.11 -y
conda activate voxcpm
```

> 沒有 conda 也可以用 venv:
>
> ```bash
> python3.11 -m venv voxcpm-env
> source voxcpm-env/bin/activate
> ```

**Step 2-2｜安裝 VoxCPM**

```bash
pip install voxcpm
```

---

## 3. 第一次測試:生出第一段語音

把以下內容存成 `test.py`。注意 `device="mps"` 讓它走 Mac GPU。

```python
from voxcpm import VoxCPM
import soundfile as sf

# 載入模型(第一次會自動從 HuggingFace 下載權重,2B 檔案較大,請耐心等)
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
```

執行:

```bash
python test.py
```

> **記憶體不足的話**:把 `"openbmb/VoxCPM2"` 換成 `"openbmb/VoxCPM-0.5B"`,需求降到約 5GB。

---

## 4. 聲音設計(Voice Design):用一句話造出新聲音

VoxCPM2 的亮點之一,**不需要參考音訊**。格式是把描述放在文字最前面的括號裡:`(描述)要合成的文字`。

```python
wav = model.generate(
    text="(一位年輕女性,溫柔甜美的聲線)歡迎使用 VoxCPM2!",
    cfg_value=2.0,
    inference_timesteps=10,
)
sf.write("design.wav", wav, model.tts_model.sample_rate)
```

描述可包含:性別、年齡、語氣、情緒、語速等。

---

## 5. 聲音克隆(Voice Cloning):複製一個人的聲音

**5-1｜可控克隆(Controllable Cloning)** — 複製音色 + 加風格指令

```python
wav = model.generate(
    text="(語速稍快、開朗的語氣)這是帶風格控制的克隆聲音。",
    reference_wav_path="path/to/voice.wav",   # 你的參考音訊
    cfg_value=2.0,
    inference_timesteps=10,
)
sf.write("clone.wav", wav, model.tts_model.sample_rate)
```

**5-2｜極致克隆(Ultimate Cloning)** — 同時給音訊與逐字稿,還原最多細節

```python
wav = model.generate(
    text="這是極致克隆的示範。",
    prompt_wav_path="path/to/voice.wav",
    prompt_text="參考音訊對應的逐字稿。",
    reference_wav_path="path/to/voice.wav",   # 選用,可提升相似度
)
sf.write("hifi_clone.wav", wav, model.tts_model.sample_rate)
```

> ⚠️ **倫理提醒**:嚴禁用於假冒、詐騙或散播假訊息。建議清楚標示 AI 生成內容。

---

## 6. 圖形介面(不想寫程式就用這個)

```bash
python app.py --device auto
# Apple Silicon 會自動使用 MPS
# 然後瀏覽器開 http://localhost:8808
```

`--device` 可選值:`auto`、`cpu`、`mps`、`cuda`、`cuda:N`。

---

## 7. 進階:Mac 上跑更快 / 更省記憶體

若覺得 MPS 太慢或記憶體吃緊,改用為 Apple 硬體優化的社群版本:

| 專案            | 用途                                            |
| --------------- | ----------------------------------------------- |
| **VoxCPMANE**   | Apple Neural Engine 後端(專為 M 晶片)           |
| **VoxCPM.cpp**  | GGML/GGUF 量化,支援 CPU / Vulkan,記憶體佔用更低 |
| **VoxCPM-ONNX** | ONNX 匯出,CPU 推論                              |

> Intel Mac mini 建議直接走 `VoxCPM.cpp`。

---

## 8. 進階:訓練專屬聲音(微調 Fine-tuning)

只要 **5~10 分鐘** 的音訊,就能用 LoRA 適配特定說話者或領域。

```bash
# LoRA 微調(參數高效,推薦)
python scripts/train_voxcpm_finetune.py \
    --config_path conf/voxcpm_v2/voxcpm_finetune_lora.yaml

# 訓練 / 推論用的 WebUI
python lora_ft_webui.py   # 開 http://localhost:7860
```

> 💡 微調很吃算力,建議在有 NVIDIA GPU 的機器或雲端做,Mac 訓練會吃力(推論留在 Mac 即可)。

---

## 9. 常見問題與小技巧

- **生出來的聲音不滿意?** 聲音設計/克隆有隨機性,同一段可生成 **1~3 次** 挑最好的。
- **下載很慢?** 可改從 ModelScope 下載:

    ```bash
    pip install modelscope
    ```

    ```python
    from modelscope import snapshot_download
    snapshot_download("OpenBMB/VoxCPM2", local_dir="./pretrained_models/VoxCPM2")
    # 之後用本地路徑載入
    model = VoxCPM.from_pretrained("./pretrained_models/VoxCPM2", load_denoiser=False)
    ```

- **想調語速 / 情緒?** 用括號描述(如 `(語速慢、平靜)`)或在克隆時加風格指令。
- **記憶體不足報錯?** 換小版本(0.5B)、關掉其他吃記憶體的程式,或改用 `VoxCPM.cpp` 量化版。

---

## 10. 建議執行順序(行動計畫)

1. 確認 Mac mini 規格(是否 M 系列、幾 GB 統一記憶體)。
2. 依記憶體選版本(<16GB 先用 0.5B;≥16GB 用 VoxCPM2)。
3. 照第 2、3 節把基本流程跑通,先生出一段 `demo.wav`。
4. 試 Voice Design(第 4 節),熟悉「用描述造聲音」。
5. 需要特定人聲再做 Voice Cloning(第 5 節)。
6. 嫌慢/吃緊 → 換 VoxCPMANE 或 VoxCPM.cpp(第 7 節)。
7. 要長期專屬聲音 → 雲端做 LoRA 微調(第 8 節)。

---

## 參考資源

- GitHub:`https://github.com/OpenBMB/VoxCPM`
- 官方文件:`https://voxcpm.readthedocs.io/en/latest/`
- 模型權重(HuggingFace):`https://huggingface.co/openbmb/VoxCPM2`
- 線上試玩(Playground):`https://huggingface.co/spaces/OpenBMB/VoxCPM-Demo`
- 授權:Apache-2.0(可商用)

---

_本指南依官方 GitHub 倉庫(VoxCPM2,2026 年版)整理。版本與指令可能更新,實作前建議對照官方文件確認。_
