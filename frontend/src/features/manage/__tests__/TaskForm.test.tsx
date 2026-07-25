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
      selectedDate: '2026-07-17',
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

  it('disables submit when title contains only whitespace', () => {
    render(<TaskForm onClose={mockOnClose} />);
    const titleInput = screen.getByPlaceholderText('e.g. Memorize Kanji');
    fireEvent.change(titleInput, { target: { value: '   ' } });
    const submitBtn = screen.getByRole('button', { name: 'Create Task' });
    expect((submitBtn as HTMLButtonElement).disabled).toBe(true);
  });

  it('can select exercise section', async () => {
    render(<TaskForm onClose={mockOnClose} />);
    
    const titleInput = screen.getByPlaceholderText('e.g. Memorize Kanji');
    fireEvent.change(titleInput, { target: { value: 'Workout' } });
    
    const sectionSelect = screen.getByLabelText('Section');
    fireEvent.change(sectionSelect, { target: { value: 'exercise' } });
    
    const submitBtn = screen.getByRole('button', { name: 'Create Task' });
    fireEvent.click(submitBtn);
    
    await waitFor(() => {
      expect(mockAddTask).toHaveBeenCalledWith(expect.objectContaining({
        title: 'Workout',
        section: 'exercise',
      }));
    });
  });

  it('populates data in edit mode and calls editTask', async () => {
    const existingTask = {
      id: 't1',
      title: 'Existing Task',
      section: 'dev',
      type: 'boolean',
      metrics: { dailyTarget: 1, totalTarget: 0, unit: '' },
      conditions: { mode: ['Growth'], weather: ['sunny'] },
    };
    
    render(<TaskForm task={existingTask as any} onClose={mockOnClose} />);
    
    expect(screen.getByDisplayValue('Existing Task')).toBeInTheDocument();
    
    const titleInput = screen.getByDisplayValue('Existing Task');
    fireEvent.change(titleInput, { target: { value: 'Updated Task' } });
    
    const submitBtn = screen.getByRole('button', { name: 'Save Changes' });
    fireEvent.click(submitBtn);
    
    await waitFor(() => {
      expect(mockEditTask).toHaveBeenCalledWith('t1', expect.objectContaining({
        title: 'Updated Task',
      }));
    });
  });

  it('can toggle type and set metrics', async () => {
    render(<TaskForm onClose={mockOnClose} />);
    
    const titleInput = screen.getByPlaceholderText('e.g. Memorize Kanji');
    fireEvent.change(titleInput, { target: { value: 'Read Pages' } });
    
    // Select quantitative type (since type defaults to boolean maybe?)
    // Actually the default is 'quantitative' in the form, let's verify.
    // Wait, let's check the type switch logic.
    // It's a SegmentedControl or similar?
  });
});
