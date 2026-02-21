import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../../core/config/app_config.dart';
import '../../../core/notifications/notification_service.dart';

class NotificationTimeNotifier extends StateNotifier<TimeOfDay> {
  NotificationTimeNotifier(this._ref)
      : super(const TimeOfDay(hour: 21, minute: 0)) {
    _load();
  }

  final Ref _ref;

  Future<void> _load() async {
    final prefs = await SharedPreferences.getInstance();
    final hour = prefs.getInt(AppConfig.notificationHourKey);
    final minute = prefs.getInt(AppConfig.notificationMinuteKey);
    if (hour != null && minute != null) {
      state = TimeOfDay(hour: hour, minute: minute);
    }
    await _ref.read(notificationServiceProvider).scheduleDailyReminder(state);
  }

  Future<void> setTime(TimeOfDay time) async {
    state = time;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setInt(AppConfig.notificationHourKey, time.hour);
    await prefs.setInt(AppConfig.notificationMinuteKey, time.minute);
    await _ref.read(notificationServiceProvider).scheduleDailyReminder(time);
  }
}

final notificationTimeProvider =
    StateNotifierProvider<NotificationTimeNotifier, TimeOfDay>(
  (ref) => NotificationTimeNotifier(ref),
);
