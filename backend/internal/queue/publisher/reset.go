package publisher

import (
	"backend/internal/payload"
	"backend/internal/queue"
	"backend/pkg/apperror"
	"context"
	"encoding/json"
)

func (s *QueuePublisher) PublishResetPasswordTask(ctx context.Context, data payload.ResetPasswordQueuePayload) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	queueId, err := s.queueRepository.Create(ctx, string(queue.QueueTaskResetPassword), jsonData, queue.ResetPasswordMaxRetries)
	if err != nil {
		return err
	}

	err = s.PublishEvent(ctx, queueId, queue.InitialGeneration, queue.ResetPasswordQueuePrefix, queue.ResetPasswordBaseDelay)
	if err != nil {
		if dbErr := s.queueRepository.SetErrorStatus(ctx, uint(queueId), err.Error(), nil); dbErr != nil {
			return apperror.Wrap(err, "Failed to set queue as failed", "dbErr", dbErr)
		}
		return err
	}

	if _, dbErr := s.queueRepository.SetQueuedStatus(ctx, uint(queueId)); dbErr != nil {
		return apperror.Wrap(dbErr, "Failed to set queue as queued")
	}

	return nil
}
