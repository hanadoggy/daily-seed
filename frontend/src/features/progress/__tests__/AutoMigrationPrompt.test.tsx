import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AutoMigrationPrompt } from '../AutoMigrationPrompt';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('AutoMigrationPrompt', () => {
  const mockMigrateTask = vi.fn();

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      tasks: [{ id: 't1', title: 'Task 1' }],
      taskProgress: [
        { taskId: 't1', title: 'Task 1', totalTarget: 10, totalCompleted: 10, percentage: 100 },
      ],
      migratingTaskIds: new Set(),
      migrateTask: mockMigrateTask,
    });
    vi.clearAllMocks();
  });

  it('renders prompt if task is 100% complete', () => {
    render(<AutoMigrationPrompt />);
    expect(screen.getByText(/You've completed 100% of the lifetime goal for "Task 1"/i)).toBeInTheDocument();
  });

  it('calls migrateTask on confirm', async () => {
    render(<AutoMigrationPrompt />);
    fireEvent.click(screen.getByRole('button', { name: 'Migrate Task' }));
    
    await waitFor(() => {
      expect(mockMigrateTask).toHaveBeenCalledWith('t1');
    });
  });

  it('does not render if no tasks are complete', () => {
    vi.mocked(useAppStore).mockReturnValue({
      tasks: [],
      taskProgress: [
        { taskId: 't1', title: 'Task 1', totalTarget: 10, totalCompleted: 5, percentage: 50 },
      ],
      migrateTask: mockMigrateTask,
    });
    
    const { container } = render(<AutoMigrationPrompt />);
    expect(container).toBeEmptyDOMElement();
  });
});
