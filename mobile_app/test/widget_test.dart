import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:health_assistant/app.dart';

void main() {
  testWidgets('App renders without errors', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(child: HealthAssistantApp()),
    );
    // App renders (loading state shown while auth initializes)
    expect(find.byType(CircularProgressIndicator), findsOneWidget);
  });
}
