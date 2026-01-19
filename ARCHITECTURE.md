# Architecture Overview

This document explains the architectural decisions and patterns used in the Health Assistant project.

## Table of Contents
- [System Architecture](#system-architecture)
- [Design Decisions](#design-decisions)
- [Data Flow](#data-flow)
- [Service Communication](#service-communication)
- [Database Design](#database-design)
- [Security Architecture](#security-architecture)
- [Scalability Considerations](#scalability-considerations)

---

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Flutter Application                       │
│              (Mobile + Web, Single Codebase)                 │
└────────────────────────┬────────────────────────────────────┘
                         │ REST API (JSON)
┌────────────────────────▼────────────────────────────────────┐
│                   API Gateway (Go)                           │
│          - Authentication (JWT)                              │
│          - Request routing & validation                      │
│          - Rate limiting                                     │
└─────┬──────────────┬──────────────┬─────────────┬──────────┘
      │              │              │             │
┌─────▼──────┐ ┌────▼──────┐ ┌─────▼──────┐ ┌───▼──────────┐
│   Data     │ │  Model    │ │ Experiment │ │  Ingestion   │
│  Service   │ │  Service  │ │  Engine    │ │  Service     │
│   (Go)     │ │  (Python) │ │   (Go)     │ │   (Go)       │
└─────┬──────┘ └────┬──────┘ └─────┬──────┘ └───┬──────────┘
      │              │              │             │
      └──────────────┴──────────────┴─────────────┘
                         │
        ┌────────────────┴────────────────┐
        │                                 │
┌───────▼────────────┐          ┌─────────▼─────┐
│   PostgreSQL +     │          │   AWS S3      │
│   TimescaleDB      │          │ (Meal photos) │
└────────────────────┘          └───────────────┘
```

### Component Responsibilities

| Component | Purpose | Technology | Port |
|-----------|---------|------------|------|
| **API Gateway** | Entry point, auth, routing | Go, Chi/Fiber | 8080 |
| **Data Service** | CRUD operations, queries | Go, pgx | 8081 |
| **Experiment Service** | Experiment design & tracking | Go | 8082 |
| **Ingestion Service** | External data ingestion | Go | 8083 |
| **Model Service** | ML models, predictions | Python, FastAPI | 8084 |
| **PostgreSQL/TimescaleDB** | Time-series data store | PostgreSQL 16 | 5432 |
| **MinIO/S3** | Object storage (photos, PDFs) | S3-compatible | 9000 |

---

## Design Decisions

### 1. Microservices vs Monolith

**Decision**: Microservices architecture with separate service binaries

**Rationale**:
- **Separation of concerns**: Each service has a clear, focused responsibility
- **Independent scaling**: Can scale model service separately (CPU-intensive)
- **Technology diversity**: Python for ML, Go for high-performance APIs
- **Fault isolation**: Failure in one service doesn't bring down entire system

**Trade-offs**:
- More complex deployment
- Network overhead between services
- Harder to debug distributed issues

**Why it's acceptable**:
- Docker Compose simplifies local development
- Small team (solo) can manage with proper tooling
- Future cloud deployment benefits from this structure

### 2. Single Go Module for Backend

**Decision**: One Go module with multiple `cmd/` binaries instead of separate modules per service

**Rationale**:
- **Code sharing**: All services share `internal/` packages (db, models, auth)
- **Dependency management**: Single `go.mod` easier to maintain
- **Refactoring**: Easy to move code between services
- **Build efficiency**: Go can optimize across the entire codebase

**Pattern**:
```
backend/
├── go.mod                    # Single module
├── cmd/
│   ├── api-gateway/main.go   # Binary 1
│   ├── data-service/main.go  # Binary 2
│   └── ...
└── internal/                 # Shared code
    ├── db/
    ├── models/
    └── auth/
```

### 3. TimescaleDB for Time-Series Data

**Decision**: PostgreSQL with TimescaleDB extension

**Rationale**:
- **Time-series optimization**: Automatic partitioning, fast queries on time ranges
- **PostgreSQL compatibility**: Standard SQL, ACID guarantees, rich ecosystem
- **Continuous aggregates**: Pre-computed rollups (e.g., daily metrics)
- **Retention policies**: Automatic data cleanup

**Alternatives considered**:
- InfluxDB: More specialized but less flexible for relational queries
- Raw PostgreSQL: Works but lacks time-series optimizations
- MongoDB: NoSQL flexibility but weaker for time-based analytics

### 4. Bayesian Approach for Models

**Decision**: PyMC for Bayesian hierarchical models

**Rationale**:
- **Uncertainty quantification**: Credible intervals, not just point estimates
- **Small sample size**: Works well with limited n=1 data
- **Prior knowledge**: Can incorporate domain knowledge via priors
- **Interpretability**: Posterior distributions are human-readable

**Trade-offs**:
- Slower inference (MCMC sampling)
- More complex than simple regression
- Requires understanding of Bayesian statistics

**Why it's worth it**:
- N=1 experimentation demands proper uncertainty handling
- False certainty is worse than admitting uncertainty
- Foundation for causal inference

### 5. Flutter for Cross-Platform App

**Decision**: Flutter instead of native iOS/Android or React Native

**Rationale**:
- **Single codebase**: iOS, Android, and Web from one codebase
- **Performance**: Near-native performance with compiled Dart
- **UI consistency**: Material Design 3 widgets out of the box
- **GenAI acceleration**: Good ecosystem, LLMs trained on Flutter code

**Alternatives considered**:
- Native: Best performance but 2-3x development time
- React Native: Larger ecosystem but JavaScript pain points
- PWA: Easier but limited mobile capabilities

---

## Data Flow

### 1. Event Ingestion Flow

```
┌───────────┐
│  Garmin   │
│    API    │
└─────┬─────┘
      │ OAuth + Poll
┌─────▼──────────┐
│   Ingestion    │ ────┐
│   Service      │     │ Parse & Normalize
└────────────────┘     │
                       ▼
┌─────────────────────────┐
│   PostgreSQL/Events     │
│   (TimescaleDB)         │
└───────┬─────────────────┘
        │
        │ Query
        ▼
┌────────────────┐
│  Model Service │ ──► Predictions, Insights
└────────────────┘
```

### 2. Meal Photo Flow

```
┌───────────┐
│  Flutter  │
│    App    │
└─────┬─────┘
      │ 1. Upload photo
┌─────▼──────────┐
│   Ingestion    │
│   Service      │
└─────┬──────────┘
      │ 2. Store in S3
      ▼
┌─────────────┐
│   MinIO/S3  │
└─────────────┘
      │ 3. Send to LLM
      ▼
┌──────────────┐
│  Claude/GPT  │ ──► Extract macros
└──────┬───────┘
       │ 4. Save event
       ▼
┌─────────────────┐
│  PostgreSQL     │
│  (meal event)   │
└─────────────────┘
```

### 3. Experiment Lifecycle

```
1. Model Service analyzes data ─► Proposes experiment
                                          │
2. User reviews in Flutter app ───────────┘
                                          │
3. User accepts ──────────────────────────┤
                                          │
4. Experiment Service tracks compliance ──┤
                                          │
5. User logs daily (supplements, etc.) ───┤
                                          │
6. Experiment completes ──────────────────┤
                                          │
7. Model Service analyzes results ────────┤
                                          │
8. Update posterior beliefs ──────────────┘
```

---

## Service Communication

### Inter-Service Communication

**Pattern**: HTTP REST APIs (synchronous)

**Why not gRPC?**
- HTTP/JSON simpler for small team
- Easier debugging (curl, browser)
- Good enough for MVP scale
- Can migrate to gRPC later if needed

**Example**:
```go
// API Gateway calls Data Service
resp, err := http.Get("http://data-service:8081/events?user_id=123")
```

**Future consideration**: Service mesh (Istio, Linkerd) if complexity grows

### API Gateway Pattern

The API Gateway:
1. **Authenticates** requests (JWT validation)
2. **Routes** to appropriate service
3. **Aggregates** responses if needed
4. **Transforms** data formats
5. **Handles** CORS and rate limiting

**Benefits**:
- Single entry point for clients
- Centralized auth and logging
- Can cache responses
- Shields internal service changes

---

## Database Design

### Schema Philosophy

**Principles**:
1. **Flexible data field**: Use JSONB for varying event structures
2. **Strong types for queries**: Typed columns for filtering (time, user_id, event_type)
3. **Indexes for performance**: Index all query patterns
4. **Time-series optimization**: TimescaleDB hypertables

### Events Table Design

```sql
CREATE TABLE events (
    time TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    source VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,           -- Flexible schema
    metadata JSONB,
    confidence FLOAT,
    PRIMARY KEY (time, user_id, event_type)
);

SELECT create_hypertable('events', 'time');
```

**Why JSONB for `data`?**
- Different event types have different schemas
- Avoids table proliferation (one events table, not 10+)
- PostgreSQL JSONB is fast and indexable
- Flexibility for schema evolution

**Trade-offs**:
- Less type safety at database level
- Complex queries require JSON operators
- Need validation in application layer

### Indexing Strategy

```sql
-- Time-range queries (most common)
CREATE INDEX idx_events_user_time ON events (user_id, time DESC);

-- Event type filtering
CREATE INDEX idx_events_type_time ON events (event_type, time DESC);

-- JSON field queries
CREATE INDEX idx_events_data_gin ON events USING GIN (data);
```

---

## Security Architecture

### Authentication Flow

```
1. User logs in ─► API Gateway validates credentials
                           │
2. Generate JWT token ─────┤
                           │
3. Return token ───────────┘
                           │
4. Client stores token ────┤
                           │
5. Include in requests ────┤
   (Authorization: Bearer <token>)
                           │
6. API Gateway validates ──┘
   - Signature check
   - Expiration check
   - User exists
```

### Data Security

**At Rest**:
- Database encryption (AWS RDS encryption)
- S3 server-side encryption (SSE)
- Secrets in AWS Secrets Manager

**In Transit**:
- HTTPS/TLS for all external APIs
- Internal services: HTTPS in production, HTTP in local dev

**Application Level**:
- Parameterized queries (SQL injection prevention)
- Input validation on all endpoints
- Rate limiting to prevent abuse
- CORS configured properly

### Secrets Management

**Development**:
```bash
# .env file (not committed)
GARMIN_API_KEY=abc123
JWT_SECRET=dev-secret
```

**Production**:
- AWS Secrets Manager
- Environment variables in ECS
- No secrets in code or config files

---

## Scalability Considerations

### Current Scale (MVP)

- **Users**: 1 (personal use)
- **Events**: ~1000/day
- **Database**: <10GB
- **Compute**: Single EC2 or ECS task per service

### Future Scaling

**Horizontal Scaling**:
```
Load Balancer
      │
      ├── API Gateway (×3)
      ├── Data Service (×2)
      ├── Model Service (×5)  ← CPU-intensive
      └── Ingestion Service (×2)
```

**Database Scaling**:
1. **Read replicas** for analytics queries
2. **Connection pooling** (pgBouncer)
3. **Caching layer** (Redis for hot data)
4. **Partitioning** (TimescaleDB handles this)

**Object Storage**:
- S3 scales automatically
- CloudFront CDN for photo delivery

### Bottlenecks to Watch

1. **Model inference** - Python/PyMC is slow
   - Solution: Cache predictions, async job queue
2. **Database writes** - High event volume
   - Solution: Batch writes, connection pooling
3. **API Gateway** - Single point of failure
   - Solution: Multiple instances behind load balancer

---

## Monitoring & Observability

### Logging Strategy

**Levels**:
- **DEBUG**: Development only
- **INFO**: Normal operations (request logs)
- **WARN**: Recoverable errors
- **ERROR**: Failures requiring attention

**Structured Logging**:
```go
log.Info("event created",
    "user_id", userID,
    "event_type", eventType,
    "duration_ms", duration)
```

### Metrics to Track

**Application**:
- Request rate (req/sec)
- Error rate (4xx, 5xx)
- Response time (p50, p95, p99)
- Database query time

**Business**:
- Events ingested per day
- Experiment completion rate
- Model prediction accuracy

**Infrastructure**:
- CPU/Memory usage
- Database connections
- Disk space

### Tools

**Development**:
- Docker Compose logs
- PostgreSQL logs

**Production**:
- CloudWatch (AWS)
- Grafana dashboards
- Alerting on errors

---

## Technology Choices Summary

| Aspect | Technology | Why |
|--------|-----------|-----|
| **Backend** | Go | Fast, concurrent, good for APIs |
| **ML/Models** | Python + PyMC | Bayesian inference, ML ecosystem |
| **Frontend** | Flutter | Cross-platform, single codebase |
| **Database** | PostgreSQL + TimescaleDB | Time-series + relational |
| **Storage** | S3 | Scalable, cheap, durable |
| **Auth** | JWT | Stateless, standard |
| **Deployment** | Docker + AWS ECS | Containerized, managed |

---

## Future Architecture Evolution

### Phase 1 (MVP): Current Architecture
- Monolithic deployment (all services together)
- Single database instance
- Local development with Docker Compose

### Phase 2 (Multi-User)
- Separate deployments per service
- Read replicas for database
- Redis caching layer
- CDN for static assets

### Phase 3 (Scale)
- Kubernetes orchestration
- Message queue (RabbitMQ/SQS) for async tasks
- Service mesh (Istio)
- Distributed tracing (Jaeger)

---

## Open Questions

1. **Should we add GraphQL?**
   - Pros: Flexible queries from frontend
   - Cons: Added complexity, not needed for MVP

2. **Event sourcing pattern?**
   - Pros: Complete audit trail, replay capability
   - Cons: Complexity, storage overhead

3. **Real-time updates?**
   - WebSockets for live experiment updates?
   - Or poll-based for simplicity?

---

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [TimescaleDB Best Practices](https://docs.timescale.com/use-timescale/latest/best-practices/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [Flutter Architecture](https://docs.flutter.dev/development/data-and-backend/state-mgmt/intro)

---

**Last Updated**: January 2026
**Status**: Architecture defined, implementation in progress
