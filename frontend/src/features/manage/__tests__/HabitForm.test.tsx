import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { HabitForm } from '../HabitForm';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('HabitForm', () => {
  const mockAddHabit = vi.fn();
  const mockEditHabit = vi.fn();
  const mockOnClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useAppStore as any).mockReturnValue({
      addHabit: mockAddHabit,
      editHabit: mockEditHabit,
    });
  });

  it('creates a new habit', async () => {
    render(<HabitForm onClose={mockOnClose} />);
    
    const titleInput = screen.getByPlaceholderText('e.g. Morning stretching routine');
    fireEvent.change(titleInput, { target: { value: 'New Habit' } });
    
    // Click on Finance category
    const financeBtn = screen.getByRole('button', { name: 'Finance' });
    fireEvent.click(financeBtn);

    const submitBtn = screen.getByRole('button', { name: 'Create Habit' });
    fireEvent.click(submitBtn);
    
    await waitFor(() => {
      expect(mockAddHabit).toHaveBeenCalledWith(expect.objectContaining({
        title: 'New Habit',
        category: 'finance',
      }));
    });
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('validates empty title', () => {
    render(<HabitForm onClose={mockOnClose} />);
    const submitBtn = screen.getByRole('button', { name: 'Create Habit' });
    expect((submitBtn as HTMLButtonElement).disabled).toBe(true);
  });

  it('populates data in edit mode and calls editHabit', async () => {
    const existingHabit = {
      id: 'h1',
      title: 'Existing Habit',
      category: 'health',
    };
    
    render(<HabitForm habit={existingHabit as any} onClose={mockOnClose} />);
    
    expect(screen.getByDisplayValue('Existing Habit')).toBeInTheDocument();
    
    const titleInput = screen.getByDisplayValue('Existing Habit');
    fireEvent.change(titleInput, { target: { value: 'Updated Habit' } });
    
    const submitBtn = screen.getByRole('button', { name: 'Save Changes' });
    fireEvent.click(submitBtn);
    
    await waitFor(() => {
      expect(mockEditHabit).toHaveBeenCalledWith('h1', expect.objectContaining({
        title: 'Updated Habit',
      }));
    });
  });
});
