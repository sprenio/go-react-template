package repository

import (
	"backend/internal/models"
	"backend/internal/response"
	"context"
	"database/sql"
	"strings"
)

type UserRepository struct {
	db DBExecutor
}

func NewUserRepository(db DBExecutor) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, email, password, registered_at, confirmed_at
        FROM users WHERE email = ?`, strings.ToLower(email))

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.RegisteredAt, &u.ConfirmedAt)
	return u, err
}

func (r *UserRepository) ExistsByEmailOrName(ctx context.Context, email, username string) (error, bool) {
	exists := false
	row := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email =?
	  UNION SELECT id FROM users WHERE name = ?`,
		email,
		username,
	)
	err := row.Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return err, false
	}
	return nil, exists
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (error, bool) {
	exists := false
	row := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email =?`,
		email,
	)
	err := row.Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return err, false
	}
	return nil, exists
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (uint, error) {
	result, err := r.db.ExecContext(ctx, `
        INSERT INTO users (name, email, password, registered_at, confirmed_at)
        VALUES (?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Password, user.RegisteredAt, user.ConfirmedAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (r *UserRepository) GetById(ctx context.Context, id uint) (models.User, error) {
	var u models.User
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, email, password, registered_at, confirmed_at
        FROM users WHERE id = ?`, id)

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.RegisteredAt, &u.ConfirmedAt)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetDataById(ctx context.Context, id uint) (response.UserResponseData, error) {
	var u models.User
	var us models.UserSettings
	langCode := ""
	row := r.db.QueryRowContext(ctx, `
        SELECT u.id, u.name, email, registered_at, confirmed_at, user_flags, app_flags, app_opt_1, app_opt_2, app_opt_3, l.code
        FROM users AS u INNER JOIN user_settings AS us ON u.id = us.user_id 
		LEFT JOIN languages AS l ON us.lang_id = l.id
		WHERE u.id = ?`, id)

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.RegisteredAt, &u.ConfirmedAt, &us.UserFlags, &us.AppFlags, &us.AppOpt1, &us.AppOpt2, &us.AppOpt3, &langCode)

	if err != nil {
		return response.UserResponseData{}, err
	}
	userData := response.UserResponseData{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		RegisteredAt: u.RegisteredAt.Format("2006-01-02 15:04:05"),
		ConfirmedAt:  u.ConfirmedAt.Format("2006-01-02 15:04:05"),
		Settings: models.UserSettingsData{
			Language:  &langCode,
			UserFlags: us.GetUserFlags(),
			AppFlags:  us.GetAppFlags(),
			AppOpt1:   &us.AppOpt1,
			AppOpt2:   &us.AppOpt2,
			AppOpt3:   &us.AppOpt3,
		},
	}
	return userData, nil
}

func (r *UserRepository) ChangeEmail(ctx context.Context, userId uint, newEmail string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET email = ? WHERE id = ?`,
		newEmail, userId)
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
