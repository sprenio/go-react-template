package cron

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestCronSchedulerRegistrationAndExecution(t *testing.T) {
	testContext := context.Background()
	scheduler := NewScheduler(testContext)

	var executionCount int32
	testJobName := "test_quick_job"
	testSchedule := "@every 1s"

	err := scheduler.AddJob(testSchedule, testJobName, func(jobContext context.Context) error {
		atomic.AddInt32(&executionCount, 1)
		return nil
	})

	if err != nil {
		t.Fatalf("Expected job registration to succeed, got error: %v", err)
	}

	scheduler.Start()
	time.Sleep(1500 * time.Millisecond)

	stopContext := scheduler.Stop()
	<-stopContext.Done()

	finalCount := atomic.LoadInt32(&executionCount)
	if finalCount < 1 {
		t.Errorf("Expected job to execute at least once, got: %d", finalCount)
	}
}

func TestRegisterScheduledJobs(t *testing.T) {
	testContext := context.Background()
	scheduler := NewScheduler(testContext)

	err := scheduler.RegisterScheduledJobs()
	if err != nil {
		t.Fatalf("Expected RegisterScheduledJobs to succeed, got: %v", err)
	}
}
