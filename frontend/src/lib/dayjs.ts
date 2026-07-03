import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.tz.setDefault('Asia/Tokyo');

export function todayJST(): string {
  return dayjs().tz('Asia/Tokyo').format('YYYY-MM-DD');
}

export function dayjsJST(date?: string) {
  return dayjs.tz(date, 'Asia/Tokyo');
}

export default dayjs;
