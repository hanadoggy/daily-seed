import { render, screen, act } from '@testing-library/react';
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
        { taskId: 't1', title: 'Task 1', type: 'quantitative', totalTarget: 10, totalCompleted: 5, percentage: 50 },
        { taskId: 't2', title: 'Task 2', type: 'quantitative', totalTarget: 5, totalCompleted: 5, percentage: 100 },
        { taskId: 't3', title: 'Task 3', type: 'boolean', totalTarget: 0, totalCompleted: 2, percentage: 0 },
      ],
      isLoading: false,
      fetchProgress: vi.fn(),
    });
    vi.clearAllMocks();
  });

  it('renders task progress bars and sections', async () => {
    await act(async () => {
      render(<ProgressTracker />);
    });
    expect(screen.getByText('Project Progress')).toBeInTheDocument();
    expect(screen.getByText('Continuous Progress')).toBeInTheDocument();

    expect(screen.getByText('Task 1')).toBeInTheDocument();
    expect(screen.getByText('5/10')).toBeInTheDocument();

    expect(screen.getByText('Task 2')).toBeInTheDocument();
    expect(screen.getByText('5/5')).toBeInTheDocument();

    expect(screen.getByText('Task 3')).toBeInTheDocument();
    // 2 completions but endless task, so only '2' is rendered
    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('shows empty state if no progress', async () => {
    vi.mocked(useAppStore).mockReturnValue({ taskProgress: [], isLoading: false, fetchProgress: vi.fn() });
    let containerElement: HTMLElement | null = null;
    await act(async () => {
      const { container } = render(<ProgressTracker />);
      containerElement = container;
    });
    expect(containerElement).toBeEmptyDOMElement();
  });
});

