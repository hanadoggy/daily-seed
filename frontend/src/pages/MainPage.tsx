import { ContextModeToggle } from '@/features/context-mode/ContextModeToggle';
import { WeatherSelector } from '@/features/context-mode/WeatherSelector';
import { DailyChecklist } from '@/features/checklist/DailyChecklist';
import { JournalSection } from '@/features/journal/JournalSection';
import { Calendar } from '@/features/calendar/Calendar';
import { ProgressTracker } from '@/features/progress/ProgressTracker';
import { useAppStore } from '@/store/useAppStore';
import { useIsReadOnly } from '@/hooks/useIsReadOnly';
import { todayJST } from '@/lib/dayjs';
import { Pencil, Eye, Unlock } from 'lucide-react';
import { cn } from '@/lib/utils';

export function MainPage() {
  const { tasks, habits, selectedDate, isAdminMode, toggleAdminMode } = useAppStore();
  const isReadOnly = useIsReadOnly();
  const isToday = selectedDate === todayJST();

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
            <div className="flex items-center justify-between mb-3">
              <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Today's Mode
              </h2>
              {!isToday && (
                <button
                  onClick={toggleAdminMode}
                  className={cn(
                    "flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-widest transition-colors",
                    isAdminMode
                      ? "bg-destructive/10 text-destructive hover:bg-destructive/20"
                      : "bg-accent/10 text-accent hover:bg-accent/20 cursor-pointer"
                  )}
                >
                  {isAdminMode ? (
                    <>
                      <Unlock className="h-3 w-3" />
                      Admin Mode
                    </>
                  ) : (
                    <>
                      <Eye className="h-3 w-3" />
                      View Mode
                    </>
                  )}
                </button>
              )}
              {isToday && (
                <div className="flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-widest bg-primary/10 text-primary">
                  <Pencil className="h-3 w-3" />
                  Edit Mode
                </div>
              )}
            </div>
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
