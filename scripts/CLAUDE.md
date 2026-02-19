# Scripts — Python Garmin Sync

## Purpose
Polls Garmin Connect API and POSTs data to the Go backend's ingest endpoints.

## Auth
Uses `X-Ingest-Secret` header (value from `GARMIN_INGEST_SECRET` env var).
No JWT — this is server-to-server communication.

## Ingest Endpoints (all POST)
| Path | Payload key |
|------|-------------|
| /api/v1/garmin/ingest/sleep | `sleep_data` |
| /api/v1/garmin/ingest/activity | `activity_data` |
| /api/v1/garmin/ingest/hrv | `hrv_data` |
| /api/v1/garmin/ingest/stress | `stress_data` |
| /api/v1/garmin/ingest/daily-stats | `daily_stats_data` |
| /api/v1/garmin/ingest/body-battery | `body_battery_data` |

All payloads include `user_id` (UUID) and `date` (YYYY-MM-DD).

## Database Migrations
`db/migrations/` — numbered SQL files, run in order against TimescaleDB.
