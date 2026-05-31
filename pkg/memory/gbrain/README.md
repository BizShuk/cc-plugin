# GBrain

# Setup embedding model

## Reset pglite

重新 Init gbrain（embedding model 需重新 sizing schema）

```bash
mv ~/.gbrain/brain.pglite ~/.gbrain/brain.pglite.bak
gbrain init --pglite --embedding-model ollama:bge-m3 --embedding-dimensions 1024

export PATH="$HOME/.bun/bin:$PATH"
gbrain doctor --json | jq '.checks[] | select(.name=="embedding_provider" or .name=="embeddings")'
```

### import markdown

```bash
mkdir -p ~/brain && git init # 或沿用既有 brain repo
gbrain import ~/brain/ --no-embed
gbrain embed --stale
```
