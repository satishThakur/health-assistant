import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/api_client.dart';
import '../domain/checkin_model.dart';
import 'checkin_api.dart';

final checkinRepositoryProvider = Provider<CheckinRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return CheckinRepository(CheckinApi(apiClient));
});

class CheckinRepository {
  final CheckinApi _api;

  CheckinRepository(this._api);

  Future<CheckinResponse> submitCheckin(CheckinModel checkin) async {
    return await _api.submitCheckin(checkin);
  }

  Future<CheckinModel?> getLatestCheckin() async {
    return await _api.getLatestCheckin();
  }

  Future<List<CheckinHistoryItem>> getCheckinHistory({int days = 30}) async {
    return await _api.getCheckinHistory(days: days);
  }
}
