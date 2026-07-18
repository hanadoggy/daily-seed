package daily

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

type MockDailyRecordRepository struct {
	mock.Mock
}

func (m *MockDailyRecordRepository) FindByDate(ctx context.Context, date string) (*DailyRecord, error) {
	args := m.Called(ctx, date)
	if args.Get(0) != nil {
		return args.Get(0).(*DailyRecord), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDailyRecordRepository) Upsert(ctx context.Context, record *DailyRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockDailyRecordRepository) PatchByDate(ctx context.Context, date string, setFields bson.M) error {
	args := m.Called(ctx, date, setFields)
	return args.Error(0)
}

func (m *MockDailyRecordRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDailyRecordRepository) FindBetweenDates(ctx context.Context, startDate string, endDate string) ([]*DailyRecord, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) != nil {
		return args.Get(0).([]*DailyRecord), args.Error(1)
	}
	return nil, args.Error(1)
}
