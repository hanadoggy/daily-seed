import { describe, it, expect } from 'vitest';
import { validateTaskForm, validateHabitForm } from '../validation';

describe('validation', () => {
  describe('validateTaskForm', () => {
    const validTask = {
      title: 'Study Japanese',
      section: 'japanese',
      type: 'quantitative',
      unit: 'pages',
      metrics: { dailyTarget: 10, totalTarget: 100 },
      conditions: { weather: ['sunny'], mode: ['Growth'] },
      startDate: '2026-08-11',
    };

    it('passes for valid task payload', () => {
      const res = validateTaskForm(validTask);
      expect(res.valid).toBe(true);
    });

    it('fails when title is empty', () => {
      const res = validateTaskForm({ ...validTask, title: '   ' });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Title is required');
    });

    it('fails when title exceeds 200 chars', () => {
      const res = validateTaskForm({ ...validTask, title: 'a'.repeat(201) });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Title must not exceed 200 characters');
    });

    it('fails when section is invalid', () => {
      const res = validateTaskForm({ ...validTask, section: 'invalid' });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Invalid section');
    });

    it('fails when quantitative dailyTarget is less than 1', () => {
      const res = validateTaskForm({
        ...validTask,
        metrics: { dailyTarget: 0, totalTarget: 100 },
      });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Daily target must be at least 1');
    });

    it('fails when startDate is invalid format', () => {
      const res = validateTaskForm({ ...validTask, startDate: '2026/08/11' });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Start date must be in YYYY-MM-DD format');
    });
  });

  describe('validateHabitForm', () => {
    const validHabit = {
      title: 'Meditation',
      category: 'mindfulness',
    };

    it('passes for valid habit payload', () => {
      const res = validateHabitForm(validHabit);
      expect(res.valid).toBe(true);
    });

    it('fails when title is empty', () => {
      const res = validateHabitForm({ ...validHabit, title: '' });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Title is required');
    });

    it('fails when category is empty', () => {
      const res = validateHabitForm({ ...validHabit, category: '  ' });
      expect(res.valid).toBe(false);
      expect(res.error).toBe('Category is required');
    });
  });
});
