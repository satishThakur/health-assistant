# Coding Standards

This document outlines coding standards and conventions for the Health Assistant project.

## General Principles

1. **Clarity over cleverness** - Code should be easy to understand
2. **Consistency** - Follow established patterns within the codebase
3. **Simplicity** - Avoid over-engineering; solve the problem at hand
4. **Testability** - Write code that's easy to test
5. **Error handling** - Always handle errors explicitly, never silently fail
6. **Documentation** - Document why, not what (code should be self-explanatory)

## Language-Specific Standards

### Go (Backend Services)

#### File Organization
```
backend/
├── cmd/              # Main applications (one per service)
├── internal/         # Private application code
│   ├── api/         # HTTP handlers (grouped by service)
│   ├── db/          # Database layer
│   ├── models/      # Domain models
│   └── ...
└── pkg/             # Public libraries (use sparingly)
```

#### Naming Conventions
- **Packages**: lowercase, single word (e.g., `config`, `models`)
- **Files**: lowercase with underscores (e.g., `event_handler.go`)
- **Interfaces**: suffix with `-er` when possible (e.g., `EventStore`, `DataFetcher`)
- **Structs**: PascalCase (e.g., `Event`, `UserPreferences`)
- **Functions**: camelCase for private, PascalCase for exported
- **Constants**: PascalCase or ALL_CAPS for package-level constants

#### Code Style

**Error Handling**:
```go
// GOOD - Explicit error handling
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// BAD - Silent failure
result, _ := doSomething()
```

**Function Length**:
- Keep functions under 50 lines when possible
- Extract complex logic into helper functions
- One level of abstraction per function

**Comments**:
```go
// GOOD - Explain why
// Cache for 5 minutes to reduce database load during peak hours
cache.Set(key, value, 5*time.Minute)

// BAD - Explain what (code already shows this)
// Set cache with key and value for 5 minutes
cache.Set(key, value, 5*time.Minute)
```

**Struct Tags**:
```go
type Event struct {
    Time       time.Time       `json:"time" db:"time"`
    UserID     string          `json:"user_id" db:"user_id"`
    EventType  string          `json:"event_type" db:"event_type"`
    Data       json.RawMessage `json:"data" db:"data"`
}
```

**Context**:
- Always pass `context.Context` as the first parameter
- Use context for cancellation, deadlines, and request-scoped values
```go
func FetchEvents(ctx context.Context, userID string) ([]Event, error) {
    // ...
}
```

**Database Queries**:
- Use prepared statements or parameterized queries (prevent SQL injection)
- Always use transactions for multi-step operations
- Log slow queries (>100ms)

#### Testing
```go
// Test file naming: <file>_test.go
// Test function naming: TestFunctionName or TestType_Method

func TestEventStore_Create(t *testing.T) {
    // Arrange
    store := NewEventStore(db)
    event := &Event{Time: time.Now(), UserID: "test"}

    // Act
    err := store.Create(context.Background(), event)

    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
}
```

#### Dependencies
- Pin versions in `go.mod`
- Minimize external dependencies
- Prefer standard library when possible

---

### Python (Model Service)

#### File Organization
```
model-service/
├── app/
│   ├── main.py        # FastAPI app
│   ├── models/        # ML models
│   ├── api/           # API routes
│   └── db.py          # Database connection
└── notebooks/         # Jupyter notebooks for exploration
```

#### Naming Conventions
- **Modules**: lowercase with underscores (e.g., `sleep_quality.py`)
- **Classes**: PascalCase (e.g., `SleepQualityModel`)
- **Functions**: snake_case (e.g., `compute_correlations`)
- **Constants**: ALL_CAPS (e.g., `MAX_SAMPLES`)
- **Private**: prefix with underscore (e.g., `_internal_helper`)

#### Code Style

**Type Hints** (Always use):
```python
from typing import Dict, List, Tuple, Optional

def predict_sleep_quality(
    features: Dict[str, float],
    model: SleepQualityModel
) -> Tuple[float, Tuple[float, float]]:
    """
    Predict sleep quality with confidence interval.

    Args:
        features: Dictionary of feature values
        model: Trained sleep quality model

    Returns:
        Tuple of (prediction, (lower_ci, upper_ci))
    """
    prediction, ci = model.predict(features)
    return prediction, ci
```

**Docstrings** (Google Style):
```python
def compute_correlations(
    user_id: str,
    target_metric: str,
    start_date: Optional[str] = None
) -> Dict[str, float]:
    """Compute time-lagged correlations for a target metric.

    Args:
        user_id: User identifier
        target_metric: Metric to analyze (e.g., 'sleep_quality')
        start_date: Optional start date for analysis

    Returns:
        Dictionary mapping feature names to correlation coefficients

    Raises:
        ValueError: If target_metric is not recognized
    """
```

**Error Handling**:
```python
# GOOD - Specific exceptions
try:
    result = process_data(data)
except KeyError as e:
    raise ValueError(f"Missing required field: {e}")
except Exception as e:
    logger.error(f"Unexpected error: {e}")
    raise

# BAD - Bare except
try:
    result = process_data(data)
except:
    pass
```

**FastAPI Routes**:
```python
@app.post("/models/sleep-quality/predict", response_model=PredictionResponse)
async def predict_sleep_quality(request: PredictionRequest) -> PredictionResponse:
    """Predict sleep quality endpoint."""
    try:
        prediction, ci = model.predict(request.features)
        return PredictionResponse(
            prediction=prediction,
            confidence_interval=list(ci),
            uncertainty=ci[1] - ci[0]
        )
    except Exception as e:
        logger.exception("Prediction failed")
        raise HTTPException(status_code=500, detail=str(e))
```

#### Testing
```python
# Test file naming: test_<module>.py
# Use pytest

def test_sleep_quality_model_predict():
    # Arrange
    model = SleepQualityModel()
    features = {'hrv': 65.0, 'exercise': 45.0}

    # Act
    prediction, ci = model.predict(features)

    # Assert
    assert 0 <= prediction <= 100
    assert ci[0] < prediction < ci[1]
```

#### Dependencies
- Use `requirements.txt` for dependencies
- Pin major and minor versions (e.g., `numpy==1.26.3`)
- Keep dependencies up to date for security

---

### Flutter (App)

#### File Organization
```
lib/
├── main.dart
├── models/          # Data models
├── services/        # API clients, business logic
├── screens/         # UI screens
├── widgets/         # Reusable widgets
└── providers/       # State management
```

#### Naming Conventions
- **Files**: lowercase with underscores (e.g., `daily_log_screen.dart`)
- **Classes**: PascalCase (e.g., `DailyLogScreen`, `MealCard`)
- **Functions/Variables**: camelCase (e.g., `fetchEvents`, `userId`)
- **Constants**: lowerCamelCase (e.g., `apiBaseUrl`)
- **Private**: prefix with underscore (e.g., `_buildHeader`)

#### Code Style

**Widget Structure**:
```dart
class DailyLogScreen extends StatelessWidget {
  const DailyLogScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Daily Log')),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    // Extract complex widgets into separate methods
    return Column(
      children: [
        _buildHeader(),
        _buildForm(),
      ],
    );
  }
}
```

**Null Safety**:
```dart
// GOOD - Handle nulls explicitly
String? getUserName(User? user) {
  return user?.name ?? 'Anonymous';
}

// BAD - Force unwrap
String getUserName(User? user) {
  return user!.name;  // Can crash!
}
```

**Async/Await**:
```dart
Future<void> fetchEvents() async {
  try {
    final events = await apiClient.getEvents();
    setState(() {
      _events = events;
    });
  } catch (e) {
    logger.error('Failed to fetch events: $e');
    _showError(context, 'Failed to load events');
  }
}
```

---

## Git Commit Standards

### Commit Message Format
```
<type>: <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring (no functional changes)
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, build, etc.)
- `perf`: Performance improvements

### Examples
```
feat: Add Garmin OAuth integration

Implement OAuth2 flow for Garmin API authentication.
Stores tokens securely in database with encryption.

Closes #12

---

fix: Handle null HRV values in sleep model

Some Garmin devices don't report HRV. Default to 0
when missing to prevent model errors.

---

refactor: Extract event validation to separate function

Reduces complexity in event handler by moving validation
logic to reusable validator.
```

### Commit Guidelines
- Keep subject line under 72 characters
- Use imperative mood ("Add" not "Added")
- Reference issue numbers when applicable
- Include Co-Authored-By for AI assistance

---

## Code Review Checklist

Before committing code, ensure:

- [ ] **No hardcoded secrets** (API keys, passwords)
- [ ] **Error handling** is explicit
- [ ] **Tests pass** (if tests exist)
- [ ] **No commented-out code** (remove or add TODO)
- [ ] **Imports organized** (Go: goimports, Python: isort)
- [ ] **Consistent formatting** (Go: gofmt, Python: black)
- [ ] **Logging** is appropriate (debug vs info vs error)
- [ ] **Documentation** updated if public API changed
- [ ] **No console.log/print** statements in production code
- [ ] **Database queries** use parameterization (no SQL injection)

---

## Security Best Practices

1. **Never commit secrets**
   - Use environment variables
   - Store in `.env` (gitignored)
   - Use secrets manager in production

2. **Input validation**
   - Validate all user input
   - Sanitize data before database insertion
   - Use prepared statements

3. **Authentication**
   - Use JWT tokens with expiration
   - Hash passwords with bcrypt
   - Implement rate limiting

4. **Dependencies**
   - Regularly update dependencies
   - Run security audits (Go: `govulncheck`, Python: `safety`)

---

## Performance Guidelines

1. **Database**
   - Add indexes for frequently queried fields
   - Use connection pooling
   - Implement caching for expensive queries
   - Log slow queries (>100ms)

2. **APIs**
   - Implement pagination for large result sets
   - Use gzip compression
   - Set appropriate cache headers
   - Timeout requests (5-30 seconds)

3. **Frontend**
   - Lazy load images
   - Debounce user input
   - Use pagination for lists
   - Cache API responses

---

## Tools & Automation

### Go
- **Formatting**: `gofmt` or `goimports`
- **Linting**: `golangci-lint`
- **Testing**: `go test -cover ./...`
- **Security**: `govulncheck`

### Python
- **Formatting**: `black`
- **Import sorting**: `isort`
- **Linting**: `pylint` or `ruff`
- **Type checking**: `mypy`
- **Testing**: `pytest`

### Flutter
- **Formatting**: `dart format`
- **Linting**: `flutter analyze`
- **Testing**: `flutter test`

---

## Questions or Clarifications?

When in doubt:
1. Check existing code for patterns
2. Refer to language-specific style guides (Effective Go, PEP 8, Dart Style Guide)
3. Prioritize readability and maintainability
4. Ask for clarification before implementing complex features

Remember: **Code is read more often than it's written.**
