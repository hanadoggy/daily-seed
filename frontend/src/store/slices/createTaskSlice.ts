import type { StateCreator } from 'zustand';
import { toast } from 'sonner';
import type { AppState, TaskSlice } from '../types';
import {
  createTask as apiCreateTask,
  updateTask as apiUpdateTask,
  deleteTask as apiDeleteTask,
  fetchTaskProgress as apiFetchProgress,
  migrateTask as apiMigrateTask,
  ApiError,
} from '../../api/client';
import { getErrorMessage } from '../../lib/errorMessages';

export const createTaskSlice: StateCreator<AppState, [], [], TaskSlice> = (set, get) => ({
  tasks: [],
  taskProgress: [],
  migratingTaskIds: new Set<string>(),

  addTask: async (task) => {
    try {
      const created = await apiCreateTask(task);
      set((state) => ({ tasks: [...state.tasks, created] }));
      get().setDateAndFetch(get().selectedDate);
      get().fetchProgress();
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('createTask', status);
      toast.error(message);
      set({ error: message });
    }
  },

  editTask: async (id, task) => {
    const previous = get().tasks;
    try {
      const updated = await apiUpdateTask(id, task);
      set((state) => ({
        tasks: state.tasks.map((t) => (t.id === id ? updated : t)),
      }));
      get().setDateAndFetch(get().selectedDate);
      get().fetchProgress();
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('updateTask', status);
      toast.error(message);
      set({ tasks: previous, error: message });
    }
  },

  archiveTask: async (id) => {
    const previous = get().tasks;
    const today = new Date().toISOString().split('T')[0];
    set((state) => ({
      tasks: state.tasks.map((t) => (t.id === id ? { ...t, status: 'archived', endDate: today } : t)),
    }));
    try {
      await apiDeleteTask(id);
      get().setDateAndFetch(get().selectedDate);
      get().fetchProgress();
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('archiveTask', status);
      toast.error(message);
      set({ tasks: previous, error: message });
    }
  },

  migrateTask: async (id) => {
    if (get().migratingTaskIds.has(id)) return;
    set((state) => ({
      migratingTaskIds: new Set([...state.migratingTaskIds, id]),
    }));
    try {
      const completionDate = get().selectedDate;
      const result = await apiMigrateTask(id, completionDate);
      set((state) => ({
        tasks: state.tasks
          .map((t) => (t.id === result.archivedTask.id ? result.archivedTask : t))
          .concat(result.newTask),
      }));
      // Refresh progress after migration.
      await get().fetchProgress();
      await get().setDateAndFetch(get().selectedDate);
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('migrateTask', status);
      toast.error(message);
      set({ error: message });
    } finally {
      set((state) => {
        const next = new Set(state.migratingTaskIds);
        next.delete(id);
        return { migratingTaskIds: next };
      });
    }
  },

  fetchProgress: async () => {
    try {
      const progress = await apiFetchProgress();
      set({ taskProgress: progress });
    } catch (err) {
      const status = err instanceof ApiError ? err.status : 0;
      const message = getErrorMessage('fetchProgress', status);
      toast.error(message);
      set({ error: message });
    }
  },
});
