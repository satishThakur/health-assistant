# Wearable Data Sources - Alternatives & Strategies

This document outlines realistic options for pulling wearable health data into your personal health assistant.

## The Problem

**Garmin Health API** is the ideal choice but has a critical limitation:
- ❌ Only approved for **commercial entities, healthcare orgs, research institutions**
- ❌ **Individual/personal developers are typically rejected**
- ❌ Not suitable for personal projects

We need practical alternatives.

---

## Option 1: Unofficial Garmin API (python-garminconnect)

### Overview
Use unofficial library that simulates Garmin Connect login and scrapes data.

### How It Works
```python
from garminconnect import Garmin

# Login with your Garmin credentials
client = Garmin("your_email@example.com", "your_password")
client.login()

# Fetch sleep data
sleep = client.get_sleep_data("2026-01-12")

# Fetch activities
activities = client.get_activities_by_date("2026-01-12", "2026-01-12")

# Fetch stats
stats = client.get_stats("2026-01-12")
```

### Available Data
- Sleep (duration, stages, quality)
- Activities (workouts with detailed metrics)
- Steps, calories, distance
- Heart rate (resting, max, zones)
- Stress scores
- Body Battery
- HRV (if available on device)

### Pros
✅ Works immediately
✅ Full access to your personal data
✅ Free
✅ Actively maintained library
✅ You already own a Garmin watch
✅ Many personal projects use this successfully

### Cons
⚠️ **Against Garmin's Terms of Service**
⚠️ Could break if Garmin changes their website
⚠️ Requires storing actual Garmin credentials (security concern)
⚠️ Rate limiting possible (account could be flagged/banned)
⚠️ Not suitable if you plan to commercialize

### Implementation Strategy

**For Go Backend:**
1. Create Python microservice for Garmin scraping
2. Your Go ingestion-service calls Python service
3. Python service uses python-garminconnect
4. Returns normalized JSON to Go service

**Architecture:**
```
Ingestion Service (Go)
         ↓
   HTTP Request
         ↓
Garmin Scraper (Python)
         ↓
python-garminconnect library
         ↓
Garmin Connect (unofficial)
```

**Security:**
- Store Garmin credentials encrypted in database
- Use environment variables in dev
- Never log credentials
- Consider using session tokens (library supports this)

### Risk Assessment

**Risk Level**: Medium

**Likelihood of Issues**:
- Account ban: Low (for personal, low-volume use)
- API breaking: Medium (happens occasionally, library updates quickly)
- Legal issues: Very Low (personal use, no commercial intent)

**Mitigation**:
- Design system to be wearable-agnostic (easy to switch)
- Keep credentials secure
- Limit request frequency (hourly polling is fine)
- Have backup plan (manual export, alternative wearable)

### Resources
- [python-garminconnect GitHub](https://github.com/cyberjunky/python-garminconnect)
- [garminconnect PyPI](https://pypi.org/project/garminconnect/)
- [garth (alternative)](https://github.com/matin/garth)

---

## Option 2: Oura Ring (Official Personal API) ⭐ RECOMMENDED

### Overview
Oura Ring provides **official API access for personal developers** - no business entity required!

### How It Works
```bash
# Get Personal Access Token from dashboard
curl https://api.ouraring.com/v2/usercollection/sleep \
  -H "Authorization: Bearer YOUR_PERSONAL_TOKEN"
```

### Available Data
- **Sleep**: Duration, stages (deep, light, REM), efficiency, latency
- **Readiness**: Daily readiness score (similar to Body Battery)
- **Activity**: Steps, calories, METs, activity levels
- **Heart Rate**: Resting HR, HRV (overnight and daytime)
- **Respiration**: Breathing rate
- **Temperature**: Skin temperature trends (early illness detection)
- **Workout**: Exercise heart rate data

### Pros
✅ **Official personal developer access**
✅ Excellent for sleep/recovery tracking (arguably better than Garmin)
✅ Clean, well-documented REST API
✅ Perfect for your use case (sleep quality, HRV, recovery focus)
✅ Great for n=1 experimentation
✅ Active community and support

### Cons
❌ Need to purchase Oura Ring (~$300-500)
❌ Requires active membership ($6/month)
❌ Less activity tracking than Garmin (no GPS, limited workout types)
❌ Personal Access Tokens being deprecated (end of 2025, moving to OAuth2)
❌ Gen3/Ring 4 users need active membership for API access

### Why Oura is Great for This Project
1. **Sleep focus**: Your primary outcome variable
2. **HRV quality**: Oura's HRV measurements are research-grade
3. **Recovery metrics**: Readiness score is excellent for tracking interventions
4. **Official API**: No TOS concerns, stable, supported
5. **N=1 research**: Designed with this use case in mind

### Implementation Strategy

**Phase 1: Personal Access Token**
```go
// backend/internal/oura/client.go
type OuraClient struct {
    token  string
    baseURL string
}

func (c *OuraClient) GetSleep(date string) (*SleepData, error) {
    url := fmt.Sprintf("%s/v2/usercollection/sleep?start_date=%s", c.baseURL, date)
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+c.token)
    // ... fetch and parse
}
```

**Phase 2: OAuth2 (post-2025)**
- Implement OAuth flow
- Store tokens in database
- Refresh tokens as needed

### Data Quality Comparison: Oura vs Garmin

| Metric | Oura | Garmin |
|--------|------|--------|
| Sleep Stages | ⭐⭐⭐⭐⭐ (Excellent) | ⭐⭐⭐⭐ (Good) |
| HRV | ⭐⭐⭐⭐⭐ (Research-grade) | ⭐⭐⭐⭐ (Good) |
| Recovery | ⭐⭐⭐⭐⭐ (Readiness Score) | ⭐⭐⭐⭐ (Body Battery) |
| Activity Tracking | ⭐⭐⭐ (Basic) | ⭐⭐⭐⭐⭐ (GPS, detailed) |
| Workout Tracking | ⭐⭐ (Manual) | ⭐⭐⭐⭐⭐ (Auto-detect, GPS) |
| API Access | ⭐⭐⭐⭐⭐ (Personal OK) | ❌ (Commercial only) |

**Verdict**: For sleep quality and recovery modeling, Oura is arguably **better** than Garmin + has official API.

### Resources
- [Oura API Documentation](https://cloud.ouraring.com/docs/)
- [Personal Use Guide](https://developer.ouraring.com/docs/personal-use)
- [Oura Developer Dashboard](https://cloud.ouraring.com/)

---

## Option 3: Terra API (Third-Party Aggregator)

### Overview
[Terra](https://tryterra.co/) is a unified API for 100+ wearables including Garmin, Oura, Fitbit, Apple, Whoop, etc.

### How It Works
```javascript
// Terra handles the integration
const data = await terra.getSleep(userId, startDate, endDate);
// Normalized data from ANY connected wearable
```

### Pros
✅ Terra is approved by Garmin (they handle the commercial relationship)
✅ One API for multiple wearables
✅ Future-proof (can switch wearables without code changes)
✅ Good if you plan to support multiple devices
✅ Good if you plan to commercialize later

### Cons
❌ **Costs money** (pricing requires contact, likely $50-200/month)
❌ Adds external dependency
❌ Overkill for single personal user
❌ Still need users to connect their wearables

### When to Use
- You want to support multiple wearables from day 1
- You plan to eventually support other users
- You have budget for this
- You want commercial-grade reliability

### Resources
- [Terra Website](https://tryterra.co/)
- [Terra Documentation](https://docs.tryterra.co/)
- [Terra Garmin Integration](https://tryterra.co/integrations/garmin)

---

## Option 4: Manual Export

### Overview
Manually export data from Garmin Connect and upload to your system.

### How It Works
1. Garmin Connect → Account → Export Data
2. Download CSV files (daily, weekly, or on-demand)
3. Upload to your system via UI or script
4. Parse and store in events table

### Pros
✅ Completely legitimate
✅ No API concerns
✅ Free

### Cons
❌ Manual effort (defeats automation)
❌ Not real-time
❌ Tedious
❌ Easy to forget/skip

### When to Use
- Fallback if APIs fail
- Historical data import (one-time)
- Testing your data pipeline

---

## Recommended Strategy

### For Your Personal Health Assistant

#### Phase 1: MVP (Months 1-3)
**Use: python-garminconnect (unofficial)**

**Why:**
- You already own Garmin watch
- Works immediately
- Free
- Good enough for MVP validation

**Risk Mitigation:**
- Design system to be wearable-agnostic (interface pattern)
- Low request volume (hourly polling)
- Encrypt credentials
- Have backup plan

**Implementation:**
```go
// Design for flexibility
type WearableClient interface {
    GetSleep(ctx context.Context, date time.Time) (*SleepData, error)
    GetActivity(ctx context.Context, date time.Time) (*ActivityData, error)
    GetHRV(ctx context.Context, date time.Time) (*HRVData, error)
}

// Multiple implementations
type GarminUnofficialClient struct { /* python-garminconnect */ }
type OuraOfficialClient struct { /* Oura API */ }
type TerraClient struct { /* Terra aggregator */ }
```

#### Phase 2: Validation (Months 3-6)
**Continue with Garmin OR switch to Oura**

**Decision Criteria:**
- Is the project proving valuable?
- Are you committed to long-term?
- Do you trust the Garmin unofficial approach?

**If switching to Oura:**
- Purchase Oura Ring (~$350)
- Subscribe to membership ($6/month)
- Implement Oura client
- Collect parallel data for 2-4 weeks
- Validate model accuracy
- Full switch

#### Phase 3: Production (Months 6+)
**Official Oura API OR Terra**

**If staying personal:**
- Oura Ring with official API
- Reliable, supported, perfect for use case

**If going commercial:**
- Terra API
- Support multiple wearables
- Scalable for users

---

## Implementation: Wearable-Agnostic Design

### Interface Pattern

```go
// backend/internal/wearables/client.go
package wearables

import (
    "context"
    "time"
)

// Client interface for all wearable data sources
type Client interface {
    GetSleep(ctx context.Context, date time.Time) (*SleepData, error)
    GetActivity(ctx context.Context, date time.Time) ([]ActivityData, error)
    GetHRV(ctx context.Context, date time.Time) (*HRVData, error)
    GetStress(ctx context.Context, date time.Time) (*StressData, error)
    GetDailySummary(ctx context.Context, date time.Time) (*DailySummary, error)
}

// Normalized data structures
type SleepData struct {
    Date              time.Time
    DurationMinutes   int
    DeepMinutes       int
    LightMinutes      int
    REMMinutes        int
    AwakeMinutes      int
    SleepScore        int     // 0-100
    AverageHRV        float64
    LowestHeartRate   int
}

type ActivityData struct {
    Date            time.Time
    ActivityType    string
    DurationMinutes int
    Calories        int
    AverageHR       int
    MaxHR           int
    Distance        float64 // meters
}

type HRVData struct {
    Date           time.Time
    AverageHRV     float64
    MinHRV         float64
    MaxHRV         float64
    MeasurementWindow string // "overnight", "all_day", "workout"
}
```

### Factory Pattern

```go
// backend/internal/wearables/factory.go
func NewClient(source string, config Config) (Client, error) {
    switch source {
    case "garmin_unofficial":
        return NewGarminUnofficialClient(config)
    case "oura":
        return NewOuraClient(config)
    case "terra":
        return NewTerraClient(config)
    default:
        return nil, fmt.Errorf("unknown wearable source: %s", source)
    }
}
```

### Configuration

```go
// backend/internal/config/config.go
type WearableConfig struct {
    Source string // "garmin_unofficial", "oura", "terra"

    // Garmin (unofficial)
    GarminEmail    string
    GarminPassword string
    GarminPythonServiceURL string

    // Oura
    OuraPersonalToken string

    // Terra
    TerraAPIKey string
    TerraDevID  string
}
```

---

## Security Considerations

### For Unofficial APIs
- **Encrypt credentials** at rest (database encryption)
- **Never log credentials**
- Use secure password hashing (bcrypt)
- Rotate credentials periodically
- Monitor for unusual activity

### For Official APIs
- **Store tokens encrypted**
- Implement token refresh logic
- Use HTTPS for all requests
- Set appropriate token scopes
- Monitor API usage

---

## Decision Matrix

| Factor | Garmin (Unofficial) | Oura (Official) | Terra |
|--------|---------------------|-----------------|-------|
| **Cost** | Free | $350 + $6/mo | $50-200/mo |
| **Official Support** | ❌ | ✅ | ✅ |
| **Sleep Quality** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Depends |
| **Activity Tracking** | ⭐⭐⭐⭐⭐ | ⭐⭐ | Depends |
| **API Stability** | ⚠️ Medium | ✅ High | ✅ High |
| **Setup Time** | Days | Days | Hours |
| **Legal Risk** | ⚠️ Medium | ✅ None | ✅ None |
| **Best For** | MVP Testing | Long-term Personal | Multi-user/Commercial |

---

## My Recommendation

**Start:** python-garminconnect (unofficial)
**Evolve:** Oura Ring (official) if project proves valuable
**Scale:** Terra API if going commercial

**Why:**
1. Validate the concept first (Garmin, free)
2. Switch to official when committed (Oura)
3. Oura is arguably better for your use case anyway (sleep/recovery focus)
4. Wearable-agnostic design makes switching painless

---

## Next Steps

1. **Week 1**: Implement wearable interface pattern
2. **Week 2**: Build Garmin unofficial client (Python service)
3. **Week 3**: Test with your personal data
4. **Month 3**: Evaluate switching to Oura
5. **Month 6**: Make long-term decision

---

## Resources

### Unofficial Garmin
- [python-garminconnect](https://github.com/cyberjunky/python-garminconnect)
- [Unofficial Garmin API Guide](https://wiki.brianturchyn.net/programming/apis/garmin/)

### Oura Official
- [Oura API Docs](https://cloud.ouraring.com/docs/)
- [Personal Use Guide](https://developer.ouraring.com/docs/personal-use)

### Terra API
- [Terra Website](https://tryterra.co/)
- [Terra Documentation](https://docs.tryterra.co/)

### Alternative Wearables
- **Apple Health**: HealthKit API (iOS only)
- **Whoop**: Official API (requires Whoop strap)
- **Fitbit**: API available but Google ecosystem

---

**Last Updated**: January 2026
**Status**: Strategy Defined, Ready to Implement
