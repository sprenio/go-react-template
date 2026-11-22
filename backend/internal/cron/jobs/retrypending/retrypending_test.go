package retrypending

import (
	"context"
	"testing"

	"backend/internal/cron/contract"
)

type mockRegistrar struct {
	registeredSchedule string
	registeredJobName  string
	registeredFunction contract.JobExecutionFunction
}

func (m *mockRegistrar) AddJob(scheduleExpression string, jobName string, executionFunction contract.JobExecutionFunction) error {
	m.registeredSchedule = scheduleExpression
	m.registeredJobName = jobName
	m.registeredFunction = executionFunction
	return nil
}

func TestRegisterRetryPendingJob(t *testing.T) {
	mock := &mockRegistrar{}
	err := Register(mock)
	if err != nil {
		t.Fatalf("expected nil error on Register, got: %v", err)
	}

	if mock.registeredSchedule != scheduleExpression {
		t.Errorf("expected schedule %q, got %q", scheduleExpression, mock.registeredSchedule)
	}

	if mock.registeredJobName != jobName {
		t.Errorf("expected job name %q, got %q", jobName, mock.registeredJobName)
	}

	if mock.registeredFunction == nil {
		t.Errorf("expected executionFunction to be non-nil")
	}
}

func TestExecuteFailsWithoutDatabaseInContext(t *testing.T) {
	emptyContext := context.Background()
	err := execute(emptyContext)
	if err == nil {
		t.Fatalf("expected error when database is missing from context, got nil")
	}
}
