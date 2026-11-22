package repository

import (
	"backend/internal/helper"
	"backend/internal/settings"
	"backend/internal/user"
	"context"
)

func (r *UserRepository) GetByNormalizedEmail(ctx context.Context, email string) (user.User, error) {
	var u user.User
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, email, normalized_email, password, registered_at, confirmed_at
        FROM users WHERE normalized_email = ?`, helper.NormalizeEmail(email))

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.NormalizedEmail, &u.Password, &u.RegisteredAt, &u.ConfirmedAt)
	return u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	var u user.User
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, email, normalized_email, password, registered_at, confirmed_at
        FROM users WHERE email = ?`, email)

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.NormalizedEmail, &u.Password, &u.RegisteredAt, &u.ConfirmedAt)
	return u, err
}

func (r *UserRepository) ExistsByName(ctx context.Context, username string) (bool, error) {
	exists := false
	row := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE name = ? LIMIT 1)`,
		username,
	)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) Create(ctx context.Context, user user.User) (uint, error) {
	result, err := r.db.ExecContext(ctx, `
        INSERT INTO users (name, email, normalized_email, password, registered_at, confirmed_at)
        VALUES (?, ?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.NormalizedEmail, user.Password, user.RegisteredAt, user.ConfirmedAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (r *UserRepository) GetById(ctx context.Context, id uint) (user.User, error) {
	var u user.User
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, email, normalized_email, password, registered_at, confirmed_at
        FROM users WHERE id = ?`, id)

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.NormalizedEmail, &u.Password, &u.RegisteredAt, &u.ConfirmedAt)
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetDataById(ctx context.Context, id uint) (user.User, settings.UserSettingsData, error) {
	var u user.User
	var us settings.UserSettings
	langCode := ""
	row := r.db.QueryRowContext(ctx, `
        SELECT u.id, u.name, email, normalized_email, registered_at, confirmed_at, user_flags, app_flags, app_opt_1, app_opt_2, app_opt_3, l.code
        FROM users AS u INNER JOIN user_settings AS us ON u.id = us.user_id 
		LEFT JOIN languages AS l ON us.lang_id = l.id
		WHERE u.id = ?`, id)

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.NormalizedEmail, &u.RegisteredAt, &u.ConfirmedAt, &us.UserFlags, &us.AppFlags, &us.AppOpt1, &us.AppOpt2, &us.AppOpt3, &langCode)

	if err != nil {
		return u, settings.UserSettingsData{}, err
	}
	settingsData := settings.UserSettingsData{
		Language:  &langCode,
		UserFlags: us.GetUserFlags(),
		AppFlags:  us.GetAppFlags(),
		AppOpt1:   &us.AppOpt1,
		AppOpt2:   &us.AppOpt2,
		AppOpt3:   &us.AppOpt3,
	}
	return u, settingsData, err
}

func (r *UserRepository) ChangeEmail(ctx context.Context, userId uint, newEmail string, normalizedEmail string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET email = ?, normalized_email = ? WHERE id = ?`,
		newEmail, normalizedEmail, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdatePasswordById(ctx context.Context, userId uint, newPassword string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET password = ? WHERE id = ?`,
		newPassword, userId)
	if err != nil {
		return err
	}
	return nil
}
