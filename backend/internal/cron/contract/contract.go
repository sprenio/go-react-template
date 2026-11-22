package contract

import "context"

type JobExecutionFunction func(jobContext context.Context) error

type JobRegistrar interface {
	AddJob(scheduleExpression string, jobName string, executionFunction JobExecutionFunction) error
}

