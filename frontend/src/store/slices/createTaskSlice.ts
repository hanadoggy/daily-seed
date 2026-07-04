import type { StateCreator } from 'zustand';
import type { AppState, TaskSlice } from '../types';
import {
  createTask as apiCreateTask,
  updateTask as apiUpdateTask,
  deleteTask as apiDeleteTask,
} from '../../api/client';

export const createTaskSlice: StateCreator<AppState, [], [], TaskSlice> = (set, get) => ({
  tasks: [],

  addTask: async (task) => {
    try {
      const created = await apiCreateTask(task);
      set((state) => ({ tasks: [...state.tasks, created] }));
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
    } catch {
      set({ tasks: previous });
    }
  },
});
