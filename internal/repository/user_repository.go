package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *sqlx.DB
}

func newUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password, role, is_verified, panti_id, kode_dipakai, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, user.ID, user.Name, user.Email, user.Password, user.Role,
		user.IsVerified, user.PantiID, user.KodeDipakai, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.GetContext(ctx, user, `
		SELECT id, name, email, password, role, is_verified, panti_id, kode_dipakai, created_at, updated_at
		FROM users WHERE email = $1
	`, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.GetContext(ctx, user, `
		SELECT id, name, email, password, role, is_verified, panti_id, kode_dipakai, created_at, updated_at
		FROM users WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE users SET name = $1, password = $2, role = $3, is_verified = $4, panti_id = $5, updated_at = $6
		WHERE id = $7
	`, user.Name, user.Password, user.Role, user.IsVerified, user.PantiID, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}