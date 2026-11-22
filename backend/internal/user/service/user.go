package service

import (
	"backend/internal/apierrors"
	"backend/internal/cache"
	confirmationService "backend/internal/confirmation/service"
	"backend/internal/helper"
	"backend/internal/payload"
	queuePublisher "backend/internal/queue/publisher"
	"backend/internal/user"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) EnqueueRegisterUser(ctx context.Context, userName, email, password, langCode string) error {

	lock := cache.NewLock(s.redisClient, "register_user_"+userName)

	// Spróbuj założyć lock na max 5 sekund
	acquired, err := lock.AcquireWithRetry(ctx, 5*time.Second, cache.DefaultTimeout, cache.DefaultInterval)
	if err != nil {
		return err
	}
	if !acquired {
		return apierrors.NewRegisterUserNameTakenError("Registration for username " + userName + " is already in progress")
	}

	defer func() {
		if err := lock.Release(ctx); err != nil {
			logger.ErrorCtx(ctx, "Failed to release lock", "error", err)
		}
	}()

	// Check if user name exists
	exists, err := s.userRepository.ExistsByName(ctx, userName)
	if err != nil {
		return err
	}
	if exists {
		logger.WarnCtx(ctx, "Username already taken", "username", userName)
		return apierrors.NewRegisterUserNameTakenError("Username already taken")
	}
	exists, err = s.queueRepository.ExistsRegisterQueueByUserName(ctx, userName)
	if err != nil {
		return err
	}

	if exists {
		logger.WarnCtx(ctx, "Queue task already exists for user", "username", userName)
		return apierrors.NewRegisterUserNameTakenError("Username already taken")
	}
	exists, err = s.confirmationTokenRepository.ExistsRegisterTokenByName(ctx, userName)
	if err != nil {
		return err
	}

	if exists {
		logger.WarnCtx(ctx, "Register token already exists for user", "username", userName)
		return apierrors.NewRegisterUserNameTakenError("Username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	payload := payload.RegisterQueuePayload{
		User: user.User{
			Name:         userName,
			Email:        email,
			RegisteredAt: time.Now(),
		},
		LanguageCode:    langCode,
		NewPassword:     string(hashedPassword),
		NormalizedEmail: helper.NormalizeEmail(email),
	}

	publisher := queuePublisher.NewPublisher(ctx)
	err = publisher.PublishRegisterUserTask(ctx, payload)
	if err != nil {
		return apperror.Wrap(err, "Failed to enqueue register user task")
	}

	return nil
}

func (s *UserService) GetUserResponseData(ctx context.Context, userId uint) (user.UserResponseData, error) {
	userData, settings, err := s.userRepository.GetDataById(ctx, userId)
	if err != nil || userData.Id == 0 {
		return user.UserResponseData{}, apperror.Wrap(err, "user not found")
	}
	userResponseData := user.UserResponseData{
		Id:           userData.Id,
		Name:         userData.Name,
		Email:        userData.Email,
		RegisteredAt: userData.RegisteredAt.Format("2006-01-02 15:04:05"),
		ConfirmedAt:  userData.ConfirmedAt.Format("2006-01-02 15:04:05"),
		Settings:     settings,
	}

	return userResponseData, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (user.UserResponseData, error) {
	userData, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil || userData.Id == 0 {
		return user.UserResponseData{}, apperror.Wrap(err, "user not found")
	}

	// Sprawdzenie hasła bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password))
	if err != nil {
		return user.UserResponseData{}, apperror.Wrap(err, "invalid password")
	}
	userFullData, settings, err := s.userRepository.GetDataById(ctx, userData.Id)
	if err != nil {
		return user.UserResponseData{}, apperror.Wrap(err, "user not found")
	}
	userResponseData := user.UserResponseData{
		Id:           userFullData.Id,
		Name:         userData.Name,
		Email:        userData.Email,
		RegisteredAt: userData.RegisteredAt.Format("2006-01-02 15:04:05"),
		ConfirmedAt:  userData.ConfirmedAt.Format("2006-01-02 15:04:05"),
		Settings:     settings,
	}

	return userResponseData, nil
}
func (s *UserService) EnqueueResetPassword(ctx context.Context, email string) error {

	payload := payload.ResetPasswordQueuePayload{
		Email: email,
	}

	publisher := queuePublisher.NewPublisher(ctx)
	err := publisher.PublishResetPasswordTask(ctx, payload)
	if err != nil {
		return apperror.Wrap(err, "Failed to enqueue reset password task")
	}

	return nil
}

func (s *UserService) EnqueueChangeEmail(ctx context.Context, userId uint, email string) error {

	confService := confirmationService.NewService(ctx)
	token, err := confService.CreateEmailChangeToken(ctx, payload.ChangeEmailQueuePayload{
		UserId:   userId,
		NewEmail: email,
	})
	if err != nil {
		return err
	}
	publisher := queuePublisher.NewPublisher(ctx)
	return publisher.PublishConfirmationEmailTask(ctx, token)
}

func (s *UserService) GetByNormalizedEmail(ctx context.Context, email string) (user.User, error) {
	return s.userRepository.GetByNormalizedEmail(ctx, email)
}

func (s *UserService) IsEmailAvailableForUser(ctx context.Context, userId uint, email string) (bool, error) {
	userData, err := s.userRepository.GetByNormalizedEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if userData.Id == userId {
		return true, nil
	}

	if userData.Id != 0 {
		return false, nil
	}
	return s.confirmationTokenRepository.IsNewEmailAvailable(ctx, email)
}
