import { create } from 'zustand';
import type { ContextMode, DailyRecord, Journal, Task, Habit } from '../types';
import {
  fetchDailyRecord,
  patchDailyRecord,
  fetchTasks,
  fetchHabits,
  createTask as apiCreateTask,
  updateTask as apiUpdateTask,
  deleteTask as apiDeleteTask,
  createHabit as apiCreateHabit,
  updateHabit as apiUpdateHabit,
  deleteHabit as apiDeleteHabit,
} from '../api/client';

interface AppState {
  selectedDate: string;
  currentMode: ContextMode;
  dailyRecord: DailyRecord | null;
  isLoading: boolean;
  error: string | null;

  tasks: Task[];
  habits: Habit[];

  setDateAndFetch: (date: string) => Promise<void>;
  fetchMasterData: () => Promise<void>;
  updateContextMode: (mode: ContextMode) => Promise<void>;
  saveJournal: (journalData: Journal) => Promise<void>;
  toggleHabitOptimistic: (habitId: string, isCompleted: boolean) => void;
  updateTaskProgressOptimistic: (taskId: string, amount: number) => void;

  addTask: (task: Omit<Task, 'id' | 'status'>) => Promise<void>;
  editTask: (id: string, task: Omit<Task, 'id' | 'status'>) => Promise<void>;
  archiveTask: (id: string) => Promise<void>;
  addHabit: (habit: Omit<Habit, 'id' | 'status'>) => Promise<void>;
  editHabit: (id: string, habit: Omit<Habit, 'id' | 'status'>) => Promise<void>;
  archiveHabit: (id: string) => Promise<void>;
}

export const useAppStore = create<AppState>((set, get) => ({
  selectedDate: '',
  currentMode: 'Growth',
  dailyRecord: null,
  isLoading: false,
  error: null,
  tasks: [],
  habits: [],

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

  fetchMasterData: async () => {
    try {
      const [tasks, habits] = await Promise.all([fetchTasks(), fetchHabits()]);
      set({ tasks, habits });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch master data';
      set({ error: message });
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

  addTask: async (task) => {
    try {
      const created = await apiCreateTask(task);
      set((state) => ({ tasks: [...state.tasks, created] }));
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create task';
      set({ error: message });
    }
  },

  editTask: async (id, task) => {
    const previous = get().tasks;
    try {
      const updated = await apiUpdateTask(id, task);
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? updated : t)),
      }));
    } catch (err) {
      set({ tasks: previous });
      const message = err instanceof Error ? err.message : 'Failed to update task';
      set({ error: message });
    }
  },

  archiveTask: async (id) => {
    const previous = get().tasks;
    // Optimistic removal from the active list
    set((state) => ({ tasks: state.tasks.filter((t) => t.id !== id) }));
    try {
      await apiDeleteTask(id);
    } catch {
      set({ tasks: previous });
    }
  },

  addHabit: async (habit) => {
    try {
      const created = await apiCreateHabit(habit);
      set((state) => ({ habits: [...state.habits, created] }));
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
    } catch (err) {
      set({ habits: previous });
      const message = err instanceof Error ? err.message : 'Failed to update habit';
      set({ error: message });
    }
  },

  archiveHabit: async (id) => {
    const previous = get().habits;
    // Optimistic removal from the active list
    set((state) => ({ habits: state.habits.filter((h) => h.id !== id) }));
    try {
      await apiDeleteHabit(id);
    } catch {
      set({ habits: previous });
    }
  },
}));
