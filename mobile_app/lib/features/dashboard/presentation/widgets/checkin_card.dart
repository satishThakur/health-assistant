import 'package:flutter/material.dart';
import '../../../../core/config/theme.dart';
import '../../../checkin/domain/checkin_model.dart';

class CheckinCard extends StatelessWidget {
  final CheckinModel? checkin;

  const CheckinCard({super.key, this.checkin});

  @override
  Widget build(BuildContext context) {
    if (checkin == null) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              Icon(
                Icons.edit_calendar,
                size: 48,
                color: Colors.grey[400],
              ),
              const SizedBox(height: 12),
              Text(
                'No check-in yet today',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      color: Colors.grey[600],
                    ),
              ),
              const SizedBox(height: 8),
              Text(
                'Tap the + button to check-in',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: Colors.grey[500],
                    ),
              ),
            ],
          ),
        ),
      );
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.check_circle, color: Colors.green),
                const SizedBox(width: 8),
                Text(
                  'Today\'s Check-in',
                  style: Theme.of(context).textTheme.titleLarge,
                ),
              ],
            ),
            const SizedBox(height: 16),
            _buildMetricRow(
              context,
              '💪',
              'Energy',
              checkin!.energy,
              AppTheme.energyColor,
            ),
            const SizedBox(height: 12),
            _buildMetricRow(
              context,
              '😊',
              'Mood',
              checkin!.mood,
              AppTheme.moodColor,
            ),
            const SizedBox(height: 12),
            _buildMetricRow(
              context,
              '🎯',
              'Focus',
              checkin!.focus,
              AppTheme.focusColor,
            ),
            const SizedBox(height: 12),
            _buildMetricRow(
              context,
              '🏃',
              'Physical',
              checkin!.physical,
              AppTheme.physicalColor,
            ),
            if (checkin!.notes != null && checkin!.notes!.isNotEmpty) ...[
              const SizedBox(height: 16),
              const Divider(),
              const SizedBox(height: 12),
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('📝 ', style: TextStyle(fontSize: 16)),
                  Expanded(
                    child: Text(
                      checkin!.notes!,
                      style: Theme.of(context).textTheme.bodyMedium,
                    ),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildMetricRow(
    BuildContext context,
    String emoji,
    String label,
    int value,
    Color color,
  ) {
    return Row(
      children: [
        Text(emoji, style: const TextStyle(fontSize: 20)),
        const SizedBox(width: 8),
        SizedBox(
          width: 80,
          child: Text(
            label,
            style: Theme.of(context).textTheme.bodyLarge,
          ),
        ),
        Expanded(
          child: Stack(
            children: [
              Container(
                height: 8,
                decoration: BoxDecoration(
                  color: color.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(4),
                ),
              ),
              FractionallySizedBox(
                widthFactor: value / 10,
                child: Container(
                  height: 8,
                  decoration: BoxDecoration(
                    color: color,
                    borderRadius: BorderRadius.circular(4),
                  ),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(width: 12),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(
            color: color.withOpacity(0.2),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text(
            '$value/10',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: color,
            ),
          ),
        ),
      ],
    );
  }
}
