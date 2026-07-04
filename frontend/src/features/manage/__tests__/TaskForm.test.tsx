import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TaskForm } from '../TaskForm';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('TaskForm', () => {
  const mockAddTask = vi.fn();
  const mockEditTask = vi.fn();
  const mockOnClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useAppStore as any).mockReturnValue({
      addTask: mockAddTask,
      editTask: mockEditTask,
    });
  });

  it('creates a new task', async () => {
    render(<TaskForm onClose={mockOnClose} />);
    
    const titleInput = screen.getByPlaceholderText('e.g. Memorize Kanji');
    fireEvent.change(titleInput, { target: { value: 'New Task' } });
    
    const submitBtn = screen.getByRole('button', { name: 'Create Task' });
    fireEvent.click(submitBtn);
    
    await waitFor(() => {
      expect(mockAddTask).toHaveBeenCalledWith(expect.objectContaining({
        title: 'New Task',
        section: 'japanese',
        type: 'quantitative',
      }));
    });
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('validates empty title', () => {
    render(<TaskForm onClose={mockOnClose} />);
    const submitBtn = screen.getByRole('button', { name: 'Create Task' });
    expect((submitBtn as HTMLButtonElement).disabled).toBe(true);
  });
});
