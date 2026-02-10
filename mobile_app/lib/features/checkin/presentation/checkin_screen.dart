import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/config/theme.dart';
import '../providers/checkin_provider.dart';
import 'widgets/feeling_slider.dart';

class CheckinScreen extends ConsumerWidget {
  const CheckinScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formState = ref.watch(checkinFormProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Daily Check-in'),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => context.pop(),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Header
              Text(
                'How are you feeling today?',
                style: Theme.of(context).textTheme.displaySmall,
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 8),
              Text(
                'Rate your current state on a scale of 1-10',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Colors.grey[600],
                    ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 32),

              // Energy Slider
              FeelingSlider(
                label: 'Energy',
                emoji: '💪',
                value: formState.energy,
                color: AppTheme.energyColor,
                onChanged: (value) {
                  ref
                      .read(checkinFormProvider.notifier)
                      .updateEnergy(value.round());
                },
              ),
              const SizedBox(height: 24),

              // Mood Slider
              FeelingSlider(
                label: 'Mood',
                emoji: '😊',
                value: formState.mood,
                color: AppTheme.moodColor,
                onChanged: (value) {
                  ref
                      .read(checkinFormProvider.notifier)
                      .updateMood(value.round());
                },
              ),
              const SizedBox(height: 24),

              // Focus Slider
              FeelingSlider(
                label: 'Focus',
                emoji: '🎯',
                value: formState.focus,
                color: AppTheme.focusColor,
                onChanged: (value) {
                  ref
                      .read(checkinFormProvider.notifier)
                      .updateFocus(value.round());
                },
              ),
              const SizedBox(height: 24),

              // Physical Slider
              FeelingSlider(
                label: 'Physical',
                emoji: '🏃',
                value: formState.physical,
                color: AppTheme.physicalColor,
                onChanged: (value) {
                  ref
                      .read(checkinFormProvider.notifier)
                      .updatePhysical(value.round());
                },
              ),
              const SizedBox(height: 32),

              // Notes TextField
              TextField(
                decoration: InputDecoration(
                  labelText: 'Notes (optional)',
                  hintText: 'How are you feeling? Any thoughts?',
                  prefixIcon: const Icon(Icons.notes),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                maxLines: 3,
                maxLength: 1000,
                onChanged: (value) {
                  ref.read(checkinFormProvider.notifier).updateNotes(value);
                },
              ),
              const SizedBox(height: 24),

              // Error Message
              if (formState.error != null)
                Padding(
                  padding: const EdgeInsets.only(bottom: 16),
                  child: Text(
                    formState.error!,
                    style: TextStyle(
                      color: Theme.of(context).colorScheme.error,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ),

              // Submit Button
              ElevatedButton(
                onPressed: formState.isSubmitting
                    ? null
                    : () async {
                        final success = await ref
                            .read(checkinFormProvider.notifier)
                            .submitCheckin();

                        if (success && context.mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Check-in saved! 🎉'),
                              backgroundColor: Colors.green,
                            ),
                          );
                          context.pop();
                        }
                      },
                child: formState.isSubmitting
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                    : const Text('Submit Check-in'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
