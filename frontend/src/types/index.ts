export type ContextMode = 'Growth' | 'Rest' | 'Office' | 'Remote';

export interface TaskMetrics {
  dailyTarget: number;
  totalTarget: number;
}

export interface TaskProgress {
  taskId: string;
  title: string;
  type: 'quantitative' | 'boolean';
  totalTarget: number;
  totalCompleted: number;
  percentage: number;
}

export interface MigrationResult {
  archivedTask: Task;
  newTask: Task;
}

export interface TaskConditions {
  weather: string[];
  mode: string[];
}

export interface Task {
  id: string;
  section: 'japanese' | 'dev' | 'self_dev' | 'exercise';
  title: string;
  type: 'quantitative' | 'boolean';
  metrics: TaskMetrics;
  conditions: TaskConditions;
  status: 'active' | 'archived';
  startDate: string;
  endDate?: string;
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

export interface HeatmapDay {
  date: string;
  total: number;
  habits: number;
  sectionCounts: Record<string, number>;
}

export interface HeatmapResponse {
  days: HeatmapDay[];
}

export interface TaskStat {
  taskId: string;
  title: string;
  section: 'japanese' | 'dev' | 'self_dev' | 'exercise' | string;
  type: 'quantitative' | 'boolean';
  rate: number;
  completed: number;
  target: number;
}

export interface TaskCompletionStats {
  overall: number;
  sections: Record<string, number>;
  perTask: TaskStat[];
}

export interface HabitStat {
  habitId: string;
  title: string;
  category: string;
  rate: number;
  completed: number;
  total: number;
}

export interface HabitCompletionStats {
  overall: number;
  categories: Record<string, number>;
  perHabit: HabitStat[];
}

export interface JournalEntry {
  date: string;
  oneLineReview: string;
  threeLineDiary: string;
}

export interface SummaryResponse {
  period: 'weekly' | 'monthly';
  startDate: string;
  endDate: string;
  totalDays: number;
  recordedDays: number;
  taskCompletion: TaskCompletionStats;
  habitCompletion: HabitCompletionStats;
  modeDistribution: Record<string, number>;
  journals: JournalEntry[];
}
