package publisher

import (
	"backend/internal/payload"
	"backend/internal/queue"
	"context"
	"encoding/json"
)

func (s *QueuePublisher) PublishConfirmationEmailTask(ctx context.Context, token string) error {
	confirmationEmailPayload := payload.ConfirmationEmailPayload{
		Token: token,
	}

	confirmationEmailPayloadJson, err := json.Marshal(confirmationEmailPayload)
	if err != nil {
		return err
	}
	emailQueuePayload := payload.EmailQueuePayload{
		Type: queue.ConfirmationEmailType,
		Data: confirmationEmailPayloadJson,
	}

	return PublishTask(ctx, s, emailQueuePayload, queue.QueueTaskEmail, queue.EmailQueuePrefix, queue.EmailBaseDelay, queue.EmailMaxRetries)
}
