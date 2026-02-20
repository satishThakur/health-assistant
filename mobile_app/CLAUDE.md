# Mobile App — Flutter

## Entry Points
- `lib/main.dart` — app bootstrap, ProviderScope
- `lib/app.dart` — MaterialApp.router, watches authProvider for loading state
- `lib/core/routing/app_router.dart` — GoRouter, auth redirect guard

## Feature Structure
```
lib/features/
  auth/
    domain/auth_state.dart      # sealed class: AuthLoading | AuthUnauthenticated | AuthAuthenticated
    data/auth_service.dart      # GoogleSignIn → backend JWT, flutter_secure_storage
    providers/auth_provider.dart # AuthNotifier, authProvider, authTokenProvider
    presentation/login_screen.dart
  checkin/
    data/checkin_api.dart        # Dio calls to /api/v1/checkin
    data/checkin_repository.dart
    data/offline_queue_service.dart  # Hive Box<String> queue — enqueue/getPending/remove
    domain/checkin_model.dart
    providers/checkin_provider.dart  # CheckinFormNotifier, SubmitResult sealed class
    providers/sync_provider.dart     # offlineQueueServiceProvider, pendingCountProvider, SyncNotifier
    presentation/checkin_screen.dart
  health/                       # dashboard, trends, insights screens
lib/core/
  network/api_client.dart          # Dio instance
  network/api_interceptor.dart     # injects Bearer token, signs out on 401
  network/connectivity_service.dart # connectivity_plus wrapper — isOnline / onStatusChange
  routing/app_router.dart
  config/app_config.dart           # base URL, storage keys, pendingCheckinsBox
```

## State Management
- **Riverpod** (`flutter_riverpod`) — all providers in `features/*/providers/`
- `authProvider` = `StateNotifierProvider<AuthNotifier, AuthState>`
- `authTokenProvider` = derived `Provider<String?>` — returns token if authenticated

## Offline Check-in
- `OfflineQueueService` stores check-ins as JSON in a Hive `Box<String>` (key = ISO8601 timestamp)
- `CheckinFormNotifier.submitCheckin()` returns `SubmitResult` (sealed: `SubmitSuccess | SubmitSavedOffline | SubmitError`)
- On network error → queues locally, returns `SubmitSavedOffline`; orange snackbar shown in `checkin_screen.dart`
- `SyncNotifier` (bootstrapped in `app.dart`) listens to `ConnectivityService.onStatusChange`; auto-replays queue when online
- `pendingCountProvider` drives the orange banner above the submit button
- Hive box opened in `main.dart`: `Hive.openBox<String>(AppConfig.pendingCheckinsBox)`

## Auth Flow
1. Cold start → `AuthLoading` → read `flutter_secure_storage`
2. Token found → `AuthAuthenticated` → `/dashboard`
3. No token → `AuthUnauthenticated` → `/login`
4. Sign in → Google → idToken → POST backend → store JWT → `AuthAuthenticated`
5. 401 received → `authProvider.notifier.signOut()` → `/login`

## Storage Keys (flutter_secure_storage)
Defined in `AppConfig`: `auth_token`, `user_id`, `user_email`, `display_name`

## Adding a New Screen
1. Create `lib/features/<feature>/presentation/<screen>.dart`
2. Add route to `app_router.dart`
3. Add Riverpod provider if needed

## Platform Config
- Android: `android/app/google-services.json` (gitignored, see `.json.example`)
- iOS: `ios/Runner/GoogleService-Info.plist` (gitignored, see `.plist.example`)
- Backend `GOOGLE_CLIENT_ID` must match the **Web** client ID

## Run
```bash
cd mobile_app
flutter pub get
flutter run
```
