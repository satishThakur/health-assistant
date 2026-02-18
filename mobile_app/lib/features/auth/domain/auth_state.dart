sealed class AuthState {}

class AuthLoading extends AuthState {}

class AuthUnauthenticated extends AuthState {}

class AuthAuthenticated extends AuthState {
  final String token;
  final String userId;
  final String email;
  final String displayName;

  AuthAuthenticated({
    required this.token,
    required this.userId,
    required this.email,
    required this.displayName,
  });
}
