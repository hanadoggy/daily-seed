import type { StateCreator } from 'zustand';
import type { AppState, DailySlice } from '../types';
import { fetchDailyRecord, patchDailyRecord, fetchExistingRecordDates } from '../../api/client';
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
      console.error('Failed to fetch existing record dates', err);
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
      const message = err instanceof Error ? err.message : 'Failed to fetch daily record';
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
    } catch {
      set({ currentMode: previousMode, dailyRecord: previousRecord });
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
    } catch {
      set({ currentWeather: previousWeather, dailyRecord: previousRecord });
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
      const message = err instanceof Error ? err.message : 'Failed to save journal';
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

    patchDailyRecord(selectedDate, { habits: updatedHabits }).catch(() => {
      set({ dailyRecord: previousRecord });
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
      .catch(() => {
        set({ dailyRecord: previousRecord });
      });
  },
});
