import 'package:json_annotation/json_annotation.dart';

part 'checkin_model.g.dart';

@JsonSerializable()
class CheckinModel {
  final int energy;
  final int mood;
  final int focus;
  final int physical;
  final String? notes;

  CheckinModel({
    required this.energy,
    required this.mood,
    required this.focus,
    required this.physical,
    this.notes,
  });

  factory CheckinModel.fromJson(Map<String, dynamic> json) =>
      _$CheckinModelFromJson(json);

  Map<String, dynamic> toJson() => _$CheckinModelToJson(this);

  CheckinModel copyWith({
    int? energy,
    int? mood,
    int? focus,
    int? physical,
    String? notes,
  }) {
    return CheckinModel(
      energy: energy ?? this.energy,
      mood: mood ?? this.mood,
      focus: focus ?? this.focus,
      physical: physical ?? this.physical,
      notes: notes ?? this.notes,
    );
  }
}

@JsonSerializable()
class CheckinResponse {
  final String status;
  final String action;
  final DateTime timestamp;
  final CheckinModel data;

  CheckinResponse({
    required this.status,
    required this.action,
    required this.timestamp,
    required this.data,
  });

  factory CheckinResponse.fromJson(Map<String, dynamic> json) =>
      _$CheckinResponseFromJson(json);

  Map<String, dynamic> toJson() => _$CheckinResponseToJson(this);
}

@JsonSerializable()
class CheckinHistoryItem {
  final String date;
  final CheckinModel checkin;

  CheckinHistoryItem({
    required this.date,
    required this.checkin,
  });

  factory CheckinHistoryItem.fromJson(Map<String, dynamic> json) =>
      _$CheckinHistoryItemFromJson(json);

  Map<String, dynamic> toJson() => _$CheckinHistoryItemToJson(this);
}
