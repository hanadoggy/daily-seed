import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAnalyticsStore } from '../useAnalyticsStore';
import * as apiClient from '@/api/client';

vi.mock('@/api/client', () => ({
  fetchHeatmap: vi.fn(),
}));

describe('useAnalyticsStore', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAnalyticsStore.setState({
      heatmapData: null,
      isLoading: false,
      error: null,
    });
  });

  it('initializes with default values', () => {
    const state = useAnalyticsStore.getState();
    expect(state.heatmapData).toBeNull();
    expect(state.isLoading).toBe(false);
    expect(state.error).toBeNull();
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

  it('handles fetchHeatmapData failure with generic fallback error message', async () => {
    vi.mocked(apiClient.fetchHeatmap).mockRejectedValueOnce(new Error('Network error'));

    await useAnalyticsStore.getState().fetchHeatmapData(2026);

    const state = useAnalyticsStore.getState();
    expect(state.isLoading).toBe(false);
    expect(state.error).toBe('Failed to fetch heatmap data');
  });
});
