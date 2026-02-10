import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/checkin_repository.dart';
import '../domain/checkin_model.dart';

// State for check-in form
class CheckinFormState {
  final int energy;
  final int mood;
  final int focus;
  final int physical;
  final String notes;
  final bool isSubmitting;
  final String? error;

  CheckinFormState({
    this.energy = 5,
    this.mood = 5,
    this.focus = 5,
    this.physical = 5,
    this.notes = '',
    this.isSubmitting = false,
    this.error,
  });

  CheckinFormState copyWith({
    int? energy,
    int? mood,
    int? focus,
    int? physical,
    String? notes,
    bool? isSubmitting,
    String? error,
  }) {
    return CheckinFormState(
      energy: energy ?? this.energy,
      mood: mood ?? this.mood,
      focus: focus ?? this.focus,
      physical: physical ?? this.physical,
      notes: notes ?? this.notes,
      isSubmitting: isSubmitting ?? this.isSubmitting,
      error: error,
    );
  }

  CheckinModel toCheckinModel() {
    return CheckinModel(
      energy: energy,
      mood: mood,
      focus: focus,
      physical: physical,
      notes: notes.isEmpty ? null : notes,
    );
  }
}

// Check-in form controller
class CheckinFormNotifier extends StateNotifier<CheckinFormState> {
  final CheckinRepository _repository;

  CheckinFormNotifier(this._repository) : super(CheckinFormState());

  void updateEnergy(int value) {
    state = state.copyWith(energy: value, error: null);
  }

  void updateMood(int value) {
    state = state.copyWith(mood: value, error: null);
  }

  void updateFocus(int value) {
    state = state.copyWith(focus: value, error: null);
  }

  void updatePhysical(int value) {
    state = state.copyWith(physical: value, error: null);
  }

  void updateNotes(String value) {
    state = state.copyWith(notes: value, error: null);
  }

  Future<bool> submitCheckin() async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final checkin = state.toCheckinModel();
      await _repository.submitCheckin(checkin);

      // Reset form after successful submission
      state = CheckinFormState();
      return true;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        error: e.toString(),
      );
      return false;
    }
  }
}

final checkinFormProvider =
    StateNotifierProvider<CheckinFormNotifier, CheckinFormState>((ref) {
  final repository = ref.watch(checkinRepositoryProvider);
  return CheckinFormNotifier(repository);
});

// Latest check-in provider
final latestCheckinProvider = FutureProvider<CheckinModel?>((ref) async {
  final repository = ref.watch(checkinRepositoryProvider);
  return await repository.getLatestCheckin();
});

// Check-in history provider
final checkinHistoryProvider =
    FutureProvider.family<List<CheckinHistoryItem>, int>((ref, days) async {
  final repository = ref.watch(checkinRepositoryProvider);
  return await repository.getCheckinHistory(days: days);
});
