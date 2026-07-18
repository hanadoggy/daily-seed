package daily

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDailyService struct {
	mock.Mock
}

func (m *MockDailyService) GetDailyRecord(ctx context.Context, date string) (*DailyRecord, error) {
	args := m.Called(ctx, date)
	if args.Get(0) != nil {
		return args.Get(0).(*DailyRecord), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDailyService) UpdateDailyRecord(ctx context.Context, date string, req *UpdateDailyRecordRequest) (*DailyRecord, error) {
	args := m.Called(ctx, date, req)
	if args.Get(0) != nil {
		return args.Get(0).(*DailyRecord), args.Error(1)
	}
	return nil, args.Error(1)
}
