# Feature Plan: Daily Check-in & Health Correlation

## 🎯 MVP Feature: Daily Check-in with Garmin Correlation

**Vision:** Help users understand how their daily habits affect how they feel by combining subjective check-ins with objective Garmin data.

**Last Updated:** 2026-02-18
**Current Status:** Auth + platform config complete — offline support + notifications are next priorities

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

---

## 🚀 Phase 1: MVP (2-3 Weeks)

### Week 1: Backend Foundation ✅ COMPLETE

**Day 1-2: Check-in API Endpoints** ✅
- ✅ `POST /api/v1/checkin` - Submit daily check-in
- ✅ `GET /api/v1/checkin/latest` - Get today's check-in
- ✅ `GET /api/v1/checkin/history?days=30` - Get check-in history

**Day 3-4: Dashboard API Endpoints** ✅
- ✅ `GET /api/v1/dashboard/today` - Today's summary (Garmin + check-in)
- ✅ `GET /api/v1/trends/week` - 7-day trends
- ⬜ `GET /api/v1/health/summary?start=X&end=Y` - Date range data (not yet implemented)

**Day 5: Simple Correlation Logic** ✅
- ✅ `GET /api/v1/insights/correlations?days=30` - Basic correlations
- ✅ Calculate averages grouped by conditions
- ✅ Example: "When sleep > 7hrs, avg energy = 8.2"

**Garmin Ingestion Endpoints** ✅
- ✅ `POST /api/v1/garmin/ingest/sleep`
- ✅ `POST /api/v1/garmin/ingest/activity`
- ✅ `POST /api/v1/garmin/ingest/hrv`
- ✅ `POST /api/v1/garmin/ingest/stress`
- ✅ `POST /api/v1/garmin/ingest/daily-stats`
- ✅ `POST /api/v1/garmin/ingest/body-battery`

**Audit / Observability** ✅
- ✅ Sync audit endpoints (POST, GET recent, GET by type, GET stats)

**Authentication** ✅ COMPLETE
- ✅ Google Sign-In → JWT issued by backend (`/api/v1/auth/google`)
- ✅ JWT middleware wires real `user_id` into all handlers
- ✅ Garmin ingest routes protected by `X-Ingest-Secret` header
- ✅ Token stored in Keychain (iOS) / Keystore (Android) via flutter_secure_storage

**Database Schema:** ✅ Uses existing `events` table with `event_type = 'subjective_feeling'`

---

### Week 2: Flutter App Foundation ✅ COMPLETE

**Day 1-2: Project Setup** ✅
- ✅ Flutter project structure created (`mobile_app/`)
- ✅ Riverpod for state management
- ✅ Dio API client configured with interceptors
- ✅ go_router routing set up
- ✅ Design system (theme, colors, typography)

**Day 3-4: Authentication & Onboarding** ✅ COMPLETE
- ✅ Login screen (Google Sign-In)
- ✅ JWT authentication (backend + Flutter)
- ❌ Onboarding flow
- ❌ Notification permission request

**Day 5: Core Features - Check-in Screen** ✅
- ✅ Check-in form (4 sliders: energy, mood, focus, physical)
- ✅ Optional notes field
- ✅ Submit to backend
- ❌ Celebration animation on submit (not yet added)

---

### Week 3: Dashboard & Insights ✅ MOSTLY COMPLETE

**Day 1-2: Today's Dashboard** ✅
- ✅ Today's check-in card
- ✅ Last night's sleep data (duration, score, deep/light/REM)
- ✅ HRV average metric card
- ✅ Stress level metric card
- ✅ Daily activity stats card (steps, calories, active minutes)
- ✅ Body Battery card
- ✅ Navigation cards → Trends, Insights

**Day 3-4: 7-Day Trends** ✅
- ✅ Trends screen with TrendChart widget
- ✅ Quick insights summary (days tracked, consistency)
- ⬜ Line charts for energy/mood specifically (TrendChart in place, depth TBD)
- ⬜ Best/worst day identification (not yet implemented)

**Day 5: Simple Insights** ✅
- ✅ Insights screen showing correlation cards
- ✅ Empty state: "Not enough data yet" with guidance
- ✅ InsightCard widget
- ⬜ Progressive insight unlocking (basic structure in place, not fully wired)

---

## 📐 App Structure — Current State

```
mobile_app/lib/
├── main.dart                          ✅
├── app.dart                           ✅
├── core/
│   ├── config/
│   │   ├── app_config.dart            ✅
│   │   └── theme.dart                 ✅
│   ├── network/
│   │   ├── api_client.dart            ✅
│   │   ├── api_endpoints.dart         ✅
│   │   └── api_interceptor.dart       ✅
│   └── routing/
│       └── app_router.dart            ✅
├── features/
│   ├── auth/                          ✅ (domain/data/providers/presentation)
│   ├── checkin/
│   │   ├── data/                      ✅
│   │   ├── domain/                    ✅
│   │   ├── presentation/              ✅ (checkin_screen + feeling_slider)
│   │   └── providers/                 ✅
│   ├── dashboard/
│   │   ├── data/                      ✅
│   │   ├── domain/                    ✅
│   │   ├── presentation/              ✅ (dashboard, trends screens + all widgets)
│   │   └── providers/                 ✅
│   └── insights/
│       └── presentation/              ✅ (insights_screen + insight_card)
│       ❌ data/ and providers/ missing (wired through dashboard_provider)
└── shared/
    └── widgets/                       ✅ (loading_indicator, error_view)
```

---

## 🔔 Key Features Status

| Feature | Status | Notes |
|---|---|---|
| Daily Notifications | ❌ Not started | `flutter_local_notifications` planned |
| Offline Support | ❌ Not started | Save & sync check-ins locally |
| Streak Counter | ❌ Not started | "7 days in a row 🔥" |
| Progress Badges | ❌ Not started | |
| Celebration Animation | ❌ Not started | On check-in submit |
| JWT Auth (backend) | ✅ Done | Google Sign-In → JWT, middleware on all routes |
| JWT Auth (Flutter) | ✅ Done | Login screen, AuthProvider, secure token storage |

---

## 🎯 Success Criteria

| Criteria | Status |
|---|---|
| Can submit daily check-in | ✅ Done |
| Can view today's dashboard | ✅ Done |
| Can view 7-day trends | ✅ Done |
| At least 1 simple correlation showing | ✅ Done |
| App deployed to TestFlight/Play Store | ❌ Not done |
| Backend deployed to AWS | ❌ Not done |
| Check-in takes < 30 seconds | ✅ UI is simple enough |
| Works offline for check-ins | ❌ Not done |
| Insights appear after 7 days of data | ⬜ Logic exists, not fully tuned |

---

## 🛣️ Immediate Next Steps (Recommended Priority)

1. **Platform Config** ✅ COMPLETE
   - ✅ Android: Google Services plugin, `applicationId`, `minSdk = 21`
   - ✅ iOS: Bundle ID, `REVERSED_CLIENT_ID` URL scheme slot, Keychain entitlements, `Runner.entitlements`
   - ✅ `.env.example` and `google-services.json.example` / `GoogleService-Info.plist.example` committed
   - ⚠️ Manual remaining: download real `google-services.json` + `GoogleService-Info.plist` from Google Cloud Console, set backend env vars

2. **Offline Check-in Support**
   - Store check-ins locally with Hive if backend is unavailable
   - Sync on reconnection

3. **Daily Notifications**
   - Morning reminder using `flutter_local_notifications`
   - Configurable time

4. **Celebration Animation**
   - Lottie or confetti animation after submitting check-in

5. **Deployment**
   - Backend to AWS ECS (infra partially in place via docker-compose)
   - Flutter to TestFlight (iOS) + Google Play (Android)

---

## 📅 Future Phases (Unstarted)

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

---

## 🔧 Development Setup

1. **Backend:** `ingestion-service` on port `:8083` (all routes live here)
2. **Flutter:** `mobile_app/` — run with `flutter run`
3. **Infra:** `infra/docker-compose.yml` for local Postgres + TimescaleDB
4. **Deployment:**
   - Backend: AWS ECS (docker infrastructure in place)
   - Flutter: TestFlight (iOS) + Google Play (Android)

## 📝 Documentation Status

- [x] Garmin integration guide (`docs/garmin-integration-guide.md`)
- [x] TimescaleDB aggregation strategy (`docs/timescaledb-aggregation-strategy.md`)
- [x] High-level design (`docs/highleveldesign.md`)
- [x] Check-in API README (`CHECKIN_API_README.md`)
- [x] Flutter app summary (`FLUTTER_APP_SUMMARY.md`)
- [ ] OpenAPI/Swagger documentation
- [ ] Flutter app setup README
- [ ] Design system documentation
- [ ] User guide / help screens in app

## 💡 Open Questions

1. ~~**Authentication:** JWT from backend — middleware not yet wired~~ ✅ Done
2. **Push notifications:** AWS SNS or Firebase Cloud Messaging?
3. **Analytics:** Firebase Analytics or custom solution?
4. **Crash reporting:** Sentry or Firebase Crashlytics?
5. **`/api/v1/health/summary`:** Still needed or covered by dashboard/trends?
6. **Insights `data/` layer:** Should correlations have their own feature folder vs piggyback on dashboard_provider?

---

**Current Phase:** Phase 1 MVP — Core features built, auth + deployment + polish remaining.
