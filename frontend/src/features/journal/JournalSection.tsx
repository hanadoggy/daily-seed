import { useCallback, useEffect, useRef, useState } from 'react';
import { BookOpen, Check, Loader2 } from 'lucide-react';
import { useAppStore } from '@/store/useAppStore';
import { useIsReadOnly } from '@/hooks/useIsReadOnly';
import type { Journal } from '@/types';

export function JournalSection() {
  const { dailyRecord, saveJournal, isLoading } = useAppStore();
  const isReadOnly = useIsReadOnly();
  const [oneLineReview, setOneLineReview] = useState('');
  const [threeLineDiary, setThreeLineDiary] = useState('');
  const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'saved'>('idle');
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Sync local state when dailyRecord changes (e.g., navigating dates).
  useEffect(() => {
    if (dailyRecord) {
      setOneLineReview(dailyRecord.journal.oneLineReview);
      setThreeLineDiary(dailyRecord.journal.threeLineDiary);
      setSaveStatus('idle');
    }
  }, [dailyRecord?.id]);

  const debouncedSave = useCallback(
    (journal: Journal) => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
      setSaveStatus('saving');
      debounceRef.current = setTimeout(async () => {
        await saveJournal(journal);
        setSaveStatus('saved');
        setTimeout(() => setSaveStatus('idle'), 2000);
      }, 1000);
    },
    [saveJournal],
  );

  // Cleanup timeout on unmount.
  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  const handleOneLineChange = (value: string) => {
    setOneLineReview(value);
    debouncedSave({ oneLineReview: value, threeLineDiary });
  };

  const handleDiaryChange = (value: string) => {
    setThreeLineDiary(value);
    debouncedSave({ oneLineReview, threeLineDiary: value });
  };

  if (isLoading || !dailyRecord) return null;

  return (
    <div className="rounded-2xl border border-border bg-card p-5 shadow-sm">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <BookOpen className="h-4 w-4 text-muted-foreground" />
          <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Reflection
          </h3>
        </div>
        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          {saveStatus === 'saving' && (
            <>
              <Loader2 className="h-3 w-3 animate-spin" />
              <span>Saving…</span>
            </>
          )}
          {saveStatus === 'saved' && (
            <>
              <Check className="h-3 w-3 text-green-500" />
              <span className="text-green-500">Saved</span>
            </>
          )}
        </div>
      </div>

      <div className="space-y-4">
        <div>
          <label
            htmlFor="one-line-review"
            className="block text-sm font-medium text-foreground mb-1.5"
          >
            One-line review
          </label>
          <input
            id="one-line-review"
            type="text"
            readOnly={isReadOnly}
            className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-1 transition-shadow read-only:opacity-60 read-only:cursor-default read-only:focus:ring-0"
            placeholder={isReadOnly ? "No review added." : "How was your day in one sentence?"}
            value={oneLineReview}
            onChange={(e) => handleOneLineChange(e.target.value)}
          />
        </div>

        <div>
          <label
            htmlFor="three-line-diary"
            className="block text-sm font-medium text-foreground mb-1.5"
          >
            Three-line diary
          </label>
          <textarea
            id="three-line-diary"
            rows={3}
            readOnly={isReadOnly}
            className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-1 resize-none transition-shadow read-only:opacity-60 read-only:cursor-default read-only:focus:ring-0"
            placeholder={isReadOnly ? "No diary added." : "Reflect on today…"}
            value={threeLineDiary}
            onChange={(e) => handleDiaryChange(e.target.value)}
          />
        </div>
      </div>
    </div>
  );
}
