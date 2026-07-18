import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DailyChecklist } from '../DailyChecklist';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: Object.assign(vi.fn(), {
    getState: vi.fn(),
  }),
}));

vi.mock('../TaskItem', () => ({
  TaskItem: ({ title }: any) => <div data-testid="task-item">{title}</div>,
}));
vi.mock('../HabitItem', () => ({
  HabitItem: ({ title }: any) => <div data-testid="habit-item">{title}</div>,
}));

describe('DailyChecklist', () => {
  const mockTasks = [
    { id: 't1', title: 'Task 1', section: 'dev', type: 'boolean', conditions: { mode: ['Growth'], weather: ['sunny'] } },
    { id: 't2', title: 'Task 2', section: 'japanese', type: 'boolean', conditions: { mode: ['Rest'], weather: ['rainy'] } },
  ];
  const mockHabits = [
    { id: 'h1', title: 'Habit 1', category: 'Health' },
  ];

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      dailyRecord: {
        tasks: [{ taskId: 't1', targetAmount: 1, actualAmount: 0, isCompleted: false }, { taskId: 't2', targetAmount: 1, actualAmount: 0, isCompleted: false }],
        habits: [{ habitId: 'h1', isCompleted: false }],
      },
      isLoading: false,
    });
    vi.mocked(useAppStore.getState as any).mockReturnValue({
      currentMode: 'Growth',
      currentWeather: 'sunny',
    });
  });

  it('renders loading skeleton', () => {
    vi.mocked(useAppStore).mockReturnValue({ isLoading: true });
    render(<DailyChecklist tasks={[]} habits={[]} />);
    // Our skeleton is just 4 divs, but no specific text. We can check for a common class.
    expect(document.querySelector('.space-y-4')).toBeInTheDocument();
  });

  it('renders empty state when no daily record', () => {
    vi.mocked(useAppStore).mockReturnValue({ dailyRecord: null, isLoading: false });
    render(<DailyChecklist tasks={[]} habits={[]} />);
    expect(screen.getByText(/Select a date/i)).toBeInTheDocument();
  });

  it('renders filtered tasks and habits', () => {
    render(<DailyChecklist tasks={mockTasks as any} habits={mockHabits as any} />);
    
    // Only t1 should be rendered because of Growth and sunny conditions
    expect(screen.getByTestId('task-item')).toHaveTextContent('Task 1');
    expect(screen.queryByText('Task 2')).not.toBeInTheDocument();

    expect(screen.getByTestId('habit-item')).toHaveTextContent('Habit 1');
  });
});
