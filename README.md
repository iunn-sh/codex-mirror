![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/iunn-sh/codex-mirror?color=00ADD8&logo=go&logoColor=white&style=for-the-badge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/iunn-sh/codex-mirror?style=for-the-badge)
![GitHub](https://img.shields.io/github/license/iunn-sh/codex-mirror?style=for-the-badge)

![GitHub CI](https://img.shields.io/github/actions/workflow/status/iunn-sh/codex-mirror/main.yml?logo=github&style=for-the-badge) 
![GitHub CD](https://img.shields.io/github/deployments/iunn-sh/codex-mirror/github-pages?logo=github&style=for-the-badge)

# Codex Mirror for Formosa & Pescadores

Hosting on Github Pages https://iunn-sh.github.io/codex-mirror

## Local Development

```bash
# Go: process data 
## with Docker
docker run --rm $(docker build -t codex-mirror -q .)
docker run --rm $(docker build -t codex-mirror --progress=plain --no-cache .) # debug
## with Golang (tested with 1.20)
go fmt ./...
go run main.go

# Mkdocs Material: host frontend
docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:9.1.5
# visit http://localhost:8000/ from browser
```

## Production Deployment

Commit or merge PR to `main` branch
