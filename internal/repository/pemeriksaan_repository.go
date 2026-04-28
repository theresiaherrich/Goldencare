package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type PemeriksaanRepository interface {
	Create(ctx context.Context, pemeriksaan *models.PemeriksaanKesehatan) error
	GetByID(ctx context.Context, id string) (*models.PemeriksaanKesehatan, error)
	GetByLansia(ctx context.Context, lansiaID string) ([]models.PemeriksaanKesehatan, error)
	GetLatestByLansia(ctx context.Context, lansiaID string) (*models.PemeriksaanKesehatan, error)
}

type pemeriksaanRepository struct {
	db *sqlx.DB
}

func newPemeriksaanRepository(db *sqlx.DB) PemeriksaanRepository {
	return &pemeriksaanRepository{db: db}
}

func (r *pemeriksaanRepository) Create(ctx context.Context, p *models.PemeriksaanKesehatan) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO pemeriksaan_kesehatan
		(id, lansia_id, pengurus_id, tekanan_darah, detak_jantung, gula_darah, suhu_tubuh, berat_badan, keluhan, status, rekomendasi, status_darurat, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, p.ID, p.LansiaID, p.PengurusID, p.TekananDarah, p.DetakJantung, p.GulaDarah,
		p.SuhuTubuh, p.BeratBadan, p.Keluhan, p.Status, p.Rekomendasi, p.StatusDarurat, p.CreatedAt)
	return err
}

func (r *pemeriksaanRepository) GetByID(ctx context.Context, id string) (*models.PemeriksaanKesehatan, error) {
	var p models.PemeriksaanKesehatan
	err := r.db.GetContext(ctx, &p, `
		SELECT id, lansia_id, pengurus_id, tekanan_darah, detak_jantung, gula_darah, suhu_tubuh, berat_badan, keluhan, status, rekomendasi, status_darurat, created_at
		FROM pemeriksaan_kesehatan WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pemeriksaanRepository) GetByLansia(ctx context.Context, lansiaID string) ([]models.PemeriksaanKesehatan, error) {
	var results []models.PemeriksaanKesehatan
	err := r.db.SelectContext(ctx, &results, `
		SELECT id, lansia_id, pengurus_id, tekanan_darah, detak_jantung, gula_darah, suhu_tubuh, berat_badan, keluhan, status, rekomendasi, status_darurat, created_at
		FROM pemeriksaan_kesehatan WHERE lansia_id = $1 ORDER BY created_at DESC LIMIT 20
	`, lansiaID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *pemeriksaanRepository) GetLatestByLansia(ctx context.Context, lansiaID string) (*models.PemeriksaanKesehatan, error) {
	var p models.PemeriksaanKesehatan
	err := r.db.GetContext(ctx, &p, `
		SELECT id, lansia_id, pengurus_id, tekanan_darah, detak_jantung, gula_darah, suhu_tubuh, berat_badan, keluhan, status, rekomendasi, status_darurat, created_at
		FROM pemeriksaan_kesehatan WHERE lansia_id = $1 ORDER BY created_at DESC LIMIT 1
	`, lansiaID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}