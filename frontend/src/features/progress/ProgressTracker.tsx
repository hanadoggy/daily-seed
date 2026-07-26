import { useEffect, useState } from 'react';
import { ArrowRightLeft, TrendingUp, Activity } from 'lucide-react';
import { useAppStore } from '@/store/useAppStore';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import type { TaskProgress } from '@/types';
import { HabitStreakWidget } from './HabitStreakWidget';

export function ProgressTracker() {
  const { taskProgress, fetchProgress, migrateTask } = useAppStore();
  const [migratingId, setMigratingId] = useState<string | null>(null);
  const [randomPcts, setRandomPcts] = useState<Record<string, number>>({});

  useEffect(() => {
    fetchProgress();
  }, [fetchProgress]);

  useEffect(() => {
    setRandomPcts((prev) => {
      const next = { ...prev };
      let changed = false;
      taskProgress.forEach((tp) => {
        const isEndless = tp.type === 'boolean' || tp.totalTarget === 0;
        if (isEndless && next[tp.taskId] === undefined) {
          next[tp.taskId] = Math.floor(Math.random() * 61) + 20; // 20 to 80
          changed = true;
        }
      });
      return changed ? next : prev;
    });
  }, [taskProgress]);

  const handleMigrate = async (taskId: string) => {
    setMigratingId(taskId);
    await migrateTask(taskId);
    setMigratingId(null);
  };

  if (taskProgress.length === 0) return null;

  const targetedTasks = taskProgress.filter(tp => tp.type === 'quantitative' && tp.totalTarget > 0);
  const endlessTasks = taskProgress.filter(tp => tp.type === 'boolean' || tp.totalTarget === 0);

  const renderTask = (tp: TaskProgress, isEndless: boolean) => {
    const isComplete = !isEndless && tp.percentage >= 100;
    let clampedPct = Math.min(tp.percentage, 100);

    if (isEndless) {
      clampedPct = tp.totalCompleted === 0 ? 0 : (randomPcts[tp.taskId] || 0);
    }

    return (
      <div key={tp.taskId} className="space-y-1.5">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium text-foreground truncate mr-2">
            {tp.title}
          </span>
          <span
            className={cn(
              'text-xs font-semibold tabular-nums',
              isComplete ? 'text-green-500' : 'text-muted-foreground',
            )}
          >
            {isEndless ? tp.totalCompleted : `${tp.totalCompleted}/${tp.totalTarget}`}
          </span>
        </div>

        {/* Progress bar */}
        <div className="relative h-2 w-full overflow-hidden rounded-full bg-secondary">
          <div
            className={cn(
              'h-full rounded-full transition-all duration-700 ease-out',
              isComplete
                ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]'
                : isEndless ? 'bg-accent' : 'bg-primary',
            )}
            style={{ width: `${clampedPct}%` }}
          />
        </div>

        {/* Percentage + migrate button */}
        <div className="flex items-center justify-between">
          <span className="text-xs text-muted-foreground min-h-[16px]">
            {!isEndless && `${Math.round(tp.percentage)}%`}
          </span>

          {isComplete && (
            <Button
              id={`migrate-${tp.taskId}`}
              variant="ghost"
              size="sm"
              className="h-6 gap-1 px-2 text-xs text-green-500 hover:text-green-400 hover:bg-green-500/10"
              disabled={migratingId === tp.taskId}
              onClick={() => handleMigrate(tp.taskId)}
            >
              <ArrowRightLeft className="h-3 w-3" />
              {migratingId === tp.taskId ? 'Migrating…' : 'Migrate'}
            </Button>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="space-y-4">
      {targetedTasks.length > 0 && (
        <div className="rounded-2xl border border-border bg-card p-4 shadow-sm">
          <div className="flex items-center gap-2 mb-4">
            <TrendingUp className="h-4 w-4 text-primary" />
            <h3 className="text-xs font-semibold uppercase tracking-wider text-primary">
              Project Progress
            </h3>
          </div>
          <div className="space-y-3">
            {targetedTasks.map(tp => renderTask(tp, false))}
          </div>
        </div>
      )}

      {endlessTasks.length > 0 && (
        <div className="rounded-2xl border border-border bg-card p-4 shadow-sm">
          <div className="flex items-center gap-2 mb-4">
            <Activity className="h-4 w-4 text-accent" />
            <h3 className="text-xs font-semibold uppercase tracking-wider text-accent">
              Continuous Progress
            </h3>
          </div>
          <div className="space-y-3">
            {endlessTasks.map(tp => renderTask(tp, true))}
          </div>
        </div>
      )}

      {/* Habit Streak Widget */}
      <HabitStreakWidget />
    </div>
  );
}
