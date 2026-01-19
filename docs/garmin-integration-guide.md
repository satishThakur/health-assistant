# Garmin Health API Integration Guide

⚠️ **IMPORTANT LIMITATION**: Garmin Health API is typically **only approved for commercial entities, healthcare organizations, and research institutions** - NOT individual/personal developers. For practical alternatives, see [wearable-data-sources.md](./wearable-data-sources.md).

This guide explains how to integrate with Garmin Health API to pull health and fitness data **if you have commercial entity approval**.

## Overview

The Garmin Health API provides access to comprehensive health and fitness data from Garmin wearables. For this project, we'll focus on data that supports our causal health modeling.

---

## What Data Can We Get?

### Core Metrics (What We Need)

| Data Type | Endpoint | Description | Use Case |
|-----------|----------|-------------|----------|
| **Sleep** | `/wellness/daily/{summaryDate}/sleep` | Duration, stages (deep, light, REM), quality score, HRV | Primary outcome for sleep quality model |
| **HRV** | `/wellness/hrv` | Heart rate variability during sleep | Key predictor for recovery and stress |
| **Activity** | `/wellness/daily/{summaryDate}/activities` | Workouts, duration, intensity, calories | Exercise impact on sleep/recovery |
| **Stress** | `/wellness/daily/{summaryDate}/stress` | All-day stress scores | Stress impact on health outcomes |
| **Body Battery** | (Included in daily summary) | Garmin's energy level metric | Recovery indicator |
| **Heart Rate** | `/wellness/daily/{summaryDate}/heartRate` | Resting HR, all-day HR | Cardiovascular health |
| **Steps** | `/wellness/daily/{summaryDate}/steps` | Daily step count | Activity level |

### Additional Metrics (Nice to Have)

- Pulse Ox (blood oxygen)
- Respiration rate
- Body composition (weight, BMI, body fat %)
- Menstrual cycle (if tracked)

---

## Authentication (OAuth 2.0 with PKCE)

### Step 1: Register App with Garmin

1. Go to [Garmin Developer Portal](https://developer.garmin.com/)
2. Create account
3. Register new application
4. Select "Health API" program
5. Get **Consumer Key** and **Consumer Secret**

### Step 2: OAuth Flow

```
1. User clicks "Connect Garmin"
   ↓
2. Redirect to Garmin authorization URL
   https://connectapi.garmin.com/oauth-service/oauth/authorize
   ?oauth_consumer_key={key}
   &oauth_callback={your_callback_url}
   ↓
3. User authorizes app
   ↓
4. Garmin redirects to callback URL with oauth_token
   ↓
5. Exchange token for access token
   POST https://connectapi.garmin.com/oauth-service/oauth/access_token
   ↓
6. Store access token and refresh token
   (Access token expires after 3 months)
```

### Step 3: Token Storage

Store in database (encrypted):
```json
{
  "user_id": "uuid",
  "access_token": "encrypted_token",
  "refresh_token": "encrypted_refresh",
  "token_expires_at": "2026-04-15T00:00:00Z"
}
```

---

## Data Fetching Strategies

### Strategy 1: Poll Architecture (Recommended for MVP)

**How it works**:
1. Cron job runs every hour
2. Call "Ping" endpoint to check for updates
3. If updates available, pull data for relevant summaries
4. Store in events table

**Implementation**:
```go
// Pseudo-code
func SyncGarminData(ctx context.Context, userID string) error {
    // 1. Get stored OAuth token
    token, err := db.GetGarminToken(ctx, userID)
    if err != nil {
        return err
    }

    // 2. Check for updates (Ping)
    updates, err := garminClient.CheckForUpdates(ctx, token, userID)
    if err != nil {
        return err
    }

    // 3. Pull each summary type that has updates
    for _, update := range updates {
        switch update.SummaryType {
        case "sleep":
            sleepData, _ := garminClient.GetSleep(ctx, token, update.Date)
            db.SaveEvent(ctx, "garmin_sleep", sleepData)
        case "activity":
            activityData, _ := garminClient.GetActivities(ctx, token, update.Date)
            db.SaveEvent(ctx, "garmin_activity", activityData)
        // ... other types
        }
    }

    return nil
}
```

**Cron Schedule**:
```yaml
# Run every hour
schedule: "0 * * * *"
```

**Pros**:
- Simple to implement
- Predictable (hourly updates)
- No webhook infrastructure needed
- Good enough for daily insights

**Cons**:
- Not real-time (up to 1 hour delay)
- Some unnecessary API calls (when no updates)

### Strategy 2: Webhook Architecture (Production)

**How it works**:
1. Register webhook URL with Garmin
2. Garmin sends POST request when data available
3. Your service receives notification
4. Pull data immediately
5. Store in events table

**Webhook Endpoint**:
```go
// POST /webhook/garmin
func HandleGarminWebhook(w http.ResponseWriter, r *http.Request) {
    // 1. Validate request signature
    if !validateGarminSignature(r) {
        http.Error(w, "Invalid signature", 401)
        return
    }

    // 2. Parse notification
    var notification GarminNotification
    json.NewDecoder(r.Body).Decode(&notification)

    // 3. Queue data fetch job (async)
    jobs.Enqueue("fetch_garmin_data", notification.UserID, notification.SummaryType)

    // 4. Respond immediately (Garmin expects 200 within seconds)
    w.WriteHeader(http.StatusOK)
}
```

**Pros**:
- Real-time (data available within minutes)
- Most efficient (only fetch when available)
- Recommended by Garmin

**Cons**:
- Requires public webhook endpoint
- Need to handle retries and failures
- More complex infrastructure

---

## Recommended Approach for This Project

### Phase 1 (MVP): Poll Architecture - Hourly

**Rationale**:
- Simpler to implement
- No need for public webhook yet (local development)
- Hourly is sufficient for daily insights and modeling
- Most Garmin users sync 1-3x per day anyway

**Implementation Plan**:
1. Create `GarminClient` in `backend/internal/garmin/`
2. Implement OAuth flow
3. Add cron job to `ingestion-service`
4. Poll every hour, pull available data
5. Store in `events` table

### Phase 2 (Production): Webhook Architecture

**Transition when**:
- Deploying to AWS
- Need real-time updates
- Ready to apply for production API keys

---

## Rate Limits & Best Practices

### Evaluation Keys (Development)
- Rate-limited (exact limits not public)
- Use ONLY for development/testing
- Do NOT use in production

### Production Keys
- Not rate-limited
- Requires Partner Verification
- Apply after MVP is working

### Best Practices
1. **Cache responses** - Don't re-fetch same data
2. **Respect 7-day retention** - Garmin only keeps data for 7 days
3. **Handle errors gracefully** - Token expiration, API downtime
4. **Log all API calls** - For debugging and quota monitoring
5. **Use Ping endpoint** - Don't blindly fetch every hour
6. **Batch requests** - Fetch multiple days if user hasn't synced in a while

---

## Data Mapping to Our Events Table

### Sleep Data
```json
{
  "time": "2026-01-12T08:00:00Z",
  "user_id": "uuid",
  "event_type": "garmin_sleep",
  "source": "garmin",
  "data": {
    "duration_minutes": 435,
    "deep_sleep_minutes": 92,
    "light_sleep_minutes": 254,
    "rem_sleep_minutes": 89,
    "awake_minutes": 12,
    "sleep_score": 82,
    "hrv_avg": 67.5
  },
  "confidence": 0.95
}
```

### Activity Data
```json
{
  "time": "2026-01-12T10:00:00Z",
  "user_id": "uuid",
  "event_type": "garmin_activity",
  "source": "garmin",
  "data": {
    "activity_type": "strength_training",
    "duration_minutes": 45,
    "calories": 285,
    "avg_hr": 132,
    "max_hr": 168,
    "intensity": "moderate"
  },
  "confidence": 0.95
}
```

### HRV Data
```json
{
  "time": "2026-01-12T08:00:00Z",
  "user_id": "uuid",
  "event_type": "garmin_hrv",
  "source": "garmin",
  "data": {
    "hrv_avg": 67.5,
    "hrv_min": 52.0,
    "hrv_max": 83.2,
    "measurement_time": "overnight"
  },
  "confidence": 0.95
}
```

---

## Implementation Checklist

### Week 1: Setup
- [ ] Register app on Garmin Developer Portal
- [ ] Get Consumer Key and Consumer Secret
- [ ] Store keys in `.env` file
- [ ] Create `GarminClient` struct in `backend/internal/garmin/`

### Week 2: OAuth Flow
- [ ] Implement OAuth authorization redirect
- [ ] Create callback handler
- [ ] Exchange code for access token
- [ ] Store tokens in database (encrypted)
- [ ] Test with your personal Garmin account

### Week 3: Data Fetching
- [ ] Implement Ping endpoint call
- [ ] Implement Sleep data fetch
- [ ] Implement Activity data fetch
- [ ] Implement HRV data fetch (if available on your device)
- [ ] Map responses to Event models

### Week 4: Ingestion Service
- [ ] Add cron scheduler to ingestion-service
- [ ] Run hourly sync job
- [ ] Handle errors (token expired, API down)
- [ ] Log all sync attempts
- [ ] Verify data in database

### Week 5: Testing
- [ ] Sync historical data (last 7 days)
- [ ] Verify data quality
- [ ] Check for missing data gaps
- [ ] Test token refresh flow

---

## Common Issues & Solutions

### Issue: Token Expired
**Solution**: Implement automatic token refresh
```go
if err == ErrTokenExpired {
    newToken, err := garminClient.RefreshToken(ctx, oldToken)
    db.UpdateGarminToken(ctx, userID, newToken)
    // Retry request
}
```

### Issue: Device Not Synced
**Solution**: Log warning, try again next hour
```go
if err == ErrNoDataAvailable {
    log.Warn("User hasn't synced device", "user_id", userID)
    return nil // Don't error, just skip
}
```

### Issue: Rate Limit Hit
**Solution**: Exponential backoff
```go
if err == ErrRateLimit {
    backoff := time.Minute * time.Duration(2^retryCount)
    time.Sleep(backoff)
    // Retry
}
```

---

## API Endpoints Reference

Base URL: `https://apis.garmin.com/wellness-api/rest`

### Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/oauth/authorize` | GET | Start OAuth flow |
| `/oauth/access_token` | POST | Get access token |
| `/user/id` | GET | Get user ID |
| `/ping` | GET | Check for data updates |
| `/daily/{summaryDate}` | GET | Get daily summary |
| `/dailies` | GET | Get multiple daily summaries |
| `/sleep/{summaryDate}` | GET | Get sleep data |
| `/activities/{summaryDate}` | GET | Get activities |
| `/hrv` | GET | Get HRV data |

---

## Security Considerations

1. **Never log tokens** - Use `[REDACTED]` in logs
2. **Encrypt tokens at rest** - Use encryption in database
3. **Use HTTPS** - All API calls must be HTTPS
4. **Validate webhooks** - Check signature if using webhook architecture
5. **Rotate secrets** - Regenerate consumer secret periodically

---

## Testing Strategy

### Manual Testing
1. Connect your personal Garmin account
2. Trigger manual sync
3. Check database for new events
4. Verify data accuracy against Garmin Connect app

### Automated Testing
```go
func TestGarminClient_GetSleep(t *testing.T) {
    client := NewGarminClient(testConfig)

    sleep, err := client.GetSleep(ctx, testToken, "2026-01-12")

    assert.NoError(t, err)
    assert.NotNil(t, sleep)
    assert.Greater(t, sleep.DurationMinutes, 0)
}
```

---

## Resources

- [Garmin Health API Documentation](https://developer.garmin.com/gc-developer-program/health-api/)
- [Garmin Developer Portal](https://developerportal.garmin.com/)
- [OAuth 2.0 PKCE Specification](https://developerportal.garmin.com/sites/default/files/OAuth2PKCE_1.pdf)
- [python-garminconnect (Reference Implementation)](https://github.com/cyberjunky/python-garminconnect)

---

## Next Steps

Once this is working:
1. Add retry logic for failed syncs
2. Implement data quality checks
3. Add Grafana dashboard for sync monitoring
4. Consider adding other data sources (Apple Health, Oura, etc.)

---

**Last Updated**: January 2026
**Status**: Design Complete, Ready to Implement
