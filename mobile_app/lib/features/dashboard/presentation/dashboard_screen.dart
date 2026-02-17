import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import '../../../shared/widgets/loading_indicator.dart';
import '../../../shared/widgets/error_view.dart';
import '../domain/dashboard_model.dart';
import '../providers/dashboard_provider.dart';
import 'widgets/checkin_card.dart';
import 'widgets/sleep_card.dart';
import 'widgets/metric_card.dart';
import 'widgets/activity_stats_card.dart';
import 'widgets/body_battery_card.dart';

class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboardAsync = ref.watch(todayDashboardProvider);

    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Dashboard'),
            Text(
              DateFormat('EEEE, MMM d').format(DateTime.now()),
              style: Theme.of(context).textTheme.bodySmall,
            ),
          ],
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(todayDashboardProvider);
        },
        child: dashboardAsync.when(
          data: (dashboard) => _buildDashboard(context, ref, dashboard),
          loading: () => const AppLoadingIndicator(),
          error: (error, stack) => AppErrorView(
            error: error.toString(),
            onRetry: () => ref.invalidate(todayDashboardProvider),
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.push('/checkin'),
        icon: const Icon(Icons.add),
        label: const Text('Check-in'),
      ),
    );
  }

  Widget _buildDashboard(
    BuildContext context,
    WidgetRef ref,
    DashboardData dashboard,
  ) {
    return SingleChildScrollView(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Greeting
          _buildGreeting(context),
          const SizedBox(height: 24),

          // Today's Check-in
          CheckinCard(checkin: dashboard.checkin),
          const SizedBox(height: 16),

          // Last Night Sleep
          if (dashboard.garmin?.sleep != null) ...[
            Text(
              'Last Night',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 12),
            SleepCard(sleep: dashboard.garmin!.sleep!),
            const SizedBox(height: 16),
          ],

          // Today's Activity Stats
          if (dashboard.garmin?.dailyStats != null) ...[
            Text(
              'Today\'s Activity',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 12),
            ActivityStatsCard(stats: dashboard.garmin!.dailyStats!),
            const SizedBox(height: 16),
          ],

          // Body Battery
          if (dashboard.garmin?.bodyBattery != null) ...[
            BodyBatteryCard(bodyBattery: dashboard.garmin!.bodyBattery!),
            const SizedBox(height: 16),
          ],

          // HRV and Stress Metrics
          if (dashboard.garmin?.hrv != null ||
              dashboard.garmin?.stress != null) ...[
            Text(
              'Recovery Metrics',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                if (dashboard.garmin!.hrv != null)
                  Expanded(
                    child: MetricCard(
                      icon: Icons.favorite,
                      label: 'HRV',
                      value: '${dashboard.garmin!.hrv!.average.round()} ms',
                      color: Colors.red,
                    ),
                  ),
                if (dashboard.garmin!.hrv != null &&
                    dashboard.garmin!.stress != null)
                  const SizedBox(width: 12),
                if (dashboard.garmin!.stress != null)
                  Expanded(
                    child: MetricCard(
                      icon: Icons.psychology,
                      label: 'Stress',
                      value:
                          '${dashboard.garmin!.stress!.average} (${dashboard.garmin!.stress!.level})',
                      color: Colors.orange,
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 16),
          ],

          // Navigation Cards
          Row(
            children: [
              Expanded(
                child: _buildNavigationCard(
                  context,
                  icon: Icons.trending_up,
                  label: 'View Trends',
                  onTap: () => context.push('/trends'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildNavigationCard(
                  context,
                  icon: Icons.lightbulb,
                  label: 'Insights',
                  onTap: () => context.push('/insights'),
                ),
              ),
            ],
          ),
          const SizedBox(height: 80), // Space for FAB
        ],
      ),
    );
  }

  Widget _buildGreeting(BuildContext context) {
    final hour = DateTime.now().hour;
    String greeting;
    String emoji;

    if (hour < 12) {
      greeting = 'Good morning';
      emoji = '☀️';
    } else if (hour < 17) {
      greeting = 'Good afternoon';
      emoji = '🌤️';
    } else {
      greeting = 'Good evening';
      emoji = '🌙';
    }

    return Row(
      children: [
        Text(
          '$emoji  ',
          style: const TextStyle(fontSize: 32),
        ),
        Text(
          greeting,
          style: Theme.of(context).textTheme.displaySmall,
        ),
      ],
    );
  }

  Widget _buildNavigationCard(
    BuildContext context, {
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              Icon(icon, size: 32, color: Theme.of(context).primaryColor),
              const SizedBox(height: 8),
              Text(
                label,
                style: Theme.of(context).textTheme.titleMedium,
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
