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

  it('displays locked indicator on creation and includes unit field', async () => {
    render(<TaskForm onClose={mockOnClose} />);
    
    // Check locked indicators
    expect(screen.getAllByText('Locked after creation')).toHaveLength(3);
    
    const titleInput = screen.getByPlaceholderText('e.g. Memorize Kanji');
    fireEvent.change(titleInput, { target: { value: 'Read Book' } });
    
    const unitInput = screen.getByPlaceholderText('e.g. pages, mins, exercises');
    fireEvent.change(unitInput, { target: { value: 'pages' } });
    
    const submitBtn = screen.getByRole('button', { name: 'Create Task' });
    fireEvent.click(submitBtn);
    
    await waitFor(() => {
      expect(mockAddTask).toHaveBeenCalledWith(expect.objectContaining({
        title: 'Read Book',
        unit: 'pages',
      }));
    });
  });

  it('disables type, weather, and mode in edit mode', () => {
    const existingTask = {
      id: 't1',
      title: 'Existing Task',
      section: 'dev',
      type: 'quantitative',
      unit: 'pages',
      metrics: { dailyTarget: 10, totalTarget: 100 },
      conditions: { mode: ['Growth'], weather: ['sunny'] },
      startDate: '2026-01-01',
    };
    
    render(<TaskForm task={existingTask as any} onClose={mockOnClose} />);
    
    // Locked after creation indicators should not be present in edit mode
    expect(screen.queryByText('Locked after creation')).not.toBeInTheDocument();

    // Type select is disabled
    const typeSelect = screen.getByLabelText('Type');
    expect((typeSelect as HTMLSelectElement).disabled).toBe(true);

    // Weather and Mode buttons are disabled
    const sunnyButton = screen.getByRole('button', { name: 'Sunny' });
    expect((sunnyButton as HTMLButtonElement).disabled).toBe(true);
  });

  it('toggles weather and mode conditions on task creation', async () => {
    render(<TaskForm onClose={mockOnClose} />);

    const titleInput = screen.getByPlaceholderText('e.g. Memorize Kanji');
    fireEvent.change(titleInput, { target: { value: 'Outdoor Run' } });

    // Toggle Rainy weather off
    const rainyButton = screen.getByRole('button', { name: 'Rainy' });
    fireEvent.click(rainyButton);

    // Toggle Rest mode off
    const restButton = screen.getByRole('button', { name: 'Rest' });
    fireEvent.click(restButton);

    const submitBtn = screen.getByRole('button', { name: 'Create Task' });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockAddTask).toHaveBeenCalledWith(expect.objectContaining({
        title: 'Outdoor Run',
        conditions: {
          weather: ['sunny'],
          mode: ['Growth', 'Office', 'Remote'],
        },
      }));
    });
  });

  it('hides dailyTarget input when task type is boolean', () => {
    render(<TaskForm onClose={mockOnClose} />);

    // Switch type to boolean
    const typeSelect = screen.getByLabelText('Type');
    fireEvent.change(typeSelect, { target: { value: 'boolean' } });

    // Daily target input should be hidden for boolean tasks
    expect(screen.queryByLabelText('Daily Target')).not.toBeInTheDocument();
  });
});

