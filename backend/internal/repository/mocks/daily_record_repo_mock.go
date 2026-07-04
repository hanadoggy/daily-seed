package mocks

import (
	"context"

	"daily-seed/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/stretchr/testify/mock"
)

type MockDailyRecordRepository struct {
	mock.Mock
}

func (m *MockDailyRecordRepository) FindByDate(ctx context.Context, date string) (*model.DailyRecord, error) {
	args := m.Called(ctx, date)
	if args.Get(0) != nil {
		return args.Get(0).(*model.DailyRecord), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDailyRecordRepository) Upsert(ctx context.Context, record *model.DailyRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockDailyRecordRepository) PatchByDate(ctx context.Context, date string, setFields bson.M) error {
	args := m.Called(ctx, date, setFields)
	return args.Error(0)
}
