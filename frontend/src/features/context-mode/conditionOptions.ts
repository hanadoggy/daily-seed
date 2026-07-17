import { Sprout, Coffee, Building, Home, Sun, CloudRain } from 'lucide-react';
import type { ContextMode } from '@/types';

export const MODE_OPTIONS: { value: ContextMode; label: string; icon: typeof Sprout; color: string }[] = [
  { value: 'Growth', label: 'Growth', icon: Sprout, color: 'text-emerald-400' },
  { value: 'Rest', label: 'Rest', icon: Coffee, color: 'text-amber-400' },
  { value: 'Office', label: 'Office', icon: Building, color: 'text-blue-400' },
  { value: 'Remote', label: 'Remote', icon: Home, color: 'text-purple-400' },
];

export const WEATHER_OPTIONS: { value: string; label: string; icon: typeof Sun; color: string }[] = [
  { value: 'sunny', label: 'Sunny', icon: Sun, color: 'text-orange-400' },
  { value: 'rainy', label: 'Rainy', icon: CloudRain, color: 'text-blue-400' },
];
