import { useEffect } from 'react';
import { HeatmapDashboard } from '@/features/analytics/HeatmapDashboard';
import { SummaryDashboard } from '@/features/analytics/SummaryDashboard';
import { StreakDashboard } from '@/features/analytics/StreakDashboard';
import { useAnalyticsStore } from '@/store/useAnalyticsStore';
import dayjs from 'dayjs';

export function AnalyticsPage() {
  const {
    heatmapData,
    isLoading,
    fetchHeatmapData,
    error,
    summaryData,
    summaryPeriod,
    isSummaryLoading,
    summaryError,
    fetchSummaryData,
    setSummaryPeriod,
    navigateSummary,
    streakData,
    isStreakLoading,
    streakError,
    fetchStreakData,
  } = useAnalyticsStore();

  const currentYear = dayjs().year();

  useEffect(() => {
    fetchHeatmapData(currentYear);
    fetchSummaryData();
    fetchStreakData();
  }, [fetchHeatmapData, fetchSummaryData, fetchStreakData, currentYear]);

  return (
    <main className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500 pb-12">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold tracking-tight">Analytics & Insights</h1>
        <div className="text-sm text-muted-foreground">{currentYear}</div>
      </div>

      {(error || summaryError || streakError) && (
        <div className="rounded-xl border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {error || summaryError || streakError}
        </div>
      )}

      {/* Annual Heatmap */}
      <HeatmapDashboard data={heatmapData} isLoading={isLoading} />

      {/* Weekly / Monthly Summary */}
      <SummaryDashboard
        data={summaryData}
        period={summaryPeriod}
        isLoading={isSummaryLoading}
        onPeriodChange={setSummaryPeriod}
        onNavigate={navigateSummary}
      />

      {/* Habit Streak Tracking */}
      <StreakDashboard data={streakData} isLoading={isStreakLoading} />
    </main>
  );
}
