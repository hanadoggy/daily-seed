import { useState } from 'react';
import { toast } from 'sonner';
import type { Task } from '@/types';
import { useAppStore } from '@/store/useAppStore';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { MODE_OPTIONS, WEATHER_OPTIONS } from '../context-mode/conditionOptions';

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
  const { addTask, editTask, selectedDate } = useAppStore();
  const [title, setTitle] = useState(task?.title ?? '');
  const [section, setSection] = useState<string>(task?.section ?? 'japanese');
  const [type, setType] = useState<string>(task?.type ?? 'quantitative');
  const [dailyTarget, setDailyTarget] = useState(task?.metrics.dailyTarget ?? 1);
  const [totalTarget, setTotalTarget] = useState(task?.metrics.totalTarget ?? 0);
  const [weather, setWeather] = useState<string[]>(task?.conditions.weather ?? ['sunny', 'rainy']);
  const [mode, setMode] = useState<string[]>(task?.conditions.mode ?? ['Growth', 'Rest', 'Office', 'Remote']);
  const [startDate, setStartDate] = useState<string>(task?.startDate ?? selectedDate);
  const [submitting, setSubmitting] = useState(false);

  const toggleWeather = (val: string) => {
    if (weather.includes(val)) {
      if (weather.length <= 1) {
        toast.error("최소 1개의 옵션을 선택해야 합니다", { duration: 5000 });
        return;
      }
      setWeather(weather.filter((w) => w !== val));
    } else {
      setWeather([...weather, val]);
    }
  };

  const toggleMode = (val: string) => {
    if (mode.includes(val)) {
      if (mode.length <= 1) {
        toast.error("최소 1개의 옵션을 선택해야 합니다", { duration: 5000 });
        return;
      }
      setMode(mode.filter((m) => m !== val));
    } else {
      setMode([...mode, val]);
    }
  };

  const isEditing = !!task;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim()) return;

    setSubmitting(true);
    if (isEditing && task?.startDate && startDate > task.startDate) {
      if (!window.confirm('시작일을 늦출 경우 이전 날짜의 체크리스트 기록이 모두 삭제됩니다. 정말 변경하시겠습니까?')) {
        setSubmitting(false);
        return;
      }
    }

    const payload = {
      title: title.trim(),
      section: section as Task['section'],
      type: type as Task['type'],
      metrics: {
        dailyTarget: type === 'boolean' ? 1 : dailyTarget,
        totalTarget,
      },
      conditions: {
        weather,
        mode,
      },
      startDate,
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

      <div className="space-y-1.5">
        <label htmlFor="task-start-date" className="text-xs font-medium text-muted-foreground">
          Start Date
        </label>
        <input
          id="task-start-date"
          type="date"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
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

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-muted-foreground">
            Weather Condition
          </label>
          <div className="grid grid-cols-1 gap-2">
            {WEATHER_OPTIONS.map(({ value, label, icon: Icon, color }) => {
              const isActive = weather.includes(value);
              return (
                <Button
                  key={value}
                  type="button"
                  variant="outline"
                  onClick={() => toggleWeather(value)}
                  className={cn(
                    'w-full gap-2 transition-all duration-300 border',
                    isActive
                      ? 'bg-mode-accent-soft border-mode-accent text-mode-accent shadow-sm'
                      : 'border-border text-muted-foreground hover:text-foreground hover:border-muted-foreground',
                  )}
                >
                  <Icon className={cn('h-4 w-4', isActive && color)} />
                  <span className="text-sm font-medium">{label}</span>
                </Button>
              );
            })}
          </div>
        </div>

        <div className="space-y-1.5">
          <label className="text-xs font-medium text-muted-foreground">
            Context Mode
          </label>
          <div className="grid grid-cols-2 gap-2">
            {MODE_OPTIONS.map(({ value, label, icon: Icon, color }) => {
              const isActive = mode.includes(value);
              return (
                <Button
                  key={value}
                  type="button"
                  variant="outline"
                  onClick={() => toggleMode(value)}
                  className={cn(
                    'w-full h-9 transition-all duration-300 border flex items-center justify-center',
                    isActive
                      ? 'bg-mode-accent-soft border-mode-accent text-mode-accent shadow-sm'
                      : 'border-border text-muted-foreground hover:text-foreground hover:border-muted-foreground',
                  )}
                  title={label}
                >
                  <Icon className={cn('h-4 w-4', isActive && color)} />
                </Button>
              );
            })}
          </div>
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
