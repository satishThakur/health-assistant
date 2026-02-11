import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../shared/widgets/loading_indicator.dart';
import '../../../shared/widgets/error_view.dart';
import '../../dashboard/domain/dashboard_model.dart';
import '../../dashboard/providers/dashboard_provider.dart';
import 'widgets/insight_card.dart';

class InsightsScreen extends ConsumerWidget {
  const InsightsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final correlationsAsync = ref.watch(correlationsProvider(30));

    return Scaffold(
      appBar: AppBar(
        title: const Text('Personalized Insights'),
      ),
      body: correlationsAsync.when(
        data: (correlations) => _buildInsights(context, correlations),
        loading: () => const AppLoadingIndicator(),
        error: (error, stack) => AppErrorView(
          error: error.toString(),
          onRetry: () => ref.invalidate(correlationsProvider(30)),
        ),
      ),
    );
  }

  Widget _buildInsights(BuildContext context, List<CorrelationInsight> correlations) {
    if (correlations.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.lightbulb_outline, size: 64, color: Colors.grey[400]),
              const SizedBox(height: 16),
              Text(
                'Not enough data yet',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      color: Colors.grey[600],
                    ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 8),
              Text(
                'Keep checking in daily and syncing your Garmin data. '
                'Insights will appear after you have at least 30 days of data.',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Colors.grey[500],
                    ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      );
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Based on ${30} days of data',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.grey[600],
                ),
          ),
          const SizedBox(height: 16),
          ...correlations.map((insight) => Padding(
                padding: const EdgeInsets.only(bottom: 16),
                child: InsightCard(insight: insight),
              )),
        ],
      ),
    );
  }
}
