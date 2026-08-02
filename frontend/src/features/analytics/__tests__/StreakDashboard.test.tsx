import { describe, it, expect } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { StreakDashboard } from '../StreakDashboard';
import type { StreakResponse } from '@/types';

const mockStreakData: StreakResponse = {
  habits: [
    {
      habitId: 'h1',
      title: 'Morning Meditation',
      category: 'Mindfulness',
      currentStreak: 7,
      longestStreak: 14,
      totalDays: 20,
      lastCompleted: '2026-07-25',
      milestones: [7],
    },
    {
      habitId: 'h2',
      title: 'Workout',
      category: 'Health',
      currentStreak: 30,
      longestStreak: 30,
      totalDays: 45,
      lastCompleted: '2026-07-25',
      milestones: [7, 30],
    },
  ],
};

describe('StreakDashboard component', () => {
  it('renders loading skeleton when isLoading is true', () => {
    const { container } = render(<StreakDashboard data={null} isLoading={true} />);
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument();
  });

  it('renders empty message when data has no active habits', () => {
    render(<StreakDashboard data={{ habits: [] }} isLoading={false} />);
    expect(screen.getByText(/No active habits tracked for streaks yet/i)).toBeInTheDocument();
  });

  it('renders habit streak cards correctly when data is provided', () => {
    render(<StreakDashboard data={mockStreakData} isLoading={false} />);

    expect(screen.getByText('Habit Streaks & Statistics')).toBeInTheDocument();
    expect(screen.getByText('Morning Meditation')).toBeInTheDocument();
    expect(screen.getByText('Workout')).toBeInTheDocument();

    // Streaks
    expect(screen.getByText('7d')).toBeInTheDocument();
    expect(screen.getByText('30d')).toBeInTheDocument();
    expect(screen.getAllByText('7d Streak').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText('30d Streak')).toBeInTheDocument();
  });

  it('opens celebration modal when clicking a milestone badge', () => {
    render(<StreakDashboard data={mockStreakData} isLoading={false} />);

    const milestoneBadge = screen.getByText('30d Streak');
    fireEvent.click(milestoneBadge);

    expect(screen.getByText('30-Day Streak!')).toBeInTheDocument();
    expect(screen.getByText('Milestone Achieved!')).toBeInTheDocument();
  });
});
