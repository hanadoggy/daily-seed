import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import * as client from '../client';
import type { Task } from '../../types';

describe('API Client', () => {
  const mockFetch = vi.fn();
  
  beforeEach(() => {
    vi.stubGlobal('fetch', mockFetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('request error handling', () => {
    it('throws ApiError on non-ok response with json body', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: () => Promise.resolve({ code: 'BAD_REQUEST', message: 'Invalid field' }),
      });

      try {
        await client.fetchDailyRecord('2023-10-10');
        expect.fail('Should have thrown an error');
      } catch (error: any) {
        expect(error.name).toBe('ApiError');
        expect(error.status).toBe(400);
        expect(error.code).toBe('BAD_REQUEST');
        expect(error.message).toBe('Invalid field');
      }
    });

    it('throws ApiError on non-ok response with no json body', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.reject(new Error('no body')),
      });

      try {
        await client.fetchDailyRecord('2023-10-10');
        expect.fail('Should have thrown an error');
      } catch (error: any) {
        expect(error.status).toBe(500);
        expect(error.code).toBe('UNKNOWN');
        expect(error.message).toBe('Request failed with status 500');
      }
    });
  });

  describe('fetchDailyRecord', () => {
    it('makes correct GET request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ id: '1' }) });
      const res = await client.fetchDailyRecord('2023-10-10');
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/daily/2023-10-10'),
        expect.objectContaining({ headers: { 'Content-Type': 'application/json' } })
      );
      expect(res).toEqual({ id: '1' });
    });
  });

  describe('patchDailyRecord', () => {
    it('makes correct PATCH request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      await client.patchDailyRecord('2023-10-10', { context: { mode: 'Growth', weather: 'sunny' } });
      
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/daily/2023-10-10'),
        expect.objectContaining({
          method: 'PATCH',
          body: JSON.stringify({ context: { mode: 'Growth', weather: 'sunny' } }),
        })
      );
    });
  });

  describe('fetchTasks', () => {
    it('makes correct GET request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) });
      await client.fetchTasks();
      expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('/tasks'), expect.anything());
    });
  });

  describe('createTask', () => {
    it('makes correct POST request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      const task: Omit<Task, 'id' | 'status'> = { title: 'Test', section: 'dev', type: 'boolean', metrics: { dailyTarget: 1, totalTarget: 0 }, conditions: { weather: ['sunny'], mode: ['Growth'] }, startDate: '2023-10-10' };
      await client.createTask(task);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/tasks'),
        expect.objectContaining({ method: 'POST', body: JSON.stringify(task) })
      );
    });
  });

  describe('updateTask', () => {
    it('makes correct PUT request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      const task: Omit<Task, 'id' | 'status'> = { title: 'Test', section: 'dev', type: 'boolean', metrics: { dailyTarget: 1, totalTarget: 0 }, conditions: { weather: ['sunny'], mode: ['Growth'] }, startDate: '2023-10-10' };
      await client.updateTask('id1', task);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/tasks/id1'),
        expect.objectContaining({ method: 'PUT', body: JSON.stringify(task) })
      );
    });
  });

  describe('deleteTask', () => {
    it('makes correct DELETE request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      await client.deleteTask('id1');
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/tasks/id1'),
        expect.objectContaining({ method: 'DELETE' })
      );
    });
  });

  describe('fetchTaskProgress', () => {
    it('makes correct GET request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) });
      await client.fetchTaskProgress();
      expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('/tasks/progress'), expect.anything());
    });
  });

  describe('migrateTask', () => {
    it('makes correct POST request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      await client.migrateTask('id1', '2023-10-10');
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/tasks/id1/migrate'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ completionDate: '2023-10-10' })
        })
      );
    });
  });

  describe('fetchHabits', () => {
    it('makes correct GET request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) });
      await client.fetchHabits();
      expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('/habits'), expect.anything());
    });
  });

  describe('createHabit', () => {
    it('makes correct POST request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      const habit = { title: 'Test', category: 'Health' };
      await client.createHabit(habit);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/habits'),
        expect.objectContaining({ method: 'POST', body: JSON.stringify(habit) })
      );
    });
  });

  describe('updateHabit', () => {
    it('makes correct PUT request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      const habit = { title: 'Test', category: 'Health' };
      await client.updateHabit('id1', habit);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/habits/id1'),
        expect.objectContaining({ method: 'PUT', body: JSON.stringify(habit) })
      );
    });
  });

  describe('deleteHabit', () => {
    it('makes correct DELETE request', async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });
      await client.deleteHabit('id1');
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/habits/id1'),
        expect.objectContaining({ method: 'DELETE' })
      );
    });
  });
});
