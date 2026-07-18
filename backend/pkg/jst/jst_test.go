package jst_test

import (
	"daily-seed/pkg/jst"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	loc := jst.Location()
	assert.NotNil(t, loc)
	assert.Equal(t, "Asia/Tokyo", loc.String())
}

func TestNow(t *testing.T) {
	now := jst.Now()
	assert.Equal(t, "Asia/Tokyo", now.Location().String())
}

func TestTodayString(t *testing.T) {
	today := jst.TodayString()
	_, err := time.Parse("2006-01-02", today)
	assert.NoError(t, err, "TodayString should return a valid YYYY-MM-DD date format")
}
