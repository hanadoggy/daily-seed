import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { WeatherSelector } from '../WeatherSelector';
import { useAppStore } from '@/store/useAppStore';

vi.mock('@/store/useAppStore', () => ({
  useAppStore: vi.fn(),
}));

describe('WeatherSelector', () => {
  const mockUpdateWeather = vi.fn();

  beforeEach(() => {
    vi.mocked(useAppStore).mockReturnValue({
      currentWeather: 'sunny',
      updateWeather: mockUpdateWeather,
    });
    vi.clearAllMocks();
  });

  it('renders all weathers', () => {
    render(<WeatherSelector />);
    expect(screen.getByTitle('Sunny')).toBeInTheDocument();
    expect(screen.getByTitle('Rainy')).toBeInTheDocument();
  });

  it('calls updateWeather on click', () => {
    render(<WeatherSelector />);
    fireEvent.click(screen.getByTitle('Rainy'));
    expect(mockUpdateWeather).toHaveBeenCalledWith('rainy');
  });
});
