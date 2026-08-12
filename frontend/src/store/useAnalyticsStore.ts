import { create } from 'zustand';
import { toast } from 'sonner';
import { fetchHeatmap, fetchSummary, fetchStreaks, ApiError } from '@/api/client';
import { getErrorMessage } from '@/lib/errorMessages';
import type { HeatmapResponse, SummaryResponse, StreakResponse } from '@/types';
import dayjs from 'dayjs';

interface AnalyticsState {
  heatmapData: HeatmapResponse | null;
  isLoading: boolean;
  error: string | null;

  summaryData: SummaryResponse | null;
  summaryPeriod: 'weekly' | 'monthly';
  summaryDate: string; // YYYY-MM-DD
  isSummaryLoading: boolean;
  summaryError: string | null;

  streakData: StreakResponse | null;
  isStreakLoading: boolean;
  streakError: string | null;

  fetchHeatmapData: (year: number) => Promise<void>;
  fetchSummaryData: (period?: 'weekly' | 'monthly', date?: string) => Promise<void>;
  fetchStreakData: () => Promise<void>;
  setSummaryPeriod: (period: 'weekly' | 'monthly') => void;
  navigateSummary: (direction: 'prev' | 'next') => void;
}

export const useAnalyticsStore = create<AnalyticsState>((set, get) => ({
  heatmapData: null,
  isLoading: false,
  error: null,

  summaryData: null,
  summaryPeriod: 'weekly',
  summaryDate: dayjs().format('YYYY-MM-DD'),
  isSummaryLoading: false,
  summaryError: null,

  streakData: null,
  isStreakLoading: false,
  streakError: null,

  fetchHeatmapData: async (year: number) => {
    set({ isLoading: true, error: null });
    try {
      const data = await fetchHeatmap(year);
      set({ heatmapData: data, isLoading: false });
    } catch (err: unknown) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchHeatmap', status);
      toast.error(message);
      set({
        isLoading: false,
        error: message,
      });
    }
  },

  fetchSummaryData: async (period, date) => {
    const currentPeriod = period || get().summaryPeriod;
    const currentDate = date || get().summaryDate;

    set({ isSummaryLoading: true, summaryError: null });
    try {
      const data = await fetchSummary(currentPeriod, currentDate);
      set({
        summaryData: data,
        summaryPeriod: currentPeriod,
        summaryDate: currentDate,
        isSummaryLoading: false,
      });
    } catch (err: unknown) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchSummary', status);
      toast.error(message);
      set({
        isSummaryLoading: false,
        summaryError: message,
      });
    }
  },

  fetchStreakData: async () => {
    set({ isStreakLoading: true, streakError: null });
    try {
      const data = await fetchStreaks();
      set({ streakData: data, isStreakLoading: false });
    } catch (err: unknown) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchStreaks', status);
      toast.error(message);
      set({
        isStreakLoading: false,
        streakError: message,
      });
    }
  },

  setSummaryPeriod: (period) => {
    set({ summaryPeriod: period });
    get().fetchSummaryData(period, get().summaryDate);
  },

  navigateSummary: (direction) => {
    const { summaryPeriod, summaryDate } = get();
    const current = dayjs(summaryDate);
    const unit = summaryPeriod === 'weekly' ? 'week' : 'month';
    const nextDate = direction === 'prev'
      ? current.subtract(1, unit).format('YYYY-MM-DD')
      : current.add(1, unit).format('YYYY-MM-DD');

    set({ summaryDate: nextDate });
    get().fetchSummaryData(summaryPeriod, nextDate);
  },
}));
