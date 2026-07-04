import { useState } from 'react';
import { X, Plus, Pencil, Trash2, ListTodo, Sparkles } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/store/useAppStore';
import { TaskForm } from './TaskForm';
import { HabitForm } from './HabitForm';
import type { Task, Habit } from '@/types';

type Tab = 'tasks' | 'habits';
type FormState =
  | { mode: 'closed' }
  | { mode: 'create-task' }
  | { mode: 'edit-task'; task: Task }
  | { mode: 'create-habit' }
  | { mode: 'edit-habit'; habit: Habit };

const SECTION_LABELS: Record<string, string> = {
  japanese: '🇯🇵 Japanese',
  dev: '💻 Dev',
  self_dev: '📚 Self Dev',
};

interface ManagePanelProps {
  onClose: () => void;
}

export function ManagePanel({ onClose }: ManagePanelProps) {
  const { tasks, habits, archiveTask, archiveHabit } = useAppStore();
  const [activeTab, setActiveTab] = useState<Tab>('tasks');
  const [formState, setFormState] = useState<FormState>({ mode: 'closed' });
  const [confirmingTaskId, setConfirmingTaskId] = useState<string | null>(null);
  const [confirmingHabitId, setConfirmingHabitId] = useState<string | null>(null);

  const activeTasks = tasks.filter((t) => t.status === 'active');
  const activeHabits = habits.filter((h) => h.status === 'active');

  const showingForm = formState.mode !== 'closed';

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />

      {/* Panel */}
      <div className="relative z-10 flex h-full w-full max-w-md flex-col border-l border-border bg-card shadow-2xl animate-in slide-in-from-right duration-300">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-base font-semibold">Manage</h2>
          <Button variant="ghost" size="icon-sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-border">
          <button
            onClick={() => {
              setActiveTab('tasks');
              setFormState({ mode: 'closed' });
              setConfirmingTaskId(null);
            }}
            className={`flex flex-1 items-center justify-center gap-1.5 px-4 py-2.5 text-sm font-medium transition-colors ${
              activeTab === 'tasks'
                ? 'border-b-2 border-mode-accent text-foreground'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <ListTodo className="h-3.5 w-3.5" />
            Tasks ({activeTasks.length})
          </button>
          <button
            onClick={() => {
              setActiveTab('habits');
              setFormState({ mode: 'closed' });
              setConfirmingHabitId(null);
            }}
            className={`flex flex-1 items-center justify-center gap-1.5 px-4 py-2.5 text-sm font-medium transition-colors ${
              activeTab === 'habits'
                ? 'border-b-2 border-mode-accent text-foreground'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <Sparkles className="h-3.5 w-3.5" />
            Habits ({activeHabits.length})
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto px-5 py-4">
          {showingForm ? (
            <div className="rounded-xl border border-border bg-background p-4">
              {(formState.mode === 'create-task' || formState.mode === 'edit-task') && (
                <TaskForm
                  task={formState.mode === 'edit-task' ? formState.task : undefined}
                  onClose={() => setFormState({ mode: 'closed' })}
                />
              )}
              {(formState.mode === 'create-habit' || formState.mode === 'edit-habit') && (
                <HabitForm
                  habit={formState.mode === 'edit-habit' ? formState.habit : undefined}
                  onClose={() => setFormState({ mode: 'closed' })}
                />
              )}
            </div>
          ) : (
            <>
              {/* Add button */}
              <Button
                variant="outline"
                className="mb-4 w-full gap-1.5"
                onClick={() =>
                  setFormState(
                    activeTab === 'tasks' ? { mode: 'create-task' } : { mode: 'create-habit' },
                  )
                }
              >
                <Plus className="h-3.5 w-3.5" />
                Add {activeTab === 'tasks' ? 'Task' : 'Habit'}
              </Button>

              {/* List */}
              {activeTab === 'tasks' && (
                <div className="space-y-2">
                  {activeTasks.length === 0 && (
                    <p className="py-8 text-center text-sm text-muted-foreground">
                      No active tasks yet. Create one to get started.
                    </p>
                  )}
                  {activeTasks.map((task) => (
                    <div
                      key={task.id}
                      className="group flex items-center gap-3 rounded-xl border border-border bg-background p-3 transition-colors hover:border-mode-accent/30"
                    >
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-medium">{task.title}</p>
                        <div className="mt-0.5 flex items-center gap-2 text-xs text-muted-foreground">
                          <span>{SECTION_LABELS[task.section] ?? task.section}</span>
                          <span>·</span>
                          <span>
                            {task.type === 'boolean'
                              ? 'Yes/No'
                              : `${task.metrics.dailyTarget}/day`}
                          </span>
                        </div>
                      </div>
                      <div className={`flex gap-1 transition-opacity ${confirmingTaskId === task.id ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'}`}>
                        {confirmingTaskId === task.id ? (
                          <div className="flex items-center gap-1 rounded-md bg-muted/50 px-1">
                            <span className="mr-1 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">Sure?</span>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              className="h-6 w-6 text-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 dark:hover:bg-emerald-950"
                              onClick={() => {
                                archiveTask(task.id);
                                setConfirmingTaskId(null);
                              }}
                            >
                              ✅
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              className="h-6 w-6 text-red-500 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950"
                              onClick={() => setConfirmingTaskId(null)}
                            >
                              ❌
                            </Button>
                          </div>
                        ) : (
                          <>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              onClick={() => setFormState({ mode: 'edit-task', task })}
                            >
                              <Pencil className="h-3 w-3" />
                            </Button>
                            <Button
                              variant="destructive"
                              size="icon-xs"
                              onClick={() => setConfirmingTaskId(task.id)}
                            >
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {activeTab === 'habits' && (
                <div className="space-y-2">
                  {activeHabits.length === 0 && (
                    <p className="py-8 text-center text-sm text-muted-foreground">
                      No active habits yet. Create one to get started.
                    </p>
                  )}
                  {activeHabits.map((habit) => (
                    <div
                      key={habit.id}
                      className="group flex items-center gap-3 rounded-xl border border-border bg-background p-3 transition-colors hover:border-mode-accent/30"
                    >
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-medium">{habit.title}</p>
                        <p className="mt-0.5 text-xs text-muted-foreground capitalize">
                          {habit.category}
                        </p>
                      </div>
                      <div className={`flex gap-1 transition-opacity ${confirmingHabitId === habit.id ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'}`}>
                        {confirmingHabitId === habit.id ? (
                          <div className="flex items-center gap-1 rounded-md bg-muted/50 px-1">
                            <span className="mr-1 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">Sure?</span>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              className="h-6 w-6 text-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 dark:hover:bg-emerald-950"
                              onClick={() => {
                                archiveHabit(habit.id);
                                setConfirmingHabitId(null);
                              }}
                            >
                              ✅
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              className="h-6 w-6 text-red-500 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950"
                              onClick={() => setConfirmingHabitId(null)}
                            >
                              ❌
                            </Button>
                          </div>
                        ) : (
                          <>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              onClick={() => setFormState({ mode: 'edit-habit', habit })}
                            >
                              <Pencil className="h-3 w-3" />
                            </Button>
                            <Button
                              variant="destructive"
                              size="icon-xs"
                              onClick={() => setConfirmingHabitId(habit.id)}
                            >
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
