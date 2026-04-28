package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type GaleriRepository interface {
	Create(ctx context.Context, galeri *models.GaleriFisik) error
	GetByID(ctx context.Context, id string) (*models.GaleriFisik, error)
	GetByLansia(ctx context.Context, lansiaID string) ([]models.GaleriFisik, error)
}

type galeriRepository struct {
	db *sqlx.DB
}

func newGaleriRepository(db *sqlx.DB) GaleriRepository {
	return &galeriRepository{db: db}
}

func (r *galeriRepository) Create(ctx context.Context, galeri *models.GaleriFisik) error {
	query := `
        INSERT INTO galeri_fisik
        (id, lansia_id, pengurus_id, foto_url, lokasi_luka, deskripsi, analisis_ai, tingkat_darurat, prediksi_penyakit, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err := r.db.ExecContext(ctx, query,
		galeri.ID, galeri.LansiaID, galeri.PengurusID, galeri.FotoURL, galeri.LokasiLuka,
		galeri.Deskripsi, galeri.AnalisisAI, galeri.TingkatDarurat, galeri.PrediksiPenyakit, galeri.CreatedAt,
	)
	return err
}

func (r *galeriRepository) GetByID(ctx context.Context, id string) (*models.GaleriFisik, error) {
	var galeri models.GaleriFisik
	query := `SELECT id, lansia_id, pengurus_id, foto_url, lokasi_luka, deskripsi, analisis_ai, tingkat_darurat, prediksi_penyakit, created_at
              FROM galeri_fisik WHERE id = $1`
	if err := r.db.GetContext(ctx, &galeri, query, id); err != nil {
		return nil, err
	}
	return &galeri, nil
}

func (r *galeriRepository) GetByLansia(ctx context.Context, lansiaID string) ([]models.GaleriFisik, error) {
	var galeris []models.GaleriFisik
	query := `SELECT id, lansia_id, pengurus_id, foto_url, lokasi_luka, deskripsi, analisis_ai, tingkat_darurat, prediksi_penyakit, created_at
              FROM galeri_fisik WHERE lansia_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &galeris, query, lansiaID); err != nil {
		return nil, err
	}
	return galeris, nil
}
