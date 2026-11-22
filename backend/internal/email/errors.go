package email

import (
	"backend/internal/apierrors"
	"fmt"

	"errors"
)

var (
	ErrFailedToGetConfirmationToken = fmt.Errorf("%w: failed to get confirmation token", apierrors.ErrUnrecoverableError)
	ErrFailedToUnmarshalPayloadData = fmt.Errorf("%w: failed to unmarshal payload data", apierrors.ErrUnrecoverableError)
	ErrInvalidEmailData             = fmt.Errorf("%w: invalid email data", apierrors.ErrUnrecoverableError)
	ErrFailedToGetEmailSender       = fmt.Errorf("%w: failed to get email sender", apierrors.ErrUnrecoverableError)
	ErrFailedToSendEmail            = errors.New("failed to send email")
	ErrInvalidUserId                = fmt.Errorf("%w: invalid user", apierrors.ErrUnrecoverableError)
)
