import { cn } from '@/lib/utils';
import type { SummaryResponse } from '@/types';
import { MODE_OPTIONS } from '@/features/context-mode/conditionOptions';
import {
  Calendar,
  ChevronLeft,
  ChevronRight,
  CheckCircle2,
  Target,
  Compass,
  BookOpen,
  TrendingUp,
  FileText,
  ListTodo,
  Smile,
} from 'lucide-react';

interface SummaryDashboardProps {
  data: SummaryResponse | null;
  period: 'weekly' | 'monthly';
  isLoading: boolean;
  onPeriodChange: (period: 'weekly' | 'monthly') => void;
  onNavigate: (direction: 'prev' | 'next') => void;
}

export function SummaryDashboard({
  data,
  period,
  isLoading,
  onPeriodChange,
  onNavigate,
}: SummaryDashboardProps) {
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="rounded-2xl border border-border bg-card p-6 shadow-sm animate-pulse space-y-4">
          <div className="h-6 w-48 bg-muted rounded" />
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="h-28 bg-muted/40 rounded-xl" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="rounded-2xl border border-border bg-card p-8 text-center text-muted-foreground shadow-sm">
        Summary data unavailable.
      </div>
    );
  }

  const recordRate =
    data.totalDays > 0 ? Math.round((data.recordedDays / data.totalDays) * 100) : 0;

  return (
    <div className="space-y-6">
      {/* Summary Header & Navigation */}
      <div className="rounded-2xl border border-border bg-card p-6 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <TrendingUp className="w-5 h-5 text-primary" />
            <h2 className="text-lg font-bold tracking-tight">Period Summary</h2>
          </div>
          <p className="text-sm text-muted-foreground">
            {data.startDate} ~ {data.endDate}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          {/* Weekly / Monthly Tab */}
          <div className="flex items-center gap-1 p-1 bg-muted/40 rounded-lg">
            <button
              onClick={() => onPeriodChange('weekly')}
              className={cn(
                'px-3 py-1.5 text-xs font-medium rounded-md transition-all',
                period === 'weekly'
                  ? 'bg-primary text-primary-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted/50',
              )}
            >
              Weekly
            </button>
            <button
              onClick={() => onPeriodChange('monthly')}
              className={cn(
                'px-3 py-1.5 text-xs font-medium rounded-md transition-all',
                period === 'monthly'
                  ? 'bg-primary text-primary-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted/50',
              )}
            >
              Monthly
            </button>
          </div>

          {/* Prev / Next Navigation */}
          <div className="flex items-center gap-1">
            <button
              onClick={() => onNavigate('prev')}
              aria-label="Previous period"
              className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded-md transition-colors"
            >
              <ChevronLeft className="w-5 h-5" />
            </button>
            <button
              onClick={() => onNavigate('next')}
              aria-label="Next period"
              className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded-md transition-colors"
            >
              <ChevronRight className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      {/* Top Stat Overview Cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {/* Task Completion */}
        <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Task Avg
            </span>
            <CheckCircle2 className="w-4 h-4 text-emerald-400" />
          </div>
          <div className="text-2xl font-extrabold">{data.taskCompletion.overall}%</div>
          <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
            <div
              className="bg-emerald-400 h-2 rounded-full transition-all"
              style={{ width: `${Math.min(100, data.taskCompletion.overall)}%` }}
            />
          </div>
        </div>

        {/* Habit Completion */}
        <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Habit Avg
            </span>
            <Target className="w-4 h-4 text-amber-500" />
          </div>
          <div className="text-2xl font-extrabold">{data.habitCompletion.overall}%</div>
          <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
            <div
              className="bg-amber-500 h-2 rounded-full transition-all"
              style={{ width: `${Math.min(100, data.habitCompletion.overall)}%` }}
            />
          </div>
        </div>

        {/* Record Rate */}
        <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Recorded Days
            </span>
            <Calendar className="w-4 h-4 text-purple-500" />
          </div>
          <div className="text-2xl font-extrabold">
            {data.recordedDays} <span className="text-sm text-muted-foreground font-normal">/ {data.totalDays} days</span>
          </div>
          <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
            <div
              className="bg-purple-500 h-2 rounded-full transition-all"
              style={{ width: `${recordRate}%` }}
            />
          </div>
        </div>

        {/* Primary Context Mode */}
        <div className="rounded-xl border border-border bg-card p-5 shadow-sm space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Modes Logged
            </span>
            <Compass className="w-4 h-4 text-amber-500" />
          </div>
          <div className="flex items-center gap-1.5 flex-wrap">
            {Object.keys(data.modeDistribution).length > 0 ? (
              Object.entries(data.modeDistribution).map(([mode, count]) => {
                const modeOpt = MODE_OPTIONS.find((m) => m.value === mode);
                const Icon = modeOpt?.icon;
                return (
                  <span
                    key={mode}
                    className={cn(
                      'flex items-center gap-1 px-2.5 py-1 text-xs font-medium rounded-md border',
                      modeOpt
                        ? `${modeOpt.bgColor} ${modeOpt.color} ${modeOpt.borderColor}`
                        : 'bg-muted/60 text-foreground border-border',
                    )}
                  >
                    {Icon && <Icon className="w-3.5 h-3.5" />}
                    <span>
                      {mode}: {count}d
                    </span>
                  </span>
                );
              })
            ) : (
              <span className="text-sm text-muted-foreground">No mode recorded</span>
            )}
          </div>
        </div>
      </div>

      {/* Detailed Stats Grid: Tasks & Habits */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Task Details */}
        <div className="rounded-2xl border border-border bg-card p-6 shadow-sm space-y-4">
          <div className="flex items-center gap-2">
            <ListTodo className="w-5 h-5 text-primary" />
            <h3 className="font-bold text-base">Task Breakdown</h3>
          </div>

          {/* Sections Breakdown */}
          {Object.keys(data.taskCompletion.sections).length > 0 && (
            <div className="space-y-2 pt-2">
              <span className="text-xs font-semibold text-muted-foreground uppercase">Section Averages</span>
              <div className="grid grid-cols-2 gap-2">
                {Object.entries(data.taskCompletion.sections).map(([sec, rate]) => (
                  <div key={sec} className="p-3 rounded-lg bg-muted/30 border border-border/50">
                    <div className="text-xs text-muted-foreground capitalize">{sec}</div>
                    <div className="text-lg font-bold text-emerald-400">{rate}%</div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Per-Task List */}
          <div className="space-y-2 pt-2">
            <span className="text-xs font-semibold text-muted-foreground uppercase">Individual Tasks</span>
            {data.taskCompletion.perTask.length > 0 ? (
              <div className="space-y-3">
                {data.taskCompletion.perTask.map((t) => (
                  <div key={t.taskId} className="space-y-1">
                    <div className="flex items-center justify-between text-xs">
                      <span className="font-medium truncate max-w-[200px]">{t.title}</span>
                      <span className="text-muted-foreground font-semibold">{t.rate}% ({t.completed}/{t.target})</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-1.5 overflow-hidden">
                      <div
                        className="bg-emerald-400 h-1.5 rounded-full transition-all"
                        style={{ width: `${Math.min(100, t.rate)}%` }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-xs text-muted-foreground py-2">No task records in this period.</p>
            )}
          </div>
        </div>

        {/* Habit Details */}
        <div className="rounded-2xl border border-border bg-card p-6 shadow-sm space-y-4">
          <div className="flex items-center gap-2">
            <Smile className="w-5 h-5 text-primary" />
            <h3 className="font-bold text-base">Habit Breakdown</h3>
          </div>

          <div className="space-y-2 pt-2">
            <span className="text-xs font-semibold text-muted-foreground uppercase">Individual Habits</span>
            {data.habitCompletion.perHabit.length > 0 ? (
              <div className="space-y-3">
                {data.habitCompletion.perHabit.map((h) => (
                  <div key={h.habitId} className="space-y-1">
                    <div className="flex items-center justify-between text-xs">
                      <span className="font-medium truncate max-w-[200px]">{h.title}</span>
                      <span className="text-muted-foreground font-semibold">{h.rate}% ({h.completed}/{h.total}d)</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-1.5 overflow-hidden">
                      <div
                        className="bg-amber-500 h-1.5 rounded-full transition-all"
                        style={{ width: `${Math.min(100, h.rate)}%` }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-xs text-muted-foreground py-2">No habit records in this period.</p>
            )}
          </div>
        </div>
      </div>

      {/* Journal Entries Timeline */}
      <div className="rounded-2xl border border-border bg-card p-6 shadow-sm space-y-4">
        <div className="flex items-center gap-2">
          <BookOpen className="w-5 h-5 text-primary" />
          <h3 className="font-bold text-base">Journal Entries</h3>
        </div>

        {data.journals.length > 0 ? (
          <div className="grid gap-4 sm:grid-cols-2">
            {data.journals.map((j) => (
              <div
                key={j.date}
                className="rounded-xl border border-border/80 bg-muted/20 p-4 space-y-2 text-sm"
              >
                <div className="flex items-center justify-between border-b border-border/50 pb-2">
                  <span className="font-semibold text-primary">{j.date}</span>
                  <FileText className="w-4 h-4 text-muted-foreground" />
                </div>
                {j.oneLineReview && (
                  <div>
                    <span className="text-xs font-semibold text-muted-foreground block">One Line Review</span>
                    <p className="italic text-foreground">{j.oneLineReview}</p>
                  </div>
                )}
                {j.threeLineDiary && (
                  <div>
                    <span className="text-xs font-semibold text-muted-foreground block">Three Line Diary</span>
                    <p className="whitespace-pre-line text-muted-foreground text-xs leading-relaxed">
                      {j.threeLineDiary}
                    </p>
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">No journal entries written during this period.</p>
        )}
      </div>
    </div>
  );
}
