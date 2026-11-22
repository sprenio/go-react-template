package consumer

import (
	confirmationService "backend/internal/confirmation/service"
	"backend/internal/payload"
	"backend/internal/queue"
	userService "backend/internal/user/service"
	"context"

	"errors"
)

func (s *QueueConsumer) ConsumeRegisterUserTask(ctx context.Context, event queue.QueueEvent, registerPayload payload.RegisterQueuePayload) error {

	usrService := userService.NewService(ctx)
	emailAvailable, err := usrService.IsEmailAvailableForUser(ctx, 0, registerPayload.NormalizedEmail)
	if err != nil {
		return err
	}
	if !emailAvailable {
		return errors.New("email is not available")
	}
	confService := confirmationService.NewService(ctx)
	registerToken, err := confService.CreateRegisterToken(ctx, registerPayload)
	if err != nil {
		return err
	}

	return s.queuePublisher.PublishConfirmationEmailTask(ctx, registerToken)
}

func (s *QueueConsumer) StartRegisterConsumer(ctx context.Context) error {
	return StartConsumer(
		ctx,
		s,
		queue.RegisterQueuePrefix,
		queue.RegisterBaseDelay,
		queue.QueueTaskRegisterUser,
		s.ConsumeRegisterUserTask,
	)
}
