# Health Assistant Mobile App

Flutter application for personal health tracking and experiment management.

## Setup

### Prerequisites

- Flutter SDK (3.16.0 or later)
- Dart SDK (3.2.0 or later)
- Android Studio / Xcode (for mobile development)
- VS Code or Android Studio with Flutter plugins

### Create Flutter Project

The Flutter app needs to be initialized using `flutter create`:

```bash
cd app
flutter create health_assistant
cd health_assistant
```

### Recommended Project Structure

```
app/health_assistant/
├── lib/
│   ├── main.dart
│   ├── models/              # Data models
│   │   ├── event.dart
│   │   ├── experiment.dart
│   │   └── user.dart
│   │
│   ├── services/            # API clients and services
│   │   ├── api_client.dart
│   │   ├── auth_service.dart
│   │   └── event_service.dart
│   │
│   ├── screens/             # UI screens
│   │   ├── dashboard/
│   │   ├── daily_log/
│   │   ├── experiments/
│   │   ├── insights/
│   │   └── timeline/
│   │
│   ├── widgets/             # Reusable widgets
│   │   ├── metric_card.dart
│   │   ├── feeling_slider.dart
│   │   └── meal_photo.dart
│   │
│   └── providers/           # State management (Riverpod/Bloc)
│       ├── auth_provider.dart
│       ├── event_provider.dart
│       └── experiment_provider.dart
│
├── test/
├── android/
├── ios/
├── web/
└── pubspec.yaml
```

## Key Features

### Screens

1. **Dashboard**
   - Today's date and Garmin sync status
   - Key metrics (sleep, HRV, steps, body battery)
   - Active experiment status
   - Quick action buttons

2. **Daily Log**
   - Morning/Evening subjective feelings (energy, mood, focus, physical)
   - Slider inputs (1-10 scale)
   - Notes field
   - Timestamp tracking

3. **Meal Logging**
   - Camera integration for meal photos
   - Photo upload to backend
   - Display extracted macros from LLM
   - Manual override/editing
   - Meal type selection

4. **Supplement Tracking**
   - Checklist of daily supplements
   - Mark as taken with timestamp
   - Track compliance
   - Reminders

5. **Experiments**
   - List of proposed experiments (swipe to accept/reject)
   - Active experiment tracking
   - Compliance monitoring
   - Results visualization

6. **Insights**
   - Model-generated insights
   - Correlation visualizations
   - Time-series charts
   - Feature importance

7. **Timeline/Data Explorer**
   - Chronological view of all events
   - Filter by type and date range
   - Search functionality

## Dependencies

Key packages to add to `pubspec.yaml`:

```yaml
dependencies:
  flutter:
    sdk: flutter

  # State Management
  flutter_riverpod: ^2.4.0  # or bloc: ^8.1.0

  # Networking
  dio: ^5.4.0
  retrofit: ^4.0.0

  # Local Storage
  shared_preferences: ^2.2.0
  hive: ^2.2.3

  # UI Components
  flutter_svg: ^2.0.0
  cached_network_image: ^3.3.0

  # Camera & Images
  image_picker: ^1.0.0
  camera: ^0.10.0

  # Charts & Visualization
  fl_chart: ^0.65.0
  syncfusion_flutter_charts: ^24.0.0

  # Date & Time
  intl: ^0.18.0
  timezone: ^0.9.0

  # Authentication
  flutter_secure_storage: ^9.0.0

  # Utilities
  logger: ^2.0.0
  uuid: ^4.0.0
```

## API Configuration

Configure the backend API endpoint in `lib/config/api_config.dart`:

```dart
class ApiConfig {
  static const String baseUrl = 'http://localhost:8080';
  static const String modelServiceUrl = 'http://localhost:8084';
}
```

## Running the App

```bash
# Get dependencies
flutter pub get

# Run on mobile
flutter run

# Run on web
flutter run -d chrome

# Run on specific device
flutter devices
flutter run -d <device-id>
```

## Building

```bash
# Android APK
flutter build apk

# iOS
flutter build ios

# Web
flutter build web
```

## State Management

Recommend using **Riverpod** for state management:
- Simple, performant, compile-safe
- Good for this project's scale
- Easy to test

Alternative: **Bloc** (if you prefer more structure)

## Testing

```bash
# Run all tests
flutter test

# Run with coverage
flutter test --coverage
```

## Next Steps

- [ ] Initialize Flutter project with `flutter create`
- [ ] Set up folder structure as above
- [ ] Add dependencies to pubspec.yaml
- [ ] Create API client service
- [ ] Implement authentication flow
- [ ] Build dashboard screen
- [ ] Build daily log screen
- [ ] Integrate camera for meal photos
- [ ] Add state management
- [ ] Connect to backend APIs
