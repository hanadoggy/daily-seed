import { describe, it, expect } from 'vitest';
import { getErrorMessage } from '../errorMessages';

describe('getErrorMessage', () => {
  it('returns network error message when status is 0', () => {
    expect(getErrorMessage('fetchDailyRecord', 0)).toBe(
      '서버에 연결할 수 없습니다. 네트워크 상태를 확인해주세요.',
    );
  });

  it('returns specific error message for action and status', () => {
    expect(getErrorMessage('createTask', 400)).toBe(
      '태스크 입력값이 올바르지 않습니다. 다시 확인해주세요.',
    );
    expect(getErrorMessage('migrateTask', 409)).toBe(
      '이미 마이그레이션이 진행 중입니다.',
    );
  });

  it('falls back to common error message if no specific message exists', () => {
    expect(getErrorMessage('fetchDailyRecord', 500)).toBe(
      '서버에 문제가 발생했습니다. 잠시 후 다시 시도해주세요.',
    );
    expect(getErrorMessage('createTask', 404)).toBe(
      '요청한 데이터를 찾을 수 없습니다.',
    );
  });

  it('returns default fallback message for unknown status code', () => {
    expect(getErrorMessage('fetchDailyRecord', 418)).toBe(
      '알 수 없는 오류가 발생했습니다. 잠시 후 다시 시도해주세요.',
    );
  });
});
