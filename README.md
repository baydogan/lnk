# lnk

A self-hosted, **CLI-only** URL shortener for developers who want to run their own
shortener from the terminal. No web UI.

Two binaries share the same `internal/*` packages:

| Binary | Role |
|--------|------|
| `lnkd` | server daemon — runs the HTTP server (Gin + MongoDB); its home is the Docker image |
| `lnk`  | client CLI — `shorten`, `ls`, `login`, … **and** `lnk init` (server bootstrap) |

## Features

- Shorten URLs with optional custom `--alias` and `--expires` TTL.
- Root-level redirects (`base_url/<code>`) resolved by code **or** alias; every hit is
  counted (302, never cached).
- QR codes for any short link — terminal ANSI or PNG (`-o file.png`), generated on demand.
- API-key auth (`Authorization: Bearer`) on `/api/v1`; the redirect and `/health` stay public.
- Split config: server settings and client credentials never collide.
- Ships as a distroless Docker image with a full `docker compose` stack.

## Requirements

- Go **1.26+** (to build the `lnk` CLI)
- Docker + Docker Compose (to run the server stack)

## Quick start

```bash
# 1. Build and install the CLI (lnk) onto your PATH
make install

# 2. Start the server stack (MongoDB + Redis + lnkd) in Docker
make up

# 3. Grab the admin key printed on first run
docker compose logs lnkd
#   Admin API key generated. Run:
#     lnk login --server http://localhost:8080 --api-key lnk_a3f9k2...

# 4. Configure the CLI with that key
lnk login --server http://localhost:8080 --api-key lnk_a3f9k2...

# 5. Shorten away
lnk shorten https://example.com
```

For local development against containerized dependencies only:

```bash
make dev   # starts mongodb + redis in Docker, runs lnkd from source
```

## CLI commands

```bash
# Server setup (admin, on the host running the server)
lnk init                          # interactive wizard, writes ~/.lnk/server.yaml
lnk init --mode single --base-url http://localhost:8080   # flags skip the prompts
lnkd                              # run the HTTP server daemon

# Client setup
lnk login --server http://localhost:8080 --api-key lnk_...
lnk logout                        # remove ~/.lnk/config.yaml (local session only)

# Links
lnk shorten <url>
lnk shorten <url> --alias mylink --expires 7d --qr
lnk ls                            # list your links
lnk rm <code>                     # delete (confirm; -y to skip)
lnk stats <code>                  # click stats (code or alias)
lnk qr <code>                     # terminal QR, or -o file.png
```

`--expires` accepts relative TTLs: `30m`, `1h`, `7d`, `2w`.

## Configuration

Config is split by responsibility into two files under `~/.lnk` (dir `0700`, files `0600`):

| File | Owner | Contents |
|------|-------|----------|
| `~/.lnk/server.yaml` | server | `mode`, `base_url`, `mongo_uri`, `redis_addr`, `admin` |
| `~/.lnk/config.yaml` | client | `server` URL + `api_key` |

The daemon resolves settings **env > `server.yaml` > built-in default**. Overrides:

| Env var | Purpose |
|---------|---------|
| `MONGO_URI` | MongoDB connection string |
| `REDIS_ADDR` | Redis address (reserved; not yet wired) |
| `BASE_URL` | Public origin used in generated short links |
| `MODE` | `single` (default) or `multi` |
| `PORT` | HTTP port (default `8080`) |
| `LNK_SERVER_CONFIG` | Override the `server.yaml` path |

`base_url` is the address the browser actually hits on a click — it must resolve to the
server in DNS. Locally that is `http://localhost:8080`.

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/v1/shorten` | required | Create short link |
| `GET` | `/api/v1/urls` | required | List links |
| `GET` | `/api/v1/urls/:code` | required | Click stats (code or alias) |
| `DELETE` | `/api/v1/:code` | required | Delete link (code or alias) |
| `GET` | `/:code` | none | Redirect (public) |
| `GET` | `/health` | none | Health check |

## Testing

```bash
make test              # unit tests — fast (~2s), no Docker required
make test-integration  # repository tests against ephemeral MongoDB (testcontainers; needs Docker)
make test-all          # everything
make test-race         # unit tests with the race detector
make cover             # unit tests with coverage
```

Repository integration tests are behind the `integration` build tag, so the default
`make test` loop stays fast and Docker-free.