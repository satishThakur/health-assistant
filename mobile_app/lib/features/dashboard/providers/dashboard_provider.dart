import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/dashboard_repository.dart';
import '../domain/dashboard_model.dart';

// Today's dashboard provider
final todayDashboardProvider = FutureProvider<DashboardData>((ref) async {
  final repository = ref.watch(dashboardRepositoryProvider);
  return await repository.getTodayDashboard();
});

// Week trends provider
final weekTrendsProvider = FutureProvider<List<TrendData>>((ref) async {
  final repository = ref.watch(dashboardRepositoryProvider);
  return await repository.getWeekTrends();
});

// Correlations provider
final correlationsProvider =
    FutureProvider.family<List<CorrelationInsight>, int>((ref, days) async {
  final repository = ref.watch(dashboardRepositoryProvider);
  return await repository.getCorrelations(days: days);
});
