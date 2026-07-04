import { useState } from 'react';
import type { Task } from '@/types';
import { useAppStore } from '@/store/useAppStore';
import { Button } from '@/components/ui/button';

const SECTION_OPTIONS = [
  { value: 'japanese', label: '🇯🇵 Japanese' },
  { value: 'dev', label: '💻 Development' },
  { value: 'self_dev', label: '📚 Self Development' },
] as const;

const TYPE_OPTIONS = [
  { value: 'quantitative', label: 'Quantitative' },
  { value: 'boolean', label: 'Boolean (Yes/No)' },
] as const;

interface TaskFormProps {
  task?: Task;
  onClose: () => void;
}

export function TaskForm({ task, onClose }: TaskFormProps) {
  const { addTask, editTask } = useAppStore();
  const [title, setTitle] = useState(task?.title ?? '');
  const [section, setSection] = useState<string>(task?.section ?? 'japanese');
  const [type, setType] = useState<string>(task?.type ?? 'quantitative');
  const [dailyTarget, setDailyTarget] = useState(task?.metrics.dailyTarget ?? 1);
  const [totalTarget, setTotalTarget] = useState(task?.metrics.totalTarget ?? 0);
  const [submitting, setSubmitting] = useState(false);

  const isEditing = !!task;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim()) return;

    setSubmitting(true);
    const payload = {
      title: title.trim(),
      section: section as Task['section'],
      type: type as Task['type'],
      metrics: {
        dailyTarget: type === 'boolean' ? 1 : dailyTarget,
        totalTarget,
      },
      conditions: {
        weather: task?.conditions.weather ?? 'any',
        mode: task?.conditions.mode ?? 'any',
      },
    };

    if (isEditing) {
      await editTask(task.id, payload);
    } else {
      await addTask(payload);
    }
    setSubmitting(false);
    onClose();
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <h3 className="text-sm font-semibold">{isEditing ? 'Edit Task' : 'New Task'}</h3>

      <div className="space-y-1.5">
        <label htmlFor="task-title" className="text-xs font-medium text-muted-foreground">
          Title
        </label>
        <input
          id="task-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="e.g. Memorize Kanji"
          required
          className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent focus:ring-1 focus:ring-mode-accent/40"
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <label htmlFor="task-section" className="text-xs font-medium text-muted-foreground">
            Section
          </label>
          <select
            id="task-section"
            value={section}
            onChange={(e) => setSection(e.target.value)}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent"
          >
            {SECTION_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>

        <div className="space-y-1.5">
          <label htmlFor="task-type" className="text-xs font-medium text-muted-foreground">
            Type
          </label>
          <select
            id="task-type"
            value={type}
            onChange={(e) => setType(e.target.value)}
            className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent"
          >
            {TYPE_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {type === 'quantitative' && (
        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1.5">
            <label htmlFor="task-daily-target" className="text-xs font-medium text-muted-foreground">
              Daily Target
            </label>
            <input
              id="task-daily-target"
              type="number"
              min={1}
              value={dailyTarget}
              onChange={(e) => setDailyTarget(Number(e.target.value))}
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent"
            />
          </div>
          <div className="space-y-1.5">
            <label htmlFor="task-total-target" className="text-xs font-medium text-muted-foreground">
              Total Target
            </label>
            <input
              id="task-total-target"
              type="number"
              min={0}
              value={totalTarget}
              onChange={(e) => setTotalTarget(Number(e.target.value))}
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-mode-accent"
            />
          </div>
        </div>
      )}

      <div className="flex gap-2 pt-2">
        <Button type="submit" disabled={!title.trim() || submitting} className="flex-1">
          {submitting ? 'Saving…' : isEditing ? 'Save Changes' : 'Create Task'}
        </Button>
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
      </div>
    </form>
  );
}
