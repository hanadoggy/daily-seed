import { create } from 'zustand';
import { toast } from 'sonner';
import type { AppState } from './types';
import { createDailySlice } from './slices/createDailySlice';
import { createTaskSlice } from './slices/createTaskSlice';
import { createHabitSlice } from './slices/createHabitSlice';
import { fetchTasks, fetchHabits, ApiError } from '../api/client';
import { getErrorMessage } from '../lib/errorMessages';

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
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchMasterData', status);
      toast.error(message);
      set({ error: message });
    }
  },
}));
