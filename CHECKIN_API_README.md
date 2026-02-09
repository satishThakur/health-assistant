# Daily Check-in API Documentation

## Overview

The Daily Check-in API allows users to submit their subjective feelings (energy, mood, focus, physical state) and view correlations with their Garmin health data.

## Base URL

```
http://localhost:8083  (local development)
https://api.your-domain.com  (production)
```

## Authentication

**Note:** Authentication is not yet implemented. Currently using a default user_id.

**TODO:** Add JWT Bearer token authentication:
```
Authorization: Bearer <your_jwt_token>
```

---

## Endpoints

### 1. Submit Daily Check-in

**POST** `/api/v1/checkin`

Submit your daily subjective feelings assessment.

**Request Body:**
```json
{
  "energy": 8,
  "mood": 7,
  "focus": 9,
  "physical": 7,
  "notes": "Felt great after morning run"
}
```

**Field Validation:**
- `energy` (required): Integer 1-10
- `mood` (required): Integer 1-10
- `focus` (required): Integer 1-10
- `physical` (required): Integer 1-10
- `notes` (optional): String, max 1000 characters

**Success Response (200 OK):**
```json
{
  "status": "success",
  "action": "inserted",
  "timestamp": "2026-02-09T14:30:00Z",
  "data": {
    "energy": 8,
    "mood": 7,
    "focus": 9,
    "physical": 7,
    "notes": "Felt great after morning run"
  }
}
```

**Note:** Only one check-in per day is allowed. Submitting again on the same day will update the existing entry (action = "updated").

**Error Response (400 Bad Request):**
```json
{
  "error": "Validation failed",
  "message": "energy must be between 1 and 10, got 11"
}
```

**Example:**
```bash
curl -X POST http://localhost:8083/api/v1/checkin \
  -H "Content-Type: application/json" \
  -d '{
    "energy": 8,
    "mood": 7,
    "focus": 9,
    "physical": 7,
    "notes": "Felt great today"
  }'
```

---

### 2. Get Today's Check-in

**GET** `/api/v1/checkin/latest`

Retrieve today's check-in if it exists.

**Success Response (200 OK):**
```json
{
  "status": "success",
  "timestamp": "2026-02-09T00:00:00Z",
  "checkin": {
    "energy": 8,
    "mood": 7,
    "focus": 9,
    "physical": 7,
    "notes": "Felt great after morning run"
  }
}
```

**No Check-in Today:**
```json
{
  "status": "success",
  "checkin": null,
  "message": "No check-in for today"
}
```

**Example:**
```bash
curl http://localhost:8083/api/v1/checkin/latest
```

---

### 3. Get Check-in History

**GET** `/api/v1/checkin/history?days=30`

Retrieve historical check-ins.

**Query Parameters:**
- `days` (optional): Number of days to retrieve (default: 30)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "count": 7,
  "history": [
    {
      "date": "2026-02-09",
      "checkin": {
        "energy": 8,
        "mood": 7,
        "focus": 9,
        "physical": 7,
        "notes": "Felt great today"
      }
    },
    {
      "date": "2026-02-08",
      "checkin": {
        "energy": 6,
        "mood": 5,
        "focus": 7,
        "physical": 6,
        "notes": "Tired after late night"
      }
    }
  ]
}
```

**Example:**
```bash
curl "http://localhost:8083/api/v1/checkin/history?days=7"
```

---

### 4. Get Today's Dashboard

**GET** `/api/v1/dashboard/today`

Retrieve today's comprehensive summary including check-in and Garmin data.

**Success Response (200 OK):**
```json
{
  "status": "success",
  "data": {
    "checkin": {
      "energy": 8,
      "mood": 7,
      "focus": 9,
      "physical": 7,
      "notes": "Felt great after morning run"
    },
    "garmin": {
      "sleep": {
        "duration_minutes": 432,
        "deep_sleep_minutes": 126,
        "light_sleep_minutes": 228,
        "rem_sleep_minutes": 78,
        "awake_minutes": 0,
        "sleep_score": 82,
        "hrv_avg": 67.5
      },
      "activity": {
        "activity_type": "running",
        "duration_minutes": 45,
        "calories": 285,
        "avg_hr": 132,
        "max_hr": 168,
        "distance": 5000
      },
      "hrv": {
        "average": 67.5
      },
      "stress": {
        "average": 32,
        "level": "low"
      }
    }
  }
}
```

**Response Notes:**
- If no check-in for today, `checkin` will be `null`
- If no Garmin data synced today, respective fields will be `null`
- Stress levels: "low" (0-25), "moderate" (26-50), "high" (51+)

**Example:**
```bash
curl http://localhost:8083/api/v1/dashboard/today
```

---

### 5. Get 7-Day Trends

**GET** `/api/v1/trends/week`

Retrieve 7-day trend data for charts and visualizations.

**Success Response (200 OK):**
```json
{
  "status": "success",
  "count": 7,
  "trends": [
    {
      "date": "2026-02-03",
      "checkin": {
        "energy": 7,
        "mood": 6,
        "focus": 8,
        "physical": 7
      },
      "sleep": {
        "duration_minutes": 420,
        "sleep_score": 78
      },
      "activity": {
        "duration_minutes": 30,
        "calories": 200
      }
    },
    {
      "date": "2026-02-04",
      "checkin": {
        "energy": 8,
        "mood": 8,
        "focus": 9,
        "physical": 8
      },
      "sleep": {
        "duration_minutes": 450,
        "sleep_score": 85
      },
      "activity": {
        "duration_minutes": 45,
        "calories": 285
      }
    }
  ]
}
```

**Use Cases:**
- Plot energy/mood trends over the week
- Compare sleep duration with feelings
- Identify best/worst days

**Example:**
```bash
curl http://localhost:8083/api/v1/trends/week
```

---

### 6. Get Correlations & Insights

**GET** `/api/v1/insights/correlations?days=30`

Calculate correlations between Garmin data and subjective feelings.

**Query Parameters:**
- `days` (optional): Number of days to analyze (default: 30)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "count": 3,
  "correlations": [
    {
      "type": "sleep_energy",
      "description": "Your energy is 15% higher when you sleep 7+ hours",
      "confidence": 0.85,
      "sample_size": 25,
      "details": {
        "condition": "sleep >= 7 hours",
        "avg_energy_with": 8.2,
        "avg_energy_without": 7.1,
        "improvement_percent": 15.5
      }
    },
    {
      "type": "activity_mood",
      "description": "Your mood improves by 12% on active days (30+ min)",
      "confidence": 0.78,
      "sample_size": 22,
      "details": {
        "condition": "activity >= 30 minutes",
        "avg_mood_with": 7.8,
        "avg_mood_without": 6.9,
        "improvement_percent": 13.0
      }
    },
    {
      "type": "sleep_focus",
      "description": "Your focus is 10% better after quality sleep (score 80+)",
      "confidence": 0.82,
      "sample_size": 20,
      "details": {
        "condition": "sleep_score >= 80",
        "avg_focus_with": 8.5,
        "avg_focus_without": 7.7,
        "improvement_percent": 10.4
      }
    }
  ]
}
```

**Correlation Types:**
- `sleep_energy`: Sleep duration vs energy levels
- `activity_mood`: Physical activity vs mood
- `sleep_focus`: Sleep quality vs focus/concentration

**Requirements:**
- At least 5 samples in each comparison group
- Improvement must be at least 5% to show insight
- Minimum 30 days of data recommended for reliable insights

**Example:**
```bash
curl "http://localhost:8083/api/v1/insights/correlations?days=30"
```

---

## Data Model

### SubjectiveFeeling

Stored in the `events` table with:
- `event_type`: "subjective_feeling"
- `source`: "manual"
- `time`: Start of the day (00:00:00)

```json
{
  "energy": 8,
  "mood": 7,
  "focus": 9,
  "physical": 7,
  "notes": "Optional notes"
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request (validation failed) |
| 405 | Method Not Allowed |
| 500 | Internal Server Error |

---

## Testing

### Quick Test

Run all endpoints:
```bash
./scripts/test-checkin-api.sh
```

### Manual Testing

1. **Submit a check-in:**
```bash
curl -X POST http://localhost:8083/api/v1/checkin \
  -H "Content-Type: application/json" \
  -d '{"energy":8,"mood":7,"focus":9,"physical":7}'
```

2. **View dashboard:**
```bash
curl http://localhost:8083/api/v1/dashboard/today | jq .
```

3. **Check trends:**
```bash
curl http://localhost:8083/api/v1/trends/week | jq .
```

4. **Get insights:**
```bash
curl "http://localhost:8083/api/v1/insights/correlations?days=30" | jq .
```

---

## Integration with Flutter App

### Example API Client (Dart)

```dart
class CheckinApi {
  final Dio dio;
  final String baseUrl;

  CheckinApi({required this.dio, required this.baseUrl});

  Future<void> submitCheckin({
    required int energy,
    required int mood,
    required int focus,
    required int physical,
    String? notes,
  }) async {
    await dio.post(
      '$baseUrl/api/v1/checkin',
      data: {
        'energy': energy,
        'mood': mood,
        'focus': focus,
        'physical': physical,
        if (notes != null) 'notes': notes,
      },
    );
  }

  Future<DashboardData> getTodayDashboard() async {
    final response = await dio.get('$baseUrl/api/v1/dashboard/today');
    return DashboardData.fromJson(response.data['data']);
  }

  Future<List<TrendData>> getWeekTrends() async {
    final response = await dio.get('$baseUrl/api/v1/trends/week');
    return (response.data['trends'] as List)
        .map((e) => TrendData.fromJson(e))
        .toList();
  }

  Future<List<Correlation>> getCorrelations({int days = 30}) async {
    final response = await dio.get(
      '$baseUrl/api/v1/insights/correlations',
      queryParameters: {'days': days},
    );
    return (response.data['correlations'] as List)
        .map((e) => Correlation.fromJson(e))
        .toList();
  }
}
```

---

## Development

### Run Service Locally

```bash
# Start database
cd infra
docker-compose up postgres -d

# Run service
cd backend
go run cmd/ingestion-service/main.go
```

Service will be available at `http://localhost:8083`

### Run Tests

```bash
# Unit tests
cd backend
go test ./internal/validation/...

# Integration test
./scripts/test-checkin-api.sh
```

---

## Future Enhancements

1. **Authentication**
   - JWT token validation
   - User-specific data isolation

2. **Advanced Correlations**
   - Multi-variate analysis
   - Machine learning predictions
   - Personalized recommendations

3. **Notifications**
   - Daily reminder to check-in
   - Insights notifications
   - Streak tracking

4. **Export**
   - CSV export of check-in history
   - PDF reports with visualizations

5. **Social Features**
   - Compare with anonymized community data
   - Share insights

---

## Support

For issues or questions:
- Check logs: `docker logs ingestion-service`
- Run health check: `curl http://localhost:8083/health`
- View database: `psql -h localhost -U healthuser -d health_assistant`

---

## License

MIT License - See LICENSE file for details
