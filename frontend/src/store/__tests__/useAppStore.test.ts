import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAppStore } from '../useAppStore';
import * as apiClient from '../../api/client';
import type { DailyRecord } from '../../types';

// Mock the API client
vi.mock('../../api/client', () => ({
  fetchDailyRecord: vi.fn(),
  patchDailyRecord: vi.fn(),
  fetchTasks: vi.fn(),
  fetchHabits: vi.fn(),
  createTask: vi.fn(),
  updateTask: vi.fn(),
  deleteTask: vi.fn(),
  createHabit: vi.fn(),
  updateHabit: vi.fn(),
  deleteHabit: vi.fn(),
  migrateTask: vi.fn(),
  fetchTaskProgress: vi.fn(),
}));

describe('useAppStore Optimistic Updates', () => {
  const initialRecord: DailyRecord = {
    id: '2023-10-10',
    date: '2023-10-10',
    context: { mode: 'Growth', weather: 'sunny' },
    tasks: [{ taskId: 'task1', targetAmount: 2, actualAmount: 0, isCompleted: false }],
    habits: [{ habitId: 'habit1', isCompleted: false }],
    journal: { oneLineReview: '', threeLineDiary: '' },
  };

  beforeEach(() => {
    // Reset store state
    useAppStore.setState({
      selectedDate: '2023-10-10',
      dailyRecord: initialRecord,
      currentMode: 'Growth',
    });
    vi.clearAllMocks();
  });

  it('updateTaskProgressOptimistic should update immediately and rollback on API failure', async () => {
    // Make patchDailyRecord reject
    vi.mocked(apiClient.patchDailyRecord).mockRejectedValueOnce(new Error('Network error'));

    // Trigger optimistic update
    useAppStore.getState().updateTaskProgressOptimistic('task1', 2);

    // Verify immediate state change
    expect(useAppStore.getState().dailyRecord?.tasks[0].actualAmount).toBe(2);
    expect(useAppStore.getState().dailyRecord?.tasks[0].isCompleted).toBe(true);

    // Wait for the promise rejection to be handled
    await new Promise((resolve) => setTimeout(resolve, 0));

    // Verify state was rolled back
    expect(useAppStore.getState().dailyRecord?.tasks[0].actualAmount).toBe(0);
    expect(useAppStore.getState().dailyRecord?.tasks[0].isCompleted).toBe(false);
  });

  it('updateTaskProgressOptimistic should NOT rollback on API success', async () => {
    vi.mocked(apiClient.patchDailyRecord).mockResolvedValueOnce({} as DailyRecord);

    useAppStore.getState().updateTaskProgressOptimistic('task1', 1);

    expect(useAppStore.getState().dailyRecord?.tasks[0].actualAmount).toBe(1);

    await new Promise((resolve) => setTimeout(resolve, 0));

    // State should remain updated
    expect(useAppStore.getState().dailyRecord?.tasks[0].actualAmount).toBe(1);
  });

  it('toggleHabitOptimistic should rollback on API failure', async () => {
    vi.mocked(apiClient.patchDailyRecord).mockRejectedValueOnce(new Error('Failed'));

    useAppStore.getState().toggleHabitOptimistic('habit1', true);

    expect(useAppStore.getState().dailyRecord?.habits[0].isCompleted).toBe(true);

    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(useAppStore.getState().dailyRecord?.habits[0].isCompleted).toBe(false);
  });

  it('updateContextMode should rollback on API failure', async () => {
    vi.mocked(apiClient.patchDailyRecord).mockRejectedValueOnce(new Error('Failed'));

    await useAppStore.getState().updateContextMode('Office');

    // Due to await, the try-catch in updateContextMode already finished rolling back.
    expect(useAppStore.getState().currentMode).toBe('Growth');
    expect(useAppStore.getState().dailyRecord?.context.mode).toBe('Growth');
  });

  it('updateContextMode should persist on API success', async () => {
    vi.mocked(apiClient.patchDailyRecord).mockResolvedValueOnce({} as DailyRecord);

    await useAppStore.getState().updateContextMode('Office');

    expect(useAppStore.getState().currentMode).toBe('Office');
    expect(useAppStore.getState().dailyRecord?.context.mode).toBe('Office');
  });

  it('updateWeather should rollback on API failure', async () => {
    vi.mocked(apiClient.patchDailyRecord).mockRejectedValueOnce(new Error('Failed'));

    await useAppStore.getState().updateWeather('rainy');

    expect(useAppStore.getState().currentWeather).toBe('sunny');
    expect(useAppStore.getState().dailyRecord?.context.weather).toBe('sunny');
  });

  it('updateWeather should persist on API success', async () => {
    vi.mocked(apiClient.patchDailyRecord).mockResolvedValueOnce({} as DailyRecord);

    await useAppStore.getState().updateWeather('rainy');

    expect(useAppStore.getState().currentWeather).toBe('rainy');
    expect(useAppStore.getState().dailyRecord?.context.weather).toBe('rainy');
  });
});

describe('TaskSlice', () => {
  beforeEach(() => {
    useAppStore.setState({ tasks: [], error: null });
    vi.clearAllMocks();
  });

  it('addTask should append to list on success', async () => {
    const mockTask = { id: 't1', title: 'New Task' };
    vi.mocked(apiClient.createTask).mockResolvedValueOnce(mockTask as any);
    vi.mocked(apiClient.fetchTasks).mockResolvedValueOnce([mockTask] as any);

    await useAppStore.getState().addTask({ title: 'New Task' } as any);

    expect(useAppStore.getState().tasks).toContainEqual(mockTask);
  });

  it('archiveTask should remove optimistically and rollback on failure', async () => {
    const existingTask = { id: 't1', title: 'Task 1', status: 'active' };
    useAppStore.setState({ tasks: [existingTask as any] });
    vi.mocked(apiClient.deleteTask).mockRejectedValueOnce(new Error('fail'));

    await useAppStore.getState().archiveTask('t1');

    // Should be rolled back
    expect(useAppStore.getState().tasks).toHaveLength(1);
  });

  it('migrateTask should replace archived and add new', async () => {
    const existingTask = { id: 't1', title: 'Task 1', status: 'active' };
    useAppStore.setState({ tasks: [existingTask as any] });
    const mockResult = {
      archivedTask: { id: 't1', title: 'Task 1', status: 'archived' },
      newTask: { id: 't2', title: 'Task 1', status: 'active' },
    };
    vi.mocked(apiClient.migrateTask).mockResolvedValueOnce(mockResult as any);
    vi.mocked(apiClient.fetchTaskProgress).mockResolvedValueOnce([]);

    await useAppStore.getState().migrateTask('t1', '2023-10-10');

    expect(useAppStore.getState().tasks).toContainEqual(mockResult.archivedTask);
    expect(useAppStore.getState().tasks).toContainEqual(mockResult.newTask);
  });
});

describe('HabitSlice', () => {
  beforeEach(() => {
    useAppStore.setState({ habits: [], error: null });
    vi.clearAllMocks();
  });

  it('addHabit should fetch habits on success', async () => {
    const mockHabit = { id: 'h1', title: 'New Habit' };
    vi.mocked(apiClient.createHabit).mockResolvedValueOnce(mockHabit as any);
    vi.mocked(apiClient.fetchHabits).mockResolvedValueOnce([mockHabit] as any);

    await useAppStore.getState().addHabit({ title: 'New Habit' } as any);

    expect(useAppStore.getState().habits).toContainEqual(mockHabit);
  });

  it('archiveHabit should rollback on failure', async () => {
    const existingHabit = { id: 'h1', title: 'Habit 1', status: 'active' };
    useAppStore.setState({ habits: [existingHabit as any] });
    vi.mocked(apiClient.deleteHabit).mockRejectedValueOnce(new Error('fail'));

    await useAppStore.getState().archiveHabit('h1');

    expect(useAppStore.getState().habits).toHaveLength(1);
  });
});

describe('DailySlice extended', () => {
  beforeEach(() => {
    useAppStore.setState({ dailyRecord: null, error: null });
    vi.clearAllMocks();
  });

  it('setDateAndFetch populates record on success', async () => {
    const mockRecord = { id: '2023-10-10', context: { mode: 'Growth', weather: 'sunny' } };
    vi.mocked(apiClient.fetchDailyRecord).mockResolvedValueOnce(mockRecord as any);

    await useAppStore.getState().setDateAndFetch('2023-10-10');

    expect(useAppStore.getState().selectedDate).toBe('2023-10-10');
    expect(useAppStore.getState().dailyRecord).toEqual(mockRecord);
    expect(useAppStore.getState().currentMode).toBe('Growth');
  });

  it('saveJournal updates record on success', async () => {
    useAppStore.setState({ dailyRecord: { id: '2023-10-10', journal: {} } as any });
    vi.mocked(apiClient.patchDailyRecord).mockResolvedValueOnce({ 
      id: '2023-10-10', 
      journal: { oneLineReview: 'test', threeLineDiary: '' } 
    } as any);

    await useAppStore.getState().saveJournal({ oneLineReview: 'test', threeLineDiary: '' });

    expect(useAppStore.getState().dailyRecord?.journal.oneLineReview).toBe('test');
  });
});
