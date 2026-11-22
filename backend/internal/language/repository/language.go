package repository

import (
	"backend/internal/database"
	"backend/internal/language"
	"context"
	"database/sql"
)

type LanguageRepository struct {
	db database.DBExecutor
}

func NewRepository(db database.DBExecutor) *LanguageRepository {
	return &LanguageRepository{db: db}
}

func (r *LanguageRepository) Get(ctx context.Context) ([]language.Language, error) {
	list := []language.Language{}
	rows, err := r.db.QueryContext(ctx, `SELECT id, code, name FROM languages`)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var lang language.Language
		if err := rows.Scan(&lang.Id, &lang.Code, &lang.Name); err != nil {
			return list, err
		}
		list = append(list, lang)
	}
	return list, nil
}
func (r *LanguageRepository) GetById(ctx context.Context, id uint8) (language.Language, error) {
	var lang language.Language
	err := r.db.QueryRowContext(ctx, `SELECT id, code, name FROM languages WHERE id = ?`, id).Scan(&lang.Id, &lang.Code, &lang.Name)
	if err == sql.ErrNoRows {
		return lang, nil
	}
	return lang, err
}
func (r *LanguageRepository) GetLangByCode(ctx context.Context, code string) (language.Language, error) {
	var lang language.Language
	err := r.db.QueryRowContext(ctx, `SELECT id, code, name FROM languages WHERE code = ?`, code).Scan(&lang.Id, &lang.Code, &lang.Name)
	if err != nil {
		return lang, err
	}
	return lang, nil
}
