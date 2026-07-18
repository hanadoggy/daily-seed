import { Minus, Plus, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { useAppStore } from '@/store/useAppStore';
import type { TaskEntry } from '@/types';
import { cn } from '@/lib/utils';

interface TaskItemProps {
  entry: TaskEntry;
  title: string;
  type: 'quantitative' | 'boolean';
}

export function TaskItem({ entry, title, type }: TaskItemProps) {
  const updateTaskProgress = useAppStore((s) => s.updateTaskProgressOptimistic);

  const progressPercent =
    entry.targetAmount > 0
      ? Math.round((entry.actualAmount / entry.targetAmount) * 100)
      : 0;

  if (type === 'boolean') {
    return (
      <div
        className={cn(
          'flex items-center gap-3 rounded-xl px-4 py-3 transition-all duration-300',
          'bg-card border border-border',
          entry.isCompleted && 'border-mode-accent bg-mode-accent-soft',
        )}
      >
        <button
          onClick={() =>
            updateTaskProgress(entry.taskId, entry.isCompleted ? 0 : entry.targetAmount)
          }
          className={cn(
            'flex h-6 w-6 shrink-0 items-center justify-center rounded-md border-2 transition-all duration-300',
            entry.isCompleted
              ? 'border-mode-accent bg-mode-accent text-white'
              : 'border-muted-foreground/40 hover:border-mode-accent',
          )}
        >
          {entry.isCompleted && <Check className="h-3.5 w-3.5" />}
        </button>
        <span
          className={cn(
            'text-sm font-medium transition-colors duration-300',
            entry.isCompleted && 'text-mode-accent',
          )}
        >
          {title}
        </span>
      </div>
    );
  }

  return (
    <div
      className={cn(
        'rounded-xl px-4 py-3 transition-all duration-300',
        'bg-card border border-border',
        entry.isCompleted && 'border-mode-accent bg-mode-accent-soft',
      )}
    >
      <div className="flex items-center justify-between mb-2">
        <span
          className={cn(
            'text-sm font-medium transition-colors duration-300',
            entry.isCompleted && 'text-mode-accent',
          )}
        >
          {title}
        </span>
        <div className="flex items-center gap-1.5">
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 rounded-lg text-cancel hover:bg-cancel/10 hover:text-cancel"
            onClick={() => updateTaskProgress(entry.taskId, entry.actualAmount - 1)}
            disabled={entry.actualAmount <= 0}
          >
            <Minus className="h-3.5 w-3.5" />
          </Button>
          <span className="text-xs font-semibold tabular-nums min-w-[3rem] text-center">
            {entry.actualAmount} / {entry.targetAmount}
          </span>
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 rounded-lg text-save hover:bg-save/10 hover:text-save"
            onClick={() => updateTaskProgress(entry.taskId, entry.actualAmount + 1)}
            disabled={entry.actualAmount >= entry.targetAmount}
          >
            <Plus className="h-3.5 w-3.5" />
          </Button>
        </div>
      </div>
      <Progress
        value={progressPercent}
        className="h-1.5 bg-muted"
      />
    </div>
  );
}
