import { Sun, CloudRain } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/store/useAppStore';
import { cn } from '@/lib/utils';

const WEATHERS: { value: string; label: string; icon: typeof Sun; color: string }[] = [
  { value: 'sunny', label: 'Sunny', icon: Sun, color: 'text-orange-400' },
  { value: 'rainy', label: 'Rainy', icon: CloudRain, color: 'text-blue-400' },
];

export function WeatherSelector() {
  const { currentWeather, updateWeather } = useAppStore();

  return (
    <div className="flex gap-2">
      {WEATHERS.map(({ value, label, icon: Icon, color }) => {
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
