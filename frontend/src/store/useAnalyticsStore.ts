import { create } from 'zustand';
import { fetchHeatmap } from '@/api/client';
import type { HeatmapResponse } from '@/types';

interface AnalyticsState {
  heatmapData: HeatmapResponse | null;
  isLoading: boolean;
  error: string | null;
  
  fetchHeatmapData: (year: number) => Promise<void>;
}

export const useAnalyticsStore = create<AnalyticsState>((set) => ({
  heatmapData: null,
  isLoading: false,
  error: null,

  fetchHeatmapData: async (year: number) => {
    set({ isLoading: true, error: null });
    try {
      const data = await fetchHeatmap(year);
      set({ heatmapData: data, isLoading: false });
    } catch (err: any) {
      set({
        isLoading: false,
        error: err.response?.data?.message || 'Failed to fetch heatmap data',
      });
    }
  },
}));
