# Feature Plan: Daily Check-in & Health Correlation

## 🎯 MVP Feature: Daily Check-in with Garmin Correlation

**Vision:** Help users understand how their daily habits affect how they feel by combining subjective check-ins with objective Garmin data.

## 📱 Technology Stack

### Frontend
- **Flutter** (mobile app - iOS & Android)
- **State Management:** Riverpod
- **HTTP Client:** Dio
- **Local Storage:** Hive/SharedPreferences
- **Charts:** fl_chart
- **Notifications:** flutter_local_notifications

### Backend
- **Go API** (existing services)
- **PostgreSQL + TimescaleDB** (existing database)
- **REST API** with JSON

## 🚀 Phase 1: MVP (2-3 Weeks)

### Week 1: Backend Foundation

**Day 1-2: Check-in API Endpoints**
- `POST /api/v1/checkin` - Submit daily check-in
- `GET /api/v1/checkin/latest` - Get today's check-in
- `GET /api/v1/checkin/history?days=30` - Get check-in history

**Day 3-4: Dashboard API Endpoints**
- `GET /api/v1/dashboard/today` - Today's summary (Garmin + check-in)
- `GET /api/v1/dashboard/week` - 7-day trends
- `GET /api/v1/health/summary?start=X&end=Y` - Date range data

**Day 5: Simple Correlation Logic**
- `GET /api/v1/insights/correlations?metric=sleep` - Basic correlations
- Calculate averages grouped by conditions
- Example: "When sleep > 7hrs, avg energy = 8.2"

**Database Schema:**
```sql
-- Already exists in events table!
-- Just need to use event_type = 'subjective_feeling'

-- Example:
{
  "time": "2026-01-29T08:00:00Z",
  "user_id": "uuid",
  "event_type": "subjective_feeling",
  "source": "manual",
  "data": {
    "energy": 8,
    "mood": 7,
    "focus": 9,
    "physical": 7,
    "notes": "Felt great after morning run"
  }
}
```

### Week 2: Flutter App Foundation

**Day 1-2: Project Setup**
- Create Flutter project structure
- Set up Riverpod for state management
- Configure API client (Dio)
- Set up routing (go_router)
- Design system (theme, colors, typography)

**Day 3-4: Authentication & Onboarding**
- Login screen
- Simple JWT authentication
- Onboarding flow (explain the concept)
- Request notification permissions

**Day 5: Core Features - Check-in Screen**
- Morning check-in form (4 sliders: energy, mood, focus, physical)
- Optional notes field
- Submit and store locally + sync to backend
- Celebration animation on submit

### Week 3: Dashboard & Insights

**Day 1-2: Today's Dashboard**
- Show today's check-in
- Display last night's Garmin data:
  - Sleep duration & score
  - HRV average
  - Activity summary
  - Stress level
- Beautiful card-based UI

**Day 3-4: 7-Day Trends**
- Line charts for energy/mood over 7 days
- Bar chart for sleep duration
- HRV trend
- Identify best/worst days

**Day 5: Simple Insights**
- Show correlation insights:
  - "Your energy is 15% higher when you sleep 7+ hours"
  - "You're most focused after active days"
  - "Your mood improves with lower stress"
- Progressive insights (unlock as data grows)

## 📐 App Structure

```
lib/
├── main.dart
├── app.dart
├── core/
│   ├── config/
│   │   ├── app_config.dart
│   │   └── theme.dart
│   ├── network/
│   │   ├── api_client.dart
│   │   ├── api_endpoints.dart
│   │   └── interceptors.dart
│   └── utils/
│       ├── date_utils.dart
│       └── validators.dart
├── features/
│   ├── auth/
│   │   ├── data/
│   │   │   ├── auth_repository.dart
│   │   │   └── auth_api.dart
│   │   ├── domain/
│   │   │   └── user.dart
│   │   ├── presentation/
│   │   │   ├── login_screen.dart
│   │   │   └── widgets/
│   │   └── providers/
│   │       └── auth_provider.dart
│   ├── checkin/
│   │   ├── data/
│   │   │   ├── checkin_repository.dart
│   │   │   └── checkin_api.dart
│   │   ├── domain/
│   │   │   └── checkin_model.dart
│   │   ├── presentation/
│   │   │   ├── checkin_screen.dart
│   │   │   └── widgets/
│   │   │       ├── feeling_slider.dart
│   │   │       └── submit_button.dart
│   │   └── providers/
│   │       └── checkin_provider.dart
│   ├── dashboard/
│   │   ├── data/
│   │   │   ├── dashboard_repository.dart
│   │   │   └── health_api.dart
│   │   ├── domain/
│   │   │   ├── garmin_summary.dart
│   │   │   └── daily_summary.dart
│   │   ├── presentation/
│   │   │   ├── dashboard_screen.dart
│   │   │   ├── trends_screen.dart
│   │   │   └── widgets/
│   │   │       ├── metric_card.dart
│   │   │       ├── sleep_card.dart
│   │   │       └── trend_chart.dart
│   │   └── providers/
│   │       └── dashboard_provider.dart
│   └── insights/
│       ├── data/
│       │   └── insights_repository.dart
│       ├── domain/
│       │   └── correlation.dart
│       ├── presentation/
│       │   ├── insights_screen.dart
│       │   └── widgets/
│       │       └── insight_card.dart
│       └── providers/
│           └── insights_provider.dart
└── shared/
    ├── widgets/
    │   ├── loading_indicator.dart
    │   ├── error_view.dart
    │   └── app_button.dart
    └── models/
        └── api_response.dart
```

## 🎨 UI/UX Design

### Color Palette
```dart
// Based on health/wellness theme
primary: Color(0xFF6C63FF),     // Vibrant purple
secondary: Color(0xFF4CAF50),   // Success green
background: Color(0xFFF5F7FA),  // Light gray
surface: Colors.white,
error: Color(0xFFE57373),       // Soft red
text: Color(0xFF2D3748),        // Dark gray

// Metric colors
sleep: Color(0xFF7C3AED),       // Purple
energy: Color(0xFFFBBF24),      // Yellow
mood: Color(0xFF3B82F6),        // Blue
focus: Color(0xFF10B981),       // Green
physical: Color(0xFFEF4444),    // Red
```

### Key Screens

#### 1. Home/Dashboard Screen
```
┌─────────────────────────────────┐
│  ☀️  Good morning, Satish!      │
│                                 │
│  ┌───────────────────────────┐ │
│  │  Today's Check-in          │ │
│  │  ─────────────────────     │ │
│  │  💪 Energy:        8/10    │ │
│  │  😊 Mood:          7/10    │ │
│  │  🎯 Focus:         9/10    │ │
│  │  🏃 Physical:      7/10    │ │
│  │                            │ │
│  │  Checked in at 8:30 AM    │ │
│  └───────────────────────────┘ │
│                                 │
│  Last Night                     │
│  ┌───────────────────────────┐ │
│  │  😴 Sleep                  │ │
│  │  7.2 hours · Score: 82     │ │
│  │  ▓▓▓░░░░░░░░ Deep 2.1h    │ │
│  │  ▓▓▓▓▓▓░░░░ Light 3.8h    │ │
│  │  ▓▓░░░░░░░░ REM 1.3h      │ │
│  └───────────────────────────┘ │
│                                 │
│  ┌──────────┐  ┌──────────┐   │
│  │ 💓 HRV    │  │ 😰 Stress │   │
│  │ 67 ms    │  │ 32 (low) │   │
│  └──────────┘  └──────────┘   │
│                                 │
│  🏃 Activity: 45 min active     │
│                                 │
│  [View Trends →]                │
└─────────────────────────────────┘
```

#### 2. Morning Check-in Screen
```
┌─────────────────────────────────┐
│  ← How are you feeling today?  │
│                                 │
│  Rate your current state:       │
│                                 │
│  💪 Energy                       │
│  ●━━━━━━━━○─────── [8]         │
│  Low              High          │
│                                 │
│  😊 Mood                         │
│  ●━━━━━━○──────── [7]          │
│  Low              High          │
│                                 │
│  🎯 Focus                        │
│  ●━━━━━━━━━○───── [9]          │
│  Low              High          │
│                                 │
│  🏃 Physical                     │
│  ●━━━━━━○──────── [7]          │
│  Low              High          │
│                                 │
│  📝 Notes (optional)             │
│  ┌─────────────────────────┐   │
│  │ Felt great after...     │   │
│  └─────────────────────────┘   │
│                                 │
│  [ Submit Check-in ]            │
└─────────────────────────────────┘
```

#### 3. 7-Day Trends Screen
```
┌─────────────────────────────────┐
│  ← Your Week                    │
│                                 │
│  Energy Levels                  │
│  10┐                 ●          │
│   9│              ●  │          │
│   8│           ●  │  │          │
│   7│        ●  │  │  │          │
│   6│     ●  │  │  │  │  ●       │
│   5└─────────────────────       │
│     M  T  W  T  F  S  S         │
│                                 │
│  Sleep Duration                 │
│  █ 8.2h  █ 7.1h  █ 7.8h        │
│  M       T       W              │
│                                 │
│  💡 Insights                     │
│  • Best day: Wednesday (9/10)   │
│  • Sleep correlation: +15%      │
│  • Most consistent: Weekdays    │
│                                 │
│  [See More Insights →]          │
└─────────────────────────────────┘
```

#### 4. Insights Screen
```
┌─────────────────────────────────┐
│  ← Personalized Insights        │
│                                 │
│  Based on 30 days of data       │
│                                 │
│  ┌───────────────────────────┐ │
│  │  😴 Sleep & Energy         │ │
│  │  ─────────────────────     │ │
│  │  Your energy is 15% higher │ │
│  │  when you sleep 7+ hours   │ │
│  │                            │ │
│  │  📊 [View Details]         │ │
│  └───────────────────────────┘ │
│                                 │
│  ┌───────────────────────────┐ │
│  │  🏃 Activity & Mood        │ │
│  │  ─────────────────────     │ │
│  │  Your mood improves by     │ │
│  │  12% on active days        │ │
│  │                            │ │
│  │  📊 [View Details]         │ │
│  └───────────────────────────┘ │
│                                 │
│  ┌───────────────────────────┐ │
│  │  💓 HRV & Recovery         │ │
│  │  ─────────────────────     │ │
│  │  Your HRV is highest on    │ │
│  │  low-stress days           │ │
│  │                            │ │
│  │  📊 [View Details]         │ │
│  └───────────────────────────┘ │
└─────────────────────────────────┘
```

## 🔔 Key Features

### 1. Daily Notifications
- Morning reminder: "Time for your daily check-in! ☀️"
- Smart timing: Learn user's wake time from Garmin data
- Configurable reminder time

### 2. Offline Support
- Save check-ins locally if offline
- Sync when connection restored
- Show cached data while loading

### 3. Simple Analytics
- Streak counter: "7 days in a row! 🔥"
- Progress badges
- Data completeness indicator

### 4. Privacy First
- Data stored locally when possible
- Clear data retention policy
- Export/delete data option

## 📊 Metrics to Track (Analytics)

User Engagement:
- Daily Active Users (DAU)
- Check-in completion rate
- Time to complete check-in
- Return rate (day 7, day 30)

Feature Usage:
- Dashboard views
- Trends views
- Insights views
- Notification engagement

Data Quality:
- Check-ins per user
- Garmin sync success rate
- Data gaps

## 🚧 Technical Implementation Details

### Backend APIs (Go)

**1. Check-in Submission**
```go
POST /api/v1/checkin
Authorization: Bearer <token>

Request:
{
  "energy": 8,
  "mood": 7,
  "focus": 9,
  "physical": 7,
  "notes": "Felt great after morning run"
}

Response:
{
  "id": "uuid",
  "timestamp": "2026-01-29T08:30:00Z",
  "data": {...}
}
```

**2. Dashboard Data**
```go
GET /api/v1/dashboard/today
Authorization: Bearer <token>

Response:
{
  "checkin": {
    "energy": 8,
    "mood": 7,
    "focus": 9,
    "physical": 7,
    "timestamp": "2026-01-29T08:30:00Z"
  },
  "garmin": {
    "sleep": {
      "duration_hours": 7.2,
      "score": 82,
      "deep_minutes": 126,
      "light_minutes": 228,
      "rem_minutes": 78,
      "awake_minutes": 0
    },
    "hrv": {
      "average": 67.5
    },
    "activity": {
      "active_minutes": 45,
      "steps": 8234,
      "calories": 2145
    },
    "stress": {
      "average": 32,
      "level": "low"
    }
  }
}
```

**3. Correlations**
```go
GET /api/v1/insights/correlations?days=30
Authorization: Bearer <token>

Response:
{
  "correlations": [
    {
      "type": "sleep_energy",
      "description": "Your energy is 15% higher when you sleep 7+ hours",
      "confidence": 0.85,
      "sample_size": 25,
      "details": {
        "condition": "sleep >= 7",
        "avg_energy_with": 8.2,
        "avg_energy_without": 7.1
      }
    },
    {
      "type": "activity_mood",
      "description": "Your mood improves by 12% on active days",
      "confidence": 0.78,
      "sample_size": 22,
      "details": {
        "condition": "active_minutes >= 30",
        "avg_mood_with": 7.8,
        "avg_mood_without": 6.9
      }
    }
  ]
}
```

### Flutter State Management (Riverpod)

```dart
// Providers
final checkinProvider = StateNotifierProvider<CheckinNotifier, CheckinState>(...);
final dashboardProvider = FutureProvider<DashboardData>(...);
final trendsProvider = FutureProvider<TrendsData>(...);
final insightsProvider = FutureProvider<List<Correlation>>(...);

// Usage in widget
final dashboard = ref.watch(dashboardProvider);

dashboard.when(
  data: (data) => DashboardView(data: data),
  loading: () => LoadingIndicator(),
  error: (error, stack) => ErrorView(error: error),
);
```

## 🎯 Success Criteria

**Week 3 Goals:**
- ✅ App deployed to TestFlight/Play Store (internal testing)
- ✅ Backend deployed to AWS
- ✅ Can submit daily check-in
- ✅ Can view today's dashboard
- ✅ Can view 7-day trends
- ✅ At least 1 simple correlation showing

**User Experience:**
- Check-in takes < 30 seconds
- App loads in < 2 seconds
- Insights appear after 7 days of data
- Works offline for check-ins

## 📅 Next Steps (After MVP)

### Phase 2: Experiment Tracking
- Create experiments: "Sleep 30min earlier for a week"
- Track compliance
- A/B test interventions
- Show before/after comparisons

### Phase 3: Advanced Analytics
- ML-based insights
- Personalized recommendations
- Predictive models: "You might feel low energy tomorrow"
- Multi-variate correlations

### Phase 4: Social & Gamification
- Share insights with friends
- Challenges and competitions
- Community insights
- Achievement system

## 🔧 Development Setup

1. **Backend:** Already exists, just need new endpoints
2. **Flutter:** New project, start from scratch
3. **Deployment:**
   - Backend: AWS ECS (already set up)
   - Flutter: TestFlight (iOS) + Google Play (Android)

## 📝 Documentation Needed

- [ ] API documentation (OpenAPI/Swagger)
- [ ] Flutter app README with setup instructions
- [ ] Design system documentation
- [ ] User guide / help screens in app

## 💡 Questions to Consider

1. **Authentication:** Use existing JWT from backend?
2. **Push notifications:** AWS SNS or Firebase Cloud Messaging?
3. **Analytics:** Firebase Analytics or custom solution?
4. **Crash reporting:** Sentry or Firebase Crashlytics?
5. **Backend changes:** New microservice or extend ingestion service?

---

**Ready to start building?** Let's begin with the backend API endpoints!
