import 'package:flutter/material.dart';
import '../../domain/dashboard_model.dart';

class SleepCard extends StatelessWidget {
  final SleepData sleep;

  const SleepCard({super.key, required this.sleep});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Text('😴', style: TextStyle(fontSize: 28)),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Sleep',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    Text(
                      '${sleep.durationHours.toStringAsFixed(1)} hours · Score: ${sleep.sleepScore}',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: Colors.grey[600],
                          ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 20),
            _buildSleepStage(
              context,
              'Deep',
              sleep.deepSleepMinutes,
              sleep.durationMinutes,
              const Color(0xFF7C3AED),
            ),
            const SizedBox(height: 12),
            _buildSleepStage(
              context,
              'Light',
              sleep.lightSleepMinutes,
              sleep.durationMinutes,
              const Color(0xFFA78BFA),
            ),
            const SizedBox(height: 12),
            _buildSleepStage(
              context,
              'REM',
              sleep.remSleepMinutes,
              sleep.durationMinutes,
              const Color(0xFF6366F1),
            ),
            if (sleep.awakeMinutes > 0) ...[
              const SizedBox(height: 12),
              _buildSleepStage(
                context,
                'Awake',
                sleep.awakeMinutes,
                sleep.durationMinutes,
                Colors.grey,
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildSleepStage(
    BuildContext context,
    String label,
    int minutes,
    int totalMinutes,
    Color color,
  ) {
    final hours = (minutes / 60).toStringAsFixed(1);
    final percentage = (minutes / totalMinutes * 100).round();

    return Column(
      children: [
        Row(
          children: [
            SizedBox(
              width: 60,
              child: Text(
                label,
                style: Theme.of(context).textTheme.bodyMedium,
              ),
            ),
            Expanded(
              child: Stack(
                children: [
                  Container(
                    height: 20,
                    decoration: BoxDecoration(
                      color: color.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(10),
                    ),
                  ),
                  FractionallySizedBox(
                    widthFactor: percentage / 100,
                    child: Container(
                      height: 20,
                      decoration: BoxDecoration(
                        color: color,
                        borderRadius: BorderRadius.circular(10),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            SizedBox(
              width: 60,
              child: Text(
                '${hours}h',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                textAlign: TextAlign.end,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
