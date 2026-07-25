import { render, screen, fireEvent, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { JournalSection } from '../JournalSection';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('JournalSection', () => {
  const mockSaveJournal = vi.fn();

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      dailyRecord: {
        journal: {
          oneLineReview: 'Great day',
          threeLineDiary: '1. Ate food\n2. Coded\n3. Slept',
        },
      },
      saveJournal: mockSaveJournal,
    });
    vi.clearAllMocks();
  });

  it('renders inputs with values from store', () => {
    render(<JournalSection />);
    expect(screen.getByDisplayValue('Great day')).toBeInTheDocument();
    expect(screen.getByDisplayValue(/Ate food/)).toBeInTheDocument();
  });

  it('calls saveJournal after typing', async () => {
    vi.useFakeTimers();
    render(<JournalSection />);
    
    const input = screen.getByDisplayValue('Great day');
    fireEvent.change(input, { target: { value: 'Awesome day' } });
    
    await act(async () => {
      vi.runAllTimers();
    });
    
    expect(mockSaveJournal).toHaveBeenCalledWith({
      oneLineReview: 'Awesome day',
      threeLineDiary: '1. Ate food\n2. Coded\n3. Slept',
    });
    
    vi.useRealTimers();
  });
});
