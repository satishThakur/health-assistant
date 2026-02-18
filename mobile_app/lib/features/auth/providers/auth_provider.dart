import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../data/auth_service.dart';
import '../domain/auth_state.dart';

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthService _authService;

  AuthNotifier(this._authService) : super(AuthLoading()) {
    _initialize();
  }

  Future<void> _initialize() async {
    try {
      final credentials = await _authService.getStoredCredentials();
      if (credentials != null) {
        state = AuthAuthenticated(
          token: credentials.token,
          userId: credentials.userId,
          email: credentials.email,
          displayName: credentials.displayName,
        );
      } else {
        state = AuthUnauthenticated();
      }
    } catch (_) {
      state = AuthUnauthenticated();
    }
  }

  Future<void> signInWithGoogle() async {
    state = AuthLoading();
    try {
      final credentials = await _authService.signInWithGoogle();
      state = AuthAuthenticated(
        token: credentials.token,
        userId: credentials.userId,
        email: credentials.email,
        displayName: credentials.displayName,
      );
    } catch (e) {
      state = AuthUnauthenticated();
      rethrow;
    }
  }

  Future<void> signOut() async {
    await _authService.signOut();
    state = AuthUnauthenticated();
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(authServiceProvider));
});

final authTokenProvider = Provider<String?>((ref) {
  final authState = ref.watch(authProvider);
  if (authState is AuthAuthenticated) {
    return authState.token;
  }
  return null;
});
