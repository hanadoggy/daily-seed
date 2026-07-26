import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAnalyticsStore } from '../useAnalyticsStore';
import * as apiClient from '@/api/client';
import type { SummaryResponse } from '@/types';

vi.mock('@/api/client', () => ({
  fetchHeatmap: vi.fn(),
  fetchSummary: vi.fn(),
  fetchStreaks: vi.fn(),
}));

describe('useAnalyticsStore', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAnalyticsStore.setState({
      heatmapData: null,
      isLoading: false,
      error: null,
      summaryData: null,
      summaryPeriod: 'weekly',
      summaryDate: '2026-07-25',
      isSummaryLoading: false,
      summaryError: null,
      streakData: null,
      isStreakLoading: false,
      streakError: null,
    });
  });

  it('initializes with default values', () => {
    const state = useAnalyticsStore.getState();
    expect(state.heatmapData).toBeNull();
    expect(state.isLoading).toBe(false);
    expect(state.error).toBeNull();
    expect(state.summaryData).toBeNull();
    expect(state.summaryPeriod).toBe('weekly');
    expect(state.streakData).toBeNull();
    expect(state.isStreakLoading).toBe(false);
  });

  it('handles fetchStreakData success', async () => {
    const mockStreaks = {
      habits: [
        {
          habitId: 'h1',
          title: 'Meditation',
          category: 'Mindfulness',
          currentStreak: 7,
          longestStreak: 7,
          totalDays: 7,
          lastCompleted: '2026-07-25',
          milestones: [7],
        },
      ],
    };
    vi.mocked(apiClient.fetchStreaks).mockResolvedValueOnce(mockStreaks);

    const promise = useAnalyticsStore.getState().fetchStreakData();
    expect(useAnalyticsStore.getState().isStreakLoading).toBe(true);

    await promise;

    const state = useAnalyticsStore.getState();
    expect(state.isStreakLoading).toBe(false);
    expect(state.streakData).toEqual(mockStreaks);
    expect(state.streakError).toBeNull();
    expect(apiClient.fetchStreaks).toHaveBeenCalled();
  });

  it('handles fetchHeatmapData success', async () => {
    const mockData = {
      days: [
        { date: '2026-01-01', total: 2, habits: 1, sectionCounts: { dev: 1 } },
      ],
    };
    vi.mocked(apiClient.fetchHeatmap).mockResolvedValueOnce(mockData);

    const promise = useAnalyticsStore.getState().fetchHeatmapData(2026);
    expect(useAnalyticsStore.getState().isLoading).toBe(true);

    await promise;

    const state = useAnalyticsStore.getState();
    expect(state.isLoading).toBe(false);
    expect(state.heatmapData).toEqual(mockData);
    expect(state.error).toBeNull();
    expect(apiClient.fetchHeatmap).toHaveBeenCalledWith(2026);
  });

  it('handles fetchHeatmapData failure with custom API error response', async () => {
    const mockError = {
      response: {
        data: {
          message: 'Server internal error occurred',
        },
      },
    };
    vi.mocked(apiClient.fetchHeatmap).mockRejectedValueOnce(mockError);

    await useAnalyticsStore.getState().fetchHeatmapData(2026);

    const state = useAnalyticsStore.getState();
    expect(state.isLoading).toBe(false);
    expect(state.heatmapData).toBeNull();
    expect(state.error).toBe('Server internal error occurred');
  });

  it('handles fetchSummaryData success and navigation', async () => {
    const mockSummary: SummaryResponse = {
      period: 'weekly',
      startDate: '2026-07-19',
      endDate: '2026-07-25',
      totalDays: 7,
      recordedDays: 2,
      taskCompletion: { overall: 100, sections: {}, perTask: [] },
      habitCompletion: { overall: 100, categories: {}, perHabit: [] },
      modeDistribution: { Growth: 2 },
      journals: [],
    };
    vi.mocked(apiClient.fetchSummary).mockResolvedValue(mockSummary);

    await useAnalyticsStore.getState().fetchSummaryData('weekly', '2026-07-25');

    let state = useAnalyticsStore.getState();
    expect(state.isSummaryLoading).toBe(false);
    expect(state.summaryData).toEqual(mockSummary);
    expect(apiClient.fetchSummary).toHaveBeenCalledWith('weekly', '2026-07-25');

    // Test period toggle
    useAnalyticsStore.getState().setSummaryPeriod('monthly');
    expect(apiClient.fetchSummary).toHaveBeenCalledWith('monthly', '2026-07-25');

    // Test navigate summary
    useAnalyticsStore.getState().navigateSummary('prev');
    expect(apiClient.fetchSummary).toHaveBeenCalled();
  });
});
