import os
import json
import glob
import soundfile as sf
from voxcpm import VoxCPM

# 載入模型(記憶體不足可改 "openbmb/VoxCPM-0.5B")
model = VoxCPM.from_pretrained("openbmb/VoxCPM2", load_denoiser=False)

# 統一試聽用的句子,可自行修改
DEMO_TEXT = "你好,這是一段聲音樣板的試聽。歡迎使用 VoxCPM 文字轉語音。"

# 建立輸出目錄
os.makedirs("voice_samples", exist_ok=True)

# 搜尋目前目錄下的所有 JSON 設定檔
current_dir = os.path.dirname(os.path.abspath(__file__))
json_files = sorted(glob.glob(os.path.join(current_dir, "*.json")))

if not json_files:
    print("找不到任何 JSON 設定檔。")
    exit(1)

for json_path in json_files:
    try:
        with open(json_path, 'r', encoding='utf-8') as f:
            data = json.load(f)
        
        name = data.get("id")
        desc = data.get("description")
        cfg = data.get("cfg_value", 2.0)
        steps = data.get("inference_timesteps", 10)
        
        if not name or not desc:
            print(f"跳過不完整的設定檔: {os.path.basename(json_path)}")
            continue

        text = f"({desc}){DEMO_TEXT}"
        print(f"正在生成: {name} (CFG: {cfg}, Steps: {steps})...")
        
        wav = model.generate(
            text=text,
            cfg_value=cfg,
            inference_timesteps=steps,
        )
        out_path = os.path.join("voice_samples", f"{name}.wav")
        sf.write(out_path, wav, model.tts_model.sample_rate)
        print(f"已生成: {out_path}")
    except Exception as e:
        print(f"生成 {os.path.basename(json_path)} 時發生錯誤: {e}")

print("全部完成!請到 voice_samples/ 試聽。")
