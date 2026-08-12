import type { StateCreator } from 'zustand';
import { toast } from 'sonner';
import type { AppState, DailySlice } from '../types';
import { fetchDailyRecord, patchDailyRecord, fetchExistingRecordDates, ApiError } from '../../api/client';
import { getErrorMessage } from '../../lib/errorMessages';
import type { ContextMode, Journal, HabitEntry, TaskEntry } from '../../types';

export const createDailySlice: StateCreator<AppState, [], [], DailySlice> = (set, get) => ({
  selectedDate: '',
  currentMode: 'Growth',
  currentWeather: 'sunny',
  isAdminMode: false,
  existingRecordDates: [],
  dailyRecord: null,
  isLoading: false,
  error: null,

  toggleAdminMode: () => set((state) => ({ isAdminMode: !state.isAdminMode })),

  fetchExistingRecordDates: async (year: number, month: number) => {
    try {
      const response = await fetchExistingRecordDates(year, month);
      set({ existingRecordDates: response.dates });
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchExistingDates', status);
      toast.error(message);
      set({ error: message });
    }
  },

  setDateAndFetch: async (date: string) => {
    set({ selectedDate: date, isAdminMode: false, isLoading: true, error: null });
    try {
      const record = await fetchDailyRecord(date);
      set({
        dailyRecord: record,
        currentMode: record.context.mode,
        currentWeather: record.context.weather || 'sunny',
        isLoading: false,
      });
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchDailyRecord', status);
      toast.error(message);
      set({ isLoading: false, error: message });
    }
  },

  updateContextMode: async (mode: ContextMode) => {
    const { dailyRecord, selectedDate } = get();
    if (!dailyRecord) return;

    const previousMode = get().currentMode;
    const previousRecord = dailyRecord;

    set({
      currentMode: mode,
      dailyRecord: {
        ...dailyRecord,
        context: { ...dailyRecord.context, mode },
      },
    });

    try {
      await patchDailyRecord(selectedDate, {
        context: { ...dailyRecord.context, mode },
      });
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('patchDailyRecord', status);
      toast.error(message);
      set({ currentMode: previousMode, dailyRecord: previousRecord, error: message });
    }
  },

  updateWeather: async (weather: string) => {
    const { dailyRecord, selectedDate } = get();
    if (!dailyRecord) return;

    const previousWeather = get().currentWeather;
    const previousRecord = dailyRecord;

    set({
      currentWeather: weather,
      dailyRecord: {
        ...dailyRecord,
        context: { ...dailyRecord.context, weather },
      },
    });

    try {
      await patchDailyRecord(selectedDate, {
        context: { ...dailyRecord.context, weather },
      });
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('patchDailyRecord', status);
      toast.error(message);
      set({ currentWeather: previousWeather, dailyRecord: previousRecord, error: message });
    }
  },

  saveJournal: async (journalData: Journal) => {
    const { dailyRecord, selectedDate } = get();
    if (!dailyRecord) return;

    try {
      const updated = await patchDailyRecord(selectedDate, {
        journal: journalData,
      });
      set({ dailyRecord: updated });
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('patchDailyRecord', status);
      toast.error(message);
      set({ error: message });
    }
  },

  toggleHabitOptimistic: (habitId: string, isCompleted: boolean) => {
    const { dailyRecord, selectedDate } = get();
    if (!dailyRecord) return;

    const previousRecord = dailyRecord;
    const updatedHabits = dailyRecord.habits.map((h: HabitEntry) =>
      h.habitId === habitId ? { ...h, isCompleted } : h,
    );

    set({ dailyRecord: { ...dailyRecord, habits: updatedHabits } });

    patchDailyRecord(selectedDate, { habits: updatedHabits }).catch((err) => {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('patchDailyRecord', status);
      toast.error(message);
      set({ dailyRecord: previousRecord, error: message });
    });
  },

  updateTaskProgressOptimistic: (taskId: string, amount: number) => {
    const { dailyRecord, selectedDate } = get();
    if (!dailyRecord) return;

    const previousRecord = dailyRecord;
    const updatedTasks = dailyRecord.tasks.map((t: TaskEntry) => {
      if (t.taskId !== taskId) return t;
      const newAmount = Math.max(0, amount);
      return {
        ...t,
        actualAmount: newAmount,
        isCompleted: newAmount >= t.targetAmount,
      };
    });

    set({ dailyRecord: { ...dailyRecord, tasks: updatedTasks } });

    patchDailyRecord(selectedDate, { tasks: updatedTasks })
      .then(() => {
        get().fetchProgress();
      })
      .catch((err) => {
        const status = err instanceof ApiError ? err.status : 0;
        const message = getErrorMessage('patchDailyRecord', status);
        toast.error(message);
        set({ dailyRecord: previousRecord, error: message });
      });
  },
});
