package service

import (
	"backend/internal/contexthelper"
	"backend/internal/language"
	"backend/internal/language/repository"
	"context"
)

type LanguageService struct {
	repo *repository.LanguageRepository
}

func NewService(ctx context.Context) *LanguageService {
	db := contexthelper.GetDb(ctx)
	repo := repository.NewRepository(db)
	return &LanguageService{repo: repo}
}

func (s *LanguageService) Get(ctx context.Context) ([]language.Language, error) {
	return s.repo.Get(ctx)
}

func (s *LanguageService) GetById(ctx context.Context, id uint8) (language.Language, error) {
	return s.repo.GetById(ctx, id)
}

func (s *LanguageService) GetLangByCode(ctx context.Context, code string) (language.Language, error) {
	return s.repo.GetLangByCode(ctx, code)
}
