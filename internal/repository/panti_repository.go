package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type PantiRepository interface {
	Create(ctx context.Context, panti *models.Panti) error
	GetByID(ctx context.Context, id string) (*models.Panti, error)
	GetAll(ctx context.Context) ([]models.Panti, error)
	Update(ctx context.Context, panti *models.Panti) error
}

type pantiRepository struct {
	db *sqlx.DB
}

func newPantiRepository(db *sqlx.DB) PantiRepository {
	return &pantiRepository{db: db}
}

func (r *pantiRepository) Create(ctx context.Context, panti *models.Panti) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO panti (id, nama, alamat, telepon, kode_undangan, pengelola_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, panti.ID, panti.Nama, panti.Alamat, panti.Telepon, panti.KodeUndangan, panti.PengelolaID, panti.CreatedAt, panti.UpdatedAt)
	return err
}

func (r *pantiRepository) GetByID(ctx context.Context, id string) (*models.Panti, error) {
	panti := &models.Panti{}
	err := r.db.GetContext(ctx, panti, `
		SELECT id, nama, alamat, telepon, kode_undangan, pengelola_id, created_at, updated_at
		FROM panti WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return panti, nil
}

func (r *pantiRepository) GetAll(ctx context.Context) ([]models.Panti, error) {
	var pantis []models.Panti
	err := r.db.SelectContext(ctx, &pantis, `
		SELECT id, nama, alamat, telepon, kode_undangan, pengelola_id, created_at, updated_at
		FROM panti ORDER BY nama ASC
	`)
	if err != nil {
		return nil, err
	}
	return pantis, nil
}

func (r *pantiRepository) Update(ctx context.Context, panti *models.Panti) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE panti SET nama = $1, alamat = $2, telepon = $3, updated_at = $4 WHERE id = $5
	`, panti.Nama, panti.Alamat, panti.Telepon, panti.UpdatedAt, panti.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("panti not found")
	}
	return nil
}