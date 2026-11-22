package consumer

import (
	"backend/internal/apierrors"
	confirmationService "backend/internal/confirmation/service"
	"backend/internal/helper"
	"backend/internal/payload"
	"backend/internal/queue"
	"backend/pkg/apperror"
	"backend/pkg/validation"
	"context"
)

func (s *QueueConsumer) ConsumeResetPasswordTask(ctx context.Context, event queue.QueueEvent, resetPasswordPayload payload.ResetPasswordQueuePayload) error {
	err := validation.IsEmailValid(resetPasswordPayload.Email)
	if err != nil {
		return apperror.Wrap(apierrors.ErrUnrecoverableError, err.Error())
	}
	email := helper.CleanEmail(resetPasswordPayload.Email)
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.Id == 0 {
		return apperror.Wrap(apierrors.ErrUnrecoverableError, "User not found")
	}
	confService := confirmationService.NewService(ctx)
	resetToken, err := confService.CreateResetPasswordToken(ctx, user.Id)
	if err != nil {
		return err
	}

	return s.queuePublisher.PublishConfirmationEmailTask(ctx, resetToken)
}

func (s *QueueConsumer) StartResetPasswordConsumer(ctx context.Context) error {
	return StartConsumer(
		ctx,
		s,
		queue.ResetPasswordQueuePrefix,
		queue.ResetPasswordBaseDelay,
		queue.QueueTaskResetPassword,
		s.ConsumeResetPasswordTask,
	)
}
