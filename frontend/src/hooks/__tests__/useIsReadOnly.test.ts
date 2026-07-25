import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useIsReadOnly } from '../useIsReadOnly';
import { useAppStore } from '@/store/useAppStore';
import { todayJST } from '@/lib/dayjs';

describe('useIsReadOnly hook', () => {
  beforeEach(() => {
    act(() => {
      useAppStore.setState({
        selectedDate: todayJST(),
        isAdminMode: false,
      });
    });
  });

  it('returns false when selectedDate is todayJST and isAdminMode is false', () => {
    const { result } = renderHook(() => useIsReadOnly());
    expect(result.current).toBe(false);
  });

  it('returns true when selectedDate is a past date and isAdminMode is false', () => {
    act(() => {
      useAppStore.setState({ selectedDate: '2020-01-01', isAdminMode: false });
    });
    const { result } = renderHook(() => useIsReadOnly());
    expect(result.current).toBe(true);
  });

  it('returns false when selectedDate is a past date BUT isAdminMode is true', () => {
    act(() => {
      useAppStore.setState({ selectedDate: '2020-01-01', isAdminMode: true });
    });
    const { result } = renderHook(() => useIsReadOnly());
    expect(result.current).toBe(false);
  });

  it('returns true when selectedDate is a future date and isAdminMode is false', () => {
    act(() => {
      useAppStore.setState({ selectedDate: '2099-12-31', isAdminMode: false });
    });
    const { result } = renderHook(() => useIsReadOnly());
    expect(result.current).toBe(true);
  });
});
