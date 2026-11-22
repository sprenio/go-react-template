package consumer

import (
	"backend/internal/apierrors"
	emailService "backend/internal/email/service"
	"backend/internal/payload"
	"backend/internal/queue"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"encoding/json"
)

func (s *QueueConsumer) ConsumeEmailTask(ctx context.Context, event queue.QueueEvent, emailPayload payload.EmailQueuePayload) error {
	logger.InfoCtx(ctx, "📧 starting handling email task", "type", emailPayload.Type, "queueId", event.QueueId, "generation", event.Generation, "retry", event.Retry)
	var err error
	emailLogId := uint(0)
	if emailLogId, err = s.emailRepository.Create(ctx, event.QueueId, event.Generation, event.Retry); err != nil {
		logger.ErrorCtx(ctx, "❌ failed to create email log", "error", err)
		return err
	}

	defer func() {
		logger.InfoCtx(ctx, "defer")
		if err == nil {
			logger.InfoCtx(ctx, "err is nil")
			return
		}
		logger.ErrorCtx(ctx, "❌ failed to handle email task", "error", err)
		if _, err = s.emailRepository.SetErrorStatus(ctx, emailLogId, err.Error()); err != nil {
			logger.ErrorCtx(ctx, "❌ failed to set email log error status", "error", err)
		}
	}()
	switch emailPayload.Type {
	case queue.ConfirmationEmailType:
		var data payload.ConfirmationEmailPayload
		err = json.Unmarshal(emailPayload.Data, &data)
		if err != nil {
			err = apperror.Wrap(err, "failed to unmarshal confirmation email payload")
		} else if data.Token == "" {
			err = apperror.Wrap(apierrors.ErrUnrecoverableError, "empty register token")
		} else {
			emailSvc := emailService.NewService(ctx)
			err = emailSvc.SendConfirmationEmail(ctx, data.Token)
		}
	default:
		err = apperror.Wrap(apierrors.ErrUnrecoverableError, "unknown email task type", "type", emailPayload.Type)
	}
	if err != nil {
		logger.InfoCtx(ctx, "return errorv", "error", err)
		return err
	}
	logger.InfoCtx(ctx, "✅ Email task handled successfully")
	if _, err = s.emailRepository.SetSentStatus(ctx, emailLogId); err != nil {
		logger.ErrorCtx(ctx, "❌ failed to set email log sent status", "error", err)
	}
	return nil
}

func (s *QueueConsumer) StartEmailConsumer(ctx context.Context) error {
	return StartConsumer(
		ctx,
		s,
		queue.EmailQueuePrefix,
		queue.EmailBaseDelay,
		queue.QueueTaskEmail,
		s.ConsumeEmailTask,
	)
}
