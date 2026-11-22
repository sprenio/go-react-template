package retrypending

import (
	"context"
	"fmt"
	"time"

	"backend/internal/contexthelper"
	"backend/internal/cron/contract"
	"backend/internal/queue"
	"backend/internal/queue/publisher"
	queueRepository "backend/internal/queue/repository"
	"backend/pkg/apperror"
	"backend/pkg/logger"

	"errors"
)

const (
	scheduleExpression         = "*/2 * * * *"
	jobName                    = "retry_pending_queues"
	maxPendingQueuesBatchLimit = 100
)

func Register(registrar contract.JobRegistrar) error {
	return registrar.AddJob(scheduleExpression, jobName, execute)
}

func execute(jobContext context.Context) error {
	databaseConnection := contexthelper.GetDb(jobContext)
	if databaseConnection == nil {
		return errors.New("database connection is not available in context")
	}

	queueRepositoryInstance := queueRepository.NewQueueRepository(databaseConnection)

	pendingQueues, fetchError := queueRepositoryInstance.GetRetryPendingQueues(jobContext, maxPendingQueuesBatchLimit)
	if fetchError != nil {
		return apperror.Wrap(fetchError, "failed to fetch retry pending queues")
	}

	if len(pendingQueues) == 0 {
		return nil
	}

	queuePublisherInstance := publisher.NewPublisher(jobContext)
	logger.InfoCtx(jobContext, fmt.Sprintf("Found %d queue(s) with status RETRY_PENDING", len(pendingQueues)))

	for _, queueEntry := range pendingQueues {
		var queuePrefix string
		var baseDelay time.Duration

		switch queueEntry.Task {
		case queue.QueueTaskRegisterUser:
			queuePrefix = queue.RegisterQueuePrefix
			baseDelay = queue.RegisterBaseDelay
		case queue.QueueTaskResetPassword:
			queuePrefix = queue.ResetPasswordQueuePrefix
			baseDelay = queue.ResetPasswordBaseDelay
		case queue.QueueTaskEmail:
			queuePrefix = queue.EmailQueuePrefix
			baseDelay = queue.EmailBaseDelay
		default:
			errorMessage := fmt.Sprintf("unsupported queue task: %s", queueEntry.Task)
			logger.ErrorCtx(jobContext, "Queue has unsupported task", "id", queueEntry.Id, "task", queueEntry.Task)
			if statusError := queueRepositoryInstance.SetErrorStatus(jobContext, queueEntry.Id, errorMessage, nil); statusError != nil {
				logger.ErrorCtx(jobContext, "Failed to set error status for queue", "id", queueEntry.Id, "error", statusError)
			}
			continue
		}

		newGeneration, resetError := queueRepositoryInstance.SetQueuedStatus(jobContext, queueEntry.Id)
		if resetError != nil {
			logger.ErrorCtx(jobContext, "Failed to increment generation for queue", "id", queueEntry.Id, "error", resetError)
			continue
		}

		if publishError := queuePublisherInstance.PublishEvent(jobContext, queueEntry.Id, newGeneration, queuePrefix, baseDelay); publishError != nil {
			logger.ErrorCtx(jobContext, "Failed to publish retry event for queue", "id", queueEntry.Id, "error", publishError)
			if statusError := queueRepositoryInstance.SetErrorStatus(jobContext, queueEntry.Id, publishError.Error(), nil); statusError != nil {
				logger.ErrorCtx(jobContext, "Failed to set error status for queue", "id", queueEntry.Id, "error", statusError)
			}
			continue
		}

		logger.InfoCtx(jobContext, "Successfully republished queue", "id", queueEntry.Id, "task", queueEntry.Task)
	}

	return nil
}
