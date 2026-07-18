import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TaskItem } from '../TaskItem';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('TaskItem', () => {
  const mockUpdateTaskProgressOptimistic = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useAppStore as any).mockReturnValue(mockUpdateTaskProgressOptimistic);
  });

  it('renders boolean task and toggles completion', () => {
    const entry = { taskId: 't1', targetAmount: 1, actualAmount: 0, isCompleted: false };
    render(<TaskItem entry={entry} title="Read Book" type="boolean" />);
    
    const btn = screen.getByRole('button');
    fireEvent.click(btn);
    expect(mockUpdateTaskProgressOptimistic).toHaveBeenCalledWith('t1', 1);
  });

  it('renders quantitative task and updates progress', () => {
    const entry = { taskId: 't2', targetAmount: 10, actualAmount: 5, isCompleted: false };
    render(<TaskItem entry={entry} title="Pages Read" type="quantitative" />);
    
    const buttons = screen.getAllByRole('button');
    // buttons[0] is minus, buttons[1] is plus based on the layout
    fireEvent.click(buttons[1]); // Plus
    expect(mockUpdateTaskProgressOptimistic).toHaveBeenCalledWith('t2', 6);
    
    fireEvent.click(buttons[0]); // Minus
    expect(mockUpdateTaskProgressOptimistic).toHaveBeenCalledWith('t2', 4);
  });

  it('does not disable plus button when actualAmount reaches targetAmount (allows exceeding)', () => {
    const entry = { taskId: 't2', targetAmount: 10, actualAmount: 10, isCompleted: true };
    render(<TaskItem entry={entry} title="Pages Read" type="quantitative" />);
    
    const buttons = screen.getAllByRole('button');
    expect((buttons[1] as HTMLButtonElement).disabled).toBe(false); // Plus button
  });

  it('hides plus and minus buttons when isReadOnly is true', () => {
    const entry = { taskId: 't3', targetAmount: 10, actualAmount: 5, isCompleted: false };
    render(<TaskItem entry={entry} title="Pages Read" type="quantitative" isReadOnly={true} />);
    
    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(0);
  });

  it('hides plus and minus buttons when isArchived is true', () => {
    const entry = { taskId: 't4', targetAmount: 10, actualAmount: 5, isCompleted: false };
    render(<TaskItem entry={entry} title="Pages Read" type="quantitative" isArchived={true} />);
    
    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(0);
  });

  it('disables boolean checkbox when isReadOnly is true', () => {
    const entry = { taskId: 't5', targetAmount: 1, actualAmount: 0, isCompleted: false };
    render(<TaskItem entry={entry} title="Read Book" type="boolean" isReadOnly={true} />);
    
    const btn = screen.getByRole('button');
    expect((btn as HTMLButtonElement).disabled).toBe(true);
  });
});
