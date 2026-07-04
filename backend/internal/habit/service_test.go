package habit

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHabitService_Create(t *testing.T) {
	tests := []struct {
		name        string
		habit       *Habit
		mockSetup   func(*MockHabitRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful creation",
			habit: &Habit{
				Title:       "Morning Walk",
				Category:    "Health",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: missing title (edge/abnormal)",
			habit: &Habit{
				Title:    "",
				Category: "Health",
			},
			mockSetup:   func(m *MockHabitRepository) {},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: whitespace title (edge)",
			habit: &Habit{
				Title:    "   ",
				Category: "Health",
			},
			mockSetup:   func(m *MockHabitRepository) {},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: missing category (edge/abnormal)",
			habit: &Habit{
				Title:    "Walk",
				Category: "",
			},
			mockSetup:   func(m *MockHabitRepository) {},
			expectError: true,
			errContains: "category is required",
		},
		{
			name: "Fail: repository error",
			habit: &Habit{
				Title:    "Walk",
				Category: "Health",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "creating habit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockHabitRepository)
			tt.mockSetup(repo)
			service := NewHabitService(repo)

			created, err := service.Create(context.Background(), tt.habit)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
				assert.Nil(t, created)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, created)
				assert.NotEmpty(t, created.ID)
				assert.Equal(t, "active", created.Status)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestHabitService_Update(t *testing.T) {
	existingHabit := &Habit{ID: "habit_1", Title: "Walk", Category: "Health", Status: "active"}

	tests := []struct {
		name        string
		habit       *Habit
		mockSetup   func(*MockHabitRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful update",
			habit: &Habit{
				ID:       "habit_1",
				Title:    "Evening Walk",
				Category: "Health",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(existingHabit, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: habit not found",
			habit: &Habit{
				ID:    "habit_unknown",
				Title: "Walk",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_unknown").Return(nil, nil)
			},
			expectError: true,
			errContains: "habit not found",
		},
		{
			name: "Fail: FindByID repo error",
			habit: &Habit{
				ID:    "habit_1",
				Title: "Walk",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding habit",
		},
		{
			name: "Fail: empty title (edge)",
			habit: &Habit{
				ID:    "habit_1",
				Title: "   ",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(existingHabit, nil)
			},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: update repo error",
			habit: &Habit{
				ID:    "habit_1",
				Title: "Walk 2",
			},
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(existingHabit, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "updating habit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockHabitRepository)
			tt.mockSetup(repo)
			service := NewHabitService(repo)

			updated, err := service.Update(context.Background(), tt.habit)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
				assert.Nil(t, updated)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updated)
				assert.Equal(t, existingHabit.Status, updated.Status) // Status should remain same
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestHabitService_Archive(t *testing.T) {
	existingHabit := &Habit{ID: "habit_1", Status: "active"}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(*MockHabitRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful archive",
			id:   "habit_1",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(existingHabit, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(h *Habit) bool {
					return h.Status == "archived"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: not found",
			id:   "habit_unknown",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_unknown").Return(nil, nil)
			},
			expectError: true,
			errContains: "habit not found",
		},
		{
			name: "Fail: FindByID error",
			id:   "habit_1",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding habit",
		},
		{
			name: "Fail: Update error",
			id:   "habit_1",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(existingHabit, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "archiving habit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockHabitRepository)
			tt.mockSetup(repo)
			service := NewHabitService(repo)

			err := service.Archive(context.Background(), tt.id)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestHabitService_Get(t *testing.T) {
	existingHabit := &Habit{ID: "habit_1", Title: "Walk"}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(*MockHabitRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful get",
			id:   "habit_1",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(existingHabit, nil)
			},
			expectError: false,
		},
		{
			name: "Fail: not found",
			id:   "habit_unknown",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_unknown").Return(nil, nil)
			},
			expectError: true,
			errContains: "habit not found",
		},
		{
			name: "Fail: db error",
			id:   "habit_1",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindByID", mock.Anything, "habit_1").Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding habit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockHabitRepository)
			tt.mockSetup(repo)
			service := NewHabitService(repo)

			h, err := service.Get(context.Background(), tt.id)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, h)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, h)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestHabitService_List(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockHabitRepository)
		expectError bool
	}{
		{
			name: "Pass: successful list",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindActiveHabits", mock.Anything).Return([]Habit{{ID: "1"}, {ID: "2"}}, nil)
			},
			expectError: false,
		},
		{
			name: "Fail: db error",
			mockSetup: func(m *MockHabitRepository) {
				m.On("FindActiveHabits", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockHabitRepository)
			tt.mockSetup(repo)
			service := NewHabitService(repo)

			list, err := service.List(context.Background())
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, list)
			} else {
				assert.NoError(t, err)
				assert.Len(t, list, 2)
			}
			repo.AssertExpectations(t)
		})
	}
}
