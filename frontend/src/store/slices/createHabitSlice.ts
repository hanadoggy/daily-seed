import type { StateCreator } from 'zustand';
import type { AppState, HabitSlice } from '../types';
import {
  createHabit as apiCreateHabit,
  updateHabit as apiUpdateHabit,
  deleteHabit as apiDeleteHabit,
} from '../../api/client';

export const createHabitSlice: StateCreator<AppState, [], [], HabitSlice> = (set, get) => ({
  habits: [],

  addHabit: async (habit) => {
    try {
      const created = await apiCreateHabit(habit);
      set((state) => ({ habits: [...state.habits, created] }));
      get().setDateAndFetch(get().selectedDate);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create habit';
      set({ error: message });
    }
  },

  editHabit: async (id, habit) => {
    const previous = get().habits;
    try {
      const updated = await apiUpdateHabit(id, habit);
      set((state) => ({
        habits: state.habits.map((h) => (h.id === id ? updated : h)),
      }));
      get().setDateAndFetch(get().selectedDate);
    } catch (err) {
      set({ habits: previous });
      const message = err instanceof Error ? err.message : 'Failed to update habit';
      set({ error: message });
    }
  },

  archiveHabit: async (id) => {
    const previous = get().habits;
    set((state) => ({ habits: state.habits.filter((h) => h.id !== id) }));
    try {
      await apiDeleteHabit(id);
      get().setDateAndFetch(get().selectedDate);
    } catch {
      set({ habits: previous });
    }
  },
});
