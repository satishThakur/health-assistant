# Quick Start Guide - Garmin Data Ingestion

## TL;DR - Get Running in 3 Steps

```bash
# 1. Configure credentials
cd infra && cp .env.example .env
# Edit .env: add GARMIN_EMAIL and GARMIN_PASSWORD

# 2. Run integration test (does everything)
./scripts/test-integration.sh

# 3. Check results
docker exec -it health-assistant-db psql -U healthuser -d health_assistant \
  -c "SELECT event_type, COUNT(*) FROM events WHERE source='garmin' GROUP BY event_type;"
```

## Manual Commands

### Start Services
```bash
cd infra
docker-compose up -d postgres ingestion-service garmin-scheduler
```

### Check Health
```bash
curl http://localhost:8083/health  # Go ingestion service
curl http://localhost:8085/health  # Python scheduler
```

### Trigger Sync
```bash
curl -X POST http://localhost:8085/sync/trigger
```

### View Data
```bash
# Connect to database
docker exec -it health-assistant-db psql -U healthuser -d health_assistant

# Count events by type
SELECT event_type, COUNT(*) as count FROM events WHERE source='garmin' GROUP BY event_type;

# View sync audit
SELECT data_type, records_fetched, records_inserted, records_updated, status
FROM sync_audit ORDER BY sync_started_at DESC LIMIT 10;

# Get sync stats
SELECT data_type, COUNT(*) as syncs, SUM(records_inserted) as inserted
FROM sync_audit WHERE status='success' GROUP BY data_type;
```

### Check Logs
```bash
docker logs health-assistant-garmin-scheduler     # Python scheduler
docker logs health-assistant-ingestion-service    # Go service
docker logs health-assistant-db                   # PostgreSQL
```

### Query Audit API
```bash
# Recent syncs
curl "http://localhost:8083/api/v1/audit/sync/recent?user_id=00000000-0000-0000-0000-000000000001&limit=10" | jq

# Sync statistics
curl "http://localhost:8083/api/v1/audit/sync/stats?user_id=00000000-0000-0000-0000-000000000001" | jq

# By data type
curl "http://localhost:8083/api/v1/audit/sync/by-type?data_type=sleep&limit=10" | jq
```

## Stop Services
```bash
cd infra
docker-compose down               # Stop all
docker-compose down -v            # Stop and remove volumes (deletes data)
```

## Troubleshooting

### No data syncing?
```bash
# Check scheduler logs for Garmin API errors
docker logs --tail 50 health-assistant-garmin-scheduler

# Verify credentials
docker exec health-assistant-garmin-scheduler env | grep GARMIN
```

### Database connection issues?
```bash
# Check if database is ready
docker exec health-assistant-db pg_isready -U healthuser -d health_assistant

# Check if tables exist
docker exec health-assistant-db psql -U healthuser -d health_assistant \
  -c "\dt"
```

### Services not healthy?
```bash
# Check container status
docker ps

# Restart problematic service
docker-compose restart ingestion-service
```

### Apply migration manually
```bash
docker exec health-assistant-db psql -U healthuser -d health_assistant \
  -f /docker-entrypoint-initdb.d/migrations/002_sync_audit.sql
```

## Configuration

### Sync Schedule (in .env)
```bash
# Every hour (default)
SYNC_CRON_HOUR=*
SYNC_CRON_MINUTE=0

# Every 6 hours
SYNC_CRON_HOUR=*/6
SYNC_CRON_MINUTE=0

# Daily at 8 AM
SYNC_CRON_HOUR=8
SYNC_CRON_MINUTE=0
```

### Multi-User Setup
```bash
# Currently uses DEFAULT_USER_ID from .env
# Future: query users table for per-user credentials
```

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                     Garmin Data Flow                         │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  Garmin API                                                  │
│      ↓                                                       │
│  Python Scheduler (port 8085)                                │
│      ├─ garminconnect library                                │
│      ├─ Fetch sleep, activity, HRV, stress                   │
│      └─ Transform to JSON                                    │
│      ↓ HTTP POST                                             │
│  Go Ingestion Service (port 8083)                            │
│      ├─ Validate payloads                                    │
│      ├─ Transform to models.Event                            │
│      └─ Track insert vs update                               │
│      ↓ pgx/v5                                                │
│  PostgreSQL + TimescaleDB                                    │
│      ├─ events (hypertable)                                  │
│      └─ sync_audit (audit logs)                              │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

## Key Files

```
infra/
  .env                          # Your credentials (create from .env.example)
  docker-compose.yml            # Container orchestration

backend/
  cmd/ingestion-service/        # Go service entry point
  internal/db/                  # Database repositories
  internal/handlers/            # HTTP handlers
  internal/validation/          # Payload validation

services/garmin-scheduler/
  app/
    main.py                     # FastAPI app
    scheduler.py                # APScheduler cron jobs
    garmin_client.py            # Garmin API wrapper
    ingestion_client.py         # HTTP client to Go

scripts/
  test-integration.sh           # End-to-end test script
  db/init.sql                   # Database schema
  db/migrations/                # Schema migrations
```

## Next Steps

1. ✅ Run integration test: `./scripts/test-integration.sh`
2. ✅ View audit data in PostgreSQL
3. ✅ Set up automatic hourly sync (already configured)
4. Monitor sync_audit table for issues
5. Build dashboard using audit API endpoints

## Support

- **Documentation**: See `GARMIN_INTEGRATION_README.md` for detailed docs
- **Script Help**: `./scripts/test-integration.sh --help`
- **Database Schema**: `scripts/db/init.sql` and `scripts/db/migrations/`
- **API Docs**: Check handlers in `backend/internal/handlers/`
