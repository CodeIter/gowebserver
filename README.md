# Go web server

Small example HTTP server written in Go. Serves HTML views, static assets, and a small API surface.

## Quickstart

Requirements:
- Go 1.26+

Build:

```bash
make build
```

Run:

```bash
make run
# or
./bin/server
# or for development
go run ./cmd/server
```

## Configuration

This project supports `.env` files (loaded automatically) and a few command-line flags.

- Environment file: copy `.env.example` → `.env` and edit values.
- Flags (highest priority):
	- `-host` (default `0.0.0.0`)
	- `-port` (default `8000`)
	- `-external` (path served at `/external/` for large files)
	- `-version` (prints binary version)

`.env` variables (examples in `.env.example`): `HOST`, `PORT` ...

## Embedded assets

`views/`, `static/`, and `public/` are embedded into the binary via `assets.go`. The server always serves those embedded files at runtime; to change the content, edit the files under their respective folders and rebuild the binary.

## Routes

- GET `/` - home page (rendered with Go `html/template`)
- GET `/*` - public assets served at `/` route
- GET `/static/*` - static assets
- GET `/health` - health check (JSON)
- GET `/ready` - readiness check (JSON)
- GET `/external/*` - served from the path set by `-external`
- GET `/robots.txt`

## Templates

Templates are located in `views/` and include:
- `layout.html` — base layout
- `meta.html` — PWA/meta includes
- `home.html` — home page content

## License

MIT License