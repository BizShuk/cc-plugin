#!/usr/bin/env python3
import sys
import os
import yaml
import datetime

# Root research directory
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CONFIDENCE_FILE = os.path.join(BASE_DIR, 'foundation', 'confidence.yaml')
CONFIG_FILE = os.path.join(BASE_DIR, 'engine', 'config.yaml')

def load_yaml(filepath):
    if not os.path.exists(filepath):
        return {}
    with open(filepath, 'r', encoding='utf-8') as f:
        return yaml.safe_load(f) or {}

def save_yaml(data, filepath):
    with open(filepath, 'w', encoding='utf-8') as f:
        yaml.safe_dump(data, f, allow_unicode=True, default_flow_style=False)

def decay():
    """Run confidence decay rules based on engine/config.yaml"""
    config = load_yaml(CONFIG_FILE)
    decay_config = config.get('confidence', {})
    decay_per_period = decay_config.get('decay_per_period', 0.05)
    demotion_threshold = decay_config.get('demotion_threshold', 0.3)
    
    registry = load_yaml(CONFIDENCE_FILE)
    entries = registry.get('entries', [])
    
    updated_entries = []
    print("Running confidence decay simulation/process...")
    for entry in entries:
        # Simulate decay: e.g., if last_cited is older or as a simple command line action,
        # we decrement confidence of entries that have not been cited recently.
        # For simplicity, this decays every entry by the configured amount unless verified today.
        old_conf = entry.get('confidence', 0.5)
        new_conf = round(max(0.0, old_conf - decay_per_period), 2)
        entry['confidence'] = new_conf
        print(f"Decayed {entry['id']} ({entry['topic']}): {old_conf} -> {new_conf}")
        
        if new_conf < demotion_threshold:
            print(f"  [WARNING] {entry['id']} is below threshold ({demotion_threshold}) and should be demoted/retired!")
        updated_entries.append(entry)
        
    registry['entries'] = updated_entries
    save_yaml(registry, CONFIDENCE_FILE)
    print("Decay complete.")

def promote(item_id, target_layer=None):
    """Promote or demote a knowledge item"""
    registry = load_yaml(CONFIDENCE_FILE)
    entries = registry.get('entries', [])
    
    target_entry = None
    for entry in entries:
        if entry['id'] == item_id:
            target_entry = entry
            break
            
    if not target_entry:
        print(f"Error: Knowledge item {item_id} not found in registry.")
        sys.exit(1)
        
    current_layer = target_entry.get('layer')
    
    # Simple promotion logic if target_layer not specified
    if not target_layer:
        if current_layer == 'hypothesis':
            target_layer = 'principle'
        elif current_layer == 'principle':
            target_layer = 'axiom'
        else:
            print(f"Item {item_id} is already in the highest layer (axiom).")
            return
            
    # Move files
    old_dir = os.path.join(BASE_DIR, 'foundation', f"{current_layer}s")
    new_dir = os.path.join(BASE_DIR, 'foundation', f"{target_layer}s")
    
    filename = f"{item_id}-{target_entry['topic']}.md"
    old_path = os.path.join(old_dir, filename)
    new_path = os.path.join(new_dir, filename)
    
    if os.path.exists(old_path):
        os.makedirs(new_dir, exist_ok=True)
        os.rename(old_path, new_path)
        print(f"Moved physical file: {old_path} -> {new_path}")
    else:
        print(f"Warning: File {old_path} not found physically. Only updating metadata.")
        
    target_entry['layer'] = target_layer
    # Boost confidence upon promotion
    config = load_yaml(CONFIG_FILE)
    boost = config.get('confidence', {}).get('verification_boost', 0.1)
    target_entry['confidence'] = min(1.0, round(target_entry.get('confidence', 0.5) + boost, 2))
    target_entry['last_verified'] = datetime.date.today().isoformat()
    
    save_yaml(registry, CONFIDENCE_FILE)
    print(f"Promoted {item_id} to {target_layer} (Confidence: {target_entry['confidence']})")

def main():
    if len(sys.argv) < 2:
        print("Usage: kb_manager.py [decay|promote] [args...]")
        sys.exit(1)
        
    cmd = sys.argv[1]
    if cmd == 'decay':
        decay()
    elif cmd == 'promote':
        if len(sys.argv) < 3:
            print("Usage: kb_manager.py promote <item_id> [target_layer]")
            sys.exit(1)
        item_id = sys.argv[2]
        target_layer = sys.argv[3] if len(sys.argv) > 3 else None
        promote(item_id, target_layer)
    else:
        print(f"Unknown command: {cmd}")

if __name__ == '__main__':
    main()
