import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Calendar } from '../Calendar';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

vi.mock('@/lib/dayjs', async (importOriginal) => {
  const mod = await importOriginal();
  return {
    ...(mod as any),
    todayJST: () => '2023-10-10',
  };
});

describe('Calendar', () => {
  const mockSetDateAndFetch = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useAppStore as any).mockReturnValue({
      selectedDate: '2023-10-10',
      setDateAndFetch: mockSetDateAndFetch,
      existingRecordDates: ['2023-10-15'],
      fetchExistingRecordDates: vi.fn(),
    });
  });

  it('renders current month and handles date click', () => {
    render(<Calendar />);
    
    expect(screen.getByText('October 2023')).toBeDefined();
    
    // Find a specific day button (15th)
    const dayBtn = screen.getByText('15');
    fireEvent.click(dayBtn);
    
    expect(mockSetDateAndFetch).toHaveBeenCalledWith('2023-10-15');
  });

  it('can navigate to previous and next month', () => {
    render(<Calendar />);
    
    expect(screen.getByText('October 2023')).toBeDefined();
    
    // Prev month
    const buttons = screen.getAllByRole('button');
    const prevBtn = buttons[0];
    const nextBtn = buttons[1];
    
    fireEvent.click(prevBtn);
    expect(screen.getByText('September 2023')).toBeDefined();
    
    // Next month x2
    fireEvent.click(nextBtn);
    fireEvent.click(nextBtn);
    expect(screen.getByText('November 2023')).toBeDefined();
  });
});
