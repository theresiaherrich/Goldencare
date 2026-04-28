package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type KodeUndanganRepository interface {
	Create(ctx context.Context, kode *models.KodeUndangan) error
	GetByKode(ctx context.Context, kode string) (*models.KodeUndangan, error)
	GetByID(ctx context.Context, id string) (*models.KodeUndangan, error)
	GetByPantiID(ctx context.Context, pantiID string) ([]models.KodeUndangan, error)
	Update(ctx context.Context, kode *models.KodeUndangan) error
	Delete(ctx context.Context, id string) error
}

type kodeUndanganRepository struct {
	db *sqlx.DB
}

func newKodeUndanganRepository(db *sqlx.DB) KodeUndanganRepository {
	return &kodeUndanganRepository{db: db}
}

func (r *kodeUndanganRepository) Create(ctx context.Context, kode *models.KodeUndangan) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO kode_undangan
		(id, kode, untuk_role, dibuat_oleh, panti_id, tipe, dipakai_count, maks_pakai, expired_at, is_aktif, catatan, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, kode.ID, kode.Kode, kode.UntukRole, kode.DibuatOleh, kode.PantiID, kode.Tipe,
		kode.DipakaiCount, kode.MaksPakai, kode.ExpiredAt, kode.IsAktif, kode.Catatan, kode.CreatedAt)
	return err
}

func (r *kodeUndanganRepository) GetByKode(ctx context.Context, kode string) (*models.KodeUndangan, error) {
	result := &models.KodeUndangan{}
	err := r.db.GetContext(ctx, result, `
		SELECT id, kode, untuk_role, dibuat_oleh, panti_id, tipe, dipakai_count, maks_pakai, expired_at, is_aktif, catatan, created_at
		FROM kode_undangan WHERE kode = $1
	`, kode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *kodeUndanganRepository) GetByID(ctx context.Context, id string) (*models.KodeUndangan, error) {
	result := &models.KodeUndangan{}
	err := r.db.GetContext(ctx, result, `
		SELECT id, kode, untuk_role, dibuat_oleh, panti_id, tipe, dipakai_count, maks_pakai, expired_at, is_aktif, catatan, created_at
		FROM kode_undangan WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *kodeUndanganRepository) GetByPantiID(ctx context.Context, pantiID string) ([]models.KodeUndangan, error) {
	var results []models.KodeUndangan
	err := r.db.SelectContext(ctx, &results, `
		SELECT id, kode, untuk_role, dibuat_oleh, panti_id, tipe, dipakai_count, maks_pakai, expired_at, is_aktif, catatan, created_at
		FROM kode_undangan WHERE panti_id = $1 ORDER BY created_at DESC
	`, pantiID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *kodeUndanganRepository) Update(ctx context.Context, kode *models.KodeUndangan) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE kode_undangan SET kode = $1, untuk_role = $2, dipakai_count = $3, is_aktif = $4 WHERE id = $5
	`, kode.Kode, kode.UntukRole, kode.DipakaiCount, kode.IsAktif, kode.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("kode undangan not found")
	}
	return nil
}

func (r *kodeUndanganRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM kode_undangan WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("kode undangan not found")
	}
	return nil
}