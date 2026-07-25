import { describe, it, expect } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { HeatmapDashboard } from '../HeatmapDashboard';

describe('HeatmapDashboard component', () => {
  it('renders loading state when isLoading is true', () => {
    render(<HeatmapDashboard data={null} isLoading={true} />);
    expect(screen.getByText('Consistency Heatmap')).toBeInTheDocument();
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('renders empty heatmap safely when data is null or empty', () => {
    render(<HeatmapDashboard data={null} isLoading={false} />);
    expect(screen.getByText('Consistency Heatmap')).toBeInTheDocument();
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
    expect(screen.getByText('Less')).toBeInTheDocument();
  });

  it('renders days correctly with filters', () => {
    const mockData = {
      days: [
        {
          date: '2026-01-01',
          total: 5,
          habits: 2,
          sectionCounts: { dev: 3, japanese: 0 },
        },
      ],
    };

    render(<HeatmapDashboard data={mockData} isLoading={false} />);

    expect(screen.getByText('Consistency Heatmap')).toBeInTheDocument();

    // Check filters existence
    const devFilterBtn = screen.getByRole('button', { name: 'Dev' });
    expect(devFilterBtn).toBeInTheDocument();

    // Click Dev filter
    fireEvent.click(devFilterBtn);
    expect(devFilterBtn).toHaveClass('bg-primary');
  });
});
