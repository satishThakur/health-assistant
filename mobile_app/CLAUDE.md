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
  health/                       # dashboard, check-in, trends screens
lib/core/
  network/api_client.dart       # Dio instance
  network/api_interceptor.dart  # injects Bearer token, signs out on 401
  routing/app_router.dart
  config/app_config.dart        # base URL, storage keys
```

## State Management
- **Riverpod** (`flutter_riverpod`) — all providers in `features/*/providers/`
- `authProvider` = `StateNotifierProvider<AuthNotifier, AuthState>`
- `authTokenProvider` = derived `Provider<String?>` — returns token if authenticated

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
