import { ContextModeToggle } from '@/features/context-mode/ContextModeToggle';
import { WeatherSelector } from '@/features/context-mode/WeatherSelector';
import { DailyChecklist } from '@/features/checklist/DailyChecklist';
import { JournalSection } from '@/features/journal/JournalSection';
import { Calendar } from '@/features/calendar/Calendar';
import { ProgressTracker } from '@/features/progress/ProgressTracker';
import { useAppStore } from '@/store/useAppStore';

export function MainPage() {
  const { tasks, habits } = useAppStore();

  return (
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
        <section className="flex gap-4">
          <div className="flex-1">
            <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
              Today's Mode
            </h2>
            <ContextModeToggle />
          </div>
          <div>
            <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3">
              Weather
            </h2>
            <WeatherSelector />
          </div>
        </section>

        <section>
          <DailyChecklist tasks={tasks} habits={habits} />
        </section>

        <section>
          <JournalSection />
        </section>
      </main>
    </div>
  );
}
