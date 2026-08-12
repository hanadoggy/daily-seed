import type { StateCreator } from 'zustand';
import { toast } from 'sonner';
import type { AppState, HabitSlice } from '../types';
import {
  createHabit as apiCreateHabit,
  updateHabit as apiUpdateHabit,
  deleteHabit as apiDeleteHabit,
  ApiError,
} from '../../api/client';
import { getErrorMessage } from '../../lib/errorMessages';

export const createHabitSlice: StateCreator<AppState, [], [], HabitSlice> = (set, get) => ({
  habits: [],

  addHabit: async (habit) => {
    try {
      const created = await apiCreateHabit(habit);
      set((state) => ({ habits: [...state.habits, created] }));
      get().setDateAndFetch(get().selectedDate);
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('createHabit', status);
      toast.error(message);
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
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('updateHabit', status);
      toast.error(message);
      set({ habits: previous, error: message });
    }
  },

  archiveHabit: async (id) => {
    const previous = get().habits;
    set((state) => ({
      habits: state.habits.map((h) => (h.id === id ? { ...h, status: 'archived' } : h)),
    }));
    try {
      await apiDeleteHabit(id);
      get().setDateAndFetch(get().selectedDate);
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('archiveHabit', status);
      toast.error(message);
      set({ habits: previous, error: message });
    }
  },
});
