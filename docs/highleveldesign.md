# Personal Health Assistant - High Level Design

## System Architecture

### Overview

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
│            │ │           │ │            │ │              │
│ - CRUD ops │ │ - PyMC    │ │ - Design   │ │ - Garmin API │
│ - Queries  │ │ - Models  │ │ - Track    │ │ - Photo→LLM  │
│ - Aggreg.  │ │ - Predict │ │ - Analyze  │ │ - Parsers    │
└─────┬──────┘ └────┬──────┘ └─────┬──────┘ └───┬──────────┘
      │              │              │             │
      └──────────────┴──────────────┴─────────────┘
                         │
        ┌────────────────┴────────────────┐
        │                                 │
┌───────▼────────────┐          ┌─────────▼─────┐
│   PostgreSQL +     │          │   AWS S3      │
│   TimescaleDB      │          │ (Meal photos) │
│                    │          └───────────────┘
│ - Events           │
│ - Experiments      │
│ - User data        │
└────────────────────┘
```

## Technology Stack

### Frontend
- **Flutter** (Dart)
  - Single codebase for mobile (iOS/Android) and web
  - Material Design 3
  - State management: Riverpod or Bloc
  - HTTP client: Dio with retry logic

### Backend Services (Microservices)

#### API Gateway (Go)
- **Framework**: Chi or Fiber
- **Auth**: JWT tokens (consider OAuth2 for Garmin)
- **Responsibilities**:
  - User authentication
  - Request routing to services
  - Input validation
  - Rate limiting
  - CORS handling

#### Data Service (Go)
- **Purpose**: Primary CRUD operations, queries, data aggregation
- **Database**: `pgx` (PostgreSQL driver)
- **API Style**: RESTful JSON
- **Key Endpoints**:
  - `GET /events` - Query time-series events with filters
  - `POST /events` - Record new event
  - `GET /metrics/summary` - Aggregated metrics for dashboard
  - `GET /timeline` - Unified timeline view

#### Model Service (Python)
- **Framework**: FastAPI
- **Libraries**:
  - PyMC for Bayesian hierarchical models
  - NumPy/Pandas for data manipulation
  - scikit-learn for feature engineering
  - Matplotlib/Plotly for visualizations
- **Deployment**: Docker container
- **API Endpoints**:
  - `POST /models/sleep-quality/predict` - Predict sleep quality
  - `GET /models/correlations` - Compute time-lagged correlations
  - `POST /models/experiment/analyze` - Analyze experiment results
  - `GET /models/insights` - Generate insights from data

#### Experiment Engine (Go)
- **Purpose**: Design, track, and analyze experiments
- **Responsibilities**:
  - Generate experiment proposals based on model insights
  - Track compliance and adherence
  - Coordinate with Model Service for analysis
  - Store results and update priors
- **Key Logic**:
  - Factorial design generation
  - Randomization schedules
  - Statistical power calculations

#### Ingestion Service (Go)
- **Purpose**: Pull data from external sources
- **Components**:
  - Garmin API client (OAuth2 flow, hourly polling)
  - LLM integration for meal photo analysis (OpenAI/Anthropic API)
  - PDF parser for lab results (using LLM)
- **Scheduling**: Cron jobs or task queue (consider Asynq)
- **Error Handling**: Retry with exponential backoff

### Data Layer

#### Database: PostgreSQL + TimescaleDB Extension

**Why TimescaleDB?**
- Optimized for time-series queries
- Automatic partitioning (hypertables)
- Continuous aggregates for fast rollups
- Retention policies for old data

**Schema Design** (see below)

#### Object Storage: AWS S3
- Meal photos (original + thumbnails)
- Lab result PDFs
- Model artifacts (trained models, plots)
- Bucket structure: `/photos/{user_id}/{date}/{uuid}.jpg`

## Data Model

### Events Table (TimescaleDB Hypertable)

```sql
CREATE TABLE events (
    time TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    source VARCHAR(50) NOT NULL,  -- garmin, manual, parsed, etc.
    data JSONB NOT NULL,
    metadata JSONB,
    confidence FLOAT CHECK (confidence >= 0 AND confidence <= 1),
    PRIMARY KEY (time, user_id, event_type)
);

SELECT create_hypertable('events', 'time');
CREATE INDEX ON events (user_id, time DESC);
CREATE INDEX ON events (event_type, time DESC);
CREATE INDEX ON events USING GIN (data);
```

**Event Types & Data Schema**:

```typescript
// Garmin metrics
event_type: "garmin_sleep"
data: {
  duration_minutes: 420,
  deep_sleep_minutes: 90,
  light_sleep_minutes: 240,
  rem_sleep_minutes: 90,
  awake_minutes: 0,
  sleep_score: 85,
  hrv_avg: 65
}

event_type: "garmin_activity"
data: {
  activity_type: "strength_training",
  duration_minutes: 45,
  calories: 250,
  avg_hr: 135,
  max_hr: 165
}

// Subjective
event_type: "subjective_feeling"
data: {
  energy: 7,      // 1-10 scale
  mood: 8,
  focus: 6,
  physical: 7,
  notes: "Felt great after morning walk"
}

// Nutrition
event_type: "meal"
data: {
  meal_type: "dinner",
  photo_url: "s3://...",
  macros: {
    calories: 650,
    protein_g: 45,
    carbs_g: 60,
    fat_g: 25,
    fiber_g: 12
  },
  confidence: 0.75,  // LLM confidence
  manually_verified: false
}

// Supplements
event_type: "supplement"
data: {
  name: "Creatine Monohydrate",
  dosage: "5g",
  taken: true,
  scheduled_time: "08:00",
  actual_time: "08:15"
}

// Lab results
event_type: "biomarker"
data: {
  test_name: "Vitamin D, 25-Hydroxy",
  value: 45.2,
  unit: "ng/mL",
  reference_range: "30-100",
  lab_name: "Quest Diagnostics"
}
```

### Experiments Table

```sql
CREATE TABLE experiments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    hypothesis TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('proposed', 'accepted', 'active', 'completed', 'abandoned')),

    -- Design
    intervention JSONB NOT NULL,  -- What we're testing
    control_condition JSONB,
    duration_days INT NOT NULL,
    start_date DATE,
    end_date DATE,

    -- Tracking
    compliance_rate FLOAT,

    -- Results
    results JSONB,  -- Statistical outcomes from Model Service
    posterior_beliefs JSONB,  -- Updated Bayesian priors

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON experiments (user_id, status);
CREATE INDEX ON experiments (start_date, end_date);
```

**Example Experiment**:
```json
{
  "name": "Creatine Effect on Recovery",
  "hypothesis": "5g daily creatine improves HRV recovery and reduces muscle soreness",
  "intervention": {
    "supplement": "creatine_monohydrate",
    "dosage": "5g",
    "timing": "post_workout"
  },
  "control_condition": {
    "supplement": "none"
  },
  "duration_days": 28,
  "results": {
    "hrv_effect": {
      "mean": 3.2,
      "std": 1.5,
      "credible_interval_95": [0.5, 5.9],
      "probability_positive": 0.94
    }
  }
}
```

### Users Table (Future)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    garmin_oauth_token JSONB,
    preferences JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Modeling Approach

### Time-Series Hierarchical Bayesian Models

**Core Philosophy**: Model individual as hierarchy, not flat regression

#### Example: Sleep Quality Model

```python
# Conceptual PyMC model structure
with pm.Model() as sleep_model:
    # Priors (population level)
    μ_sleep = pm.Normal('mu_sleep', mu=80, sigma=10)

    # Individual-level effects
    β_hrv = pm.Normal('beta_hrv', mu=0, sigma=1)
    β_exercise = pm.Normal('beta_exercise', mu=0, sigma=1)
    β_meal_timing = pm.Normal('beta_meal_timing', mu=0, sigma=1)
    β_supplement = pm.Normal('beta_supplement', mu=0, sigma=1)

    # Autoregressive component (yesterday affects today)
    ρ = pm.Beta('rho', alpha=2, beta=2)

    # Likelihood
    sleep_quality = (μ_sleep +
                     β_hrv * hrv_data +
                     β_exercise * exercise_data +
                     β_meal_timing * meal_timing_data +
                     β_supplement * supplement_data +
                     ρ * sleep_quality_lag1 +
                     pm.Normal('noise', mu=0, sigma=5))

    # Observed data
    pm.Normal('obs', mu=sleep_quality, sigma=3, observed=observed_sleep)
```

**Key Features**:
- Time-lagged effects (X yesterday → Y today)
- Hierarchical structure for regularization
- Posterior distributions (not point estimates)
- Credible intervals on all parameters

### Causal Inference Strategy

1. **Observational Phase**:
   - Build correlational models
   - Identify candidate causal factors
   - Use time-precedence (cause must precede effect)

2. **Experimental Phase**:
   - Randomized n-of-1 trials
   - A/B or factorial designs
   - Washout periods to avoid carryover

3. **Bayesian Updating**:
   - Start with weak priors
   - Update beliefs as experiments complete
   - Track uncertainty reduction

## API Design

### RESTful Endpoints (Summary)

#### Events
- `GET /api/v1/events?start=<timestamp>&end=<timestamp>&type=<type>`
- `POST /api/v1/events`
- `GET /api/v1/events/summary?period=week|month`

#### Models
- `GET /api/v1/models/insights`
- `POST /api/v1/models/predict`
- `GET /api/v1/models/correlations?target=sleep_quality`

#### Experiments
- `GET /api/v1/experiments`
- `POST /api/v1/experiments` (user accepts proposal)
- `PATCH /api/v1/experiments/:id/compliance`
- `GET /api/v1/experiments/:id/results`

#### Ingestion
- `POST /api/v1/ingest/garmin/sync` (manual trigger)
- `POST /api/v1/ingest/meal` (upload photo)
- `POST /api/v1/ingest/lab-result` (upload PDF)

## Deployment Architecture

### Local Development
- Docker Compose with:
  - PostgreSQL + TimescaleDB
  - Go services (live reload with Air)
  - Python service (uvicorn with reload)
  - LocalStack (S3 emulation) or MinIO

### Production (AWS)
- **Compute**: ECS Fargate (containers)
- **Database**: RDS PostgreSQL with TimescaleDB
- **Storage**: S3 (photos, PDFs)
- **Secrets**: AWS Secrets Manager
- **Monitoring**: CloudWatch + Grafana
- **CI/CD**: GitHub Actions → ECR → ECS

## Security Considerations

1. **Authentication**:
   - JWT tokens for API
   - Refresh token rotation
   - Garmin OAuth2 tokens encrypted at rest

2. **Data Privacy**:
   - All PII encrypted in database
   - Photos encrypted in S3 (SSE)
   - No third-party analytics

3. **API Security**:
   - Rate limiting per user
   - Input validation on all endpoints
   - SQL injection prevention (parameterized queries)

## Scalability Notes

**Current Scope**: Single user (MVP)

**Future Considerations**:
- Multi-tenancy: Partition by user_id
- Model caching: Store predictions in Redis
- Background jobs: Task queue (Asynq, Temporal)
- Read replicas for analytics queries

## Open Questions / Decisions Needed

1. **Authentication**: Build custom or use service (Auth0, Supabase)?
2. **LLM Provider**: OpenAI (GPT-4V) or Anthropic (Claude) for food analysis?
3. **Monitoring**: Self-hosted (Prometheus/Grafana) or managed (Datadog)?
4. **Flutter state management**: Riverpod or Bloc?
5. **Model versioning**: How to track model evolution over time?

## Success Metrics (Technical)

- **Data Pipeline**: 99% successful Garmin syncs (hourly)
- **LLM Accuracy**: >80% meal macro estimation within 20% of ground truth
- **API Latency**: p95 < 200ms for all GET requests
- **Model Inference**: Sleep quality prediction in <5 seconds
- **Uptime**: 99.5% availability (MVP acceptable)
