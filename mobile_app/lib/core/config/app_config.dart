class AppConfig {
  // API Configuration
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8083',
  );

  // Timeouts
  static const Duration connectTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);

  // Storage Keys
  static const String tokenKey = 'auth_token';
  static const String userIdKey = 'user_id';
  static const String lastCheckinKey = 'last_checkin_date';

  // Default Values
  static const String defaultUserId = '00000000-0000-0000-0000-000000000001';

  // Feature Flags
  static const bool enableOfflineMode = true;
  static const bool enableNotifications = true;

  // Correlation Settings
  static const int minSamplesForInsight = 5;
  static const double minImprovementPercent = 5.0;
  static const int defaultCorrelationDays = 30;
}
