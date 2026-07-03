import { useEffect, useState } from 'react';
import { Leaf } from 'lucide-react';
import { Calendar } from '@/features/calendar/Calendar';
import { ContextModeToggle } from '@/features/context-mode/ContextModeToggle';
import { DailyChecklist } from '@/features/checklist/DailyChecklist';
import { useAppStore } from '@/store/useAppStore';
import { todayJST } from '@/lib/dayjs';
import { cn } from '@/lib/utils';
import type { Task, Habit } from '@/types';

// Placeholder master data until we have CRUD endpoints.
// These would normally come from a GET /api/v1/tasks?status=active call.
const SAMPLE_TASKS: Task[] = [
  {
    id: 'task_001',
    section: 'japanese',
    title: 'Memorize Kanji',
    type: 'quantitative',
    metrics: { dailyTarget: 10, totalTarget: 500 },
    conditions: { weather: 'any', mode: 'any' },
    status: 'active',
  },
  {
    id: 'task_002',
    section: 'japanese',
    title: 'Read NHK News',
    type: 'boolean',
    metrics: { dailyTarget: 1, totalTarget: 0 },
    conditions: { weather: 'any', mode: 'any' },
    status: 'active',
  },
  {
    id: 'task_003',
    section: 'dev',
    title: 'LeetCode Problems',
    type: 'quantitative',
    metrics: { dailyTarget: 3, totalTarget: 100 },
    conditions: { weather: 'any', mode: 'any' },
    status: 'active',
  },
  {
    id: 'task_004',
    section: 'self_dev',
    title: 'Read Book (pages)',
    type: 'quantitative',
    metrics: { dailyTarget: 20, totalTarget: 300 },
    conditions: { weather: 'any', mode: 'any' },
    status: 'active',
  },
];

const SAMPLE_HABITS: Habit[] = [
  { id: 'habit_001', title: 'Meditate instead of using smartphone', category: 'mindfulness', status: 'active' },
  { id: 'habit_002', title: 'Morning stretching routine', category: 'health', status: 'active' },
  { id: 'habit_003', title: 'Write gratitude note', category: 'mindfulness', status: 'active' },
];

const MODE_CLASS_MAP = {
  Growth: 'mode-growth',
  Rest: 'mode-rest',
  Work: 'mode-work',
} as const;

function App() {
  const { currentMode, selectedDate, setDateAndFetch, error } = useAppStore();
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    if (!initialized) {
      setDateAndFetch(todayJST());
      setInitialized(true);
    }
  }, [initialized, setDateAndFetch]);

  const modeClass = MODE_CLASS_MAP[currentMode];
  const formattedDate = selectedDate
    ? new Date(selectedDate + 'T00:00:00+09:00').toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : '';

  return (
    <div className={cn('min-h-screen bg-background transition-colors duration-500', modeClass)}>
      <div className="mx-auto max-w-6xl px-4 py-6 lg:px-8">
        {/* Header */}
        <header className="mb-8">
          <div className="flex items-center gap-3 mb-1">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-mode-accent-soft">
              <Leaf className="h-5 w-5 text-mode-accent" />
            </div>
            <h1 className="text-xl font-bold tracking-tight">Daily Seed</h1>
          </div>
          {selectedDate && (
            <p className="text-sm text-muted-foreground ml-12">{formattedDate}</p>
          )}
        </header>

        {/* Error banner */}
        {error && (
          <div className="mb-6 rounded-xl border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
            {error}
          </div>
        )}

        {/* Main layout */}
        <div className="grid gap-8 lg:grid-cols-[280px_1fr]">
          {/* Sidebar */}
          <aside className="space-y-6">
            <div className="rounded-2xl border border-border bg-card p-4 shadow-sm">
              <Calendar />
            </div>
          </aside>

          {/* Content */}
          <main className="space-y-6">
            <section>
              <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
                Today's Mode
              </h2>
              <ContextModeToggle />
            </section>

            <section>
              <DailyChecklist tasks={SAMPLE_TASKS} habits={SAMPLE_HABITS} />
            </section>
          </main>
        </div>
      </div>
    </div>
  );
}

export default App;
