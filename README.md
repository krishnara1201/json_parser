# JSON Parser (Go)

A small JSON lexer + parser project written in Go, with a minimal web UI to test parsing from the browser.

## Features

- Custom JSON lexer (`lexer/`)
- Custom JSON parser (`parser/`)
- AST value types (`ast/`)
- Browser UI for input and parse output (`main.go`)
- Health endpoint for hosting checks (`/healthz`)

## Project Structure

```text
ast/
lexer/
parser/
main.go
go.mod
Dockerfile
.dockerignore
```

## Requirements

- Go `1.26.1` (matches `go.mod`)

## Run Locally

```bash
go run main.go
```

Open:

- `http://localhost:8080`

If you want a custom port:

```bash
PORT=9090 go run main.go
```

Then open `http://localhost:9090`.

## Web Endpoints

- `GET /` -> parser UI
- `POST /parse` -> parse submitted JSON
- `GET /healthz` -> health check (`ok`)

## Run Tests

```bash
go test ./...
```

## Docker

Build image:

```bash
docker build -t json-parser .
```

Run container:

```bash
docker run --rm -p 8080:8080 -e PORT=8080 json-parser
```

Open `http://localhost:8080`.

## Deploy (GitHub + Render)

1. Push your latest code to GitHub.
2. In Render, create a new **Web Service** from your GitHub repo.
3. Choose Docker runtime (Render will use your `Dockerfile`).
4. Set health check path to `/healthz`.
5. Deploy.

Render provides a public URL with HTTPS.

## Notes

- The app reads `PORT` from environment for cloud hosting compatibility.
- If parser errors occur, they are shown in the UI under **Errors**.
