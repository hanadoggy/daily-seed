import { MODE_OPTIONS } from './conditionOptions';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/store/useAppStore';
import { useIsReadOnly } from '@/hooks/useIsReadOnly';
import { cn } from '@/lib/utils';

export function ContextModeToggle() {
  const { currentMode, updateContextMode } = useAppStore();
  const isReadOnly = useIsReadOnly();

  return (
    <div className="flex gap-2">
      {MODE_OPTIONS.map(({ value, label, icon: Icon, color }) => {
        const isActive = currentMode === value;
        return (
          <Button
            key={value}
            variant="outline"
            disabled={isReadOnly}
            onClick={() => !isReadOnly && updateContextMode(value)}
            className={cn(
              'flex-1 gap-2 h-[62px] transition-all duration-300 border',
              isActive
                ? 'bg-mode-accent-soft border-mode-accent text-mode-accent shadow-sm'
                : 'border-border text-muted-foreground hover:text-foreground hover:border-muted-foreground',
              isReadOnly && 'opacity-60 cursor-default hover:text-muted-foreground hover:border-border'
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
