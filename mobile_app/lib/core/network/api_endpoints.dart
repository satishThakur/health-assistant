class ApiEndpoints {
  // Base paths
  static const String apiV1 = '/api/v1';

  // Check-in endpoints
  static const String checkin = '$apiV1/checkin';
  static const String checkinLatest = '$apiV1/checkin/latest';
  static const String checkinHistory = '$apiV1/checkin/history';

  // Dashboard endpoints
  static const String dashboardToday = '$apiV1/dashboard/today';
  static const String trendsWeek = '$apiV1/trends/week';

  // Insights endpoints
  static const String correlations = '$apiV1/insights/correlations';

  // Health check
  static const String health = '/health';
}
