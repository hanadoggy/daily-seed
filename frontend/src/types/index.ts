export type ContextMode = 'Growth' | 'Rest' | 'Work';

export interface TaskMetrics {
  dailyTarget: number;
  totalTarget: number;
}

export interface TaskProgress {
  taskId: string;
  title: string;
  totalTarget: number;
  totalCompleted: number;
  percentage: number;
}

export interface MigrationResult {
  archivedTask: Task;
  newTask: Task;
}

export interface TaskConditions {
  weather: string;
  mode: string;
}

export interface Task {
  id: string;
  section: 'japanese' | 'dev' | 'self_dev';
  title: string;
  type: 'quantitative' | 'boolean';
  metrics: TaskMetrics;
  conditions: TaskConditions;
  status: 'active' | 'archived';
}

export interface Habit {
  id: string;
  title: string;
  category: string;
  status: 'active' | 'archived';
}

export interface TaskEntry {
  taskId: string;
  targetAmount: number;
  actualAmount: number;
  isCompleted: boolean;
}

export interface HabitEntry {
  habitId: string;
  isCompleted: boolean;
}

export interface DayContext {
  mode: ContextMode;
  weather: string;
}

export interface Journal {
  oneLineReview: string;
  threeLineDiary: string;
}

export interface DailyRecord {
  id: string;
  date: string;
  context: DayContext;
  tasks: TaskEntry[];
  habits: HabitEntry[];
  journal: Journal;
}
