export type ErrorAction =
  | 'fetchDailyRecord'
  | 'patchDailyRecord'
  | 'fetchMasterData'
  | 'createTask'
  | 'updateTask'
  | 'archiveTask'
  | 'migrateTask'
  | 'fetchProgress'
  | 'createHabit'
  | 'updateHabit'
  | 'archiveHabit'
  | 'fetchHeatmap'
  | 'fetchSummary'
  | 'fetchStreaks'
  | 'fetchExistingDates';

const SPECIFIC_MESSAGES: Partial<Record<ErrorAction, Record<number, string>>> = {
  createTask: {
    400: '태스크 입력값이 올바르지 않습니다. 다시 확인해주세요.',
  },
  updateTask: {
    400: '태스크 수정 내용이 올바르지 않습니다.',
    404: '수정하려는 태스크를 찾을 수 없습니다.',
  },
  archiveTask: {
    400: '이미 아카이브된 태스크입니다.',
    404: '삭제하려는 태스크를 찾을 수 없습니다.',
  },
  migrateTask: {
    400: '마이그레이션 조건이 충족되지 않았습니다.',
    404: '마이그레이션 대상 태스크를 찾을 수 없습니다.',
    409: '이미 마이그레이션이 진행 중입니다.',
  },
  createHabit: {
    400: '습관 입력값이 올바르지 않습니다. 다시 확인해주세요.',
  },
  updateHabit: {
    400: '습관 수정 내용이 올바르지 않습니다.',
    404: '수정하려는 습관을 찾을 수 없습니다.',
  },
  archiveHabit: {
    404: '삭제하려는 습관을 찾을 수 없습니다.',
  },
  patchDailyRecord: {
    400: '기록 업데이트 내용이 올바르지 않습니다.',
    404: '해당 날짜의 기록을 찾을 수 없습니다.',
  },
  fetchDailyRecord: {
    400: '날짜 형식이 올바르지 않습니다.',
    404: '해당 날짜의 기록을 찾을 수 없습니다.',
  },
};

const COMMON_MESSAGES: Record<number, string> = {
  400: '요청이 올바르지 않습니다.',
  404: '요청한 데이터를 찾을 수 없습니다.',
  409: '다른 요청과 충돌이 발생했습니다. 다시 시도해주세요.',
  500: '서버에 문제가 발생했습니다. 잠시 후 다시 시도해주세요.',
};

export function getErrorMessage(action: ErrorAction, status: number): string {
  if (status === 0) {
    return '서버에 연결할 수 없습니다. 네트워크 상태를 확인해주세요.';
  }

  const actionMap = SPECIFIC_MESSAGES[action];
  if (actionMap && actionMap[status]) {
    return actionMap[status];
  }

  if (COMMON_MESSAGES[status]) {
    return COMMON_MESSAGES[status];
  }

  return '알 수 없는 오류가 발생했습니다. 잠시 후 다시 시도해주세요.';
}
