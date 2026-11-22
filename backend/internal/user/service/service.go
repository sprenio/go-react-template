package service

import (
	confirmationRepository "backend/internal/confirmation/repository"
	"backend/internal/contexthelper"
	languageRepository "backend/internal/language/repository"
	queueRepository "backend/internal/queue/repository"
	userRepository "backend/internal/user/repository"
	"context"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	confirmationTokenRepository *confirmationRepository.ConfirmationTokenRepository
	queueRepository             *queueRepository.QueueRepository
	languageRepository          *languageRepository.LanguageRepository
	userRepository              *userRepository.UserRepository
	redisClient                 *redis.Client
}

func NewService(ctx context.Context) *UserService {
	db := contexthelper.GetDb(ctx)
	confirmationTokenRepository := confirmationRepository.NewRepository(db)
	languageRepository := languageRepository.NewRepository(db)
	userRepository := userRepository.NewRepository(db)
	redisClient := contexthelper.GetRedisClient(ctx)
	queueRepository := queueRepository.NewQueueRepository(db)
	return &UserService{confirmationTokenRepository: confirmationTokenRepository, languageRepository: languageRepository, userRepository: userRepository, redisClient: redisClient, queueRepository: queueRepository}
}
