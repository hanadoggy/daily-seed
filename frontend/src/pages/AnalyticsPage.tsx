import { useEffect } from 'react';
import { HeatmapDashboard } from '@/features/analytics/HeatmapDashboard';
import { useAnalyticsStore } from '@/store/useAnalyticsStore';
import dayjs from 'dayjs';

export function AnalyticsPage() {
  const { heatmapData, isLoading, fetchHeatmapData, error } = useAnalyticsStore();
  const currentYear = dayjs().year();

  useEffect(() => {
    fetchHeatmapData(currentYear);
  }, [fetchHeatmapData, currentYear]);

  return (
    <main className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold tracking-tight">Analytics</h1>
        <div className="text-sm text-muted-foreground">
          {currentYear}
        </div>
      </div>

      {error && (
        <div className="rounded-xl border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <HeatmapDashboard data={heatmapData} isLoading={isLoading} />
      
      {/* Placeholder for future analytics features */}
      <div className="grid gap-6 md:grid-cols-2 mt-6">
        <div className="rounded-2xl border border-border bg-card p-6 shadow-sm flex items-center justify-center min-h-[200px] text-muted-foreground">
          Future: Streak Tracking
        </div>
        <div className="rounded-2xl border border-border bg-card p-6 shadow-sm flex items-center justify-center min-h-[200px] text-muted-foreground">
          Future: Task Correlations
        </div>
      </div>
    </main>
  );
}
