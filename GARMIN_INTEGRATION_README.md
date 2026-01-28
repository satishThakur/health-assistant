# Garmin Data Ingestion Setup

This document describes how to set up and use the Garmin data ingestion system.

## Architecture

```
Python Scheduler Service (Port 8085)
    ↓ Fetch from Garmin API (python-garminconnect)
    ↓ Transform to JSON
    ↓ HTTP POST to Go service
Go Ingestion Service (Port 8083)
    ↓ Validate payloads
    ↓ Transform using internal/models
    ↓ Store in PostgreSQL via internal/db
```

## Quick Start

### 1. Set up environment variables

Copy the example environment file:
```bash
cd infra
cp .env.example .env
```

Edit `.env` and add your Garmin credentials:
```bash
GARMIN_EMAIL=your_email@example.com
GARMIN_PASSWORD=your_password
DEFAULT_USER_ID=00000000-0000-0000-0000-000000000001
```

### 2. Start the services

```bash
cd infra
docker-compose up -d postgres ingestion-service garmin-scheduler
```

### 3. Check service health

```bash
# Check ingestion service
curl http://localhost:8083/health

# Check garmin scheduler
curl http://localhost:8085/health
```

### 4. Trigger manual sync (for testing)

```bash
curl -X POST http://localhost:8085/sync/trigger
```

### 5. Verify data in database

```bash
docker exec -it health-assistant-db psql -U healthuser -d health_assistant

# Query events
SELECT event_type, COUNT(*) as count, MAX(time) as latest
FROM events
WHERE source = 'garmin'
GROUP BY event_type
ORDER BY event_type;

# View recent sleep data
SELECT time, user_id, data
FROM events
WHERE event_type = 'garmin_sleep'
ORDER BY time DESC
LIMIT 5;
```

## API Endpoints

### Go Ingestion Service (Port 8083)

- `GET /health` - Health check
- `POST /api/v1/garmin/ingest/sleep` - Ingest sleep data
- `POST /api/v1/garmin/ingest/activity` - Ingest activity data
- `POST /api/v1/garmin/ingest/hrv` - Ingest HRV data
- `POST /api/v1/garmin/ingest/stress` - Ingest stress data

### Python Scheduler Service (Port 8085)

- `GET /health` - Health check
- `GET /` - Service info
- `POST /sync/trigger` - Manually trigger sync

## Data Flow

### Sleep Data Example

```json
POST http://localhost:8083/api/v1/garmin/ingest/sleep
{
  "user_id": "uuid",
  "date": "2026-01-28",
  "sleep_data": {
    "sleep_time_seconds": 26100,
    "deep_sleep_seconds": 5520,
    "light_sleep_seconds": 15240,
    "rem_sleep_seconds": 5340,
    "awake_seconds": 720,
    "sleep_scores": {"overall_score": 82},
    "average_hrv": 67.5,
    "sleep_end_timestamp_gmt": "2026-01-28T07:45:00Z"
  }
}
```

### Activity Data Example

```json
POST http://localhost:8083/api/v1/garmin/ingest/activity
{
  "user_id": "uuid",
  "date": "2026-01-28",
  "activity_data": {
    "activity_type": "running",
    "start_time_gmt": "2026-01-28T10:00:00Z",
    "duration_seconds": 2700,
    "distance_meters": 5000,
    "calories": 285,
    "average_heart_rate": 132,
    "max_heart_rate": 168
  }
}
```

## Scheduler Configuration

The scheduler runs on a cron schedule. Configure via environment variables:

```bash
# Every hour at minute 0 (default)
SYNC_CRON_HOUR=*
SYNC_CRON_MINUTE=0

# Every day at 8:00 AM
SYNC_CRON_HOUR=8
SYNC_CRON_MINUTE=0

# Every 6 hours
SYNC_CRON_HOUR=*/6
SYNC_CRON_MINUTE=0
```

## Troubleshooting

### Check logs

```bash
# Ingestion service logs
docker logs health-assistant-ingestion-service

# Scheduler service logs
docker logs health-assistant-garmin-scheduler
```

### Common issues

**Authentication failed:**
- Verify your Garmin credentials in `.env`
- Check if your account requires 2FA (not supported)

**No data syncing:**
- Check if you have data in Garmin Connect for the target dates
- Try manual sync: `curl -X POST http://localhost:8085/sync/trigger`
- Check scheduler logs for errors

**Database connection errors:**
- Ensure PostgreSQL is healthy: `docker ps`
- Check database credentials in docker-compose.yml

## Development

### Test ingestion endpoint directly

```bash
curl -X POST http://localhost:8083/api/v1/garmin/ingest/sleep \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "00000000-0000-0000-0000-000000000001",
    "date": "2026-01-28",
    "sleep_data": {
      "sleep_time_seconds": 28800,
      "deep_sleep_seconds": 7200,
      "light_sleep_seconds": 14400,
      "rem_sleep_seconds": 7200,
      "awake_seconds": 0,
      "sleep_scores": {"overall_score": 85}
    }
  }'
```

### Run Python service locally

```bash
cd services/garmin-scheduler

# Create virtual environment
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows

# Install dependencies
pip install -r requirements.txt

# Set environment variables
export GARMIN_EMAIL=your_email@example.com
export GARMIN_PASSWORD=your_password
export DEFAULT_USER_ID=00000000-0000-0000-0000-000000000001
export INGESTION_SERVICE_URL=http://localhost:8083

# Run the service
uvicorn app.main:app --host 0.0.0.0 --port 8085
```

### Run Go service locally

```bash
cd backend

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=healthuser
export DB_PASSWORD=healthpass
export DB_NAME=health_assistant
export DB_SSLMODE=disable

# Run the service
go run cmd/ingestion-service/main.go
```

## Files Reference

### Go Service
- `backend/internal/db/postgres.go` - Database connection pool
- `backend/internal/db/events.go` - Event repository
- `backend/internal/validation/garmin_validator.go` - Payload validation
- `backend/internal/handlers/garmin_ingestion.go` - HTTP handlers
- `backend/cmd/ingestion-service/main.go` - Service entry point

### Python Service
- `services/garmin-scheduler/app/config.py` - Configuration
- `services/garmin-scheduler/app/garmin_client.py` - Garmin API wrapper
- `services/garmin-scheduler/app/ingestion_client.py` - HTTP client to Go
- `services/garmin-scheduler/app/scheduler.py` - APScheduler setup
- `services/garmin-scheduler/app/main.py` - FastAPI application

### Docker
- `backend/Dockerfile.ingestion-service` - Go service Dockerfile
- `services/garmin-scheduler/Dockerfile` - Python service Dockerfile
- `infra/docker-compose.yml` - Container orchestration

## Next Steps

1. Add more data types (respiration, SpO2, body composition)
2. Multi-user support (query users table for credentials)
3. Webhook support (replace polling with push)
4. Retry queues for failed syncs
5. Metrics dashboard
6. Additional wearable adapters (Oura, Apple Health, etc.)
