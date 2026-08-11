const VALID_SECTIONS = ['japanese', 'dev', 'self_dev', 'exercise'] as const;
const VALID_TYPES = ['quantitative', 'boolean'] as const;
const VALID_WEATHERS = ['sunny', 'rainy'] as const;
const VALID_MODES = ['Growth', 'Rest', 'Office', 'Remote'] as const;
const DATE_REGEX = /^\d{4}-\d{2}-\d{2}$/;

export interface ValidationResult {
  valid: boolean;
  error?: string;
}

function fail(message: string): ValidationResult {
  return { valid: false, error: message };
}

const ok: ValidationResult = { valid: true };

export interface TaskFormPayload {
  title: string;
  section: string;
  type: string;
  unit: string;
  metrics: { dailyTarget: number; totalTarget: number };
  conditions: { weather: string[]; mode: string[] };
  startDate: string;
}

export function validateTaskForm(data: TaskFormPayload): ValidationResult {
  if (!data.title.trim()) return fail('Title is required');
  if (data.title.length > 200) return fail('Title must not exceed 200 characters');

  if (!VALID_SECTIONS.includes(data.section as (typeof VALID_SECTIONS)[number])) {
    return fail('Invalid section');
  }
  if (!VALID_TYPES.includes(data.type as (typeof VALID_TYPES)[number])) {
    return fail('Invalid type');
  }

  if (data.type === 'quantitative') {
    if (!data.unit.trim()) return fail('Unit is required');
    if (data.unit.length > 50) return fail('Unit must not exceed 50 characters');
    if (data.metrics.dailyTarget < 1) return fail('Daily target must be at least 1');
  }
  if (data.metrics.totalTarget < 0) return fail('Total target cannot be negative');

  if (data.conditions.weather.length === 0) return fail('At least one weather condition is required');
  if (data.conditions.mode.length === 0) return fail('At least one mode is required');

  for (const w of data.conditions.weather) {
    if (!VALID_WEATHERS.includes(w as (typeof VALID_WEATHERS)[number])) {
      return fail(`Invalid weather value: ${w}`);
    }
  }
  for (const m of data.conditions.mode) {
    if (!VALID_MODES.includes(m as (typeof VALID_MODES)[number])) {
      return fail(`Invalid mode value: ${m}`);
    }
  }

  if (!DATE_REGEX.test(data.startDate)) return fail('Start date must be in YYYY-MM-DD format');

  return ok;
}

export interface HabitFormPayload {
  title: string;
  category: string;
}

export function validateHabitForm(data: HabitFormPayload): ValidationResult {
  if (!data.title.trim()) return fail('Title is required');
  if (data.title.length > 200) return fail('Title must not exceed 200 characters');
  if (!data.category.trim()) return fail('Category is required');
  return ok;
}
