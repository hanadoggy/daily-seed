import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import App from '../App';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/api/client', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/api/client')>();
  return {
    ...actual,
    fetchDailyRecord: vi.fn().mockImplementation((date: string) =>
      Promise.resolve({
        id: 'rec1',
        date: date,
        context: { mode: 'Growth', weather: 'sunny' },
        tasks: [],
        habits: [],
        journal: { oneLineReview: '', threeLineDiary: '' },
      }),
    ),
    fetchExistingRecordDates: vi.fn().mockResolvedValue({ dates: ['2026-07-25'] }),
    fetchTasks: vi.fn().mockResolvedValue([]),
    fetchHabits: vi.fn().mockResolvedValue([]),
    fetchTaskProgress: vi.fn().mockResolvedValue([]),
    fetchHeatmap: vi.fn().mockResolvedValue({ days: [{ date: '2026-01-01', total: 0, habits: 0, sectionCounts: {} }] }),
    fetchSummary: vi.fn().mockResolvedValue({
      period: 'weekly',
      startDate: '2026-07-19',
      endDate: '2026-07-25',
      totalDays: 7,
      recordedDays: 1,
      taskCompletion: { overall: 0, sections: {}, perTask: [] },
      habitCompletion: { overall: 0, perHabit: [] },
      modeDistribution: {},
      journals: [],
    }),
    fetchStreaks: vi.fn().mockResolvedValue({ habits: [] }),
  };
});

describe('App component & Header mode button', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAppStore.setState({
      selectedDate: '2026-07-25',
      isAdminMode: false,
      currentMode: 'Growth',
      currentWeather: 'sunny',
      dailyRecord: null,
      tasks: [],
      habits: [],
      taskProgress: [],
      existingRecordDates: ['2026-07-25'],
      isLoading: false,
      error: null,
    });
  });

  it('renders Edit Mode button in header when selectedDate is today', async () => {
    await act(async () => {
      render(
        <MemoryRouter initialEntries={['/']}>
          <App />
        </MemoryRouter>,
      );
    });

    const editModeBtn = await screen.findByRole('button', { name: /edit mode/i });
    expect(editModeBtn).toBeInTheDocument();
  });

  it('renders View Mode button when selectedDate is past and isAdminMode is false', async () => {
    await act(async () => {
      render(
        <MemoryRouter initialEntries={['/']}>
          <App />
        </MemoryRouter>,
      );
    });

    // Simulate user selecting a past date after App initialization
    await act(async () => {
      useAppStore.setState({ selectedDate: '2026-01-01', isAdminMode: false });
    });

    const viewModeBtn = await screen.findByRole('button', { name: /view mode/i });
    expect(viewModeBtn).toBeInTheDocument();

    // Click View Mode button to toggle to Admin Mode
    await act(async () => {
      fireEvent.click(viewModeBtn);
    });

    const adminModeBtn = await screen.findByRole('button', { name: /admin mode/i });
    expect(adminModeBtn).toBeInTheDocument();
  });

  it('does not render Mode button in header when on /analytics route', async () => {
    await act(async () => {
      render(
        <MemoryRouter initialEntries={['/analytics']}>
          <App />
        </MemoryRouter>,
      );
    });

    expect(screen.queryByRole('button', { name: /edit mode/i })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /view mode/i })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /admin mode/i })).not.toBeInTheDocument();
  });
});

