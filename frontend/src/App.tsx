import { useEffect, useState } from 'react';
import { Leaf, Settings } from 'lucide-react';
import { Calendar } from '@/features/calendar/Calendar';
import { ContextModeToggle } from '@/features/context-mode/ContextModeToggle';
import { DailyChecklist } from '@/features/checklist/DailyChecklist';
import { JournalSection } from '@/features/journal/JournalSection';
import { ProgressTracker } from '@/features/progress/ProgressTracker';
import { ManagePanel } from '@/features/manage/ManagePanel';
import { useAppStore } from '@/store/useAppStore';
import { todayJST } from '@/lib/dayjs';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { ThemeSwitcher } from '@/components/theme-switcher';

const MODE_CLASS_MAP = {
  Growth: 'mode-growth',
  Rest: 'mode-rest',
  Work: 'mode-work',
} as const;

function App() {
  const { currentMode, selectedDate, setDateAndFetch, fetchMasterData, tasks, habits, error } =
    useAppStore();
  const [initialized, setInitialized] = useState(false);
  const [showManage, setShowManage] = useState(false);

  useEffect(() => {
    if (!initialized) {
      setDateAndFetch(todayJST());
      fetchMasterData();
      setInitialized(true);
    }
  }, [initialized, setDateAndFetch, fetchMasterData]);

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
            <div className="ml-auto flex items-center gap-2">
              <ThemeSwitcher />
              <Button
                id="manage-button"
                variant="outline"
                size="sm"
                className="gap-1.5"
                onClick={() => setShowManage(true)}
              >
                <Settings className="h-3.5 w-3.5" />
                Manage
              </Button>
            </div>
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
            <ProgressTracker />
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
              <DailyChecklist tasks={tasks} habits={habits} />
            </section>

            <section>
              <JournalSection />
            </section>
          </main>
        </div>
      </div>

      {/* Management Panel */}
      {showManage && <ManagePanel onClose={() => setShowManage(false)} />}
    </div>
  );
}

export default App;
