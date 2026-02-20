import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/connectivity_service.dart';
import '../data/checkin_repository.dart';
import '../data/offline_queue_service.dart';

enum SyncState { idle, syncing }

final offlineQueueServiceProvider = Provider<OfflineQueueService>(
  (ref) => OfflineQueueService(),
);

final pendingCountProvider = StateProvider<int>((ref) {
  return ref.read(offlineQueueServiceProvider).pendingCount;
});

class SyncNotifier extends StateNotifier<SyncState> {
  final Ref _ref;
  StreamSubscription<bool>? _subscription;

  SyncNotifier(this._ref) : super(SyncState.idle) {
    _startListening();
  }

  void _startListening() {
    final connectivityService = _ref.read(connectivityServiceProvider);
    _subscription = connectivityService.onStatusChange.listen((isOnline) {
      if (isOnline && _ref.read(pendingCountProvider) > 0) {
        syncPending();
      }
    });
  }

  Future<void> syncPending() async {
    if (state == SyncState.syncing) return;
    state = SyncState.syncing;

    try {
      final queue = _ref.read(offlineQueueServiceProvider);
      final repository = _ref.read(checkinRepositoryProvider);
      final pending = queue.getPending();

      for (final entry in pending.entries) {
        try {
          await repository.submitCheckin(entry.value);
          await queue.remove(entry.key);
          _ref.read(pendingCountProvider.notifier).update((n) => n > 0 ? n - 1 : 0);
        } on DioException catch (e) {
          if (_isNetworkError(e)) {
            break; // network gone — stop and keep remaining items
          }
          // 5xx or other server error — skip this item for now
        } catch (_) {
          // Exception from 400 → already submitted today — remove silently
          await queue.remove(entry.key);
          _ref.read(pendingCountProvider.notifier).update((n) => n > 0 ? n - 1 : 0);
        }
      }
    } finally {
      state = SyncState.idle;
    }
  }

  bool _isNetworkError(DioException e) {
    return e.type == DioExceptionType.connectionError ||
        e.type == DioExceptionType.connectionTimeout ||
        e.type == DioExceptionType.receiveTimeout ||
        e.type == DioExceptionType.sendTimeout ||
        (e.type == DioExceptionType.unknown && e.response == null);
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }
}

final syncNotifierProvider =
    StateNotifierProvider<SyncNotifier, SyncState>(
  (ref) => SyncNotifier(ref),
);
