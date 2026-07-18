import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ProgressTracker } from '../ProgressTracker';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('ProgressTracker', () => {
  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      taskProgress: [
        { taskId: 't1', title: 'Task 1', totalTarget: 10, totalCompleted: 5, percentage: 50 },
        { taskId: 't2', title: 'Task 2', totalTarget: 5, totalCompleted: 5, percentage: 100 },
      ],
      isLoading: false,
      fetchProgress: vi.fn(),
    });
    vi.clearAllMocks();
  });

  it('renders task progress bars', () => {
    render(<ProgressTracker />);
    expect(screen.getByText('Task 1')).toBeInTheDocument();
    expect(screen.getByText('5/10')).toBeInTheDocument();
    
    expect(screen.getByText('Task 2')).toBeInTheDocument();
    expect(screen.getByText('5/5')).toBeInTheDocument();
  });

  it('shows empty state if no progress', () => {
    vi.mocked(useAppStore).mockReturnValue({ taskProgress: [], isLoading: false, fetchProgress: vi.fn() });
    const { container } = render(<ProgressTracker />);
    expect(container).toBeEmptyDOMElement();
  });
});
