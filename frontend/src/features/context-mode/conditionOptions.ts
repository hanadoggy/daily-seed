import { Sprout, Coffee, Building, Home, Sun, CloudRain } from 'lucide-react';
import type { ContextMode } from '@/types';

export const MODE_OPTIONS: { value: ContextMode; label: string; icon: typeof Sprout; color: string; bgColor: string; borderColor: string }[] = [
  { value: 'Growth', label: 'Growth', icon: Sprout, color: 'text-emerald-400', bgColor: 'bg-emerald-400/15', borderColor: 'border-emerald-400/60' },
  { value: 'Rest', label: 'Rest', icon: Coffee, color: 'text-amber-400', bgColor: 'bg-amber-400/15', borderColor: 'border-amber-400/60' },
  { value: 'Office', label: 'Office', icon: Building, color: 'text-blue-400', bgColor: 'bg-blue-400/15', borderColor: 'border-blue-400/60' },
  { value: 'Remote', label: 'Remote', icon: Home, color: 'text-purple-400', bgColor: 'bg-purple-400/15', borderColor: 'border-purple-400/60' },
];

export const WEATHER_OPTIONS: { value: string; label: string; icon: typeof Sun; color: string; bgColor: string; borderColor: string }[] = [
  { value: 'sunny', label: 'Sunny', icon: Sun, color: 'text-orange-400', bgColor: 'bg-orange-400/15', borderColor: 'border-orange-400/60' },
  { value: 'rainy', label: 'Rainy', icon: CloudRain, color: 'text-blue-400', bgColor: 'bg-blue-400/15', borderColor: 'border-blue-400/60' },
];
