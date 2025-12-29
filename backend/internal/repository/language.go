package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

type LanguageRepository struct {
	db DBExecutor
}

func NewLanguageRepository(db DBExecutor) *LanguageRepository {
	return &LanguageRepository{db: db}
}
func (r *LanguageRepository) Get(ctx context.Context) ([]models.Language, error) {
	list := []models.Language{}
	rows, err := r.db.QueryContext(ctx, `SELECT id, code, name FROM languages`)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var lang models.Language
		if err := rows.Scan(&lang.Id, &lang.Code, &lang.Name); err != nil {
			return list, err
		}
		list = append(list, lang)
	}
	return list, nil
}
func (r *LanguageRepository) GetById(ctx context.Context, id uint8) (models.Language, error) {
	var lang models.Language
	err := r.db.QueryRowContext(ctx, `SELECT id, code, i18n_code, name FROM languages WHERE id = ?`, id).Scan(&lang.Id, &lang.Code, &lang.I18nCode, &lang.Name)
	if err == sql.ErrNoRows {
		return lang, nil
	}
	return lang, err
}
func (r *LanguageRepository) GetLangByCode(ctx context.Context, code string) (models.Language, error) {
	var lang models.Language
	err := r.db.QueryRowContext(ctx, `SELECT id, code, i18n_code, name FROM languages WHERE code = ?`, code).Scan(&lang.Id, &lang.Code, &lang.I18nCode, &lang.Name)
	if err != nil {
		return lang, err
	}
	return lang, nil
}