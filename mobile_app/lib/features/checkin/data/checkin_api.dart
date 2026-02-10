import 'package:dio/dio.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/api_endpoints.dart';
import '../domain/checkin_model.dart';

class CheckinApi {
  final ApiClient _client;

  CheckinApi(this._client);

  Future<CheckinResponse> submitCheckin(CheckinModel checkin) async {
    try {
      final response = await _client.post(
        ApiEndpoints.checkin,
        data: checkin.toJson(),
      );
      return CheckinResponse.fromJson(response.data);
    } on DioException catch (e) {
      if (e.response?.statusCode == 400) {
        throw Exception(
          e.response?.data['message'] ?? 'Validation failed',
        );
      }
      rethrow;
    }
  }

  Future<CheckinModel?> getLatestCheckin() async {
    try {
      final response = await _client.get(ApiEndpoints.checkinLatest);

      if (response.data['checkin'] == null) {
        return null;
      }

      return CheckinModel.fromJson(response.data['checkin']);
    } catch (e) {
      throw Exception('Failed to fetch latest check-in: $e');
    }
  }

  Future<List<CheckinHistoryItem>> getCheckinHistory({int days = 30}) async {
    try {
      final response = await _client.get(
        ApiEndpoints.checkinHistory,
        queryParameters: {'days': days},
      );

      final history = (response.data['history'] as List)
          .map((item) => CheckinHistoryItem.fromJson(item))
          .toList();

      return history;
    } catch (e) {
      throw Exception('Failed to fetch check-in history: $e');
    }
  }
}
