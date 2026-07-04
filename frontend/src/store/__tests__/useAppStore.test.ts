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
