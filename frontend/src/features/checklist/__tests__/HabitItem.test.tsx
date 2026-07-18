import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { HabitItem } from '../HabitItem';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('HabitItem', () => {
  const mockToggleHabit = vi.fn();

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue(mockToggleHabit);
    vi.clearAllMocks();
  });

  const entry = { habitId: 'h1', isCompleted: false };

  it('renders habit title', () => {
    render(<HabitItem entry={entry} title="Drink Water" category="Health" />);
    expect(screen.getByText('Drink Water')).toBeInTheDocument();
  });

  it('calls toggleHabitOptimistic on click', () => {
    render(<HabitItem entry={entry} title="Drink Water" category="Health" />);
    fireEvent.click(screen.getByRole('button'));
    expect(mockToggleHabit).toHaveBeenCalledWith('h1', true);
  });

  it('renders completed state correctly', () => {
    const completedEntry = { habitId: 'h1', isCompleted: true };
    render(<HabitItem entry={completedEntry} title="Drink Water" category="Health" />);
    const button = screen.getByRole('button');
    expect(button.className).toContain('border-mode-accent');
  });
});
