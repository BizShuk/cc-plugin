#!/bin/bash
# Scaffolds a new execution cycle file for the Knowledge Engine
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_DIR="$(dirname "$SCRIPT_DIR")"

CURRENT_MONTH=$(date +"%Y-%m")
CURRENT_DATE=$(date +"%Y-%m-%d")
TARGET_DIR="$BASE_DIR/cycles/$CURRENT_MONTH"
mkdir -p "$TARGET_DIR"

# Count existing cycle files for today to determine NNN
COUNT=$(find "$TARGET_DIR" -name "${CURRENT_DATE}-*.md" | wc -l)
NEXT_NUM=$(printf "%03d" $((COUNT + 1)))
TARGET_FILE="$TARGET_DIR/${CURRENT_DATE}-${NEXT_NUM}.md"

echo "Initializing new execution cycle file: $TARGET_FILE"

cat << EOF > "$TARGET_FILE"
---
🔍【探・未知】
  盲點：
  層次：
  探索脈絡：

🧬【合・辯證】
  策略 A（和諧融入）：
  策略 B（辯證衝突）：
  策略 C（速度優先）：
  策略 D（漸進驗證）：
  推薦策略：

⚖️【衡・四軸】
  系統價值： /10 — 
  業務價值： /10 — 
  認知價值： /10 — 
  反脆弱性： /10 — 
  總分： /40
  判定：

🧱【沉・基礎】
  目標層級：
  內容：
  信心度：
  驗證方式：
---
EOF

echo "Done. Please edit the generated cycle template to complete this execution cycle."
