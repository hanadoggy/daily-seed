import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ContextModeToggle } from '../ContextModeToggle';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('ContextModeToggle', () => {
  const mockUpdateMode = vi.fn();

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      currentMode: 'Growth',
      updateContextMode: mockUpdateMode,
    });
    vi.clearAllMocks();
  });

  it('renders all modes', () => {
    render(<ContextModeToggle />);
    expect(screen.getByText('Growth')).toBeInTheDocument();
    expect(screen.getByText('Rest')).toBeInTheDocument();
    expect(screen.getByText('Office')).toBeInTheDocument();
    expect(screen.getByText('Remote')).toBeInTheDocument();
  });

  it('calls updateContextMode on click', () => {
    render(<ContextModeToggle />);
    fireEvent.click(screen.getByText('Rest'));
    expect(mockUpdateMode).toHaveBeenCalledWith('Rest');
  });
});
