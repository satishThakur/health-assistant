# Health Assistant — Knowledge Graph

## Project Overview
Personal health assistant that syncs Garmin wearable data, accepts daily check-ins, and surfaces correlations between metrics.

## Services (2 total)
| Service | Language | Port | Entry Point |
|---------|----------|------|-------------|
| Go Backend | Go 1.21+ | 8083 | `backend/cmd/server/main.go` |
| Garmin Scraper | Python | N/A | `scripts/garmin_sync.py` |

## Repository Structure
```
health-assistant/
  backend/          # Go server (see backend/CLAUDE.md)
  mobile_app/       # Flutter app (see mobile_app/CLAUDE.md)
  scripts/          # Python Garmin sync (see scripts/CLAUDE.md)
  scripts/db/       # SQL migrations
```

## Key Architectural Decisions
- **Single Go binary** — all HTTP routes served from `cmd/server`; no microservices
- **Feature-based packages** — each domain owns its handler + repository + validator (`checkin/`, `garmin/`, `audit/`, `auth/`)
- **Event store pattern** — all health data stored in `events` table as JSON blobs, keyed by `(user_id, event_type, time)`
- **TimescaleDB** — Postgres extension for time-series; `events` is a hypertable
- **JWT auth** — HS256, 24h expiry; Google Sign-In → JWT exchange at `POST /api/v1/auth/google`
- **Ingest secret** — Garmin routes use `X-Ingest-Secret` header (server-to-server), not JWT

## Auth Flow
```
Google Sign-In → idToken → POST /api/v1/auth/google → JWT
JWT → Authorization: Bearer header → all user API calls
```

## Database Migrations
Located at `scripts/db/migrations/`. Run in order (001 → 005+).

## Environment Variables
| Var | Service | Notes |
|-----|---------|-------|
| `DATABASE_URL` | backend | Postgres connection string |
| `JWT_SECRET` | backend | min 32 chars |
| `GOOGLE_CLIENT_ID` | backend | empty = dev mode (skip audience check) |
| `GARMIN_INGEST_SECRET` | backend + scripts | shared secret |

## Updating CLAUDE.md Files
- On every commit: update the **local** `CLAUDE.md` for the module you changed
- If you add a new package, new env var, or change architecture: update this root file too
- Keep entries concise — these files are loaded into every session context
