import type { StateCreator } from 'zustand';
import type { AppState, TaskSlice } from '../types';
import {
  createTask as apiCreateTask,
  updateTask as apiUpdateTask,
  deleteTask as apiDeleteTask,
  fetchTaskProgress as apiFetchProgress,
  migrateTask as apiMigrateTask,
} from '../../api/client';

export const createTaskSlice: StateCreator<AppState, [], [], TaskSlice> = (set, get) => ({
  tasks: [],
  taskProgress: [],

  addTask: async (task) => {
    try {
      const created = await apiCreateTask(task);
      set((state) => ({ tasks: [...state.tasks, created] }));
      get().setDateAndFetch(get().selectedDate);
      get().fetchProgress();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create task';
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
      set({ tasks: previous });
      const message = err instanceof Error ? err.message : 'Failed to update task';
      set({ error: message });
    }
  },

  archiveTask: async (id) => {
    const previous = get().tasks;
    set((state) => ({ tasks: state.tasks.filter((t) => t.id !== id) }));
    try {
      await apiDeleteTask(id);
      get().setDateAndFetch(get().selectedDate);
      get().fetchProgress();
    } catch {
      set({ tasks: previous });
    }
  },

  migrateTask: async (id) => {
    try {
      const result = await apiMigrateTask(id);
      set((state) => ({
        tasks: state.tasks
          .map((t) => (t.id === result.archivedTask.id ? result.archivedTask : t))
          .concat(result.newTask),
      }));
      // Refresh progress after migration.
      await get().fetchProgress();
      await get().setDateAndFetch(get().selectedDate);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to migrate task';
      set({ error: message });
    }
  },

  fetchProgress: async () => {
    try {
      const progress = await apiFetchProgress();
      set({ taskProgress: progress });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch progress';
      set({ error: message });
    }
  },
});
