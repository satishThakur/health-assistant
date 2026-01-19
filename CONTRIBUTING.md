# Contributing to Health Assistant

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

Ensure you have the following installed:
- **Docker & Docker Compose** (for local infrastructure)
- **Go 1.22+** (for backend development)
- **Python 3.11+** (for model service)
- **Flutter SDK** (for mobile app)
- **Git**

### Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/satishThakur/health-assistant.git
   cd health-assistant
   ```

2. **Set up environment variables**
   ```bash
   cd infra
   cp .env.example .env
   # Edit .env with your API keys
   ```

3. **Start infrastructure**
   ```bash
   docker-compose up -d postgres minio
   ```

4. **Verify setup**
   ```bash
   # Check database
   docker-compose exec postgres psql -U healthuser -d health_assistant -c "\dt"

   # Check MinIO
   curl http://localhost:9000/minio/health/live
   ```

## Development Workflow

### Branch Strategy

- `main` - Production-ready code
- `develop` - Integration branch (future)
- `feature/feature-name` - Feature branches
- `fix/bug-description` - Bug fix branches

### Working on a Feature

1. **Create a feature branch**
   ```bash
   git checkout -b feature/garmin-integration
   ```

2. **Make your changes**
   - Follow [CODING_STANDARDS.md](./CODING_STANDARDS.md)
   - Write tests for new functionality
   - Update documentation if needed

3. **Test your changes**
   ```bash
   # Backend (Go)
   cd backend
   go test ./...

   # Model Service (Python)
   cd model-service
   pytest

   # Flutter
   cd app/health_assistant
   flutter test
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: Add Garmin OAuth integration"
   ```
   Follow [commit message conventions](#commit-messages)

5. **Push and create PR** (if working with others)
   ```bash
   git push origin feature/garmin-integration
   ```

### Running Services Locally

**Backend Services**:
```bash
cd backend

# Run specific service
go run ./cmd/api-gateway
go run ./cmd/data-service

# Or build and run
go build -o bin/api-gateway ./cmd/api-gateway
./bin/api-gateway
```

**Model Service**:
```bash
cd model-service
pip install -r requirements.txt
python app/main.py
```

**Flutter App**:
```bash
cd app/health_assistant
flutter pub get
flutter run
```

## Coding Standards

Please read and follow [CODING_STANDARDS.md](./CODING_STANDARDS.md) for:
- Language-specific conventions (Go, Python, Flutter)
- Error handling patterns
- Testing guidelines
- Security best practices

### Key Highlights

**Go**:
- Use `gofmt` or `goimports` for formatting
- Always handle errors explicitly
- Use context for cancellation and timeouts
- Write table-driven tests

**Python**:
- Use type hints everywhere
- Format with `black` and `isort`
- Write docstrings (Google style)
- Use `mypy` for type checking

**Flutter**:
- Follow Dart style guide
- Extract complex widgets into methods/classes
- Use const constructors when possible
- Handle null safety properly

## Commit Messages

Follow conventional commits format:

```
<type>: <subject>

<optional body>

<optional footer>
```

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `refactor` - Code refactoring
- `docs` - Documentation
- `test` - Tests
- `chore` - Maintenance

**Examples**:
```
feat: Add meal photo upload endpoint

Implements S3 upload with presigned URLs.
Supports JPEG and PNG formats up to 10MB.

---

fix: Handle missing HRV values in sleep model

Default to 0 when HRV is not available from Garmin.
Prevents model prediction errors.

Closes #45
```

## Testing Guidelines

### Backend (Go)

**Unit Tests**:
```go
// backend/internal/models/event_test.go
func TestEvent_Validate(t *testing.T) {
    tests := []struct {
        name    string
        event   Event
        wantErr bool
    }{
        {
            name:    "valid event",
            event:   Event{Time: time.Now(), UserID: "test"},
            wantErr: false,
        },
        {
            name:    "missing user ID",
            event:   Event{Time: time.Now()},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.event.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Integration Tests** (with database):
```go
func TestEventStore_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()

    store := NewEventStore(db)
    event := &Event{Time: time.Now(), UserID: "test"}

    err := store.Create(context.Background(), event)
    if err != nil {
        t.Fatalf("Create() failed: %v", err)
    }
}
```

### Model Service (Python)

**Unit Tests**:
```python
# model-service/tests/test_sleep_quality.py
def test_sleep_quality_model_predict():
    model = SleepQualityModel()
    features = {'hrv': 65.0, 'exercise': 45.0}

    prediction, ci = model.predict(features)

    assert 0 <= prediction <= 100
    assert ci[0] < prediction < ci[1]
```

**API Tests**:
```python
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_predict_endpoint():
    response = client.post(
        "/models/sleep-quality/predict",
        json={
            "user_id": "test",
            "features": {"hrv": 65.0}
        }
    )
    assert response.status_code == 200
    assert "prediction" in response.json()
```

### Flutter

**Widget Tests**:
```dart
testWidgets('DailyLogScreen displays feeling sliders', (tester) async {
  await tester.pumpWidget(const MaterialApp(home: DailyLogScreen()));

  expect(find.text('Energy'), findsOneWidget);
  expect(find.text('Mood'), findsOneWidget);
  expect(find.byType(Slider), findsNWidgets(4));
});
```

## Database Migrations

When adding database changes:

1. **Create migration file**
   ```bash
   # Format: YYYYMMDDHHMMSS_description.sql
   touch scripts/db/migrations/20260112120000_add_experiment_tags.sql
   ```

2. **Write migration**
   ```sql
   -- Add experiment tags
   ALTER TABLE experiments ADD COLUMN tags JSONB;
   CREATE INDEX idx_experiments_tags ON experiments USING GIN (tags);
   ```

3. **Test migration**
   ```bash
   docker-compose exec postgres psql -U healthuser -d health_assistant -f /docker-entrypoint-initdb.d/migrations/20260112120000_add_experiment_tags.sql
   ```

4. **Update init.sql** if needed for clean setup

## Documentation

Update documentation when:
- Adding new API endpoints → Update service README
- Changing database schema → Update highleveldesign.md
- Adding new features → Update main README
- Making architectural decisions → Consider adding to docs/

### API Documentation

Document endpoints in code:

**Go**:
```go
// GetEvents retrieves events for a user within a time range
// GET /api/v1/events?user_id={id}&start={ts}&end={ts}
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

**Python**:
```python
@app.post("/models/sleep-quality/predict")
async def predict_sleep_quality(request: PredictionRequest) -> PredictionResponse:
    """
    Predict sleep quality based on input features.

    Args:
        request: Prediction request with user_id and features

    Returns:
        Prediction with confidence interval
    """
```

## Security Guidelines

**Never commit**:
- API keys or secrets
- Passwords
- OAuth tokens
- AWS credentials
- `.env` files

**Always**:
- Use environment variables for secrets
- Validate and sanitize user input
- Use parameterized database queries
- Implement rate limiting on APIs
- Log security-relevant events

**Sensitive Data**:
```go
// GOOD - Use environment variables
apiKey := os.Getenv("GARMIN_API_KEY")

// BAD - Hardcoded secret
apiKey := "abc123xyz"
```

## Code Review Process

When reviewing code, check:

- [ ] Follows coding standards
- [ ] Tests pass and cover new code
- [ ] No secrets committed
- [ ] Error handling is proper
- [ ] Documentation updated
- [ ] Database migrations tested
- [ ] Performance considerations addressed
- [ ] Security implications reviewed

## Getting Help

- **Questions?** Open an issue with `question` label
- **Bug?** Open an issue with `bug` label and reproduction steps
- **Feature idea?** Open an issue with `enhancement` label

## Project Structure Reference

```
health-assistant/
├── backend/              # Go backend services
│   ├── cmd/             # Service binaries
│   ├── internal/        # Private code
│   └── go.mod
├── model-service/       # Python ML service
│   ├── app/
│   └── requirements.txt
├── app/                 # Flutter mobile app
│   └── health_assistant/
├── infra/               # Docker & infrastructure
│   └── docker-compose.yml
├── scripts/             # Helper scripts
│   └── db/             # Database scripts
└── docs/                # Documentation
```

## Useful Commands

**Backend**:
```bash
# Format code
gofmt -w .
goimports -w .

# Lint
golangci-lint run

# Test with coverage
go test -cover ./...

# Build all services
go build ./cmd/...
```

**Model Service**:
```bash
# Format code
black .
isort .

# Type check
mypy app/

# Test
pytest --cov=app
```

**Flutter**:
```bash
# Format
dart format lib/

# Analyze
flutter analyze

# Test
flutter test
```

**Docker**:
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api-gateway

# Rebuild service
docker-compose build api-gateway

# Stop all
docker-compose down
```

---

## License

TBD - Will be determined when project is open sourced

---

Thank you for contributing to Health Assistant! 🚀
