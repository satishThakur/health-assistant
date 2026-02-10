import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import '../../../../core/config/theme.dart';
import '../../domain/dashboard_model.dart';

class TrendChart extends StatefulWidget {
  final List<TrendData> trends;

  const TrendChart({super.key, required this.trends});

  @override
  State<TrendChart> createState() => _TrendChartState();
}

class _TrendChartState extends State<TrendChart> {
  String selectedMetric = 'energy';

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Metric Selector
            SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children: [
                  _buildMetricChip('Energy', 'energy', AppTheme.energyColor),
                  const SizedBox(width: 8),
                  _buildMetricChip('Mood', 'mood', AppTheme.moodColor),
                  const SizedBox(width: 8),
                  _buildMetricChip('Focus', 'focus', AppTheme.focusColor),
                  const SizedBox(width: 8),
                  _buildMetricChip(
                      'Physical', 'physical', AppTheme.physicalColor),
                  const SizedBox(width: 8),
                  _buildMetricChip('Sleep', 'sleep', AppTheme.sleepColor),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // Chart
            SizedBox(
              height: 200,
              child: _buildChart(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMetricChip(String label, String value, Color color) {
    final isSelected = selectedMetric == value;

    return FilterChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (selected) {
        setState(() {
          selectedMetric = value;
        });
      },
      backgroundColor: color.withOpacity(0.1),
      selectedColor: color.withOpacity(0.3),
      checkmarkColor: color,
      labelStyle: TextStyle(
        color: isSelected ? color : Colors.grey[700],
        fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
      ),
    );
  }

  Widget _buildChart() {
    final spots = _getChartSpots();

    if (spots.isEmpty) {
      return Center(
        child: Text(
          'No data available for this metric',
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: Colors.grey[600],
              ),
        ),
      );
    }

    return LineChart(
      LineChartData(
        gridData: FlGridData(
          show: true,
          drawVerticalLine: false,
          horizontalInterval: 2,
        ),
        titlesData: FlTitlesData(
          leftTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              interval: 2,
              reservedSize: 32,
            ),
          ),
          bottomTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              reservedSize: 32,
              getTitlesWidget: (value, meta) {
                if (value.toInt() >= 0 &&
                    value.toInt() < widget.trends.length) {
                  final date = DateTime.parse(
                      widget.trends[value.toInt()].date);
                  return Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: Text(
                      ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'][
                          date.weekday - 1],
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  );
                }
                return const Text('');
              },
            ),
          ),
          topTitles: const AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
          rightTitles: const AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
        ),
        borderData: FlBorderData(show: false),
        minY: selectedMetric == 'sleep' ? 0 : 0,
        maxY: selectedMetric == 'sleep' ? 10 : 10,
        lineBarsData: [
          LineChartBarData(
            spots: spots,
            isCurved: true,
            color: _getMetricColor(),
            barWidth: 3,
            isStrokeCapRound: true,
            dotData: const FlDotData(show: true),
            belowBarData: BarAreaData(
              show: true,
              color: _getMetricColor().withOpacity(0.2),
            ),
          ),
        ],
      ),
    );
  }

  List<FlSpot> _getChartSpots() {
    final spots = <FlSpot>[];

    for (var i = 0; i < widget.trends.length; i++) {
      final trend = widget.trends[i];
      double? value;

      switch (selectedMetric) {
        case 'energy':
          value = trend.checkin?.energy.toDouble();
          break;
        case 'mood':
          value = trend.checkin?.mood.toDouble();
          break;
        case 'focus':
          value = trend.checkin?.focus.toDouble();
          break;
        case 'physical':
          value = trend.checkin?.physical.toDouble();
          break;
        case 'sleep':
          value = trend.sleep != null
              ? (trend.sleep!.durationMinutes / 60.0)
              : null;
          break;
      }

      if (value != null) {
        spots.add(FlSpot(i.toDouble(), value));
      }
    }

    return spots;
  }

  Color _getMetricColor() {
    switch (selectedMetric) {
      case 'energy':
        return AppTheme.energyColor;
      case 'mood':
        return AppTheme.moodColor;
      case 'focus':
        return AppTheme.focusColor;
      case 'physical':
        return AppTheme.physicalColor;
      case 'sleep':
        return AppTheme.sleepColor;
      default:
        return AppTheme.primaryColor;
    }
  }
}
