package schedule

import (
	"testing"
	"time"
)

type MockTask struct {
}

func (m *MockTask) Exec() {
}

func TestAdd(t *testing.T) {
	Add(&MockTask{}, time.Now())
	Add(&MockTask{}, time.Now())

	if len(s.jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(s.jobs))
	}
}

func TestRunSchedule(t *testing.T) {
	RunSchedule()

	time.Sleep(1000 * time.Millisecond)

	if len(s.jobs) != 0 {
		t.Errorf("Expected 0 jobs, got %d", len(s.jobs))
	}
}
