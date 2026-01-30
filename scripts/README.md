# Scripts

Utility scripts for the Health Assistant project.

## Integration Testing

### test-integration.sh

End-to-end integration test for the Garmin data ingestion system.

**What it does:**
1. ✓ Checks environment configuration (.env file)
2. ✓ Starts all required containers (PostgreSQL, Ingestion Service, Scheduler)
3. ✓ Waits for all services to become healthy
4. ✓ Gets baseline data counts
5. ✓ Triggers manual data sync from Garmin
6. ✓ Validates data was pulled and stored
7. ✓ Checks sync audit logs
8. ✓ Tests API endpoints
9. ✓ Displays summary report

**Prerequisites:**
- Docker and docker-compose installed
- `.env` file in `infra/` with valid Garmin credentials
- jq (optional, for prettier JSON output)

**Usage:**

```bash
# Run full integration test
./scripts/test-integration.sh

# Run test and cleanup containers after
./scripts/test-integration.sh --cleanup

# Run test without showing logs
./scripts/test-integration.sh --skip-logs

# Show help
./scripts/test-integration.sh --help
```

**Example Output:**

```
========================================
Integration Test Summary
========================================

┌─────────────────────────────────────────┐
│        Integration Test Results         │
├─────────────────────────────────────────┤
│ Total Events:                      12 │
│ Total Sync Runs:                    8 │
│ Successful Syncs:                   7 │
│ Records Fetched:                   12 │
│ Records Inserted:                   5 │
│ Records Updated:                    7 │
└─────────────────────────────────────────┘

[SUCCESS] Integration test PASSED ✓
```

**What it validates:**

- ✓ All containers start successfully
- ✓ Database schema is initialized (events, sync_audit tables)
- ✓ Services respond to health checks
- ✓ Garmin API connection works
- ✓ Data flows from Garmin → Python → Go → PostgreSQL
- ✓ Audit logging captures sync metrics
- ✓ API endpoints return data
- ✓ Insert vs update logic works correctly

**Troubleshooting:**

If the test fails, check:

1. **Environment configuration:**
   ```bash
   cat infra/.env
   # Ensure GARMIN_EMAIL and GARMIN_PASSWORD are set
   ```

2. **Container logs:**
   ```bash
   docker logs health-assistant-garmin-scheduler
   docker logs health-assistant-ingestion-service
   docker logs health-assistant-db
   ```

3. **Service health:**
   ```bash
   curl http://localhost:8083/health
   curl http://localhost:8085/health
   ```

4. **Database connection:**
   ```bash
   docker exec health-assistant-db pg_isready -U healthuser -d health_assistant
   ```

## Database Scripts

### db/init.sql

Initial database schema with:
- TimescaleDB extension
- Users table
- Events table (converted to hypertable)
- Experiments table
- Continuous aggregates
- Test data

### db/migrations/002_sync_audit.sql

Migration to add sync_audit table for tracking data ingestion runs.

**Apply manually if needed:**
```bash
docker exec health-assistant-db psql -U healthuser -d health_assistant -f /docker-entrypoint-initdb.d/migrations/002_sync_audit.sql
```

## Continuous Integration

The `test-integration.sh` script can be used in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Integration Tests
  run: |
    ./scripts/test-integration.sh --cleanup
  env:
    GARMIN_EMAIL: ${{ secrets.GARMIN_EMAIL }}
    GARMIN_PASSWORD: ${{ secrets.GARMIN_PASSWORD }}
```

**Note:** Use test Garmin credentials in CI, not production credentials.
