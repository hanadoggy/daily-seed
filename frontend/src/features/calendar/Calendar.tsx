import { useState } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAppStore } from '@/store/useAppStore';
import { dayjsJST, todayJST } from '@/lib/dayjs';
import { cn } from '@/lib/utils';

const WEEKDAY_LABELS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

export function Calendar() {
  const { selectedDate, setDateAndFetch } = useAppStore();
  const today = todayJST();

  const [viewMonth, setViewMonth] = useState(() => {
    const d = selectedDate || today;
    return dayjsJST(d).startOf('month');
  });

  const year = viewMonth.year();
  const month = viewMonth.month();
  const daysInMonth = viewMonth.daysInMonth();

  // Sunday=0 based: startOf('month').day() returns 0=Sun
  const firstDayOfWeek = viewMonth.day(); // Sun=0

  const prevMonth = () => setViewMonth(viewMonth.subtract(1, 'month'));
  const nextMonth = () => setViewMonth(viewMonth.add(1, 'month'));

  const handleDateClick = (day: number) => {
    const dateStr = viewMonth.date(day).format('YYYY-MM-DD');
    setDateAndFetch(dateStr);
  };

  const cells: (number | null)[] = [];
  for (let i = 0; i < firstDayOfWeek; i++) cells.push(null);
  for (let d = 1; d <= daysInMonth; d++) cells.push(d);

  return (
    <div className="select-none">
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <Button variant="ghost" size="icon" onClick={prevMonth} className="h-8 w-8">
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <span className="text-sm font-semibold tracking-wide">
          {viewMonth.format('MMMM YYYY')}
        </span>
        <Button variant="ghost" size="icon" onClick={nextMonth} className="h-8 w-8">
          <ChevronRight className="h-4 w-4" />
        </Button>
      </div>

      {/* Weekday labels */}
      <div className="grid grid-cols-7 gap-1 mb-1">
        {WEEKDAY_LABELS.map((label) => (
          <div
            key={label}
            className={cn(
              "text-center text-[11px] font-medium py-1",
              label === 'Sat' ? 'text-cal-sat' : label === 'Sun' ? 'text-cal-sun' : 'text-muted-foreground'
            )}
          >
            {label}
          </div>
        ))}
      </div>

      {/* Day grid */}
      <div className="grid grid-cols-7 gap-1">
        {cells.map((day, i) => {
          if (day === null) {
            return <div key={`empty-${i}`} />;
          }

          const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
          const isToday = dateStr === today;
          const isSelected = dateStr === selectedDate;

          return (
            <button
              key={dateStr}
              onClick={() => handleDateClick(day)}
              className={cn(
                'relative h-9 w-full rounded-lg text-sm font-medium transition-all duration-200',
                'hover:bg-accent hover:text-accent-foreground',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
                isSelected && 'bg-mode-accent text-white shadow-md hover:bg-mode-accent',
                !isSelected && isToday && 'ring-1 ring-mode-accent text-mode-accent',
                !isSelected && !isToday && i % 7 === 6 && 'text-cal-sat',
                !isSelected && !isToday && i % 7 === 0 && 'text-cal-sun',
              )}
            >
              {day}
            </button>
          );
        })}
      </div>
    </div>
  );
}
