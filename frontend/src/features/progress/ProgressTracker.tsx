import { useEffect, useState } from 'react';
import { ArrowRightLeft, TrendingUp } from 'lucide-react';
import { useAppStore } from '@/store/useAppStore';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

export function ProgressTracker() {
  const { taskProgress, fetchProgress, migrateTask } = useAppStore();
  const [migratingId, setMigratingId] = useState<string | null>(null);

  useEffect(() => {
    fetchProgress();
  }, [fetchProgress]);

  const handleMigrate = async (taskId: string) => {
    setMigratingId(taskId);
    await migrateTask(taskId);
    setMigratingId(null);
  };

  if (taskProgress.length === 0) return null;

  return (
    <div className="rounded-2xl border border-border bg-card p-4 shadow-sm">
      <div className="flex items-center gap-2 mb-4">
        <TrendingUp className="h-4 w-4 text-muted-foreground" />
        <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          Project Progress
        </h3>
      </div>

      <div className="space-y-3">
        {taskProgress.map((tp) => {
          const isComplete = tp.percentage >= 100;
          const clampedPct = Math.min(tp.percentage, 100);

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
                  {tp.totalCompleted}/{tp.totalTarget}
                </span>
              </div>

              {/* Progress bar */}
              <div className="relative h-2 w-full overflow-hidden rounded-full bg-secondary">
                <div
                  className={cn(
                    'h-full rounded-full transition-all duration-700 ease-out',
                    isComplete
                      ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]'
                      : 'bg-primary',
                  )}
                  style={{ width: `${clampedPct}%` }}
                />
              </div>

              {/* Percentage + migrate button */}
              <div className="flex items-center justify-between">
                <span className="text-xs text-muted-foreground">
                  {Math.round(tp.percentage)}%
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
        })}
      </div>
    </div>
  );
}
