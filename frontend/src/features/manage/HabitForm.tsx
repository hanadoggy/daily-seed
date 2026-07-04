import { useState } from 'react';
import type { Habit } from '@/types';
import { useAppStore } from '@/store/useAppStore';
import { Button } from '@/components/ui/button';

const CATEGORY_OPTIONS = [
  'mindfulness',
  'health',
  'productivity',
  'learning',
  'social',
  'other',
] as const;

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
        <label htmlFor="habit-category" className="text-xs font-medium text-muted-foreground">
          Category
        </label>
        <select
          id="habit-category"
          value={category}
          onChange={(e) => setCategory(e.target.value)}
          className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent"
        >
          {CATEGORY_OPTIONS.map((cat) => (
            <option key={cat} value={cat}>
              {cat.charAt(0).toUpperCase() + cat.slice(1)}
            </option>
          ))}
        </select>
      </div>

      <div className="flex gap-2 pt-2">
        <Button type="submit" disabled={!title.trim() || submitting} className="flex-1">
          {submitting ? 'Saving…' : isEditing ? 'Save Changes' : 'Create Habit'}
        </Button>
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
      </div>
    </form>
  );
}
