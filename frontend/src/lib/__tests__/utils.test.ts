import { describe, it, expect } from 'vitest';
import { cn } from '../utils';

describe('cn utility', () => {
  it('combines class names correctly', () => {
    expect(cn('bg-red-500', 'text-white')).toBe('bg-red-500 text-white');
  });

  it('merges tailwind conflict classes correctly', () => {
    expect(cn('px-2 py-1', 'p-4')).toBe('p-4');
    expect(cn('text-red-500', 'text-blue-500')).toBe('text-blue-500');
  });

  it('handles conditional class names, undefined, null, and empty strings', () => {
    const isTrue = true;
    const isFalse = false;
    expect(cn('base', isTrue && 'active', isFalse && 'disabled', null, undefined, '')).toBe('base active');
  });

  it('handles arrays and nested arrays of classes', () => {
    expect(cn(['btn', 'btn-primary'], ['shadow-md'])).toBe('btn btn-primary shadow-md');
  });
});
