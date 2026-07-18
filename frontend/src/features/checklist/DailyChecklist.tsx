import { ListChecks, Sparkles } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { useAppStore } from '@/store/useAppStore';
import { TaskItem } from './TaskItem';
import { HabitItem } from './HabitItem';
import type { Task, Habit } from '@/types';
import { HABIT_CATEGORIES } from '../manage/categoryOptions';

import { SECTION_LABELS } from '@/lib/constants';

interface DailyChecklistProps {
  tasks: Task[];
  habits: Habit[];
}

export function DailyChecklist({ tasks, habits }: DailyChecklistProps) {
  const { dailyRecord, isLoading } = useAppStore();

  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-16 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  if (!dailyRecord) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-muted-foreground">
        <ListChecks className="h-10 w-10 mb-3 opacity-40" />
        <p className="text-sm">Select a date to view your checklist</p>
      </div>
    );
  }

  // Group task entries by section using the master Task data
  const taskMap = new Map(tasks.map((t) => [t.id, t]));
  const habitMap = new Map(habits.map((h) => [h.id, h]));

  const tasksBySection = new Map<string, typeof dailyRecord.tasks>();
  for (const entry of dailyRecord.tasks) {
    const task = taskMap.get(entry.taskId);
    if (!task) continue;

    const { currentMode, currentWeather } = useAppStore.getState();

    if (!task.conditions.mode.includes(currentMode)) continue;
    if (!task.conditions.weather.includes(currentWeather)) continue;

    const section = task.section ?? 'other';
    const existing = tasksBySection.get(section) ?? [];
    existing.push(entry);
    tasksBySection.set(section, existing);
  }

  const sectionOrder = ['japanese', 'dev', 'self_dev', 'exercise'];

  return (
    <div className="space-y-6">
      {/* Tasks by section */}
      {sectionOrder.map((section) => {
        const entries = tasksBySection.get(section);
        if (!entries || entries.length === 0) return null;

        return (
          <div key={section}>
            <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3 flex items-center gap-2">
              {SECTION_LABELS[section] ?? section}
            </h3>
            <div className="space-y-2">
              {entries.map((entry) => {
                const task = taskMap.get(entry.taskId);
                return (
                  <TaskItem
                    key={entry.taskId}
                    entry={entry}
                    title={task?.title ?? 'Unknown Task'}
                    type={task?.type ?? 'boolean'}
                  />
                );
              })}
            </div>
          </div>
        );
      })}

      {/* Habits */}
      {dailyRecord.habits.length > 0 && (
        <div>
          <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-3 flex items-center gap-2">
            <Sparkles className="h-3.5 w-3.5" />
            Habits
          </h3>
          <div className="space-y-2">
            {[...dailyRecord.habits]
              .sort((a, b) => {
                const habitA = habitMap.get(a.habitId);
                const habitB = habitMap.get(b.habitId);
                if (!habitA || !habitB) return 0;
                
                const catOrderA = HABIT_CATEGORIES.findIndex((c) => c.value === habitA.category);
                const catOrderB = HABIT_CATEGORIES.findIndex((c) => c.value === habitB.category);
                
                if (catOrderA !== catOrderB) {
                  if (catOrderA === -1) return 1;
                  if (catOrderB === -1) return -1;
                  return catOrderA - catOrderB;
                }
                return habitA.title.localeCompare(habitB.title);
              })
              .map((entry) => {
              const habit = habitMap.get(entry.habitId);
              return (
                <HabitItem
                  key={entry.habitId}
                  entry={entry}
                  title={habit?.title ?? 'Unknown Habit'}
                  category={habit?.category}
                />
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
