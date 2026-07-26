import { useEffect } from 'react';
import { Flame, Trophy } from 'lucide-react';
import { useAnalyticsStore } from '@/store/useAnalyticsStore';

export function HabitStreakWidget() {
  const { streakData, fetchStreakData, isStreakLoading } = useAnalyticsStore();

  useEffect(() => {
    fetchStreakData();
  }, [fetchStreakData]);

  if (isStreakLoading) {
    return (
      <div className="rounded-2xl border border-border bg-card p-4 shadow-sm space-y-3 animate-pulse">
        <div className="h-4 w-32 bg-muted rounded" />
        <div className="h-10 bg-muted/40 rounded-xl" />
      </div>
    );
  }

  if (!streakData || streakData.habits.length === 0) return null;

  return (
    <div className="rounded-2xl border border-border bg-card p-4 shadow-sm space-y-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Flame className="h-4 w-4 text-orange-500" />
          <h3 className="text-xs font-semibold uppercase tracking-wider text-orange-500">
            Habit Streaks
          </h3>
        </div>
        <span className="text-[10px] text-muted-foreground font-medium">
          {streakData.habits.length} habits
        </span>
      </div>

      <div className="space-y-2">
        {streakData.habits.map((h) => (
          <div
            key={h.habitId}
            className="flex items-center justify-between p-2 rounded-xl bg-muted/20 border border-border/40 text-xs"
          >
            <div className="truncate font-medium text-foreground max-w-[150px]">
              {h.title}
            </div>
            <div className="flex items-center gap-2">
              <div className="flex items-center gap-1 font-bold text-orange-500">
                <Flame className="w-3.5 h-3.5 fill-orange-500" />
                <span>{h.currentStreak}d</span>
              </div>
              <div className="flex items-center gap-1 text-muted-foreground text-[10px]" title="Longest streak">
                <Trophy className="w-3 h-3 text-amber-500" />
                <span>{h.longestStreak}d</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
