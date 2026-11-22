package publisher

import (
	"backend/internal/payload"
	"backend/internal/queue"
	"context"
)

func (s *QueuePublisher) PublishRegisterUserTask(ctx context.Context, data payload.RegisterQueuePayload) error {
	return PublishTask(ctx, s, data, queue.QueueTaskRegisterUser, queue.RegisterQueuePrefix, queue.RegisterBaseDelay, queue.RegisterMaxRetries)
}
