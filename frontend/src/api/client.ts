import type { DailyRecord, Task, Habit, TaskProgress, MigrationResult } from '../types';

const BASE_URL = import.meta.env.VITE_API_URL ?? '/api/v1';

export class ApiError extends Error {
  status: number;

  constructor(status: number) {
    super(`Request failed with status ${status}`);
    this.name = 'ApiError';
    this.status = status;
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (!res.ok) {
    throw new ApiError(res.status);
  }

  return res.json();
}

// --- Daily Record ---

const patchQueues = new Map<string, Promise<DailyRecord>>();

export function fetchDailyRecord(date: string): Promise<DailyRecord> {
  return request<DailyRecord>(`/daily/${date}`);
}

export function fetchExistingRecordDates(year: number, month: number): Promise<{dates: string[]}> {
  return request<{dates: string[]}>(`/daily/exists?year=${year}&month=${month}`);
}

export function patchDailyRecord(
  date: string,
  body: Partial<DailyRecord>,
): Promise<DailyRecord> {
  const prev = patchQueues.get(date) ?? Promise.resolve(null as unknown as DailyRecord);
  const next = prev
    .catch(() => {})
    .then(() =>
      request<DailyRecord>(`/daily/${date}`, {
        method: 'PATCH',
        body: JSON.stringify(body),
      }),
    );
  patchQueues.set(date, next);
  next.finally(() => {
    if (patchQueues.get(date) === next) {
      patchQueues.delete(date);
    }
  });
  return next;
}

// --- Tasks ---

export function fetchTasks(): Promise<Task[]> {
  return request<Task[]>('/tasks');
}

export function createTask(task: Omit<Task, 'id' | 'status'>): Promise<Task> {
  return request<Task>('/tasks', {
    method: 'POST',
    body: JSON.stringify(task),
  });
}

export function updateTask(id: string, task: Omit<Task, 'id' | 'status'>): Promise<Task> {
  return request<Task>(`/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(task),
  });
}

export function deleteTask(id: string): Promise<void> {
  return request(`/tasks/${id}`, { method: 'DELETE' });
}

export function fetchTaskProgress(): Promise<TaskProgress[]> {
  return request<TaskProgress[]>('/tasks/progress');
}

export function migrateTask(id: string, completionDate: string): Promise<MigrationResult> {
  return request<MigrationResult>(`/tasks/${id}/migrate`, {
    method: 'POST',
    body: JSON.stringify({ completionDate }),
  });
}

// --- Habits ---

export function fetchHabits(): Promise<Habit[]> {
  return request<Habit[]>('/habits');
}

export function createHabit(habit: Omit<Habit, 'id' | 'status'>): Promise<Habit> {
  return request<Habit>('/habits', {
    method: 'POST',
    body: JSON.stringify(habit),
  });
}

export function updateHabit(id: string, habit: Omit<Habit, 'id' | 'status'>): Promise<Habit> {
  return request<Habit>(`/habits/${id}`, {
    method: 'PUT',
    body: JSON.stringify(habit),
  });
}

export function deleteHabit(id: string): Promise<void> {
  return request(`/habits/${id}`, { method: 'DELETE' });
}

// --- Analytics ---

export function fetchHeatmap(year: number): Promise<import('../types').HeatmapResponse> {
  return request<import('../types').HeatmapResponse>(`/analytics/heatmap?year=${year}`);
}

export function fetchSummary(
  period: 'weekly' | 'monthly',
  date: string,
): Promise<import('../types').SummaryResponse> {
  return request<import('../types').SummaryResponse>(
    `/analytics/summary?period=${period}&date=${date}`,
  );
}

export function fetchStreaks(): Promise<import('../types').StreakResponse> {
  return request<import('../types').StreakResponse>('/analytics/streaks');
}
