import { create } from 'zustand';
import type { ContextMode, DailyRecord, Journal } from '../types';
import { fetchDailyRecord, patchDailyRecord } from '../api/client';

interface AppState {
  selectedDate: string;
  currentMode: ContextMode;
  dailyRecord: DailyRecord | null;
  isLoading: boolean;
  error: string | null;

  setDateAndFetch: (date: string) => Promise<void>;
  updateContextMode: (mode: ContextMode) => Promise<void>;
  saveJournal: (journalData: Journal) => Promise<void>;
  toggleHabitOptimistic: (habitId: string, isCompleted: boolean) => void;
  updateTaskProgressOptimistic: (taskId: string, amount: number) => void;
}

export const useAppStore = create<AppState>((set, get) => ({
  selectedDate: '',
  currentMode: 'Growth',
  dailyRecord: null,
  isLoading: false,
  error: null,

  setDateAndFetch: async (date: string) => {
    set({ selectedDate: date, isLoading: true, error: null });
    try {
      const record = await fetchDailyRecord(date);
      set({
        dailyRecord: record,
        currentMode: record.context.mode as ContextMode,
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

    // Optimistic update
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
      // Rollback on failure
      set({ currentMode: previousMode, dailyRecord: previousRecord });
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

    const updatedHabits = dailyRecord.habits.map((h) =>
      h.habitId === habitId ? { ...h, isCompleted } : h,
    );

    // Optimistic update
    set({ dailyRecord: { ...dailyRecord, habits: updatedHabits } });

    patchDailyRecord(selectedDate, { habits: updatedHabits }).catch(() => {
      // Rollback on failure
      set({ dailyRecord: previousRecord });
    });
  },

  updateTaskProgressOptimistic: (taskId: string, amount: number) => {
    const { dailyRecord, selectedDate } = get();
    if (!dailyRecord) return;

    const previousRecord = dailyRecord;

    const updatedTasks = dailyRecord.tasks.map((t) => {
      if (t.taskId !== taskId) return t;
      const newAmount = Math.max(0, Math.min(amount, t.targetAmount));
      return {
        ...t,
        actualAmount: newAmount,
        isCompleted: newAmount >= t.targetAmount,
      };
    });

    // Optimistic update
    set({ dailyRecord: { ...dailyRecord, tasks: updatedTasks } });

    patchDailyRecord(selectedDate, { tasks: updatedTasks }).catch(() => {
      // Rollback on failure
      set({ dailyRecord: previousRecord });
    });
  },
}));
