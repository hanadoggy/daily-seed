import { Check } from 'lucide-react';
import { useAppStore } from '@/store/useAppStore';
import type { HabitEntry } from '@/types';
import { cn } from '@/lib/utils';
import { HABIT_CATEGORIES } from '../manage/categoryOptions';

interface HabitItemProps {
  entry: HabitEntry;
  title: string;
  category?: string;
  isReadOnly?: boolean;
}

export function HabitItem({ entry, title, category, isReadOnly }: HabitItemProps) {
  const toggleHabit = useAppStore((s) => s.toggleHabitOptimistic);

  const categoryOption = HABIT_CATEGORIES.find(c => c.value === category);
  const CategoryIcon = categoryOption?.icon;

  return (
    <button
      onClick={() => {
        if (!isReadOnly) toggleHabit(entry.habitId, !entry.isCompleted);
      }}
      disabled={isReadOnly}
      className={cn(
        'flex w-full items-center gap-3 rounded-xl px-4 py-3 transition-all duration-300 text-left',
        'bg-card border border-border',
        !isReadOnly && 'hover:border-muted-foreground/50',
        entry.isCompleted && 'border-mode-accent bg-mode-accent-soft',
        isReadOnly && 'opacity-80 cursor-not-allowed',
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
      <div className="flex-1 min-w-0 flex items-center justify-between">
        <span
          className={cn(
            'text-sm font-medium transition-colors duration-300',
            entry.isCompleted && 'text-mode-accent',
          )}
        >
          {title}
        </span>
        {CategoryIcon && (
          <CategoryIcon 
            className={cn(
              'h-4 w-4 opacity-70 transition-opacity',
              entry.isCompleted ? 'opacity-100 text-mode-accent' : categoryOption.color
            )} 
          />
        )}
      </div>
    </button>
  );
}
