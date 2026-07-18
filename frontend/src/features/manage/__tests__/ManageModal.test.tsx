import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ManageModal } from '../ManageModal';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('ManageModal', () => {
  const mockArchiveTask = vi.fn();
  const mockArchiveHabit = vi.fn();

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      tasks: [
        { id: 't1', title: 'Active Task', status: 'active', section: 'dev', type: 'boolean', conditions: { mode: [], weather: [] } },
        { id: 't2', title: 'Archived Task', status: 'archived', section: 'dev', type: 'boolean', conditions: { mode: [], weather: [] } },
      ],
      habits: [
        { id: 'h1', title: 'Health Habit', category: 'health', status: 'active' },
      ],
      archiveTask: mockArchiveTask,
      archiveHabit: mockArchiveHabit,
    });
    vi.clearAllMocks();
  });

  it('renders tasks tab by default and switches to habits', () => {
    render(<ManageModal open={true} onOpenChange={vi.fn()} />);
    
    expect(screen.getByText('Active Task')).toBeInTheDocument();
    expect(screen.queryByText('Health Habit')).not.toBeInTheDocument();
    
    // Switch to habits
    fireEvent.click(screen.getByText(/Habits/));
    
    expect(screen.getByText('Health Habit')).toBeInTheDocument();
    expect(screen.queryByText('Active Task')).not.toBeInTheDocument();
  });

  it('filters active and archived tasks', () => {
    render(<ManageModal open={true} onOpenChange={vi.fn()} />);
    
    expect(screen.getByText('Active Task')).toBeInTheDocument();
    expect(screen.queryByText('Archived Task')).not.toBeInTheDocument();
    
    // Switch to archived
    fireEvent.click(screen.getByText('Archived'));
    
    expect(screen.getByText('Archived Task')).toBeInTheDocument();
    expect(screen.queryByText('Active Task')).not.toBeInTheDocument();
  });

  it('opens create forms when clicking add buttons', () => {
    render(<ManageModal open={true} onOpenChange={vi.fn()} />);
    
    // Click Add Task
    fireEvent.click(screen.getByText('Add Task'));
    
    // Check if TaskForm is rendered (has specific inputs)
    expect(screen.getByPlaceholderText('e.g. Memorize Kanji')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create Task' })).toBeInTheDocument();
    
    // Switch back to habits tab which should reset the form state
    fireEvent.click(screen.getByText(/Habits/));
    fireEvent.click(screen.getByText('Add Habit'));
    
    // Check if HabitForm is rendered
    expect(screen.getByPlaceholderText('e.g. Morning stretching routine')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Create Habit' })).toBeInTheDocument();
  });

  it('handles delete confirmation for tasks', async () => {
    render(<ManageModal open={true} onOpenChange={vi.fn()} />);
    
    // The trash icon is inside a button, there's only one active task so 1 trash button
    // It's the second button (first is edit)
    const buttons = screen.getAllByRole('button');
    const trashBtn = buttons.find(b => b.innerHTML.includes('lucide-trash'));
    if (!trashBtn) throw new Error('Trash button not found');
    
    fireEvent.click(trashBtn);
    
    // The 'Sure?' text should appear
    expect(screen.getByText('Sure?')).toBeInTheDocument();
    
    // Click confirm (check icon)
    const confirmButtons = screen.getAllByRole('button').filter(b => b.innerHTML.includes('lucide-check'));
    fireEvent.click(confirmButtons[0]);
    
    await waitFor(() => {
      expect(mockArchiveTask).toHaveBeenCalledWith('t1');
    });
  });
});
