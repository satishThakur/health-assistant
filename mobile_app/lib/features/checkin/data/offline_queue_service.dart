import 'dart:convert';

import 'package:hive/hive.dart';

import '../../../core/config/app_config.dart';
import '../domain/checkin_model.dart';

class OfflineQueueService {
  Box<String> get _box => Hive.box<String>(AppConfig.pendingCheckinsBox);

  Future<void> enqueue(CheckinModel checkin) async {
    final key = DateTime.now().toIso8601String();
    await _box.put(key, jsonEncode(checkin.toJson()));
  }

  Map<String, CheckinModel> getPending() {
    final result = <String, CheckinModel>{};
    for (final dynamic key in _box.keys) {
      final value = _box.get(key as String);
      if (value != null) {
        try {
          result[key] = CheckinModel.fromJson(
            jsonDecode(value) as Map<String, dynamic>,
          );
        } catch (_) {
          // Corrupted entry — skip
        }
      }
    }
    return result;
  }

  Future<void> remove(String key) async {
    await _box.delete(key);
  }

  int get pendingCount => _box.length;
}
