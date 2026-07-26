import { useState } from 'react';
import type { StreakResponse, HabitStreak } from '@/types';
import { MilestoneCelebrationModal } from './MilestoneCelebrationModal';
import {
  Flame,
  Trophy,
  CalendarCheck,
  Award,
  Sparkles,
  Smile,
  Medal,
  Clock,
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface StreakDashboardProps {
  data: StreakResponse | null;
  isLoading: boolean;
}

export function StreakDashboard({ data, isLoading }: StreakDashboardProps) {
  const [celebration, setCelebration] = useState<{
    isOpen: boolean;
    habitTitle: string;
    milestone: number;
  }>({
    isOpen: false,
    habitTitle: '',
    milestone: 0,
  });

  if (isLoading) {
    return (
      <div className="rounded-2xl border border-border bg-card p-6 shadow-sm space-y-4 animate-pulse">
        <div className="h-6 w-48 bg-muted rounded" />
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-40 bg-muted/40 rounded-xl" />
          ))}
        </div>
      </div>
    );
  }

  if (!data || data.habits.length === 0) {
    return (
      <div className="rounded-2xl border border-border bg-card p-6 shadow-sm text-center text-muted-foreground">
        <div className="flex items-center justify-center gap-2 mb-2 text-primary">
          <Flame className="w-5 h-5" />
          <h3 className="font-bold text-base text-foreground">Habit Streaks</h3>
        </div>
        <p className="text-sm">No active habits tracked for streaks yet. Add active habits to start tracking streaks!</p>
      </div>
    );
  }

  const handleMilestoneClick = (habitTitle: string, milestone: number) => {
    setCelebration({
      isOpen: true,
      habitTitle,
      milestone,
    });
  };

  const getMilestoneIcon = (milestone: number) => {
    if (milestone >= 365) return Trophy;
    if (milestone >= 100) return Medal;
    if (milestone >= 30) return Award;
    return Sparkles;
  };

  return (
    <div className="rounded-2xl border border-border bg-card p-6 shadow-sm space-y-6">
      {/* Celebration Modal Popup (Option C) */}
      <MilestoneCelebrationModal
        isOpen={celebration.isOpen}
        onClose={() => setCelebration((prev) => ({ ...prev, isOpen: false }))}
        habitTitle={celebration.habitTitle}
        milestone={celebration.milestone}
      />

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Flame className="w-5 h-5 text-orange-500" />
          <h2 className="text-lg font-bold tracking-tight">Habit Streaks & Statistics</h2>
        </div>
        <div className="text-xs text-muted-foreground font-medium bg-muted/40 px-2.5 py-1 rounded-md">
          {data.habits.length} Active {data.habits.length === 1 ? 'Habit' : 'Habits'}
        </div>
      </div>

      {/* Habits Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {data.habits.map((h: HabitStreak) => (
          <div
            key={h.habitId}
            className="group relative rounded-xl border border-border/80 bg-card p-5 shadow-sm hover:border-orange-500/50 hover:shadow-md transition-all space-y-4"
          >
            {/* Top Bar: Title & Category */}
            <div className="flex items-start justify-between gap-2">
              <div>
                <h3 className="font-bold text-base text-foreground group-hover:text-orange-500 transition-colors">
                  {h.title}
                </h3>
                <div className="flex items-center gap-1.5 mt-1 text-xs text-muted-foreground">
                  <Smile className="w-3.5 h-3.5" />
                  <span>{h.category || 'General'}</span>
                </div>
              </div>
              <div className="flex items-center gap-1 px-2.5 py-1 rounded-full bg-orange-500/10 text-orange-500 border border-orange-500/20 text-xs font-extrabold">
                <Flame className="w-3.5 h-3.5 fill-orange-500" />
                <span>{h.currentStreak}d</span>
              </div>
            </div>

            {/* Main Stats Row */}
            <div className="grid grid-cols-3 gap-2 pt-1 border-t border-border/50 text-center">
              <div className="p-2 rounded-lg bg-muted/20">
                <div className="flex items-center justify-center gap-1 text-[10px] uppercase font-semibold text-muted-foreground">
                  <Flame className="w-3 h-3 text-orange-500" />
                  Current
                </div>
                <div className="text-base font-extrabold text-foreground mt-0.5">
                  {h.currentStreak} <span className="text-[10px] font-normal text-muted-foreground">days</span>
                </div>
              </div>

              <div className="p-2 rounded-lg bg-muted/20">
                <div className="flex items-center justify-center gap-1 text-[10px] uppercase font-semibold text-muted-foreground">
                  <Trophy className="w-3 h-3 text-amber-500" />
                  Longest
                </div>
                <div className="text-base font-extrabold text-foreground mt-0.5">
                  {h.longestStreak} <span className="text-[10px] font-normal text-muted-foreground">days</span>
                </div>
              </div>

              <div className="p-2 rounded-lg bg-muted/20">
                <div className="flex items-center justify-center gap-1 text-[10px] uppercase font-semibold text-muted-foreground">
                  <CalendarCheck className="w-3 h-3 text-emerald-500" />
                  Total
                </div>
                <div className="text-base font-extrabold text-foreground mt-0.5">
                  {h.totalDays} <span className="text-[10px] font-normal text-muted-foreground">days</span>
                </div>
              </div>
            </div>

            {/* Milestones Section */}
            {h.milestones.length > 0 && (
              <div className="space-y-1.5 pt-1">
                <span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground block">
                  Milestones Achieved
                </span>
                <div className="flex flex-wrap gap-1.5">
                  {h.milestones.map((m) => {
                    const Icon = getMilestoneIcon(m);
                    return (
                      <button
                        key={m}
                        onClick={() => handleMilestoneClick(h.title, m)}
                        className={cn(
                          'flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-semibold border transition-transform hover:scale-105',
                          m >= 365
                            ? 'bg-purple-500/10 text-purple-400 border-purple-500/30'
                            : m >= 100
                            ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30'
                            : m >= 30
                            ? 'bg-amber-500/10 text-amber-400 border-amber-500/30'
                            : 'bg-blue-500/10 text-blue-400 border-blue-500/30',
                        )}
                        title={`Click to celebrate ${m}-day milestone`}
                      >
                        <Icon className="w-3 h-3" />
                        <span>{m}d Streak</span>
                      </button>
                    );
                  })}
                </div>
              </div>
            )}

            {/* Footer: Last Completed */}
            {h.lastCompleted && (
              <div className="flex items-center justify-between text-[11px] text-muted-foreground pt-1 border-t border-border/40">
                <span className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  Last record
                </span>
                <span className="font-medium text-foreground">{h.lastCompleted}</span>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
