import { useState } from 'react';
import type { Habit } from '@/types';
import { useAppStore } from '@/store/useAppStore';
import { Button } from '@/components/ui/button';

import { HABIT_CATEGORIES } from './categoryOptions';
import { cn } from '@/lib/utils';

interface HabitFormProps {
  habit?: Habit;
  onClose: () => void;
}

export function HabitForm({ habit, onClose }: HabitFormProps) {
  const { addHabit, editHabit } = useAppStore();
  const [title, setTitle] = useState(habit?.title ?? '');
  const [category, setCategory] = useState(habit?.category ?? 'mindfulness');
  const [submitting, setSubmitting] = useState(false);

  const isEditing = !!habit;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim()) return;

    setSubmitting(true);
    const payload = {
      title: title.trim(),
      category,
    };

    if (isEditing) {
      await editHabit(habit.id, payload);
    } else {
      await addHabit(payload);
    }
    setSubmitting(false);
    onClose();
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <h3 className="text-sm font-semibold">{isEditing ? 'Edit Habit' : 'New Habit'}</h3>

      <div className="space-y-1.5">
        <label htmlFor="habit-title" className="text-xs font-medium text-muted-foreground">
          Title
        </label>
        <input
          id="habit-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="e.g. Morning stretching routine"
          required
          className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent focus:ring-1 focus:ring-mode-accent/40"
        />
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-medium text-muted-foreground">
          Category
        </label>
        <div className="grid grid-cols-2 gap-2">
          {HABIT_CATEGORIES.map((cat) => {
            const isSelected = category === cat.value;
            return (
              <button
                key={cat.value}
                type="button"
                onClick={() => setCategory(cat.value)}
                className={cn(
                  'flex items-center gap-2 rounded-lg border px-3 py-2 text-sm transition-colors',
                  isSelected
                    ? `${cat.bgColor} ${cat.borderColor} text-foreground`
                    : 'border-border bg-background text-muted-foreground hover:bg-muted/50 hover:text-foreground'
                )}
              >
                <cat.icon className={cn('h-4 w-4', cat.color)} />
                {cat.label}
              </button>
            );
          })}
        </div>
      </div>

      <div className="flex gap-2 pt-2">
        <Button type="submit" disabled={!title.trim() || submitting} className="flex-1 bg-save hover:bg-save/90 text-white">
          {submitting ? 'Saving…' : isEditing ? 'Save Changes' : 'Create Habit'}
        </Button>
        <Button type="button" variant="outline" onClick={onClose} className="bg-cancel/15 text-cancel border-cancel/40 hover:bg-cancel/25 hover:text-cancel">
          Cancel
        </Button>
      </div>
    </form>
  );
}
