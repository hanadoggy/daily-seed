import { useState, useEffect } from 'react';
import { useAppStore } from '@/store/useAppStore';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';

export function AutoMigrationPrompt() {
  const { taskProgress, migrateTask, migratingTaskIds } = useAppStore();
  const [isOpen, setIsOpen] = useState(false);
  const [targetTask, setTargetTask] = useState<{ id: string; title: string } | null>(null);
  const [dismissedTaskIds, setDismissedTaskIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    // Find first task at 100% that hasn't been dismissed
    const readyTask = taskProgress.find(
      (tp) => tp.percentage >= 100 && !dismissedTaskIds.has(tp.taskId)
    );

    if (readyTask && !isOpen) {
      setTargetTask({ id: readyTask.taskId, title: readyTask.title });
      setIsOpen(true);
    }
  }, [taskProgress, dismissedTaskIds, isOpen]);

  const isMigrating = targetTask ? (migratingTaskIds?.has(targetTask.id) ?? false) : false;

  const handleConfirm = async () => {
    if (targetTask) {
      if (migratingTaskIds?.has(targetTask.id)) return;
      setDismissedTaskIds((prev) => new Set(prev).add(targetTask.id));
      await migrateTask(targetTask.id);
      setIsOpen(false);
      setTargetTask(null);
    }
  };

  const handleCancel = () => {
    if (targetTask) {
      // Add to dismissed so we don't prompt again for this specific task ID in this session
      setDismissedTaskIds((prev) => new Set(prev).add(targetTask.id));
      setIsOpen(false);
      setTargetTask(null);
    }
  };

  return (
    <AlertDialog open={isOpen} onOpenChange={setIsOpen}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Task Goal Reached! 🎉</AlertDialogTitle>
          <AlertDialogDescription>
            You've completed 100% of the lifetime goal for "{targetTask?.title}". 
            Would you like to migrate this to a new active goal?
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={handleCancel}>Not Yet</AlertDialogCancel>
          <AlertDialogAction
            onClick={handleConfirm}
            disabled={isMigrating}
            className="bg-green-500 hover:bg-green-600 disabled:opacity-50"
          >
            {isMigrating ? 'Migrating…' : 'Migrate Task'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
