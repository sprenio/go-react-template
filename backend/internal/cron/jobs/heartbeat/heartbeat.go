package heartbeat

import (
	"context"
	"time"

	"backend/internal/contexthelper"
	"backend/internal/cron/contract"
	"backend/pkg/logger"
)

const (
	scheduleExpression = "@every 5m"
	jobName            = "heartbeat"
	pingTimeout        = 5 * time.Second
)

func Register(registrar contract.JobRegistrar) error {
	return registrar.AddJob(scheduleExpression, jobName, execute)
}

func execute(jobContext context.Context) error {
	databaseConnection := contexthelper.GetDb(jobContext)
	if databaseConnection != nil {
		pingTimeoutContext, cancelPing := context.WithTimeout(jobContext, pingTimeout)
		defer cancelPing()

		if pingError := databaseConnection.PingContext(pingTimeoutContext); pingError != nil {
			return pingError
		}
	}

	logger.InfoCtx(jobContext, "Heartbeat check passed: database connection is alive")
	return nil
}
