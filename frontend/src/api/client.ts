import type { DailyRecord } from '../types';

const BASE_URL = import.meta.env.VITE_API_URL ?? '/api/v1';

class ApiError extends Error {
  status: number;
  code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new ApiError(
      res.status,
      body.code ?? 'UNKNOWN',
      body.message ?? `Request failed with status ${res.status}`,
    );
  }

  return res.json();
}

export function fetchDailyRecord(date: string): Promise<DailyRecord> {
  return request<DailyRecord>(`/daily/${date}`);
}

export function patchDailyRecord(
  date: string,
  body: Partial<DailyRecord>,
): Promise<DailyRecord> {
  return request<DailyRecord>(`/daily/${date}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}
