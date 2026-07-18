import { Flower2, HeartPulse, Target, PiggyBank, Users } from 'lucide-react';

export const HABIT_CATEGORIES = [
  { value: 'mindfulness', label: 'Mindfulness', icon: Flower2, color: 'text-purple-400', bgColor: 'bg-purple-400/15', borderColor: 'border-purple-400/60' },
  { value: 'health', label: 'Health', icon: HeartPulse, color: 'text-emerald-400', bgColor: 'bg-emerald-400/15', borderColor: 'border-emerald-400/60' },
  { value: 'productivity', label: 'Productivity', icon: Target, color: 'text-blue-400', bgColor: 'bg-blue-400/15', borderColor: 'border-blue-400/60' },
  { value: 'finance', label: 'Finance', icon: PiggyBank, color: 'text-amber-400', bgColor: 'bg-amber-400/15', borderColor: 'border-amber-400/60' },
  { value: 'social', label: 'Social', icon: Users, color: 'text-pink-400', bgColor: 'bg-pink-400/15', borderColor: 'border-pink-400/60' },
];
