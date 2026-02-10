import 'package:flutter/material.dart';
import '../../../dashboard/domain/dashboard_model.dart';

class InsightCard extends StatelessWidget {
  final CorrelationInsight insight;

  const InsightCard({super.key, required this.insight});

  @override
  Widget build(BuildContext context) {
    final icon = _getIcon();
    final color = _getColor();

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: color.withOpacity(0.2),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(icon, color: color, size: 28),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Text(
                    _getTitle(),
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Text(
              insight.description,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                _buildMetric(
                  context,
                  'Confidence',
                  '${(insight.confidence * 100).round()}%',
                ),
                const SizedBox(width: 16),
                _buildMetric(
                  context,
                  'Sample Size',
                  '${insight.sampleSize} days',
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMetric(BuildContext context, String label, String value) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: Colors.grey[600],
              ),
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
      ],
    );
  }

  String _getTitle() {
    switch (insight.type) {
      case 'sleep_energy':
        return '😴 Sleep & Energy';
      case 'activity_mood':
        return '🏃 Activity & Mood';
      case 'sleep_focus':
        return '😴 Sleep & Focus';
      default:
        return 'Insight';
    }
  }

  IconData _getIcon() {
    switch (insight.type) {
      case 'sleep_energy':
        return Icons.bedtime;
      case 'activity_mood':
        return Icons.directions_run;
      case 'sleep_focus':
        return Icons.psychology;
      default:
        return Icons.lightbulb;
    }
  }

  Color _getColor() {
    switch (insight.type) {
      case 'sleep_energy':
        return const Color(0xFF7C3AED);
      case 'activity_mood':
        return const Color(0xFF3B82F6);
      case 'sleep_focus':
        return const Color(0xFF10B981);
      default:
        return Colors.blue;
    }
  }
}
