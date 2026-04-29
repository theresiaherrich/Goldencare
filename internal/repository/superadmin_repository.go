package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type SuperadminRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.Superadmin, error)
	Create(ctx context.Context, superadmin *models.Superadmin) error
}

type superadminRepository struct {
	db *sqlx.DB
}

func newSuperadminRepository(db *sqlx.DB) SuperadminRepository {
	return &superadminRepository{db: db}
}

func (r *superadminRepository) GetByEmail(ctx context.Context, email string) (*models.Superadmin, error) {
	superadmin := &models.Superadmin{}
	err := r.db.GetContext(ctx, superadmin, `
		SELECT id, email, password, name, created_at FROM superadmin WHERE email = $1
	`, email)
	if err != nil {
		return nil, err
	}
	return superadmin, nil
}

func (r *superadminRepository) Create(ctx context.Context, superadmin *models.Superadmin) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO superadmin (id, email, password, name, created_at) VALUES ($1, $2, $3, $4, $5)
	`, superadmin.ID, superadmin.Email, superadmin.Password, superadmin.Name, superadmin.CreatedAt)
	return err
}

