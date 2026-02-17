import 'package:flutter/material.dart';
import '../../domain/dashboard_model.dart';

class BodyBatteryCard extends StatelessWidget {
  final BodyBatteryData bodyBattery;

  const BodyBatteryCard({
    super.key,
    required this.bodyBattery,
  });

  @override
  Widget build(BuildContext context) {
    final netEnergy = bodyBattery.netEnergy;
    final isPositive = netEnergy >= 0;
    final batteryColor = _getBatteryColor(netEnergy);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  Icons.battery_charging_full,
                  color: batteryColor,
                ),
                const SizedBox(width: 8),
                Text(
                  'Body Battery',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // Net Energy Display
            Center(
              child: Column(
                children: [
                  Text(
                    isPositive ? '+$netEnergy' : '$netEnergy',
                    style: Theme.of(context).textTheme.displayMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: batteryColor,
                        ),
                  ),
                  Text(
                    'Net Energy',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: Colors.grey[600],
                        ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // Visual Bar
            Container(
              height: 8,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(4),
                gradient: LinearGradient(
                  colors: [
                    Colors.red.shade300,
                    Colors.yellow.shade300,
                    Colors.green.shade300,
                  ],
                  stops: const [0.0, 0.5, 1.0],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Charged and Drained Row
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildEnergyMetric(
                  context,
                  icon: Icons.battery_charging_full,
                  label: 'Charged',
                  value: bodyBattery.charged,
                  color: Colors.green,
                ),
                Container(
                  height: 40,
                  width: 1,
                  color: Colors.grey[300],
                ),
                _buildEnergyMetric(
                  context,
                  icon: Icons.battery_alert,
                  label: 'Drained',
                  value: bodyBattery.drained,
                  color: Colors.orange,
                ),
              ],
            ),

            // High/Low values if available
            if (bodyBattery.highestValue != null &&
                bodyBattery.lowestValue != null) ...[
              const SizedBox(height: 12),
              Divider(color: Colors.grey[300]),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  _buildSmallMetric(
                    context,
                    label: 'Highest',
                    value: bodyBattery.highestValue!,
                  ),
                  _buildSmallMetric(
                    context,
                    label: 'Lowest',
                    value: bodyBattery.lowestValue!,
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildEnergyMetric(
    BuildContext context, {
    required IconData icon,
    required String label,
    required int value,
    required Color color,
  }) {
    return Column(
      children: [
        Icon(icon, color: color, size: 28),
        const SizedBox(height: 4),
        Text(
          value.toString(),
          style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: color,
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

  Widget _buildSmallMetric(
    BuildContext context, {
    required String label,
    required int value,
  }) {
    return Column(
      children: [
        Text(
          value.toString(),
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
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

  Color _getBatteryColor(int netEnergy) {
    if (netEnergy > 20) {
      return Colors.green;
    } else if (netEnergy > 0) {
      return Colors.lightGreen;
    } else if (netEnergy > -20) {
      return Colors.orange;
    } else {
      return Colors.red;
    }
  }
}
