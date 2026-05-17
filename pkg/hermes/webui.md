<https://github.com/nesquena/hermes-webui>

```bash
git clone https://github.com/nesquena/hermes-webui
cd hermes-webui
cp .env.docker.example .env
# Edit .env if your host UID isn't 1000 (e.g. macOS where UIDs start at 501)
brew install docker-compose docker-buildx
# docker buildx build .
docker-compose up -d
# Open http://localhost:8787




```
