import 'package:json_annotation/json_annotation.dart';
import '../../checkin/domain/checkin_model.dart';

part 'dashboard_model.g.dart';

@JsonSerializable()
class DashboardData {
  final CheckinModel? checkin;
  final GarminData? garmin;

  DashboardData({
    this.checkin,
    this.garmin,
  });

  factory DashboardData.fromJson(Map<String, dynamic> json) =>
      _$DashboardDataFromJson(json);

  Map<String, dynamic> toJson() => _$DashboardDataToJson(this);
}

@JsonSerializable()
class GarminData {
  final SleepData? sleep;
  final ActivityData? activity;
  final HRVData? hrv;
  final StressData? stress;

  @JsonKey(name: 'daily_stats')
  final DailyStatsData? dailyStats;

  @JsonKey(name: 'body_battery')
  final BodyBatteryData? bodyBattery;

  GarminData({
    this.sleep,
    this.activity,
    this.hrv,
    this.stress,
    this.dailyStats,
    this.bodyBattery,
  });

  factory GarminData.fromJson(Map<String, dynamic> json) =>
      _$GarminDataFromJson(json);

  Map<String, dynamic> toJson() => _$GarminDataToJson(this);
}

@JsonSerializable()
class SleepData {
  @JsonKey(name: 'duration_minutes')
  final int durationMinutes;

  @JsonKey(name: 'deep_sleep_minutes')
  final int deepSleepMinutes;

  @JsonKey(name: 'light_sleep_minutes')
  final int lightSleepMinutes;

  @JsonKey(name: 'rem_sleep_minutes')
  final int remSleepMinutes;

  @JsonKey(name: 'awake_minutes')
  final int awakeMinutes;

  @JsonKey(name: 'sleep_score')
  final int sleepScore;

  @JsonKey(name: 'hrv_avg')
  final double? hrvAvg;

  SleepData({
    required this.durationMinutes,
    required this.deepSleepMinutes,
    required this.lightSleepMinutes,
    required this.remSleepMinutes,
    required this.awakeMinutes,
    required this.sleepScore,
    this.hrvAvg,
  });

  double get durationHours => durationMinutes / 60.0;

  factory SleepData.fromJson(Map<String, dynamic> json) =>
      _$SleepDataFromJson(json);

  Map<String, dynamic> toJson() => _$SleepDataToJson(this);
}

@JsonSerializable()
class ActivityData {
  @JsonKey(name: 'activity_type')
  final String activityType;

  @JsonKey(name: 'duration_minutes')
  final int durationMinutes;

  final int calories;

  @JsonKey(name: 'avg_hr')
  final int? avgHr;

  @JsonKey(name: 'max_hr')
  final int? maxHr;

  final double? distance;

  ActivityData({
    required this.activityType,
    required this.durationMinutes,
    required this.calories,
    this.avgHr,
    this.maxHr,
    this.distance,
  });

  factory ActivityData.fromJson(Map<String, dynamic> json) =>
      _$ActivityDataFromJson(json);

  Map<String, dynamic> toJson() => _$ActivityDataToJson(this);
}

@JsonSerializable()
class HRVData {
  final double average;

  HRVData({required this.average});

  factory HRVData.fromJson(Map<String, dynamic> json) =>
      _$HRVDataFromJson(json);

  Map<String, dynamic> toJson() => _$HRVDataToJson(this);
}

@JsonSerializable()
class StressData {
  final int average;
  final String level;

  StressData({
    required this.average,
    required this.level,
  });

  factory StressData.fromJson(Map<String, dynamic> json) =>
      _$StressDataFromJson(json);

  Map<String, dynamic> toJson() => _$StressDataToJson(this);
}

@JsonSerializable()
class DailyStatsData {
  final int steps;
  final int calories;

  @JsonKey(name: 'distance_meters')
  final int distanceMeters;

  @JsonKey(name: 'active_calories')
  final int? activeCalories;

  @JsonKey(name: 'bmr_calories')
  final int? bmrCalories;

  @JsonKey(name: 'min_heart_rate')
  final int? minHeartRate;

  @JsonKey(name: 'max_heart_rate')
  final int? maxHeartRate;

  @JsonKey(name: 'resting_heart_rate')
  final int? restingHeartRate;

  @JsonKey(name: 'moderate_intensity_minutes')
  final int? moderateIntensityMinutes;

  @JsonKey(name: 'vigorous_intensity_minutes')
  final int? vigorousIntensityMinutes;

  DailyStatsData({
    required this.steps,
    required this.calories,
    required this.distanceMeters,
    this.activeCalories,
    this.bmrCalories,
    this.minHeartRate,
    this.maxHeartRate,
    this.restingHeartRate,
    this.moderateIntensityMinutes,
    this.vigorousIntensityMinutes,
  });

  double get distanceKm => distanceMeters / 1000.0;

  factory DailyStatsData.fromJson(Map<String, dynamic> json) =>
      _$DailyStatsDataFromJson(json);

  Map<String, dynamic> toJson() => _$DailyStatsDataToJson(this);
}

@JsonSerializable()
class BodyBatteryData {
  final int charged;
  final int drained;

  @JsonKey(name: 'highest_value')
  final int? highestValue;

  @JsonKey(name: 'lowest_value')
  final int? lowestValue;

  BodyBatteryData({
    required this.charged,
    required this.drained,
    this.highestValue,
    this.lowestValue,
  });

  int get netEnergy => charged - drained;

  factory BodyBatteryData.fromJson(Map<String, dynamic> json) =>
      _$BodyBatteryDataFromJson(json);

  Map<String, dynamic> toJson() => _$BodyBatteryDataToJson(this);
}

@JsonSerializable()
class TrendData {
  final String date;
  final CheckinModel? checkin;
  final SleepData? sleep;
  final ActivityData? activity;

  TrendData({
    required this.date,
    this.checkin,
    this.sleep,
    this.activity,
  });

  factory TrendData.fromJson(Map<String, dynamic> json) =>
      _$TrendDataFromJson(json);

  Map<String, dynamic> toJson() => _$TrendDataToJson(this);
}

@JsonSerializable()
class CorrelationInsight {
  final String type;
  final String description;
  final double confidence;

  @JsonKey(name: 'sample_size')
  final int sampleSize;

  final Map<String, dynamic> details;

  CorrelationInsight({
    required this.type,
    required this.description,
    required this.confidence,
    required this.sampleSize,
    required this.details,
  });

  factory CorrelationInsight.fromJson(Map<String, dynamic> json) =>
      _$CorrelationInsightFromJson(json);

  Map<String, dynamic> toJson() => _$CorrelationInsightToJson(this);
}
