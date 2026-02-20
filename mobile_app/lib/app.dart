import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/config/theme.dart';
import 'core/routing/app_router.dart';
import 'features/auth/domain/auth_state.dart';
import 'features/auth/providers/auth_provider.dart';
import 'features/checkin/providers/sync_provider.dart';

class HealthAssistantApp extends ConsumerWidget {
  const HealthAssistantApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Keep connectivity listener alive for the app lifetime
    ref.watch(syncNotifierProvider);

    final authState = ref.watch(authProvider);

    // Show a plain loading screen while reading stored credentials to
    // avoid flashing the login screen on authenticated cold-starts.
    if (authState is AuthLoading) {
      return const MaterialApp(
        debugShowCheckedModeBanner: false,
        home: Scaffold(
          body: Center(child: CircularProgressIndicator()),
        ),
      );
    }

    final router = ref.watch(appRouterProvider);

    return MaterialApp.router(
      title: 'Health Assistant',
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
