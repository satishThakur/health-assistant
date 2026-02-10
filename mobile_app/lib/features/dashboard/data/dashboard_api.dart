import '../../../core/network/api_client.dart';
import '../../../core/network/api_endpoints.dart';
import '../domain/dashboard_model.dart';

class DashboardApi {
  final ApiClient _client;

  DashboardApi(this._client);

  Future<DashboardData> getTodayDashboard() async {
    try {
      final response = await _client.get(ApiEndpoints.dashboardToday);
      return DashboardData.fromJson(response.data['data']);
    } catch (e) {
      throw Exception('Failed to fetch dashboard: $e');
    }
  }

  Future<List<TrendData>> getWeekTrends() async {
    try {
      final response = await _client.get(ApiEndpoints.trendsWeek);
      return (response.data['trends'] as List)
          .map((item) => TrendData.fromJson(item))
          .toList();
    } catch (e) {
      throw Exception('Failed to fetch trends: $e');
    }
  }

  Future<List<CorrelationInsight>> getCorrelations({int days = 30}) async {
    try {
      final response = await _client.get(
        ApiEndpoints.correlations,
        queryParameters: {'days': days},
      );

      return (response.data['correlations'] as List)
          .map((item) => CorrelationInsight.fromJson(item))
          .toList();
    } catch (e) {
      throw Exception('Failed to fetch correlations: $e');
    }
  }
}
