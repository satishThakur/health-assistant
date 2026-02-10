import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/api_client.dart';
import '../domain/dashboard_model.dart';
import 'dashboard_api.dart';

final dashboardRepositoryProvider = Provider<DashboardRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return DashboardRepository(DashboardApi(apiClient));
});

class DashboardRepository {
  final DashboardApi _api;

  DashboardRepository(this._api);

  Future<DashboardData> getTodayDashboard() async {
    return await _api.getTodayDashboard();
  }

  Future<List<TrendData>> getWeekTrends() async {
    return await _api.getWeekTrends();
  }

  Future<List<CorrelationInsight>> getCorrelations({int days = 30}) async {
    return await _api.getCorrelations(days: days);
  }
}
