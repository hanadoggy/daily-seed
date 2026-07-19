import { useAppStore } from '@/store/useAppStore';
import { todayJST } from '@/lib/dayjs';

export function useIsReadOnly() {
  const selectedDate = useAppStore((state) => state.selectedDate);
  const isAdminMode = useAppStore((state) => state.isAdminMode);
  return selectedDate !== todayJST() && !isAdminMode;
}
