import { Check } from 'lucide-react';
import { useAppStore } from '@/store/useAppStore';
import type { HabitEntry } from '@/types';
import { cn } from '@/lib/utils';

interface HabitItemProps {
  entry: HabitEntry;
  title: string;
}

export function HabitItem({ entry, title }: HabitItemProps) {
  const toggleHabit = useAppStore((s) => s.toggleHabitOptimistic);

  return (
    <button
      onClick={() => toggleHabit(entry.habitId, !entry.isCompleted)}
      className={cn(
        'flex w-full items-center gap-3 rounded-xl px-4 py-3 transition-all duration-300 text-left',
        'bg-card border border-border',
        'hover:border-muted-foreground/50',
        entry.isCompleted && 'border-mode-accent bg-mode-accent-soft',
      )}
    >
      <div
        className={cn(
          'flex h-6 w-6 shrink-0 items-center justify-center rounded-full border-2 transition-all duration-300',
          entry.isCompleted
            ? 'border-mode-accent bg-mode-accent text-white scale-110'
            : 'border-muted-foreground/40',
        )}
      >
        {entry.isCompleted && <Check className="h-3.5 w-3.5" />}
      </div>
      <span
        className={cn(
          'text-sm font-medium transition-colors duration-300',
          entry.isCompleted && 'text-mode-accent',
        )}
      >
        {title}
      </span>
    </button>
  );
}
