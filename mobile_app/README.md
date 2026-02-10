# Health Assistant - Flutter Mobile App

A Flutter mobile app for daily health check-ins with personalized insights based on Garmin data.

## Features

- **Daily Check-in**: Submit subjective feelings (energy, mood, focus, physical) on a 1-10 scale
- **Dashboard**: View today's check-in alongside Garmin health data (sleep, HRV, stress, activity)
- **7-Day Trends**: Interactive charts showing patterns over the week
- **Personalized Insights**: Automatic correlations between your feelings and health metrics

## Architecture

### Clean Architecture + Riverpod

```
lib/
├── main.dart                 # App entry point
├── app.dart                  # Main app widget
├── core/                     # Core functionality
│   ├── config/              # App configuration & theme
│   ├── network/             # API client & networking
│   └── routing/             # Navigation setup
├── features/                 # Feature modules
│   ├── checkin/            # Daily check-in feature
│   │   ├── data/           # API & repository
│   │   ├── domain/         # Models
│   │   ├── presentation/   # UI screens & widgets
│   │   └── providers/      # Riverpod state management
│   ├── dashboard/          # Dashboard & trends
│   └── insights/           # Correlation insights
└── shared/                  # Shared widgets
```

## Getting Started

### Prerequisites

- Flutter SDK 3.0+
- Dart 3.0+
- Backend API running (see main README)

### Installation

1. **Install dependencies:**
```bash
cd mobile_app
flutter pub get
```

2. **Generate code:**
```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

3. **Configure API endpoint:**

Edit `lib/core/config/app_config.dart`:
```dart
static const String baseUrl = 'http://YOUR_API_URL:8083';
```

Or use environment variable:
```bash
flutter run --dart-define=API_BASE_URL=http://localhost:8083
```

### Running the App

**iOS Simulator:**
```bash
flutter run -d ios
```

**Android Emulator:**
```bash
flutter run -d android
```

**Web (for testing):**
```bash
flutter run -d chrome
```

## Project Structure

### Core

- **Theme**: Material 3 design with custom health-related colors
- **API Client**: Dio-based HTTP client with interceptors
- **Routing**: go_router for declarative navigation
- **Config**: Centralized app configuration

### Features

#### Check-in Module
- Form with 4 sliders (energy, mood, focus, physical)
- Optional notes field
- Validation (1-10 scale, max 1000 chars)
- Upsert logic (one check-in per day)

#### Dashboard Module
- Today's check-in summary
- Last night's sleep data
- Current HRV and stress levels
- Quick navigation to trends and insights

#### Trends Module
- Interactive line charts (using fl_chart)
- Switch between metrics (energy, mood, focus, physical, sleep)
- 7-day historical view

#### Insights Module
- Automatic correlation detection
- Sleep vs Energy
- Activity vs Mood
- Sleep Quality vs Focus

## State Management

Using **Riverpod** for state management:

### Providers

```dart
// Check-in form state
final checkinFormProvider = StateNotifierProvider<CheckinFormNotifier, CheckinFormState>

// Today's dashboard data
final todayDashboardProvider = FutureProvider<DashboardData>

// 7-day trends
final weekTrendsProvider = FutureProvider<List<TrendData>>

// Correlations
final correlationsProvider = FutureProvider.family<List<CorrelationInsight>, int>
```

## API Integration

### Endpoints Used

```dart
POST   /api/v1/checkin              // Submit check-in
GET    /api/v1/checkin/latest       // Get today's check-in
GET    /api/v1/dashboard/today      // Dashboard data
GET    /api/v1/trends/week          // 7-day trends
GET    /api/v1/insights/correlations // Insights
```

### Models

All API models use `json_serializable` for automatic JSON parsing:

```dart
@JsonSerializable()
class CheckinModel {
  final int energy;
  final int mood;
  final int focus;
  final int physical;
  final String? notes;
}
```

## UI Components

### Custom Widgets

- **FeelingSlider**: Customizable slider for feelings (1-10)
- **CheckinCard**: Display today's check-in
- **SleepCard**: Visualize sleep stages
- **MetricCard**: Show individual metrics (HRV, stress)
- **TrendChart**: Interactive line chart
- **InsightCard**: Display correlation insights

### Theme

Using Material 3 with custom color palette:

```dart
Primary:   #6C63FF (Vibrant purple)
Secondary: #4CAF50 (Success green)

Metrics:
Sleep:    #7C3AED (Purple)
Energy:   #FBBF24 (Yellow)
Mood:     #3B82F6 (Blue)
Focus:    #10B981 (Green)
Physical: #EF4444 (Red)
```

## Dependencies

### Core
- `flutter_riverpod`: State management
- `dio`: HTTP client
- `go_router`: Navigation
- `hive`: Local storage

### UI
- `fl_chart`: Charts and graphs
- `intl`: Date formatting

### Code Generation
- `json_serializable`: JSON parsing
- `riverpod_generator`: Provider generation
- `build_runner`: Code generation runner

## Development

### Code Generation

When modifying models or providers:

```bash
# Watch for changes
flutter pub run build_runner watch

# One-time build
flutter pub run build_runner build --delete-conflicting-outputs
```

### Linting

```bash
flutter analyze
```

### Testing

```bash
# Run all tests
flutter test

# Run with coverage
flutter test --coverage
```

## Building

### Android

```bash
flutter build apk --release
# Output: build/app/outputs/flutter-apk/app-release.apk
```

### iOS

```bash
flutter build ios --release
# Then open Xcode to archive and upload
```

### Web

```bash
flutter build web --release
# Output: build/web/
```

## Configuration

### Environment Variables

```bash
# Development
flutter run --dart-define=API_BASE_URL=http://localhost:8083

# Production
flutter run --dart-define=API_BASE_URL=https://api.yourapp.com
```

### Build Variants

Create `lib/core/config/env.dart`:

```dart
enum Environment { dev, staging, prod }

const environment = Environment.dev; // Change per build
```

## TODO

- [ ] Add authentication (JWT tokens)
- [ ] Implement offline mode with local caching
- [ ] Add push notifications for daily reminders
- [ ] Export data functionality
- [ ] Dark mode toggle
- [ ] Onboarding screens
- [ ] Unit tests for providers
- [ ] Widget tests for screens
- [ ] Integration tests

## Troubleshooting

### API Connection Issues

1. Check backend is running: `curl http://localhost:8083/health`
2. Use correct IP (not `localhost`) for physical devices
3. Check network permissions in Android/iOS config

### Build Errors

1. Clean build: `flutter clean && flutter pub get`
2. Regenerate code: `flutter pub run build_runner build --delete-conflicting-outputs`
3. Update dependencies: `flutter pub upgrade`

### iOS Simulator Issues

1. Reset simulator: Device → Erase All Content and Settings
2. Check Xcode version: `xcode-select --version`

## Resources

- [Flutter Documentation](https://docs.flutter.dev/)
- [Riverpod Documentation](https://riverpod.dev/)
- [fl_chart Examples](https://github.com/imaNNeoFighT/fl_chart)
- [Backend API Docs](../CHECKIN_API_README.md)

## License

MIT License - See LICENSE file for details
