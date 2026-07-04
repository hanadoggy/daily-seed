import type { ContextMode, DailyRecord, Journal, Task, Habit, TaskProgress } from '../types';

export interface DailySlice {
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

export interface TaskSlice {
  tasks: Task[];
  taskProgress: TaskProgress[];
  addTask: (task: Omit<Task, 'id' | 'status'>) => Promise<void>;
  editTask: (id: string, task: Omit<Task, 'id' | 'status'>) => Promise<void>;
  archiveTask: (id: string) => Promise<void>;
  migrateTask: (id: string) => Promise<void>;
  fetchProgress: () => Promise<void>;
}

export interface HabitSlice {
  habits: Habit[];
  addHabit: (habit: Omit<Habit, 'id' | 'status'>) => Promise<void>;
  editHabit: (id: string, habit: Omit<Habit, 'id' | 'status'>) => Promise<void>;
  archiveHabit: (id: string) => Promise<void>;
}

export interface AppState extends DailySlice, TaskSlice, HabitSlice {
  fetchMasterData: () => Promise<void>;
}
