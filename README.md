![GitHub build](https://img.shields.io/github/actions/workflow/status/iunn-sh/codex-mirror/main.yml?logo=github&style=for-the-badge) ![GitHub deployments](https://img.shields.io/github/deployments/iunn-sh/codex-mirror/github-pages?logo=github&style=for-the-badge)

# Codex Mirror for Formosa & Pescadores

:construction: Work in progress

Hosting on Github Pages https://iunn-sh.github.io/codex-mirror

## Local Development

```bash
# Go: process data
docker run --rm $(docker build -t codex-mirror -q .)
docker run --rm $(docker build -t codex-mirror --progress=plain --no-cache .) # debug
# or `go run main.go` -> tested with Golang 1.20

# Mkdocs Material: host frontend
docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:9.1.5
# visit http://localhost:8000/ from browser
```

## Production Deployment

Commit or merge PR to `main` branch
