package service

import (
	"backend/internal/models"
	"backend/internal/repository"

	"github.com/pkg/errors"
	"context"
)

type LanguageService struct {
	repo *repository.LanguageRepository
}

func NewLanguageService(repo *repository.LanguageRepository) *LanguageService {
	return &LanguageService{repo: repo}
}

func (s *LanguageService) GetLanguages(ctx context.Context) ([]models.Language, error) {
	languages, err := s.repo.Get(ctx)
	if err != nil {
		return []models.Language{}, errors.Wrap(err, "languages not found")
	}
	if len(languages) == 0 {
		return []models.Language{}, errors.New("languages not found")
	}
	return languages, nil
}
