![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/iunn-sh/codex-mirror?color=00ADD8&logo=go&logoColor=white&style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/iunn-sh/codex-mirror?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/iunn-sh/codex-mirror?style=flat-square)](https://goreportcard.com/report/iunn-sh/codex-mirror)
![License](https://img.shields.io/github/license/iunn-sh/codex-mirror?style=flat-square)
![GitHub CI](https://img.shields.io/github/actions/workflow/status/iunn-sh/codex-mirror/main.yml?logo=github&style=flat-square)
![GitHub CD](https://img.shields.io/github/deployments/iunn-sh/codex-mirror/github-pages?logo=github&style=flat-square)
![Website](https://img.shields.io/website?style=flat-square&url=https%3A%2F%2Fiunn-sh.github.io%2Fcodex-mirror)

# Codex Mirror for Formosa & Pescadores

:warning: Personal project to track change of laws. USE AT YOUR OWN DISGRESSION.

Please refer to [全國法規資料庫](https://law.moj.gov.tw/) for the most accurate and up-to-date material.

Hosting on Github Pages https://iunn-sh.github.io/codex-mirror

| dir  		| content						| format					|
| --------: | ----------------------------- | :------------------------ |
| `raw`		| downloaded and unzipped files	| `.zip` / `.csv` / `.json`	|
| `depot`	| parsed data from raw			| `.json`					|
| `docs`	| processed data for frontend	| `.md`						|

## Local Development

```bash
# Go: process data 
## with Docker
docker run --rm $(docker build -t codex-mirror -q .)
docker run --rm $(docker build -t codex-mirror --progress=plain --no-cache .) # debug
## with Golang (tested with 1.20)
go fmt ./...
go mod tidy
go run .

# Mkdocs Material: host frontend
docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material:9.1.5
# visit http://localhost:8000/ from browser
```

## Production Deployment

Commit or merge PR to `main` branch
