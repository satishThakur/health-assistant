import 'package:flutter/material.dart';
import '../../domain/dashboard_model.dart';

class ActivityStatsCard extends StatelessWidget {
  final DailyStatsData stats;

  const ActivityStatsCard({
    super.key,
    required this.stats,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  Icons.directions_walk,
                  color: Theme.of(context).primaryColor,
                ),
                const SizedBox(width: 8),
                Text(
                  'Daily Activity',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // Steps and Calories Row
            Row(
              children: [
                Expanded(
                  child: _buildStat(
                    context,
                    icon: Icons.directions_walk,
                    label: 'Steps',
                    value: _formatNumber(stats.steps),
                    color: Colors.blue,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStat(
                    context,
                    icon: Icons.local_fire_department,
                    label: 'Calories',
                    value: _formatNumber(stats.calories),
                    color: Colors.orange,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),

            // Distance and Heart Rate Row
            Row(
              children: [
                Expanded(
                  child: _buildStat(
                    context,
                    icon: Icons.straighten,
                    label: 'Distance',
                    value: '${stats.distanceKm.toStringAsFixed(1)} km',
                    color: Colors.green,
                  ),
                ),
                const SizedBox(width: 12),
                if (stats.restingHeartRate != null)
                  Expanded(
                    child: _buildStat(
                      context,
                      icon: Icons.favorite,
                      label: 'Resting HR',
                      value: '${stats.restingHeartRate} bpm',
                      color: Colors.red,
                    ),
                  )
                else
                  const Expanded(child: SizedBox()),
              ],
            ),

            // Heart Rate Details (if available)
            if (stats.minHeartRate != null && stats.maxHeartRate != null) ...[
              const SizedBox(height: 12),
              Divider(color: Colors.grey[300]),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  _buildSmallStat(
                    context,
                    label: 'Min HR',
                    value: '${stats.minHeartRate}',
                  ),
                  _buildSmallStat(
                    context,
                    label: 'Max HR',
                    value: '${stats.maxHeartRate}',
                  ),
                  if (stats.moderateIntensityMinutes != null)
                    _buildSmallStat(
                      context,
                      label: 'Active',
                      value: '${stats.moderateIntensityMinutes}m',
                    ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildStat(
    BuildContext context, {
    required IconData icon,
    required String label,
    required String value,
    required Color color,
  }) {
    return Column(
      children: [
        Icon(icon, color: color, size: 24),
        const SizedBox(height: 4),
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: Colors.grey[600],
              ),
        ),
        const SizedBox(height: 2),
        Text(
          value,
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
      ],
    );
  }

  Widget _buildSmallStat(
    BuildContext context, {
    required String label,
    required String value,
  }) {
    return Column(
      children: [
        Text(
          value,
          style: Theme.of(context).textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: Colors.grey[600],
              ),
        ),
      ],
    );
  }

  String _formatNumber(int number) {
    if (number >= 1000) {
      return '${(number / 1000).toStringAsFixed(1)}k';
    }
    return number.toString();
  }
}
