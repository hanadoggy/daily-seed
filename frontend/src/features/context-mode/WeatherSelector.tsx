import { WEATHER_OPTIONS } from './conditionOptions';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/store/useAppStore';
import { cn } from '@/lib/utils';

export function WeatherSelector() {
  const { currentWeather, updateWeather } = useAppStore();

  return (
    <div className="flex gap-2">
      {WEATHER_OPTIONS.map(({ value, label, icon: Icon, color }) => {
        const isActive = currentWeather === value;
        return (
          <Button
            key={value}
            variant="outline"
            size="icon"
            onClick={() => updateWeather(value)}
            className={cn(
              'h-[62px] w-[62px] transition-all duration-300 border',
              isActive
                ? 'bg-mode-accent-soft border-mode-accent text-mode-accent shadow-sm'
                : 'border-border text-muted-foreground hover:text-foreground hover:border-muted-foreground',
            )}
            title={label}
          >
            <Icon className={cn('h-5 w-5', isActive && color)} />
          </Button>
        );
      })}
    </div>
  );
}
