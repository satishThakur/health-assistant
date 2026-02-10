import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class ApiInterceptor extends Interceptor {
  final Ref ref;

  ApiInterceptor(this.ref);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    // TODO: Add authentication token when implemented
    // final token = ref.read(authTokenProvider);
    // if (token != null) {
    //   options.headers['Authorization'] = 'Bearer $token';
    // }

    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    // Handle common errors
    if (err.response?.statusCode == 401) {
      // TODO: Handle unauthorized (logout user)
      // ref.read(authProvider.notifier).logout();
    }

    handler.next(err);
  }
}
