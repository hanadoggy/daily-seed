import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { SummaryDashboard } from '../SummaryDashboard';
import type { SummaryResponse } from '@/types';

const mockSummaryData: SummaryResponse = {
  period: 'weekly',
  startDate: '2026-07-19',
  endDate: '2026-07-25',
  totalDays: 7,
  recordedDays: 3,
  taskCompletion: {
    overall: 80.5,
    sections: { dev: 90, japanese: 70 },
    perTask: [
      {
        taskId: 't1',
        title: 'Study Go',
        section: 'dev',
        type: 'quantitative',
        rate: 90,
        completed: 9,
        target: 10,
      },
    ],
  },
  habitCompletion: {
    overall: 66.7,
    perHabit: [
      {
        habitId: 'h1',
        title: 'Exercise',
        category: 'Health',
        rate: 66.7,
        completed: 2,
        total: 3,
      },
    ],
  },
  modeDistribution: { Growth: 2, Rest: 1 },
  journals: [
    {
      date: '2026-07-20',
      oneLineReview: 'Productive day!',
      threeLineDiary: 'Woke up early.\nCoded 4 hours.\nRead book.',
    },
  ],
};

describe('SummaryDashboard component', () => {
  it('renders loading skeleton when isLoading is true', () => {
    const { container } = render(
      <SummaryDashboard
        data={null}
        period="weekly"
        isLoading={true}
        onPeriodChange={() => {}}
        onNavigate={() => {}}
      />,
    );
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument();
  });

  it('renders data correctly when summary data is provided', () => {
    render(
      <SummaryDashboard
        data={mockSummaryData}
        period="weekly"
        isLoading={false}
        onPeriodChange={() => {}}
        onNavigate={() => {}}
      />,
    );

    expect(screen.getByText('Period Summary')).toBeInTheDocument();
    expect(screen.getByText('2026-07-19 ~ 2026-07-25')).toBeInTheDocument();

    // Overall stats
    expect(screen.getByText('80.5%')).toBeInTheDocument();
    expect(screen.getByText('66.7%')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument(); // recorded days
    expect(screen.getByText(/Growth: 2d/)).toBeInTheDocument();

    // Individual breakdown
    expect(screen.getByText('Study Go')).toBeInTheDocument();
    expect(screen.getByText('Exercise')).toBeInTheDocument();

    // Journal Timeline
    expect(screen.getByText('Productive day!')).toBeInTheDocument();
    expect(screen.getByText(/Woke up early/)).toBeInTheDocument();
  });

  it('triggers onPeriodChange when clicking Monthly tab', () => {
    const onPeriodChange = vi.fn();
    render(
      <SummaryDashboard
        data={mockSummaryData}
        period="weekly"
        isLoading={false}
        onPeriodChange={onPeriodChange}
        onNavigate={() => {}}
      />,
    );

    const monthlyBtn = screen.getByRole('button', { name: 'Monthly' });
    fireEvent.click(monthlyBtn);
    expect(onPeriodChange).toHaveBeenCalledWith('monthly');
  });

  it('triggers onNavigate when clicking prev/next arrows', () => {
    const onNavigate = vi.fn();
    render(
      <SummaryDashboard
        data={mockSummaryData}
        period="weekly"
        isLoading={false}
        onPeriodChange={() => {}}
        onNavigate={onNavigate}
      />,
    );

    const prevBtn = screen.getByRole('button', { name: 'Previous period' });
    const nextBtn = screen.getByRole('button', { name: 'Next period' });

    fireEvent.click(prevBtn);
    expect(onNavigate).toHaveBeenCalledWith('prev');

    fireEvent.click(nextBtn);
    expect(onNavigate).toHaveBeenCalledWith('next');
  });
});
