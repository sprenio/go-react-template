package cron

import (
	"context"
	"fmt"
	"time"

	"backend/internal/contexthelper"
	"backend/internal/cron/contract"
	"backend/internal/cron/jobs/heartbeat"
	"backend/internal/cron/jobs/retrypending"
	"backend/pkg/logger"

	robfigCron "github.com/robfig/cron/v3"
)

type CronScheduler struct {
	appContext    context.Context
	cronRunner    *robfigCron.Cron
	loggerAdapter *cronLoggerAdapter
}

func NewScheduler(appContext context.Context) *CronScheduler {
	loggerAdapter := newCronLoggerAdapter(appContext)

	cronRunnerInstance := robfigCron.New(
		robfigCron.WithChain(
			robfigCron.Recover(loggerAdapter),
		),
		robfigCron.WithLogger(loggerAdapter),
	)

	return &CronScheduler{
		appContext:    appContext,
		cronRunner:    cronRunnerInstance,
		loggerAdapter: loggerAdapter,
	}
}

func (scheduler *CronScheduler) RegisterScheduledJobs() error {
	if err := heartbeat.Register(scheduler); err != nil {
		return err
	}

	if err := retrypending.Register(scheduler); err != nil {
		return err
	}

	return nil
}

func (scheduler *CronScheduler) AddJob(
	scheduleExpression string,
	jobName string,
	jobExecutionFunction contract.JobExecutionFunction,
) error {
	entryID, err := scheduler.cronRunner.AddFunc(scheduleExpression, func() {
		scheduler.runJobSafely(jobName, jobExecutionFunction)
	})

	if err != nil {
		logger.ErrorCtx(scheduler.appContext, "Failed to schedule job", "jobName", jobName, "scheduleExpression", scheduleExpression, "error", err)
		return err
	}

	scheduler.loggerAdapter.RegisterJob(int(entryID), jobName)

	logger.InfoCtx(scheduler.appContext, "Registered cron job", "jobName", jobName, "scheduleExpression", scheduleExpression)
	return nil
}

func (scheduler *CronScheduler) runJobSafely(
	jobName string,
	jobExecutionFunction contract.JobExecutionFunction,
) {
	jobContext := contexthelper.SetRequestID(scheduler.appContext, "cron:"+jobName)
	startTime := time.Now()

	logger.InfoCtx(jobContext, "Job started", "jobName", jobName)

	executionError := jobExecutionFunction(jobContext)
	executionDuration := time.Since(startTime)

	if executionError != nil {
		logger.ErrorCtx(jobContext, fmt.Sprintf("Job [%s] failed after %v", jobName, executionDuration), "error", executionError)
		return
	}

	logger.InfoCtx(jobContext, fmt.Sprintf("Job [%s] completed successfully in %v", jobName, executionDuration))
}

func (scheduler *CronScheduler) Start() {
	scheduler.cronRunner.Start()
}

func (scheduler *CronScheduler) Stop() context.Context {
	return scheduler.cronRunner.Stop()
}
