package repository

import (
	"backend/internal/models"
	"context"
)

type UserSettingsRepository struct {
	db DBExecutor
}

func NewUserSettingsRepository(db DBExecutor) *UserSettingsRepository {
	return &UserSettingsRepository{db: db}
}
func (r *UserSettingsRepository) Create(ctx context.Context, userId uint, languageId uint8) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_settings (user_id, lang_id)
		VALUES (?, ?)`,
		userId, languageId)
	return err
}

func (r *UserSettingsRepository) GetByUserId(ctx context.Context, userId uint) (models.UserSettings, error) {
	var us models.UserSettings
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, lang_id, user_flags, app_flags, app_opt_1, app_opt_2, app_opt_3, updated_at FROM user_settings WHERE user_id = ?`, userId).Scan(&us.Id, &us.UserId, &us.LangId, &us.UserFlags, &us.AppFlags, &us.AppOpt1, &us.AppOpt2, &us.AppOpt3, &us.UpdatedAt)
	return us, err
}

func (r *UserSettingsRepository) Update(ctx context.Context, us models.UserSettings) error {
	if us.Id == 0 {
		return nil
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_settings SET user_id=?, lang_id=?, user_flags=?, app_flags=?, app_opt_1=?, app_opt_2=?, app_opt_3=?, updated_at=? WHERE id=?`,
			us.UserId,
			us.LangId,
			us.UserFlags,
			us.AppFlags,
			us.AppOpt1,
			us.AppOpt2,
			us.AppOpt3,
			us.UpdatedAt,
			us.Id,
	)
	return err
}