import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../shared/widgets/loading_indicator.dart';
import '../../../shared/widgets/error_view.dart';
import '../domain/dashboard_model.dart';
import '../providers/dashboard_provider.dart';
import 'widgets/trend_chart.dart';

class TrendsScreen extends ConsumerWidget {
  const TrendsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final trendsAsync = ref.watch(weekTrendsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Your Week'),
      ),
      body: trendsAsync.when(
        data: (trends) => _buildTrends(context, trends),
        loading: () => const AppLoadingIndicator(),
        error: (error, stack) => AppErrorView(
          error: error.toString(),
          onRetry: () => ref.invalidate(weekTrendsProvider),
        ),
      ),
    );
  }

  Widget _buildTrends(BuildContext context, List<TrendData> trends) {
    if (trends.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.timeline, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              'No trend data available yet',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                    color: Colors.grey[600],
                  ),
            ),
            const SizedBox(height: 8),
            Text(
              'Check in daily to see your trends',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey[500],
                  ),
            ),
          ],
        ),
      );
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '7-Day Trends',
            style: Theme.of(context).textTheme.displaySmall,
          ),
          const SizedBox(height: 24),
          TrendChart(trends: trends),
          const SizedBox(height: 24),
          _buildInsightsSummary(context, trends),
        ],
      ),
    );
  }

  Widget _buildInsightsSummary(BuildContext context, List<TrendData> trends) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Text('💡', style: TextStyle(fontSize: 24)),
                const SizedBox(width: 8),
                Text(
                  'Quick Insights',
                  style: Theme.of(context).textTheme.titleLarge,
                ),
              ],
            ),
            const SizedBox(height: 16),
            _buildInsightItem(
              context,
              'Days tracked',
              '${trends.length} days',
            ),
            const SizedBox(height: 12),
            _buildInsightItem(
              context,
              'Consistency',
              trends.length >= 7 ? 'Excellent! 🔥' : 'Keep going! 💪',
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInsightItem(BuildContext context, String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodyLarge,
        ),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
      ],
    );
  }
}
