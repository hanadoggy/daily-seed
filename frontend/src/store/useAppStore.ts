import { create } from 'zustand';
import type { AppState } from './types';
import { createDailySlice } from './slices/createDailySlice';
import { createTaskSlice } from './slices/createTaskSlice';
import { createHabitSlice } from './slices/createHabitSlice';
import { fetchTasks, fetchHabits } from '../api/client';

export const useAppStore = create<AppState>()((...a) => ({
  ...createDailySlice(...a),
  ...createTaskSlice(...a),
  ...createHabitSlice(...a),

  fetchMasterData: async () => {
    const set = a[0];
    try {
      const [tasks, habits] = await Promise.all([fetchTasks(), fetchHabits()]);
      set({ tasks, habits });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch master data';
      set({ error: message });
    }
  },
}));
