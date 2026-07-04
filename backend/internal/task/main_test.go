package task_test

import (
	"daily-seed/internal/testutil"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.RunWithDB(m)
}
