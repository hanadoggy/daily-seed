import { Sprout, Coffee, Building, Home } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/store/useAppStore';
import type { ContextMode } from '@/types';
import { cn } from '@/lib/utils';

const MODES: { value: ContextMode; label: string; icon: typeof Sprout; color: string }[] = [
  { value: 'Growth', label: 'Growth', icon: Sprout, color: 'text-emerald-400' },
  { value: 'Rest', label: 'Rest', icon: Coffee, color: 'text-amber-400' },
  { value: 'Office', label: 'Office', icon: Building, color: 'text-blue-400' },
  { value: 'Remote', label: 'Remote', icon: Home, color: 'text-purple-400' },
];

export function ContextModeToggle() {
  const { currentMode, updateContextMode } = useAppStore();

  return (
    <div className="flex gap-2">
      {MODES.map(({ value, label, icon: Icon, color }) => {
        const isActive = currentMode === value;
        return (
          <Button
            key={value}
            variant="outline"
            onClick={() => updateContextMode(value)}
            className={cn(
              'flex-1 gap-2 py-5 transition-all duration-300 border',
              isActive
                ? 'bg-mode-accent-soft border-mode-accent text-mode-accent shadow-sm'
                : 'border-border text-muted-foreground hover:text-foreground hover:border-muted-foreground',
            )}
          >
            <Icon className={cn('h-4 w-4', isActive && color)} />
            <span className="text-sm font-medium">{label}</span>
          </Button>
        );
      })}
    </div>
  );
}
