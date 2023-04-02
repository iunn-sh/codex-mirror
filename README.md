![GitHub build](https://img.shields.io/github/actions/workflow/status/iunn-sh/codex-mirror/main.yml?logo=github&style=for-the-badge) ![GitHub deployments](https://img.shields.io/github/deployments/iunn-sh/codex-mirror/github-pages?logo=github&style=for-the-badge)

# Codex Mirror for Formosa & Pescadores

:construction: Work in progress

Hosting on Github Pages https://iunn-sh.github.io/codex-mirror

## Local Development

```bash
# Go: backend
docker run --rm $(docker build -t codex-mirror -q .) # build
docker run --rm $(docker build -t codex-mirror --progress=plain --no-cache .) # debug

# Mkdocs Material: frontend
docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:9.1.2
# visit http://localhost:8000/ from browser
```

## Production Deployment

Commit or merge PR to `main` branch
