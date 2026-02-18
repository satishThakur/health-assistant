import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:google_sign_in/google_sign_in.dart';

import '../../../core/config/app_config.dart';
import '../../../core/network/api_client.dart';

class AuthResponse {
  final String token;
  final String userId;
  final String email;
  final String displayName;

  AuthResponse({
    required this.token,
    required this.userId,
    required this.email,
    required this.displayName,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      token: json['token'] as String,
      userId: json['user_id'] as String,
      email: json['email'] as String,
      displayName: (json['display_name'] as String?) ?? '',
    );
  }
}

class AuthService {
  final ApiClient _apiClient;
  final FlutterSecureStorage _secureStorage;
  final GoogleSignIn _googleSignIn;

  AuthService(this._apiClient, this._secureStorage, this._googleSignIn);

  Future<AuthResponse> signInWithGoogle() async {
    final googleUser = await _googleSignIn.signIn();
    if (googleUser == null) {
      throw Exception('Google sign-in was cancelled');
    }

    final googleAuth = await googleUser.authentication;
    final idToken = googleAuth.idToken;
    if (idToken == null) {
      throw Exception('Failed to obtain Google ID token');
    }

    final response = await _apiClient.post<Map<String, dynamic>>(
      '/api/v1/auth/google',
      data: {'id_token': idToken},
    );

    final authResponse = AuthResponse.fromJson(response.data!);

    await Future.wait([
      _secureStorage.write(key: AppConfig.tokenKey, value: authResponse.token),
      _secureStorage.write(key: AppConfig.userIdKey, value: authResponse.userId),
      _secureStorage.write(key: 'email', value: authResponse.email),
      _secureStorage.write(key: 'display_name', value: authResponse.displayName),
    ]);

    return authResponse;
  }

  Future<void> signOut() async {
    await Future.wait([
      _googleSignIn.signOut(),
      _secureStorage.deleteAll(),
    ]);
  }

  Future<AuthResponse?> getStoredCredentials() async {
    final token = await _secureStorage.read(key: AppConfig.tokenKey);
    if (token == null) return null;

    final userId = await _secureStorage.read(key: AppConfig.userIdKey);
    final email = await _secureStorage.read(key: 'email');
    final displayName = await _secureStorage.read(key: 'display_name');

    if (userId == null || email == null) return null;

    return AuthResponse(
      token: token,
      userId: userId,
      email: email,
      displayName: displayName ?? '',
    );
  }

  Future<String?> getStoredToken() async {
    return _secureStorage.read(key: AppConfig.tokenKey);
  }
}

final _secureStorageProvider = Provider<FlutterSecureStorage>((ref) {
  return const FlutterSecureStorage(
    aOptions: AndroidOptions(encryptedSharedPreferences: true),
    iOptions: IOSOptions(
      accessibility: KeychainAccessibility.first_unlock_this_device,
    ),
  );
});

final _googleSignInProvider = Provider<GoogleSignIn>((ref) {
  return GoogleSignIn();
});

final authServiceProvider = Provider<AuthService>((ref) {
  return AuthService(
    ref.watch(apiClientProvider),
    ref.watch(_secureStorageProvider),
    ref.watch(_googleSignInProvider),
  );
});
