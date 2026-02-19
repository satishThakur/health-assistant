# Backend — Go Server

## Entry Point
`cmd/server/main.go` — wires all dependencies and starts HTTP server on `:8083`

## Package Structure
```
internal/
  auth/       Handler, UserRepository, TokenService, GoogleVerifier
  checkin/    Handler, Repository, Payload validator
  dashboard/  Handler (imports checkin.Repository, no separate repo)
  garmin/     Handler, payload validators (6 types), transform functions
  audit/      Handler, Repository (SyncAudit)
  db/         Database (pool), EventRepository (shared event CRUD)
  middleware/ WithAuth (JWT), WithIngestSecret (Garmin routes)
  models/     Event, User, GarminSleep/Activity/HRV/Stress/DailyStats/BodyBattery
  config/     Config loader (env vars)
```

## Conventions
- Each feature package uses `Handler` and `Repository` (not `XxxHandler`)
- Callers reference as `checkin.NewHandler(...)`, `garmin.NewHandler(...)`
- `db.EventRepository` is shared — stays in `db/` to avoid circular imports
- No circular imports: feature packages → `db`; `db` → nothing internal

## Adding a New Endpoint
1. Add handler method to the relevant feature package's `handler.go`
2. Register route in `cmd/server/main.go` with appropriate middleware
3. Add repository method if DB access needed

## Adding a New Feature Package
1. Create `internal/<feature>/handler.go` and `repository.go`
2. Import from `cmd/server/main.go`
3. Update root `CLAUDE.md` if it changes architecture

## API Route Map
| Method | Path | Middleware | Handler |
|--------|------|------------|---------|
| POST | /api/v1/auth/google | public | auth.HandleGoogleAuth |
| POST | /api/v1/garmin/ingest/* | ingest-secret | garmin.Handle* |
| POST | /api/v1/checkin | JWT | checkin.HandleSubmission |
| GET | /api/v1/checkin/latest | JWT | checkin.HandleGetLatest |
| GET | /api/v1/checkin/history | JWT | checkin.HandleGetHistory |
| GET | /api/v1/dashboard/today | JWT | dashboard.HandleGetToday |
| GET | /api/v1/trends/week | JWT | dashboard.HandleGetWeekTrends |
| GET | /api/v1/insights/correlations | JWT | dashboard.HandleGetCorrelations |
| POST | /api/v1/audit/sync | JWT | audit.HandlePostSyncAudit |
| GET | /api/v1/audit/sync/recent | JWT | audit.HandleGetRecentSyncAudits |
| GET | /api/v1/audit/sync/by-type | JWT | audit.HandleGetSyncAuditsByType |
| GET | /api/v1/audit/sync/stats | JWT | audit.HandleGetSyncAuditStats |

## Build & Test
```bash
cd backend
go build ./...
go test ./...
```
