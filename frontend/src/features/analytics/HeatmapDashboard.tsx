import { useState, useMemo } from 'react';
import { cn } from '@/lib/utils';
import type { HeatmapResponse } from '@/types';
import dayjs from 'dayjs';

interface HeatmapDashboardProps {
  data: HeatmapResponse | null;
  isLoading: boolean;
}

type FilterType = 'overall' | 'dev' | 'japanese' | 'self_dev' | 'exercise';

const FILTERS: { value: FilterType; label: string }[] = [
  { value: 'overall', label: 'Overall' },
  { value: 'dev', label: 'Dev' },
  { value: 'japanese', label: 'Japanese' },
  { value: 'self_dev', label: 'Self Dev' },
  { value: 'exercise', label: 'Exercise' },
];

export function HeatmapDashboard({ data, isLoading }: HeatmapDashboardProps) {
  const [activeFilter, setActiveFilter] = useState<FilterType>('overall');

  const { maxIntensity, weeks } = useMemo(() => {
    if (!data || !data.days) return { maxIntensity: 0, weeks: [] };

    let max = 1;
    const allDays = data.days;
    
    // Calculate intensity per day based on active filter
    const daysWithIntensity = allDays.map(day => {
      let intensity = 0;
      if (activeFilter === 'overall') {
        intensity = day.total;
      } else {
        intensity = day.sectionCounts?.[activeFilter] || 0;
      }
      if (intensity > max) max = intensity;
      return { ...day, intensity };
    });

    // Group by weeks
    // Find the first day of the year to pad the first week if necessary
    const firstDay = dayjs(daysWithIntensity[0].date);
    const firstDayOfWeek = firstDay.day(); // 0 = Sunday, 1 = Monday...

    const paddedDays: any[] = [];
    // Pad start with empty days so it aligns with Sunday
    for (let i = 0; i < firstDayOfWeek; i++) {
      paddedDays.push(null);
    }
    paddedDays.push(...daysWithIntensity);

    const groupedWeeks = [];
    for (let i = 0; i < paddedDays.length; i += 7) {
      groupedWeeks.push(paddedDays.slice(i, i + 7));
    }

    return { maxIntensity: max, weeks: groupedWeeks };
  }, [data, activeFilter]);

  const getColorClass = (intensity: number) => {
    if (intensity === 0) return 'bg-muted/30';
    
    // Calculate intensity level 1-4
    const ratio = intensity / maxIntensity;
    if (ratio <= 0.25) return 'bg-mode-accent/30';
    if (ratio <= 0.5) return 'bg-mode-accent/60';
    if (ratio <= 0.75) return 'bg-mode-accent/80';
    return 'bg-mode-accent';
  };

  return (
    <div className="rounded-2xl border border-border bg-card p-6 shadow-sm overflow-hidden flex flex-col">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h2 className="text-lg font-bold tracking-tight">Consistency Heatmap</h2>
          <p className="text-sm text-muted-foreground">Your daily progress across different areas</p>
        </div>
        
        {/* Filters */}
        <div className="flex flex-wrap items-center gap-1.5 p-1 bg-muted/40 rounded-lg w-fit">
          {FILTERS.map((f) => (
            <button
              key={f.value}
              onClick={() => setActiveFilter(f.value)}
              className={cn(
                'px-3 py-1.5 text-xs font-medium rounded-md transition-all',
                activeFilter === f.value
                  ? 'bg-primary text-primary-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
              )}
            >
              {f.label}
            </button>
          ))}
        </div>
      </div>

      <div className="relative w-full overflow-x-auto pt-8 pb-6 custom-scrollbar">
        {isLoading ? (
          <div className="flex h-[120px] items-center justify-center text-muted-foreground text-sm">
            Loading...
          </div>
        ) : (
          <div className="flex gap-1 min-w-max mx-auto w-fit">
            {weeks.map((week, wIdx) => (
              <div key={wIdx} className="flex flex-col gap-1">
                {week.map((day, dIdx) => {
                  if (!day) {
                    return <div key={`empty-${wIdx}-${dIdx}`} className="w-3 h-3 rounded-sm bg-transparent" />;
                  }
                  
                  const isTopRow = dIdx < 3;

                  return (
                    <div
                      key={day.date}
                      className="group relative"
                    >
                      <div
                        className={cn(
                          "w-3 h-3 rounded-[2px] transition-colors hover:ring-1 hover:ring-ring hover:ring-offset-1 hover:ring-offset-background",
                          getColorClass(day.intensity)
                        )}
                      />
                      {/* Tooltip */}
                      <div
                        className={cn(
                          "pointer-events-none absolute left-1/2 z-[100] -translate-x-1/2 whitespace-nowrap rounded-md bg-zinc-900 px-2 py-1 text-xs font-medium text-zinc-50 opacity-0 transition-opacity group-hover:opacity-100 dark:bg-zinc-100 dark:text-zinc-900 shadow-xl",
                          isTopRow ? "top-full mt-2" : "bottom-full mb-2"
                        )}
                      >
                        {day.date}: {day.intensity} completions
                        <div
                          className={cn(
                            "absolute left-1/2 -translate-x-1/2 border-4 border-transparent",
                            isTopRow
                              ? "-top-1 border-b-zinc-900 dark:border-b-zinc-100"
                              : "-bottom-1 border-t-zinc-900 dark:border-t-zinc-100"
                          )}
                        />
                      </div>
                    </div>
                  );
                })}
              </div>
            ))}
          </div>
        )}
      </div>
      
      {/* Legend */}
      <div className="mt-4 flex items-center justify-end gap-2 text-xs text-muted-foreground">
        <span>Less</span>
        <div className="flex gap-1">
          <div className="w-3 h-3 rounded-[2px] bg-muted/30" />
          <div className="w-3 h-3 rounded-[2px] bg-mode-accent/30" />
          <div className="w-3 h-3 rounded-[2px] bg-mode-accent/60" />
          <div className="w-3 h-3 rounded-[2px] bg-mode-accent/80" />
          <div className="w-3 h-3 rounded-[2px] bg-mode-accent" />
        </div>
        <span>More</span>
      </div>
    </div>
  );
}
