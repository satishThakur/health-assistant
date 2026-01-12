# Backend Services

Single Go module with multiple service binaries for the Personal Health Assistant.

## Architecture

This backend follows the **single Go module, multiple binaries** pattern:
- All services share code from `internal/`
- Each service has its own `cmd/` entry point
- Shared types, database access, auth, and utilities

## Services

| Service | Port | Purpose |
|---------|------|---------|
| **api-gateway** | 8080 | Main API gateway, authentication, routing |
| **data-service** | 8081 | CRUD operations for events and queries |
| **experiment-service** | 8082 | Experiment design, tracking, and management |
| **ingestion-service** | 8083 | Garmin sync, photo processing, external data ingestion |

## Project Structure

```
backend/
├── cmd/                      # Service entry points (binaries)
│   ├── api-gateway/
│   ├── data-service/
│   ├── experiment-service/
│   └── ingestion-service/
│
├── internal/                 # Private application code
│   ├── api/                  # HTTP handlers per service
│   ├── db/                   # Database layer
│   ├── auth/                 # JWT authentication
│   ├── models/               # Domain models
│   ├── garmin/               # Garmin API client
│   ├── llm/                  # LLM integrations
│   └── config/               # Configuration
│
└── pkg/                      # Public libraries (optional)
```

## Running Services

### Development

Run individual services:
```bash
# From backend/ directory
go run ./cmd/api-gateway
go run ./cmd/data-service
go run ./cmd/experiment-service
go run ./cmd/ingestion-service
```

### Building

Build all services:
```bash
go build -o bin/api-gateway ./cmd/api-gateway
go build -o bin/data-service ./cmd/data-service
go build -o bin/experiment-service ./cmd/experiment-service
go build -o bin/ingestion-service ./cmd/ingestion-service
```

Or use the build script:
```bash
cd backend
go build ./cmd/...
```

### Testing

```bash
go test ./...
```

## Configuration

Services are configured via environment variables. See `internal/config/config.go` for all options.

Key variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET`
- `GARMIN_CONSUMER_KEY`, `GARMIN_CONSUMER_SECRET`
- `AWS_REGION`, `AWS_S3_BUCKET`

For local development, create a `.env` file (not committed) or use docker-compose.

## Dependencies

Core dependencies (to be added as needed):
- **Database**: `github.com/jackc/pgx/v5` - PostgreSQL driver
- **Router**: `github.com/go-chi/chi/v5` - HTTP router
- **Auth**: `github.com/golang-jwt/jwt/v5` - JWT tokens
- **AWS SDK**: `github.com/aws/aws-sdk-go-v2` - S3 for photos

## API Documentation

API documentation will be available via Swagger/OpenAPI once endpoints are implemented.

## Next Steps

- [ ] Set up database connection in `internal/db/`
- [ ] Implement JWT authentication in `internal/auth/`
- [ ] Add router and handlers in `internal/api/`
- [ ] Integrate Garmin API in `internal/garmin/`
- [ ] Add LLM integration in `internal/llm/`
