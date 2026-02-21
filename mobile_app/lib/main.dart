import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_timezone/flutter_timezone.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:timezone/data/latest.dart' as tz;
import 'package:timezone/timezone.dart' as tz;

import 'app.dart';
import 'core/config/app_config.dart';
import 'core/notifications/notification_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize Hive for local storage
  await Hive.initFlutter();
  await Hive.openBox<String>(AppConfig.pendingCheckinsBox);

  // Timezone (required for zonedSchedule)
  tz.initializeTimeZones();
  tz.setLocalLocation(
    tz.getLocation(await FlutterTimezone.getLocalTimezone()),
  );

  // Notification service — initialize before ProviderScope
  final notificationService = NotificationService();
  await notificationService.initialize();

  runApp(
    ProviderScope(
      overrides: [
        notificationServiceProvider.overrideWithValue(notificationService),
      ],
      child: const HealthAssistantApp(),
    ),
  );
}
